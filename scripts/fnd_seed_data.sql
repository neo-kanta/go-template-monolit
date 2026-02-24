-- FND API Test Data Seeding
-- Target Schema: TA_STD_TH

-- Ensure schema exists
CREATE SCHEMA IF NOT EXISTS "TA_STD_TH";

-- ═══════════════════════════════════════════════════════════════════
-- FNDM001: Fund Info
-- ═══════════════════════════════════════════════════════════════════
INSERT INTO "TA_STD_TH"."TA_FND_FundInfo" ("SysCoID", "PrtFundCode", "FundInShName")
VALUES ('C01', 'FND001', 'Global Equity Fund') ON CONFLICT DO NOTHING;

INSERT INTO "TA_STD_TH"."DTA_FND_FundInfo" ("SysCoID", "PrtFundCode", "SitcaFundType")
VALUES ('C01', 'FND001', 'EQ') ON CONFLICT DO NOTHING;

INSERT INTO "TA_STD_TH"."DTA_FND_FundDetail" ("SysCoID", "PrtFundCode", "FundCode")
VALUES ('C01', 'FND001', 'F01') ON CONFLICT DO NOTHING;

INSERT INTO "TA_STD_TH"."DTA_FND_FundCryGroup" ("SysCoID", "PrtFundCode", "IssueCry", "IsRsp")
VALUES ('C01', 'FND001', 'USD', 'Y') ON CONFLICT DO NOTHING;

INSERT INTO "TA_STD_TH"."DTA_FND_FundTxCry" ("SysCoID", "PrtFundCode", "CryID")
VALUES ('C01', 'FND001', 'THB') ON CONFLICT DO NOTHING;

-- ═══════════════════════════════════════════════════════════════════
-- FNDM002: Fund Fee
-- ═══════════════════════════════════════════════════════════════════
INSERT INTO "TA_STD_TH"."DTA_FND_FundFee" ("SysCoID", "PrtFundCode")
VALUES ('C01', 'FND001') ON CONFLICT DO NOTHING;

INSERT INTO "TA_STD_TH"."DTA_FND_FundFeeSub" ("SysCoID", "PrtFundCode", "FundCode", "CryID")
VALUES ('C01', 'FND001', 'F01', 'THB') ON CONFLICT DO NOTHING;

INSERT INTO "TA_STD_TH"."DTA_FND_FundFeeSubDtl" ("SysCoID", "PrtFundCode", "FundCode", "CryID", "RangeAmtAbove", "SubsFeeRate")
VALUES ('C01', 'FND001', 'F01', 'THB', 10000.00, 1.50) ON CONFLICT DO NOTHING;

-- ═══════════════════════════════════════════════════════════════════
-- FNDM003: Switch Fund
-- ═══════════════════════════════════════════════════════════════════
INSERT INTO "TA_STD_TH"."DTA_FND_Switch" ("SysCoID", "PrtFundCode", "IsSwitchIn")
VALUES ('C01', 'FND001', 'Y') ON CONFLICT DO NOTHING;

INSERT INTO "TA_STD_TH"."DTA_FND_SwitchFund" ("SysCoID", "PrtFundCode", "SwOPrtFundCode", "SwDateType")
VALUES ('C01', 'FND001', 'FND002', 'T+1') ON CONFLICT DO NOTHING;

-- ═══════════════════════════════════════════════════════════════════
-- FNDM004: Fund Agent
-- ═══════════════════════════════════════════════════════════════════
INSERT INTO "TA_STD_TH"."DTA_FND_FundAgent" ("SysCoID", "PrtFundCode")
VALUES ('C01', 'FND001') ON CONFLICT DO NOTHING;

INSERT INTO "TA_STD_TH"."DTA_FND_FundAgentDtl" ("SysCoID", "PrtFundCode", "AgentType", "AgentCode", "IsMAgentOPType", "FundCrySet", "TermDate")
VALUES ('C01', 'FND001', 'AGN', 'A001', true, 'THB,USD', '2030-12-31') ON CONFLICT DO NOTHING;

-- ═══════════════════════════════════════════════════════════════════
-- FNDM005: Fund Calendar
-- ═══════════════════════════════════════════════════════════════════
INSERT INTO "TA_STD_TH"."TA_FND_FundCal" ("SysCoID", "CalYear", "FNDCalType", "FundCry")
VALUES ('C01', '2026', 'HOL', 'THB') ON CONFLICT DO NOTHING;

INSERT INTO "TA_STD_TH"."TA_FND_FundCalMemo" ("SysCoID", "CalYear", "FNDCalType", "FundCry", "CalDate", "Remark")
VALUES ('C01', '2026', 'HOL', 'THB', '2026-04-13', 'Songkran Festival') ON CONFLICT DO NOTHING;

-- ═══════════════════════════════════════════════════════════════════
-- FNDM008: Pause Txn
-- ═══════════════════════════════════════════════════════════════════
INSERT INTO "TA_STD_TH"."DTA_FND_PauseTxn" ("SysCoID", "PrtFundCode", "PTxnBegDate", "PTxnEndDate", "Remark")
VALUES ('C01', 'FND001', '2026-01-01', '2026-12-31', 'System Maintenance') ON CONFLICT DO NOTHING;

-- ═══════════════════════════════════════════════════════════════════
-- FNDM009: Fav Disc
-- ═══════════════════════════════════════════════════════════════════
INSERT INTO "TA_STD_TH"."TA_FND_FavDisc" ("SysCoID", "FavGrpCode", "FavDiscName")
VALUES ('C01', 'VIP', 'VIP Discount') ON CONFLICT DO NOTHING;

INSERT INTO "TA_STD_TH"."TA_FND_FavDiscFund" ("SysCoID", "FavGrpCode", "PrtFundCode")
VALUES ('C01', 'VIP', 'FND001') ON CONFLICT DO NOTHING;

INSERT INTO "TA_STD_TH"."TA_FND_FavDiscTypeDtl" ("SysCoID", "FavGrpCode", "PrtFundCode", "CusIDCode", "DiscItem", "DiscRate")
VALUES ('C01', 'VIP', 'FND001', 'C001', 'SUB', 0.5) ON CONFLICT DO NOTHING;

-- ═══════════════════════════════════════════════════════════════════
-- FNDM010: IShare Fund Fee
-- ═══════════════════════════════════════════════════════════════════
INSERT INTO "TA_STD_TH"."DTA_FND_IShareFundFee" ("SysCoID", "PrtFundCode")
VALUES ('C01', 'FND001') ON CONFLICT DO NOTHING;

INSERT INTO "TA_STD_TH"."DTA_FND_IShareFundFeeSub" ("SysCoID", "PrtFundCode", "FundCode", "IShareType", "SubsFeeRate")
VALUES ('C01', 'FND001', 'F01', 'NOR', 1.25) ON CONFLICT DO NOTHING;
