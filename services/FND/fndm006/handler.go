package fndm006

import (
	"context"
	"log/slog"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
)

// Handler implements the FNDM006 gRPC service interface.
type Handler struct {
	fndv1.UnimplementedFNDM006ServiceServer
	log *slog.Logger
	svc *Service
}

// NewHandler creates a new FNDM006 gRPC handler.
func NewHandler(log *slog.Logger, svc *Service) *Handler {
	return &Handler{
		log: log.With(slog.String("module", "fndm006")),
		svc: svc,
	}
}

// Service exposes the underlying business service.
func (h *Handler) Service() *Service { return h.svc }

// TAFNDRPFeeChgType - GET query endpoint
func (h *Handler) TAFNDRPFeeChgType(ctx context.Context, req *fndv1.TAFNDRPFeeChgTypeRequest) (*fndv1.TAFNDRPFeeChgTypeResponse, error) {
	h.log.Info("TAFNDRPFeeChgType",
		slog.String("sys_co_id", req.GetSysCoId()),
		slog.String("prt_fund_code", req.GetPrtFundCode()),
	)
	return h.svc.TAFNDRPFeeChgType(ctx, req)
}

// SaveTAFNDRPFeeChgType - POST save (add) endpoint
func (h *Handler) SaveTAFNDRPFeeChgType(ctx context.Context, req *fndv1.TAFNDRPFeeChgTypeRequest) (*fndv1.SaveResponse, error) {
	h.log.Info("SaveTAFNDRPFeeChgType",
		slog.String("sys_co_id", req.GetSysCoId()),
		slog.String("prt_fund_code", req.GetPrtFundCode()),
	)
	return h.svc.SaveTAFNDRPFeeChgType(ctx, req)
}

// UpdateTAFNDRPFeeChgType - PUT save (modify) endpoint
func (h *Handler) UpdateTAFNDRPFeeChgType(ctx context.Context, req *fndv1.TAFNDRPFeeChgTypeRequest) (*fndv1.SaveResponse, error) {
	h.log.Info("UpdateTAFNDRPFeeChgType",
		slog.String("sys_co_id", req.GetSysCoId()),
		slog.String("prt_fund_code", req.GetPrtFundCode()),
	)
	return h.svc.UpdateTAFNDRPFeeChgType(ctx, req)
}

// DeleteTAFNDRPFeeChgType - DELETE endpoint
func (h *Handler) DeleteTAFNDRPFeeChgType(ctx context.Context, req *fndv1.DeleteRequest) (*fndv1.SaveResponse, error) {
	h.log.Info("DeleteTAFNDRPFeeChgType",
		slog.String("data_id", req.GetDataId()),
		slog.String("data_flag", req.GetDataFlag()),
	)
	return h.svc.DeleteTAFNDRPFeeChgType(ctx, req)
}

// ApproveTAFNDRPFeeChgType handles the Approval POST endpoint.
func (h *Handler) ApproveTAFNDRPFeeChgType(ctx context.Context, req *fndv1.ApproveTAFNDRPFeeChgTypeRequest) (*fndv1.SaveResponse, error) {
	h.log.Info("ApproveTAFNDRPFeeChgType called", slog.String("data_id", req.GetDataId()))
	return h.svc.ApproveTAFNDRPFeeChgType(ctx, req)
}
