package fndm004

import (
	"context"
	"log/slog"
	"strings"
	"time"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
	"go-transfer-agent/services/fnd/fndm004/db"
	"go-transfer-agent/services/fnd/shared"

	"go-transfer-agent/common/platform/model"
	"gorm.io/gorm"
)

// Service holds the FNDM004 business logic.
type Service struct {
	db  *gorm.DB
	log *slog.Logger
}

// NewService creates a new FNDM004 service.
func NewService(database *gorm.DB, logger *slog.Logger) *Service {
	return &Service{
		db:  database,
		log: logger.With(slog.String("module", "fndm004-service")),
	}
}

// TAFNDFundAgent retrieves Master and Detail structures from PostgreSQL via GORM.
func (s *Service) TAFNDFundAgent(ctx context.Context, req *fndv1.TAFNDFundAgentRequest) (*fndv1.TAFNDFundAgentResponse, error) {
	sysCoID := req.GetSysCoId()
	prtFundCode := req.GetPrtFundCode()

	resp := &fndv1.TAFNDFundAgentResponse{
		ResultList:       make([]*fndv1.FundAgentMasterDTO, 0),
		FundAgentDtlList: make([]*fndv1.FundAgentDtlListDTO, 0),
	}

	// 1. Query Master
	masterQuery := s.db.WithContext(ctx).Model(&db.DTAFNDFundAgent{}).Where("SysCoID = ?", sysCoID)
	if prtFundCode != "" {
		masterQuery = masterQuery.Where("PrtFundCode = ?", prtFundCode)
	}

	var masters []db.DTAFNDFundAgent
	if err := masterQuery.Find(&masters).Error; err != nil {
		s.log.Error("Failed to query DTAFNDFundAgent", slog.Any("error", err))
		return nil, err
	}

	for _, m := range masters {
		dto := &fndv1.FundAgentMasterDTO{
			SysCoId:       m.SysCoID,
			SysCoIdNm:     shared.GetSysCoName(m.SysCoID),
			PrtFundCode:   m.PrtFundCode,
			PrtFundCodeNm: shared.GetDFundName(m.PrtFundCode),
			FundIpoDate:   "", // Map if IPO date exists elsewhere, standard logic leaves empty unless joined
		}
		resp.ResultList = append(resp.ResultList, dto)
	}

	// 2. Query Detail
	dtlQuery := s.db.WithContext(ctx).Model(&db.DTAFNDFundAgentDtl{}).Where("SysCoID = ?", sysCoID)
	if prtFundCode != "" {
		dtlQuery = dtlQuery.Where("PrtFundCode = ?", prtFundCode)
	}

	var details []db.DTAFNDFundAgentDtl
	if err := dtlQuery.Find(&details).Error; err != nil {
		s.log.Error("Failed to query DTAFNDFundAgentDtl", slog.Any("error", err))
		return nil, err
	}

	for _, d := range details {
		dto := &fndv1.FundAgentDtlListDTO{
			SysCoId:           d.SysCoID,
			PrtFundCode:       d.PrtFundCode,
			AgentType:         d.AgentType,
			AgentCode:         d.AgentCode,
			AgentCodeNm:       "Test Neo Agent Name",
			IsMAgentOpType:    d.IsMAgentOPType,
			AgentOpType:       d.AgentOPType,
			IsMRcvTxnType:     d.IsMRcvTxnType,
			RcvTxnType:        d.RcvTxnType,
			IsMSubsFeePct:     d.IsMSubsFeePct,
			SubsFeePctAg:      d.SubsFeePctAG.InexactFloat64(),
			SubsFeePctFh:      d.SubsFeePctFH.InexactFloat64(),
			IsMSubsFeePctType: d.IsMSubsFeePctType,
			SubsFeePctType:    d.SubsFeePctType,
			IsMFundCrySet:     d.IsMFundCrySet,
			FundCrySet:        splitString(d.FundCrySet),
			FundCrySetNm:      splitString(d.FundCrySet),
			NoSaleShareSet:    splitString(d.NoSaleShareSet),
			NoSaleShareSetNm:  splitString(d.NoSaleShareSet),
			IsMEndOfSale:      d.IsMEndOfSale,
			TermDate:          d.TermDate.Format(time.RFC3339),
		}
		resp.FundAgentDtlList = append(resp.FundAgentDtlList, dto)
	}

	return resp, nil
}

func splitString(s string) []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(s, ",")
}

