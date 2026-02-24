package fndm002

import (
	"context"
	"log/slog"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
	"go-transfer-agent/services/fnd/fndm002/db"
	"go-transfer-agent/services/fnd/shared"

	"gorm.io/gorm"
)

// Service holds the FNDM002 business logic.
type Service struct {
	db  *gorm.DB
	log *slog.Logger
}

// NewService creates a new FNDM002 service.
func NewService(database *gorm.DB, logger *slog.Logger) *Service {
	return &Service{
		db:  database,
		log: logger.With(slog.String("module", "fndm002-service")),
	}
}

// TAFNDFundFee retrieves Master, Sub/SubDtl, CDSC/CDSCDtl, Back/BackDtl structures from PostgreSQL via GORM.
func (s *Service) TAFNDFundFee(ctx context.Context, req *fndv1.TAFNDFundFeeRequest) (*fndv1.TAFNDFundFeeResponse, error) {
	sysCoID := req.GetSysCoId()
	prtFundCode := req.GetPrtFundCode()

	resp := &fndv1.TAFNDFundFeeResponse{
		ResultList:      make([]*fndv1.FundFeeMasterDTO, 0),
		FundFeeSubList:  make([]*fndv1.FundFeeSubListResDTO, 0),
		FundFeeCdscList: make([]*fndv1.FundFeeCDSCListResDTO, 0),
		FundFeeBackList: make([]*fndv1.FundFeeBackListResDTO, 0),
	}

	// 1. Query Master
	masterQuery := s.db.WithContext(ctx).Model(&db.DTAFNDFundFee{}).Where("SysCoID = ?", sysCoID)
	if prtFundCode != "" {
		masterQuery = masterQuery.Where("PrtFundCode = ?", prtFundCode)
	}

	var masters []db.DTAFNDFundFee
	if err := masterQuery.Find(&masters).Error; err != nil {
		s.log.Error("Failed to query DTAFNDFundFee", slog.Any("error", err))
		return nil, err
	}

	for _, m := range masters {
		dto := &fndv1.FundFeeMasterDTO{
			SysCoId:       m.SysCoID,
			SysCoIdNm:     shared.GetSysCoName(m.SysCoID),
			PrtFundCode:   m.PrtFundCode,
			PrtFundCodeNm: shared.GetDFundName(m.PrtFundCode),
		}
		resp.ResultList = append(resp.ResultList, dto)
	}

	// 2. Query FundFeeSub and FundFeeSubDtl
	var subs []db.DTAFNDFundFeeSub
	var subDtls []db.DTAFNDFundFeeSubDtl
	subQuery := s.db.WithContext(ctx).Model(&db.DTAFNDFundFeeSub{}).Where("SysCoID = ?", sysCoID)
	subDtlQuery := s.db.WithContext(ctx).Model(&db.DTAFNDFundFeeSubDtl{}).Where("SysCoID = ?", sysCoID)

	if prtFundCode != "" {
		subQuery = subQuery.Where("PrtFundCode = ?", prtFundCode)
		subDtlQuery = subDtlQuery.Where("PrtFundCode = ?", prtFundCode)
	}

	if err := subQuery.Find(&subs).Error; err == nil {
		_ = subDtlQuery.Find(&subDtls).Error

		// Group SubDtls by FundCode + CryID
		dtlMap := make(map[string][]*fndv1.FundFeeSubDtlListResDTO)
		for _, sd := range subDtls {
			key := sd.FundCode + "|" + sd.CryID
			dto := &fndv1.FundFeeSubDtlListResDTO{
				SysCoId:       sd.SysCoID,
				PrtFundCode:   sd.PrtFundCode,
				FundCode:      sd.FundCode,
				CryId:         sd.CryID,
				RangeAmtAbove: sd.RangeAmtAbove,
				SubsFeeRate:   sd.SubsFeeRate,
			}
			dtlMap[key] = append(dtlMap[key], dto)
		}

		for _, sub := range subs {
			key := sub.FundCode + "|" + sub.CryID
			dto := &fndv1.FundFeeSubListResDTO{
				SysCoId:           sub.SysCoID,
				PrtFundCode:       sub.PrtFundCode,
				FundCode:          sub.FundCode,
				CryId:             sub.CryID,
				CryIdNm:           shared.GetCryName(sub.CryID),
				FundFeeSubDtlList: dtlMap[key],
			}
			if dto.FundFeeSubDtlList == nil {
				dto.FundFeeSubDtlList = make([]*fndv1.FundFeeSubDtlListResDTO, 0)
			}
			resp.FundFeeSubList = append(resp.FundFeeSubList, dto)
		}
	} else {
		s.log.Error("Failed to query DTAFNDFundFeeSub", slog.Any("error", err))
	}

	// 3. Query FundFeeCDSC and FundFeeCDSCDtl
	var cdscs []db.DTAFNDFundFeeCDSC
	var cdscDtls []db.DTAFNDFundFeeCDSCDtl
	cdscQuery := s.db.WithContext(ctx).Model(&db.DTAFNDFundFeeCDSC{}).Where("SysCoID = ?", sysCoID)
	cdscDtlQuery := s.db.WithContext(ctx).Model(&db.DTAFNDFundFeeCDSCDtl{}).Where("SysCoID = ?", sysCoID)

	if prtFundCode != "" {
		cdscQuery = cdscQuery.Where("PrtFundCode = ?", prtFundCode)
		cdscDtlQuery = cdscDtlQuery.Where("PrtFundCode = ?", prtFundCode)
	}

	if err := cdscQuery.Find(&cdscs).Error; err == nil {
		_ = cdscDtlQuery.Find(&cdscDtls).Error

		dtlMap := make(map[string][]*fndv1.FundFeeCDSCDtlListResDTO)
		for _, cd := range cdscDtls {
			key := cd.FundCode + "|" + cd.MatureYear
			dto := &fndv1.FundFeeCDSCDtlListResDTO{
				SysCoId:     cd.SysCoID,
				PrtFundCode: cd.PrtFundCode,
				FundCode:    cd.FundCode,
				MatureYear:  cd.MatureYear,
				HoldBegDay:  int32(cd.HoldBegDay),
				CdscFeeRate: cd.CDSCFeeRate,
			}
			dtlMap[key] = append(dtlMap[key], dto)
		}

		for _, c := range cdscs {
			key := c.FundCode + "|" + c.MatureYear
			dto := &fndv1.FundFeeCDSCListResDTO{
				SysCoId:            c.SysCoID,
				PrtFundCode:        c.PrtFundCode,
				FundCode:           c.FundCode,
				MatureYear:         c.MatureYear,
				SubsCalcId:         c.SubsCalcID,
				FundFeeCdscDtlList: dtlMap[key],
			}
			if dto.FundFeeCdscDtlList == nil {
				dto.FundFeeCdscDtlList = make([]*fndv1.FundFeeCDSCDtlListResDTO, 0)
			}
			resp.FundFeeCdscList = append(resp.FundFeeCdscList, dto)
		}
	} else {
		s.log.Error("Failed to query DTAFNDFundFeeCDSC", slog.Any("error", err))
	}

	// 4. Query FundFeeBack and FundFeeBackDtl
	var backs []db.DTAFNDFundFeeBack
	var backDtls []db.DTAFNDFundFeeBackDtl
	backQuery := s.db.WithContext(ctx).Model(&db.DTAFNDFundFeeBack{}).Where("SysCoID = ?", sysCoID)
	backDtlQuery := s.db.WithContext(ctx).Model(&db.DTAFNDFundFeeBackDtl{}).Where("SysCoID = ?", sysCoID)

	if prtFundCode != "" {
		backQuery = backQuery.Where("PrtFundCode = ?", prtFundCode)
		backDtlQuery = backDtlQuery.Where("PrtFundCode = ?", prtFundCode)
	}

	if err := backQuery.Find(&backs).Error; err == nil {
		_ = backDtlQuery.Find(&backDtls).Error

		dtlMap := make(map[string][]*fndv1.FundFeeBackDtlListResDTO)
		for _, bd := range backDtls {
			key := bd.FundCode + "|" + bd.HoldPeriodYear
			dto := &fndv1.FundFeeBackDtlListResDTO{
				SysCoId:        bd.SysCoID,
				PrtFundCode:    bd.PrtFundCode,
				FundCode:       bd.FundCode,
				HoldPeriodYear: bd.HoldPeriodYear,
				HoldBegDay:     int32(bd.HoldBegDay),
				BackFeeRate:    bd.BackFeeRate,
			}
			dtlMap[key] = append(dtlMap[key], dto)
		}

		for _, b := range backs {
			key := b.FundCode + "|" + b.HoldPeriodYear
			dto := &fndv1.FundFeeBackListResDTO{
				SysCoId:            b.SysCoID,
				PrtFundCode:        b.PrtFundCode,
				FundCode:           b.FundCode,
				HoldPeriodYear:     b.HoldPeriodYear,
				SubsCalcId:         b.SubsCalcID,
				FundFeeBackDtlList: dtlMap[key],
			}
			if dto.FundFeeBackDtlList == nil {
				dto.FundFeeBackDtlList = make([]*fndv1.FundFeeBackDtlListResDTO, 0)
			}
			resp.FundFeeBackList = append(resp.FundFeeBackList, dto)
		}
	} else {
		s.log.Error("Failed to query DTAFNDFundFeeBack", slog.Any("error", err))
	}

	return resp, nil
}
