package fndm009

import (
	"context"
	"log/slog"
	"strings"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
	"go-transfer-agent/services/fnd/fndm009/db"
	"go-transfer-agent/services/fnd/shared"

	"gorm.io/gorm"
)

// Service holds the FNDM009 business logic.
type Service struct {
	db  *gorm.DB
	log *slog.Logger
}

// NewService creates a new FNDM009 service.
func NewService(database *gorm.DB, logger *slog.Logger) *Service {
	return &Service{
		db:  database,
		log: logger.With(slog.String("module", "fndm009-service")),
	}
}

// TAFNDFavDisc retrieves Master and Detail structures from PostgreSQL via GORM.
func (s *Service) TAFNDFavDisc(ctx context.Context, req *fndv1.TAFNDFavDiscRequest) (*fndv1.TAFNDFavDiscResponse, error) {
	sysCoID := req.GetSysCoId()
	cusIDCode := req.GetCusIdCode()

	resp := &fndv1.TAFNDFavDiscResponse{
		ResultList:           make([]*fndv1.FavDiscMasterDTO, 0),
		PrtFundList:          make([]*fndv1.FavDiscPrtFundListDTO, 0),
		AllotDiscTypeDtlList: make([]*fndv1.FavDiscTypeDtlListDTO, 0),
		RspDiscTypeDtlList:   make([]*fndv1.FavDiscTypeDtlListDTO, 0),
		SwDiscTypeDtlList:    make([]*fndv1.FavDiscTypeDtlListDTO, 0),
	}

	// 1. Query Master
	masterQuery := s.db.WithContext(ctx).Model(&db.TAFNDFavDisc{}).Where("SysCoID = ?", sysCoID)
	if cusIDCode != "" {
		masterQuery = masterQuery.Where("CusIDCode = ?", cusIDCode)
	}

	var masters []db.TAFNDFavDisc
	if err := masterQuery.Find(&masters).Error; err != nil {
		s.log.Error("Failed to query TAFNDFavDisc", slog.Any("error", err))
		return nil, err
	}

	for _, m := range masters {
		dto := &fndv1.FavDiscMasterDTO{
			SysCoId:      m.SysCoID,
			SysCoIdNm:    shared.GetSysCoName(m.SysCoID),
			CusIdCode:    m.CusIDCode,
			CusIdCodeNm:  shared.GetCusName(m.CusIDCode),
			DiscItemSet:  strings.Split(m.DiscItemSet, ","),
			DiscFundType: m.DiscFundType,
		}
		resp.ResultList = append(resp.ResultList, dto)
	}

	// 2. Query Fund List
	fundQuery := s.db.WithContext(ctx).Model(&db.TAFNDFavDiscFund{}).Where("SysCoID = ?", sysCoID)
	if cusIDCode != "" {
		fundQuery = fundQuery.Where("CusIDCode = ?", cusIDCode)
	}

	var funds []db.TAFNDFavDiscFund
	if err := fundQuery.Find(&funds).Error; err != nil {
		s.log.Error("Failed to query TAFNDFavDiscFund", slog.Any("error", err))
		return nil, err
	}

	for _, f := range funds {
		dto := &fndv1.FavDiscPrtFundListDTO{
			SysCoId:       f.SysCoID,
			CusIdCode:     f.CusIDCode,
			PrtFundCode:   f.PrtFundCode,
			PrtFundCodeNm: shared.GetDFundName(f.PrtFundCode),
		}
		resp.PrtFundList = append(resp.PrtFundList, dto)
	}

	// 3. Query Type Dtl
	dtlQuery := s.db.WithContext(ctx).Model(&db.TAFNDFavDiscTypeDtl{}).Where("SysCoID = ?", sysCoID)
	if cusIDCode != "" {
		dtlQuery = dtlQuery.Where("CusIDCode = ?", cusIDCode)
	}

	var details []db.TAFNDFavDiscTypeDtl
	if err := dtlQuery.Find(&details).Error; err != nil {
		s.log.Error("Failed to query TAFNDFavDiscTypeDtl", slog.Any("error", err))
		return nil, err
	}

	for _, d := range details {
		dto := &fndv1.FavDiscTypeDtlListDTO{
			SysCoId:       d.SysCoID,
			CusIdCode:     d.CusIDCode,
			DiscItem:      d.DiscItem,
			TxCry:         d.TxCry,
			TxCryNm:       shared.GetCryName(d.TxCry),
			RangeAmtAbove: d.RangeAmtAbove,
			RangeFeeRate:  d.RangeFeeRate,
		}

		// Map by DiscItem conceptually
		// In a real system, `Allot`=1, `RSP`=2, `SW`=3 or similar.
		if d.DiscItem == "1" || d.DiscItem == "Allot" {
			resp.AllotDiscTypeDtlList = append(resp.AllotDiscTypeDtlList, dto)
		} else if d.DiscItem == "2" || d.DiscItem == "RSP" {
			resp.RspDiscTypeDtlList = append(resp.RspDiscTypeDtlList, dto)
		} else if d.DiscItem == "3" || d.DiscItem == "SW" {
			resp.SwDiscTypeDtlList = append(resp.SwDiscTypeDtlList, dto)
		} else {
			// default or fallback, mapping generic
			resp.AllotDiscTypeDtlList = append(resp.AllotDiscTypeDtlList, dto)
		}
	}

	return resp, nil
}
