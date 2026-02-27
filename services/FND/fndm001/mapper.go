package fndm001

import (
	"go-transfer-agent/common/platform/thaidate"
	"go-transfer-agent/services/fnd/fndm001/db"
	"time"

	"github.com/shopspring/decimal"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
)

// parseTime is a helper for parsing RFC3339 timestamps from protobuf strings.
func parseTime(s string) time.Time {
	t, _ := thaidate.ParseBE(s)
	return t
}

// ═══════════════════════════════════════════════════════════════════
// Proto → DB Mappers (for AUD operations)
// ═══════════════════════════════════════════════════════════════════

// mapMasterEdit maps the TAFNDFundInfo proto to its _Edit DB model.
func mapMasterEdit(m *fndv1.TAFNDFundInfo) db.TAFNDFundInfoEdit {
	var r db.TAFNDFundInfoEdit
	if m == nil {
		return r
	}
	r.SysCoID = m.GetSysCoId()
	r.PrtFundCode = m.GetPrtFundCode()
	r.UniCode = m.GetUniCode()
	r.FundInShName = m.GetFundInShName()
	r.FundMName = m.GetFundMName()
	r.FundShMName = m.GetFundShMName()
	r.FundSName = m.GetFundSName()
	r.FundShSName = m.GetFundShSName()
	r.FundRiskLevel = m.GetFundRiskLevel()
	r.IssueBaseCry = m.GetIssueBaseCry()
	r.FundSetupDate = parseTime(m.GetFundSetupDate())
	r.GIINNo = m.GetGiinNo()
	r.FundStatus = m.GetFundStatus()
	r.FundWarningMsg = m.GetFundWarningMsg()
	r.UnitDotLen = int16(m.GetUnitDotLen())
	r.FundSubsWay = m.GetFundSubsWay()
	r.FundRdmWay = m.GetFundRdmWay()
	r.RdmFNavDay = int16(m.GetRdmFNavDay())
	r.RdmFNavWay = m.GetRdmFNavWay()
	r.RdmFNavDRate = decimal.NewFromFloat(m.GetRdmFNavDRate()).RoundBank(4).RoundBank(4).RoundBank(4)
	r.RdmFNavDVal = decimal.NewFromFloat(m.GetRdmFNavDVal()).RoundBank(4).RoundBank(4).RoundBank(4)
	r.IsShortFee = m.GetIsShortFee()
	r.IsAntiDilFee = m.GetIsAntiDilFee()
	r.IsETF = m.GetIsEtf()
	r.ETFCode = m.GetEtfCode()
	r.IsDIM = m.GetIsDim()
	r.ShoreID = m.GetShoreId()
	r.OFDFHCode = m.GetOfdFhCode()
	return r
}

