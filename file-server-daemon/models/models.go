package models

import (
    "context"
    "log/slog"
    
    pb "file-server-daemon/proto"
)

type Server struct {
    pb.UnimplementedDaemonServiceServer
}

func (s *Server) GetStatus(ctx context.Context, req *pb.StatusRequest) (*pb.StatusResponse, error) {
    slog.Info("Received request for server: %s", "Server ID", req.ServerId)
    return &pb.StatusResponse{Status: "Running"}, nil
}
