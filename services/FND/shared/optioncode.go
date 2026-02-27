// Package shared – Option code constants ported from C# OptionCode.cs.
// Only the codes actively used in FNDM001-013 services are included.
// DB-based DIC_Option lookups are handled separately via lookup.go.
package shared

// ═══════════════════════════════════════════════════════════════════
// 000000 - Yes / No
// ═══════════════════════════════════════════════════════════════════

const (
	MapYN_Yes = "Y"
	MapYN_No  = "N"
)

// ═══════════════════════════════════════════════════════════════════
// 000001 - Onshore / Offshore Fund Identifier
// ═══════════════════════════════════════════════════════════════════

const (
	ShoreID_DTA = "DTA" // Onshore Fund
	ShoreID_OTA = "OTA" // Offshore Fund
)

// ═══════════════════════════════════════════════════════════════════
// 000006 - Onshore Fund Offering Type
// ═══════════════════════════════════════════════════════════════════

const (
	OfferType_Public  = "Y" // Public Offering
	OfferType_Private = "N" // Private Placement
)

// ═══════════════════════════════════════════════════════════════════
// 000007 - Fund Status
// ═══════════════════════════════════════════════════════════════════

const (
	FundStatus_Active          = "0" // Active / Normal
	FundStatus_Termination     = "1" // Terminated (used by offshore funds)
	FundStatus_Liquidation     = "2" // Liquidated
	FundStatus_Merger          = "3" // Merged
	FundStatus_LiquidationPrep = "8" // Preparing for Liquidation
	FundStatus_MergerPrep      = "9" // Preparing for Merger
)

// ═══════════════════════════════════════════════════════════════════
// 000009 - PG Fund Duration
// ═══════════════════════════════════════════════════════════════════

const (
	PGFStay_Fixed = "F" // Fixed type
)

// ═══════════════════════════════════════════════════════════════════
// 000010 - Target Maturity Kind
// ═══════════════════════════════════════════════════════════════════

const (
	TMFKind_TM  = "1" // Target Maturity
	TMFKind_TMC = "2" // Target Maturity with Duration
	TMFKind_LM  = "3" // Step-up Maturity
	TMFKind_LFM = "4" // Lock-in / Flexible Maturity
)

// ═══════════════════════════════════════════════════════════════════
// 000014 - Remittance/Postal Fee Payer
// ═══════════════════════════════════════════════════════════════════

const (
	RPFeeObj_Beneficiary = "BF" // Paid by Beneficiary
	RPFeeObj_Company     = "C"  // Paid by Company
	RPFeeObj_Bank        = "B"  // Paid by Custodian Bank
)

// ═══════════════════════════════════════════════════════════════════
// 000018 - Short-Term Range Type
// ═══════════════════════════════════════════════════════════════════

const (
	ShortRangeType_Day   = "D"
	ShortRangeType_Month = "M"
)

// ═══════════════════════════════════════════════════════════════════
// 000019 - Short-Term Calculation Calendar Type
// ═══════════════════════════════════════════════════════════════════

const (
	ShortDateType_Calendar = "D"  // Calendar Day
	ShortDateType_FundNAV  = "FB" // Fund NAV Calendar
	ShortDateType_Company  = "CB" // Company Calendar
	ShortDateType_Bank     = "BB" // Bank Calendar
)

// ═══════════════════════════════════════════════════════════════════
// 000025 - Subscription Fee Type
// ═══════════════════════════════════════════════════════════════════

const (
	SubsFeeType_Front    = "B" // Front-end
	SubsFeeType_Deferred = "C" // Deferred
	SubsFeeType_None     = "N" // None
	SubsFeeType_Back     = "P" // Back-end
)

// ═══════════════════════════════════════════════════════════════════
// 000027 - Discount Fee Fund Option
// ═══════════════════════════════════════════════════════════════════

const (
	DiscFundType_AllPrtFund = "0" // All Parent Funds
	DiscFundType_PrtFund    = "1" // Specific Parent Fund
	DiscFundType_Fund       = "2" // Specific Fund
)

// ═══════════════════════════════════════════════════════════════════
// 000037 - Redemption Payment Calendar Identifier
// ═══════════════════════════════════════════════════════════════════

const (
	RdmPayDateType_FundNAV = "FB"  // Fund NAV Calendar
	RdmPayDateType_Bank    = "BB"  // Bank Calendar
	RdmPayDateType_Custom  = "SPB" // Custom Payment Calendar
)

// ═══════════════════════════════════════════════════════════════════
// 000038 - Validity Code
// ═══════════════════════════════════════════════════════════════════

const (
	MapValid_Valid   = "Y"
	MapValid_Invalid = "N"
)

// ═══════════════════════════════════════════════════════════════════
// 000040 - Feature Open / Closed Status
// ═══════════════════════════════════════════════════════════════════

const (
	FeatureStatus_Open   = "Y"
	FeatureStatus_Closed = "N"
)

// ═══════════════════════════════════════════════════════════════════
// 000043 - Basic Calendar Type
// ═══════════════════════════════════════════════════════════════════

const (
	BasicCalType_Company = "CB" // Company Calendar
	BasicCalType_Bank    = "BB" // Bank Calendar
	BasicCalType_Stock   = "SB" // Stock Market Calendar
)

