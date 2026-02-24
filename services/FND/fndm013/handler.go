package fndm013

import (
	"context"
	"log/slog"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"

	"google.golang.org/grpc"
)

// Handler implements the FNDM013 gRPC methods.
type Handler struct {
	log *slog.Logger
	svc *Service
}

// NewHandler creates a new FNDM013 gRPC handler.
func NewHandler(log *slog.Logger, svc *Service) *Handler {
	return &Handler{
		log: log.With(slog.String("module", "fndm013")),
		svc: svc,
	}
}

// Service exposes the underlying business service.
func (h *Handler) Service() *Service { return h.svc }

// TAFNDTMFundFeeRdm - GET endpoint implementation
func (h *Handler) TAFNDTMFundFeeRdm(ctx context.Context, req *fndv1.TAFNDTMFundFeeRdmRequest) (*fndv1.TAFNDTMFundFeeRdmResponse, error) {
	h.log.Info("TAFNDTMFundFeeRdm",
		slog.String("sys_co_id", req.GetSysCoId()),
		slog.String("prt_fund_code", req.GetPrtFundCode()),
	)
	return h.svc.TAFNDTMFundFeeRdm(ctx, req)
}

// RegisterHandlers placeholder for gRPC server registration
func RegisterHandlers(_ *grpc.Server) {
	// fndv1.RegisterFNDM013ServiceServer(srv, NewHandler(log, svc))
}
