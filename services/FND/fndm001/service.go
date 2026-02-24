package fndm001

import (
	"fmt"
	"strings"
	"time"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
	"go-transfer-agent/services/fnd/shared"
)

// ═══════════════════════════════════════════════════════════════════
// In-memory store (stub DB — replace with GORM repository)
// ═══════════════════════════════════════════════════════════════════

// StoredFund holds a full SaveFundInfoRequest payload in memory.
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

// Service holds the FNDM001 business logic.
type Service struct {
	Funds map[string]*StoredFund // key = "SysCoID|PrtFundCode"
}

// NewService creates a new FNDM001 service.
func NewService() *Service {
	return &Service{
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
func (s *Service) QueryFundInfoByDataID(sysCoID, prtFundCode string) (*fndv1.QueryFundInfoByDataIDResponse, error) {
	key := fundKey(sysCoID, prtFundCode)
	f, ok := s.Funds[key]
	if !ok {
		return nil, fmt.Errorf("fund not found: %s/%s", sysCoID, prtFundCode)
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

// SaveFundInfo — APIFNDM001Post/Put: persists all sub-lists.
// Field names in SaveFundInfoRequest match TAFNDFundInfoRequest.
func (s *Service) SaveFundInfo(req *fndv1.SaveFundInfoRequest) *fndv1.SaveFundInfoResponse {
	master := req.GetMaster()
	if master == nil || master.SysCoId == "" || master.PrtFundCode == "" {
		return &fndv1.SaveFundInfoResponse{
			Success: false, Message: "SysCoID and PrtFundCode are required", ReturnCode: "VALIDATION_ERROR",
		}
	}

	key := fundKey(master.SysCoId, master.PrtFundCode)
	// Proto field names → generated getters:
	// fund_info → GetFundInfo(), cust_contacts → GetCustContacts(), special → GetSpecial()
	// cry_groups → GetCryGroups(), fund_details → GetFundDetails(), fund_accounts → GetFundAccounts()
	// fund_rel_groups → GetFundRelGroups(), fund_dis_fees → GetFundDisFees(), fund_tx_crys → GetFundTxCrys()
	// short_infs → GetShortInfs(), short_dtls → GetShortDtls(), anti_dils → GetAntiDils()
	// mgt_fees → GetMgtFees(), mgt_fee_dtls → GetMgtFeeDtls()
	fundInfoSlice := []*fndv1.DTAFNDFundInfo{}
	if req.GetFundInfo() != nil {
		fundInfoSlice = append(fundInfoSlice, req.GetFundInfo())
	}
	specialSlice := []*fndv1.DTAFNDFundSpecial{}
	if req.GetSpecial() != nil {
		specialSlice = append(specialSlice, req.GetSpecial())
	}
	s.Funds[key] = &StoredFund{
		Master:              master,
		DTAFundInfoList:     fundInfoSlice,
		DTACustContactList:  req.GetCustContacts(),
		DTASpecialList:      specialSlice,
		DTACryGroupList:     req.GetCryGroups(),
		DTAFundDetailList:   req.GetFundDetails(),
		DTAFundAccountList:  req.GetFundAccounts(),
		DTAFundRelGroupList: req.GetFundRelGroups(),
		DTAFundDisFeeList:   req.GetFundDisFees(),
		DTAFundTxCryList:    req.GetFundTxCrys(),
		DTAFundShortInfList: req.GetShortInfs(),
		DTAShortDtlList:     req.GetShortDtls(),
		DTAFundAntiDilList:  req.GetAntiDils(),
		DTAFundMGTFeeList:   req.GetMgtFees(),
		DTAMGTFeeDtlList:    req.GetMgtFeeDtls(),
		SavedAt:             time.Now(),
	}

	return &fndv1.SaveFundInfoResponse{
		Success: true, Message: "Fund info saved successfully", ReturnCode: "OK",
	}
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
			{"TWD", "Y"}, {"USD", "Y"}, {"CNY", "N"}, {"EUR", "Y"}, {"JPY", "N"}, {"AUD", "Y"}, {"ZAR", "N"},
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
		for _, d := range []string{"TWD", "USD", "EUR", "JPY", "AUD", "CNY", "ZAR", "BRL"} {
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
		case "B", "P", "C", "N":
			if fd.SubsFeeType == feeType {
				return true
			}
		case "BN":
			if fd.SubsFeeType == "B" || fd.SubsFeeType == "N" {
				return true
			}
		case "PC":
			if fd.SubsFeeType == "P" || fd.SubsFeeType == "C" {
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
