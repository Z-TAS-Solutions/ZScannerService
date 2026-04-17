package main

import (
	"github.com/Z-TAS-Solutions/ZScannerService/api/grpc/zpi_server"
	zpi_camera "github.com/Z-TAS-Solutions/ZScannerService/internal/pkg/zpi_cam"
)

func main() {
	camera := zpi_camera.NewCameraProcess()

	camera.Start(720, 720, 113, "tcp://0.0.0.0:8888")

	zpi_server.RunZPiServer()

}
