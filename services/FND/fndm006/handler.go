package fndm006

import (
	"context"
	"log/slog"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"

	"google.golang.org/grpc"
)

// Handler implements the FNDM006 gRPC methods.
type Handler struct {
	log *slog.Logger
	svc *Service
}

// NewHandler creates a new FNDM006 gRPC handler.
func NewHandler(log *slog.Logger, svc *Service) *Handler {
	return &Handler{
		log: log.With(slog.String("module", "fndm006")),
		svc: svc,
	}
}

// Service exposes the underlying business service.
func (h *Handler) Service() *Service { return h.svc }

// TAFNDRPFeeChgType - GET endpoint implementation
func (h *Handler) TAFNDRPFeeChgType(ctx context.Context, req *fndv1.TAFNDRPFeeChgTypeRequest) (*fndv1.TAFNDRPFeeChgTypeResponse, error) {
	h.log.Info("TAFNDRPFeeChgType",
		slog.String("sys_co_id", req.GetSysCoId()),
		slog.String("prt_fund_code", req.GetPrtFundCode()),
	)
	return h.svc.TAFNDRPFeeChgType(ctx, req)
}

// RegisterHandlers placeholder for gRPC server registration
func RegisterHandlers(_ *grpc.Server) {
	// fndv1.RegisterFNDM006ServiceServer(srv, NewHandler(log, svc))
}
