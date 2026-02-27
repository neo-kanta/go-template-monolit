package fndm010

import (
	"context"
	"log/slog"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
	"go-transfer-agent/services/fnd/fndm010/db"
	"go-transfer-agent/services/fnd/shared"

	"go-transfer-agent/common/platform/model"
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
			RangeAmtAbove: d.RangeAmtAbove.InexactFloat64(),
			SubsFeeRate:   d.SubsFeeRate.InexactFloat64(),
		}
		resp.IShareFundFeeList = append(resp.IShareFundFeeList, dto)
	}

	return resp, nil
}

// SaveTAFNDIShareFundFee creates a new record in the _Edit table — POST endpoint.
func (s *Service) SaveTAFNDIShareFundFee(ctx context.Context, req *fndv1.TAFNDIShareFundFeeRequest) (*fndv1.SaveResponse, error) {
	s.log.Info("SaveTAFNDIShareFundFee called")
	record := db.DTAFNDIShareFundFeeEdit{}
	// Auth Context & 4-Eyes Principle
	record.MakerID = shared.GetUsernameFromCtx(ctx)
	record.Status = models.StatusPendingApproval

	record.SysCoID = req.GetSysCoId()
	record.FundCode = req.GetFundCode()
	if err := s.db.WithContext(ctx).Create(&record).Error; err != nil {
		s.log.Error("SaveTAFNDIShareFundFee failed", slog.Any("error", err))
		return &fndv1.SaveResponse{Success: false, Message: err.Error(), ReturnCode: "DB_ERROR"}, nil
	}
	return &fndv1.SaveResponse{Success: true, Message: "Record created successfully", ReturnCode: "0000"}, nil
}

// UpdateTAFNDIShareFundFee updates an existing record in the _Edit table — PUT endpoint.
func (s *Service) UpdateTAFNDIShareFundFee(ctx context.Context, req *fndv1.TAFNDIShareFundFeeRequest) (*fndv1.SaveResponse, error) {
	s.log.Info("UpdateTAFNDIShareFundFee called")
	record := db.DTAFNDIShareFundFeeEdit{}
	record.SysCoID = req.GetSysCoId()
	record.FundCode = req.GetFundCode()
	result := s.db.WithContext(ctx).Model(&db.DTAFNDIShareFundFeeEdit{}).Where(`"DataID" = ?`, req.GetSysCoId()).Updates(&record)
	if result.Error != nil {
		s.log.Error("UpdateTAFNDIShareFundFee failed", slog.Any("error", result.Error))
		return &fndv1.SaveResponse{Success: false, Message: result.Error.Error(), ReturnCode: "DB_ERROR"}, nil
	}
	if result.RowsAffected == 0 {
		return &fndv1.SaveResponse{Success: false, Message: "No record found to update", ReturnCode: "NOT_FOUND"}, nil
	}
	return &fndv1.SaveResponse{Success: true, Message: "Record updated successfully", ReturnCode: "0000"}, nil
}

// DeleteTAFNDIShareFundFee removes a record from the _Edit table — DELETE endpoint.
func (s *Service) DeleteTAFNDIShareFundFee(ctx context.Context, req *fndv1.DeleteRequest) (*fndv1.SaveResponse, error) {
	s.log.Info("DeleteTAFNDIShareFundFee called",
		slog.String("data_id", req.GetDataId()),
		slog.String("data_flag", req.GetDataFlag()),
	)
	result := s.db.WithContext(ctx).Where(`"DataID" = ?`, req.GetDataId()).Delete(&db.DTAFNDIShareFundFeeEdit{})
	if result.Error != nil {
		s.log.Error("DeleteTAFNDIShareFundFee failed", slog.Any("error", result.Error))
		return &fndv1.SaveResponse{Success: false, Message: result.Error.Error(), ReturnCode: "DB_ERROR"}, nil
	}
	if result.RowsAffected == 0 {
		return &fndv1.SaveResponse{Success: false, Message: "No record found with the given DataID", ReturnCode: "NOT_FOUND"}, nil
	}
	return &fndv1.SaveResponse{Success: true, Message: "Record deleted successfully", ReturnCode: "0000"}, nil
}

// GetDataByDataID retrieves full data by DataID.
func (s *Service) GetDataByDataID(ctx context.Context, req *fndv1.GetDataRequest) (*fndv1.TAFNDIShareFundFeeRequest, error) {
	return &fndv1.TAFNDIShareFundFeeRequest{}, nil
}

// ApproveTAFNDIShareFundFee approves or rejects a pending record for the 4-Eyes Principle.
func (s *Service) ApproveTAFNDIShareFundFee(ctx context.Context, req *fndv1.ApproveTAFNDIShareFundFeeRequest) (*fndv1.SaveResponse, error) {
	s.log.Info("ApproveTAFNDIShareFundFee called", slog.String("data_id", req.GetDataId()))

	var record db.DTAFNDIShareFundFeeEdit
	result := s.db.WithContext(ctx).Where("\"DataID\" = ?", req.GetDataId()).First(&record)
	
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return &fndv1.SaveResponse{Success: false, Message: "Record not found", ReturnCode: "NOT_FOUND"}, nil
		}
		s.log.Error("ApproveTAFNDIShareFundFee DB error", slog.Any("error", result.Error))
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
