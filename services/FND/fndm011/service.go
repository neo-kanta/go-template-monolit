package fndm011

import (
	"context"
	"log/slog"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
	"go-transfer-agent/services/fnd/fndm011/db"
	"go-transfer-agent/services/fnd/shared"

	"gorm.io/gorm"
)

// Service holds the FNDM011 business logic.
type Service struct {
	db  *gorm.DB
	log *slog.Logger
}

// NewService creates a new FNDM011 service.
func NewService(database *gorm.DB, logger *slog.Logger) *Service {
	return &Service{
		db:  database,
		log: logger.With(slog.String("module", "fndm011-service")),
	}
}

// TAFNDIShareFundFeeRdm retrieves Master structures from PostgreSQL via GORM.
func (s *Service) TAFNDIShareFundFeeRdm(ctx context.Context, req *fndv1.TAFNDIShareFundFeeRdmRequest) (*fndv1.TAFNDIShareFundFeeRdmResponse, error) {
	sysCoID := req.GetSysCoId()
	fundCode := req.GetFundCode()
	feeName := req.GetFeeName()

	resp := &fndv1.TAFNDIShareFundFeeRdmResponse{
		ResultList: make([]*fndv1.IShareFundFeeRdmMasterDTO, 0),
	}

	// 1. Query Master
	masterQuery := s.db.WithContext(ctx).Model(&db.DTAFNDIShareFundFeeRdm{}).Where("SysCoID = ?", sysCoID)
	if fundCode != "" {
		masterQuery = masterQuery.Where("FundCode = ?", fundCode)
	}
	if feeName != "" {
		masterQuery = masterQuery.Where("FeeName = ?", feeName)
	}

	// Add other optional filters if provided logically
	if req.GetRdmRangeType() != "" {
		masterQuery = masterQuery.Where("RdmRangeType = ?", req.GetRdmRangeType())
	}
	if req.GetRdmDateType() != "" {
		masterQuery = masterQuery.Where("RdmDateType = ?", req.GetRdmDateType())
	}

	var masters []db.DTAFNDIShareFundFeeRdm
	if err := masterQuery.Find(&masters).Error; err != nil {
		s.log.Error("Failed to query DTAFNDIShareFundFeeRdm", slog.Any("error", err))
		return nil, err
	}

	for _, m := range masters {
		dto := &fndv1.IShareFundFeeRdmMasterDTO{
			SysCoId:        m.SysCoID,
			SysCoIdNm:      shared.GetSysCoName(m.SysCoID),
			FundCode:       m.FundCode,
			FundCodeNm:     shared.GetDFundName(m.FundCode),
			FeeName:        m.FeeName,
			RdmRangeType:   m.RdmRangeType,
			RdmDateType:    m.RdmDateType,
			RdmBaseId:      m.RdmBaseID,
			SubsBaseId:     m.SubsBaseID,
			RdmCalcId:      m.RdmCalcID,
			ShouldHoldDays: int32(m.ShouldHoldDays),
			FeeRate:        m.FeeRate,
		}
		resp.ResultList = append(resp.ResultList, dto)
	}

	return resp, nil
}
