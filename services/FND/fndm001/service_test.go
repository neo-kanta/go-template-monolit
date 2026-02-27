package fndm001_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
	"go-transfer-agent/services/fnd/fndm001"
	"go-transfer-agent/services/fnd/testutil"
)

type ServiceTestSuite struct {
	suite.Suite
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

// ═══════════════════════════════════════════════════════════════════
// 1. AUD / Save Operations (APIFNDM001Post/Put)
// ═══════════════════════════════════════════════════════════════════

func (suite *ServiceTestSuite) TestSaveFundInfo_Success() {
	svc := fndm001.NewService(testutil.SetupTestDB(suite.T()), testutil.NewTestLogger())
	req := &fndv1.SaveFundInfoRequest{
		Master: &fndv1.TAFNDFundInfo{
			SysCoId:     "SWSTD",
			PrtFundCode: "A001",
			FundMName:   "Test Parent Fund",
		},
		FundInfo: &fndv1.DTAFNDFundInfo{
			SysCoId:       "SWSTD",
			PrtFundCode:   "A001",
			SitcaFundType: "TYPE01",
		},
	}

	res := svc.SaveFundInfo(context.Background(), req)
	if !res.Success {
		suite.T().Fatalf("Expected check to succeed, got %v", res.Message)
	}
	assert.Equal(suite.T(), "OK", res.ReturnCode)

	// Verify it was stored
	qRes, err := svc.GetDataByDataID(context.Background(), &fndv1.GetDataRequest{DataId: "PF001"})
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Test Parent Fund", qRes.Master.FundMName)
}

func (suite *ServiceTestSuite) TestSaveFundInfo_MissingKeys() {
	svc := fndm001.NewService(testutil.SetupTestDB(suite.T()), testutil.NewTestLogger())
	req := &fndv1.SaveFundInfoRequest{
		Master: &fndv1.TAFNDFundInfo{
			SysCoId:     "", // Missing
			PrtFundCode: "PF001",
		},
	}

	res := svc.SaveFundInfo(context.Background(), req)
	if res.Success {
		suite.T().Fatal("Expected check to fail due to missing SysCoID")
	}
	assert.Equal(suite.T(), "VALIDATION_ERROR", res.ReturnCode)
}

// ═══════════════════════════════════════════════════════════════════
// 2. Query Operations (APIFNDM001Get)
// ═══════════════════════════════════════════════════════════════════

func (suite *ServiceTestSuite) TestQueryFundInfo_APIFNDM001Get() {
	svc := fndm001.NewService(testutil.SetupTestDB(suite.T()), testutil.NewTestLogger())
	// Pre-populate
	svc.SaveFundInfo(context.Background(), &fndv1.SaveFundInfoRequest{
		Master: &fndv1.TAFNDFundInfo{
			SysCoId: "C1", PrtFundCode: "P1", FundInShName: "Fund1",
		},
	})
	svc.SaveFundInfo(context.Background(), &fndv1.SaveFundInfoRequest{
		Master: &fndv1.TAFNDFundInfo{
			SysCoId: "C1", PrtFundCode: "P2", FundInShName: "Fund2",
		},
	})
	svc.SaveFundInfo(context.Background(), &fndv1.SaveFundInfoRequest{
		Master: &fndv1.TAFNDFundInfo{
			SysCoId: "C2", PrtFundCode: "P1", FundInShName: "Fund3",
		},
	})

	res := svc.QueryFundInfo("C1")
	if len(res.ResultList) != 2 {
		suite.T().Fatalf("Expected 2 results for SysCoId C1, got %d", len(res.ResultList))
	}

	foundP1 := false
	for _, f := range res.ResultList {
		if f.PrtFundCode == "P1" {
			foundP1 = true
			assert.Equal(suite.T(), "Fund1", f.FundInShName)
		}
	}
	assert.True(suite.T(), foundP1)
}

// ═══════════════════════════════════════════════════════════════════
// 3. Check Functions (KFNDM)
// ═══════════════════════════════════════════════════════════════════

func (suite *ServiceTestSuite) TestTACKCSDPrtFundSetlDate_KFNDM00101() {
	svc := fndm001.NewService(testutil.SetupTestDB(suite.T()), testutil.NewTestLogger())

	// Valid date (future)
	validDate := time.Now().AddDate(0, 0, 5).Format("2006-01-02")
	res := svc.TACKCSDPrtFundSetlDate("COMP1", validDate, nil)
	assert.Equal(suite.T(), "OK", res.ReturnCode)

	// Invalid date (past/too soon)
	invalidDate := time.Now().Format("2006-01-02")

	// Without FundCode (returns FNDM001009)
	resCode009 := svc.TACKCSDPrtFundSetlDate("COMP1", invalidDate, nil)
	assert.Equal(suite.T(), "FNDM001009", resCode009.ReturnCode)

	// With FundCode (returns FNDM001014)
	fc := "F123"
	resCode014 := svc.TACKCSDPrtFundSetlDate("COMP1", invalidDate, &fc)
	assert.Equal(suite.T(), "FNDM001014", resCode014.ReturnCode)
}

func (suite *ServiceTestSuite) TestTACKFundCodeExist_KFNDM00103() {
	svc := fndm001.NewService(testutil.SetupTestDB(suite.T()), testutil.NewTestLogger())
	svc.SaveFundInfo(context.Background(), &fndv1.SaveFundInfoRequest{
		Master: &fndv1.TAFNDFundInfo{SysCoId: "C1", PrtFundCode: "P1"},
		FundDetails: []*fndv1.DTAFNDFundDetail{
			{SysCoId: "C1", PrtFundCode: "P1", FundCode: "EXISTING_FUND"},
		},
	})

	// Check existing
	resExist := svc.TACKFundCodeExist("C1", "EXISTING_FUND")
	assert.Equal(suite.T(), "TA00030", resExist.ReturnCode)

	// Check non-existing in same company
	resNew := svc.TACKFundCodeExist("C1", "NEW_FUND")
	assert.Equal(suite.T(), "OK", resNew.ReturnCode)

	// Check existing fund code but in different company
	resOtherCo := svc.TACKFundCodeExist("C2", "EXISTING_FUND")
	assert.Equal(suite.T(), "OK", resOtherCo.ReturnCode)
}

// ═══════════════════════════════════════════════════════════════════
// 4. Search/Info Functions (GFNDM)
// ═══════════════════════════════════════════════════════════════════

func (suite *ServiceTestSuite) TestTAGetDPrtFundIssueCry_GFNDM00101() {
	svc := fndm001.NewService(testutil.SetupTestDB(suite.T()), testutil.NewTestLogger())
	svc.SaveFundInfo(context.Background(), &fndv1.SaveFundInfoRequest{
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
		suite.T().Fatalf("Expected 2 currencies, got %d", len(res.Master))
	}

	// Filter by FeeChargeType (B = TWD only)
	resFeeB := svc.TAGetDPrtFundIssueCry("C1", false, "P1", "B", "", "SMP")
	if len(resFeeB.Master) != 1 || resFeeB.Master[0].CryId != "TWD" {
		suite.T().Errorf("Expected TWD only for FeeChargeType B, got %v", resFeeB.Master)
	}

	// With IsAll = true
	resAll := svc.TAGetDPrtFundIssueCry("C1", true, "P1", "", "", "SMP")
	if len(resAll.Master) != 3 {
		suite.T().Fatalf("Expected 3 results (including ALL), got %d", len(resAll.Master))
	}
	if resAll.Master[0].CryId != "*" {
		suite.T().Errorf("Expected first result to be ALL (*), got %v", resAll.Master[0].CryId)
	}
}

func (suite *ServiceTestSuite) TestTAGetDTxCry_GFNDM00104() {
	svc := fndm001.NewService(testutil.SetupTestDB(suite.T()), testutil.NewTestLogger())
	svc.SaveFundInfo(context.Background(), &fndv1.SaveFundInfoRequest{
		Master: &fndv1.TAFNDFundInfo{SysCoId: "C1", PrtFundCode: "P1"},
		FundTxCrys: []*fndv1.DTAFNDFundTxCry{
			{SysCoId: "C1", PrtFundCode: "P1", CryId: "TWD", CryDataSrc: "1"},
			{SysCoId: "C1", PrtFundCode: "P1", CryId: "USD", CryDataSrc: "2"},
		},
	})

	// Basic query
	res := svc.TAGetDTxCry("C1", false, "P1", "", "", false, "SMP")
	if len(res.Master) != 2 {
		suite.T().Fatalf("Expected 2 tx currencies, got %d", len(res.Master))
	}

	// Filter by CryDataSrc
	resSrc1 := svc.TAGetDTxCry("C1", false, "P1", "", "1", false, "SMP")
	if len(resSrc1.Master) != 1 || resSrc1.Master[0].CryId != "TWD" {
		suite.T().Errorf("Expected TWD only for CryDataSrc 1, got %v", resSrc1.Master)
	}

	// Test IsExcludeCCY (assuming CenterCurrency is TWD)
	resExcl := svc.TAGetDTxCry("C1", false, "P1", "", "", true, "SMP")
	if len(resExcl.Master) != 1 || resExcl.Master[0].CryId != "USD" {
		suite.T().Errorf("Expected USD only when excluding center currency, got %v", resExcl.Master)
	}
}
