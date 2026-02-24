package fndm001_test

import (
	"testing"
	"time"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
	"go-transfer-agent/services/fnd/fndm001"
)

// ═══════════════════════════════════════════════════════════════════
// 1. AUD / Save Operations (APIFNDM001Post/Put)
// ═══════════════════════════════════════════════════════════════════

func TestService_SaveFundInfo_Success(t *testing.T) {
	svc := fndm001.NewService()
	req := &fndv1.SaveFundInfoRequest{
		Master: &fndv1.TAFNDFundInfo{
			SysCoId:     "COMP01",
			PrtFundCode: "PF001",
			FundMName:   "Test Parent Fund",
		},
		FundInfo: &fndv1.DTAFNDFundInfo{
			SysCoId:       "COMP01",
			PrtFundCode:   "PF001",
			SitcaFundType: "TYPE01",
		},
	}

	res := svc.SaveFundInfo(req)
	if !res.Success {
		t.Fatalf("Expected check to succeed, got %v", res.Message)
	}
	if res.ReturnCode != "OK" {
		t.Errorf("Expected ReturnCode OK, got %v", res.ReturnCode)
	}

	// Verify it was stored
	qRes, err := svc.QueryFundInfoByDataID("COMP01", "PF001")
	if err != nil {
		t.Fatalf("Failed to retrieve saved fund: %v", err)
	}
	if qRes.Master.FundMName != "Test Parent Fund" {
		t.Errorf("Expected FundMName 'Test Parent Fund', got '%v'", qRes.Master.FundMName)
	}
}

func TestService_SaveFundInfo_MissingKeys(t *testing.T) {
	svc := fndm001.NewService()
	req := &fndv1.SaveFundInfoRequest{
		Master: &fndv1.TAFNDFundInfo{
			SysCoId:     "", // Missing
			PrtFundCode: "PF001",
		},
	}

	res := svc.SaveFundInfo(req)
	if res.Success {
		t.Fatal("Expected check to fail due to missing SysCoID")
	}
	if res.ReturnCode != "VALIDATION_ERROR" {
		t.Errorf("Expected ReturnCode VALIDATION_ERROR, got %v", res.ReturnCode)
	}
}

// ═══════════════════════════════════════════════════════════════════
// 2. Query Operations (APIFNDM001Get)
// ═══════════════════════════════════════════════════════════════════

func TestService_QueryFundInfo_APIFNDM001Get(t *testing.T) {
	svc := fndm001.NewService()
	// Pre-populate
	svc.SaveFundInfo(&fndv1.SaveFundInfoRequest{
		Master: &fndv1.TAFNDFundInfo{
			SysCoId: "C1", PrtFundCode: "P1", FundInShName: "Fund1",
		},
	})
	svc.SaveFundInfo(&fndv1.SaveFundInfoRequest{
		Master: &fndv1.TAFNDFundInfo{
			SysCoId: "C1", PrtFundCode: "P2", FundInShName: "Fund2",
		},
	})
	svc.SaveFundInfo(&fndv1.SaveFundInfoRequest{
		Master: &fndv1.TAFNDFundInfo{
			SysCoId: "C2", PrtFundCode: "P1", FundInShName: "Fund3",
		},
	})

	res := svc.QueryFundInfo("C1")
	if len(res.ResultList) != 2 {
		t.Fatalf("Expected 2 results for SysCoId C1, got %d", len(res.ResultList))
	}

	foundP1 := false
	for _, f := range res.ResultList {
		if f.PrtFundCode == "P1" {
			foundP1 = true
			if f.FundInShName != "Fund1" {
				t.Errorf("Expected FundInShName Fund1, got %v", f.FundInShName)
			}
		}
	}
	if !foundP1 {
		t.Error("Expected to find PrtFundCode P1 in results")
	}
}

// ═══════════════════════════════════════════════════════════════════
// 3. Check Functions (KFNDM)
// ═══════════════════════════════════════════════════════════════════

func TestService_TACKCSDPrtFundSetlDate_KFNDM00101(t *testing.T) {
	svc := fndm001.NewService()

	// Valid date (future)
	validDate := time.Now().AddDate(0, 0, 5).Format("2006-01-02")
	res := svc.TACKCSDPrtFundSetlDate("COMP1", validDate, nil)
	if res.ReturnCode != "OK" {
		t.Errorf("Expected OK for valid future date, got %v", res.ReturnCode)
	}

	// Invalid date (past/too soon)
	invalidDate := time.Now().Format("2006-01-02")

	// Without FundCode (returns FNDM001009)
	resCode009 := svc.TACKCSDPrtFundSetlDate("COMP1", invalidDate, nil)
	if resCode009.ReturnCode != "FNDM001009" {
		t.Errorf("Expected ReturnCode FNDM001009 without FundCode, got %v", resCode009.ReturnCode)
	}

	// With FundCode (returns FNDM001014)
	fc := "F123"
	resCode014 := svc.TACKCSDPrtFundSetlDate("COMP1", invalidDate, &fc)
	if resCode014.ReturnCode != "FNDM001014" {
		t.Errorf("Expected ReturnCode FNDM001014 with FundCode, got %v", resCode014.ReturnCode)
	}
}

