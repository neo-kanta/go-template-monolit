package fndm006

import (
	"context"
	"log/slog"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
	"go-transfer-agent/services/fnd/fndm006/db"
	"go-transfer-agent/services/fnd/shared"

	"go-transfer-agent/common/platform/model"
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

// SaveTAFNDRPFeeChgType creates a new record in the _Edit table — POST endpoint.
func (s *Service) SaveTAFNDRPFeeChgType(ctx context.Context, req *fndv1.TAFNDRPFeeChgTypeRequest) (*fndv1.SaveResponse, error) {
	s.log.Info("SaveTAFNDRPFeeChgType called")
	record := db.DTAFNDRPFeeChgTypeEdit{}
	// Auth Context & 4-Eyes Principle
	record.MakerID = shared.GetUsernameFromCtx(ctx)
	record.Status = models.StatusPendingApproval

	record.SysCoID = req.GetSysCoId()
	record.PrtFundCode = req.GetPrtFundCode()
	if err := s.db.WithContext(ctx).Create(&record).Error; err != nil {
		s.log.Error("SaveTAFNDRPFeeChgType failed", slog.Any("error", err))
		return &fndv1.SaveResponse{Success: false, Message: err.Error(), ReturnCode: "DB_ERROR"}, nil
	}
	return &fndv1.SaveResponse{Success: true, Message: "Record created successfully", ReturnCode: "0000"}, nil
}

// UpdateTAFNDRPFeeChgType updates an existing record in the _Edit table — PUT endpoint.
func (s *Service) UpdateTAFNDRPFeeChgType(ctx context.Context, req *fndv1.TAFNDRPFeeChgTypeRequest) (*fndv1.SaveResponse, error) {
	s.log.Info("UpdateTAFNDRPFeeChgType called")
	record := db.DTAFNDRPFeeChgTypeEdit{}
	record.SysCoID = req.GetSysCoId()
	record.PrtFundCode = req.GetPrtFundCode()
	result := s.db.WithContext(ctx).Model(&db.DTAFNDRPFeeChgTypeEdit{}).Where(`"DataID" = ?`, req.GetSysCoId()).Updates(&record)
	if result.Error != nil {
		s.log.Error("UpdateTAFNDRPFeeChgType failed", slog.Any("error", result.Error))
		return &fndv1.SaveResponse{Success: false, Message: result.Error.Error(), ReturnCode: "DB_ERROR"}, nil
	}
	if result.RowsAffected == 0 {
		return &fndv1.SaveResponse{Success: false, Message: "No record found to update", ReturnCode: "NOT_FOUND"}, nil
	}
	return &fndv1.SaveResponse{Success: true, Message: "Record updated successfully", ReturnCode: "0000"}, nil
}

// DeleteTAFNDRPFeeChgType removes a record from the _Edit table — DELETE endpoint.
func (s *Service) DeleteTAFNDRPFeeChgType(ctx context.Context, req *fndv1.DeleteRequest) (*fndv1.SaveResponse, error) {
	s.log.Info("DeleteTAFNDRPFeeChgType called",
		slog.String("data_id", req.GetDataId()),
		slog.String("data_flag", req.GetDataFlag()),
	)
	result := s.db.WithContext(ctx).Where(`"DataID" = ?`, req.GetDataId()).Delete(&db.DTAFNDRPFeeChgTypeEdit{})
	if result.Error != nil {
		s.log.Error("DeleteTAFNDRPFeeChgType failed", slog.Any("error", result.Error))
		return &fndv1.SaveResponse{Success: false, Message: result.Error.Error(), ReturnCode: "DB_ERROR"}, nil
	}
	if result.RowsAffected == 0 {
		return &fndv1.SaveResponse{Success: false, Message: "No record found with the given DataID", ReturnCode: "NOT_FOUND"}, nil
	}
	return &fndv1.SaveResponse{Success: true, Message: "Record deleted successfully", ReturnCode: "0000"}, nil
}

// GetDataByDataID retrieves full data by DataID.
func (s *Service) GetDataByDataID(ctx context.Context, req *fndv1.GetDataRequest) (*fndv1.TAFNDRPFeeChgTypeRequest, error) {
	return &fndv1.TAFNDRPFeeChgTypeRequest{}, nil
}

// ApproveTAFNDRPFeeChgType approves or rejects a pending record for the 4-Eyes Principle.
func (s *Service) ApproveTAFNDRPFeeChgType(ctx context.Context, req *fndv1.ApproveTAFNDRPFeeChgTypeRequest) (*fndv1.SaveResponse, error) {
	s.log.Info("ApproveTAFNDRPFeeChgType called", slog.String("data_id", req.GetDataId()))

	var record db.DTAFNDRPFeeChgTypeEdit
	result := s.db.WithContext(ctx).Where("\"DataID\" = ?", req.GetDataId()).First(&record)
	
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return &fndv1.SaveResponse{Success: false, Message: "Record not found", ReturnCode: "NOT_FOUND"}, nil
		}
		s.log.Error("ApproveTAFNDRPFeeChgType DB error", slog.Any("error", result.Error))
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
