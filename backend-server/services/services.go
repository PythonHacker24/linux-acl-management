package services

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"backend-server/models"
	pb "backend-server/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ConnectToServer(server models.Server) error {
    if server.Remote == nil {
        return fmt.Errorf("remote server configuration is nil")
    }
    address := fmt.Sprintf("%s:%d", server.Remote.Host, server.Remote.Port)
    conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        slog.Error("Failed to connect to daemon via gRPC", "host", server.Remote.Host, "port", server.Remote.Port, "error", err.Error())
        return err
    }
    defer conn.Close()
    
    client := pb.NewDaemonServiceClient(conn)
    for {
        ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
        response, err := client.GetStatus(ctx, &pb.StatusRequest{ServerId: server.Remote.Host})
        cancel()
        if err != nil {
            slog.Error("Error calling GetStatus", "host", server.Remote.Host, "port", server.Remote.Port, "error", err.Error())
        } else {
            slog.Info("Status", "host", server.Remote.Host, "status", response.Status)
        }
        time.Sleep(10 * time.Second)
    }
    return nil
}
