package fndm009_test

import (
	"context"
	"testing"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
	"go-transfer-agent/services/fnd/fndm009"
	"go-transfer-agent/services/fnd/testutil"
)

func TestHandler_TAFNDFavDisc(t *testing.T) {
	database := testutil.SetupTestDB(t)
	logger := testutil.NewTestLogger()
	svc := fndm009.NewService(database, logger)
	handler := fndm009.NewHandler(logger, svc)

	ctx := context.Background()
	req := &fndv1.TAFNDFavDiscRequest{
		SysCoId:   "UNKNOWN",
		CusIdCode: "UNKNOWN",
	}

	res, err := handler.TAFNDFavDisc(ctx, req)
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
