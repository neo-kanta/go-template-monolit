package fndm003

import (
	"context"
	"log/slog"
	"strings"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
	"go-transfer-agent/services/fnd/fndm003/db"
	"go-transfer-agent/services/fnd/shared"

	"go-transfer-agent/common/platform/model"
	"gorm.io/gorm"
)

// Service holds the FNDM003 business logic.
type Service struct {
	db  *gorm.DB
	log *slog.Logger
}

// NewService creates a new FNDM003 service.
func NewService(database *gorm.DB, logger *slog.Logger) *Service {
	return &Service{
		db:  database,
		log: logger.With(slog.String("module", "fndm003-service")),
	}
}

// TAFNDSwitch retrieves Master and Detail structures from PostgreSQL via GORM.
func (s *Service) TAFNDSwitch(ctx context.Context, req *fndv1.TAFNDSwitchRequest) (*fndv1.TAFNDSwitchResponse, error) {
	sysCoID := req.GetSysCoId()
	prtFundCode := req.GetPrtFundCode()
	isSwitchIn := req.GetIsSwitchIn()

	resp := &fndv1.TAFNDSwitchResponse{
		ResultList:        make([]*fndv1.SwitchMasterDTO, 0),
		SwitchFundList:    make([]*fndv1.SwitchFundListDTO, 0),
		FundSwitchFeeList: make([]*fndv1.FundSwitchFeeListDTO, 0),
		SwitchCryList:     make([]*fndv1.SwitchCryListDTO, 0),
	}

	// 1. Query Master
	masterQuery := s.db.WithContext(ctx).Model(&db.DTAFNDSwitch{}).Where("SysCoID = ?", sysCoID)
	if prtFundCode != "" {
		masterQuery = masterQuery.Where("PrtFundCode = ?", prtFundCode)
	}
	if isSwitchIn != "" {
		masterQuery = masterQuery.Where("IsSwitchIn = ?", isSwitchIn)
	}

	var masters []db.DTAFNDSwitch
	if err := masterQuery.Find(&masters).Error; err != nil {
		s.log.Error("Failed to query DTAFNDSwitch", slog.Any("error", err))
		return nil, err
	}

	for _, m := range masters {
		dto := &fndv1.SwitchMasterDTO{
			SysCoId:       m.SysCoID,
			SysCoIdNm:     shared.GetSysCoName(m.SysCoID),
			PrtFundCode:   m.PrtFundCode,
			PrtFundCodeNm: shared.GetDFundName(m.PrtFundCode),
			IsSwitchIn:    m.IsSwitchIn,
		}
		resp.ResultList = append(resp.ResultList, dto)
	}

	// 2. Query SwitchFund
	fundQuery := s.db.WithContext(ctx).Model(&db.DTAFNDSwitchFund{}).Where("SysCoID = ?", sysCoID)
	if prtFundCode != "" {
		fundQuery = fundQuery.Where("PrtFundCode = ?", prtFundCode)
	}

	var switchFunds []db.DTAFNDSwitchFund
	if err := fundQuery.Find(&switchFunds).Error; err != nil {
		s.log.Error("Failed to query DTAFNDSwitchFund", slog.Any("error", err))
		return nil, err
	}

	for _, f := range switchFunds {
		dto := &fndv1.SwitchFundListDTO{
			SysCoId:          f.SysCoID,
			PrtFundCode:      f.PrtFundCode,
			SwOPrtFundCode:   f.SwOPrtFundCode,
			SwOPrtFundCodeNm: shared.GetDFundName(f.SwOPrtFundCode),
			SwDateType:       f.SwDateType,
		}
		resp.SwitchFundList = append(resp.SwitchFundList, dto)
	}

	// 3. Query FundSwitchFee
	feeQuery := s.db.WithContext(ctx).Model(&db.DTAFNDFundFeeSwitch{}).Where("SysCoID = ?", sysCoID)
	if prtFundCode != "" {
		feeQuery = feeQuery.Where("PrtFundCode = ?", prtFundCode)
	}

	var fees []db.DTAFNDFundFeeSwitch
	if err := feeQuery.Find(&fees).Error; err != nil {
		s.log.Error("Failed to query DTAFNDFundFeeSwitch", slog.Any("error", err))
		return nil, err
	}

	for _, f := range fees {
		dto := &fndv1.FundSwitchFeeListDTO{
			SysCoId:     f.SysCoID,
			PrtFundCode: f.PrtFundCode,
			FundCode:    f.FundCode,
			SwFundType:  f.SwFundType,
			SwDiscType:  f.SwDiscType,
			SwitchRate:  f.SwitchRate.InexactFloat64(),
		}
		resp.FundSwitchFeeList = append(resp.FundSwitchFeeList, dto)
	}

	// 4. Query SwitchCry
	cryQuery := s.db.WithContext(ctx).Model(&db.DTAFNDSwitchCry{}).Where("SysCoID = ?", sysCoID)
	if prtFundCode != "" {
		cryQuery = cryQuery.Where("PrtFundCode = ?", prtFundCode)
	}

	var cries []db.DTAFNDSwitchCry
	if err := cryQuery.Find(&cries).Error; err != nil {
		s.log.Error("Failed to query DTAFNDSwitchCry", slog.Any("error", err))
		return nil, err
	}

	for _, c := range cries {
		dto := &fndv1.SwitchCryListDTO{
			SysCoId:         c.SysCoID,
			PrtFundCode:     c.PrtFundCode,
			SwIFundCry:      c.SwIFundCry,
			SwIFundCryNm:    shared.GetCryName(c.SwIFundCry),
			SwOFundCrySet:   c.SwOFundCrySet,
			SwOFundCrySetNm: splitString(c.SwOFundCrySet), // mapping string slices
		}
		resp.SwitchCryList = append(resp.SwitchCryList, dto)
	}

	return resp, nil
}

