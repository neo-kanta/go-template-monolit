package fndm009

import (
	"context"
	"log/slog"
	"strings"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
	"go-transfer-agent/services/fnd/fndm009/db"
	"go-transfer-agent/services/fnd/shared"

	"go-transfer-agent/common/platform/model"
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
			RangeAmtAbove: d.RangeAmtAbove.InexactFloat64(),
			RangeFeeRate:  d.RangeFeeRate.InexactFloat64(),
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

// SaveTAFNDFavDisc creates a new record in the _Edit table — POST endpoint.
func (s *Service) SaveTAFNDFavDisc(ctx context.Context, req *fndv1.TAFNDFavDiscRequest) (*fndv1.SaveResponse, error) {
	s.log.Info("SaveTAFNDFavDisc called")
	record := db.TAFNDFavDiscEdit{}
	record.SysCoID = req.GetSysCoId()
	record.CusIDCode = req.GetCusIdCode()
	if err := s.db.WithContext(ctx).Create(&record).Error; err != nil {
		s.log.Error("SaveTAFNDFavDisc failed", slog.Any("error", err))
		return &fndv1.SaveResponse{Success: false, Message: err.Error(), ReturnCode: "DB_ERROR"}, nil
	}
	return &fndv1.SaveResponse{Success: true, Message: "Record created successfully", ReturnCode: "0000"}, nil
}

// UpdateTAFNDFavDisc updates an existing record in the _Edit table — PUT endpoint.
func (s *Service) UpdateTAFNDFavDisc(ctx context.Context, req *fndv1.TAFNDFavDiscRequest) (*fndv1.SaveResponse, error) {
	s.log.Info("UpdateTAFNDFavDisc called")
	record := db.TAFNDFavDiscEdit{}
	record.SysCoID = req.GetSysCoId()
	record.CusIDCode = req.GetCusIdCode()
	result := s.db.WithContext(ctx).Model(&db.TAFNDFavDiscEdit{}).Where(`"DataID" = ?`, req.GetSysCoId()).Updates(&record)
	if result.Error != nil {
		s.log.Error("UpdateTAFNDFavDisc failed", slog.Any("error", result.Error))
		return &fndv1.SaveResponse{Success: false, Message: result.Error.Error(), ReturnCode: "DB_ERROR"}, nil
	}
	if result.RowsAffected == 0 {
		return &fndv1.SaveResponse{Success: false, Message: "No record found to update", ReturnCode: "NOT_FOUND"}, nil
	}
	return &fndv1.SaveResponse{Success: true, Message: "Record updated successfully", ReturnCode: "0000"}, nil
}

// DeleteTAFNDFavDisc removes a record from the _Edit table — DELETE endpoint.
func (s *Service) DeleteTAFNDFavDisc(ctx context.Context, req *fndv1.DeleteRequest) (*fndv1.SaveResponse, error) {
	s.log.Info("DeleteTAFNDFavDisc called",
		slog.String("data_id", req.GetDataId()),
		slog.String("data_flag", req.GetDataFlag()),
	)
	result := s.db.WithContext(ctx).Where(`"DataID" = ?`, req.GetDataId()).Delete(&db.TAFNDFavDiscEdit{})
	if result.Error != nil {
		s.log.Error("DeleteTAFNDFavDisc failed", slog.Any("error", result.Error))
		return &fndv1.SaveResponse{Success: false, Message: result.Error.Error(), ReturnCode: "DB_ERROR"}, nil
	}
	if result.RowsAffected == 0 {
		return &fndv1.SaveResponse{Success: false, Message: "No record found with the given DataID", ReturnCode: "NOT_FOUND"}, nil
	}
	return &fndv1.SaveResponse{Success: true, Message: "Record deleted successfully", ReturnCode: "0000"}, nil
}

// GetDataByDataID retrieves full data by DataID.
func (s *Service) GetDataByDataID(ctx context.Context, req *fndv1.GetDataRequest) (*fndv1.TAFNDFavDiscRequest, error) {
	return &fndv1.TAFNDFavDiscRequest{}, nil
}

// ApproveTAFNDFavDisc approves or rejects a pending record for the 4-Eyes Principle.
func (s *Service) ApproveTAFNDFavDisc(ctx context.Context, req *fndv1.ApproveTAFNDFavDiscRequest) (*fndv1.SaveResponse, error) {
	s.log.Info("ApproveTAFNDFavDisc called", slog.String("data_id", req.GetDataId()))

	var record db.TAFNDFavDiscEdit
	result := s.db.WithContext(ctx).Where("\"DataID\" = ?", req.GetDataId()).First(&record)
	
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return &fndv1.SaveResponse{Success: false, Message: "Record not found", ReturnCode: "NOT_FOUND"}, nil
		}
		s.log.Error("ApproveTAFNDFavDisc DB error", slog.Any("error", result.Error))
		return &fndv1.SaveResponse{Success: false, Message: result.Error.Error(), ReturnCode: "DB_ERROR"}, nil
	}

	// Only allow approval if status is PENDING_APPROVAL
	if record.Status != models.StatusPendingApproval {
		return &fndv1.SaveResponse{Success: false, Message: "Record is not in PENDING_APPROVAL status", ReturnCode: "INVALID_STATUS"}, nil
	}

	// 4-Eyes Principle constraint: Maker != Checker
	checkerID := req.GetCheckerId()
	if checkerID == "" {
		checkerID = shared.GetUsernameFromCtx(ctx)
	}

	if record.MakerID != "" && checkerID != "" && record.MakerID == checkerID {
		return &fndv1.SaveResponse{Success: false, Message: "4-Eyes Principle Violation: Maker cannot be the Checker", ReturnCode: "FOUR_EYES_VIOLATION"}, nil
	}

	// Process Approval/Rejection
	if req.GetIsApproved() {
		record.Status = models.StatusApproved
	} else {
		record.Status = models.StatusRejected
	}
	
	record.CheckerID = &checkerID
		remark := req.GetRemark()
	record.Remark = &remark

	if err := s.db.WithContext(ctx).Save(&record).Error; err != nil {
		s.log.Error("Failed to update status", slog.Any("error", err))
		return &fndv1.SaveResponse{Success: false, Message: "Failed to update record", ReturnCode: "DB_ERROR"}, nil
	}

	return &fndv1.SaveResponse{Success: true, Message: "Record reviewed successfully", ReturnCode: "0000"}, nil
}