func mapDTAFNDFundInfoEditList(srcs []*fndv1.DTAFNDFundInfo) []db.DTAFNDFundInfoEdit {
	res := make([]db.DTAFNDFundInfoEdit, len(srcs))
	for i, s := range srcs {
		res[i].SysCoID = s.GetSysCoId()
		res[i].PrtFundCode = s.GetPrtFundCode()
		res[i].SITCAFundType = s.GetSitcaFundType()
		res[i].DFundType = s.GetDFundType()
		res[i].DInvArea = s.GetDInvArea()
		res[i].OfferType = s.GetOfferType()
		res[i].BCFundQuoUpAmt = decimal.NewFromFloat(s.GetBcFundQuoUpAmt()).RoundBank(2).RoundBank(2).RoundBank(2)
		res[i].BCFundQuoMinAmt = decimal.NewFromFloat(s.GetBcFundQuoMinAmt()).RoundBank(2).RoundBank(2).RoundBank(2)
		res[i].FundQuoLmtRate = decimal.NewFromFloat(s.GetFundQuoLmtRate()).RoundBank(4).RoundBank(4).RoundBank(4)
		res[i].TradeWarningAmt = decimal.NewFromFloat(s.GetTradeWarningAmt()).RoundBank(2).RoundBank(2).RoundBank(2)
		res[i].ExchUnitLmt = s.GetExchUnitLmt()
		res[i].FMNetQuoUpUnit = decimal.NewFromFloat(s.GetFmNetQuoUpUnit()).RoundBank(4).RoundBank(4).RoundBank(4)
		res[i].ExchNetQuoUpUnit = decimal.NewFromFloat(s.GetExchNetQuoUpUnit()).RoundBank(4).RoundBank(4).RoundBank(4)
		res[i].FundIPODate = parseTime(s.GetFundIpoDate())
		res[i].BegRdmDate = parseTime(s.GetBegRdmDate())
		res[i].FundIssueType = s.GetFundIssueType()
		res[i].FundIssueDate = parseTime(s.GetFundIssueDate())
		res[i].CSDOnlineDate = parseTime(s.GetCsdOnlineDate())
		res[i].CSDPrtFundSetlDate = parseTime(s.GetCsdPrtFundSetlDate())
		res[i].ApprovelDate = parseTime(s.GetApprovelDate())
		res[i].ClearDate = parseTime(s.GetClearDate())
		res[i].LastSubsDate = parseTime(s.GetLastSubsDate())
		res[i].LastRdmDate = parseTime(s.GetLastRdmDate())
		res[i].BIZDayDef = s.GetBizDayDef()
		res[i].SubsDateType = s.GetSubsDateType()
		res[i].RdmDateType = s.GetRdmDateType()
		res[i].RdmPayDateType = s.GetRdmPayDateType()
		res[i].SubsNavDay = int16(s.GetSubsNavDay())
		res[i].RdmNavDay = int16(s.GetRdmNavDay())
		res[i].RdmPayDay = int16(s.GetRdmPayDay())
		res[i].CustBankHeadID = s.GetCustBankHeadId()
		res[i].FundAccMName = s.GetFundAccMName()
		res[i].FundAccSName = s.GetFundAccSName()
		res[i].CheckDiffDay = int16(s.GetCheckDiffDay())
		res[i].CashDiffDay = int16(s.GetCashDiffDay())
		res[i].FASubsInDay = int16(s.GetFaSubsInDay())
		res[i].FARdmOutDay = int16(s.GetFaRdmOutDay())
		res[i].FASubsSetlDay = int16(s.GetFaSubsSetlDay())
		res[i].FARdmSetlDay = int16(s.GetFaRdmSetlDay())
		res[i].SubsFeeGroup = s.GetSubsFeeGroup()
		res[i].CDSCFeeGroup = s.GetCdscFeeGroup()
		res[i].DivFundGroup = s.GetDivFundGroup()
		res[i].CSDSubsDeadline = s.GetCsdSubsDeadline()
		res[i].CSDRdmDeadline = s.GetCsdRdmDeadline()
		res[i].BeginShortDate = parseTime(s.GetBeginShortDate())
		res[i].ShortCalcMemo = s.GetShortCalcMemo()
		res[i].ShortTxMemo = s.GetShortTxMemo()
		res[i].ShortMemo = s.GetShortMemo()
		res[i].RemitChargeCD = s.GetRemitChargeCd()
		res[i].IsIPO = s.GetIsIpo()
	}
	return res
}

func mapDTAFNDFundCustContactEditList(srcs []*fndv1.DTAFNDFundCustContact) []db.DTAFNDFundCustContactEdit {
	res := make([]db.DTAFNDFundCustContactEdit, len(srcs))
	for i, s := range srcs {
		res[i].SysCoID = s.GetSysCoId()
		res[i].PrtFundCode = s.GetPrtFundCode()
		res[i].DataSeq = int16(s.GetDataSeq())
		res[i].IsMainContact = s.GetIsMainContact()
		res[i].ContactMan = s.GetContactMan()
		res[i].ContactTel = s.GetContactTel()
		res[i].Fax1 = s.GetFax1()
		res[i].Fax2 = s.GetFax2()
		res[i].WorkItem = s.GetWorkItem()
	}
	return res
}

