package fndm003

import (
	"context"
	"log/slog"
	"strings"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
	"go-transfer-agent/services/fnd/fndm003/db"
	"go-transfer-agent/services/fnd/shared"

	"gorm.io/gorm"
)

// Service holds the FNDM003 business logic.
type Service struct {
	db  *gorm.DB
	log *slog.Logger
}

// NewService creates a new FNDM003 service.
func NewService(database *gorm.DB, logger *slog.Logger) *Service {
	return &Service{
		db:  database,
		log: logger.With(slog.String("module", "fndm003-service")),
	}
}

// TAFNDSwitch retrieves Master and Detail structures from PostgreSQL via GORM.
func (s *Service) TAFNDSwitch(ctx context.Context, req *fndv1.TAFNDSwitchRequest) (*fndv1.TAFNDSwitchResponse, error) {
	sysCoID := req.GetSysCoId()
	prtFundCode := req.GetPrtFundCode()
	isSwitchIn := req.GetIsSwitchIn()

	resp := &fndv1.TAFNDSwitchResponse{
		ResultList:        make([]*fndv1.SwitchMasterDTO, 0),
		SwitchFundList:    make([]*fndv1.SwitchFundListDTO, 0),
		FundSwitchFeeList: make([]*fndv1.FundSwitchFeeListDTO, 0),
		SwitchCryList:     make([]*fndv1.SwitchCryListDTO, 0),
	}

	// 1. Query Master
	masterQuery := s.db.WithContext(ctx).Model(&db.DTAFNDSwitch{}).Where("SysCoID = ?", sysCoID)
	if prtFundCode != "" {
		masterQuery = masterQuery.Where("PrtFundCode = ?", prtFundCode)
	}
	if isSwitchIn != "" {
		masterQuery = masterQuery.Where("IsSwitchIn = ?", isSwitchIn)
	}

	var masters []db.DTAFNDSwitch
	if err := masterQuery.Find(&masters).Error; err != nil {
		s.log.Error("Failed to query DTAFNDSwitch", slog.Any("error", err))
		return nil, err
	}

	for _, m := range masters {
		dto := &fndv1.SwitchMasterDTO{
			SysCoId:       m.SysCoID,
			SysCoIdNm:     shared.GetSysCoName(m.SysCoID),
			PrtFundCode:   m.PrtFundCode,
			PrtFundCodeNm: shared.GetDFundName(m.PrtFundCode),
			IsSwitchIn:    m.IsSwitchIn,
		}
		resp.ResultList = append(resp.ResultList, dto)
	}

	// 2. Query SwitchFund
	fundQuery := s.db.WithContext(ctx).Model(&db.DTAFNDSwitchFund{}).Where("SysCoID = ?", sysCoID)
	if prtFundCode != "" {
		fundQuery = fundQuery.Where("PrtFundCode = ?", prtFundCode)
	}

	var switchFunds []db.DTAFNDSwitchFund
	if err := fundQuery.Find(&switchFunds).Error; err != nil {
		s.log.Error("Failed to query DTAFNDSwitchFund", slog.Any("error", err))
		return nil, err
	}

	for _, f := range switchFunds {
		dto := &fndv1.SwitchFundListDTO{
			SysCoId:          f.SysCoID,
			PrtFundCode:      f.PrtFundCode,
			SwOPrtFundCode:   f.SwOPrtFundCode,
			SwOPrtFundCodeNm: shared.GetDFundName(f.SwOPrtFundCode),
			SwDateType:       f.SwDateType,
		}
		resp.SwitchFundList = append(resp.SwitchFundList, dto)
	}

	// 3. Query FundSwitchFee
	feeQuery := s.db.WithContext(ctx).Model(&db.DTAFNDFundFeeSwitch{}).Where("SysCoID = ?", sysCoID)
	if prtFundCode != "" {
		feeQuery = feeQuery.Where("PrtFundCode = ?", prtFundCode)
	}

	var fees []db.DTAFNDFundFeeSwitch
	if err := feeQuery.Find(&fees).Error; err != nil {
		s.log.Error("Failed to query DTAFNDFundFeeSwitch", slog.Any("error", err))
		return nil, err
	}

	for _, f := range fees {
		dto := &fndv1.FundSwitchFeeListDTO{
			SysCoId:     f.SysCoID,
			PrtFundCode: f.PrtFundCode,
			FundCode:    f.FundCode,
			SwFundType:  f.SwFundType,
			SwDiscType:  f.SwDiscType,
			SwitchRate:  f.SwitchRate,
		}
		resp.FundSwitchFeeList = append(resp.FundSwitchFeeList, dto)
	}

	// 4. Query SwitchCry
	cryQuery := s.db.WithContext(ctx).Model(&db.DTAFNDSwitchCry{}).Where("SysCoID = ?", sysCoID)
	if prtFundCode != "" {
		cryQuery = cryQuery.Where("PrtFundCode = ?", prtFundCode)
	}

	var cries []db.DTAFNDSwitchCry
	if err := cryQuery.Find(&cries).Error; err != nil {
		s.log.Error("Failed to query DTAFNDSwitchCry", slog.Any("error", err))
		return nil, err
	}

	for _, c := range cries {
		dto := &fndv1.SwitchCryListDTO{
			SysCoId:         c.SysCoID,
			PrtFundCode:     c.PrtFundCode,
			SwIFundCry:      c.SwIFundCry,
			SwIFundCryNm:    shared.GetCryName(c.SwIFundCry),
			SwOFundCrySet:   c.SwOFundCrySet,
			SwOFundCrySetNm: splitString(c.SwOFundCrySet), // mapping string slices
		}
		resp.SwitchCryList = append(resp.SwitchCryList, dto)
	}

	return resp, nil
}

func splitString(s string) []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(s, ",")
}
