// Package main is the entry point for the FND gRPC service.
package main

import (
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go-transfer-agent/common/platform/config"
	"go-transfer-agent/common/platform/logger"
	"go-transfer-agent/services/FND/internal/adapter/grpc"

	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.LogLevel)

	lis, err := net.Listen("tcp", cfg.GRPCAddr)
	if err != nil {
		log.Error("failed to listen", slog.String("addr", cfg.GRPCAddr), slog.String("error", err.Error()))
		os.Exit(1)
	}

	srv := grpc.NewServer(log)

	// Enable gRPC server reflection for tools like grpcurl.
	reflection.Register(srv)

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigCh
		log.Info("shutting down", slog.String("signal", sig.String()))
		srv.GracefulStop()
	}()

	log.Info("🚀 FND service starting", slog.String("addr", cfg.GRPCAddr))
	if err := srv.Serve(lis); err != nil {
		log.Error("server failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
