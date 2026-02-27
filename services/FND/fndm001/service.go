package fndm001

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
	models "go-transfer-agent/common/platform/model"
	fnddb "go-transfer-agent/services/fnd/fndm001/db"
	"go-transfer-agent/services/fnd/shared"

	"gorm.io/gorm"
)

// ═══════════════════════════════════════════════════════════════════
// Service
// ═══════════════════════════════════════════════════════════════════

// StoredFund holds a full SaveFundInfoRequest payload in memory (used by Query methods).
type StoredFund struct {
	Master              *fndv1.TAFNDFundInfo
	DTAFundInfoList     []*fndv1.DTAFNDFundInfo
	DTACustContactList  []*fndv1.DTAFNDFundCustContact
	DTASpecialList      []*fndv1.DTAFNDFundSpecial
	DTACryGroupList     []*fndv1.DTAFNDFundCryGroup
	DTAFundDetailList   []*fndv1.DTAFNDFundDetail
	DTAFundAccountList  []*fndv1.DTAFNDFundAccount
	DTAFundRelGroupList []*fndv1.DTAFNDFundRelGroup
	DTAFundDisFeeList   []*fndv1.DTAFNDFundDisclosureFee
	DTAFundTxCryList    []*fndv1.DTAFNDFundTxCry
	DTAFundShortInfList []*fndv1.DTAFNDFundShortInf
	DTAShortDtlList     []*fndv1.DTAFNDFundShortDtl
	DTAFundAntiDilList  []*fndv1.DTAFNDFundAntiDilution
	DTAFundMGTFeeList   []*fndv1.DTAFNDFundMGTFeeInf
	DTAMGTFeeDtlList    []*fndv1.DTAFNDFundMGTFeeDtl
	SavedAt             time.Time
}

// Service holds the FNDM001 business logic with GORM database access.
type Service struct {
	db    *gorm.DB
	log   *slog.Logger
	Funds map[string]*StoredFund // key = "SysCoID|PrtFundCode" — kept for Query backward compat
}

// NewService creates a new FNDM001 service with GORM database connection.
func NewService(database *gorm.DB, logger *slog.Logger) *Service {
	return &Service{
		db:    database,
		log:   logger.With(slog.String("module", "fndm001-service")),
		Funds: make(map[string]*StoredFund),
	}
}

func fundKey(sysCoID, prtFundCode string) string {
	return sysCoID + "|" + prtFundCode
}

// ═══════════════════════════════════════════════════════════════════
// APIFNDM001: Main Operations
// ═══════════════════════════════════════════════════════════════════

// QueryFundInfo — APIFNDM001Get: 一般查詢, returns TAFNDFundInfoDTO list.
func (s *Service) QueryFundInfo(sysCoID string) *fndv1.QueryFundInfoResponse {
	var result []*fndv1.TAFNDFundInfoDTO
	for _, f := range s.Funds {
		if f.Master.SysCoId == sysCoID {
			dto := &fndv1.TAFNDFundInfoDTO{
				SysCoId:         f.Master.SysCoId,
				IsDeletable:     true,
				PrtFundCode:     f.Master.PrtFundCode,
				UniCode:         f.Master.UniCode,
				FundInShName:    f.Master.FundInShName,
				FundMName:       f.Master.FundMName,
				FundShMName:     f.Master.FundShMName,
				FundSName:       f.Master.FundSName,
				FundShSName:     f.Master.FundShSName,
				FundRiskLevelNm: shared.GetCodeName("000003", f.Master.FundRiskLevel),
				IssueBaseCry:    f.Master.IssueBaseCry,
				IssueBaseCryNm:  shared.GetCryName(f.Master.IssueBaseCry),
				FundSetupDate:   f.Master.FundSetupDate,
				FundStatusNm:    shared.GetCodeName("000007", f.Master.FundStatus),
			}
			result = append(result, dto)
		}
	}
	return &fndv1.QueryFundInfoResponse{ResultList: result}
}