func mapDTAFNDFundSpecialEditList(srcs []*fndv1.DTAFNDFundSpecial) []db.DTAFNDFundSpecialEdit {
	res := make([]db.DTAFNDFundSpecialEdit, len(srcs))
	for i, s := range srcs {
		res[i].SysCoID = s.GetSysCoId()
		res[i].PrtFundCode = s.GetPrtFundCode()
		res[i].IsHighRisk = s.GetIsHighRisk()
		res[i].IsPGF = s.GetIsPgf()
		res[i].PGFStay = s.GetPgfStay()
		res[i].PGFMatureYear = int16(s.GetPgfMatureYear())
		res[i].PGFStopSubsDay = int16(s.GetPgfStopSubsDay())
		res[i].PGFExpiryDate = parseTime(s.GetPgfExpiryDate())
		res[i].IsPGFRdmFee = s.GetIsPgfRdmFee()
		res[i].IsTMF = s.GetIsTmf()
		res[i].TMFKind = s.GetTmfKind()
		res[i].TMFStopSubsDay = int16(s.GetTmfStopSubsDay())
		res[i].TMFTargetDate = parseTime(s.GetTmfTargetDate())
		res[i].TMFStay = int16(s.GetTmfStay())
		res[i].TMFExpiryDate = parseTime(s.GetTmfExpiryDate())
		res[i].TMFRetunrBC = s.GetTmfReturnBc()
		res[i].TMFReturnStarr = s.GetTmfReturnStarr()
		res[i].IsTMFRdmFee = s.GetIsTmfRdmFee()
		res[i].IsPFF = s.GetIsPff()
		res[i].PvtFundLimit = int16(s.GetPvtFundLimit())
		res[i].IsPerformanceFee = s.GetIsPerformanceFee()
		res[i].IsUmbrella = s.GetIsUmbrella()
		res[i].UmbrellaName = s.GetUmbrellaName()
	}
	return res
}

