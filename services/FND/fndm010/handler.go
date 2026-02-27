package fndm010

import (
	"context"
	"log/slog"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
)

// Handler implements the FNDM010 gRPC service interface.
type Handler struct {
	fndv1.UnimplementedFNDM010ServiceServer
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

// TAFNDIShareFundFee - GET query endpoint
func (h *Handler) TAFNDIShareFundFee(ctx context.Context, req *fndv1.TAFNDIShareFundFeeRequest) (*fndv1.TAFNDIShareFundFeeResponse, error) {
	h.log.Info("TAFNDIShareFundFee",
		slog.String("sys_co_id", req.GetSysCoId()),
		slog.String("fund_code", req.GetFundCode()),
	)
	return h.svc.TAFNDIShareFundFee(ctx, req)
}

// SaveTAFNDIShareFundFee - POST save (add) endpoint
func (h *Handler) SaveTAFNDIShareFundFee(ctx context.Context, req *fndv1.TAFNDIShareFundFeeRequest) (*fndv1.SaveResponse, error) {
	h.log.Info("SaveTAFNDIShareFundFee",
		slog.String("sys_co_id", req.GetSysCoId()),
		slog.String("fund_code", req.GetFundCode()),
	)
	return h.svc.SaveTAFNDIShareFundFee(ctx, req)
}

// UpdateTAFNDIShareFundFee - PUT save (modify) endpoint
func (h *Handler) UpdateTAFNDIShareFundFee(ctx context.Context, req *fndv1.TAFNDIShareFundFeeRequest) (*fndv1.SaveResponse, error) {
	h.log.Info("UpdateTAFNDIShareFundFee",
		slog.String("sys_co_id", req.GetSysCoId()),
		slog.String("fund_code", req.GetFundCode()),
	)
	return h.svc.UpdateTAFNDIShareFundFee(ctx, req)
}

// DeleteTAFNDIShareFundFee - DELETE endpoint
func (h *Handler) DeleteTAFNDIShareFundFee(ctx context.Context, req *fndv1.DeleteRequest) (*fndv1.SaveResponse, error) {
	h.log.Info("DeleteTAFNDIShareFundFee",
		slog.String("data_id", req.GetDataId()),
		slog.String("data_flag", req.GetDataFlag()),
	)
	return h.svc.DeleteTAFNDIShareFundFee(ctx, req)
}

// ApproveTAFNDIShareFundFee handles the Approval POST endpoint.
func (h *Handler) ApproveTAFNDIShareFundFee(ctx context.Context, req *fndv1.ApproveTAFNDIShareFundFeeRequest) (*fndv1.SaveResponse, error) {
	h.log.Info("ApproveTAFNDIShareFundFee called", slog.String("data_id", req.GetDataId()))
	return h.svc.ApproveTAFNDIShareFundFee(ctx, req)
}
