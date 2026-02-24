package fndm010

import (
	"context"
	"log/slog"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
	"go-transfer-agent/services/fnd/fndm010/db"
	"go-transfer-agent/services/fnd/shared"

	"gorm.io/gorm"
)

// Service holds the FNDM010 business logic.
type Service struct {
	db  *gorm.DB
	log *slog.Logger
}

// NewService creates a new FNDM010 service.
func NewService(database *gorm.DB, logger *slog.Logger) *Service {
	return &Service{
		db:  database,
		log: logger.With(slog.String("module", "fndm010-service")),
	}
}

// TAFNDIShareFundFee retrieves Master and Detail structures from PostgreSQL via GORM.
func (s *Service) TAFNDIShareFundFee(ctx context.Context, req *fndv1.TAFNDIShareFundFeeRequest) (*fndv1.TAFNDIShareFundFeeResponse, error) {
	sysCoID := req.GetSysCoId()
	fundCode := req.GetFundCode()

	resp := &fndv1.TAFNDIShareFundFeeResponse{
		ResultList:        make([]*fndv1.IShareFundFeeMasterDTO, 0),
		IShareFundFeeList: make([]*fndv1.IShareFundFeeListDTO, 0),
	}

	// 1. Query Master
	masterQuery := s.db.WithContext(ctx).Model(&db.DTAFNDIShareFundFee{}).Where("SysCoID = ?", sysCoID)
	if fundCode != "" {
		masterQuery = masterQuery.Where("FundCode = ?", fundCode)
	}

	var masters []db.DTAFNDIShareFundFee
	if err := masterQuery.Find(&masters).Error; err != nil {
		s.log.Error("Failed to query DTAFNDIShareFundFee", slog.Any("error", err))
		return nil, err
	}

	for _, m := range masters {
		dto := &fndv1.IShareFundFeeMasterDTO{
			SysCoId:    m.SysCoID,
			SysCoIdNm:  shared.GetSysCoName(m.SysCoID),
			FundCode:   m.FundCode,
			FundCodeNm: shared.GetDFundName(m.FundCode),
		}
		resp.ResultList = append(resp.ResultList, dto)
	}

	// 2. Query Detail
	dtlQuery := s.db.WithContext(ctx).Model(&db.DTAFNDIShareFundFeeSub{}).Where("SysCoID = ?", sysCoID)
	if fundCode != "" {
		dtlQuery = dtlQuery.Where("FundCode = ?", fundCode)
	}

	var details []db.DTAFNDIShareFundFeeSub
	if err := dtlQuery.Find(&details).Error; err != nil {
		s.log.Error("Failed to query DTAFNDIShareFundFeeSub", slog.Any("error", err))
		return nil, err
	}

	for _, d := range details {
		dto := &fndv1.IShareFundFeeListDTO{
			SysCoId:       d.SysCoID,
			FundCode:      d.FundCode,
			RangeAmtAbove: d.RangeAmtAbove,
			SubsFeeRate:   d.SubsFeeRate,
		}
		resp.IShareFundFeeList = append(resp.IShareFundFeeList, dto)
	}

	return resp, nil
}
