package main

import (
	"log/slog"
	"net"

    "file-server-daemon/models"
    pb "file-server-daemon/proto"

	"google.golang.org/grpc"
)

func main() {
    lis, err := net.Listen("tcp", ":4444")
    if err != nil {
        slog.Error("Failed to listen", "Error", err.Error())
        return
    }

    grpcServer := grpc.NewServer()
    pb.RegisterDaemonServiceServer(grpcServer, &models.Server{})

    slog.Info("gRPC server listening on :4444")
    if err := grpcServer.Serve(lis); err != nil {
        slog.Error("Failed to serve", "Error", err.Error())
    }
}