func mapDTAFNDFundCryGroupEditList(srcs []*fndv1.DTAFNDFundCryGroup) []db.DTAFNDFundCryGroupEdit {
	res := make([]db.DTAFNDFundCryGroupEdit, len(srcs))
	for i, s := range srcs {
		res[i].SysCoID = s.GetSysCoId()
		res[i].PrtFundCode = s.GetPrtFundCode()
		res[i].IssueCry = s.GetIssueCry()
		res[i].FundFaceAmt = decimal.NewFromFloat(s.GetFundFaceAmt()).RoundBank(2).RoundBank(2).RoundBank(2)
		res[i].NavDotLen = int16(s.GetNavDotLen())
		res[i].IsOthTxCry = s.GetIsOthTxCry()
		res[i].OthTxCry = s.GetOthTxCry()
		res[i].FundQuoAmtBC = decimal.NewFromFloat(s.GetFundQuoAmtBc()).RoundBank(2).RoundBank(2).RoundBank(2)
		res[i].UnitConvRate = decimal.NewFromFloat(s.GetUnitConvRate()).RoundBank(4).RoundBank(4).RoundBank(4)
		res[i].FundQuoUnit = decimal.NewFromFloat(s.GetFundQuoUnit()).RoundBank(4).RoundBank(4).RoundBank(4)
		res[i].IPOSubsFeeDisType = s.GetIpoSubsFeeDisType()
		res[i].IPOFixRate = decimal.NewFromFloat(s.GetIpoFixRate()).RoundBank(4).RoundBank(4).RoundBank(4)
		res[i].IPODiscRateOff = decimal.NewFromFloat(s.GetIpoDiscRateOff()).RoundBank(4).RoundBank(4).RoundBank(4)
		res[i].FundMBankBrh = s.GetFundMBankBrh()
		res[i].FundMAccount = s.GetFundMAccount()
		res[i].CRemitFeeID = s.GetCRemitFeeId()
		res[i].CRemitFeeFM = decimal.NewFromFloat(s.GetCRemitFeeFm()).RoundBank(2).RoundBank(2).RoundBank(2)
		res[i].CRemitFeeOtherID = s.GetCRemitFeeOtherId()
		res[i].CRemitFeeOtherFM = decimal.NewFromFloat(s.GetCRemitFeeOtherFm()).RoundBank(2).RoundBank(2).RoundBank(2)
		res[i].IsRSP = s.GetIsRsp()
		res[i].BeginRspDate = parseTime(s.GetBeginRspDate())
		res[i].RSPMinAmt = decimal.NewFromFloat(s.GetRspMinAmt()).RoundBank(2).RoundBank(2).RoundBank(2)
		res[i].RSPAmtApc = decimal.NewFromFloat(s.GetRspAmtApc()).RoundBank(2).RoundBank(2).RoundBank(2)
		res[i].RSPMaxAmt = decimal.NewFromFloat(s.GetRspMaxAmt()).RoundBank(2).RoundBank(2).RoundBank(2)
		res[i].RSPMinAmtTxC = decimal.NewFromFloat(s.GetRspMinAmtTxC()).RoundBank(2).RoundBank(2).RoundBank(2)
		res[i].RSPAmtApcTxC = decimal.NewFromFloat(s.GetRspAmtApcTxC()).RoundBank(2).RoundBank(2).RoundBank(2)
		res[i].RSPMaxAmtTxC = decimal.NewFromFloat(s.GetRspMaxAmtTxC()).RoundBank(2).RoundBank(2).RoundBank(2)
		res[i].RSPConFailCnt = int16(s.GetRspConFailCnt())
		res[i].RSPFeeRate = decimal.NewFromFloat(s.GetRspFeeRate()).RoundBank(4).RoundBank(4).RoundBank(4)
		res[i].IsRSPChgAmt = s.GetIsRspChgAmt()
		res[i].IsRSPShareSet = s.GetIsRspShareSet()
		res[i].IsDividend = s.GetIsDividend()
		res[i].DivType = s.GetDivType()
		res[i].DivMinAmt = decimal.NewFromFloat(s.GetDivMinAmt()).RoundBank(2).RoundBank(2).RoundBank(2)
		res[i].DivCycle = s.GetDivCycle()
		res[i].DivPayPlatform = s.GetDivPayPlatform()
		res[i].TaxFormat = s.GetTaxFormat()
		res[i].LakhSubsAmt = decimal.NewFromFloat(s.GetLakhSubsAmt()).RoundBank(2).RoundBank(2).RoundBank(2)
		res[i].LakhRdmAmt = decimal.NewFromFloat(s.GetLakhRdmAmt()).RoundBank(2).RoundBank(2).RoundBank(2)
	}
	return res
}