// QueryFundInfoByDataID — APIFNDM001GetMaintain: full detail with name lookups.
func (s *Service) GetDataByDataID(ctx context.Context, req *fndv1.GetDataRequest) (*fndv1.QueryFundInfoByDataIDResponse, error) {
	var f *StoredFund
	for _, fund := range s.Funds {
		// Mock: just take the first fund, or ideally match on something mapped to DataID
		f = fund
		break
	}
	if f == nil {
		return nil, fmt.Errorf("fund not found for data_id: %s", req.GetDataId())
	}

	// ─── Enrich master (TAFNDFundInfo) ───
	master := cloneMaster(f.Master)
	master.SysCoIdNm = shared.GetSysCoName(master.SysCoId)
	master.IsDeletable = true
	master.IssueBaseCryNm = shared.GetCryName(master.IssueBaseCry)
	master.FundRiskLevelNm = shared.GetCodeName("000003", master.FundRiskLevel)
	master.FundStatusNm = shared.GetCodeName("000007", master.FundStatus)

	// ─── Enrich DTAFNDFundInfo ───
	var fundInfo *fndv1.DTAFNDFundInfo
	if len(f.DTAFundInfoList) > 0 {
		fundInfo = cloneFundInfo(f.DTAFundInfoList[0])
		if fundInfo != nil {
			bank := shared.GetBankInfo(fundInfo.CustBankHeadId)
			fundInfo.CustBankHeadIdNm = bank.Name
			fundInfo.SwiftCode = bank.SWIFTCode
			fundInfo.BankSName = bank.SName
			fundInfo.Addr = bank.Addr
			fundInfo.SitcaFundTypeNm = shared.GetCodeName("001", fundInfo.SitcaFundType)
			fundInfo.SubsFeeGroupNm = shared.GetCodeName("021", fundInfo.SubsFeeGroup)
			fundInfo.CdscFeeGroupNm = shared.GetCodeName("021", fundInfo.CdscFeeGroup)
			fundInfo.DivFundGroupNm = shared.GetCodeName("030", fundInfo.DivFundGroup)
		}
	}

	// ─── Enrich DTAFNDFundSpecial ───
	var special *fndv1.DTAFNDFundSpecial
	if len(f.DTASpecialList) > 0 {
		special = cloneSpecial(f.DTASpecialList[0])
		if special != nil && special.TmfReturnBc != "" {
			special.TmfReturnBcNm = shared.GetCryName(special.TmfReturnBc)
		}
	}

	// ─── Enrich DTAFNDFundCryGroup ───
	var cryGroups []*fndv1.DTAFNDFundCryGroup
	for _, cg := range f.DTACryGroupList {
		c := cloneCryGroup(cg)
		c.IssueCryNm = shared.GetCryName(c.IssueCry)
		if c.OthTxCry != "" {
			c.OthTxCryNm = shared.GetCryName(c.OthTxCry)
		}
		if c.FundMBankBrh != "" {
			brhBank := shared.GetBankInfo(c.FundMBankBrh)
			c.FundMBankBrhNm = brhBank.Name
			c.SwiftCode = brhBank.SWIFTCode
		}
		cryGroups = append(cryGroups, c)
	}

	// ─── Enrich DTAFNDFundDetail ───
	var fundDetails []*fndv1.DTAFNDFundDetail
	for _, fd := range f.DTAFundDetailList {
		d := cloneFundDetail(fd)
		d.FundCryNm = shared.GetCryName(d.FundCry)
		d.FhFundShareNm = shared.GetCodeName("031", d.FhFundShare)
		if d.MatureFund != "" {
			d.MatureFundNm = shared.GetDFundName(d.MatureFund)
		}
		fundDetails = append(fundDetails, d)
	}

	// ─── Enrich DTAFNDFundAccount ───
	var fundAccounts []*fndv1.DTAFNDFundAccount
	for _, fa := range f.DTAFundAccountList {
		a := cloneFundAccount(fa)
		if a.CusBankBrh != "" {
			brhBank := shared.GetBankInfo(a.CusBankBrh)
			a.CusBankBrhNm = brhBank.Name
			a.SwiftCode = brhBank.SWIFTCode
		}
		fundAccounts = append(fundAccounts, a)
	}

	// ─── Enrich DTAFNDFundRelGroup ───
	var fundRelGroups []*fndv1.DTAFNDFundRelGroup
	for _, rg := range f.DTAFundRelGroupList {
		g := cloneFundRelGroup(rg)
		g.FundGroupNoNm = shared.GetCodeTypeName(g.FundGroupNo)
		g.GroupTypeNoNm = shared.GetCodeName(g.FundGroupNo, g.GroupTypeNo)
		fundRelGroups = append(fundRelGroups, g)
	}

	// ─── Enrich DTAFNDFundMGTFeeInf ───
	var mgtFees []*fndv1.DTAFNDFundMGTFeeInf
	for _, mf := range f.DTAFundMGTFeeList {
		m := cloneMGTFee(mf)
		m.FhFundShareNm = shared.GetCodeName("031", m.FhFundShare)
		mgtFees = append(mgtFees, m)
	}

	// ─── Split DTAFNDFundAntiDilution by type (spec §4.1/4.2) ───
	var antiDilMechs []*fndv1.DTAFNDFundAntiDilMech
	var antiDilAdjs []*fndv1.DTAFNDFundAntiDilAdj
	for _, ad := range f.DTAFundAntiDilList {
		if ad.AntiDilSetType == "1" {
			antiDilMechs = append(antiDilMechs, &fndv1.DTAFNDFundAntiDilMech{
				SysCoId: ad.SysCoId, PrtFundCode: ad.PrtFundCode,
				EffDate: ad.EffDate, AntiDilTrigger: ad.AntiDilTrigger,
				AntiDilFeeRate: ad.AntiDilFeeRate, Remark: ad.Remark,
			})
		} else {
			antiDilAdjs = append(antiDilAdjs, &fndv1.DTAFNDFundAntiDilAdj{
				SysCoId: ad.SysCoId, PrtFundCode: ad.PrtFundCode,
				AntiDilSetType: ad.AntiDilSetType, EffDate: ad.EffDate,
				TermDate: ad.TermDate, AntiDilTrigger: ad.AntiDilTrigger,
				AntiDilFeeRate: ad.AntiDilFeeRate, AdjRsn: ad.AdjRsn, Remark: ad.Remark,
			})
		}
	}

	return &fndv1.QueryFundInfoByDataIDResponse{
		Master:        master,
		FundInfo:      fundInfo,
		CustContacts:  f.DTACustContactList,
		Special:       special,
		CryGroups:     cryGroups,
		FundDetails:   fundDetails,
		FundAccounts:  fundAccounts,
		FundRelGroups: fundRelGroups,
		FundDisFees:   f.DTAFundDisFeeList,
		FundTxCrys:    f.DTAFundTxCryList,
		ShortInfs:     f.DTAFundShortInfList,
		ShortDtls:     f.DTAShortDtlList,
		AntiDils:      f.DTAFundAntiDilList,
		AntiDilMechs:  antiDilMechs,
		AntiDilAdjs:   antiDilAdjs,
		MgtFees:       mgtFees,
		MgtFeeDtls:    f.DTAMGTFeeDtlList,
	}, nil
}

