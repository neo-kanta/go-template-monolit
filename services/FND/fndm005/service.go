package fndm005

import (
	"context"
	"log/slog"
	"time"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
	"go-transfer-agent/services/fnd/fndm005/db"
	"go-transfer-agent/services/fnd/shared"

	"gorm.io/gorm"
)

// Service holds the FNDM005 business logic.
type Service struct {
	db  *gorm.DB
	log *slog.Logger
}

// NewService creates a new FNDM005 service.
func NewService(database *gorm.DB, logger *slog.Logger) *Service {
	return &Service{
		db:  database,
		log: logger.With(slog.String("module", "fndm005-service")),
	}
}

// TAFNDFundCalDate retrieves Master, Memo, and Fund structures from PostgreSQL via GORM.
func (s *Service) TAFNDFundCalDate(ctx context.Context, req *fndv1.TAFNDFundCalDateRequest) (*fndv1.TAFNDFundCalDateResponse, error) {
	sysCoID := req.GetSysCoId()
	calYear := req.GetCalYear()
	fndCalType := req.GetFndCalType()
	fundCry := req.GetFundCry()

	resp := &fndv1.TAFNDFundCalDateResponse{
		ResultList:             make([]*fndv1.FundCalMasterDTO, 0),
		FundClosedMemoList:     make([]*fndv1.FundClosedDateMemoListDTO, 0),
		FundClosedDateFundList: make([]*fndv1.FundClosedDateFundListDTO, 0),
	}

	// 1. Query Master
	masterQuery := s.db.WithContext(ctx).Model(&db.TAFNDFundCal{}).Where("SysCoID = ?", sysCoID)
	if calYear != "" {
		masterQuery = masterQuery.Where("CalYear = ?", calYear)
	}
	if fndCalType != "" {
		masterQuery = masterQuery.Where("FNDCalType = ?", fndCalType)
	}
	if fundCry != "" {
		masterQuery = masterQuery.Where("FundCry = ?", fundCry)
	}

	var masters []db.TAFNDFundCal
	if err := masterQuery.Find(&masters).Error; err != nil {
		s.log.Error("Failed to query TAFNDFundCal", slog.Any("error", err))
		return nil, err
	}

	for _, m := range masters {
		dto := &fndv1.FundCalMasterDTO{
			SysCoId:    m.SysCoID,
			SysCoIdNm:  shared.GetSysCoName(m.SysCoID),
			CalYear:    m.CalYear,
			FndCalType: m.FNDCalType,
			FundCry:    m.FundCry,
			FundCryNm:  shared.GetCryName(m.FundCry),
		}
		resp.ResultList = append(resp.ResultList, dto)
	}

	// 2. Query Memo
	memoQuery := s.db.WithContext(ctx).Model(&db.TAFNDFundCalMemo{}).Where("SysCoID = ?", sysCoID)
	if calYear != "" {
		memoQuery = memoQuery.Where("CalYear = ?", calYear)
	}
	if fndCalType != "" {
		memoQuery = memoQuery.Where("FNDCalType = ?", fndCalType)
	}
	if fundCry != "" {
		memoQuery = memoQuery.Where("FundCry = ?", fundCry)
	}

	var memos []db.TAFNDFundCalMemo
	if err := memoQuery.Find(&memos).Error; err != nil {
		s.log.Error("Failed to query TAFNDFundCalMemo", slog.Any("error", err))
		return nil, err
	}

	for _, m := range memos {
		dto := &fndv1.FundClosedDateMemoListDTO{
			SysCoId:    m.SysCoID,
			CalYear:    m.CalYear,
			FndCalType: m.FNDCalType,
			FundCry:    m.FundCry,
			CalDate:    m.CalDate.Format(time.RFC3339),
			Remark:     m.Remark,
		}
		resp.FundClosedMemoList = append(resp.FundClosedMemoList, dto)
	}

	// 3. Query Fund List
	fundQuery := s.db.WithContext(ctx).Model(&db.TAFNDFundCalDtl{}).Where("SysCoID = ?", sysCoID)
	if calYear != "" {
		fundQuery = fundQuery.Where("CalYear = ?", calYear)
	}
	if fndCalType != "" {
		fundQuery = fundQuery.Where("FNDCalType = ?", fndCalType)
	}
	if fundCry != "" {
		fundQuery = fundQuery.Where("FundCry = ?", fundCry)
	}

	var funds []db.TAFNDFundCalDtl
	if err := fundQuery.Find(&funds).Error; err != nil {
		s.log.Error("Failed to query TAFNDFundCalDtl", slog.Any("error", err))
		return nil, err
	}

	for _, f := range funds {
		dto := &fndv1.FundClosedDateFundListDTO{
			SysCoId:       f.SysCoID,
			CalYear:       f.CalYear,
			FndCalType:    f.FNDCalType,
			FundCry:       f.FundCry,
			PrtFundCode:   f.PrtFundCode,
			PrtFundCodeNm: shared.GetDFundName(f.PrtFundCode),
			CalDate:       f.CalDate.Format(time.RFC3339),
		}
		resp.FundClosedDateFundList = append(resp.FundClosedDateFundList, dto)
	}

	return resp, nil
}
