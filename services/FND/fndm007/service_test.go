package fndm007_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
	"go-transfer-agent/services/fnd/fndm007"
	"go-transfer-agent/services/fnd/fndm007/db"
	"go-transfer-agent/services/fnd/testutil"
)

func TestService_TAFNDCustGroup_Integration(t *testing.T) {
	database := testutil.SetupTestDB(t)
	logger := testutil.NewTestLogger()

	// 1. AutoMigrate schemas
	err := database.AutoMigrate(
		&db.TAFNDCustGroup{},
		&db.TAFNDCustGroupDtl{},
	)
	require.NoError(t, err)

	// 2. Clean up
	database.Exec(`DELETE FROM "TA_STD_TH"."TA_FND_CustGroupDtl"`)
	database.Exec(`DELETE FROM "TA_STD_TH"."TA_FND_CustGroup"`)

	// 3. Seed data
	master := db.TAFNDCustGroup{
		SysCoID:      "C01",
		CustGrpCode:  "GRP1",
		CustGrpMName: "Group 1 Main",
	}
	if err := database.Create(&master).Error; err != nil {
		t.Fatalf("Failed to seed master: %v", err)
	}

	dtls := []db.TAFNDCustGroupDtl{
		{SysCoID: "C01", CustGrpCode: "GRP1", PrtFundCode: "F1"},
		{SysCoID: "C01", CustGrpCode: "GRP1", PrtFundCode: "F2"},
	}
	if err := database.Create(&dtls).Error; err != nil {
		t.Fatalf("Failed to seed dtls: %v", err)
	}

	// 4. Test Service
	svc := fndm007.NewService(database, logger)
	ctx := context.Background()

	req := &fndv1.TAFNDCustGroupRequest{
		SysCoId:     "C01",
		CustGrpCode: "GRP1",
	}

	res, err := svc.TAFNDCustGroup(ctx, req)
	require.NoError(t, err)

	// 5. Assertions
	if len(res.ResultList) != 1 {
		t.Fatalf("Expected 1 master result, got %d", len(res.ResultList))
	}
	if res.ResultList[0].CustGrpMName != "Group 1 Main" {
		t.Errorf("Expected MName 'Group 1 Main', got '%v'", res.ResultList[0].CustGrpMName)
	}

	if len(res.CustFundGroupList) != 2 {
		t.Fatalf("Expected 2 dtls, got %d", len(res.CustFundGroupList))
	}
	fMap := map[string]bool{}
	for _, d := range res.CustFundGroupList {
		fMap[d.PrtFundCode] = true
	}
	if !fMap["F1"] || !fMap["F2"] {
		t.Errorf("Expected F1 and F2, got: %v", fMap)
	}
}
