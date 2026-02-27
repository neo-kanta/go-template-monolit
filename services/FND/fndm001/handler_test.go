package fndm001_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
	"go-transfer-agent/services/fnd/fndm001"
	"go-transfer-agent/services/fnd/testutil"
)

type HandlerTestSuite struct {
	suite.Suite
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func setupHandler(t *testing.T) *fndm001.Handler {
	logger := testutil.NewTestLogger()
	db := testutil.SetupTestDB(t)
	svc := fndm001.NewService(db, logger)
	return fndm001.NewHandler(logger, svc)
}

// ═══════════════════════════════════════════════════════════════════
// 1. AUD / Save Operations (APIFNDM001Post/Put)
// ═══════════════════════════════════════════════════════════════════

func (suite *HandlerTestSuite) TestSaveFundInfo_Success() {
	h := setupHandler(suite.T())
	ctx := context.Background()

	req := &fndv1.SaveFundInfoRequest{
		Master: &fndv1.TAFNDFundInfo{
			SysCoId:     "SWSTD",
			PrtFundCode: "P1",
		},
	}

	res, err := h.SaveFundInfo(ctx, req)
	require.NoError(suite.T(), err)
	assert.True(suite.T(), res.Success)
}

// ═══════════════════════════════════════════════════════════════════
// 2. Query Operations (APIFNDM001Get)
// ═══════════════════════════════════════════════════════════════════

func (suite *HandlerTestSuite) TestQueryFundInfo_APIFNDM001Get() {
	h := setupHandler(suite.T())
	ctx := context.Background()

	h.SaveFundInfo(ctx, &fndv1.SaveFundInfoRequest{
		Master: &fndv1.TAFNDFundInfo{SysCoId: "SWSTD", PrtFundCode: "P1", FundInShName: "Fund1"},
	})
	h.SaveFundInfo(ctx, &fndv1.SaveFundInfoRequest{
		Master: &fndv1.TAFNDFundInfo{SysCoId: "SWSTD", PrtFundCode: "P2", FundInShName: "Fund2"},
	})

	res, err := h.QueryFundInfo(ctx, &fndv1.QueryFundInfoRequest{SysCoId: "SWSTD"})
	require.NoError(suite.T(), err)

	if len(res.ResultList) != 1 {
		suite.T().Fatalf("Expected 1 result for Company C1, got %d", len(res.ResultList))
	}
	if res.ResultList[0].PrtFundCode != "P1" {
		suite.T().Errorf("Expected PrtFundCode P1, got %v", res.ResultList[0].PrtFundCode)
	}
}

func (suite *HandlerTestSuite) TestQueryFundInfoByDataID_APIFNDM001GetMaintain() {
	h := setupHandler(suite.T())
	ctx := context.Background()

	h.SaveFundInfo(ctx, &fndv1.SaveFundInfoRequest{
		Master:   &fndv1.TAFNDFundInfo{SysCoId: "SWSTD", PrtFundCode: "P1", FundInShName: "Fund1"},
		FundInfo: &fndv1.DTAFNDFundInfo{SysCoId: "SWSTD", PrtFundCode: "P1", SitcaFundType: "T1"},
	})

	res, err := h.GetDataByDataID(ctx, &fndv1.GetDataRequest{
		DataId: "P1", // using DataId since PrtFundCode serves as data id here
	})
	require.NoError(suite.T(), err)

	if res.Master == nil || res.Master.FundInShName != "Fund1" {
		suite.T().Errorf("Expected Master FundInShName Fund1, got %v", res.Master)
	}
	if res.FundInfo == nil || res.FundInfo.SitcaFundType != "T1" {
		suite.T().Errorf("Expected FundInfo SitcaFundType T1, got %v", res.FundInfo)
	}
}

// ═══════════════════════════════════════════════════════════════════
// 3. Check Functions (KFNDM)
// ═══════════════════════════════════════════════════════════════════

func (suite *HandlerTestSuite) TestTACKCSDPrtFundSetlDate_KFNDM00101() {
	h := setupHandler(suite.T())
	ctx := context.Background()

	req := &fndv1.TACKCSDPrtFundSetlDateRequest{
		SysCoId:            "C1",
		CsdPrtFundSetlDate: "3000-01-01", // Way into the future, should be valid
	}

	res, err := h.TACKCSDPrtFundSetlDate(ctx, req)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), "OK", res.ReturnCode)
}

func (suite *HandlerTestSuite) TestTACKFundCodeExist_KFNDM00103() {
	h := setupHandler(suite.T())
	ctx := context.Background()

	h.SaveFundInfo(ctx, &fndv1.SaveFundInfoRequest{
		Master: &fndv1.TAFNDFundInfo{SysCoId: "SWSTD", PrtFundCode: "P1"},
		FundDetails: []*fndv1.DTAFNDFundDetail{
			{SysCoId: "SWSTD", PrtFundCode: "P1", FundCode: "EXIST"},
		},
	})

	res, err := h.TACKFundCodeExist(ctx, &fndv1.TACKFundCodeExistRequest{
		SysCoId:  "SWSTD",
		FundCode: "EXIST",
	})
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), "TA00030", res.ReturnCode)
}

// ═══════════════════════════════════════════════════════════════════
// 4. Search/Info Functions (GFNDM)
// ═══════════════════════════════════════════════════════════════════

func (suite *HandlerTestSuite) TestTAGetDPrtFundIssueCry_GFNDM00101() {
	h := setupHandler(suite.T())
	ctx := context.Background()

	h.SaveFundInfo(ctx, &fndv1.SaveFundInfoRequest{
		Master: &fndv1.TAFNDFundInfo{SysCoId: "SWSTD", PrtFundCode: "P1"},
		CryGroups: []*fndv1.DTAFNDFundCryGroup{
			{SysCoId: "SWSTD", PrtFundCode: "P1", IssueCry: "USD", IsRsp: "Y"},
		},
	})

	res, err := h.TAGetDPrtFundIssueCry(ctx, &fndv1.TAGetDPrtFundIssueCryRequest{
		SysCoId:     "SWSTD",
		PrtFundCode: "P1",
		QueryType:   "ALL",
	})
	require.NoError(suite.T(), err)

	if len(res.Master) != 1 {
		suite.T().Fatalf("Expected 1 currency group, got %d", len(res.Master))
	}
	if res.Master[0].CryId != "USD" {
		suite.T().Errorf("Expected USD, got %v", res.Master[0].CryId)
	}
}

func (suite *HandlerTestSuite) TestTAGetDTxCry_GFNDM00104() {
	h := setupHandler(suite.T())
	ctx := context.Background()

	h.SaveFundInfo(ctx, &fndv1.SaveFundInfoRequest{
		Master: &fndv1.TAFNDFundInfo{SysCoId: "SWSTD", PrtFundCode: "P1"},
		FundTxCrys: []*fndv1.DTAFNDFundTxCry{
			{SysCoId: "SWSTD", PrtFundCode: "P1", CryId: "EUR"},
		},
	})

	res, err := h.TAGetDTxCry(ctx, &fndv1.TAGetDTxCryRequest{
		SysCoId:     "SWSTD",
		PrtFundCode: "P1",
		QueryType:   "ALL",
	})
	require.NoError(suite.T(), err)

	if len(res.Master) != 1 {
		suite.T().Fatalf("Expected 1 transaction currency, got %d", len(res.Master))
	}
	if res.Master[0].CryId != "EUR" {
		suite.T().Errorf("Expected EUR, got %v", res.Master[0].CryId)
	}
}
