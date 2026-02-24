package fndm004_test

import (
	"context"
	"testing"
	"time"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
	"go-transfer-agent/services/fnd/fndm004"
	"go-transfer-agent/services/fnd/fndm004/db"
	"go-transfer-agent/services/fnd/testutil"
)

func TestService_TAFNDFundAgent_Integration(t *testing.T) {
	database := testutil.SetupTestDB(t)
	logger := testutil.NewTestLogger()

	// 1. AutoMigrate schemas
	err := database.AutoMigrate(
		&db.DTAFNDFundAgent{},
		&db.DTAFNDFundAgentDtl{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate FNDM004 schemas: %v", err)
	}

	// 2. Clean up
	database.Exec(`DELETE FROM "TA_STD_TH"."DTA_FND_FundAgentDtl"`)
	database.Exec(`DELETE FROM "TA_STD_TH"."DTA_FND_FundAgent"`)

	// 3. Seed data
	master := db.DTAFNDFundAgent{
		SysCoID:     "C01",
		PrtFundCode: "F1",
	}
	if err := database.Create(&master).Error; err != nil {
		t.Fatalf("Failed to seed master: %v", err)
	}

	termDate, _ := time.Parse("2006-01-02", "2030-12-31")
	dtls := []db.DTAFNDFundAgentDtl{
		{
			SysCoID:        "C01",
			PrtFundCode:    "F1",
			AgentType:      "TYPE1",
			AgentCode:      "AGT1",
			SubsFeePctAG:   0.5,
			SubsFeePctFH:   0.5,
			FundCrySet:     "THB,USD",
			NoSaleShareSet: "S1,S2",
			TermDate:       termDate,
		},
	}
	if err := database.Create(&dtls).Error; err != nil {
		t.Fatalf("Failed to seed dtls: %v", err)
	}

	// 4. Test Service
	svc := fndm004.NewService(database, logger)
	ctx := context.Background()

	req := &fndv1.TAFNDFundAgentRequest{
		SysCoId:     "C01",
		PrtFundCode: "F1",
	}

	res, err := svc.TAFNDFundAgent(ctx, req)
	if err != nil {
		t.Fatalf("Service error: %v", err)
	}

	// 5. Assertions
	if len(res.ResultList) != 1 {
		t.Fatalf("Expected 1 master result, got %d", len(res.ResultList))
	}
	if res.ResultList[0].PrtFundCode != "F1" {
		t.Errorf("Expected PrtFundCode 'F1', got '%v'", res.ResultList[0].PrtFundCode)
	}

	if len(res.FundAgentDtlList) != 1 {
		t.Fatalf("Expected 1 dtl, got %d", len(res.FundAgentDtlList))
	}
	dtl := res.FundAgentDtlList[0]
	if dtl.AgentCode != "AGT1" {
		t.Errorf("Expected AgentCode 'AGT1', got '%v'", dtl.AgentCode)
	}
	if dtl.SubsFeePctAg != 0.5 {
		t.Errorf("Expected SubsFeePctAG 0.5, got '%v'", dtl.SubsFeePctAg)
	}

	// Check split sets
	if len(dtl.FundCrySet) != 2 || dtl.FundCrySet[0] != "THB" || dtl.FundCrySet[1] != "USD" {
		t.Errorf("FundCrySet mapping failed: %v", dtl.FundCrySet)
	}
	if len(dtl.NoSaleShareSet) != 2 || dtl.NoSaleShareSet[0] != "S1" || dtl.NoSaleShareSet[1] != "S2" {
		t.Errorf("NoSaleShareSet mapping failed: %v", dtl.NoSaleShareSet)
	}
}
