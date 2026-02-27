package fndm012

import (
	"context"
	"log/slog"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
)

// Handler implements the FNDM012 gRPC service interface.
type Handler struct {
	fndv1.UnimplementedFNDM012ServiceServer
	log *slog.Logger
	svc *Service
}

// NewHandler creates a new FNDM012 gRPC handler.
func NewHandler(log *slog.Logger, svc *Service) *Handler {
	return &Handler{
		log: log.With(slog.String("module", "fndm012")),
		svc: svc,
	}
}

// Service exposes the underlying business service.
func (h *Handler) Service() *Service { return h.svc }

// TAFNDPGFundFeeRdm - GET query endpoint
func (h *Handler) TAFNDPGFundFeeRdm(ctx context.Context, req *fndv1.TAFNDPGFundFeeRdmRequest) (*fndv1.TAFNDPGFundFeeRdmResponse, error) {
	h.log.Info("TAFNDPGFundFeeRdm",
		slog.String("sys_co_id", req.GetSysCoId()),
		slog.String("prt_fund_code", req.GetPrtFundCode()),
	)
	return h.svc.TAFNDPGFundFeeRdm(ctx, req)
}

// SaveTAFNDPGFundFeeRdm - POST save (add) endpoint
func (h *Handler) SaveTAFNDPGFundFeeRdm(ctx context.Context, req *fndv1.TAFNDPGFundFeeRdmRequest) (*fndv1.SaveResponse, error) {
	h.log.Info("SaveTAFNDPGFundFeeRdm",
		slog.String("sys_co_id", req.GetSysCoId()),
		slog.String("prt_fund_code", req.GetPrtFundCode()),
	)
	return h.svc.SaveTAFNDPGFundFeeRdm(ctx, req)
}

// UpdateTAFNDPGFundFeeRdm - PUT save (modify) endpoint
func (h *Handler) UpdateTAFNDPGFundFeeRdm(ctx context.Context, req *fndv1.TAFNDPGFundFeeRdmRequest) (*fndv1.SaveResponse, error) {
	h.log.Info("UpdateTAFNDPGFundFeeRdm",
		slog.String("sys_co_id", req.GetSysCoId()),
		slog.String("prt_fund_code", req.GetPrtFundCode()),
	)
	return h.svc.UpdateTAFNDPGFundFeeRdm(ctx, req)
}

// DeleteTAFNDPGFundFeeRdm - DELETE endpoint
func (h *Handler) DeleteTAFNDPGFundFeeRdm(ctx context.Context, req *fndv1.DeleteRequest) (*fndv1.SaveResponse, error) {
	h.log.Info("DeleteTAFNDPGFundFeeRdm",
		slog.String("data_id", req.GetDataId()),
		slog.String("data_flag", req.GetDataFlag()),
	)
	return h.svc.DeleteTAFNDPGFundFeeRdm(ctx, req)
}

// ApproveTAFNDPGFundFeeRdm handles the Approval POST endpoint.
func (h *Handler) ApproveTAFNDPGFundFeeRdm(ctx context.Context, req *fndv1.ApproveTAFNDPGFundFeeRdmRequest) (*fndv1.SaveResponse, error) {
	h.log.Info("ApproveTAFNDPGFundFeeRdm called", slog.String("data_id", req.GetDataId()))
	return h.svc.ApproveTAFNDPGFundFeeRdm(ctx, req)
}
