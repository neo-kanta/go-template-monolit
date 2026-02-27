package fndm005

import (
	"context"
	"log/slog"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
)

// Handler implements the FNDM005 gRPC service interface.
type Handler struct {
	fndv1.UnimplementedFNDM005ServiceServer
	log *slog.Logger
	svc *Service
}

// NewHandler creates a new FNDM005 gRPC handler.
func NewHandler(log *slog.Logger, svc *Service) *Handler {
	return &Handler{
		log: log.With(slog.String("module", "fndm005")),
		svc: svc,
	}
}

// Service exposes the underlying business service.
func (h *Handler) Service() *Service { return h.svc }

// TAFNDFundCalDate - GET query endpoint
func (h *Handler) TAFNDFundCalDate(ctx context.Context, req *fndv1.TAFNDFundCalDateRequest) (*fndv1.TAFNDFundCalDateResponse, error) {
	h.log.Info("TAFNDFundCalDate",
		slog.String("sys_co_id", req.GetSysCoId()),
		slog.String("cal_year", req.GetCalYear()),
	)
	return h.svc.TAFNDFundCalDate(ctx, req)
}

// SaveTAFNDFundCalDate - POST save (add) endpoint
func (h *Handler) SaveTAFNDFundCalDate(ctx context.Context, req *fndv1.TAFNDFundCalDateRequest) (*fndv1.SaveResponse, error) {
	h.log.Info("SaveTAFNDFundCalDate",
		slog.String("sys_co_id", req.GetSysCoId()),
		slog.String("cal_year", req.GetCalYear()),
	)
	return h.svc.SaveTAFNDFundCalDate(ctx, req)
}

// UpdateTAFNDFundCalDate - PUT save (modify) endpoint
func (h *Handler) UpdateTAFNDFundCalDate(ctx context.Context, req *fndv1.TAFNDFundCalDateRequest) (*fndv1.SaveResponse, error) {
	h.log.Info("UpdateTAFNDFundCalDate",
		slog.String("sys_co_id", req.GetSysCoId()),
		slog.String("cal_year", req.GetCalYear()),
	)
	return h.svc.UpdateTAFNDFundCalDate(ctx, req)
}

// DeleteTAFNDFundCalDate - DELETE endpoint
func (h *Handler) DeleteTAFNDFundCalDate(ctx context.Context, req *fndv1.DeleteRequest) (*fndv1.SaveResponse, error) {
	h.log.Info("DeleteTAFNDFundCalDate",
		slog.String("data_id", req.GetDataId()),
		slog.String("data_flag", req.GetDataFlag()),
	)
	return h.svc.DeleteTAFNDFundCalDate(ctx, req)
}

// TACKFNDCalEdit - KFNDM00502: check for pending calendar edit data
func (h *Handler) TACKFNDCalEdit(ctx context.Context, req *fndv1.TACKFNDCalEditRequest) (*fndv1.CheckResponse, error) {
	h.log.Info("TACKFNDCalEdit",
		slog.String("sys_co_id", req.GetSysCoId()),
		slog.String("check_time_point", req.GetCheckTimePoint()),
	)
	return h.svc.TACKFNDCalEdit(ctx, req)
}

// ApproveTAFNDFundCalDate handles the Approval POST endpoint.
func (h *Handler) ApproveTAFNDFundCalDate(ctx context.Context, req *fndv1.ApproveTAFNDFundCalDateRequest) (*fndv1.SaveResponse, error) {
	h.log.Info("ApproveTAFNDFundCalDate called", slog.String("data_id", req.GetDataId()))
	return h.svc.ApproveTAFNDFundCalDate(ctx, req)
}
