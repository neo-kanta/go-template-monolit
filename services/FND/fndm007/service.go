package fndm007

import (
	"context"
	"log/slog"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
	"go-transfer-agent/services/fnd/fndm007/db"
	"go-transfer-agent/services/fnd/shared"

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
