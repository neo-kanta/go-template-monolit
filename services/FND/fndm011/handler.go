package fndm011

import (
	"context"
	"log/slog"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
)

// Handler implements the FNDM011 gRPC service interface.
type Handler struct {
	fndv1.UnimplementedFNDM011ServiceServer
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

// TAFNDIShareFundFeeRdm - GET query endpoint
func (h *Handler) TAFNDIShareFundFeeRdm(ctx context.Context, req *fndv1.TAFNDIShareFundFeeRdmRequest) (*fndv1.TAFNDIShareFundFeeRdmResponse, error) {
	h.log.Info("TAFNDIShareFundFeeRdm",
		slog.String("sys_co_id", req.GetSysCoId()),
		slog.String("fund_code", req.GetFundCode()),
	)
	return h.svc.TAFNDIShareFundFeeRdm(ctx, req)
}

// SaveTAFNDIShareFundFeeRdm - POST save (add) endpoint
func (h *Handler) SaveTAFNDIShareFundFeeRdm(ctx context.Context, req *fndv1.TAFNDIShareFundFeeRdmRequest) (*fndv1.SaveResponse, error) {
	h.log.Info("SaveTAFNDIShareFundFeeRdm",
		slog.String("sys_co_id", req.GetSysCoId()),
		slog.String("fund_code", req.GetFundCode()),
	)
	return h.svc.SaveTAFNDIShareFundFeeRdm(ctx, req)
}

// UpdateTAFNDIShareFundFeeRdm - PUT save (modify) endpoint
func (h *Handler) UpdateTAFNDIShareFundFeeRdm(ctx context.Context, req *fndv1.TAFNDIShareFundFeeRdmRequest) (*fndv1.SaveResponse, error) {
	h.log.Info("UpdateTAFNDIShareFundFeeRdm",
		slog.String("sys_co_id", req.GetSysCoId()),
		slog.String("fund_code", req.GetFundCode()),
	)
	return h.svc.UpdateTAFNDIShareFundFeeRdm(ctx, req)
}

// DeleteTAFNDIShareFundFeeRdm - DELETE endpoint
func (h *Handler) DeleteTAFNDIShareFundFeeRdm(ctx context.Context, req *fndv1.DeleteRequest) (*fndv1.SaveResponse, error) {
	h.log.Info("DeleteTAFNDIShareFundFeeRdm",
		slog.String("data_id", req.GetDataId()),
		slog.String("data_flag", req.GetDataFlag()),
	)
	return h.svc.DeleteTAFNDIShareFundFeeRdm(ctx, req)
}

// ApproveTAFNDIShareFundFeeRdm handles the Approval POST endpoint.
func (h *Handler) ApproveTAFNDIShareFundFeeRdm(ctx context.Context, req *fndv1.ApproveTAFNDIShareFundFeeRdmRequest) (*fndv1.SaveResponse, error) {
	h.log.Info("ApproveTAFNDIShareFundFeeRdm called", slog.String("data_id", req.GetDataId()))
	return h.svc.ApproveTAFNDIShareFundFeeRdm(ctx, req)
}
