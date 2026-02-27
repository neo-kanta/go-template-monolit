package fndm011

import (
	"context"
	"log/slog"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
	"go-transfer-agent/services/fnd/fndm011/db"
	"go-transfer-agent/services/fnd/shared"

	"go-transfer-agent/common/platform/model"
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
			FeeRate:        m.FeeRate.InexactFloat64(),
		}
		resp.ResultList = append(resp.ResultList, dto)
	}

	return resp, nil
}

// SaveTAFNDIShareFundFeeRdm creates a new record in the _Edit table — POST endpoint.
func (s *Service) SaveTAFNDIShareFundFeeRdm(ctx context.Context, req *fndv1.TAFNDIShareFundFeeRdmRequest) (*fndv1.SaveResponse, error) {
	s.log.Info("SaveTAFNDIShareFundFeeRdm called")
	record := db.DTAFNDIShareFundFeeRdmEdit{}
	// Auth Context & 4-Eyes Principle
	record.MakerID = shared.GetUsernameFromCtx(ctx)
	record.Status = models.StatusPendingApproval

	record.SysCoID = req.GetSysCoId()
	record.FundCode = req.GetFundCode()
	record.FeeName = req.GetFeeName()
	if err := s.db.WithContext(ctx).Create(&record).Error; err != nil {
		s.log.Error("SaveTAFNDIShareFundFeeRdm failed", slog.Any("error", err))
		return &fndv1.SaveResponse{Success: false, Message: err.Error(), ReturnCode: "DB_ERROR"}, nil
	}
	return &fndv1.SaveResponse{Success: true, Message: "Record created successfully", ReturnCode: "0000"}, nil
}

// UpdateTAFNDIShareFundFeeRdm updates an existing record in the _Edit table — PUT endpoint.
func (s *Service) UpdateTAFNDIShareFundFeeRdm(ctx context.Context, req *fndv1.TAFNDIShareFundFeeRdmRequest) (*fndv1.SaveResponse, error) {
	s.log.Info("UpdateTAFNDIShareFundFeeRdm called")
	record := db.DTAFNDIShareFundFeeRdmEdit{}
	record.SysCoID = req.GetSysCoId()
	record.FundCode = req.GetFundCode()
	record.FeeName = req.GetFeeName()
	result := s.db.WithContext(ctx).Model(&db.DTAFNDIShareFundFeeRdmEdit{}).Where(`"DataID" = ?`, req.GetSysCoId()).Updates(&record)
	if result.Error != nil {
		s.log.Error("UpdateTAFNDIShareFundFeeRdm failed", slog.Any("error", result.Error))
		return &fndv1.SaveResponse{Success: false, Message: result.Error.Error(), ReturnCode: "DB_ERROR"}, nil
	}
	if result.RowsAffected == 0 {
		return &fndv1.SaveResponse{Success: false, Message: "No record found to update", ReturnCode: "NOT_FOUND"}, nil
	}
	return &fndv1.SaveResponse{Success: true, Message: "Record updated successfully", ReturnCode: "0000"}, nil
}

// DeleteTAFNDIShareFundFeeRdm removes a record from the _Edit table — DELETE endpoint.
func (s *Service) DeleteTAFNDIShareFundFeeRdm(ctx context.Context, req *fndv1.DeleteRequest) (*fndv1.SaveResponse, error) {
	s.log.Info("DeleteTAFNDIShareFundFeeRdm called",
		slog.String("data_id", req.GetDataId()),
		slog.String("data_flag", req.GetDataFlag()),
	)
	result := s.db.WithContext(ctx).Where(`"DataID" = ?`, req.GetDataId()).Delete(&db.DTAFNDIShareFundFeeRdmEdit{})
	if result.Error != nil {
		s.log.Error("DeleteTAFNDIShareFundFeeRdm failed", slog.Any("error", result.Error))
		return &fndv1.SaveResponse{Success: false, Message: result.Error.Error(), ReturnCode: "DB_ERROR"}, nil
	}
	if result.RowsAffected == 0 {
		return &fndv1.SaveResponse{Success: false, Message: "No record found with the given DataID", ReturnCode: "NOT_FOUND"}, nil
	}
	return &fndv1.SaveResponse{Success: true, Message: "Record deleted successfully", ReturnCode: "0000"}, nil
}

// GetDataByDataID retrieves full data by DataID.
func (s *Service) GetDataByDataID(ctx context.Context, req *fndv1.GetDataRequest) (*fndv1.TAFNDIShareFundFeeRdmRequest, error) {
	return &fndv1.TAFNDIShareFundFeeRdmRequest{}, nil
}

// ApproveTAFNDIShareFundFeeRdm approves or rejects a pending record for the 4-Eyes Principle.
func (s *Service) ApproveTAFNDIShareFundFeeRdm(ctx context.Context, req *fndv1.ApproveTAFNDIShareFundFeeRdmRequest) (*fndv1.SaveResponse, error) {
	s.log.Info("ApproveTAFNDIShareFundFeeRdm called", slog.String("data_id", req.GetDataId()))

	var record db.DTAFNDIShareFundFeeRdmEdit
	result := s.db.WithContext(ctx).Where("\"DataID\" = ?", req.GetDataId()).First(&record)
	
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return &fndv1.SaveResponse{Success: false, Message: "Record not found", ReturnCode: "NOT_FOUND"}, nil
		}
		s.log.Error("ApproveTAFNDIShareFundFeeRdm DB error", slog.Any("error", result.Error))
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
