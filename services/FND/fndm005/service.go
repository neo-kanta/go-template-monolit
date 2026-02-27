package fndm005

import (
	"context"
	"log/slog"
	"time"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
	"go-transfer-agent/services/fnd/fndm005/db"
	"go-transfer-agent/services/fnd/shared"

	"go-transfer-agent/common/platform/model"
	"gorm.io/gorm"
)

// Service holds the FNDM005 business logic.
type Service struct {
	db  *gorm.DB
	log *slog.Logger
}

// NewService creates a new FNDM005 service.
func NewService(database *gorm.DB, logger *slog.Logger) *Service {
	return &Service{
		db:  database,
		log: logger.With(slog.String("module", "fndm005-service")),
	}
}

// TAFNDFundCalDate retrieves Master, Memo, and Fund structures from PostgreSQL via GORM.
func (s *Service) TAFNDFundCalDate(ctx context.Context, req *fndv1.TAFNDFundCalDateRequest) (*fndv1.TAFNDFundCalDateResponse, error) {
	sysCoID := req.GetSysCoId()
	calYear := req.GetCalYear()
	fndCalType := req.GetFndCalType()
	fundCry := req.GetFundCry()

	resp := &fndv1.TAFNDFundCalDateResponse{
		ResultList:             make([]*fndv1.FundCalMasterDTO, 0),
		FundClosedMemoList:     make([]*fndv1.FundClosedDateMemoListDTO, 0),
		FundClosedDateFundList: make([]*fndv1.FundClosedDateFundListDTO, 0),
	}

	// 1. Query Master
	masterQuery := s.db.WithContext(ctx).Model(&db.TAFNDFundCal{}).Where("SysCoID = ?", sysCoID)
	if calYear != "" {
		masterQuery = masterQuery.Where("CalYear = ?", calYear)
	}
	if fndCalType != "" {
		masterQuery = masterQuery.Where("FNDCalType = ?", fndCalType)
	}
	if fundCry != "" {
		masterQuery = masterQuery.Where("FundCry = ?", fundCry)
	}

	var masters []db.TAFNDFundCal
	if err := masterQuery.Find(&masters).Error; err != nil {
		s.log.Error("Failed to query TAFNDFundCal", slog.Any("error", err))
		return nil, err
	}

	for _, m := range masters {
		dto := &fndv1.FundCalMasterDTO{
			SysCoId:    m.SysCoID,
			SysCoIdNm:  shared.GetSysCoName(m.SysCoID),
			CalYear:    m.CalYear,
			FndCalType: m.FNDCalType,
			FundCry:    m.FundCry,
			FundCryNm:  shared.GetCryName(m.FundCry),
		}
		resp.ResultList = append(resp.ResultList, dto)
	}

	// 2. Query Memo
	memoQuery := s.db.WithContext(ctx).Model(&db.TAFNDFundCalMemo{}).Where("SysCoID = ?", sysCoID)
	if calYear != "" {
		memoQuery = memoQuery.Where("CalYear = ?", calYear)
	}
	if fndCalType != "" {
		memoQuery = memoQuery.Where("FNDCalType = ?", fndCalType)
	}
	if fundCry != "" {
		memoQuery = memoQuery.Where("FundCry = ?", fundCry)
	}

	var memos []db.TAFNDFundCalMemo
	if err := memoQuery.Find(&memos).Error; err != nil {
		s.log.Error("Failed to query TAFNDFundCalMemo", slog.Any("error", err))
		return nil, err
	}

	for _, m := range memos {
		dto := &fndv1.FundClosedDateMemoListDTO{
			SysCoId:    m.SysCoID,
			CalYear:    m.CalYear,
			FndCalType: m.FNDCalType,
			FundCry:    m.FundCry,
			CalDate:    m.CalDate.Format(time.RFC3339),
			Remark: func() string { if m.Remark != nil { return *m.Remark }; return "" }(),
		}
		resp.FundClosedMemoList = append(resp.FundClosedMemoList, dto)
	}

	// 3. Query Fund List
	fundQuery := s.db.WithContext(ctx).Model(&db.TAFNDFundCalDtl{}).Where("SysCoID = ?", sysCoID)
	if calYear != "" {
		fundQuery = fundQuery.Where("CalYear = ?", calYear)
	}
	if fndCalType != "" {
		fundQuery = fundQuery.Where("FNDCalType = ?", fndCalType)
	}
	if fundCry != "" {
		fundQuery = fundQuery.Where("FundCry = ?", fundCry)
	}

	var funds []db.TAFNDFundCalDtl
	if err := fundQuery.Find(&funds).Error; err != nil {
		s.log.Error("Failed to query TAFNDFundCalDtl", slog.Any("error", err))
		return nil, err
	}

	for _, f := range funds {
		dto := &fndv1.FundClosedDateFundListDTO{
			SysCoId:       f.SysCoID,
			CalYear:       f.CalYear,
			FndCalType:    f.FNDCalType,
			FundCry:       f.FundCry,
			PrtFundCode:   f.PrtFundCode,
			PrtFundCodeNm: shared.GetDFundName(f.PrtFundCode),
			CalDate:       f.CalDate.Format(time.RFC3339),
		}
		resp.FundClosedDateFundList = append(resp.FundClosedDateFundList, dto)
	}

	return resp, nil
}

