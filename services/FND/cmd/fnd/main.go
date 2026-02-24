// Package main is the entry point for the FND gRPC + HTTP service.
package main

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"go-transfer-agent/common/platform/config"
	"go-transfer-agent/common/platform/logger"
	"go-transfer-agent/services/fnd/internal/adapter/gateway"
	fndgrpc "go-transfer-agent/services/fnd/internal/adapter/grpc"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.LogLevel)

	// ─── gRPC Server ───
	lis, err := net.Listen("tcp", cfg.GRPCAddr)
	if err != nil {
		log.Error("failed to listen (gRPC)", slog.String("addr", cfg.GRPCAddr), slog.String("error", err.Error()))
		os.Exit(1)
	}

	grpcSrv := fndgrpc.NewServer(cfg, log)

	// Start gRPC server in background (must be listening before gateway dials it)
	go func() {
		log.Info("🚀 FND gRPC service starting", slog.String("addr", cfg.GRPCAddr))
		if err := grpcSrv.Serve(lis); err != nil {
			log.Error("gRPC server failed", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	// Give gRPC a moment to start accepting connections
	time.Sleep(200 * time.Millisecond)

	// ─── HTTP/REST Gateway + Swagger ───
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	httpHandler, err := gateway.NewHTTPServer(ctx, cfg.GRPCAddr, log)
	if err != nil {
		log.Error("failed to create HTTP gateway", slog.String("error", err.Error()))
		os.Exit(1)
	}

	httpSrv := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: httpHandler,
	}

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigCh
		log.Info("shutting down", slog.String("signal", sig.String()))
		grpcSrv.GracefulStop()
		httpSrv.Shutdown(context.Background())
	}()

	// Start HTTP server (blocking)
	colonIdx := strings.LastIndex(cfg.HTTPAddr, ":")
	port := cfg.HTTPAddr[colonIdx:]
	log.Info("📖 HTTP gateway + Swagger UI starting",
		slog.String("addr", cfg.HTTPAddr),
		slog.String("swagger", "http://localhost"+port+"/swagger/"),
	)
	if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error("HTTP server failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