func mapDTAFNDFundDetailEditList(srcs []*fndv1.DTAFNDFundDetail) []db.DTAFNDFundDetailEdit {
	res := make([]db.DTAFNDFundDetailEdit, len(srcs))
	for i, s := range srcs {
		res[i].SysCoID = s.GetSysCoId()
		res[i].PrtFundCode = s.GetPrtFundCode()
		res[i].FundCode = s.GetFundCode()
		res[i].FundCry = s.GetFundCry()
		res[i].FundInShName = s.GetFundInShName()
		res[i].FundShMName = s.GetFundShMName()
		res[i].FundShSName = s.GetFundShSName()
		res[i].FHFundShare = s.GetFhFundShare()
		res[i].ShareFirSaleDate = parseTime(s.GetShareFirSaleDate())
		res[i].SMARTFundCode = s.GetSmartFundCode()
		res[i].CSDFundCode = s.GetCsdFundCode()
		res[i].SITCAFundCode = s.GetSitcaFundCode()
		res[i].FAFundCode = s.GetFaFundCode()
		res[i].FISFeeCode = s.GetFisFeeCode()
		res[i].ISINCode = s.GetIsinCode()
		res[i].IsTISA = s.GetIsTisa()
		res[i].IsIShare = s.GetIsIShare()
		res[i].DivCalcType = s.GetDivCalcType()
		res[i].SubsFeeType = s.GetSubsFeeType()
		res[i].SubsMinAmt = decimal.NewFromFloat(s.GetSubsMinAmt()).RoundBank(2).RoundBank(2).RoundBank(2)
		res[i].SubsAmtApc = decimal.NewFromFloat(s.GetSubsAmtApc()).RoundBank(2).RoundBank(2).RoundBank(2)
		res[i].SwMinAmt = decimal.NewFromFloat(s.GetSwMinAmt()).RoundBank(2).RoundBank(2).RoundBank(2)
		res[i].IsDistributionFee = s.GetIsDistributionFee()
		res[i].MatureYear = int16(s.GetMatureYear())
		res[i].MatureFund = s.GetMatureFund()
		res[i].CSDSubsSetlDate = parseTime(s.GetCsdSubsSetlDate())
		res[i].RdmMinUnitSur = decimal.NewFromFloat(s.GetRdmMinUnitSur()).RoundBank(4).RoundBank(4).RoundBank(4)
		res[i].RdmMinUnit = decimal.NewFromFloat(s.GetRdmMinUnit()).RoundBank(4).RoundBank(4).RoundBank(4)
		res[i].RdmMinAmt = decimal.NewFromFloat(s.GetRdmMinAmt()).RoundBank(2).RoundBank(2).RoundBank(2)
		res[i].CSDRdmSetlDate = parseTime(s.GetCsdRdmSetlDate())
		res[i].IsShareRSP = s.GetIsShareRsp()
		res[i].RSPMinAmt = decimal.NewFromFloat(s.GetRspMinAmt()).RoundBank(2).RoundBank(2).RoundBank(2)
		res[i].RSPAmtApc = decimal.NewFromFloat(s.GetRspAmtApc()).RoundBank(2).RoundBank(2).RoundBank(2)
		res[i].RSPMaxAmt = decimal.NewFromFloat(s.GetRspMaxAmt()).RoundBank(2).RoundBank(2).RoundBank(2)
		res[i].SubsMinAmtTxC = decimal.NewFromFloat(s.GetSubsMinAmtTxC()).RoundBank(2).RoundBank(2).RoundBank(2)
		res[i].SubsAmtApcTxC = decimal.NewFromFloat(s.GetSubsAmtApcTxC()).RoundBank(2).RoundBank(2).RoundBank(2)
		res[i].SwMinAmtTxC = decimal.NewFromFloat(s.GetSwMinAmtTxC()).RoundBank(2).RoundBank(2).RoundBank(2)
		res[i].RSPMinAmtTxC = decimal.NewFromFloat(s.GetRspMinAmtTxC()).RoundBank(2).RoundBank(2).RoundBank(2)
		res[i].RSPAmtApcTxC = decimal.NewFromFloat(s.GetRspAmtApcTxC()).RoundBank(2).RoundBank(2).RoundBank(2)
		res[i].RSPMaxAmtTxC = decimal.NewFromFloat(s.GetRspMaxAmtTxC()).RoundBank(2).RoundBank(2).RoundBank(2)
		res[i].IsAdditional = s.GetIsAdditional()
		res[i].IsOld = s.GetIsOld()
	}
	return res
}