func TestService_TACKFundCodeExist_KFNDM00103(t *testing.T) {
	svc := fndm001.NewService()
	svc.SaveFundInfo(&fndv1.SaveFundInfoRequest{
		Master: &fndv1.TAFNDFundInfo{SysCoId: "C1", PrtFundCode: "P1"},
		FundDetails: []*fndv1.DTAFNDFundDetail{
			{SysCoId: "C1", PrtFundCode: "P1", FundCode: "EXISTING_FUND"},
		},
	})

	// Check existing
	resExist := svc.TACKFundCodeExist("C1", "EXISTING_FUND")
	if resExist.ReturnCode != "TA00030" {
		t.Errorf("Expected ReturnCode TA00030 for existing fund, got %v", resExist.ReturnCode)
	}

	// Check non-existing in same company
	resNew := svc.TACKFundCodeExist("C1", "NEW_FUND")
	if resNew.ReturnCode != "OK" {
		t.Errorf("Expected OK for new fund code, got %v", resNew.ReturnCode)
	}

	// Check existing fund code but in different company
	resOtherCo := svc.TACKFundCodeExist("C2", "EXISTING_FUND")
	if resOtherCo.ReturnCode != "OK" {
		t.Errorf("Expected OK for existing fund code in different company, got %v", resOtherCo.ReturnCode)
	}
}

// ═══════════════════════════════════════════════════════════════════
// 4. Search/Info Functions (GFNDM)
// ═══════════════════════════════════════════════════════════════════

func TestService_TAGetDPrtFundIssueCry_GFNDM00101(t *testing.T) {
	svc := fndm001.NewService()
	svc.SaveFundInfo(&fndv1.SaveFundInfoRequest{
		Master: &fndv1.TAFNDFundInfo{SysCoId: "C1", PrtFundCode: "P1"},
		CryGroups: []*fndv1.DTAFNDFundCryGroup{
			{SysCoId: "C1", PrtFundCode: "P1", IssueCry: "TWD", IsRsp: "Y"},
			{SysCoId: "C1", PrtFundCode: "P1", IssueCry: "USD", IsRsp: "N"},
		},
		FundDetails: []*fndv1.DTAFNDFundDetail{
			{SysCoId: "C1", PrtFundCode: "P1", FundCry: "TWD", SubsFeeType: "B"},
			{SysCoId: "C1", PrtFundCode: "P1", FundCry: "USD", SubsFeeType: "P"},
		},
	})

	// Basic query
	res := svc.TAGetDPrtFundIssueCry("C1", false, "P1", "", "", "SMP")
	if len(res.Master) != 2 {
		t.Fatalf("Expected 2 currencies, got %d", len(res.Master))
	}

	// Filter by FeeChargeType (B = TWD only)
	resFeeB := svc.TAGetDPrtFundIssueCry("C1", false, "P1", "B", "", "SMP")
	if len(resFeeB.Master) != 1 || resFeeB.Master[0].CryId != "TWD" {
		t.Errorf("Expected TWD only for FeeChargeType B, got %v", resFeeB.Master)
	}

	// With IsAll = true
	resAll := svc.TAGetDPrtFundIssueCry("C1", true, "P1", "", "", "SMP")
	if len(resAll.Master) != 3 {
		t.Fatalf("Expected 3 results (including ALL), got %d", len(resAll.Master))
	}
	if resAll.Master[0].CryId != "*" {
		t.Errorf("Expected first result to be ALL (*), got %v", resAll.Master[0].CryId)
	}
}

func TestService_TAGetDTxCry_GFNDM00104(t *testing.T) {
	svc := fndm001.NewService()
	svc.SaveFundInfo(&fndv1.SaveFundInfoRequest{
		Master: &fndv1.TAFNDFundInfo{SysCoId: "C1", PrtFundCode: "P1"},
		FundTxCrys: []*fndv1.DTAFNDFundTxCry{
			{SysCoId: "C1", PrtFundCode: "P1", CryId: "TWD", CryDataSrc: "1"},
			{SysCoId: "C1", PrtFundCode: "P1", CryId: "USD", CryDataSrc: "2"},
		},
	})

	// Basic query
	res := svc.TAGetDTxCry("C1", false, "P1", "", "", false, "SMP")
	if len(res.Master) != 2 {
		t.Fatalf("Expected 2 tx currencies, got %d", len(res.Master))
	}

	// Filter by CryDataSrc
	resSrc1 := svc.TAGetDTxCry("C1", false, "P1", "", "1", false, "SMP")
	if len(resSrc1.Master) != 1 || resSrc1.Master[0].CryId != "TWD" {
		t.Errorf("Expected TWD only for CryDataSrc 1, got %v", resSrc1.Master)
	}

	// Test IsExcludeCCY (assuming CenterCurrency is TWD)
	resExcl := svc.TAGetDTxCry("C1", false, "P1", "", "", true, "SMP")
	if len(resExcl.Master) != 1 || resExcl.Master[0].CryId != "USD" {
		t.Errorf("Expected USD only when excluding center currency, got %v", resExcl.Master)
	}
}
