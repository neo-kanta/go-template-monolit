package fndm002_test

import (
	"context"
	"testing"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
	"go-transfer-agent/services/fnd/fndm002"
	"go-transfer-agent/services/fnd/fndm002/db"
	"go-transfer-agent/services/fnd/testutil"
)

func TestService_TAFNDFundFee_Integration(t *testing.T) {
	database := testutil.SetupTestDB(t)
	logger := testutil.NewTestLogger()

	// 1. AutoMigrate schemas
	err := database.AutoMigrate(
		&db.DTAFNDFundFee{},
		&db.DTAFNDFundFeeSub{},
		&db.DTAFNDFundFeeSubDtl{},
		&db.DTAFNDFundFeeCDSC{},
		&db.DTAFNDFundFeeCDSCDtl{},
		&db.DTAFNDFundFeeBack{},
		&db.DTAFNDFundFeeBackDtl{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate FNDM002 schemas: %v", err)
	}

	// 2. Clean up
	database.Exec(`DELETE FROM "TA_STD_TH"."DTA_FND_FundFeeBackDtl"`)
	database.Exec(`DELETE FROM "TA_STD_TH"."DTA_FND_FundFeeBack"`)
	database.Exec(`DELETE FROM "TA_STD_TH"."DTA_FND_FundFeeCDSCDtl"`)
	database.Exec(`DELETE FROM "TA_STD_TH"."DTA_FND_FundFeeCDSC"`)
	database.Exec(`DELETE FROM "TA_STD_TH"."DTA_FND_FundFeeSubDtl"`)
	database.Exec(`DELETE FROM "TA_STD_TH"."DTA_FND_FundFeeSub"`)
	database.Exec(`DELETE FROM "TA_STD_TH"."DTA_FND_FundFee"`)

	// 3. Seed data
	master := db.DTAFNDFundFee{
		SysCoID:     "C01",
		PrtFundCode: "F1",
	}
	if err := database.Create(&master).Error; err != nil {
		t.Fatalf("Failed to seed master: %v", err)
	}

	subs := []db.DTAFNDFundFeeSub{
		{SysCoID: "C01", PrtFundCode: "F1", FundCode: "F2", CryID: "THB"},
	}
	if err := database.Create(&subs).Error; err != nil {
		t.Fatalf("Failed to seed subs: %v", err)
	}

	subDtls := []db.DTAFNDFundFeeSubDtl{
		{SysCoID: "C01", PrtFundCode: "F1", FundCode: "F2", CryID: "THB", SubsFeeRate: 1.5},
	}
	if err := database.Create(&subDtls).Error; err != nil {
		t.Fatalf("Failed to seed sub dtls: %v", err)
	}

	cdsc := []db.DTAFNDFundFeeCDSC{
		{SysCoID: "C01", PrtFundCode: "F1", FundCode: "F2", MatureYear: "2026", SubsCalcID: "CALC1"},
	}
	if err := database.Create(&cdsc).Error; err != nil {
		t.Fatalf("Failed to seed cdsc: %v", err)
	}

	back := []db.DTAFNDFundFeeBack{
		{SysCoID: "C01", PrtFundCode: "F1", FundCode: "F2", HoldPeriodYear: "3", SubsCalcID: "CALC2"},
	}
	if err := database.Create(&back).Error; err != nil {
		t.Fatalf("Failed to seed back: %v", err)
	}

	// 4. Test Service
	svc := fndm002.NewService(database, logger)
	ctx := context.Background()

	req := &fndv1.TAFNDFundFeeRequest{
		SysCoId:     "C01",
		PrtFundCode: "F1",
	}

	res, err := svc.TAFNDFundFee(ctx, req)
	if err != nil {
		t.Fatalf("Service error: %v", err)
	}

	// 5. Assertions
	if len(res.ResultList) != 1 {
		t.Fatalf("Expected 1 master result, got %d", len(res.ResultList))
	}

	if len(res.FundFeeSubList) != 1 {
		t.Fatalf("Expected 1 sub, got %d", len(res.FundFeeSubList))
	}
	if len(res.FundFeeSubList[0].FundFeeSubDtlList) != 1 || res.FundFeeSubList[0].FundFeeSubDtlList[0].SubsFeeRate != 1.5 {
		t.Errorf("Sub / SubDtl nesting failed")
	}

	if len(res.FundFeeCdscList) != 1 {
		t.Fatalf("Expected 1 CDSC, got %d", len(res.FundFeeCdscList))
	}
	if res.FundFeeCdscList[0].MatureYear != "2026" {
		t.Errorf("CDSC mapping failed")
	}

	if len(res.FundFeeBackList) != 1 {
		t.Fatalf("Expected 1 BACK, got %d", len(res.FundFeeBackList))
	}
	if res.FundFeeBackList[0].HoldPeriodYear != "3" {
		t.Errorf("BACK mapping failed")
	}
}
