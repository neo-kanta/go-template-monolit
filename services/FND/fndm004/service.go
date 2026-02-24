package fndm004

import (
	"context"
	"log/slog"
	"strings"
	"time"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
	"go-transfer-agent/services/fnd/fndm004/db"
	"go-transfer-agent/services/fnd/shared"

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
			AgentCodeNm:       "Agent Name", // Placeholder for lookup
			IsMAgentOpType:    d.IsMAgentOPType,
			AgentOpType:       d.AgentOPType,
			IsMRcvTxnType:     d.IsMRcvTxnType,
			RcvTxnType:        d.RcvTxnType,
			IsMSubsFeePct:     d.IsMSubsFeePct,
			SubsFeePctAg:      d.SubsFeePctAG,
			SubsFeePctFh:      d.SubsFeePctFH,
			IsMSubsFeePctType: d.IsMSubsFeePctType,
			SubsFeePctType:    d.SubsFeePctType,
			IsMFundCrySet:     d.IsMFundCrySet,
			FundCrySet:        splitString(d.FundCrySet),
			FundCrySetNm:      splitString(d.FundCrySet), // Same or mapped lookup
			NoSaleShareSet:    splitString(d.NoSaleShareSet),
			NoSaleShareSetNm:  splitString(d.NoSaleShareSet), // Same or mapped lookup
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
