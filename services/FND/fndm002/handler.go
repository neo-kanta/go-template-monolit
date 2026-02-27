package fndm002

import (
	"context"
	"log/slog"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
)

// Handler implements the FNDM002 gRPC service interface.
type Handler struct {
	fndv1.UnimplementedFNDM002ServiceServer
	log *slog.Logger
	svc *Service
}

// NewHandler creates a new FNDM002 gRPC handler.
func NewHandler(log *slog.Logger, svc *Service) *Handler {
	return &Handler{
		log: log.With(slog.String("module", "fndm002")),
		svc: svc,
	}
}

// Service exposes the underlying business service.
func (h *Handler) Service() *Service { return h.svc }

// TAFNDFundFee - GET query endpoint
func (h *Handler) TAFNDFundFee(ctx context.Context, req *fndv1.TAFNDFundFeeRequest) (*fndv1.TAFNDFundFeeResponse, error) {
	h.log.Info("TAFNDFundFee",
		slog.String("sys_co_id", req.GetSysCoId()),
		slog.String("prt_fund_code", req.GetPrtFundCode()),
	)
	return h.svc.TAFNDFundFee(ctx, req)
}

// SaveTAFNDFundFee - POST save (add) endpoint
func (h *Handler) SaveTAFNDFundFee(ctx context.Context, req *fndv1.TAFNDFundFeeRequest) (*fndv1.SaveResponse, error) {
	h.log.Info("SaveTAFNDFundFee",
		slog.String("sys_co_id", req.GetSysCoId()),
		slog.String("prt_fund_code", req.GetPrtFundCode()),
	)
	return h.svc.SaveTAFNDFundFee(ctx, req)
}

// UpdateTAFNDFundFee - PUT save (modify) endpoint
func (h *Handler) UpdateTAFNDFundFee(ctx context.Context, req *fndv1.TAFNDFundFeeRequest) (*fndv1.SaveResponse, error) {
	h.log.Info("UpdateTAFNDFundFee",
		slog.String("sys_co_id", req.GetSysCoId()),
		slog.String("prt_fund_code", req.GetPrtFundCode()),
	)
	return h.svc.UpdateTAFNDFundFee(ctx, req)
}

// DeleteTAFNDFundFee - DELETE endpoint
func (h *Handler) DeleteTAFNDFundFee(ctx context.Context, req *fndv1.DeleteRequest) (*fndv1.SaveResponse, error) {
	h.log.Info("DeleteTAFNDFundFee",
		slog.String("data_id", req.GetDataId()),
		slog.String("data_flag", req.GetDataFlag()),
	)
	return h.svc.DeleteTAFNDFundFee(ctx, req)
}

// ApproveTAFNDFundFee handles the Approval POST endpoint.
func (h *Handler) ApproveTAFNDFundFee(ctx context.Context, req *fndv1.ApproveTAFNDFundFeeRequest) (*fndv1.SaveResponse, error) {
	h.log.Info("ApproveTAFNDFundFee called", slog.String("data_id", req.GetDataId()))
	return h.svc.ApproveTAFNDFundFee(ctx, req)
}
