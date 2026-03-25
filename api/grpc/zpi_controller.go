package zpi_controller

import (
	"context"

	zpi_indicator "github.com/Z-TAS-Solutions/Z-QryptGIN/internal/pkg"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/pkg/zproto"
)

type ZPiControllerServer struct {
	zproto.UnimplementedZPiControllerServer
	led *zpi_indicator.ZLED
}

func NewZPiControllerServer(ledModule *zpi_indicator.ZLED) *ZPiControllerServer {
	return &ZPiControllerServer{
		led: ledModule,
	}
}

func (s *ZPiControllerServer) SetLED(ctx context.Context, req *zproto.LEDRequest) (*zproto.Status, error) {
	err := s.led.Set(req.Red, req.Green, req.Blue)
	if err != nil {
		return &zproto.Status{Success: false, Message: err.Error()}, nil
	}
	return &zproto.Status{Success: true, Message: "LED updated"}, nil
}

func (s *ZPiControllerServer) GetLED(ctx context.Context, _ *zproto.Empty) (*zproto.LEDState, error) {
	r, g, b := s.led.Get()
	return &zproto.LEDState{Red: r, Green: g, Blue: b}, nil
}
