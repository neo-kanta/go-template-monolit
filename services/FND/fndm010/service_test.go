package fndm010_test

import (
	"context"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
	"go-transfer-agent/services/fnd/fndm010"
	"go-transfer-agent/services/fnd/fndm010/db"
	"go-transfer-agent/services/fnd/testutil"
)

func TestService_TAFNDIShareFundFee_Integration(t *testing.T) {
	database := testutil.SetupTestDB(t)
	logger := testutil.NewTestLogger()

	// 1. AutoMigrate schemas
	err := database.AutoMigrate(
		&db.DTAFNDIShareFundFee{},
		&db.DTAFNDIShareFundFeeSub{},
	)
	require.NoError(t, err)

	// 2. Clean up
	database.Exec(`DELETE FROM "TA_STD_TH"."DTA_FND_IShareFundFeeSub"`)
	database.Exec(`DELETE FROM "TA_STD_TH"."DTA_FND_IShareFundFee"`)

	// 3. Seed data
	master := db.DTAFNDIShareFundFee{
		SysCoID:  "C01",
		FundCode: "F1",
	}
	if err := database.Create(&master).Error; err != nil {
		t.Fatalf("Failed to seed master: %v", err)
	}

	sub1 := db.DTAFNDIShareFundFeeSub{
		SysCoID:       "C01",
		FundCode:      "F1",
		RangeAmtAbove: decimal.NewFromFloat(1000).RoundBank(2),
		SubsFeeRate:   decimal.NewFromFloat(1.0).RoundBank(4),
	}
	sub2 := db.DTAFNDIShareFundFeeSub{
		SysCoID:       "C01",
		FundCode:      "F1",
		RangeAmtAbove: decimal.NewFromFloat(5000).RoundBank(2),
		SubsFeeRate:   decimal.NewFromFloat(0.5).RoundBank(4),
	}

	if err := database.Create(&[]db.DTAFNDIShareFundFeeSub{sub1, sub2}).Error; err != nil {
		t.Fatalf("Failed to seed sub details: %v", err)
	}

	// 4. Test Service
	svc := fndm010.NewService(database, logger)
	ctx := context.Background()

	req := &fndv1.TAFNDIShareFundFeeRequest{
		SysCoId:  "C01",
		FundCode: "F1",
	}

	res, err := svc.TAFNDIShareFundFee(ctx, req)
	require.NoError(t, err)

	// 5. Assertions
	if len(res.ResultList) != 1 {
		t.Fatalf("Expected 1 master result, got %d", len(res.ResultList))
	}

	resultMaster := res.ResultList[0]
	assert.Equal(t, "C01", resultMaster.SysCoId)
	assert.Equal(t, "F1", resultMaster.FundCode)

	if len(res.IShareFundFeeList) != 2 {
		t.Fatalf("Expected 2 sub details, got %d", len(res.IShareFundFeeList))
	}

	rates := map[float64]bool{}
	for _, s := range res.IShareFundFeeList {
		rates[s.SubsFeeRate] = true
	}

	if !rates[1.0] || !rates[0.5] {
		t.Errorf("Expected fee rates 1.0 and 0.5, got: %v", rates)
	}
}
