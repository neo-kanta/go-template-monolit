package fndm012_test

import (
	"context"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
	"go-transfer-agent/services/fnd/fndm012"
	"go-transfer-agent/services/fnd/fndm012/db"
	"go-transfer-agent/services/fnd/testutil"
)

func TestService_TAFNDPGFundFeeRdm_Integration(t *testing.T) {
	database := testutil.SetupTestDB(t)
	logger := testutil.NewTestLogger()

	// 1. AutoMigrate schemas
	err := database.AutoMigrate(&db.DTAFNDPGFundFeeRdm{})
	require.NoError(t, err)

	// 2. Clean up
	database.Exec(`DELETE FROM "TA_STD_TH"."DTA_FND_PGFundFeeRdm"`)

	// 3. Seed data
	now := time.Now().UTC()
	dtl1 := db.DTAFNDPGFundFeeRdm{
		SysCoID:        "C01",
		PrtFundCode:    "F1",
		RdmCalcBegDate: now.AddDate(-1, 0, 0),
		FeeRate:        decimal.NewFromFloat(0.75).RoundBank(4),
	}
	dtl2 := db.DTAFNDPGFundFeeRdm{
		SysCoID:        "C01",
		PrtFundCode:    "F1",
		RdmCalcBegDate: now.AddDate(1, 0, 0),
		FeeRate:        decimal.NewFromFloat(1.0).RoundBank(4),
	}

	if err := database.Create(&[]db.DTAFNDPGFundFeeRdm{dtl1, dtl2}).Error; err != nil {
		t.Fatalf("Failed to seed data: %v", err)
	}

	// 4. Test Service
	svc := fndm012.NewService(database, logger)
	ctx := context.Background()

	req := &fndv1.TAFNDPGFundFeeRdmRequest{
		SysCoId:     "C01",
		PrtFundCode: "F1",
	}

	res, err := svc.TAFNDPGFundFeeRdm(ctx, req)
	require.NoError(t, err)

	// 5. Assertions
	if len(res.ResultList) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(res.ResultList))
	}

	fees := map[float64]bool{}
	for _, r := range res.ResultList {
		fees[r.FeeRate] = true
		assert.Equal(t, "C01", r.SysCoId)
		assert.Equal(t, "F1", r.PrtFundCode)
	}

	if !fees[0.75] || !fees[1.0] {
		t.Errorf("Expected fee rates 0.75 and 1.0, got: %v", fees)
	}
}
