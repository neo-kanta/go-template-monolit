package fndm013_test

import (
	"context"
	"testing"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
	"go-transfer-agent/services/fnd/fndm013"
	"go-transfer-agent/services/fnd/testutil"
)

func TestHandler_TAFNDTMFundFeeRdm(t *testing.T) {
	database := testutil.SetupTestDB(t)
	logger := testutil.NewTestLogger()
	svc := fndm013.NewService(database, logger)
	handler := fndm013.NewHandler(logger, svc)

	ctx := context.Background()
	req := &fndv1.TAFNDTMFundFeeRdmRequest{
		SysCoId:     "UNKNOWN",
		PrtFundCode: "UNKNOWN",
	}

	// Just a basic integration ping through the handler to verify wiring
	res, err := handler.TAFNDTMFundFeeRdm(ctx, req)
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