// SaveFundInfo — APIFNDM001Post: persists master + all child tables to _Edit tables in a transaction.
// 4-Eyes Principle: Automatically sets MakerID from JWT context and Status to PENDING_APPROVAL.
func (s *Service) SaveFundInfo(ctx context.Context, req *fndv1.SaveFundInfoRequest) *fndv1.SaveFundInfoResponse {
	master := req.GetMaster()
	if master == nil || master.SysCoId == "" || master.PrtFundCode == "" {
		return &fndv1.SaveFundInfoResponse{
			Success: false, Message: "SysCoID and PrtFundCode are required", ReturnCode: "VALIDATION_ERROR",
		}
	}

	// Extract authenticated user from JWT context (Maker)
	makerID := shared.GetUsernameFromCtx(ctx)
	if makerID == "" {
		makerID = "system" // Fallback for tests or unauthenticated contexts
	}

	s.log.Info("SaveFundInfo called",
		slog.String("sys_co_id", master.GetSysCoId()),
		slog.String("prt_fund_code", master.GetPrtFundCode()),
		slog.String("maker_id", makerID),
	)

	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 1. Master table — auto-set 4-Eyes fields
		masterEdit := mapMasterEdit(master)
		masterEdit.MakerID = makerID
		masterEdit.Status = models.StatusPendingApproval
		if err := tx.Create(&masterEdit).Error; err != nil {
			return fmt.Errorf("master: %w", err)
		}

		// 2. DTAFNDFundInfo (detail fund info)
		if info := req.GetFundInfo(); info != nil {
			records := mapDTAFNDFundInfoEditList([]*fndv1.DTAFNDFundInfo{info})
			if len(records) > 0 {
				if err := tx.Create(&records).Error; err != nil {
					return fmt.Errorf("fund_info: %w", err)
				}
			}
		}

		// 3. CustContact
		if items := req.GetCustContacts(); len(items) > 0 {
			records := mapDTAFNDFundCustContactEditList(items)
			if err := tx.Create(&records).Error; err != nil {
				return fmt.Errorf("cust_contacts: %w", err)
			}
		}

		// 4. Special
		if sp := req.GetSpecial(); sp != nil {
			records := mapDTAFNDFundSpecialEditList([]*fndv1.DTAFNDFundSpecial{sp})
			if len(records) > 0 {
				if err := tx.Create(&records).Error; err != nil {
					return fmt.Errorf("special: %w", err)
				}
			}
		}

		// 5. CryGroup
		if items := req.GetCryGroups(); len(items) > 0 {
			records := mapDTAFNDFundCryGroupEditList(items)
			if err := tx.Create(&records).Error; err != nil {
				return fmt.Errorf("cry_groups: %w", err)
			}
		}

		// 6. FundDetail
		if items := req.GetFundDetails(); len(items) > 0 {
			records := mapDTAFNDFundDetailEditList(items)
			if err := tx.Create(&records).Error; err != nil {
				return fmt.Errorf("fund_details: %w", err)
			}
		}

		// 7. FundAccount
		if items := req.GetFundAccounts(); len(items) > 0 {
			records := mapDTAFNDFundAccountEditList(items)
			if err := tx.Create(&records).Error; err != nil {
				return fmt.Errorf("fund_accounts: %w", err)
			}
		}

		// 8. FundRelGroup
		if items := req.GetFundRelGroups(); len(items) > 0 {
			records := mapDTAFNDFundRelGroupEditList(items)
			if err := tx.Create(&records).Error; err != nil {
				return fmt.Errorf("fund_rel_groups: %w", err)
			}
		}

		// 9. FundDisclosureFee
		if items := req.GetFundDisFees(); len(items) > 0 {
			records := mapDTAFNDFundDisclosureFeeEditList(items)
			if err := tx.Create(&records).Error; err != nil {
				return fmt.Errorf("fund_dis_fees: %w", err)
			}
		}

		// 10. FundTxCry
		if items := req.GetFundTxCrys(); len(items) > 0 {
			records := mapDTAFNDFundTxCryEditList(items)
			if err := tx.Create(&records).Error; err != nil {
				return fmt.Errorf("fund_tx_crys: %w", err)
			}
		}

		// 11. FundShortInf
		if items := req.GetShortInfs(); len(items) > 0 {
			records := mapDTAFNDFundShortInfEditList(items)
			if err := tx.Create(&records).Error; err != nil {
				return fmt.Errorf("short_infs: %w", err)
			}
		}

		// 12. AntiDilution
		if items := req.GetAntiDils(); len(items) > 0 {
			records := mapDTAFNDFundAntiDilutionEditList(items)
			if err := tx.Create(&records).Error; err != nil {
				return fmt.Errorf("anti_dils: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		s.log.Error("SaveFundInfo failed", slog.Any("error", err))
		return &fndv1.SaveFundInfoResponse{
			Success: false, Message: err.Error(), ReturnCode: "DB_ERROR",
		}
	}

	return &fndv1.SaveFundInfoResponse{
		Success: true, Message: "Record created successfully", ReturnCode: "0000",
	}
}

// UpdateFundInfo — APIFNDM001Put: updates an existing record in _Edit tables.
func (s *Service) UpdateFundInfo(req *fndv1.SaveFundInfoRequest) *fndv1.SaveFundInfoResponse {
	master := req.GetMaster()
	if master == nil || master.SysCoId == "" || master.PrtFundCode == "" {
		return &fndv1.SaveFundInfoResponse{
			Success: false, Message: "SysCoID and PrtFundCode are required", ReturnCode: "VALIDATION_ERROR",
		}
	}

	s.log.Info("UpdateFundInfo called",
		slog.String("sys_co_id", master.GetSysCoId()),
		slog.String("prt_fund_code", master.GetPrtFundCode()),
	)

	masterEdit := mapMasterEdit(master)
	result := s.db.Model(&fnddb.TAFNDFundInfoEdit{}).
		Where(`"SysCoID" = ? AND "PrtFundCode" = ?`, master.GetSysCoId(), master.GetPrtFundCode()).Updates(&masterEdit)
	if result.Error != nil {
		s.log.Error("UpdateFundInfo failed", slog.Any("error", result.Error))
		return &fndv1.SaveFundInfoResponse{
			Success: false, Message: result.Error.Error(), ReturnCode: "DB_ERROR",
		}
	}
	if result.RowsAffected == 0 {
		return &fndv1.SaveFundInfoResponse{
			Success: false, Message: "No record found to update", ReturnCode: "NOT_FOUND",
		}
	}

	return &fndv1.SaveFundInfoResponse{
		Success: true, Message: "Record updated successfully", ReturnCode: "0000",
	}
}

// ApproveFundInfo — APIFNDM001Approve: 4-Eyes Principle Approval
// Inspired by legacy C# FlowAPI ApproveFlowCommandHandler:
//   - Validates current status is PENDING_APPROVAL (cannot approve already approved/rejected records)
//   - Enforces MakerID != CheckerID (4-eyes constraint)
//   - Extracts CheckerID from JWT context if not provided explicitly
func (s *Service) ApproveFundInfo(ctx context.Context, req *fndv1.ApproveFundInfoRequest) *fndv1.SaveFundInfoResponse {
	if req.SysCoId == "" || req.PrtFundCode == "" {
		return &fndv1.SaveFundInfoResponse{
			Success: false, Message: "SysCoID and PrtFundCode are required", ReturnCode: "VALIDATION_ERROR",
		}
	}

	// Extract CheckerID from JWT context if not explicitly provided
	checkerID := req.CheckerId
	if checkerID == "" {
		checkerID = shared.GetUsernameFromCtx(ctx)
	}
	if checkerID == "" {
		return &fndv1.SaveFundInfoResponse{
			Success: false, Message: "CheckerID is required (provide in request or authenticate)", ReturnCode: "VALIDATION_ERROR",
		}
	}

	s.log.Info("ApproveFundInfo called",
		slog.String("sys_co_id", req.SysCoId),
		slog.String("prt_fund_code", req.PrtFundCode),
		slog.String("checker_id", checkerID),
	)

	var masterEdit fnddb.TAFNDFundInfoEdit
	result := s.db.Where(`"SysCoID" = ? AND "PrtFundCode" = ?`, req.SysCoId, req.PrtFundCode).First(&masterEdit)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return &fndv1.SaveFundInfoResponse{Success: false, Message: "Record not found", ReturnCode: "NOT_FOUND"}
		}
		return &fndv1.SaveFundInfoResponse{Success: false, Message: result.Error.Error(), ReturnCode: "DB_ERROR"}
	}

	// Validate current status — only PENDING_APPROVAL records can be approved/rejected
	if masterEdit.Status != models.StatusPendingApproval {
		return &fndv1.SaveFundInfoResponse{
			Success:    false,
			Message:    fmt.Sprintf("Cannot approve/reject: current status is %s, expected PENDING_APPROVAL", masterEdit.Status),
			ReturnCode: "INVALID_STATUS",
		}
	}

	// 4-Eyes Principle: Checker cannot be the same as Maker
	if masterEdit.MakerID == checkerID {
		return &fndv1.SaveFundInfoResponse{
			Success:    false,
			Message:    "4-Eyes Violation: Checker cannot be the same as Maker (" + masterEdit.MakerID + ")",
			ReturnCode: "FOUR_EYES_VIOLATION",
		}
	}

	newStatus := models.StatusRejected
	if req.IsApproved {
		newStatus = models.StatusApproved
	}

	updates := map[string]interface{}{
		"CheckerID": checkerID,
		"Status":    newStatus,
		"Remark":    req.Remark,
	}

	if err := s.db.Model(&fnddb.TAFNDFundInfoEdit{}).
		Where(`"SysCoID" = ? AND "PrtFundCode" = ?`, req.SysCoId, req.PrtFundCode).
		Updates(updates).Error; err != nil {
		return &fndv1.SaveFundInfoResponse{
			Success: false, Message: err.Error(), ReturnCode: "DB_ERROR",
		}
	}

	return &fndv1.SaveFundInfoResponse{
		Success: true, Message: fmt.Sprintf("Record %s successfully", newStatus), ReturnCode: "0000",
	}
}