func splitString(s string) []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(s, ",")
}

// SaveTAFNDSwitch creates a new record in the _Edit table — POST endpoint.
func (s *Service) SaveTAFNDSwitch(ctx context.Context, req *fndv1.TAFNDSwitchRequest) (*fndv1.SaveResponse, error) {
	s.log.Info("SaveTAFNDSwitch called")
	record := db.DTAFNDSwitchEdit{}
	// Auth Context & 4-Eyes Principle
	record.MakerID = shared.GetUsernameFromCtx(ctx)
	record.Status = models.StatusPendingApproval

	record.SysCoID = req.GetSysCoId()
	record.PrtFundCode = req.GetPrtFundCode()
	if err := s.db.WithContext(ctx).Create(&record).Error; err != nil {
		s.log.Error("SaveTAFNDSwitch failed", slog.Any("error", err))
		return &fndv1.SaveResponse{Success: false, Message: err.Error(), ReturnCode: "DB_ERROR"}, nil
	}
	return &fndv1.SaveResponse{Success: true, Message: "Record created successfully", ReturnCode: "0000"}, nil
}

// UpdateTAFNDSwitch updates an existing record in the _Edit table — PUT endpoint.
func (s *Service) UpdateTAFNDSwitch(ctx context.Context, req *fndv1.TAFNDSwitchRequest) (*fndv1.SaveResponse, error) {
	s.log.Info("UpdateTAFNDSwitch called")
	record := db.DTAFNDSwitchEdit{}
	record.SysCoID = req.GetSysCoId()
	record.PrtFundCode = req.GetPrtFundCode()
	result := s.db.WithContext(ctx).Model(&db.DTAFNDSwitchEdit{}).Where(`"DataID" = ?`, req.GetSysCoId()).Updates(&record)
	if result.Error != nil {
		s.log.Error("UpdateTAFNDSwitch failed", slog.Any("error", result.Error))
		return &fndv1.SaveResponse{Success: false, Message: result.Error.Error(), ReturnCode: "DB_ERROR"}, nil
	}
	if result.RowsAffected == 0 {
		return &fndv1.SaveResponse{Success: false, Message: "No record found to update", ReturnCode: "NOT_FOUND"}, nil
	}
	return &fndv1.SaveResponse{Success: true, Message: "Record updated successfully", ReturnCode: "0000"}, nil
}

// DeleteTAFNDSwitch removes a record from the _Edit table — DELETE endpoint.
func (s *Service) DeleteTAFNDSwitch(ctx context.Context, req *fndv1.DeleteRequest) (*fndv1.SaveResponse, error) {
	s.log.Info("DeleteTAFNDSwitch called",
		slog.String("data_id", req.GetDataId()),
		slog.String("data_flag", req.GetDataFlag()),
	)
	result := s.db.WithContext(ctx).Where(`"DataID" = ?`, req.GetDataId()).Delete(&db.DTAFNDSwitchEdit{})
	if result.Error != nil {
		s.log.Error("DeleteTAFNDSwitch failed", slog.Any("error", result.Error))
		return &fndv1.SaveResponse{Success: false, Message: result.Error.Error(), ReturnCode: "DB_ERROR"}, nil
	}
	if result.RowsAffected == 0 {
		return &fndv1.SaveResponse{Success: false, Message: "No record found with the given DataID", ReturnCode: "NOT_FOUND"}, nil
	}
	return &fndv1.SaveResponse{Success: true, Message: "Record deleted successfully", ReturnCode: "0000"}, nil
}

// GetDataByDataID retrieves full data by DataID.
func (s *Service) GetDataByDataID(ctx context.Context, req *fndv1.GetDataRequest) (*fndv1.TAFNDSwitchRequest, error) {
	return &fndv1.TAFNDSwitchRequest{}, nil
}

// ApproveTAFNDSwitch approves or rejects a pending record for the 4-Eyes Principle.
func (s *Service) ApproveTAFNDSwitch(ctx context.Context, req *fndv1.ApproveTAFNDSwitchRequest) (*fndv1.SaveResponse, error) {
	s.log.Info("ApproveTAFNDSwitch called", slog.String("data_id", req.GetDataId()))

	var record db.DTAFNDSwitchEdit
	result := s.db.WithContext(ctx).Where("\"DataID\" = ?", req.GetDataId()).First(&record)
	
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return &fndv1.SaveResponse{Success: false, Message: "Record not found", ReturnCode: "NOT_FOUND"}, nil
		}
		s.log.Error("ApproveTAFNDSwitch DB error", slog.Any("error", result.Error))
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
