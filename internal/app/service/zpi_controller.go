import (
	"context"

	"github.com/Z-TAS-Solutions/ZScannerService/internal/pkg/zproto"
)

type ControllerServer struct {
	zproto.UnimplementedZPiControllerServer
	led *zpi_indicator.ZLED
}

func NewControllerServer(l *zpi_indicator.ZLED) *ControllerServer {
	return &ControllerServer{led: l}
}

func (s *ControllerServer) SetLED(ctx context.Context, req *zproto.LEDRequest) (*zproto.Status, error) {
	err := s.led.Set(req.Red, req.Green, req.Blue)
	if err != nil {
		return &zproto.Status{Success: false, Message: err.Error()}, nil
	}
	return &zproto.Status{Success: true, Message: "LED updated"}, nil
}

func (s *ControllerServer) GetLED(ctx context.Context, _ *zproto.Empty) (*zproto.LEDState, error) {
	r, g, b := s.led.Get()
	return &zproto.LEDState{Red: r, Green: g, Blue: b}, nil
}

