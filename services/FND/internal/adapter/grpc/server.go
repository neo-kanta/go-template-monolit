package grpc

import (
	"context"
	"log/slog"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
	"go-transfer-agent/common/platform/config"
	"go-transfer-agent/services/fnd/fndm001"
	"go-transfer-agent/services/fnd/fndm002"
	"go-transfer-agent/services/fnd/fndm003"
	"go-transfer-agent/services/fnd/fndm004"
	"go-transfer-agent/services/fnd/fndm005"
	"go-transfer-agent/services/fnd/fndm006"
	"go-transfer-agent/services/fnd/fndm007"
	"go-transfer-agent/services/fnd/fndm008"
	"go-transfer-agent/services/fnd/fndm009"
	"go-transfer-agent/services/fnd/fndm010"
	"go-transfer-agent/services/fnd/fndm011"
	"go-transfer-agent/services/fnd/fndm012"
	"go-transfer-agent/services/fnd/fndm013"
	"go-transfer-agent/services/fnd/internal/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// fndServer implements fndv1.FNDServiceServer by delegating to per-module handlers.
type fndServer struct {
	fndv1.UnimplementedFNDServiceServer
	log *slog.Logger

	// ─── Module handlers ───
	m001 *fndm001.Handler
}

// ═══════════════════════════════════════════════════════════════════
// FNDM001 delegates
// ═══════════════════════════════════════════════════════════════════

func (s *fndServer) QueryFundInfo(ctx context.Context, req *fndv1.QueryFundInfoRequest) (*fndv1.QueryFundInfoResponse, error) {
	return s.m001.QueryFundInfo(ctx, req)
}

func (s *fndServer) GetDataByDataID(ctx context.Context, req *fndv1.GetDataRequest) (*fndv1.QueryFundInfoByDataIDResponse, error) {
	return s.m001.GetDataByDataID(ctx, req)
}

func (s *fndServer) SaveFundInfo(ctx context.Context, req *fndv1.SaveFundInfoRequest) (*fndv1.SaveFundInfoResponse, error) {
	return s.m001.SaveFundInfo(ctx, req)
}

func (s *fndServer) UpdateFundInfo(ctx context.Context, req *fndv1.SaveFundInfoRequest) (*fndv1.SaveFundInfoResponse, error) {
	return s.m001.UpdateFundInfo(ctx, req)
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
// Server Factory
// ═══════════════════════════════════════════════════════════════════

// NewServer creates a gRPC server and registers all FND modules.
func NewServer(cfg *config.Config, log *slog.Logger) *grpc.Server {
	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			AuthInterceptor(cfg),
			THBValidationInterceptor(),
			ErrorInterceptor(log),
		),
	)

	// ─── Database (optional — modules degrade gracefully without DB) ───
	var db *gorm.DB
	if cfg.DBDsn != "" {
		var err error
		db, err = gorm.Open(postgres.Open(cfg.DBDsn), &gorm.Config{})
		if err != nil {
			log.Warn("⚠️  Database connection failed — modules will use stub data",
				slog.String("error", err.Error()),
				slog.String("suggestion", "Set DB_DSN environment variable to a valid PostgreSQL DSN"),
			)
		} else {
			log.Info("✅ Database connected successfully")

			// Auto-migrate auth tables and seed POC users
			if err := auth.AutoMigrateAndSeed(db, log); err != nil {
				log.Warn("⚠️  Auth table migration failed — user validation disabled",
					slog.String("error", err.Error()),
				)
			}
		}
	} else {
		log.Info("⚠️  No DB_DSN configured — modules with DB dependencies will return empty results")
	}

	// ─── Wire FNDM002-013 (GORM-backed, require DB) ───
	if db != nil {
		fndSrv := &fndServer{
			log:  log,
			m001: fndm001.NewHandler(log, fndm001.NewService(db, log)),
		}
		fndv1.RegisterFNDServiceServer(srv, fndSrv)
		fndv1.RegisterFNDM002ServiceServer(srv, fndm002.NewHandler(log, fndm002.NewService(db, log)))
		fndv1.RegisterFNDM003ServiceServer(srv, fndm003.NewHandler(log, fndm003.NewService(db, log)))
		fndv1.RegisterFNDM004ServiceServer(srv, fndm004.NewHandler(log, fndm004.NewService(db, log)))
		fndv1.RegisterFNDM005ServiceServer(srv, fndm005.NewHandler(log, fndm005.NewService(db, log)))
		fndv1.RegisterFNDM006ServiceServer(srv, fndm006.NewHandler(log, fndm006.NewService(db, log)))
		fndv1.RegisterFNDM007ServiceServer(srv, fndm007.NewHandler(log, fndm007.NewService(db, log)))
		fndv1.RegisterFNDM008ServiceServer(srv, fndm008.NewHandler(log, fndm008.NewService(db, log)))
		fndv1.RegisterFNDM009ServiceServer(srv, fndm009.NewHandler(log, fndm009.NewService(db, log)))
		fndv1.RegisterFNDM010ServiceServer(srv, fndm010.NewHandler(log, fndm010.NewService(db, log)))
		fndv1.RegisterFNDM011ServiceServer(srv, fndm011.NewHandler(log, fndm011.NewService(db, log)))
		fndv1.RegisterFNDM012ServiceServer(srv, fndm012.NewHandler(log, fndm012.NewService(db, log)))
		fndv1.RegisterFNDM013ServiceServer(srv, fndm013.NewHandler(log, fndm013.NewService(db, log)))
		log.Info("registered fnd.v1 services",
			slog.Int("modules", 13),
		)
	} else {
		log.Info("registered fnd.v1 services (FNDM001 only — no DB)",
			slog.Int("modules", 1),
		)
	}

	// ─── Health + reflection ───
	healthSrv := health.NewServer()
	healthSrv.SetServingStatus("fnd.v1.FNDService", healthpb.HealthCheckResponse_SERVING)
	healthSrv.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(srv, healthSrv)

	reflection.Register(srv)

	return srv
}
