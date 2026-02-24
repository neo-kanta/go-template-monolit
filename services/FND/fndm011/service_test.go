package fndm011_test

import (
	"context"
	"testing"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
	"go-transfer-agent/services/fnd/fndm011"
	"go-transfer-agent/services/fnd/fndm011/db"
	"go-transfer-agent/services/fnd/testutil"
)

func TestService_TAFNDIShareFundFeeRdm_Integration(t *testing.T) {
	database := testutil.SetupTestDB(t)
	logger := testutil.NewTestLogger()

	// 1. AutoMigrate schemas
	err := database.AutoMigrate(&db.DTAFNDIShareFundFeeRdm{})
	if err != nil {
		t.Fatalf("Failed to migrate FNDM011 schemas: %v", err)
	}

	// 2. Clean up
	database.Exec(`DELETE FROM "TA_STD_TH"."DTA_FND_IShareFundFeeRdm"`)

	// 3. Seed data
	dtl1 := db.DTAFNDIShareFundFeeRdm{
		SysCoID:        "C01",
		FundCode:       "F1",
		FeeName:        "FeeA",
		RdmRangeType:   "RT1",
		RdmDateType:    "DT1",
		ShouldHoldDays: 10,
		FeeRate:        1.2,
	}
	dtl2 := db.DTAFNDIShareFundFeeRdm{
		SysCoID:        "C01",
		FundCode:       "F1",
		FeeName:        "FeeB",
		RdmRangeType:   "RT2",
		RdmDateType:    "DT2",
		ShouldHoldDays: 20,
		FeeRate:        2.4,
	}

	if err := database.Create(&[]db.DTAFNDIShareFundFeeRdm{dtl1, dtl2}).Error; err != nil {
		t.Fatalf("Failed to seed data: %v", err)
	}

	// 4. Test Service
	svc := fndm011.NewService(database, logger)
	ctx := context.Background()

	req := &fndv1.TAFNDIShareFundFeeRdmRequest{
		SysCoId:  "C01",
		FundCode: "F1",
	}

	res, err := svc.TAFNDIShareFundFeeRdm(ctx, req)
	if err != nil {
		t.Fatalf("Service error: %v", err)
	}

	// 5. Assertions
	if len(res.ResultList) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(res.ResultList))
	}

	fees := map[float64]bool{}
	for _, r := range res.ResultList {
		fees[r.FeeRate] = true
		if r.SysCoId != "C01" {
			t.Errorf("Expected SysCoID 'C01', got '%v'", r.SysCoId)
		}
	}

	if !fees[1.2] || !fees[2.4] {
		t.Errorf("Expected fee rates 1.2 and 2.4, got: %v", fees)
	}
}
