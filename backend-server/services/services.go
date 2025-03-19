package services

import (
	"context"
	"log/slog"
	"time"

	"backend-server/models"
    pb "backend-server/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ConnectToServer(server models.ServerConfig) {
	conn, err := grpc.Dial(server.Address, grpc.WithTransportCredentials(insecure.NewCredentials())) // Use WithTransportCredentials for TLS
	if err != nil {
		slog.Error("Failed to connect to daemon via gRPC)", server.Name, server.Address, "Error", err.Error())
	}
	defer conn.Close()

	client := pb.NewDaemonServiceClient(conn)

	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		response, err := client.GetStatus(ctx, &pb.StatusRequest{ServerId: server.Name})
		cancel()

		if err != nil {
			slog.Error("Error calling GetStatus", server.Name, server.Address, "Error", err.Error())
		} else {
			slog.Info("Status", server.Name, response.Status)
		}

		time.Sleep(10 * time.Second)
	}
}
