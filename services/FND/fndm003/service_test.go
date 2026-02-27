package fndm003_test

import (
	"context"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
	"go-transfer-agent/services/fnd/fndm003"
	"go-transfer-agent/services/fnd/fndm003/db"
	"go-transfer-agent/services/fnd/testutil"
)

func TestService_TAFNDSwitch_Integration(t *testing.T) {
	database := testutil.SetupTestDB(t)
	logger := testutil.NewTestLogger()

	// 1. AutoMigrate schemas
	err := database.AutoMigrate(
		&db.DTAFNDSwitch{},
		&db.DTAFNDSwitchFund{},
		&db.DTAFNDFundFeeSwitch{},
		&db.DTAFNDSwitchCry{},
	)
	require.NoError(t, err)

	// 2. Clean up
	database.Exec(`DELETE FROM "TA_STD_TH"."DTA_FND_SwitchCry"`)
	database.Exec(`DELETE FROM "TA_STD_TH"."DTA_FND_FundFeeSwitch"`)
	database.Exec(`DELETE FROM "TA_STD_TH"."DTA_FND_SwitchFund"`)
	database.Exec(`DELETE FROM "TA_STD_TH"."DTA_FND_Switch"`)

	// 3. Seed data
	master := db.DTAFNDSwitch{
		SysCoID:     "C01",
		PrtFundCode: "F1",
		IsSwitchIn:  "Y",
	}
	if err := database.Create(&master).Error; err != nil {
		t.Fatalf("Failed to seed master: %v", err)
	}

	funds := []db.DTAFNDSwitchFund{
		{SysCoID: "C01", PrtFundCode: "F1", SwOPrtFundCode: "F2"},
	}
	if err := database.Create(&funds).Error; err != nil {
		t.Fatalf("Failed to seed funds: %v", err)
	}

	fees := []db.DTAFNDFundFeeSwitch{
		{SysCoID: "C01", PrtFundCode: "F1", FundCode: "F2", SwFundType: "TYPE1", SwDiscType: "DISC1", SwitchRate: decimal.NewFromFloat(0.5).RoundBank(4)},
	}
	if err := database.Create(&fees).Error; err != nil {
		t.Fatalf("Failed to seed fees: %v", err)
	}

	crys := []db.DTAFNDSwitchCry{
		{SysCoID: "C01", PrtFundCode: "F1", SwIFundCry: "THB", SwOFundCrySet: "USD,EUR"},
	}
	if err := database.Create(&crys).Error; err != nil {
		t.Fatalf("Failed to seed crys: %v", err)
	}

	// 4. Test Service
	svc := fndm003.NewService(database, logger)
	ctx := context.Background()

	req := &fndv1.TAFNDSwitchRequest{
		SysCoId:     "C01",
		PrtFundCode: "F1",
		IsSwitchIn:  "Y",
	}

	res, err := svc.TAFNDSwitch(ctx, req)
	require.NoError(t, err)

	// 5. Assertions
	if len(res.ResultList) != 1 {
		t.Fatalf("Expected 1 master result, got %d", len(res.ResultList))
	}
	if res.ResultList[0].IsSwitchIn != "Y" {
		t.Errorf("Expected IsSwitchIn 'Y', got '%v'", res.ResultList[0].IsSwitchIn)
	}

	if len(res.SwitchFundList) != 1 || res.SwitchFundList[0].SwOPrtFundCode != "F2" {
		t.Errorf("SwitchFundList mapping failed")
	}

	if len(res.FundSwitchFeeList) != 1 || res.FundSwitchFeeList[0].SwitchRate != 0.5 {
		t.Errorf("FundSwitchFeeList mapping failed")
	}

	if len(res.SwitchCryList) != 1 {
		t.Fatalf("SwitchCryList mapping failed")
	}
	sc := res.SwitchCryList[0]
	assert.Equal(t, "THB", sc.SwIFundCry)
	if len(sc.SwOFundCrySetNm) != 2 || sc.SwOFundCrySetNm[0] != "USD" || sc.SwOFundCrySetNm[1] != "EUR" {
		t.Errorf("SwOFundCrySet split mapping failed: %v", sc.SwOFundCrySetNm) // Based on set splitting logic
	}
}
