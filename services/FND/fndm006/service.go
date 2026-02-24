package fndm006

import (
	"context"
	"log/slog"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
	"go-transfer-agent/services/fnd/fndm006/db"
	"go-transfer-agent/services/fnd/shared"

	"gorm.io/gorm"
)

// Service holds the FNDM006 business logic.
type Service struct {
	db  *gorm.DB
	log *slog.Logger
}

// NewService creates a new FNDM006 service.
func NewService(database *gorm.DB, logger *slog.Logger) *Service {
	return &Service{
		db:  database,
		log: logger.With(slog.String("module", "fndm006-service")),
	}
}

// TAFNDRPFeeChgType retrieves Master and Detail structures from PostgreSQL via GORM.
func (s *Service) TAFNDRPFeeChgType(ctx context.Context, req *fndv1.TAFNDRPFeeChgTypeRequest) (*fndv1.TAFNDRPFeeChgTypeResponse, error) {
	sysCoID := req.GetSysCoId()
	prtFundCode := req.GetPrtFundCode()
	pmtTxnType := req.GetPmtTxnType()

	resp := &fndv1.TAFNDRPFeeChgTypeResponse{
		ResultList:  make([]*fndv1.RPFeeChgTypeMasterDTO, 0),
		ChgTypeList: make([]*fndv1.ChgTypeListDTO, 0),
	}

	// 1. Query Master
	masterQuery := s.db.WithContext(ctx).Model(&db.DTAFNDRPFeeChgType{}).Where("SysCoID = ?", sysCoID)
	if prtFundCode != "" {
		masterQuery = masterQuery.Where("PrtFundCode = ?", prtFundCode)
	}
	if pmtTxnType != "" {
		masterQuery = masterQuery.Where("PmtTxnType = ?", pmtTxnType)
	}

	var masters []db.DTAFNDRPFeeChgType
	if err := masterQuery.Find(&masters).Error; err != nil {
		s.log.Error("Failed to query DTAFNDRPFeeChgType", slog.Any("error", err))
		return nil, err
	}

	for _, m := range masters {
		dto := &fndv1.RPFeeChgTypeMasterDTO{
			SysCoId:       m.SysCoID,
			SysCoIdNm:     shared.GetSysCoName(m.SysCoID),
			PmtTxnType:    m.PmtTxnType,
			PrtFundCode:   m.PrtFundCode,
			PrtFundCodeNm: shared.GetDFundName(m.PrtFundCode),
		}
		resp.ResultList = append(resp.ResultList, dto)
	}

	// 2. Query Detail
	dtlQuery := s.db.WithContext(ctx).Model(&db.DTAFNDRPFeeChgTypeDtl{}).Where("SysCoID = ?", sysCoID)
	if prtFundCode != "" {
		dtlQuery = dtlQuery.Where("PrtFundCode = ?", prtFundCode)
	}
	if pmtTxnType != "" {
		dtlQuery = dtlQuery.Where("PmtTxnType = ?", pmtTxnType)
	}

	var details []db.DTAFNDRPFeeChgTypeDtl
	if err := dtlQuery.Find(&details).Error; err != nil {
		s.log.Error("Failed to query DTAFNDRPFeeChgTypeDtl", slog.Any("error", err))
		return nil, err
	}

	for _, d := range details {
		dto := &fndv1.ChgTypeListDTO{
			SysCoId:         d.SysCoID,
			PmtTxnType:      d.PmtTxnType,
			PrtFundCode:     d.PrtFundCode,
			CryKind:         d.CryKind,
			RemitFeeObj:     d.RemitFeeObj,
			RemitFeeChgType: d.RemitFeeChgType,
			PostFeeObj:      d.PostFeeObj,
		}
		resp.ChgTypeList = append(resp.ChgTypeList, dto)
	}

	return resp, nil
}
