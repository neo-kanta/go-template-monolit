package fndm013

import (
	"context"
	"log/slog"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
)

// Handler implements the FNDM013 gRPC service interface.
type Handler struct {
	fndv1.UnimplementedFNDM013ServiceServer
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

// TAFNDTMFundFeeRdm - GET query endpoint
func (h *Handler) TAFNDTMFundFeeRdm(ctx context.Context, req *fndv1.TAFNDTMFundFeeRdmRequest) (*fndv1.TAFNDTMFundFeeRdmResponse, error) {
	h.log.Info("TAFNDTMFundFeeRdm",
		slog.String("sys_co_id", req.GetSysCoId()),
		slog.String("prt_fund_code", req.GetPrtFundCode()),
	)
	return h.svc.TAFNDTMFundFeeRdm(ctx, req)
}

// SaveTAFNDTMFundFeeRdm - POST save (add) endpoint
func (h *Handler) SaveTAFNDTMFundFeeRdm(ctx context.Context, req *fndv1.TAFNDTMFundFeeRdmRequest) (*fndv1.SaveResponse, error) {
	h.log.Info("SaveTAFNDTMFundFeeRdm",
		slog.String("sys_co_id", req.GetSysCoId()),
		slog.String("prt_fund_code", req.GetPrtFundCode()),
	)
	return h.svc.SaveTAFNDTMFundFeeRdm(ctx, req)
}

// UpdateTAFNDTMFundFeeRdm - PUT save (modify) endpoint
func (h *Handler) UpdateTAFNDTMFundFeeRdm(ctx context.Context, req *fndv1.TAFNDTMFundFeeRdmRequest) (*fndv1.SaveResponse, error) {
	h.log.Info("UpdateTAFNDTMFundFeeRdm",
		slog.String("sys_co_id", req.GetSysCoId()),
		slog.String("prt_fund_code", req.GetPrtFundCode()),
	)
	return h.svc.UpdateTAFNDTMFundFeeRdm(ctx, req)
}

// DeleteTAFNDTMFundFeeRdm - DELETE endpoint
func (h *Handler) DeleteTAFNDTMFundFeeRdm(ctx context.Context, req *fndv1.DeleteRequest) (*fndv1.SaveResponse, error) {
	h.log.Info("DeleteTAFNDTMFundFeeRdm",
		slog.String("data_id", req.GetDataId()),
		slog.String("data_flag", req.GetDataFlag()),
	)
	return h.svc.DeleteTAFNDTMFundFeeRdm(ctx, req)
}

// ApproveTAFNDTMFundFeeRdm handles the Approval POST endpoint.
func (h *Handler) ApproveTAFNDTMFundFeeRdm(ctx context.Context, req *fndv1.ApproveTAFNDTMFundFeeRdmRequest) (*fndv1.SaveResponse, error) {
	h.log.Info("ApproveTAFNDTMFundFeeRdm called", slog.String("data_id", req.GetDataId()))
	return h.svc.ApproveTAFNDTMFundFeeRdm(ctx, req)
}