// DeleteTAFNDFundInfo deletes data by DataID — DELETE endpoint.
func (s *Service) DeleteTAFNDFundInfo(ctx context.Context, req *fndv1.DeleteRequest) (*fndv1.SaveResponse, error) {
	s.log.Info("DeleteTAFNDFundInfo called",
		slog.String("data_id", req.GetDataId()),
		slog.String("data_flag", req.GetDataFlag()),
	)

	result := s.db.WithContext(ctx).
		Where(`"DataID" = ?`, req.GetDataId()).
		Delete(&fnddb.TAFNDFundInfoEdit{})
	if result.Error != nil {
		s.log.Error("DeleteTAFNDFundInfo failed", slog.Any("error", result.Error))
		return &fndv1.SaveResponse{
			Success: false, Message: result.Error.Error(), ReturnCode: "DB_ERROR",
		}, nil
	}
	if result.RowsAffected == 0 {
		return &fndv1.SaveResponse{
			Success: false, Message: "No record found with the given DataID", ReturnCode: "NOT_FOUND",
		}, nil
	}

	return &fndv1.SaveResponse{
		Success: true, Message: "Record deleted successfully", ReturnCode: "0000",
	}, nil
}

// ═══════════════════════════════════════════════════════════════════
// KFNDM: Check Functions
// ═══════════════════════════════════════════════════════════════════