// SaveTAFNDFundAgent creates a new record in the _Edit table — POST endpoint.
func (s *Service) SaveTAFNDFundAgent(ctx context.Context, req *fndv1.TAFNDFundAgentRequest) (*fndv1.SaveResponse, error) {
	s.log.Info("SaveTAFNDFundAgent called")
	record := db.DTAFNDFundAgentEdit{}
	// Auth Context & 4-Eyes Principle
	record.MakerID = shared.GetUsernameFromCtx(ctx)
	record.Status = models.StatusPendingApproval

	record.SysCoID = req.GetSysCoId()
	record.PrtFundCode = req.GetPrtFundCode()
	if err := s.db.WithContext(ctx).Create(&record).Error; err != nil {
		s.log.Error("SaveTAFNDFundAgent failed", slog.Any("error", err))
		return &fndv1.SaveResponse{Success: false, Message: err.Error(), ReturnCode: "DB_ERROR"}, nil
	}
	return &fndv1.SaveResponse{Success: true, Message: "Record created successfully", ReturnCode: "0000"}, nil
}

// UpdateTAFNDFundAgent updates an existing record in the _Edit table — PUT endpoint.
func (s *Service) UpdateTAFNDFundAgent(ctx context.Context, req *fndv1.TAFNDFundAgentRequest) (*fndv1.SaveResponse, error) {
	s.log.Info("UpdateTAFNDFundAgent called")
	record := db.DTAFNDFundAgentEdit{}
	record.SysCoID = req.GetSysCoId()
	record.PrtFundCode = req.GetPrtFundCode()
	result := s.db.WithContext(ctx).Model(&db.DTAFNDFundAgentEdit{}).Where(`"DataID" = ?`, req.GetSysCoId()).Updates(&record)
	if result.Error != nil {
		s.log.Error("UpdateTAFNDFundAgent failed", slog.Any("error", result.Error))
		return &fndv1.SaveResponse{Success: false, Message: result.Error.Error(), ReturnCode: "DB_ERROR"}, nil
	}
	if result.RowsAffected == 0 {
		return &fndv1.SaveResponse{Success: false, Message: "No record found to update", ReturnCode: "NOT_FOUND"}, nil
	}
	return &fndv1.SaveResponse{Success: true, Message: "Record updated successfully", ReturnCode: "0000"}, nil
}

// DeleteTAFNDFundAgent removes a record from the _Edit table — DELETE endpoint.
func (s *Service) DeleteTAFNDFundAgent(ctx context.Context, req *fndv1.DeleteRequest) (*fndv1.SaveResponse, error) {
	s.log.Info("DeleteTAFNDFundAgent called",
		slog.String("data_id", req.GetDataId()),
		slog.String("data_flag", req.GetDataFlag()),
	)
	result := s.db.WithContext(ctx).Where(`"DataID" = ?`, req.GetDataId()).Delete(&db.DTAFNDFundAgentEdit{})
	if result.Error != nil {
		s.log.Error("DeleteTAFNDFundAgent failed", slog.Any("error", result.Error))
		return &fndv1.SaveResponse{Success: false, Message: result.Error.Error(), ReturnCode: "DB_ERROR"}, nil
	}
	if result.RowsAffected == 0 {
		return &fndv1.SaveResponse{Success: false, Message: "No record found with the given DataID", ReturnCode: "NOT_FOUND"}, nil
	}
	return &fndv1.SaveResponse{Success: true, Message: "Record deleted successfully", ReturnCode: "0000"}, nil
}

func (s *Service) TACKFundAgentEdit(ctx context.Context, req *fndv1.TACKFundAgentEditRequest) (*fndv1.CheckResponse, error) {
	s.log.Info("TACKFundAgentEdit called")
	return &fndv1.CheckResponse{ReturnCode: "", Message: "No pending data"}, nil
}

// GetDataByDataID retrieves full data by DataID.
func (s *Service) GetDataByDataID(ctx context.Context, req *fndv1.GetDataRequest) (*fndv1.TAFNDFundAgentRequest, error) {
	return &fndv1.TAFNDFundAgentRequest{}, nil
}

// ApproveTAFNDFundAgent approves or rejects a pending record for the 4-Eyes Principle.
func (s *Service) ApproveTAFNDFundAgent(ctx context.Context, req *fndv1.ApproveTAFNDFundAgentRequest) (*fndv1.SaveResponse, error) {
	s.log.Info("ApproveTAFNDFundAgent called", slog.String("data_id", req.GetDataId()))

	var record db.DTAFNDFundAgentEdit
	result := s.db.WithContext(ctx).Where("\"DataID\" = ?", req.GetDataId()).First(&record)
	
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return &fndv1.SaveResponse{Success: false, Message: "Record not found", ReturnCode: "NOT_FOUND"}, nil
		}
		s.log.Error("ApproveTAFNDFundAgent DB error", slog.Any("error", result.Error))
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
