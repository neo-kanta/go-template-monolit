package fndm008

import (
	"context"
	"log/slog"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"

	"google.golang.org/grpc"
)

// Handler implements the FNDM008 gRPC methods.
type Handler struct {
	log *slog.Logger
	svc *Service
}

// NewHandler creates a new FNDM008 gRPC handler.
func NewHandler(log *slog.Logger, svc *Service) *Handler {
	return &Handler{
		log: log.With(slog.String("module", "fndm008")),
		svc: svc,
	}
}

// Service exposes the underlying business service.
func (h *Handler) Service() *Service { return h.svc }

// TAFNDPauseTxn - GET endpoint implementation
func (h *Handler) TAFNDPauseTxn(ctx context.Context, req *fndv1.TAFNDPauseTxnRequest) (*fndv1.TAFNDPauseTxnResponse, error) {
	h.log.Info("TAFNDPauseTxn",
		slog.String("sys_co_id", req.GetSysCoId()),
		slog.String("prt_fund_code", req.GetPrtFundCode()),
	)
	return h.svc.TAFNDPauseTxn(ctx, req)
}

// RegisterHandlers placeholder for gRPC server registration
func RegisterHandlers(_ *grpc.Server) {
	// fndv1.RegisterFNDM008ServiceServer(srv, NewHandler(log, svc))
}
