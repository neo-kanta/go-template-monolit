package fndm007

import (
	"context"
	"log/slog"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
	"go-transfer-agent/services/fnd/fndm007/db"
	"go-transfer-agent/services/fnd/shared"

	"go-transfer-agent/common/platform/model"
	"gorm.io/gorm"
)

// Service holds the FNDM007 business logic.
type Service struct {
	db  *gorm.DB
	log *slog.Logger
}

// NewService creates a new FNDM007 service.
func NewService(database *gorm.DB, logger *slog.Logger) *Service {
	return &Service{
		db:  database,
		log: logger.With(slog.String("module", "fndm007-service")),
	}
}

// TAFNDCustGroup retrieves Master and Detail structures from PostgreSQL via GORM.
func (s *Service) TAFNDCustGroup(ctx context.Context, req *fndv1.TAFNDCustGroupRequest) (*fndv1.TAFNDCustGroupResponse, error) {
	sysCoID := req.GetSysCoId()
	custGrpCode := req.GetCustGrpCode()
	custGrpMName := req.GetCustGrpMName()
	custGrpSName := req.GetCustGrpSName()

	resp := &fndv1.TAFNDCustGroupResponse{
		ResultList:        make([]*fndv1.CustGroupMasterDTO, 0),
		CustFundGroupList: make([]*fndv1.CustFundGroupListDTO, 0),
	}

	// 1. Query Master
	masterQuery := s.db.WithContext(ctx).Model(&db.TAFNDCustGroup{}).Where("SysCoID = ?", sysCoID)
	if custGrpCode != "" {
		masterQuery = masterQuery.Where("CustGrpCode = ?", custGrpCode)
	}
	if custGrpMName != "" {
		masterQuery = masterQuery.Where("CustGrpMName LIKE ?", "%"+custGrpMName+"%")
	}
	if custGrpSName != "" {
		masterQuery = masterQuery.Where("CustGrpSName LIKE ?", "%"+custGrpSName+"%")
	}

	var masters []db.TAFNDCustGroup
	if err := masterQuery.Find(&masters).Error; err != nil {
		s.log.Error("Failed to query TAFNDCustGroup", slog.Any("error", err))
		return nil, err
	}

	for _, m := range masters {
		dto := &fndv1.CustGroupMasterDTO{
			SysCoId:      m.SysCoID,
			SysCoIdNm:    shared.GetSysCoName(m.SysCoID),
			CustGrpCode:  m.CustGrpCode,
			CustGrpMName: m.CustGrpMName,
			CustGrpSName: m.CustGrpSName,
		}
		resp.ResultList = append(resp.ResultList, dto)
	}

	// 2. Query Detail
	dtlQuery := s.db.WithContext(ctx).Model(&db.TAFNDCustGroupDtl{}).Where("SysCoID = ?", sysCoID)
	if custGrpCode != "" {
		dtlQuery = dtlQuery.Where("CustGrpCode = ?", custGrpCode)
	}

	var details []db.TAFNDCustGroupDtl
	if err := dtlQuery.Find(&details).Error; err != nil {
		s.log.Error("Failed to query TAFNDCustGroupDtl", slog.Any("error", err))
		return nil, err
	}

	for _, d := range details {
		dto := &fndv1.CustFundGroupListDTO{
			SysCoId:       d.SysCoID,
			CustGrpCode:   d.CustGrpCode,
			PrtFundCode:   d.PrtFundCode,
			PrtFundCodeNm: shared.GetDFundName(d.PrtFundCode),
		}
		resp.CustFundGroupList = append(resp.CustFundGroupList, dto)
	}

	return resp, nil
}

// SaveTAFNDCustGroup creates a new record in the _Edit table — POST endpoint.
func (s *Service) SaveTAFNDCustGroup(ctx context.Context, req *fndv1.TAFNDCustGroupRequest) (*fndv1.SaveResponse, error) {
	s.log.Info("SaveTAFNDCustGroup called")
	record := db.TAFNDCustGroupEdit{}
	record.SysCoID = req.GetSysCoId()
	record.CustGrpCode = req.GetCustGrpCode()
	record.CustGrpMName = req.GetCustGrpMName()
	record.CustGrpSName = req.GetCustGrpSName()
	if err := s.db.WithContext(ctx).Create(&record).Error; err != nil {
		s.log.Error("SaveTAFNDCustGroup failed", slog.Any("error", err))
		return &fndv1.SaveResponse{Success: false, Message: err.Error(), ReturnCode: "DB_ERROR"}, nil
	}
	return &fndv1.SaveResponse{Success: true, Message: "Record created successfully", ReturnCode: "0000"}, nil
}

