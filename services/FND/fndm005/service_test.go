package fndm005_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
	require.NoError(t, err)

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
		{SysCoID: "C01", CalYear: "2026", FNDCalType: "HLDY", FundCry: "THB", CalDate: calDate},
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
	require.NoError(t, err)

	// 5. Assertions
	require.Len(t, res.ResultList, 1, "Expected 1 master result")
	assert.Equal(t, "2026", res.ResultList[0].CalYear, "Expected CalYear '2026'")

	require.Len(t, res.FundClosedMemoList, 1, "Expected 1 memo")
	assert.Equal(t, "Christmas", res.FundClosedMemoList[0].Remark, "Expected Remark 'Christmas'")

	require.Len(t, res.FundClosedDateFundList, 2, "Expected 2 fund closes")
	fMap := map[string]bool{}
	for _, f := range res.FundClosedDateFundList {
		fMap[f.PrtFundCode] = true
	}
	assert.True(t, fMap["F1"] && fMap["F2"], "Expected F1 and F2")
}