func mapDTAFNDFundAccountEditList(srcs []*fndv1.DTAFNDFundAccount) []db.DTAFNDFundAccountEdit {
	res := make([]db.DTAFNDFundAccountEdit, len(srcs))
	for i, s := range srcs {
		res[i].SysCoID = s.GetSysCoId()
		res[i].PrtFundCode = s.GetPrtFundCode()
		res[i].FundCode = s.GetFundCode()
		res[i].DataSeq = int16(s.GetDataSeq())
		res[i].CusBankBrh = s.GetCusBankBrh()
		res[i].FundAccount = s.GetFundAccount()
		res[i].SRemitVCode = s.GetSRemitVCode()
		res[i].AccountUseSet = s.GetAccountUseSet()
		res[i].Remark = s.GetRemark()
	}
	return res
}

func mapDTAFNDFundRelGroupEditList(srcs []*fndv1.DTAFNDFundRelGroup) []db.DTAFNDFundRelGroupEdit {
	res := make([]db.DTAFNDFundRelGroupEdit, len(srcs))
	for i, s := range srcs {
		res[i].SysCoID = s.GetSysCoId()
		res[i].PrtFundCode = s.GetPrtFundCode()
		res[i].FundCode = s.GetFundCode()
		res[i].FundGroupNo = s.GetFundGroupNo()
		res[i].GroupTypeNo = s.GetGroupTypeNo()
	}
	return res
}

func mapDTAFNDFundDisclosureFeeEditList(srcs []*fndv1.DTAFNDFundDisclosureFee) []db.DTAFNDFundDisclosureFeeEdit {
	res := make([]db.DTAFNDFundDisclosureFeeEdit, len(srcs))
	for i := range srcs {
		// DTAFNDFundDisclosureFee has minimal fields — PK columns only
		// Future: map additional fields as spec evolves
		_ = i
	}
	return res
}

func mapDTAFNDFundTxCryEditList(srcs []*fndv1.DTAFNDFundTxCry) []db.DTAFNDFundTxCryEdit {
	res := make([]db.DTAFNDFundTxCryEdit, len(srcs))
	for i, s := range srcs {
		res[i].SysCoID = s.GetSysCoId()
		res[i].PrtFundCode = s.GetPrtFundCode()
		res[i].FundCode = s.GetFundCode()
		res[i].CryID = s.GetCryId()
		res[i].CryDataSrc = s.GetCryDataSrc()
	}
	return res
}

func mapDTAFNDFundShortInfEditList(srcs []*fndv1.DTAFNDFundShortInf) []db.DTAFNDFundShortInfEdit {
	res := make([]db.DTAFNDFundShortInfEdit, len(srcs))
	for i, s := range srcs {
		res[i].SysCoID = s.GetSysCoId()
		res[i].PrtFundCode = s.GetPrtFundCode()
		res[i].RdmRangeType = s.GetRdmRangeType()
		res[i].ShortDateType = s.GetShortDateType()
		res[i].RdmBaseID = s.GetRdmBaseId()
		res[i].SubsBaseID = s.GetSubsBaseId()
		res[i].ShortCalcID = s.GetShortCalcId()
	}
	return res
}

func mapDTAFNDFundAntiDilutionEditList(srcs []*fndv1.DTAFNDFundAntiDilution) []db.DTAFNDFundAntiDilutionEdit {
	res := make([]db.DTAFNDFundAntiDilutionEdit, len(srcs))
	for i, s := range srcs {
		res[i].SysCoID = s.GetSysCoId()
		res[i].PrtFundCode = s.GetPrtFundCode()
		res[i].AntiDilSetType = s.GetAntiDilSetType()
		res[i].EffDate = parseTime(s.GetEffDate())
		res[i].TermDate = parseTime(s.GetTermDate())
		res[i].AntiDilTrigger = decimal.NewFromFloat(s.GetAntiDilTrigger()).RoundBank(2).RoundBank(2).RoundBank(2)
		res[i].AntiDilFeeRate = decimal.NewFromFloat(s.GetAntiDilFeeRate()).RoundBank(4).RoundBank(4).RoundBank(4)
		res[i].AdjRsn = s.GetAdjRsn()
		res[i].Remark = s.GetRemark()
	}
	return res
}
