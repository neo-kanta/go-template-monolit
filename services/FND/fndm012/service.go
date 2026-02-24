package fndm012

import (
	"context"
	"log/slog"
	"time"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
	"go-transfer-agent/services/fnd/fndm012/db"
	"go-transfer-agent/services/fnd/shared"

	"gorm.io/gorm"
)

// Service holds the FNDM012 business logic.
type Service struct {
	db  *gorm.DB
	log *slog.Logger
}

// NewService creates a new FNDM012 service.
func NewService(database *gorm.DB, logger *slog.Logger) *Service {
	return &Service{
		db:  database,
		log: logger.With(slog.String("module", "fndm012-service")),
	}
}

// TAFNDPGFundFeeRdm retrieves Master structures from PostgreSQL via GORM.
func (s *Service) TAFNDPGFundFeeRdm(ctx context.Context, req *fndv1.TAFNDPGFundFeeRdmRequest) (*fndv1.TAFNDPGFundFeeRdmResponse, error) {
	sysCoID := req.GetSysCoId()
	prtFundCode := req.GetPrtFundCode()
	rdmCalcBegDate := req.GetRdmCalcBegDate()
	feeRate := req.GetFeeRate()

	resp := &fndv1.TAFNDPGFundFeeRdmResponse{
		ResultList: make([]*fndv1.PGFundFeeRdmMasterDTO, 0),
	}

	// 1. Query Master
	masterQuery := s.db.WithContext(ctx).Model(&db.DTAFNDPGFundFeeRdm{}).Where("SysCoID = ?", sysCoID)
	if prtFundCode != "" {
		masterQuery = masterQuery.Where("PrtFundCode = ?", prtFundCode)
	}
	if rdmCalcBegDate != "" {
		masterQuery = masterQuery.Where("RdmCalcBegDate = ?", rdmCalcBegDate)
	}
	if feeRate > 0 {
		masterQuery = masterQuery.Where("FeeRate = ?", feeRate)
	}

	var masters []db.DTAFNDPGFundFeeRdm
	if err := masterQuery.Find(&masters).Error; err != nil {
		s.log.Error("Failed to query DTAFNDPGFundFeeRdm", slog.Any("error", err))
		return nil, err
	}

	for _, m := range masters {
		dto := &fndv1.PGFundFeeRdmMasterDTO{
			SysCoId:        m.SysCoID,
			SysCoIdNm:      shared.GetSysCoName(m.SysCoID),
			PrtFundCode:    m.PrtFundCode,
			PrtFundCodeNm:  shared.GetDFundName(m.PrtFundCode),
			RdmCalcBegDate: m.RdmCalcBegDate.Format(time.RFC3339),
			FeeRate:        m.FeeRate,
		}
		resp.ResultList = append(resp.ResultList, dto)
	}

	return resp, nil
}
