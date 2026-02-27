package fndm002_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
	"go-transfer-agent/services/fnd/fndm002"
	"go-transfer-agent/services/fnd/testutil"
)

func TestHandler_TAFNDFundFee(t *testing.T) {
	database := testutil.SetupTestDB(t)
	logger := testutil.NewTestLogger()
	svc := fndm002.NewService(database, logger)
	handler := fndm002.NewHandler(logger, svc)

	ctx := context.Background()
	req := &fndv1.TAFNDFundFeeRequest{
		SysCoId:     "UNKNOWN",
		PrtFundCode: "UNKNOWN",
	}

	res, err := handler.TAFNDFundFee(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, res)
	if len(res.ResultList) != 0 {
		t.Errorf("Expected 0 results for unknown ID, got %d", len(res.ResultList))
	}
}
