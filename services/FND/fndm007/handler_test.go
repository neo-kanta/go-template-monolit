package fndm007_test

import (
	"context"
	"testing"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
	"go-transfer-agent/services/fnd/fndm007"
	"go-transfer-agent/services/fnd/testutil"
)

func TestHandler_TAFNDCustGroup(t *testing.T) {
	database := testutil.SetupTestDB(t)
	logger := testutil.NewTestLogger()
	svc := fndm007.NewService(database, logger)
	handler := fndm007.NewHandler(logger, svc)

	ctx := context.Background()
	req := &fndv1.TAFNDCustGroupRequest{
		SysCoId:     "UNKNOWN",
		CustGrpCode: "UNKNOWN",
	}

	res, err := handler.TAFNDCustGroup(ctx, req)
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
