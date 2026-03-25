package zpi_server

import (
	"log"
	"net"

	zpi_controller "github.com/Z-TAS-Solutions/ZScannerService/internal/app/service"
	"github.com/Z-TAS-Solutions/ZScannerService/internal/pkg/zpi_indicator"
	zpi_trigger "github.com/Z-TAS-Solutions/ZScannerService/internal/pkg/zpi_tof"
	"github.com/Z-TAS-Solutions/ZScannerService/internal/pkg/zscanproto"
	"google.golang.org/grpc"
)

func RunZPiServer() {
	grpcServer := grpc.NewServer()

	indicatorModule := zpi_indicator.NewLED(14, 15, 18)
	triggerModule, error := zpi_trigger.NewZToF()
	if error != nil {
		log.Println("Failed To Initialize Trigger Module !")
	}

	controllerServer := zpi_controller.NewControllerServer(indicatorModule, triggerModule)
	zscanproto.RegisterZPiControllerServer(grpcServer, controllerServer)

	zpi_listener, error := net.Listen("tcp", ":50051")
	if error != nil {
		log.Fatalf("failed to listen: %v", error)
	}

	log.Print("ZPi GRPC Running !")
	grpcServer.Serve(zpi_listener)

}
