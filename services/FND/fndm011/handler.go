package fndm011

import (
	"context"
	"log/slog"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"

	"google.golang.org/grpc"
)

// Handler implements the FNDM011 gRPC methods.
type Handler struct {
	log *slog.Logger
	svc *Service
}

// NewHandler creates a new FNDM011 gRPC handler.
func NewHandler(log *slog.Logger, svc *Service) *Handler {
	return &Handler{
		log: log.With(slog.String("module", "fndm011")),
		svc: svc,
	}
}

// Service exposes the underlying business service.
func (h *Handler) Service() *Service { return h.svc }

// TAFNDIShareFundFeeRdm - GET endpoint implementation
func (h *Handler) TAFNDIShareFundFeeRdm(ctx context.Context, req *fndv1.TAFNDIShareFundFeeRdmRequest) (*fndv1.TAFNDIShareFundFeeRdmResponse, error) {
	h.log.Info("TAFNDIShareFundFeeRdm",
		slog.String("sys_co_id", req.GetSysCoId()),
		slog.String("fund_code", req.GetFundCode()),
	)
	return h.svc.TAFNDIShareFundFeeRdm(ctx, req)
}

// RegisterHandlers placeholder for gRPC server registration
func RegisterHandlers(_ *grpc.Server) {
	// fndv1.RegisterFNDM011ServiceServer(srv, NewHandler(log, svc))
}
