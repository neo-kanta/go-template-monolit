package grpc

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// NewServer creates a gRPC server with the health check service registered.
func NewServer(log *slog.Logger) *grpc.Server {
	srv := grpc.NewServer()

	// Register the standard gRPC health check service.
	healthSrv := health.NewServer()
	healthSrv.SetServingStatus("fnd.v1.FNDService", healthpb.HealthCheckResponse_SERVING)
	healthSrv.SetServingStatus("", healthpb.HealthCheckResponse_SERVING) // overall
	healthpb.RegisterHealthServer(srv, healthSrv)

	log.Info("registered gRPC health service")

	// Future: register FNDServiceServer here once proto is generated.
	// fndv1.RegisterFNDServiceServer(srv, &fndHandler{...})

	_ = context.Background() // placeholder to avoid unused import in future

	return srv
}
