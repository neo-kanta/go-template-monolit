package fndm011_test

import (
	"context"
	"testing"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
	"go-transfer-agent/services/fnd/fndm011"
	"go-transfer-agent/services/fnd/testutil"
)

func TestHandler_TAFNDIShareFundFeeRdm(t *testing.T) {
	database := testutil.SetupTestDB(t)
	logger := testutil.NewTestLogger()
	svc := fndm011.NewService(database, logger)
	handler := fndm011.NewHandler(logger, svc)

	ctx := context.Background()
	req := &fndv1.TAFNDIShareFundFeeRdmRequest{
		SysCoId:  "UNKNOWN",
		FundCode: "UNKNOWN",
	}

	res, err := handler.TAFNDIShareFundFeeRdm(ctx, req)
	if err != nil {
		t.Fatalf("Handler error: %v", err)
	}
	if res == nil {
		t.Fatal("Expected non-nil response")
	}
	if len(res.ResultList) != 0 {
		t.Errorf("Expected 0 results for unknown ID, got %d", len(res.ResultList))
	}
}