// ═══════════════════════════════════════════════════════════════════
// 000044 - Fund Calendar Type
// ═══════════════════════════════════════════════════════════════════

const (
	FNDCalType_FundNAV    = "FB"  // Fund NAV Calendar
	FNDCalType_CustomSubs = "SSB" // Custom Subscription Calendar
	FNDCalType_CustomRdm  = "SRB" // Custom Redemption Calendar
	FNDCalType_CustomPay  = "SPB" // Custom Payment Calendar
)

// ═══════════════════════════════════════════════════════════════════
// 000047 - Discount Type
// ═══════════════════════════════════════════════════════════════════

const (
	DiscType_FixedRate    = "F" // Fixed Rate (%)
	DiscType_AnnDiscount  = "D" // Announced Rate Discount (%)
	DiscType_ContDiscount = "C" // Contract Rate Discount (%)
	DiscType_AmountRange  = "M" // By Amount Range
)

// ═══════════════════════════════════════════════════════════════════
// 000098 - Fund Pause Transaction Type
// ═══════════════════════════════════════════════════════════════════

const (
	PauseTxnType_Subs     = "1"  // Single Subscription
	PauseTxnType_Rdm      = "2"  // Redemption (includes switch-out)
	PauseTxnType_Switch   = "3"  // Switch-in Subscription
	PauseTxnType_RSP      = "4"  // RSP Submission and Modification
	PauseTxnType_RSPDebit = "5"  // RSP Deduction
	PauseTxnType_ETFSubs  = "E1" // ETF Cash Subscription
	PauseTxnType_ETFRdm   = "E4" // ETF Cash Redemption
)

// ═══════════════════════════════════════════════════════════════════
// 000110 - Discount Fee Item
// ═══════════════════════════════════════════════════════════════════

const (
	DiscItem_Subs   = "A" // Single Subscription
	DiscItem_RSP    = "R" // RSP (Regular Savings Plan)
	DiscItem_Switch = "S" // Switch Subscription
)

// ═══════════════════════════════════════════════════════════════════
// 000187 - Sales Contract Agent Operation Type
// ═══════════════════════════════════════════════════════════════════

const (
	AgentOPType_Both     = "0" // Own Name + On Behalf
	AgentOPType_OwnName  = "1" // Own Name
	AgentOPType_OnBehalf = "2" // On Behalf
)

// ═══════════════════════════════════════════════════════════════════
// Other - Fee Charge Type (Composite Groups)
// ═══════════════════════════════════════════════════════════════════

const (
	FeeChargeType_Front    = "B"  // Front-end
	FeeChargeType_Back     = "P"  // Back-end
	FeeChargeType_Deferred = "C"  // Deferred
	FeeChargeType_None     = "N"  // None
	FeeChargeType_FrontNil = "BN" // Front-end Type (B+N)
	FeeChargeType_BackDef  = "PC" // Back-end Type (P+C)
)

// ═══════════════════════════════════════════════════════════════════
// Other - Fund Kind
// ═══════════════════════════════════════════════════════════════════

const (
	FundKind_TA  = "TA"  // TA Fund
	FundKind_ETF = "ETF" // ETF Fund
	FundKind_DIM = "DIM" // Discretionary Fund
)

// ═══════════════════════════════════════════════════════════════════
// Other - Transaction Type
// ═══════════════════════════════════════════════════════════════════

const (
	TxType_Subs = "Subs" // Subscription
	TxType_Rdm  = "Rdm"  // Redemption
	TxType_All  = "All"  // All
)

// ═══════════════════════════════════════════════════════════════════
// 000108 - Customer ID Type
// ═══════════════════════════════════════════════════════════════════

const (
	CusIDType_Discount  = "D" // Discount Status
	CusIDType_Marketing = "M" // Marketing Status
)

// ═══════════════════════════════════════════════════════════════════
// 000109 - Data Entry Method
// ═══════════════════════════════════════════════════════════════════

const (
	DataSrc_File   = "C" // File Import
	DataSrc_Manual = "M" // Manual Entry
)

// ═══════════════════════════════════════════════════════════════════
// 000143 - Agent Receive Transaction Type
// ═══════════════════════════════════════════════════════════════════

const (
	RcvTxnType_All  = "All"  // All
	RcvTxnType_Subs = "Subs" // Subscription
	RcvTxnType_Rdm  = "Rdm"  // Redemption
)

// ═══════════════════════════════════════════════════════════════════
// 000148 - Currency Classification
// ═══════════════════════════════════════════════════════════════════

const (
	CurrencySort_Central = "CCY" // Center Currency (THB)
	CurrencySort_Foreign = "FCY" // Foreign Currency
)

// ═══════════════════════════════════════════════════════════════════
// 000189 - Debit Institution Type
// ═══════════════════════════════════════════════════════════════════

const (
	DebitType_InvestmentTrust = "1" // Investment Trust
	DebitType_Custodian       = "2" // Custodian
	DebitType_Both            = "3" // Investment Trust + Custodian
)

// ═══════════════════════════════════════════════════════════════════
// 900035 - System Mapping Source
// ═══════════════════════════════════════════════════════════════════

const (
	MapType_Option = "1" // Option (DIC_Option)
	MapType_Code   = "2" // Code
)
