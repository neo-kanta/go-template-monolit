package samplegrpc

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"go-transfer-agent/common/platform/config"
	pb "go-transfer-agent/common/gen/sample/v1"
)

type Server struct {
	pb.UnimplementedSampleServiceServer
	log *slog.Logger
}

func NewServer(cfg *config.Config, log *slog.Logger) *grpc.Server {
	grpcServer := grpc.NewServer()
	
	srv := &Server{
		log: log,
	}

	pb.RegisterSampleServiceServer(grpcServer, srv)
	reflection.Register(grpcServer)

	return grpcServer
}

func (s *Server) Echo(ctx context.Context, req *pb.EchoRequest) (*pb.EchoResponse, error) {
	s.log.Info("Echo request received", slog.String("message", req.GetMessage()))
	return &pb.EchoResponse{
		Message: "Echo: " + req.GetMessage(),
	}, nil
}

func (s *Server) HealthCheck(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	return &pb.HealthCheckResponse{
		Status: "OK",
	}, nil
}