// TACKCSDPrtFundSetlDate — KFNDM00101: validates CSD settlement date.
func (s *Service) TACKCSDPrtFundSetlDate(sysCoID, csdDateStr string, fundCode *string) *fndv1.CheckResponse {
	setlDate, err := time.Parse(time.RFC3339, csdDateStr)
	if err != nil {
		setlDate, err = time.Parse("2006-01-02", csdDateStr)
		if err != nil {
			return &fndv1.CheckResponse{ReturnCode: "INVALID_DATE", Message: "Invalid date format. Use RFC3339 or YYYY-MM-DD"}
		}
	}

	// Stub: SysCode=000021 → 3 business days offset
	effDate := time.Now().AddDate(0, 0, 3)
	if setlDate.Before(effDate) {
		code := "FNDM001009"
		if fundCode != nil && *fundCode != "" {
			code = "FNDM001014"
		}
		return &fndv1.CheckResponse{
			ReturnCode: code,
			Message:    fmt.Sprintf("Settlement date %s is before effective date %s", setlDate.Format("2006-01-02"), effDate.Format("2006-01-02")),
		}
	}
	return &fndv1.CheckResponse{ReturnCode: "OK", Message: "CSD settlement date is valid"}
}

// TACKFundCodeExist — KFNDM00103: checks if fund code exists in same company.
func (s *Service) TACKFundCodeExist(sysCoID, fundCode string) *fndv1.CheckResponse {
	for _, f := range s.Funds {
		if f.Master.SysCoId == sysCoID {
			for _, fd := range f.DTAFundDetailList {
				if fd.FundCode == fundCode {
					return &fndv1.CheckResponse{
						ReturnCode: "TA00030",
						Message:    fmt.Sprintf("Fund code '%s' already exists in company '%s'", fundCode, sysCoID),
					}
				}
			}
		}
	}
	return &fndv1.CheckResponse{ReturnCode: "OK", Message: fmt.Sprintf("Fund code '%s' is available", fundCode)}
}

