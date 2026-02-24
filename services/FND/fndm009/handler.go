package fndm009

import (
	"context"
	"log/slog"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"

	"google.golang.org/grpc"
)

// Handler implements the FNDM009 gRPC methods.
type Handler struct {
	log *slog.Logger
	svc *Service
}

// NewHandler creates a new FNDM009 gRPC handler.
func NewHandler(log *slog.Logger, svc *Service) *Handler {
	return &Handler{
		log: log.With(slog.String("module", "fndm009")),
		svc: svc,
	}
}

// Service exposes the underlying business service.
func (h *Handler) Service() *Service { return h.svc }

// TAFNDFavDisc - GET endpoint implementation
func (h *Handler) TAFNDFavDisc(ctx context.Context, req *fndv1.TAFNDFavDiscRequest) (*fndv1.TAFNDFavDiscResponse, error) {
	h.log.Info("TAFNDFavDisc",
		slog.String("sys_co_id", req.GetSysCoId()),
		slog.String("cus_id_code", req.GetCusIdCode()),
	)
	return h.svc.TAFNDFavDisc(ctx, req)
}

// RegisterHandlers placeholder for gRPC server registration
func RegisterHandlers(_ *grpc.Server) {
	// fndv1.RegisterFNDM009ServiceServer(srv, NewHandler(log, svc))
}
