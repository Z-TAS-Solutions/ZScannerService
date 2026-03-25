package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
)

func RunLocalGRPC() {

	grpcServer := grpc.NewServer()

	localServer, error := net.Listen("tcp", ":50051")
	if error != nil {
		log.Fatalf("failed to listen: %v", error)
	}

	grpcServer.Serve(localServer)

}

func main() {
	RunLocalGRPC()
}