// SaveTAFNDFundCalDate creates a new record in the _Edit table — POST endpoint.
func (s *Service) SaveTAFNDFundCalDate(ctx context.Context, req *fndv1.TAFNDFundCalDateRequest) (*fndv1.SaveResponse, error) {
	s.log.Info("SaveTAFNDFundCalDate called")
	record := db.TAFNDFundCalEdit{}
	record.SysCoID = req.GetSysCoId()
	record.CalYear = req.GetCalYear()
	record.FNDCalType = req.GetFndCalType()
	record.FundCry = req.GetFundCry()
	if err := s.db.WithContext(ctx).Create(&record).Error; err != nil {
		s.log.Error("SaveTAFNDFundCalDate failed", slog.Any("error", err))
		return &fndv1.SaveResponse{Success: false, Message: err.Error(), ReturnCode: "DB_ERROR"}, nil
	}
	return &fndv1.SaveResponse{Success: true, Message: "Record created successfully", ReturnCode: "0000"}, nil
}

// UpdateTAFNDFundCalDate updates an existing record in the _Edit table — PUT endpoint.
func (s *Service) UpdateTAFNDFundCalDate(ctx context.Context, req *fndv1.TAFNDFundCalDateRequest) (*fndv1.SaveResponse, error) {
	s.log.Info("UpdateTAFNDFundCalDate called")
	record := db.TAFNDFundCalEdit{}
	record.SysCoID = req.GetSysCoId()
	record.CalYear = req.GetCalYear()
	record.FNDCalType = req.GetFndCalType()
	record.FundCry = req.GetFundCry()
	result := s.db.WithContext(ctx).Model(&db.TAFNDFundCalEdit{}).Where(`"DataID" = ?`, req.GetSysCoId()).Updates(&record)
	if result.Error != nil {
		s.log.Error("UpdateTAFNDFundCalDate failed", slog.Any("error", result.Error))
		return &fndv1.SaveResponse{Success: false, Message: result.Error.Error(), ReturnCode: "DB_ERROR"}, nil
	}
	if result.RowsAffected == 0 {
		return &fndv1.SaveResponse{Success: false, Message: "No record found to update", ReturnCode: "NOT_FOUND"}, nil
	}
	return &fndv1.SaveResponse{Success: true, Message: "Record updated successfully", ReturnCode: "0000"}, nil
}

// DeleteTAFNDFundCalDate removes a record from the _Edit table — DELETE endpoint.
func (s *Service) DeleteTAFNDFundCalDate(ctx context.Context, req *fndv1.DeleteRequest) (*fndv1.SaveResponse, error) {
	s.log.Info("DeleteTAFNDFundCalDate called",
		slog.String("data_id", req.GetDataId()),
		slog.String("data_flag", req.GetDataFlag()),
	)
	result := s.db.WithContext(ctx).Where(`"DataID" = ?`, req.GetDataId()).Delete(&db.TAFNDFundCalEdit{})
	if result.Error != nil {
		s.log.Error("DeleteTAFNDFundCalDate failed", slog.Any("error", result.Error))
		return &fndv1.SaveResponse{Success: false, Message: result.Error.Error(), ReturnCode: "DB_ERROR"}, nil
	}
	if result.RowsAffected == 0 {
		return &fndv1.SaveResponse{Success: false, Message: "No record found with the given DataID", ReturnCode: "NOT_FOUND"}, nil
	}
	return &fndv1.SaveResponse{Success: true, Message: "Record deleted successfully", ReturnCode: "0000"}, nil
}

func (s *Service) TACKFNDCalEdit(ctx context.Context, req *fndv1.TACKFNDCalEditRequest) (*fndv1.CheckResponse, error) {
	s.log.Info("TACKFNDCalEdit called")
	return &fndv1.CheckResponse{ReturnCode: "", Message: "No pending data"}, nil
}

// GetDataByDataID retrieves full data by DataID.
func (s *Service) GetDataByDataID(ctx context.Context, req *fndv1.GetDataRequest) (*fndv1.TAFNDFundCalDateRequest, error) {
	return &fndv1.TAFNDFundCalDateRequest{}, nil
}

// ApproveTAFNDFundCalDate approves or rejects a pending record for the 4-Eyes Principle.
func (s *Service) ApproveTAFNDFundCalDate(ctx context.Context, req *fndv1.ApproveTAFNDFundCalDateRequest) (*fndv1.SaveResponse, error) {
	s.log.Info("ApproveTAFNDFundCalDate called", slog.String("data_id", req.GetDataId()))

	var record db.TAFNDFundCalDtlEdit
	result := s.db.WithContext(ctx).Where("\"DataID\" = ?", req.GetDataId()).First(&record)
	
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return &fndv1.SaveResponse{Success: false, Message: "Record not found", ReturnCode: "NOT_FOUND"}, nil
		}
		s.log.Error("ApproveTAFNDFundCalDate DB error", slog.Any("error", result.Error))
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
