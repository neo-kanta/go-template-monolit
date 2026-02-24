package fndm010

import (
	"context"
	"log/slog"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"

	"google.golang.org/grpc"
)

// Handler implements the FNDM010 gRPC methods.
type Handler struct {
	log *slog.Logger
	svc *Service
}

// NewHandler creates a new FNDM010 gRPC handler.
func NewHandler(log *slog.Logger, svc *Service) *Handler {
	return &Handler{
		log: log.With(slog.String("module", "fndm010")),
		svc: svc,
	}
}

// Service exposes the underlying business service.
func (h *Handler) Service() *Service { return h.svc }

// TAFNDIShareFundFee - GET endpoint implementation
func (h *Handler) TAFNDIShareFundFee(ctx context.Context, req *fndv1.TAFNDIShareFundFeeRequest) (*fndv1.TAFNDIShareFundFeeResponse, error) {
	h.log.Info("TAFNDIShareFundFee",
		slog.String("sys_co_id", req.GetSysCoId()),
		slog.String("fund_code", req.GetFundCode()),
	)
	return h.svc.TAFNDIShareFundFee(ctx, req)
}

// RegisterHandlers placeholder for gRPC server registration
func RegisterHandlers(_ *grpc.Server) {
	// fndv1.RegisterFNDM010ServiceServer(srv, NewHandler(log, svc))
}
