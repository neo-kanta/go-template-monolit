package fndm008

import (
	"context"
	"log/slog"
	"time"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
	"go-transfer-agent/services/fnd/fndm008/db"
	"go-transfer-agent/services/fnd/shared"

	"gorm.io/gorm"
)

// Service holds the FNDM008 business logic.
type Service struct {
	db  *gorm.DB
	log *slog.Logger
}

// NewService creates a new FNDM008 service.
func NewService(database *gorm.DB, logger *slog.Logger) *Service {
	return &Service{
		db:  database,
		log: logger.With(slog.String("module", "fndm008-service")),
	}
}

// TAFNDPauseTxn retrieves Master and Detail structures from PostgreSQL via GORM.
func (s *Service) TAFNDPauseTxn(ctx context.Context, req *fndv1.TAFNDPauseTxnRequest) (*fndv1.TAFNDPauseTxnResponse, error) {
	sysCoID := req.GetSysCoId()
	prtFundCode := req.GetPrtFundCode()

	resp := &fndv1.TAFNDPauseTxnResponse{
		ResultList:       make([]*fndv1.PauseTxnMasterDTO, 0),
		PauseTxnCryList:  make([]*fndv1.PauseTxnCryListDTO, 0),
		PauseTxnTypeList: make([]*fndv1.PauseTxnTypeListDTO, 0),
	}

	// 1. Query Master
	masterQuery := s.db.WithContext(ctx).Model(&db.DTAFNDPauseTxn{}).Where("SysCoID = ?", sysCoID)
	if prtFundCode != "" {
		masterQuery = masterQuery.Where("PrtFundCode = ?", prtFundCode)
	}

	var masters []db.DTAFNDPauseTxn
	if err := masterQuery.Find(&masters).Error; err != nil {
		s.log.Error("Failed to query DTAFNDPauseTxn", slog.Any("error", err))
		return nil, err
	}

	for _, m := range masters {
		dto := &fndv1.PauseTxnMasterDTO{
			SysCoId:       m.SysCoID,
			SysCoIdNm:     shared.GetSysCoName(m.SysCoID),
			PrtFundCode:   m.PrtFundCode,
			PrtFundCodeNm: shared.GetDFundName(m.PrtFundCode),
			PTxnBegDate:   m.PTxnBegDate.Format(time.RFC3339),
			PTxnEndDate:   m.PTxnEndDate.Format(time.RFC3339),
			Remark:        m.Remark,
		}
		resp.ResultList = append(resp.ResultList, dto)
	}

	// 2. Query Cry List
	cryQuery := s.db.WithContext(ctx).Model(&db.DTAFNDPauseTxnCry{}).Where("SysCoID = ?", sysCoID)
	if prtFundCode != "" {
		cryQuery = cryQuery.Where("PrtFundCode = ?", prtFundCode)
	}

	var cries []db.DTAFNDPauseTxnCry
	if err := cryQuery.Find(&cries).Error; err != nil {
		s.log.Error("Failed to query DTAFNDPauseTxnCry", slog.Any("error", err))
		return nil, err
	}

	for _, c := range cries {
		dto := &fndv1.PauseTxnCryListDTO{
			SysCoId:     c.SysCoID,
			PrtFundCode: c.PrtFundCode,
			PTxnBegDate: c.PTxnBegDate.Format(time.RFC3339),
			CryId:       c.CryID,
			CryIdNm:     shared.GetCryName(c.CryID),
		}
		resp.PauseTxnCryList = append(resp.PauseTxnCryList, dto)
	}

	// 3. Query Type Dtl
	dtlQuery := s.db.WithContext(ctx).Model(&db.DTAFNDPauseTxnDtl{}).Where("SysCoID = ?", sysCoID)
	if prtFundCode != "" {
		dtlQuery = dtlQuery.Where("PrtFundCode = ?", prtFundCode)
	}

	var details []db.DTAFNDPauseTxnDtl
	if err := dtlQuery.Find(&details).Error; err != nil {
		s.log.Error("Failed to query DTAFNDPauseTxnDtl", slog.Any("error", err))
		return nil, err
	}

	for _, d := range details {
		dto := &fndv1.PauseTxnTypeListDTO{
			SysCoId:      d.SysCoID,
			PrtFundCode:  d.PrtFundCode,
			PTxnBegDate:  d.PTxnBegDate.Format(time.RFC3339),
			PauseTxnType: d.PauseTxnType,
		}
		resp.PauseTxnTypeList = append(resp.PauseTxnTypeList, dto)
	}

	return resp, nil
}
