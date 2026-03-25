package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	grpcServer := grpc.NewServer()

	zpi_listener, error := net.Listen("tcp", ":50051")
	if error != nil {
		log.Fatalf("failed to listen: %v", error)
	}

	log.Print("ZPi GRPC Running !")
	grpcServer.Serve(zpi_listener)

}
