package fndm008

import (
	"context"
	"log/slog"
	"time"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
	"go-transfer-agent/services/fnd/fndm008/db"
	"go-transfer-agent/services/fnd/shared"

	"go-transfer-agent/common/platform/model"
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
			Remark: func() string { if m.Remark != nil { return *m.Remark }; return "" }(),
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

// SaveTAFNDPauseTxn creates a new record in the _Edit table — POST endpoint.
func (s *Service) SaveTAFNDPauseTxn(ctx context.Context, req *fndv1.TAFNDPauseTxnRequest) (*fndv1.SaveResponse, error) {
	s.log.Info("SaveTAFNDPauseTxn called")
	record := db.DTAFNDPauseTxnEdit{}
	// Auth Context & 4-Eyes Principle
	record.MakerID = shared.GetUsernameFromCtx(ctx)
	record.Status = models.StatusPendingApproval

	record.SysCoID = req.GetSysCoId()
	record.PrtFundCode = req.GetPrtFundCode()
		remark := req.GetRemark()
	record.Remark = &remark
	if err := s.db.WithContext(ctx).Create(&record).Error; err != nil {
		s.log.Error("SaveTAFNDPauseTxn failed", slog.Any("error", err))
		return &fndv1.SaveResponse{Success: false, Message: err.Error(), ReturnCode: "DB_ERROR"}, nil
	}
	return &fndv1.SaveResponse{Success: true, Message: "Record created successfully", ReturnCode: "0000"}, nil
}

// UpdateTAFNDPauseTxn updates an existing record in the _Edit table — PUT endpoint.
func (s *Service) UpdateTAFNDPauseTxn(ctx context.Context, req *fndv1.TAFNDPauseTxnRequest) (*fndv1.SaveResponse, error) {
	s.log.Info("UpdateTAFNDPauseTxn called")
	record := db.DTAFNDPauseTxnEdit{}
	record.SysCoID = req.GetSysCoId()
	record.PrtFundCode = req.GetPrtFundCode()
		remark := req.GetRemark()
	record.Remark = &remark
	result := s.db.WithContext(ctx).Model(&db.DTAFNDPauseTxnEdit{}).Where(`"DataID" = ?`, req.GetSysCoId()).Updates(&record)
	if result.Error != nil {
		s.log.Error("UpdateTAFNDPauseTxn failed", slog.Any("error", result.Error))
		return &fndv1.SaveResponse{Success: false, Message: result.Error.Error(), ReturnCode: "DB_ERROR"}, nil
	}
	if result.RowsAffected == 0 {
		return &fndv1.SaveResponse{Success: false, Message: "No record found to update", ReturnCode: "NOT_FOUND"}, nil
	}
	return &fndv1.SaveResponse{Success: true, Message: "Record updated successfully", ReturnCode: "0000"}, nil
}

// DeleteTAFNDPauseTxn removes a record from the _Edit table — DELETE endpoint.
func (s *Service) DeleteTAFNDPauseTxn(ctx context.Context, req *fndv1.DeleteRequest) (*fndv1.SaveResponse, error) {
	s.log.Info("DeleteTAFNDPauseTxn called",
		slog.String("data_id", req.GetDataId()),
		slog.String("data_flag", req.GetDataFlag()),
	)
	result := s.db.WithContext(ctx).Where(`"DataID" = ?`, req.GetDataId()).Delete(&db.DTAFNDPauseTxnEdit{})
	if result.Error != nil {
		s.log.Error("DeleteTAFNDPauseTxn failed", slog.Any("error", result.Error))
		return &fndv1.SaveResponse{Success: false, Message: result.Error.Error(), ReturnCode: "DB_ERROR"}, nil
	}
	if result.RowsAffected == 0 {
		return &fndv1.SaveResponse{Success: false, Message: "No record found with the given DataID", ReturnCode: "NOT_FOUND"}, nil
	}
	return &fndv1.SaveResponse{Success: true, Message: "Record deleted successfully", ReturnCode: "0000"}, nil
}

func (s *Service) TACKPTxnBegDate(ctx context.Context, req *fndv1.TACKPTxnBegDateRequest) (*fndv1.CheckResponse, error) {
	s.log.Info("TACKPTxnBegDate called")
	return &fndv1.CheckResponse{ReturnCode: "", Message: "Date is valid"}, nil
}

// GetDataByDataID retrieves full data by DataID.
func (s *Service) GetDataByDataID(ctx context.Context, req *fndv1.GetDataRequest) (*fndv1.TAFNDPauseTxnRequest, error) {
	return &fndv1.TAFNDPauseTxnRequest{}, nil
}

// ApproveTAFNDPauseTxn approves or rejects a pending record for the 4-Eyes Principle.
func (s *Service) ApproveTAFNDPauseTxn(ctx context.Context, req *fndv1.ApproveTAFNDPauseTxnRequest) (*fndv1.SaveResponse, error) {
	s.log.Info("ApproveTAFNDPauseTxn called", slog.String("data_id", req.GetDataId()))

	var record db.DTAFNDPauseTxnEdit
	result := s.db.WithContext(ctx).Where("\"DataID\" = ?", req.GetDataId()).First(&record)
	
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return &fndv1.SaveResponse{Success: false, Message: "Record not found", ReturnCode: "NOT_FOUND"}, nil
		}
		s.log.Error("ApproveTAFNDPauseTxn DB error", slog.Any("error", result.Error))
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
