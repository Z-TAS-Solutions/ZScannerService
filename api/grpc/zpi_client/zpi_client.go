package zpi_client

import (
	"log"

	"github.com/Z-TAS-Solutions/ZScannerService/internal/pkg/zscanproto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func RunZPiClient(ip string) zscanproto.ZPiControllerClient {
	conn, err := grpc.Dial(ip, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	client := zscanproto.NewZPiControllerClient(conn)
	log.Print("Connected to ZPi GRPC Host!")

	return client
}
