package fndm004_test

import (
	"context"
	"testing"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
	"go-transfer-agent/services/fnd/fndm004"
	"go-transfer-agent/services/fnd/testutil"
)

func TestHandler_TAFNDFundAgent(t *testing.T) {
	database := testutil.SetupTestDB(t)
	logger := testutil.NewTestLogger()
	svc := fndm004.NewService(database, logger)
	handler := fndm004.NewHandler(logger, svc)

	ctx := context.Background()
	req := &fndv1.TAFNDFundAgentRequest{
		SysCoId:     "UNKNOWN",
		PrtFundCode: "UNKNOWN",
	}

	res, err := handler.TAFNDFundAgent(ctx, req)
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
