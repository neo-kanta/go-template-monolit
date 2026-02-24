package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
	"go-transfer-agent/common/platform/config"
	"go-transfer-agent/services/fnd/fndm001"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// fndServer implements fndv1.FNDServiceServer by delegating to per-module handlers.
//
// When adding a new module (e.g., fndm002):
//  1. Add a field: fndm002 *fndm002.handler
//  2. Init it in NewServer()
//  3. Add delegate methods below
type fndServer struct {
	fndv1.UnimplementedFNDServiceServer
	log *slog.Logger

	// ─── Module handlers ───
	m001 *fndm001.Handler
	// m002 *fndm002.Handler   // ← uncomment when fndm002 is ready
}

// ═══════════════════════════════════════════════════════════════════
// FNDM001 delegates
// ═══════════════════════════════════════════════════════════════════

func (s *fndServer) QueryFundInfo(ctx context.Context, req *fndv1.QueryFundInfoRequest) (*fndv1.QueryFundInfoResponse, error) {
	return s.m001.QueryFundInfo(ctx, req)
}

func (s *fndServer) QueryFundInfoByDataID(ctx context.Context, req *fndv1.QueryFundInfoByDataIDRequest) (*fndv1.QueryFundInfoByDataIDResponse, error) {
	return s.m001.QueryFundInfoByDataID(ctx, req)
}

func (s *fndServer) SaveFundInfo(ctx context.Context, req *fndv1.SaveFundInfoRequest) (*fndv1.SaveFundInfoResponse, error) {
	return s.m001.SaveFundInfo(ctx, req)
}

func (s *fndServer) TACKCSDPrtFundSetlDate(ctx context.Context, req *fndv1.TACKCSDPrtFundSetlDateRequest) (*fndv1.CheckResponse, error) {
	return s.m001.TACKCSDPrtFundSetlDate(ctx, req)
}

func (s *fndServer) TACKFundCodeExist(ctx context.Context, req *fndv1.TACKFundCodeExistRequest) (*fndv1.CheckResponse, error) {
	return s.m001.TACKFundCodeExist(ctx, req)
}

func (s *fndServer) TAGetDPrtFundIssueCry(ctx context.Context, req *fndv1.TAGetDPrtFundIssueCryRequest) (*fndv1.TAGetDPrtFundIssueCryResponse, error) {
	return s.m001.TAGetDPrtFundIssueCry(ctx, req)
}

func (s *fndServer) TAGetDTxCry(ctx context.Context, req *fndv1.TAGetDTxCryRequest) (*fndv1.TAGetDTxCryResponse, error) {
	return s.m001.TAGetDTxCry(ctx, req)
}

// ═══════════════════════════════════════════════════════════════════
// POC stubs (kept for compatibility — will move to their own module)
// ═══════════════════════════════════════════════════════════════════

func (s *fndServer) CreateTransaction(ctx context.Context, req *fndv1.CreateTransactionRequest) (*fndv1.CreateTransactionResponse, error) {
	txnID := fmt.Sprintf("txn_%d", time.Now().UnixMilli())
	s.log.Info("CreateTransaction", slog.String("txn_id", txnID))
	return &fndv1.CreateTransactionResponse{
		TransactionId: txnID, Status: "pending", CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}, nil
}

func (s *fndServer) ListTransactions(ctx context.Context, req *fndv1.ListTransactionsRequest) (*fndv1.ListTransactionsResponse, error) {
	s.log.Info("ListTransactions", slog.Int("page_size", int(req.GetPageSize())))
	return &fndv1.ListTransactionsResponse{
		Transactions: []*fndv1.Transaction{{
			TransactionId: "txn_sample001", FromAccount: "1234567890", ToAccount: "0987654321",
			AmountMinor: 250000, Currency: "THB", Memo: "Sample stub transaction",
			Status: "approved", CreatedBy: "admin", CreatedAt: "2026-02-18T09:00:00Z",
		}},
		Total: 1,
	}, nil
}

// ═══════════════════════════════════════════════════════════════════
// Server Factory
// ═══════════════════════════════════════════════════════════════════
// NewServer creates a gRPC server and registers all FND modules.
func NewServer(cfg *config.Config, log *slog.Logger) *grpc.Server {
	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			AuthInterceptor(cfg.JWTSecret),
			ErrorInterceptor(log),
		),
	)

	// ─── Wire modules ───
	fndSrv := &fndServer{
		log:  log,
		m001: fndm001.NewHandler(log, fndm001.NewService()),
		// New Handler here
	}
	fndv1.RegisterFNDServiceServer(srv, fndSrv)
	log.Info("registered fnd.v1.FNDService",
		slog.Int("modules", 1), // increment when adding modules
	)

	// ─── Health + reflection ───
	healthSrv := health.NewServer()
	healthSrv.SetServingStatus("fnd.v1.FNDService", healthpb.HealthCheckResponse_SERVING)
	healthSrv.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(srv, healthSrv)

	reflection.Register(srv)

	return srv
}
