package fndm013_test

import (
	"context"
	"testing"
	"time"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
	"go-transfer-agent/services/fnd/fndm013"
	"go-transfer-agent/services/fnd/fndm013/db"
	"go-transfer-agent/services/fnd/testutil"
)

func TestService_TAFNDTMFundFeeRdm_Integration(t *testing.T) {
	database := testutil.SetupTestDB(t)
	logger := testutil.NewTestLogger()

	// 1. AutoMigrate schemas for this module
	err := database.AutoMigrate(
		&db.DTAFNDTMFundFeeRdm{},
		&db.DTAFNDTMFundFeeRdmDtl{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate FNDM013 schemas: %v", err)
	}

	// 2. Clean up before test (in case of dirty DB)
	database.Exec(`DELETE FROM "TA_STD_TH"."DTA_FND_TMFundFeeRdmDtl"`)
	database.Exec(`DELETE FROM "TA_STD_TH"."DTA_FND_TMFundFeeRdm"`)

	// 3. Seed data
	now := time.Now().UTC()
	begDate1 := now.AddDate(-1, 0, 0)
	endDate1 := now.AddDate(1, 0, 0)

	master := db.DTAFNDTMFundFeeRdm{
		SysCoID:        "C01",
		PrtFundCode:    "F1",
		RdmCalcBegDate: begDate1,
	}
	if err := database.Create(&master).Error; err != nil {
		t.Fatalf("Failed to seed master: %v", err)
	}

	dtl1 := db.DTAFNDTMFundFeeRdmDtl{
		SysCoID:        "C01",
		PrtFundCode:    "F1",
		RdmCalcBegDate: begDate1,
		RdmCalcEndDate: endDate1,
		FeeRate:        1.25,
	}
	dtl2 := db.DTAFNDTMFundFeeRdmDtl{
		SysCoID:        "C01",
		PrtFundCode:    "F1",
		RdmCalcBegDate: begDate1,
		RdmCalcEndDate: now.AddDate(2, 0, 0),
		FeeRate:        1.50,
	}
	if err := database.Create(&[]db.DTAFNDTMFundFeeRdmDtl{dtl1, dtl2}).Error; err != nil {
		t.Fatalf("Failed to seed dtls: %v", err)
	}

	// 4. Test Service
	svc := fndm013.NewService(database, logger)
	ctx := context.Background()

	req := &fndv1.TAFNDTMFundFeeRdmRequest{
		SysCoId:     "C01",
		PrtFundCode: "F1",
	}

	res, err := svc.TAFNDTMFundFeeRdm(ctx, req)
	if err != nil {
		t.Fatalf("Service error: %v", err)
	}

	// 5. Assertions
	if len(res.ResultList) != 1 {
		t.Fatalf("Expected 1 master result, got %d", len(res.ResultList))
	}
	masterRes := res.ResultList[0]
	if masterRes.SysCoId != "C01" {
		t.Errorf("Expected SysCoID 'C01', got '%v'", masterRes.SysCoId)
	}
	if masterRes.PrtFundCode != "F1" {
		t.Errorf("Expected PrtFundCode 'F1', got '%v'", masterRes.PrtFundCode)
	}

	// Check details
	if len(res.FundFeeDtlList) != 2 {
		t.Fatalf("Expected 2 details, got %d", len(res.FundFeeDtlList))
	}

	// Verify fee mapping
	fees := map[float64]bool{}
	for _, d := range res.FundFeeDtlList {
		fees[d.FeeRate] = true
	}
	if !fees[1.25] || !fees[1.50] {
		t.Errorf("Expected fee rates 1.25 and 1.50, got: %v", fees)
	}
}
