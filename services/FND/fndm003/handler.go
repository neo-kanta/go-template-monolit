package fndm003

import (
	"context"
	"log/slog"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"

	"google.golang.org/grpc"
)

// Handler implements the FNDM003 gRPC methods.
type Handler struct {
	log *slog.Logger
	svc *Service
}

// NewHandler creates a new FNDM003 gRPC handler.
func NewHandler(log *slog.Logger, svc *Service) *Handler {
	return &Handler{
		log: log.With(slog.String("module", "fndm003")),
		svc: svc,
	}
}

// Service exposes the underlying business service.
func (h *Handler) Service() *Service { return h.svc }

// TAFNDSwitch - GET endpoint implementation
func (h *Handler) TAFNDSwitch(ctx context.Context, req *fndv1.TAFNDSwitchRequest) (*fndv1.TAFNDSwitchResponse, error) {
	h.log.Info("TAFNDSwitch",
		slog.String("sys_co_id", req.GetSysCoId()),
		slog.String("prt_fund_code", req.GetPrtFundCode()),
	)
	return h.svc.TAFNDSwitch(ctx, req)
}

// RegisterHandlers placeholder for gRPC server registration
func RegisterHandlers(_ *grpc.Server) {
	// fndv1.RegisterFNDM003ServiceServer(srv, NewHandler(log, svc))
}
