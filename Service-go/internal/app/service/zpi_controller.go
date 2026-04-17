package zpi_controller

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	zpi_camera "github.com/Z-TAS-Solutions/ZScannerService/internal/pkg/zpi_cam"
	"github.com/Z-TAS-Solutions/ZScannerService/internal/pkg/zpi_indicator"
	zpi_trigger "github.com/Z-TAS-Solutions/ZScannerService/internal/pkg/zpi_tof"
	"github.com/Z-TAS-Solutions/ZScannerService/internal/pkg/zscanproto"
)

type ControllerServer struct {
	zscanproto.UnimplementedZPiControllerServer
	indicator       *zpi_indicator.ZLED
	indicatorStatus zscanproto.LEDStatus

	trigger             *zpi_trigger.ZToF
	triggerThreshold    uint16
	triggerStatus       zscanproto.ToFState
	triggerDeactivation context.CancelFunc

	camProcess *zpi_camera.CameraProcess
	camStatus  zscanproto.CamState
	camConfig  *zscanproto.CameraConfig
	camMutex   sync.Mutex

	eventStatus    bool
	tofEventStream zscanproto.ZPiController_ToFEventStreamServer
}

func NewControllerServer(indicator *zpi_indicator.ZLED, trigger *zpi_trigger.ZToF) *ControllerServer {
	return &ControllerServer{indicator: indicator, trigger: trigger}
}

var statusColors = map[zscanproto.LEDStatus][3]uint32{
	zscanproto.LEDStatus_VOID:    {0, 0, 0},
	zscanproto.LEDStatus_SUCCESS: {0, 255, 0},
	zscanproto.LEDStatus_PENDING: {0, 0, 255},
	zscanproto.LEDStatus_FAILED:  {255, 0, 0},
}

func (s *ControllerServer) SetLEDStatus(ctx context.Context, req *zscanproto.LEDStatusRequest) (*zscanproto.Status, error) {
	colors, ok := statusColors[req.Status]
	if !ok {
		return &zscanproto.Status{Success: false, Message: "unknown status"}, nil
	}

	err := s.indicator.Set(colors[0], colors[1], colors[2])
	if err != nil {
		return &zscanproto.Status{Success: false, Message: err.Error()}, nil
	}
	return &zscanproto.Status{Success: true, Message: fmt.Sprintf("LED set to %v", req.Status)}, nil
}

func (s *ControllerServer) GetLED(ctx context.Context, _ *zscanproto.Empty) (*zscanproto.LEDState, error) {
	r, g, b := s.indicator.Get()

	return &zscanproto.LEDState{Red: r, Green: g, Blue: b}, nil
}

func (s *ControllerServer) ToFEventStream(stream zscanproto.ZPiController_ToFEventStreamServer) error {
	ctx := stream.Context()
	fmt.Println("Client Registered to ToF Event Stream")
	s.tofEventStream = stream
	s.StartToFMonitor(ctx, stream)
	return nil
}

func (s *ControllerServer) ConfigureToF(ctx context.Context, req *zscanproto.ToFConfig) (*zscanproto.Status, error) {
	previousState := s.triggerStatus

	if previousState == zscanproto.ToFState_ToFActive {
		if s.triggerDeactivation != nil {
			s.triggerDeactivation()
		}
		s.triggerStatus = zscanproto.ToFState_ToFInactive
	}

	s.triggerThreshold = uint16(req.Threshold)

	if previousState == zscanproto.ToFState_ToFActive {
		loopCtx, cancel := context.WithCancel(context.Background())
		s.triggerDeactivation = cancel
		s.triggerStatus = zscanproto.ToFState_ToFActive
		go s.StartToFMonitor(loopCtx, s.tofEventStream)
	}

	return &zscanproto.Status{
		Success: true,
		Message: fmt.Sprintf("ToF configured : Threshold set to (threshold=%d)mm", req.Threshold),
	}, nil
}

func (s *ControllerServer) EnableToF(ctx context.Context, _ *zscanproto.Empty) (*zscanproto.Status, error) {
	if s.triggerStatus != zscanproto.ToFState_ToFDisabled {
		return &zscanproto.Status{Success: true, Message: "Already enabled"}, nil
	}

	triggerModule, error := zpi_trigger.NewZToF()
	if error != nil {
		log.Println("Failed To Initialize Trigger Module !")
		return &zscanproto.Status{Success: false, Message: "Failed to initialize ToF module"}, nil
	}
	s.trigger = triggerModule

	s.triggerStatus = zscanproto.ToFState_ToFInactive

	return &zscanproto.Status{Success: true, Message: "ToF enabled"}, nil
}

func (s *ControllerServer) DisableToF(ctx context.Context, _ *zscanproto.Empty) (*zscanproto.Status, error) {
	if s.triggerStatus == zscanproto.ToFState_ToFDisabled {
		return &zscanproto.Status{Success: true, Message: "Already disabled"}, nil
	}

	if s.triggerStatus == zscanproto.ToFState_ToFActive {
		if s.triggerDeactivation != nil {
			s.triggerDeactivation()
		}
	}

	s.triggerStatus = zscanproto.ToFState_ToFDisabled

	return &zscanproto.Status{Success: true, Message: "ToF disabled"}, nil
}

