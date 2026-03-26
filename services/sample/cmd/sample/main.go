package main

import (
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go-transfer-agent/common/platform/config"
	"go-transfer-agent/common/platform/logger"
	samplegrpc "go-transfer-agent/services/sample/internal/adapter/grpc"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.LogLevel)

	// ─── gRPC Server ───
	// Use SAMPLE_SERVICE_ADDR or fallback to GRPC_ADDR
	grpcAddr := cfg.GRPCAddr
	// In a real scenario, you'd add this to config: grpcAddr := cfg.SampleServiceGRPCAddr
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Error("failed to listen (gRPC)", slog.String("addr", grpcAddr), slog.String("error", err.Error()))
		os.Exit(1)
	}

	grpcSrv := samplegrpc.NewServer(cfg, log)

	// Start gRPC server in background
	go func() {
		log.Info("🚀 Sample gRPC service starting", slog.String("addr", grpcAddr))
		if err := grpcSrv.Serve(lis); err != nil {
			log.Error("gRPC server failed", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigCh
	log.Info("shutting down", slog.String("signal", sig.String()))
	
	grpcSrv.GracefulStop()
}
