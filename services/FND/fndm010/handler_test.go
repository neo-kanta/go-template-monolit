package fndm010_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
	"go-transfer-agent/services/fnd/fndm010"
	"go-transfer-agent/services/fnd/testutil"
)

func TestHandler_TAFNDIShareFundFee(t *testing.T) {
	database := testutil.SetupTestDB(t)
	logger := testutil.NewTestLogger()
	svc := fndm010.NewService(database, logger)
	handler := fndm010.NewHandler(logger, svc)

	ctx := context.Background()
	req := &fndv1.TAFNDIShareFundFeeRequest{
		SysCoId:  "SWSTD",
		FundCode: "A001",
	}

	res, err := handler.TAFNDIShareFundFee(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, res)
	if len(res.ResultList) != 0 {
		t.Errorf("Expected 0 results for unknown ID, got %d", len(res.ResultList))
	}
}
