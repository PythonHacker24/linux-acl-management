package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"

	"file-server-daemon/models"
	pb "file-server-daemon/proto"

	"google.golang.org/grpc"
)

func main() {
    if len(os.Args) < 2 {
		slog.Error("Usage: go run main.go <port>")
		os.Exit(1)
	}

    port := os.Args[1]

    lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
    if err != nil {
        slog.Error("Failed to listen", "Error", err.Error())
        return
    }

    grpcServer := grpc.NewServer()
    pb.RegisterDaemonServiceServer(grpcServer, &models.Server{})

    slog.Info("gRPC server started", "port", port)
    if err := grpcServer.Serve(lis); err != nil {
        slog.Error("Failed to serve", "Error", err.Error())
    }
}
