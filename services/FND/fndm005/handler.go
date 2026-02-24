package fndm005

import (
	"context"
	"log/slog"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"

	"google.golang.org/grpc"
)

// Handler implements the FNDM005 gRPC methods.
type Handler struct {
	log *slog.Logger
	svc *Service
}

// NewHandler creates a new FNDM005 gRPC handler.
func NewHandler(log *slog.Logger, svc *Service) *Handler {
	return &Handler{
		log: log.With(slog.String("module", "fndm005")),
		svc: svc,
	}
}

// Service exposes the underlying business service.
func (h *Handler) Service() *Service { return h.svc }

// TAFNDFundCalDate - GET endpoint implementation
func (h *Handler) TAFNDFundCalDate(ctx context.Context, req *fndv1.TAFNDFundCalDateRequest) (*fndv1.TAFNDFundCalDateResponse, error) {
	h.log.Info("TAFNDFundCalDate",
		slog.String("sys_co_id", req.GetSysCoId()),
		slog.String("cal_year", req.GetCalYear()),
	)
	return h.svc.TAFNDFundCalDate(ctx, req)
}

// RegisterHandlers placeholder for gRPC server registration
func RegisterHandlers(_ *grpc.Server) {
	// fndv1.RegisterFNDM005ServiceServer(srv, NewHandler(log, svc))
}
