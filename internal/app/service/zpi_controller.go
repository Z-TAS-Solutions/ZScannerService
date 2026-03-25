package zpi_controller

import (
	"context"
	"fmt"

	"github.com/Z-TAS-Solutions/ZScannerService/internal/pkg/zpi_indicator"
	"github.com/Z-TAS-Solutions/ZScannerService/internal/pkg/zscanproto"
)

type ControllerServer struct {
	zscanproto.UnimplementedZPiControllerServer
	indicator *zpi_indicator.ZLED
}

func NewControllerServer(l *zpi_indicator.ZLED) *ControllerServer {
	return &ControllerServer{indicator: l}
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
