// Package fndm001 — gRPC handler for APIFNDM001 (TAFNDFundInfo).
//
// Method naming follows the C# spec exactly:
//
//	AUD:   SaveFundInfo  (Post/Put)
//	Query: QueryFundInfo, QueryFundInfoByDataID
//	KFNDM: TACKCSDPrtFundSetlDate, TACKFundCodeExist
//	GFNDM: TAGetDPrtFundIssueCry, TAGetDTxCry
package fndm001

import (
	"context"
	"log/slog"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"

	"google.golang.org/grpc"
)

// Handler implements the FNDM001 gRPC methods.
// It is NOT a standalone gRPC service — the central fndServer (adapter/grpc) delegates to it.
type Handler struct {
	log *slog.Logger
	svc *Service
}

// NewHandler creates a new FNDM001 gRPC handler.
func NewHandler(log *slog.Logger, svc *Service) *Handler {
	return &Handler{
		log: log.With(slog.String("module", "fndm001")),
		svc: svc,
	}
}

// Service exposes the underlying business service (for cross-module queries).
func (h *Handler) Service() *Service { return h.svc }

// ─── APIFNDM001: AUD / Query ───

func (h *Handler) QueryFundInfo(ctx context.Context, req *fndv1.QueryFundInfoRequest) (*fndv1.QueryFundInfoResponse, error) {
	h.log.Info("QueryFundInfo", slog.String("sys_co_id", req.GetSysCoId()))
	return h.svc.QueryFundInfo(req.GetSysCoId()), nil
}

func (h *Handler) GetDataByDataID(ctx context.Context, req *fndv1.GetDataRequest) (*fndv1.QueryFundInfoByDataIDResponse, error) {
	h.log.Info("GetDataByDataID", slog.String("data_id", req.GetDataId()))
	return h.svc.GetDataByDataID(ctx, req)
}

func (h *Handler) SaveFundInfo(ctx context.Context, req *fndv1.SaveFundInfoRequest) (*fndv1.SaveFundInfoResponse, error) {
	h.log.Info("SaveFundInfo",
		slog.String("sys_co_id", req.GetMaster().GetSysCoId()),
		slog.String("prt_fund_code", req.GetMaster().GetPrtFundCode()),
	)
	return h.svc.SaveFundInfo(ctx, req), nil
}

func (h *Handler) UpdateFundInfo(ctx context.Context, req *fndv1.SaveFundInfoRequest) (*fndv1.SaveFundInfoResponse, error) {
	h.log.Info("UpdateFundInfo",
		slog.String("sys_co_id", req.GetMaster().GetSysCoId()),
		slog.String("prt_fund_code", req.GetMaster().GetPrtFundCode()),
	)
	return h.svc.UpdateFundInfo(req), nil
}

func (h *Handler) ApproveFundInfo(ctx context.Context, req *fndv1.ApproveFundInfoRequest) (*fndv1.SaveFundInfoResponse, error) {
	h.log.Info("ApproveFundInfo",
		slog.String("sys_co_id", req.GetSysCoId()),
		slog.String("prt_fund_code", req.GetPrtFundCode()),
		slog.String("checker_id", req.GetCheckerId()),
	)
	return h.svc.ApproveFundInfo(ctx, req), nil
}

func (h *Handler) DeleteTAFNDFundInfo(ctx context.Context, req *fndv1.DeleteRequest) (*fndv1.SaveResponse, error) {
	h.log.Info("DeleteTAFNDFundInfo",
		slog.String("data_id", req.GetDataId()),
		slog.String("data_flag", req.GetDataFlag()),
	)
	return h.svc.DeleteTAFNDFundInfo(ctx, req)
}

// ─── KFNDM: Check Functions ───

// TACKCSDPrtFundSetlDate — KFNDM00101: ChkData/TACKCSDPrtFundSetlDate
func (h *Handler) TACKCSDPrtFundSetlDate(ctx context.Context, req *fndv1.TACKCSDPrtFundSetlDateRequest) (*fndv1.CheckResponse, error) {
	h.log.Info("TACKCSDPrtFundSetlDate",
		slog.String("sys_co_id", req.GetSysCoId()),
		slog.String("csd_prt_fund_setl_date", req.GetCsdPrtFundSetlDate()),
	)
	return h.svc.TACKCSDPrtFundSetlDate(req.GetSysCoId(), req.GetCsdPrtFundSetlDate(), req.FundCode), nil
}

// TACKFundCodeExist — KFNDM00103: ChkData/TACKFundCodeExist
func (h *Handler) TACKFundCodeExist(ctx context.Context, req *fndv1.TACKFundCodeExistRequest) (*fndv1.CheckResponse, error) {
	h.log.Info("TACKFundCodeExist",
		slog.String("sys_co_id", req.GetSysCoId()),
		slog.String("fund_code", req.GetFundCode()),
	)
	return h.svc.TACKFundCodeExist(req.GetSysCoId(), req.GetFundCode()), nil
}

// ─── GFNDM: Get/Search Functions ───

// TAGetDPrtFundIssueCry — GFNDM00101: GetInfo/TAGetDPrtFundIssueCry
func (h *Handler) TAGetDPrtFundIssueCry(ctx context.Context, req *fndv1.TAGetDPrtFundIssueCryRequest) (*fndv1.TAGetDPrtFundIssueCryResponse, error) {
	h.log.Info("TAGetDPrtFundIssueCry", slog.String("sys_co_id", req.GetSysCoId()))
	return h.svc.TAGetDPrtFundIssueCry(
		req.GetSysCoId(), req.GetIsAll(), req.GetPrtFundCode(),
		req.GetFeeChargeType(), req.GetFilterItem(), req.GetQueryType(),
	), nil
}

// TAGetDTxCry — GFNDM00104: GetInfo/TAGetDTxCry
func (h *Handler) TAGetDTxCry(ctx context.Context, req *fndv1.TAGetDTxCryRequest) (*fndv1.TAGetDTxCryResponse, error) {
	h.log.Info("TAGetDTxCry", slog.String("sys_co_id", req.GetSysCoId()))
	return h.svc.TAGetDTxCry(
		req.GetSysCoId(), req.GetIsAll(), req.GetPrtFundCode(), req.GetCryId(),
		req.GetCryDataSrc(), req.GetIsExcludeCcy(), req.GetQueryType(),
	), nil
}

// ─── Registration (placeholder for future per-module proto split) ───

func RegisterHandlers(_ *grpc.Server) {
	// Future: fndv1.RegisterTAFNDFundInfoServiceServer(srv, NewHandler(log))
}
