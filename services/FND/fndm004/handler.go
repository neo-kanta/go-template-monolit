package fndm004

import (
	"context"
	"log/slog"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"

	"google.golang.org/grpc"
)

// Handler implements the FNDM004 gRPC methods.
type Handler struct {
	log *slog.Logger
	svc *Service
}

// NewHandler creates a new FNDM004 gRPC handler.
func NewHandler(log *slog.Logger, svc *Service) *Handler {
	return &Handler{
		log: log.With(slog.String("module", "fndm004")),
		svc: svc,
	}
}

// Service exposes the underlying business service.
func (h *Handler) Service() *Service { return h.svc }

// TAFNDFundAgent - GET endpoint implementation
func (h *Handler) TAFNDFundAgent(ctx context.Context, req *fndv1.TAFNDFundAgentRequest) (*fndv1.TAFNDFundAgentResponse, error) {
	h.log.Info("TAFNDFundAgent",
		slog.String("sys_co_id", req.GetSysCoId()),
		slog.String("prt_fund_code", req.GetPrtFundCode()),
	)
	return h.svc.TAFNDFundAgent(ctx, req)
}

// RegisterHandlers placeholder for gRPC server registration
func RegisterHandlers(_ *grpc.Server) {
	// fndv1.RegisterFNDM004ServiceServer(srv, NewHandler(log, svc))
}
