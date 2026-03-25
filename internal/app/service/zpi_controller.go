package zpi_controller

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Z-TAS-Solutions/ZScannerService/internal/pkg/zpi_indicator"
	zpi_trigger "github.com/Z-TAS-Solutions/ZScannerService/internal/pkg/zpi_tof"
	"github.com/Z-TAS-Solutions/ZScannerService/internal/pkg/zscanproto"
)

type ControllerServer struct {
	zscanproto.UnimplementedZPiControllerServer
	indicator       *zpi_indicator.ZLED
	indicatorStatus zscanproto.LEDStatus
	trigger         *zpi_trigger.ZToF
	peerClient      zscanproto.ZPiControllerClient
}

func NewControllerServer(indicator *zpi_indicator.ZLED, trigger *zpi_trigger.ZToF) *ControllerServer {
	return &ControllerServer{indicator: indicator, trigger: trigger}
}

var statusColors = map[zscanproto.LEDStatus][3]uint32{
	zscanproto.LEDStatus_VOID:    {0, 0, 0},
	zscanproto.LEDStatus_SUCCESS: {0, 255, 0},
	zscanproto.LEDStatus_PENDING: {255, 255, 0},
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

func (s *ControllerServer) SetPeerClient(client zscanproto.ZPiControllerClient) {
	s.peerClient = client
}

func (s *ControllerServer) StartToFMonitor(threshold uint16) {
	go func() {
		for {
			distance, err := s.trigger.Read()
			if err != nil {
				log.Println("ToF read error:", err)
				time.Sleep(100 * time.Millisecond)
				continue
			}

			var newState zscanproto.LEDStatus
			if distance < threshold {
				newState = zscanproto.LEDStatus_VOID
			} else {
				newState = zscanproto.LEDStatus_FAILED
			}

			if newState != s.indicatorStatus {
				_, _ = s.SetLEDStatus(context.Background(), &zscanproto.LEDStatusRequest{
					Status: newState,
				})
				s.indicatorStatus = newState

				if s.peerClient != nil {
					ctx, cancel := context.WithTimeout(context.Background(), time.Second)
					_, err := s.peerClient.SetLEDStatus(ctx, &zscanproto.LEDStatusRequest{
						Status: newState,
					})
					cancel()
					if err != nil {
						log.Println("Failed to notify peer:", err)
					}
				}
			}

			time.Sleep(50 * time.Millisecond)
		}
	}()
}
