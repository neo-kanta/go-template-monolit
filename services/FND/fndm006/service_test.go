package fndm006_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
	"go-transfer-agent/services/fnd/fndm006"
	"go-transfer-agent/services/fnd/fndm006/db"
	"go-transfer-agent/services/fnd/testutil"
)

func TestService_TAFNDRPFeeChgType_Integration(t *testing.T) {
	database := testutil.SetupTestDB(t)
	logger := testutil.NewTestLogger()

	// 1. AutoMigrate schemas
	err := database.AutoMigrate(
		&db.DTAFNDRPFeeChgType{},
		&db.DTAFNDRPFeeChgTypeDtl{},
	)
	require.NoError(t, err)

	// 2. Clean up
	database.Exec(`DELETE FROM "TA_STD_TH"."DTA_FND_RPFeeChgTypeDtl"`)
	database.Exec(`DELETE FROM "TA_STD_TH"."DTA_FND_RPFeeChgType"`)

	// 3. Seed data
	master := db.DTAFNDRPFeeChgType{
		SysCoID:     "C01",
		PmtTxnType:  "PMT1",
		PrtFundCode: "F1",
	}
	if err := database.Create(&master).Error; err != nil {
		t.Fatalf("Failed to seed master: %v", err)
	}

	dtls := []db.DTAFNDRPFeeChgTypeDtl{
		{SysCoID: "C01", PmtTxnType: "PMT1", PrtFundCode: "F1", CryKind: "THB", RemitFeeObj: "OBJ1"},
		{SysCoID: "C01", PmtTxnType: "PMT1", PrtFundCode: "F1", CryKind: "USD", RemitFeeObj: "OBJ2"},
	}
	if err := database.Create(&dtls).Error; err != nil {
		t.Fatalf("Failed to seed dtls: %v", err)
	}

	// 4. Test Service
	svc := fndm006.NewService(database, logger)
	ctx := context.Background()

	req := &fndv1.TAFNDRPFeeChgTypeRequest{
		SysCoId:    "C01",
		PmtTxnType: "PMT1",
	}

	res, err := svc.TAFNDRPFeeChgType(ctx, req)
	require.NoError(t, err)

	// 5. Assertions
	if len(res.ResultList) != 1 {
		t.Fatalf("Expected 1 master result, got %d", len(res.ResultList))
	}
	if res.ResultList[0].PmtTxnType != "PMT1" {
		t.Errorf("Expected PmtTxnType 'PMT1', got '%v'", res.ResultList[0].PmtTxnType)
	}

	if len(res.ChgTypeList) != 2 {
		t.Fatalf("Expected 2 dtls, got %d", len(res.ChgTypeList))
	}
	cMap := map[string]bool{}
	for _, d := range res.ChgTypeList {
		cMap[d.CryKind] = true
	}
	if !cMap["THB"] || !cMap["USD"] {
		t.Errorf("Expected THB and USD, got: %v", cMap)
	}
}
