package fndm003

import (
	"context"
	"log/slog"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
)

// Handler implements the FNDM003 gRPC service interface.
type Handler struct {
	fndv1.UnimplementedFNDM003ServiceServer
	log *slog.Logger
	svc *Service
}

// NewHandler creates a new FNDM003 gRPC handler.
func NewHandler(log *slog.Logger, svc *Service) *Handler {
	return &Handler{
		log: log.With(slog.String("module", "fndm003")),
		svc: svc,
	}
}

// Service exposes the underlying business service.
func (h *Handler) Service() *Service { return h.svc }

// TAFNDSwitch - GET query endpoint
func (h *Handler) TAFNDSwitch(ctx context.Context, req *fndv1.TAFNDSwitchRequest) (*fndv1.TAFNDSwitchResponse, error) {
	h.log.Info("TAFNDSwitch",
		slog.String("sys_co_id", req.GetSysCoId()),
		slog.String("prt_fund_code", req.GetPrtFundCode()),
	)
	return h.svc.TAFNDSwitch(ctx, req)
}

// SaveTAFNDSwitch - POST save (add) endpoint
func (h *Handler) SaveTAFNDSwitch(ctx context.Context, req *fndv1.TAFNDSwitchRequest) (*fndv1.SaveResponse, error) {
	h.log.Info("SaveTAFNDSwitch",
		slog.String("sys_co_id", req.GetSysCoId()),
		slog.String("prt_fund_code", req.GetPrtFundCode()),
	)
	return h.svc.SaveTAFNDSwitch(ctx, req)
}

// UpdateTAFNDSwitch - PUT save (modify) endpoint
func (h *Handler) UpdateTAFNDSwitch(ctx context.Context, req *fndv1.TAFNDSwitchRequest) (*fndv1.SaveResponse, error) {
	h.log.Info("UpdateTAFNDSwitch",
		slog.String("sys_co_id", req.GetSysCoId()),
		slog.String("prt_fund_code", req.GetPrtFundCode()),
	)
	return h.svc.UpdateTAFNDSwitch(ctx, req)
}

// DeleteTAFNDSwitch - DELETE endpoint
func (h *Handler) DeleteTAFNDSwitch(ctx context.Context, req *fndv1.DeleteRequest) (*fndv1.SaveResponse, error) {
	h.log.Info("DeleteTAFNDSwitch",
		slog.String("data_id", req.GetDataId()),
		slog.String("data_flag", req.GetDataFlag()),
	)
	return h.svc.DeleteTAFNDSwitch(ctx, req)
}

// ApproveTAFNDSwitch handles the Approval POST endpoint.
func (h *Handler) ApproveTAFNDSwitch(ctx context.Context, req *fndv1.ApproveTAFNDSwitchRequest) (*fndv1.SaveResponse, error) {
	h.log.Info("ApproveTAFNDSwitch called", slog.String("data_id", req.GetDataId()))
	return h.svc.ApproveTAFNDSwitch(ctx, req)
}
