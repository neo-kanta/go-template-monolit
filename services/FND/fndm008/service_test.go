package fndm008_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
	"go-transfer-agent/services/fnd/fndm008"
	"go-transfer-agent/services/fnd/fndm008/db"
	"go-transfer-agent/services/fnd/testutil"
)

func TestService_TAFNDPauseTxn_Integration(t *testing.T) {
	database := testutil.SetupTestDB(t)
	logger := testutil.NewTestLogger()

	// 1. AutoMigrate schemas
	err := database.AutoMigrate(
		&db.DTAFNDPauseTxn{},
		&db.DTAFNDPauseTxnCry{},
		&db.DTAFNDPauseTxnDtl{},
	)
	require.NoError(t, err)

	// 2. Clean up
	database.Exec(`DELETE FROM "TA_STD_TH"."DTA_FND_PauseTxnDtl"`)
	database.Exec(`DELETE FROM "TA_STD_TH"."DTA_FND_PauseTxnCry"`)
	database.Exec(`DELETE FROM "TA_STD_TH"."DTA_FND_PauseTxn"`)

	// 3. Seed data
	begDate, _ := time.Parse("2006-01-02", "2026-01-01")
	master := db.DTAFNDPauseTxn{
		SysCoID:     "C01",
		PrtFundCode: "F1",
		PTxnBegDate: begDate,
		
	}
	if err := database.Create(&master).Error; err != nil {
		t.Fatalf("Failed to seed master: %v", err)
	}

	crys := []db.DTAFNDPauseTxnCry{
		{SysCoID: "C01", PrtFundCode: "F1", PTxnBegDate: begDate, CryID: "THB"},
		{SysCoID: "C01", PrtFundCode: "F1", PTxnBegDate: begDate, CryID: "USD"},
	}
	if err := database.Create(&crys).Error; err != nil {
		t.Fatalf("Failed to seed crys: %v", err)
	}

	dtls := []db.DTAFNDPauseTxnDtl{
		{SysCoID: "C01", PrtFundCode: "F1", PTxnBegDate: begDate, PauseTxnType: "SUB"},
		{SysCoID: "C01", PrtFundCode: "F1", PTxnBegDate: begDate, PauseTxnType: "RDM"},
	}
	if err := database.Create(&dtls).Error; err != nil {
		t.Fatalf("Failed to seed dtls: %v", err)
	}

	// 4. Test Service
	svc := fndm008.NewService(database, logger)
	ctx := context.Background()

	req := &fndv1.TAFNDPauseTxnRequest{
		SysCoId:     "C01",
		PrtFundCode: "F1",
	}

	res, err := svc.TAFNDPauseTxn(ctx, req)
	require.NoError(t, err)

	// 5. Assertions
	if len(res.ResultList) != 1 {
		t.Fatalf("Expected 1 master result, got %d", len(res.ResultList))
	}
	if res.ResultList[0].Remark != "Test Pause" {
		t.Errorf("Expected Remark 'Test Pause', got '%v'", res.ResultList[0].Remark)
	}

	if len(res.PauseTxnCryList) != 2 {
		t.Fatalf("Expected 2 crys, got %d", len(res.PauseTxnCryList))
	}
	cMap := map[string]bool{}
	for _, c := range res.PauseTxnCryList {
		cMap[c.CryId] = true
	}
	if !cMap["THB"] || !cMap["USD"] {
		t.Errorf("Expected THB and USD crys, got: %v", cMap)
	}

	if len(res.PauseTxnTypeList) != 2 {
		t.Fatalf("Expected 2 types, got %d", len(res.PauseTxnTypeList))
	}
	tMap := map[string]bool{}
	for _, tp := range res.PauseTxnTypeList {
		tMap[tp.PauseTxnType] = true
	}
	if !tMap["SUB"] || !tMap["RDM"] {
		t.Errorf("Expected SUB and RDM types, got: %v", tMap)
	}
}
