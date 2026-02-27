package fndm008

import (
	"context"
	"log/slog"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
)

// Handler implements the FNDM008 gRPC service interface.
type Handler struct {
	fndv1.UnimplementedFNDM008ServiceServer
	log *slog.Logger
	svc *Service
}

// NewHandler creates a new FNDM008 gRPC handler.
func NewHandler(log *slog.Logger, svc *Service) *Handler {
	return &Handler{
		log: log.With(slog.String("module", "fndm008")),
		svc: svc,
	}
}

// Service exposes the underlying business service.
func (h *Handler) Service() *Service { return h.svc }

// TAFNDPauseTxn - GET query endpoint
func (h *Handler) TAFNDPauseTxn(ctx context.Context, req *fndv1.TAFNDPauseTxnRequest) (*fndv1.TAFNDPauseTxnResponse, error) {
	h.log.Info("TAFNDPauseTxn",
		slog.String("sys_co_id", req.GetSysCoId()),
		slog.String("prt_fund_code", req.GetPrtFundCode()),
	)
	return h.svc.TAFNDPauseTxn(ctx, req)
}

// SaveTAFNDPauseTxn - POST save (add) endpoint
func (h *Handler) SaveTAFNDPauseTxn(ctx context.Context, req *fndv1.TAFNDPauseTxnRequest) (*fndv1.SaveResponse, error) {
	h.log.Info("SaveTAFNDPauseTxn",
		slog.String("sys_co_id", req.GetSysCoId()),
		slog.String("prt_fund_code", req.GetPrtFundCode()),
	)
	return h.svc.SaveTAFNDPauseTxn(ctx, req)
}

// UpdateTAFNDPauseTxn - PUT save (modify) endpoint
func (h *Handler) UpdateTAFNDPauseTxn(ctx context.Context, req *fndv1.TAFNDPauseTxnRequest) (*fndv1.SaveResponse, error) {
	h.log.Info("UpdateTAFNDPauseTxn",
		slog.String("sys_co_id", req.GetSysCoId()),
		slog.String("prt_fund_code", req.GetPrtFundCode()),
	)
	return h.svc.UpdateTAFNDPauseTxn(ctx, req)
}

// DeleteTAFNDPauseTxn - DELETE endpoint
func (h *Handler) DeleteTAFNDPauseTxn(ctx context.Context, req *fndv1.DeleteRequest) (*fndv1.SaveResponse, error) {
	h.log.Info("DeleteTAFNDPauseTxn",
		slog.String("data_id", req.GetDataId()),
		slog.String("data_flag", req.GetDataFlag()),
	)
	return h.svc.DeleteTAFNDPauseTxn(ctx, req)
}

// TACKPTxnBegDate - KFNDM00801: validate pause txn begin date
func (h *Handler) TACKPTxnBegDate(ctx context.Context, req *fndv1.TACKPTxnBegDateRequest) (*fndv1.CheckResponse, error) {
	h.log.Info("TACKPTxnBegDate",
		slog.String("sys_co_id", req.GetSysCoId()),
		slog.String("prt_fund_code", req.GetPrtFundCode()),
		slog.String("p_txn_beg_date", req.GetPTxnBegDate()),
	)
	return h.svc.TACKPTxnBegDate(ctx, req)
}

// ApproveTAFNDPauseTxn handles the Approval POST endpoint.
func (h *Handler) ApproveTAFNDPauseTxn(ctx context.Context, req *fndv1.ApproveTAFNDPauseTxnRequest) (*fndv1.SaveResponse, error) {
	h.log.Info("ApproveTAFNDPauseTxn called", slog.String("data_id", req.GetDataId()))
	return h.svc.ApproveTAFNDPauseTxn(ctx, req)
}
