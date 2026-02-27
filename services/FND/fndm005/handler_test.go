package fndm005_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
	"go-transfer-agent/services/fnd/fndm005"
	"go-transfer-agent/services/fnd/testutil"
)

func TestHandler_TAFNDFundCalDate(t *testing.T) {
	database := testutil.SetupTestDB(t)
	logger := testutil.NewTestLogger()
	svc := fndm005.NewService(database, logger)
	handler := fndm005.NewHandler(logger, svc)

	ctx := context.Background()
	req := &fndv1.TAFNDFundCalDateRequest{
		SysCoId:    "SWSTD",
		CalYear:    "2026",
		FndCalType: "A001",
		FundCry:    "THB",
	}

	res, err := handler.TAFNDFundCalDate(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Empty(t, res.ResultList, "Expected 0 results for unknown ID")
}
