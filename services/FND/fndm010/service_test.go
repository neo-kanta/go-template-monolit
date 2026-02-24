package fndm010_test

import (
	"context"
	"testing"

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
	if err != nil {
		t.Fatalf("Failed to migrate FNDM010 schemas: %v", err)
	}

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
		RangeAmtAbove: 1000,
		SubsFeeRate:   1.0,
	}
	sub2 := db.DTAFNDIShareFundFeeSub{
		SysCoID:       "C01",
		FundCode:      "F1",
		RangeAmtAbove: 5000,
		SubsFeeRate:   0.5,
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
	if err != nil {
		t.Fatalf("Service error: %v", err)
	}

	// 5. Assertions
	if len(res.ResultList) != 1 {
		t.Fatalf("Expected 1 master result, got %d", len(res.ResultList))
	}

	resultMaster := res.ResultList[0]
	if resultMaster.SysCoId != "C01" {
		t.Errorf("Expected SysCoID 'C01', got '%v'", resultMaster.SysCoId)
	}
	if resultMaster.FundCode != "F1" {
		t.Errorf("Expected FundCode 'F1', got '%v'", resultMaster.FundCode)
	}

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