// UpdateTAFNDCustGroup updates an existing record in the _Edit table — PUT endpoint.
func (s *Service) UpdateTAFNDCustGroup(ctx context.Context, req *fndv1.TAFNDCustGroupRequest) (*fndv1.SaveResponse, error) {
	s.log.Info("UpdateTAFNDCustGroup called")
	record := db.TAFNDCustGroupEdit{}
	record.SysCoID = req.GetSysCoId()
	record.CustGrpCode = req.GetCustGrpCode()
	record.CustGrpMName = req.GetCustGrpMName()
	record.CustGrpSName = req.GetCustGrpSName()
	result := s.db.WithContext(ctx).Model(&db.TAFNDCustGroupEdit{}).Where(`"DataID" = ?`, req.GetSysCoId()).Updates(&record)
	if result.Error != nil {
		s.log.Error("UpdateTAFNDCustGroup failed", slog.Any("error", result.Error))
		return &fndv1.SaveResponse{Success: false, Message: result.Error.Error(), ReturnCode: "DB_ERROR"}, nil
	}
	if result.RowsAffected == 0 {
		return &fndv1.SaveResponse{Success: false, Message: "No record found to update", ReturnCode: "NOT_FOUND"}, nil
	}
	return &fndv1.SaveResponse{Success: true, Message: "Record updated successfully", ReturnCode: "0000"}, nil
}

// DeleteTAFNDCustGroup removes a record from the _Edit table — DELETE endpoint.
func (s *Service) DeleteTAFNDCustGroup(ctx context.Context, req *fndv1.DeleteRequest) (*fndv1.SaveResponse, error) {
	s.log.Info("DeleteTAFNDCustGroup called",
		slog.String("data_id", req.GetDataId()),
		slog.String("data_flag", req.GetDataFlag()),
	)
	result := s.db.WithContext(ctx).Where(`"DataID" = ?`, req.GetDataId()).Delete(&db.TAFNDCustGroupEdit{})
	if result.Error != nil {
		s.log.Error("DeleteTAFNDCustGroup failed", slog.Any("error", result.Error))
		return &fndv1.SaveResponse{Success: false, Message: result.Error.Error(), ReturnCode: "DB_ERROR"}, nil
	}
	if result.RowsAffected == 0 {
		return &fndv1.SaveResponse{Success: false, Message: "No record found with the given DataID", ReturnCode: "NOT_FOUND"}, nil
	}
	return &fndv1.SaveResponse{Success: true, Message: "Record deleted successfully", ReturnCode: "0000"}, nil
}

func (s *Service) TAGetCustGroup(ctx context.Context, req *fndv1.TAGetCustGroupRequest) (*fndv1.TAGetCustGroupResponse, error) {
	s.log.Info("TAGetCustGroup called")
	return &fndv1.TAGetCustGroupResponse{}, nil
}

// GetDataByDataID retrieves full data by DataID.
func (s *Service) GetDataByDataID(ctx context.Context, req *fndv1.GetDataRequest) (*fndv1.TAFNDCustGroupRequest, error) {
	return &fndv1.TAFNDCustGroupRequest{}, nil
}

// ApproveTAFNDCustGroup approves or rejects a pending record for the 4-Eyes Principle.
func (s *Service) ApproveTAFNDCustGroup(ctx context.Context, req *fndv1.ApproveTAFNDCustGroupRequest) (*fndv1.SaveResponse, error) {
	s.log.Info("ApproveTAFNDCustGroup called", slog.String("data_id", req.GetDataId()))

	var record db.TAFNDCustGroupEdit
	result := s.db.WithContext(ctx).Where("\"DataID\" = ?", req.GetDataId()).First(&record)
	
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return &fndv1.SaveResponse{Success: false, Message: "Record not found", ReturnCode: "NOT_FOUND"}, nil
		}
		s.log.Error("ApproveTAFNDCustGroup DB error", slog.Any("error", result.Error))
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
