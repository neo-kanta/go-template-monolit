package fndm007

import (
	"context"
	"log/slog"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
)

// Handler implements the FNDM007 gRPC service interface.
type Handler struct {
	fndv1.UnimplementedFNDM007ServiceServer
	log *slog.Logger
	svc *Service
}

// NewHandler creates a new FNDM007 gRPC handler.
func NewHandler(log *slog.Logger, svc *Service) *Handler {
	return &Handler{
		log: log.With(slog.String("module", "fndm007")),
		svc: svc,
	}
}

// Service exposes the underlying business service.
func (h *Handler) Service() *Service { return h.svc }

// TAFNDCustGroup - GET query endpoint
func (h *Handler) TAFNDCustGroup(ctx context.Context, req *fndv1.TAFNDCustGroupRequest) (*fndv1.TAFNDCustGroupResponse, error) {
	h.log.Info("TAFNDCustGroup",
		slog.String("sys_co_id", req.GetSysCoId()),
		slog.String("cust_grp_code", req.GetCustGrpCode()),
	)
	return h.svc.TAFNDCustGroup(ctx, req)
}

// SaveTAFNDCustGroup - POST save (add) endpoint
func (h *Handler) SaveTAFNDCustGroup(ctx context.Context, req *fndv1.TAFNDCustGroupRequest) (*fndv1.SaveResponse, error) {
	h.log.Info("SaveTAFNDCustGroup",
		slog.String("sys_co_id", req.GetSysCoId()),
		slog.String("cust_grp_code", req.GetCustGrpCode()),
	)
	return h.svc.SaveTAFNDCustGroup(ctx, req)
}

// UpdateTAFNDCustGroup - PUT save (modify) endpoint
func (h *Handler) UpdateTAFNDCustGroup(ctx context.Context, req *fndv1.TAFNDCustGroupRequest) (*fndv1.SaveResponse, error) {
	h.log.Info("UpdateTAFNDCustGroup",
		slog.String("sys_co_id", req.GetSysCoId()),
		slog.String("cust_grp_code", req.GetCustGrpCode()),
	)
	return h.svc.UpdateTAFNDCustGroup(ctx, req)
}

// DeleteTAFNDCustGroup - DELETE endpoint
func (h *Handler) DeleteTAFNDCustGroup(ctx context.Context, req *fndv1.DeleteRequest) (*fndv1.SaveResponse, error) {
	h.log.Info("DeleteTAFNDCustGroup",
		slog.String("data_id", req.GetDataId()),
		slog.String("data_flag", req.GetDataFlag()),
	)
	return h.svc.DeleteTAFNDCustGroup(ctx, req)
}

// TAGetCustGroup - GFNDM00701: get custom fund group codes (for Searcher)
func (h *Handler) TAGetCustGroup(ctx context.Context, req *fndv1.TAGetCustGroupRequest) (*fndv1.TAGetCustGroupResponse, error) {
	h.log.Info("TAGetCustGroup",
		slog.String("sys_co_id", req.GetSysCoId()),
		slog.String("cust_grp_code", req.GetCustGrpCode()),
	)
	return h.svc.TAGetCustGroup(ctx, req)
}

// ApproveTAFNDCustGroup handles the Approval POST endpoint.
func (h *Handler) ApproveTAFNDCustGroup(ctx context.Context, req *fndv1.ApproveTAFNDCustGroupRequest) (*fndv1.SaveResponse, error) {
	h.log.Info("ApproveTAFNDCustGroup called", slog.String("data_id", req.GetDataId()))
	return h.svc.ApproveTAFNDCustGroup(ctx, req)
}