// ═══════════════════════════════════════════════════════════════════
// GFNDM: Get/Search Functions
// ═══════════════════════════════════════════════════════════════════

// TAGetDPrtFundIssueCry — GFNDM00101.
func (s *Service) TAGetDPrtFundIssueCry(sysCoID string, isAll bool, prtFundCode, feeChargeType, filterItem, queryType string) *fndv1.TAGetDPrtFundIssueCryResponse {
	var result []*fndv1.IssueCurrencyDTO

	if isAll {
		result = append(result, &fndv1.IssueCurrencyDTO{CryId: "*", CryName: "全部(ALL)", IsRsp: ""})
	}

	seen := map[string]bool{}
	for _, f := range s.Funds {
		if sysCoID != "" && f.Master.SysCoId != sysCoID {
			continue
		}
		if prtFundCode != "" && f.Master.PrtFundCode != prtFundCode {
			continue
		}
		for _, cg := range f.DTACryGroupList {
			if seen[cg.IssueCry] {
				continue
			}
			if feeChargeType != "" && !matchFeeChargeType(f.DTAFundDetailList, cg.IssueCry, feeChargeType) {
				continue
			}
			if filterItem != "" && cg.IssueCry == filterItem {
				continue
			}
			seen[cg.IssueCry] = true
			entry := &fndv1.IssueCurrencyDTO{CryId: cg.IssueCry, CryName: shared.GetCryName(cg.IssueCry)}
			if queryType == "ALL" {
				entry.IsRsp = cg.IsRsp
			}
			result = append(result, entry)
		}
	}

	if len(seen) == 0 {
		for _, d := range []struct {
			id  string
			rsp string
		}{
			{"THB", "Y"}, {"USD", "Y"}, {"CNY", "N"}, {"EUR", "Y"}, {"JPY", "N"}, {"AUD", "Y"}, {"ZAR", "N"},
		} {
			if filterItem != "" && d.id == filterItem {
				continue
			}
			e := &fndv1.IssueCurrencyDTO{CryId: d.id, CryName: shared.GetCryName(d.id)}
			if queryType == "ALL" {
				e.IsRsp = d.rsp
			}
			result = append(result, e)
		}
	}

	return &fndv1.TAGetDPrtFundIssueCryResponse{Master: result}
}

