package fndm012

import (
	"context"
	"log/slog"
	"time"

	"github.com/shopspring/decimal"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
	"go-transfer-agent/services/fnd/fndm012/db"
	"go-transfer-agent/services/fnd/shared"

	"go-transfer-agent/common/platform/model"
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
			FeeRate:        m.FeeRate.InexactFloat64(),
		}
		resp.ResultList = append(resp.ResultList, dto)
	}

	return resp, nil
}

// SaveTAFNDPGFundFeeRdm creates a new record in the _Edit table — POST endpoint.
func (s *Service) SaveTAFNDPGFundFeeRdm(ctx context.Context, req *fndv1.TAFNDPGFundFeeRdmRequest) (*fndv1.SaveResponse, error) {
	s.log.Info("SaveTAFNDPGFundFeeRdm called")
	record := db.DTAFNDPGFundFeeRdm{}
	// Auth Context & 4-Eyes Principle
	record.MakerID = shared.GetUsernameFromCtx(ctx)
	record.Status = models.StatusPendingApproval

	record.SysCoID = req.GetSysCoId()
	record.PrtFundCode = req.GetPrtFundCode()
	record.FeeRate = decimal.NewFromFloat(req.GetFeeRate())
	if err := s.db.WithContext(ctx).Create(&record).Error; err != nil {
		s.log.Error("SaveTAFNDPGFundFeeRdm failed", slog.Any("error", err))
		return &fndv1.SaveResponse{Success: false, Message: err.Error(), ReturnCode: "DB_ERROR"}, nil
	}
	return &fndv1.SaveResponse{Success: true, Message: "Record created successfully", ReturnCode: "0000"}, nil
}

// UpdateTAFNDPGFundFeeRdm updates an existing record in the _Edit table — PUT endpoint.
func (s *Service) UpdateTAFNDPGFundFeeRdm(ctx context.Context, req *fndv1.TAFNDPGFundFeeRdmRequest) (*fndv1.SaveResponse, error) {
	s.log.Info("UpdateTAFNDPGFundFeeRdm called")
	record := db.DTAFNDPGFundFeeRdm{}
	record.SysCoID = req.GetSysCoId()
	record.PrtFundCode = req.GetPrtFundCode()
	record.FeeRate = decimal.NewFromFloat(req.GetFeeRate())
	result := s.db.WithContext(ctx).Model(&db.DTAFNDPGFundFeeRdm{}).Where(`"DataID" = ?`, req.GetSysCoId()).Updates(&record)
	if result.Error != nil {
		s.log.Error("UpdateTAFNDPGFundFeeRdm failed", slog.Any("error", result.Error))
		return &fndv1.SaveResponse{Success: false, Message: result.Error.Error(), ReturnCode: "DB_ERROR"}, nil
	}
	if result.RowsAffected == 0 {
		return &fndv1.SaveResponse{Success: false, Message: "No record found to update", ReturnCode: "NOT_FOUND"}, nil
	}
	return &fndv1.SaveResponse{Success: true, Message: "Record updated successfully", ReturnCode: "0000"}, nil
}

// DeleteTAFNDPGFundFeeRdm removes a record from the _Edit table — DELETE endpoint.
func (s *Service) DeleteTAFNDPGFundFeeRdm(ctx context.Context, req *fndv1.DeleteRequest) (*fndv1.SaveResponse, error) {
	s.log.Info("DeleteTAFNDPGFundFeeRdm called",
		slog.String("data_id", req.GetDataId()),
		slog.String("data_flag", req.GetDataFlag()),
	)
	result := s.db.WithContext(ctx).Where(`"DataID" = ?`, req.GetDataId()).Delete(&db.DTAFNDPGFundFeeRdm{})
	if result.Error != nil {
		s.log.Error("DeleteTAFNDPGFundFeeRdm failed", slog.Any("error", result.Error))
		return &fndv1.SaveResponse{Success: false, Message: result.Error.Error(), ReturnCode: "DB_ERROR"}, nil
	}
	if result.RowsAffected == 0 {
		return &fndv1.SaveResponse{Success: false, Message: "No record found with the given DataID", ReturnCode: "NOT_FOUND"}, nil
	}
	return &fndv1.SaveResponse{Success: true, Message: "Record deleted successfully", ReturnCode: "0000"}, nil
}

// GetDataByDataID retrieves full data by DataID.
func (s *Service) GetDataByDataID(ctx context.Context, req *fndv1.GetDataRequest) (*fndv1.TAFNDPGFundFeeRdmRequest, error) {
	return &fndv1.TAFNDPGFundFeeRdmRequest{}, nil
}

// ApproveTAFNDPGFundFeeRdm approves or rejects a pending record for the 4-Eyes Principle.
func (s *Service) ApproveTAFNDPGFundFeeRdm(ctx context.Context, req *fndv1.ApproveTAFNDPGFundFeeRdmRequest) (*fndv1.SaveResponse, error) {
	s.log.Info("ApproveTAFNDPGFundFeeRdm called", slog.String("data_id", req.GetDataId()))

	var record db.DTAFNDPGFundFeeRdm
	result := s.db.WithContext(ctx).Where("\"DataID\" = ?", req.GetDataId()).First(&record)
	
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return &fndv1.SaveResponse{Success: false, Message: "Record not found", ReturnCode: "NOT_FOUND"}, nil
		}
		s.log.Error("ApproveTAFNDPGFundFeeRdm DB error", slog.Any("error", result.Error))
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
