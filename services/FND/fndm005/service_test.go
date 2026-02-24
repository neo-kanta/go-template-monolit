package fndm005_test

import (
	"context"
	"testing"
	"time"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
	"go-transfer-agent/services/fnd/fndm005"
	"go-transfer-agent/services/fnd/fndm005/db"
	"go-transfer-agent/services/fnd/testutil"
)

func TestService_TAFNDFundCalDate_Integration(t *testing.T) {
	database := testutil.SetupTestDB(t)
	logger := testutil.NewTestLogger()

	// 1. AutoMigrate schemas
	err := database.AutoMigrate(
		&db.TAFNDFundCal{},
		&db.TAFNDFundCalMemo{},
		&db.TAFNDFundCalDtl{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate FNDM005 schemas: %v", err)
	}

	// 2. Clean up
	database.Exec(`DELETE FROM "TA_STD_TH"."TA_FND_FundCalDtl"`)
	database.Exec(`DELETE FROM "TA_STD_TH"."TA_FND_FundCalMemo"`)
	database.Exec(`DELETE FROM "TA_STD_TH"."TA_FND_FundCal"`)

	// 3. Seed data
	master := db.TAFNDFundCal{
		SysCoID:    "C01",
		CalYear:    "2026",
		FNDCalType: "HLDY",
		FundCry:    "THB",
	}
	if err := database.Create(&master).Error; err != nil {
		t.Fatalf("Failed to seed master: %v", err)
	}

	calDate, _ := time.Parse("2006-01-02", "2026-12-25")
	memos := []db.TAFNDFundCalMemo{
		{SysCoID: "C01", CalYear: "2026", FNDCalType: "HLDY", FundCry: "THB", CalDate: calDate, Remark: "Christmas"},
	}
	if err := database.Create(&memos).Error; err != nil {
		t.Fatalf("Failed to seed memos: %v", err)
	}

	dtls := []db.TAFNDFundCalDtl{
		{SysCoID: "C01", CalYear: "2026", FNDCalType: "HLDY", FundCry: "THB", PrtFundCode: "F1", CalDate: calDate},
		{SysCoID: "C01", CalYear: "2026", FNDCalType: "HLDY", FundCry: "THB", PrtFundCode: "F2", CalDate: calDate},
	}
	if err := database.Create(&dtls).Error; err != nil {
		t.Fatalf("Failed to seed dtls: %v", err)
	}

	// 4. Test Service
	svc := fndm005.NewService(database, logger)
	ctx := context.Background()

	req := &fndv1.TAFNDFundCalDateRequest{
		SysCoId:    "C01",
		CalYear:    "2026",
		FndCalType: "HLDY",
		FundCry:    "THB",
	}

	res, err := svc.TAFNDFundCalDate(ctx, req)
	if err != nil {
		t.Fatalf("Service error: %v", err)
	}

	// 5. Assertions
	if len(res.ResultList) != 1 {
		t.Fatalf("Expected 1 master result, got %d", len(res.ResultList))
	}
	if res.ResultList[0].CalYear != "2026" {
		t.Errorf("Expected CalYear '2026', got '%v'", res.ResultList[0].CalYear)
	}

	if len(res.FundClosedMemoList) != 1 {
		t.Fatalf("Expected 1 memo, got %d", len(res.FundClosedMemoList))
	}
	if res.FundClosedMemoList[0].Remark != "Christmas" {
		t.Errorf("Expected Remark 'Christmas', got '%v'", res.FundClosedMemoList[0].Remark)
	}

	if len(res.FundClosedDateFundList) != 2 {
		t.Fatalf("Expected 2 fund closes, got %d", len(res.FundClosedDateFundList))
	}
	fMap := map[string]bool{}
	for _, f := range res.FundClosedDateFundList {
		fMap[f.PrtFundCode] = true
	}
	if !fMap["F1"] || !fMap["F2"] {
		t.Errorf("Expected F1 and F2, got: %v", fMap)
	}
}