// TAGetDTxCry — GFNDM00104.
func (s *Service) TAGetDTxCry(sysCoID string, isAll bool, prtFundCode, cryID, cryDataSrc string, isExcludeCCY bool, queryType string) *fndv1.TAGetDTxCryResponse {
	var result []*fndv1.TransactionCurrencyDTO

	if isAll {
		result = append(result, &fndv1.TransactionCurrencyDTO{CryId: "*", CryName: "全部(ALL)", DecLen: 0, IsMcy: ""})
	}

	cryFilter := map[string]bool{}
	if cryID != "" {
		for _, c := range strings.Split(cryID, ",") {
			cryFilter[strings.TrimSpace(c)] = true
		}
	}

	seen := map[string]bool{}
	for _, f := range s.Funds {
		if sysCoID != "" && f.Master.SysCoId != sysCoID {
			continue
		}
		if prtFundCode != "" && f.Master.PrtFundCode != prtFundCode {
			continue
		}
		for _, tc := range f.DTAFundTxCryList {
			if cryDataSrc != "" && tc.CryDataSrc != cryDataSrc {
				continue
			}
			if len(cryFilter) > 0 && !cryFilter[tc.CryId] {
				continue
			}
			if isExcludeCCY && tc.CryId == shared.CenterCurrency {
				continue
			}
			if seen[tc.CryId] {
				continue
			}
			seen[tc.CryId] = true

			entry := &fndv1.TransactionCurrencyDTO{CryId: tc.CryId, CryName: shared.GetCryName(tc.CryId), DecLen: shared.GetCryDecLen(tc.CryId)}
			if queryType == "ALL" {
				if shared.IsMiscCurrency(tc.CryId) {
					entry.IsMcy = "Y"
				} else {
					entry.IsMcy = "N"
				}
			}
			result = append(result, entry)
		}
	}

	if len(seen) == 0 {
		for _, d := range []string{"THB", "USD", "EUR", "JPY", "AUD", "CNY", "ZAR", "BRL"} {
			if isExcludeCCY && d == shared.CenterCurrency {
				continue
			}
			if len(cryFilter) > 0 && !cryFilter[d] {
				continue
			}
			e := &fndv1.TransactionCurrencyDTO{CryId: d, CryName: shared.GetCryName(d), DecLen: shared.GetCryDecLen(d)}
			if queryType == "ALL" {
				if shared.IsMiscCurrency(d) {
					e.IsMcy = "Y"
				} else {
					e.IsMcy = "N"
				}
			}
			result = append(result, e)
		}
	}
	return &fndv1.TAGetDTxCryResponse{Master: result}
}

