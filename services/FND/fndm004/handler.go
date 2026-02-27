package fndm004

import (
	"context"
	"log/slog"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
)

// Handler implements the FNDM004 gRPC service interface.
type Handler struct {
	fndv1.UnimplementedFNDM004ServiceServer
	log *slog.Logger
	svc *Service
}

// NewHandler creates a new FNDM004 gRPC handler.
func NewHandler(log *slog.Logger, svc *Service) *Handler {
	return &Handler{
		log: log.With(slog.String("module", "fndm004")),
		svc: svc,
	}
}

// Service exposes the underlying business service.
func (h *Handler) Service() *Service { return h.svc }

// TAFNDFundAgent - GET query endpoint
func (h *Handler) TAFNDFundAgent(ctx context.Context, req *fndv1.TAFNDFundAgentRequest) (*fndv1.TAFNDFundAgentResponse, error) {
	h.log.Info("TAFNDFundAgent",
		slog.String("sys_co_id", req.GetSysCoId()),
		slog.String("prt_fund_code", req.GetPrtFundCode()),
	)
	return h.svc.TAFNDFundAgent(ctx, req)
}

// SaveTAFNDFundAgent - POST save (add) endpoint
func (h *Handler) SaveTAFNDFundAgent(ctx context.Context, req *fndv1.TAFNDFundAgentRequest) (*fndv1.SaveResponse, error) {
	h.log.Info("SaveTAFNDFundAgent",
		slog.String("sys_co_id", req.GetSysCoId()),
		slog.String("prt_fund_code", req.GetPrtFundCode()),
	)
	return h.svc.SaveTAFNDFundAgent(ctx, req)
}

// UpdateTAFNDFundAgent - PUT save (modify) endpoint
func (h *Handler) UpdateTAFNDFundAgent(ctx context.Context, req *fndv1.TAFNDFundAgentRequest) (*fndv1.SaveResponse, error) {
	h.log.Info("UpdateTAFNDFundAgent",
		slog.String("sys_co_id", req.GetSysCoId()),
		slog.String("prt_fund_code", req.GetPrtFundCode()),
	)
	return h.svc.UpdateTAFNDFundAgent(ctx, req)
}

// DeleteTAFNDFundAgent - DELETE endpoint
func (h *Handler) DeleteTAFNDFundAgent(ctx context.Context, req *fndv1.DeleteRequest) (*fndv1.SaveResponse, error) {
	h.log.Info("DeleteTAFNDFundAgent",
		slog.String("data_id", req.GetDataId()),
		slog.String("data_flag", req.GetDataFlag()),
	)
	return h.svc.DeleteTAFNDFundAgent(ctx, req)
}

// TACKFundAgentEdit - KFNDM00402: check for pending edit data
func (h *Handler) TACKFundAgentEdit(ctx context.Context, req *fndv1.TACKFundAgentEditRequest) (*fndv1.CheckResponse, error) {
	h.log.Info("TACKFundAgentEdit",
		slog.String("sys_co_id", req.GetSysCoId()),
		slog.String("check_time_point", req.GetCheckTimePoint()),
	)
	return h.svc.TACKFundAgentEdit(ctx, req)
}

// ApproveTAFNDFundAgent handles the Approval POST endpoint.
func (h *Handler) ApproveTAFNDFundAgent(ctx context.Context, req *fndv1.ApproveTAFNDFundAgentRequest) (*fndv1.SaveResponse, error) {
	h.log.Info("ApproveTAFNDFundAgent called", slog.String("data_id", req.GetDataId()))
	return h.svc.ApproveTAFNDFundAgent(ctx, req)
}
