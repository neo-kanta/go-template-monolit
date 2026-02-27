package fndm009

import (
	"context"
	"log/slog"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
)

// Handler implements the FNDM009 gRPC service interface.
type Handler struct {
	fndv1.UnimplementedFNDM009ServiceServer
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

// TAFNDFavDisc - GET query endpoint
func (h *Handler) TAFNDFavDisc(ctx context.Context, req *fndv1.TAFNDFavDiscRequest) (*fndv1.TAFNDFavDiscResponse, error) {
	h.log.Info("TAFNDFavDisc",
		slog.String("sys_co_id", req.GetSysCoId()),
		slog.String("cus_id_code", req.GetCusIdCode()),
	)
	return h.svc.TAFNDFavDisc(ctx, req)
}

// SaveTAFNDFavDisc - POST save (add) endpoint
func (h *Handler) SaveTAFNDFavDisc(ctx context.Context, req *fndv1.TAFNDFavDiscRequest) (*fndv1.SaveResponse, error) {
	h.log.Info("SaveTAFNDFavDisc",
		slog.String("sys_co_id", req.GetSysCoId()),
		slog.String("cus_id_code", req.GetCusIdCode()),
	)
	return h.svc.SaveTAFNDFavDisc(ctx, req)
}

// UpdateTAFNDFavDisc - PUT save (modify) endpoint
func (h *Handler) UpdateTAFNDFavDisc(ctx context.Context, req *fndv1.TAFNDFavDiscRequest) (*fndv1.SaveResponse, error) {
	h.log.Info("UpdateTAFNDFavDisc",
		slog.String("sys_co_id", req.GetSysCoId()),
		slog.String("cus_id_code", req.GetCusIdCode()),
	)
	return h.svc.UpdateTAFNDFavDisc(ctx, req)
}

// DeleteTAFNDFavDisc - DELETE endpoint
func (h *Handler) DeleteTAFNDFavDisc(ctx context.Context, req *fndv1.DeleteRequest) (*fndv1.SaveResponse, error) {
	h.log.Info("DeleteTAFNDFavDisc",
		slog.String("data_id", req.GetDataId()),
		slog.String("data_flag", req.GetDataFlag()),
	)
	return h.svc.DeleteTAFNDFavDisc(ctx, req)
}

// ApproveTAFNDFavDisc handles the Approval POST endpoint.
func (h *Handler) ApproveTAFNDFavDisc(ctx context.Context, req *fndv1.ApproveTAFNDFavDiscRequest) (*fndv1.SaveResponse, error) {
	h.log.Info("ApproveTAFNDFavDisc called", slog.String("data_id", req.GetDataId()))
	return h.svc.ApproveTAFNDFavDisc(ctx, req)
}