// ═══════════════════════════════════════════════════════════════════
// Helpers
// ═══════════════════════════════════════════════════════════════════

func matchFeeChargeType(details []*fndv1.DTAFNDFundDetail, issueCry, feeType string) bool {
	for _, fd := range details {
		if fd.FundCry != issueCry {
			continue
		}
		switch feeType {
		case shared.SubsFeeType_Front, shared.SubsFeeType_Back, shared.SubsFeeType_Deferred, shared.SubsFeeType_None:
			// Single fee type: B, P, C, or N
			if fd.SubsFeeType == feeType {
				return true
			}
		case shared.FeeChargeType_FrontNil: // "BN" = Front + None
			if fd.SubsFeeType == shared.SubsFeeType_Front || fd.SubsFeeType == shared.SubsFeeType_None {
				return true
			}
		case shared.FeeChargeType_BackDef: // "PC" = Back + Deferred
			if fd.SubsFeeType == shared.SubsFeeType_Back || fd.SubsFeeType == shared.SubsFeeType_Deferred {
				return true
			}
		}
	}
	return len(details) == 0
}

func cloneMaster(m *fndv1.TAFNDFundInfo) *fndv1.TAFNDFundInfo {
	if m == nil {
		return nil
	}
	c := *m
	return &c
}
func cloneFundInfo(f *fndv1.DTAFNDFundInfo) *fndv1.DTAFNDFundInfo {
	if f == nil {
		return nil
	}
	c := *f
	return &c
}
func cloneSpecial(s *fndv1.DTAFNDFundSpecial) *fndv1.DTAFNDFundSpecial {
	if s == nil {
		return nil
	}
	c := *s
	return &c
}
func cloneCryGroup(cg *fndv1.DTAFNDFundCryGroup) *fndv1.DTAFNDFundCryGroup     { c := *cg; return &c }
func cloneFundDetail(fd *fndv1.DTAFNDFundDetail) *fndv1.DTAFNDFundDetail       { c := *fd; return &c }
func cloneFundAccount(fa *fndv1.DTAFNDFundAccount) *fndv1.DTAFNDFundAccount    { c := *fa; return &c }
func cloneFundRelGroup(rg *fndv1.DTAFNDFundRelGroup) *fndv1.DTAFNDFundRelGroup { c := *rg; return &c }
func cloneMGTFee(mf *fndv1.DTAFNDFundMGTFeeInf) *fndv1.DTAFNDFundMGTFeeInf     { c := *mf; return &c }
