package fndm009_test

import (
	"context"
	"testing"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
	"go-transfer-agent/services/fnd/fndm009"
	"go-transfer-agent/services/fnd/fndm009/db"
	"go-transfer-agent/services/fnd/testutil"
)

func TestService_TAFNDFavDisc_Integration(t *testing.T) {
	database := testutil.SetupTestDB(t)
	logger := testutil.NewTestLogger()

	// 1. AutoMigrate schemas
	err := database.AutoMigrate(
		&db.TAFNDFavDisc{},
		&db.TAFNDFavDiscFund{},
		&db.TAFNDFavDiscType{},
		&db.TAFNDFavDiscTypeDtl{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate FNDM009 schemas: %v", err)
	}

	// 2. Clean up
	database.Exec(`DELETE FROM "TA_STD_TH"."TA_FND_FavDiscTypeDtl"`)
	database.Exec(`DELETE FROM "TA_STD_TH"."TA_FND_FavDiscType"`)
	database.Exec(`DELETE FROM "TA_STD_TH"."TA_FND_FavDiscFund"`)
	database.Exec(`DELETE FROM "TA_STD_TH"."TA_FND_FavDisc"`)

	// 3. Seed data
	master := db.TAFNDFavDisc{
		SysCoID:      "C01",
		CusIDCode:    "CUST1",
		DiscItemSet:  "ITEM1,ITEM2",
		DiscFundType: "PF",
	}
	if err := database.Create(&master).Error; err != nil {
		t.Fatalf("Failed to seed master: %v", err)
	}

	fund := db.TAFNDFavDiscFund{
		SysCoID:     "C01",
		CusIDCode:   "CUST1",
		PrtFundCode: "F1",
	}
	if err := database.Create(&fund).Error; err != nil {
		t.Fatalf("Failed to seed fund: %v", err)
	}

	types := []db.TAFNDFavDiscType{
		{SysCoID: "C01", CusIDCode: "CUST1", DiscItem: "01", DiscType: "ALLOT"},
		{SysCoID: "C01", CusIDCode: "CUST1", DiscItem: "02", DiscType: "RSP"},
		{SysCoID: "C01", CusIDCode: "CUST1", DiscItem: "03", DiscType: "SW"},
	}
	if err := database.Create(&types).Error; err != nil {
		t.Fatalf("Failed to seed types: %v", err)
	}

	dtls := []db.TAFNDFavDiscTypeDtl{
		{SysCoID: "C01", CusIDCode: "CUST1", DiscItem: "01", TxCry: "THB", RangeFeeRate: 0.5},
		{SysCoID: "C01", CusIDCode: "CUST1", DiscItem: "02", TxCry: "THB", RangeFeeRate: 0.6},
		{SysCoID: "C01", CusIDCode: "CUST1", DiscItem: "03", TxCry: "THB", RangeFeeRate: 0.7},
	}
	if err := database.Create(&dtls).Error; err != nil {
		t.Fatalf("Failed to seed dtls: %v", err)
	}

	// 4. Test Service
	svc := fndm009.NewService(database, logger)
	ctx := context.Background()

	req := &fndv1.TAFNDFavDiscRequest{
		SysCoId:   "C01",
		CusIdCode: "CUST1",
	}

	res, err := svc.TAFNDFavDisc(ctx, req)
	if err != nil {
		t.Fatalf("Service error: %v", err)
	}

	// 5. Assertions
	if len(res.ResultList) != 1 {
		t.Fatalf("Expected 1 master result, got %d", len(res.ResultList))
	}
	if res.ResultList[0].DiscFundType != "PF" {
		t.Errorf("Expected DiscFundType 'PF', got '%v'", res.ResultList[0].DiscFundType)
	}
	if len(res.ResultList[0].DiscItemSet) != 2 {
		t.Fatalf("Expected 2 disc items, got %v", res.ResultList[0].DiscItemSet)
	}

	if len(res.PrtFundList) != 1 {
		t.Fatalf("Expected 1 prt fund, got %d", len(res.PrtFundList))
	}

	if len(res.AllotDiscTypeDtlList) != 1 || res.AllotDiscTypeDtlList[0].RangeFeeRate != 0.5 {
		t.Errorf("Allot disc type mapping failed")
	}
	if len(res.RspDiscTypeDtlList) != 1 || res.RspDiscTypeDtlList[0].RangeFeeRate != 0.6 {
		t.Errorf("Rsp disc type mapping failed")
	}
	if len(res.SwDiscTypeDtlList) != 1 || res.SwDiscTypeDtlList[0].RangeFeeRate != 0.7 {
		t.Errorf("Sw disc type mapping failed")
	}
}