func (s *ControllerServer) ActivateToF(ctx context.Context, _ *zscanproto.Empty) (*zscanproto.Status, error) {
	if s.triggerStatus != zscanproto.ToFState_ToFInactive {
		return &zscanproto.Status{Success: false, Message: "ToF not ready for activation"}, nil
	}

	loopCtx, cancel := context.WithCancel(context.Background())
	s.triggerDeactivation = cancel

	s.triggerStatus = zscanproto.ToFState_ToFActive

	go s.StartToFMonitor(loopCtx, s.tofEventStream)

	return &zscanproto.Status{Success: true, Message: "ToF activated"}, nil
}

func (s *ControllerServer) DeactivateToF(ctx context.Context, _ *zscanproto.Empty) (*zscanproto.Status, error) {
	if s.triggerStatus != zscanproto.ToFState_ToFActive {
		return &zscanproto.Status{Success: true, Message: "Already inactive"}, nil
	}

	s.triggerDeactivation()

	s.triggerStatus = zscanproto.ToFState_ToFInactive

	return &zscanproto.Status{Success: true, Message: "ToF deactivated"}, nil
}

func (s *ControllerServer) StartToFMonitor(ctx context.Context, stream zscanproto.ZPiController_ToFEventStreamServer) {
	const clientTimeout = 10 * time.Second

	for {
		select {
		case <-ctx.Done():
			log.Println("ToF monitor stopped")
			return
		default:
			distance, err := s.trigger.Read()
			if err != nil {
				time.Sleep(100 * time.Millisecond)
				continue
			}

			if !s.eventStatus && distance < s.triggerThreshold {
				s.eventStatus = true
				s.SetLEDStatus(ctx, &zscanproto.LEDStatusRequest{
					Status: zscanproto.LEDStatus_FAILED,
				})

				trigger := &zscanproto.ToFEvent{
					Type: zscanproto.ToFEvent_TRIGGER,
				}
				if err := stream.Send(trigger); err != nil {
					log.Println("Failed to send ToF trigger:", err)
					s.eventStatus = false
					continue
				}
				log.Println("ToF trigger sent, waiting for client response...")

				completion := make(chan struct{})
				go func() {
					defer close(completion)
					for {
						clientEvent, err := stream.Recv()
						if err != nil {
							log.Println("Client disconnected or error:", err)
							return
						}

						switch clientEvent.Type {
						case zscanproto.ToFEvent_PENDING:
							_, _ = s.SetLEDStatus(ctx, &zscanproto.LEDStatusRequest{
								Status: zscanproto.LEDStatus_PENDING,
							})
							log.Println("Pending Auth")

						case zscanproto.ToFEvent_RESULT:
							_, _ = s.SetLEDStatus(ctx, &zscanproto.LEDStatusRequest{
								Status: clientEvent.LedStatus,
							})
							log.Println("Auth Result Received")

							time.Sleep(2 * time.Second)

							s.SetLEDStatus(ctx, &zscanproto.LEDStatusRequest{
								Status: zscanproto.LEDStatus_FAILED,
							})
							return
						}
					}
				}()

				select {
				case <-completion:
					log.Println("Client session completed")
				case <-time.After(clientTimeout):
					log.Println("Client did not respond within timeout, resetting session")
				}

				s.eventStatus = false
			}

			time.Sleep(50 * time.Millisecond)
		}
	}
}

func (c *ControllerServer) ActivateCamera(ctx context.Context, _ *zscanproto.Empty) (*zscanproto.Status, error) {
	c.camMutex.Lock()
	defer c.camMutex.Unlock()

	if c.camStatus == zscanproto.CamState_CAMACTIVE {
		return &zscanproto.Status{Success: true, Message: "Camera already active"}, nil
	}

	c.camProcess.Start(720, 720, 114, "tcp://0.0.0.0:8888")

	c.camStatus = zscanproto.CamState_CAMACTIVE
	return &zscanproto.Status{Success: true, Message: "Camera activation sequence started"}, nil
}

func (c *ControllerServer) DeactivateCamera(ctx context.Context, _ *zscanproto.Empty) (*zscanproto.Status, error) {
	c.camMutex.Lock()
	defer c.camMutex.Unlock()

	if c.camStatus != zscanproto.CamState_CAMACTIVE {
		return &zscanproto.Status{
			Success: true,
			Message: "Camera not active",
		}, nil
	}

	if err := c.camProcess.Stop(); err != nil {
		return &zscanproto.Status{
			Success: false,
			Message: fmt.Sprintf("Failed to stop camera: %v", err),
		}, nil
	}

	c.camStatus = zscanproto.CamState_CAMINACTIVE
	return &zscanproto.Status{
		Success: true,
		Message: "Camera deactivated",
	}, nil
}

func (c *ControllerServer) ConfigureCamera(ctx context.Context, cfg *zscanproto.CameraConfig) (*zscanproto.Status, error) {
	c.camMutex.Lock()
	defer c.camMutex.Unlock()

	_ = c.camProcess.Stop()

	time.Sleep(200 * time.Millisecond)

	c.camProcess.Start(cfg.Width, cfg.Height, cfg.Fps, "tcp://0.0.0.0:8888")

	c.camStatus = zscanproto.CamState_CAMACTIVE
	return &zscanproto.Status{
		Success: true,
		Message: fmt.Sprintf("Camera reconfigured: %dx%d @ %dfps", cfg.Width, cfg.Height, cfg.Fps),
	}, nil
}
