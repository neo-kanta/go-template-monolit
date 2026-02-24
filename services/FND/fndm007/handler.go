package fndm007

import (
	"context"
	"log/slog"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"

	"google.golang.org/grpc"
)

// Handler implements the FNDM007 gRPC methods.
type Handler struct {
	log *slog.Logger
	svc *Service
}

// NewHandler creates a new FNDM007 gRPC handler.
func NewHandler(log *slog.Logger, svc *Service) *Handler {
	return &Handler{
		log: log.With(slog.String("module", "fndm007")),
		svc: svc,
	}
}

// Service exposes the underlying business service.
func (h *Handler) Service() *Service { return h.svc }

// TAFNDCustGroup - GET endpoint implementation
func (h *Handler) TAFNDCustGroup(ctx context.Context, req *fndv1.TAFNDCustGroupRequest) (*fndv1.TAFNDCustGroupResponse, error) {
	h.log.Info("TAFNDCustGroup",
		slog.String("sys_co_id", req.GetSysCoId()),
		slog.String("cust_grp_code", req.GetCustGrpCode()),
	)
	return h.svc.TAFNDCustGroup(ctx, req)
}

// RegisterHandlers placeholder for gRPC server registration
func RegisterHandlers(_ *grpc.Server) {
	// fndv1.RegisterFNDM007ServiceServer(srv, NewHandler(log, svc))
}
