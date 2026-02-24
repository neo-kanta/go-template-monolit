package fndm013

import (
	"context"
	"log/slog"
	"time"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
	"go-transfer-agent/services/fnd/fndm013/db"
	"go-transfer-agent/services/fnd/shared"

	"gorm.io/gorm"
)

// Service holds the FNDM013 business logic.
type Service struct {
	db  *gorm.DB
	log *slog.Logger
}

// NewService creates a new FNDM013 service overriding the in-memory mock.
func NewService(database *gorm.DB, logger *slog.Logger) *Service {
	return &Service{
		db:  database,
		log: logger.With(slog.String("module", "fndm013-service")),
	}
}

// TAFNDTMFundFeeRdm retrieves Master and Detail structures from PostgreSQL via GORM.
func (s *Service) TAFNDTMFundFeeRdm(ctx context.Context, req *fndv1.TAFNDTMFundFeeRdmRequest) (*fndv1.TAFNDTMFundFeeRdmResponse, error) {
	sysCoID := req.GetSysCoId()
	prtFundCode := req.GetPrtFundCode()
	rdmCalcBegDate := req.GetRdmCalcBegDate()

	resp := &fndv1.TAFNDTMFundFeeRdmResponse{
		ResultList:     make([]*fndv1.TMFundFeeRdmMasterDTO, 0),
		FundFeeDtlList: make([]*fndv1.TMFundFeeRdmDtlListDTO, 0),
	}

	// 1. Query Master
	masterQuery := s.db.WithContext(ctx).Model(&db.DTAFNDTMFundFeeRdm{}).Where("SysCoID = ?", sysCoID)
	if prtFundCode != "" {
		masterQuery = masterQuery.Where("PrtFundCode = ?", prtFundCode)
	}
	if rdmCalcBegDate != "" {
		masterQuery = masterQuery.Where("RdmCalcBegDate = ?", rdmCalcBegDate)
	}

	var masters []db.DTAFNDTMFundFeeRdm
	if err := masterQuery.Find(&masters).Error; err != nil {
		s.log.Error("Failed to query DTAFNDTMFundFeeRdm", slog.Any("error", err))
		return nil, err
	}

	for _, m := range masters {
		dto := &fndv1.TMFundFeeRdmMasterDTO{
			SysCoId:        m.SysCoID,
			SysCoIdNm:      shared.GetSysCoName(m.SysCoID),
			PrtFundCode:    m.PrtFundCode,
			PrtFundCodeNm:  shared.GetDFundName(m.PrtFundCode),
			RdmCalcBegDate: m.RdmCalcBegDate.Format(time.RFC3339),
		}
		resp.ResultList = append(resp.ResultList, dto)
	}

	// 2. Query Detail
	dtlQuery := s.db.WithContext(ctx).Model(&db.DTAFNDTMFundFeeRdmDtl{}).Where("SysCoID = ?", sysCoID)
	if prtFundCode != "" {
		dtlQuery = dtlQuery.Where("PrtFundCode = ?", prtFundCode)
	}
	if rdmCalcBegDate != "" {
		dtlQuery = dtlQuery.Where("RdmCalcBegDate = ?", rdmCalcBegDate)
	}

	var details []db.DTAFNDTMFundFeeRdmDtl
	if err := dtlQuery.Find(&details).Error; err != nil {
		s.log.Error("Failed to query DTAFNDTMFundFeeRdmDtl", slog.Any("error", err))
		return nil, err
	}

	for _, d := range details {
		dto := &fndv1.TMFundFeeRdmDtlListDTO{
			SysCoId:        d.SysCoID,
			PrtFundCode:    d.PrtFundCode,
			RdmCalcBegDate: d.RdmCalcBegDate.Format(time.RFC3339),
			RdmCalcEndDate: d.RdmCalcEndDate.Format(time.RFC3339),
			FeeRate:        d.FeeRate,
		}
		resp.FundFeeDtlList = append(resp.FundFeeDtlList, dto)
	}

	return resp, nil
}
