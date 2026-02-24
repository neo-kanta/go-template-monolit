package fndm001_test

import (
	"context"
	"log/slog"
	"os"
	"testing"

	fndv1 "go-transfer-agent/common/gen/fnd/v1"
	"go-transfer-agent/services/fnd/fndm001"
)

func setupHandler() *fndm001.Handler {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	svc := fndm001.NewService()
	return fndm001.NewHandler(logger, svc)
}

// ═══════════════════════════════════════════════════════════════════
// 1. AUD / Save Operations (APIFNDM001Post/Put)
// ═══════════════════════════════════════════════════════════════════

func TestHandler_SaveFundInfo_Success(t *testing.T) {
	h := setupHandler()
	ctx := context.Background()

	req := &fndv1.SaveFundInfoRequest{
		Master: &fndv1.TAFNDFundInfo{
			SysCoId:     "C1",
			PrtFundCode: "P1",
		},
	}

	res, err := h.SaveFundInfo(ctx, req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !res.Success {
		t.Errorf("Expected SaveFundInfo to succeed, got false. Msg: %v", res.Message)
	}
}

// ═══════════════════════════════════════════════════════════════════
// 2. Query Operations (APIFNDM001Get)
// ═══════════════════════════════════════════════════════════════════

func TestHandler_QueryFundInfo_APIFNDM001Get(t *testing.T) {
	h := setupHandler()
	ctx := context.Background()

	h.SaveFundInfo(ctx, &fndv1.SaveFundInfoRequest{
		Master: &fndv1.TAFNDFundInfo{SysCoId: "C1", PrtFundCode: "P1", FundInShName: "Fund1"},
	})
	h.SaveFundInfo(ctx, &fndv1.SaveFundInfoRequest{
		Master: &fndv1.TAFNDFundInfo{SysCoId: "C2", PrtFundCode: "P2", FundInShName: "Fund2"},
	})

	res, err := h.QueryFundInfo(ctx, &fndv1.QueryFundInfoRequest{SysCoId: "C1"})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(res.ResultList) != 1 {
		t.Fatalf("Expected 1 result for Company C1, got %d", len(res.ResultList))
	}
	if res.ResultList[0].PrtFundCode != "P1" {
		t.Errorf("Expected PrtFundCode P1, got %v", res.ResultList[0].PrtFundCode)
	}
}

func TestHandler_QueryFundInfoByDataID_APIFNDM001GetMaintain(t *testing.T) {
	h := setupHandler()
	ctx := context.Background()

	h.SaveFundInfo(ctx, &fndv1.SaveFundInfoRequest{
		Master:   &fndv1.TAFNDFundInfo{SysCoId: "C1", PrtFundCode: "P1", FundInShName: "Fund1"},
		FundInfo: &fndv1.DTAFNDFundInfo{SysCoId: "C1", PrtFundCode: "P1", SitcaFundType: "T1"},
	})

	res, err := h.QueryFundInfoByDataID(ctx, &fndv1.QueryFundInfoByDataIDRequest{
		SysCoId:     "C1",
		PrtFundCode: "P1",
	})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if res.Master == nil || res.Master.FundInShName != "Fund1" {
		t.Errorf("Expected Master FundInShName Fund1, got %v", res.Master)
	}
	if res.FundInfo == nil || res.FundInfo.SitcaFundType != "T1" {
		t.Errorf("Expected FundInfo SitcaFundType T1, got %v", res.FundInfo)
	}
}

// ═══════════════════════════════════════════════════════════════════
// 3. Check Functions (KFNDM)
// ═══════════════════════════════════════════════════════════════════

func TestHandler_TACKCSDPrtFundSetlDate_KFNDM00101(t *testing.T) {
	h := setupHandler()
	ctx := context.Background()

	req := &fndv1.TACKCSDPrtFundSetlDateRequest{
		SysCoId:            "C1",
		CsdPrtFundSetlDate: "3000-01-01", // Way into the future, should be valid
	}

	res, err := h.TACKCSDPrtFundSetlDate(ctx, req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if res.ReturnCode != "OK" {
		t.Errorf("Expected OK, got %v", res.ReturnCode)
	}
}

func TestHandler_TACKFundCodeExist_KFNDM00103(t *testing.T) {
	h := setupHandler()
	ctx := context.Background()

	h.SaveFundInfo(ctx, &fndv1.SaveFundInfoRequest{
		Master: &fndv1.TAFNDFundInfo{SysCoId: "C1", PrtFundCode: "P1"},
		FundDetails: []*fndv1.DTAFNDFundDetail{
			{SysCoId: "C1", PrtFundCode: "P1", FundCode: "EXIST"},
		},
	})

	res, err := h.TACKFundCodeExist(ctx, &fndv1.TACKFundCodeExistRequest{
		SysCoId:  "C1",
		FundCode: "EXIST",
	})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if res.ReturnCode != "TA00030" {
		t.Errorf("Expected TA00030 for existing fund, got %v", res.ReturnCode)
	}
}

// ═══════════════════════════════════════════════════════════════════
// 4. Search/Info Functions (GFNDM)
// ═══════════════════════════════════════════════════════════════════

func TestHandler_TAGetDPrtFundIssueCry_GFNDM00101(t *testing.T) {
	h := setupHandler()
	ctx := context.Background()

	h.SaveFundInfo(ctx, &fndv1.SaveFundInfoRequest{
		Master: &fndv1.TAFNDFundInfo{SysCoId: "C1", PrtFundCode: "P1"},
		CryGroups: []*fndv1.DTAFNDFundCryGroup{
			{SysCoId: "C1", PrtFundCode: "P1", IssueCry: "USD", IsRsp: "Y"},
		},
	})

	res, err := h.TAGetDPrtFundIssueCry(ctx, &fndv1.TAGetDPrtFundIssueCryRequest{
		SysCoId:     "C1",
		PrtFundCode: "P1",
		QueryType:   "ALL",
	})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(res.Master) != 1 {
		t.Fatalf("Expected 1 currency group, got %d", len(res.Master))
	}
	if res.Master[0].CryId != "USD" {
		t.Errorf("Expected USD, got %v", res.Master[0].CryId)
	}
}

func TestHandler_TAGetDTxCry_GFNDM00104(t *testing.T) {
	h := setupHandler()
	ctx := context.Background()

	h.SaveFundInfo(ctx, &fndv1.SaveFundInfoRequest{
		Master: &fndv1.TAFNDFundInfo{SysCoId: "C1", PrtFundCode: "P1"},
		FundTxCrys: []*fndv1.DTAFNDFundTxCry{
			{SysCoId: "C1", PrtFundCode: "P1", CryId: "EUR"},
		},
	})

	res, err := h.TAGetDTxCry(ctx, &fndv1.TAGetDTxCryRequest{
		SysCoId:     "C1",
		PrtFundCode: "P1",
		QueryType:   "ALL",
	})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(res.Master) != 1 {
		t.Fatalf("Expected 1 transaction currency, got %d", len(res.Master))
	}
	if res.Master[0].CryId != "EUR" {
		t.Errorf("Expected EUR, got %v", res.Master[0].CryId)
	}
}
