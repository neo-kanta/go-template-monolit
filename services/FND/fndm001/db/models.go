package db

import (
	"github.com/shopspring/decimal"

	models "go-transfer-agent/common/platform/model"

	"time"
)

// ═══════════════════════════════════════════════════════════════════
// TA_STD_TH Schema Models
// ═══════════════════════════════════════════════════════════════════

// AuditFields contains the standard audit columns present on every table.
type AuditFields struct {
	ValidFrom   time.Time `gorm:"column:ValidFrom;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	ValidTo     time.Time `gorm:"column:ValidTo;type:timestamp;not null;default:'9999-12-31 23:59:59.99'"`
	DataID      string    `gorm:"column:DataID;type:uuid;not null;default:gen_random_uuid()"`
	CreateID    string    `gorm:"column:CreateID;type:varchar(30);not null;default:''"`
	CreateDate  time.Time `gorm:"column:CreateDate;type:timestamp with time zone;not null;default:'1900-01-01 00:00:00+08'"`
	Inspect     int16     `gorm:"column:Inspect;type:smallint;not null;default:1"`
	FlowID      int64     `gorm:"column:FlowID;type:bigint;not null;default:0"`
	InspectID   string    `gorm:"column:InspectID;type:varchar(30);not null;default:''"`
	InspectDate time.Time `gorm:"column:InspectDate;type:timestamp with time zone;not null;default:'1900-01-01 00:00:00+08'"`
	UpdateID    string    `gorm:"column:UpdateID;type:varchar(30);not null;default:''"`
	UpdateDate  time.Time `gorm:"column:UpdateDate;type:timestamp with time zone;not null;default:'1900-01-01 00:00:00+08'"`
	DataFlag    []byte    `gorm:"column:DataFlag;type:bytea"`
	DiffColumns string    `gorm:"column:DiffColumns;type:text;not null;default:''"`
	models.MakerCheckerFields
}

// ═══════════════════════════════════════════════════════════════════
// Master: TA_FND_FundInfo
// ═══════════════════════════════════════════════════════════════════

// TAFNDFundInfo represents the TA_FND_FundInfo table.
type TAFNDFundInfo struct {
	SysCoID        string          `gorm:"column:SysCoID;type:varchar(20);primaryKey;not null"`
	PrtFundCode    string          `gorm:"column:PrtFundCode;type:varchar(10);primaryKey;not null"`
	UniCode        string          `gorm:"column:UniCode;type:varchar(30);not null;default:''"`
	FundInShName   string          `gorm:"column:FundInShName;type:varchar(100);not null;default:''"`
	FundMName      string          `gorm:"column:FundMName;type:varchar(200);not null;default:''"`
	FundShMName    string          `gorm:"column:FundShMName;type:varchar(150);not null;default:''"`
	FundSName      string          `gorm:"column:FundSName;type:varchar(200);not null;default:''"`
	FundShSName    string          `gorm:"column:FundShSName;type:varchar(150);not null;default:''"`
	FundRiskLevel  string          `gorm:"column:FundRiskLevel;type:varchar(6);not null;default:''"`
	IssueBaseCry   string          `gorm:"column:IssueBaseCry;type:varchar(6);not null;default:''"`
	FundSetupDate  time.Time       `gorm:"column:FundSetupDate;type:timestamp with time zone;not null;default:'1900-01-01 00:00:00+08'"`
	GIINNo         string          `gorm:"column:GIINNo;type:varchar(20);not null;default:''"`
	FundStatus     string          `gorm:"column:FundStatus;type:varchar(6);not null;default:''"`
	FundWarningMsg string          `gorm:"column:FundWarningMsg;type:varchar(200);not null;default:''"`
	UnitDotLen     int16           `gorm:"column:UnitDotLen;type:smallint;not null;default:0"`
	FundSubsWay    string          `gorm:"column:FundSubsWay;type:varchar(6);not null;default:''"`
	FundRdmWay     string          `gorm:"column:FundRdmWay;type:varchar(6);not null;default:''"`
	RdmFNavDay     int16           `gorm:"column:RdmFNavDay;type:smallint;not null;default:0"`
	RdmFNavWay     string          `gorm:"column:RdmFNavWay;type:varchar(6);not null;default:''"`
	RdmFNavDRate   decimal.Decimal `gorm:"column:RdmFNavDRate;type:numeric(5,2);not null;default:0"`
	RdmFNavDVal    decimal.Decimal `gorm:"column:RdmFNavDVal;type:numeric(20,6);not null;default:0"`
	IsShortFee     string          `gorm:"column:IsShortFee;type:varchar(6);not null;default:''"`
	IsAntiDilFee   string          `gorm:"column:IsAntiDilFee;type:varchar(1);not null;default:''"`
	IsETF          string          `gorm:"column:IsETF;type:varchar(1);not null;default:''"`
	ETFCode        string          `gorm:"column:ETFCode;type:varchar(10);not null;default:''"`
	IsDIM          string          `gorm:"column:IsDIM;type:varchar(1);not null;default:''"`
	ShoreID        string          `gorm:"column:ShoreID;type:varchar(6);not null;default:''"`
	OFDFHCode      string          `gorm:"column:OFDFHCode;type:varchar(10);not null;default:''"`
	AuditFields
}

func (TAFNDFundInfo) TableName() string { return "TA_STD_TH.TA_FND_FundInfo" }

type TAFNDFundInfoEdit struct{ TAFNDFundInfo }

func (TAFNDFundInfoEdit) TableName() string { return "TA_STD_TH.TA_FND_FundInfo_Edit" }

// ═══════════════════════════════════════════════════════════════════
// Child: DTA_FND_FundInfo (Detail fund information)
// ═══════════════════════════════════════════════════════════════════

type DTAFNDFundInfo struct {
	SysCoID            string          `gorm:"column:SysCoID;type:varchar(20);primaryKey;not null"`
	PrtFundCode        string          `gorm:"column:PrtFundCode;type:varchar(10);primaryKey;not null"`
	SITCAFundType      string          `gorm:"column:SITCAFundType;type:varchar(10);not null;default:''"`
	DFundType          string          `gorm:"column:DFundType;type:varchar(10);not null;default:''"`
	DInvArea           string          `gorm:"column:DInvArea;type:varchar(10);not null;default:''"`
	OfferType          string          `gorm:"column:OfferType;type:varchar(6);not null;default:''"`
	BCFundQuoUpAmt     decimal.Decimal `gorm:"column:BCFundQuoUpAmt;type:numeric;not null;default:0"`
	BCFundQuoMinAmt    decimal.Decimal `gorm:"column:BCFundQuoMinAmt;type:numeric;not null;default:0"`
	FundQuoLmtRate     decimal.Decimal `gorm:"column:FundQuoLmtRate;type:numeric;not null;default:0"`
	TradeWarningAmt    decimal.Decimal `gorm:"column:TradeWarningAmt;type:numeric;not null;default:0"`
	ExchUnitLmt        string          `gorm:"column:ExchUnitLmt;type:varchar(20);not null;default:''"`
	FMNetQuoUpUnit     decimal.Decimal `gorm:"column:FMNetQuoUpUnit;type:numeric;not null;default:0"`
	ExchNetQuoUpUnit   decimal.Decimal `gorm:"column:ExchNetQuoUpUnit;type:numeric;not null;default:0"`
	FundIPODate        time.Time       `gorm:"column:FundIPODate;type:timestamp with time zone;not null;default:'1900-01-01 00:00:00+08'"`
	BegRdmDate         time.Time       `gorm:"column:BegRdmDate;type:timestamp with time zone;not null;default:'1900-01-01 00:00:00+08'"`
	FundIssueType      string          `gorm:"column:FundIssueType;type:varchar(6);not null;default:''"`
	FundIssueDate      time.Time       `gorm:"column:FundIssueDate;type:timestamp with time zone;not null;default:'1900-01-01 00:00:00+08'"`
	CSDOnlineDate      time.Time       `gorm:"column:CSDOnlineDate;type:timestamp with time zone;not null;default:'1900-01-01 00:00:00+08'"`
	CSDPrtFundSetlDate time.Time       `gorm:"column:CSDPrtFundSetlDate;type:timestamp with time zone;not null;default:'1900-01-01 00:00:00+08'"`
	ApprovelDate       time.Time       `gorm:"column:ApprovelDate;type:timestamp with time zone;not null;default:'1900-01-01 00:00:00+08'"`
	ClearDate          time.Time       `gorm:"column:ClearDate;type:timestamp with time zone;not null;default:'1900-01-01 00:00:00+08'"`
	LastSubsDate       time.Time       `gorm:"column:LastSubsDate;type:timestamp with time zone;not null;default:'1900-01-01 00:00:00+08'"`
	LastRdmDate        time.Time       `gorm:"column:LastRdmDate;type:timestamp with time zone;not null;default:'1900-01-01 00:00:00+08'"`
	BIZDayDef          string          `gorm:"column:BIZDayDef;type:varchar(6);not null;default:''"`
	SubsDateType       string          `gorm:"column:SubsDateType;type:varchar(6);not null;default:''"`
	RdmDateType        string          `gorm:"column:RdmDateType;type:varchar(6);not null;default:''"`
	RdmPayDateType     string          `gorm:"column:RdmPayDateType;type:varchar(6);not null;default:''"`
	SubsNavDay         int16           `gorm:"column:SubsNavDay;type:smallint;not null;default:0"`
	RdmNavDay          int16           `gorm:"column:RdmNavDay;type:smallint;not null;default:0"`
	RdmPayDay          int16           `gorm:"column:RdmPayDay;type:smallint;not null;default:0"`
	CustBankHeadID     string          `gorm:"column:CustBankHeadID;type:varchar(10);not null;default:''"`
	FundAccMName       string          `gorm:"column:FundAccMName;type:varchar(200);not null;default:''"`
	FundAccSName       string          `gorm:"column:FundAccSName;type:varchar(200);not null;default:''"`
	CheckDiffDay       int16           `gorm:"column:CheckDiffDay;type:smallint;not null;default:0"`
	CashDiffDay        int16           `gorm:"column:CashDiffDay;type:smallint;not null;default:0"`
	FASubsInDay        int16           `gorm:"column:FASubsInDay;type:smallint;not null;default:0"`
	FARdmOutDay        int16           `gorm:"column:FARdmOutDay;type:smallint;not null;default:0"`
	FASubsSetlDay      int16           `gorm:"column:FASubsSetlDay;type:smallint;not null;default:0"`
	FARdmSetlDay       int16           `gorm:"column:FARdmSetlDay;type:smallint;not null;default:0"`
	SubsFeeGroup       string          `gorm:"column:SubsFeeGroup;type:varchar(6);not null;default:''"`
	CDSCFeeGroup       string          `gorm:"column:CDSCFeeGroup;type:varchar(6);not null;default:''"`
	DivFundGroup       string          `gorm:"column:DivFundGroup;type:varchar(6);not null;default:''"`
	CSDSubsDeadline    string          `gorm:"column:CSDSubsDeadline;type:varchar(10);not null;default:''"`
	CSDRdmDeadline     string          `gorm:"column:CSDRdmDeadline;type:varchar(10);not null;default:''"`
	BeginShortDate     time.Time       `gorm:"column:BeginShortDate;type:timestamp with time zone;not null;default:'1900-01-01 00:00:00+08'"`
	ShortCalcMemo      string          `gorm:"column:ShortCalcMemo;type:varchar(500);not null;default:''"`
	ShortTxMemo        string          `gorm:"column:ShortTxMemo;type:varchar(500);not null;default:''"`
	ShortMemo          string          `gorm:"column:ShortMemo;type:varchar(500);not null;default:''"`
	RemitChargeCD      string          `gorm:"column:RemitChargeCD;type:varchar(6);not null;default:''"`
	IsIPO              string          `gorm:"column:IsIPO;type:varchar(1);not null;default:''"`
	AuditFields
}

func (DTAFNDFundInfo) TableName() string { return "TA_STD_TH.DTA_FND_FundInfo" }

type DTAFNDFundInfoEdit struct{ DTAFNDFundInfo }

func (DTAFNDFundInfoEdit) TableName() string { return "TA_STD_TH.DTA_FND_FundInfo_Edit" }

// ═══════════════════════════════════════════════════════════════════
// Child: DTA_FND_FundCustContact (Custodian bank contacts)
// ═══════════════════════════════════════════════════════════════════

type DTAFNDFundCustContact struct {
	SysCoID       string `gorm:"column:SysCoID;type:varchar(20);primaryKey;not null"`
	PrtFundCode   string `gorm:"column:PrtFundCode;type:varchar(10);primaryKey;not null"`
	DataSeq       int16  `gorm:"column:DataSeq;type:smallint;primaryKey;not null;default:0"`
	IsMainContact string `gorm:"column:IsMainContact;type:varchar(1);not null;default:''"`
	ContactMan    string `gorm:"column:ContactMan;type:varchar(50);not null;default:''"`
	ContactTel    string `gorm:"column:ContactTel;type:varchar(50);not null;default:''"`
	Fax1          string `gorm:"column:Fax1;type:varchar(30);not null;default:''"`
	Fax2          string `gorm:"column:Fax2;type:varchar(30);not null;default:''"`
	WorkItem      string `gorm:"column:WorkItem;type:varchar(200);not null;default:''"`
	AuditFields
}

func (DTAFNDFundCustContact) TableName() string { return "TA_STD_TH.DTA_FND_FundCustContact" }

type DTAFNDFundCustContactEdit struct{ DTAFNDFundCustContact }

func (DTAFNDFundCustContactEdit) TableName() string {
	return "TA_STD_TH.DTA_FND_FundCustContact_Edit"
}

// ═══════════════════════════════════════════════════════════════════
// Child: DTA_FND_FundSpecial (Special fund types — PG, TM, PFF)
// ═══════════════════════════════════════════════════════════════════

type DTAFNDFundSpecial struct {
	SysCoID          string    `gorm:"column:SysCoID;type:varchar(20);primaryKey;not null"`
	PrtFundCode      string    `gorm:"column:PrtFundCode;type:varchar(10);primaryKey;not null"`
	IsHighRisk       string    `gorm:"column:IsHighRisk;type:varchar(1);not null;default:''"`
	IsPGF            string    `gorm:"column:IsPGF;type:varchar(1);not null;default:''"`
	PGFStay          string    `gorm:"column:PGFStay;type:varchar(6);not null;default:''"`
	PGFMatureYear    int16     `gorm:"column:PGFMatureYear;type:smallint;not null;default:0"`
	PGFStopSubsDay   int16     `gorm:"column:PGFStopSubsDay;type:smallint;not null;default:0"`
	PGFExpiryDate    time.Time `gorm:"column:PGFExpiryDate;type:timestamp with time zone;not null;default:'1900-01-01 00:00:00+08'"`
	IsPGFRdmFee      string    `gorm:"column:IsPGFRdmFee;type:varchar(1);not null;default:''"`
	IsTMF            string    `gorm:"column:IsTMF;type:varchar(1);not null;default:''"`
	TMFKind          string    `gorm:"column:TMFKind;type:varchar(6);not null;default:''"`
	TMFStopSubsDay   int16     `gorm:"column:TMFStopSubsDay;type:smallint;not null;default:0"`
	TMFTargetDate    time.Time `gorm:"column:TMFTargetDate;type:timestamp with time zone;not null;default:'1900-01-01 00:00:00+08'"`
	TMFStay          int16     `gorm:"column:TMFStay;type:smallint;not null;default:0"`
	TMFExpiryDate    time.Time `gorm:"column:TMFExpiryDate;type:timestamp with time zone;not null;default:'1900-01-01 00:00:00+08'"`
	TMFRetunrBC      string    `gorm:"column:TMFRetunrBC;type:varchar(6);not null;default:''"`
	TMFReturnStarr   string    `gorm:"column:TMFReturnStarr;type:varchar(10);not null;default:''"`
	IsTMFRdmFee      string    `gorm:"column:IsTMFRdmFee;type:varchar(1);not null;default:''"`
	IsPFF            string    `gorm:"column:IsPFF;type:varchar(1);not null;default:''"`
	PvtFundLimit     int16     `gorm:"column:PvtFundLimit;type:smallint;not null;default:0"`
	IsPerformanceFee string    `gorm:"column:IsPerformanceFee;type:varchar(1);not null;default:''"`
	IsUmbrella       string    `gorm:"column:IsUmbrella;type:varchar(1);not null;default:''"`
	UmbrellaName     string    `gorm:"column:UmbrellaName;type:varchar(200);not null;default:''"`
	AuditFields
}

func (DTAFNDFundSpecial) TableName() string { return "TA_STD_TH.DTA_FND_FundSpecial" }

type DTAFNDFundSpecialEdit struct{ DTAFNDFundSpecial }

func (DTAFNDFundSpecialEdit) TableName() string { return "TA_STD_TH.DTA_FND_FundSpecial_Edit" }

// ═══════════════════════════════════════════════════════════════════
// Child: DTA_FND_FundCryGroup (Issue currency groups)
// ═══════════════════════════════════════════════════════════════════

type DTAFNDFundCryGroup struct {
	SysCoID              string          `gorm:"column:SysCoID;type:varchar(20);primaryKey;not null"`
	PrtFundCode          string          `gorm:"column:PrtFundCode;type:varchar(10);primaryKey;not null"`
	IssueCry             string          `gorm:"column:IssueCry;type:varchar(6);primaryKey;not null"`
	FundFaceAmt          decimal.Decimal `gorm:"column:FundFaceAmt;type:numeric;not null;default:0"`
	NavDotLen            int16           `gorm:"column:NavDotLen;type:smallint;not null;default:0"`
	IsOthTxCry           string          `gorm:"column:IsOthTxCry;type:varchar(1);not null;default:''"`
	OthTxCry             string          `gorm:"column:OthTxCry;type:varchar(6);not null;default:''"`
	FundQuoAmtBC         decimal.Decimal `gorm:"column:FundQuoAmtBC;type:numeric;not null;default:0"`
	UnitConvRate         decimal.Decimal `gorm:"column:UnitConvRate;type:numeric;not null;default:0"`
	FundQuoUnit          decimal.Decimal `gorm:"column:FundQuoUnit;type:numeric;not null;default:0"`
	IPOSubsFeeDisType    string          `gorm:"column:IPOSubsFeeDisType;type:varchar(6);not null;default:''"`
	IPOFixRate           decimal.Decimal `gorm:"column:IPOFixRate;type:numeric;not null;default:0"`
	IPODiscRateOff       decimal.Decimal `gorm:"column:IPODiscRateOff;type:numeric;not null;default:0"`
	FundMBankBrh         string          `gorm:"column:FundMBankBrh;type:varchar(10);not null;default:''"`
	FundMAccount         string          `gorm:"column:FundMAccount;type:varchar(30);not null;default:''"`
	CRemitFeeID          string          `gorm:"column:CRemitFeeID;type:varchar(10);not null;default:''"`
	CRemitFeeFM          decimal.Decimal `gorm:"column:CRemitFeeFM;type:numeric;not null;default:0"`
	CRemitFeeOtherID     string          `gorm:"column:CRemitFeeOtherID;type:varchar(10);not null;default:''"`
	CRemitFeeOtherFM     decimal.Decimal `gorm:"column:CRemitFeeOtherFM;type:numeric;not null;default:0"`
	IsRSP                string          `gorm:"column:IsRSP;type:varchar(1);not null;default:''"`
	BeginRspDate         time.Time       `gorm:"column:BeginRspDate;type:timestamp with time zone;not null;default:'1900-01-01 00:00:00+08'"`
	RSPMinAmt            decimal.Decimal `gorm:"column:RSPMinAmt;type:numeric;not null;default:0"`
	RSPAmtApc            decimal.Decimal `gorm:"column:RSPAmtApc;type:numeric;not null;default:0"`
	RSPMaxAmt            decimal.Decimal `gorm:"column:RSPMaxAmt;type:numeric;not null;default:0"`
	RSPMinAmtTxC         decimal.Decimal `gorm:"column:RSPMinAmtTxC;type:numeric;not null;default:0"`
	RSPAmtApcTxC         decimal.Decimal `gorm:"column:RSPAmtApcTxC;type:numeric;not null;default:0"`
	RSPMaxAmtTxC         decimal.Decimal `gorm:"column:RSPMaxAmtTxC;type:numeric;not null;default:0"`
	RSPConFailCnt        int16           `gorm:"column:RSPConFailCnt;type:smallint;not null;default:0"`
	RSPFeeRate           decimal.Decimal `gorm:"column:RSPFeeRate;type:numeric;not null;default:0"`
	IsRSPChgAmt          string          `gorm:"column:IsRSPChgAmt;type:varchar(1);not null;default:''"`
	IsRSPShareSet        string          `gorm:"column:IsRSPShareSet;type:varchar(1);not null;default:''"`
	IsDividend           string          `gorm:"column:IsDividend;type:varchar(1);not null;default:''"`
	DivType              string          `gorm:"column:DivType;type:varchar(6);not null;default:''"`
	DivMinAmt            decimal.Decimal `gorm:"column:DivMinAmt;type:numeric;not null;default:0"`
	DivShortfallHandling string          `gorm:"column:DivShortfallHandling;type:varchar(6);not null;default:''"`
	DivCycle             string          `gorm:"column:DivCycle;type:varchar(10);not null;default:''"`
	DivPayPlatform       string          `gorm:"column:DivPayPlatform;type:varchar(6);not null;default:''"`
	TaxFormat            string          `gorm:"column:TaxFormat;type:varchar(10);not null;default:''"`
	LakhSubsAmt          decimal.Decimal `gorm:"column:LakhSubsAmt;type:numeric;not null;default:0"`
	LakhRdmAmt           decimal.Decimal `gorm:"column:LakhRdmAmt;type:numeric;not null;default:0"`
	AuditFields
}

func (DTAFNDFundCryGroup) TableName() string { return "TA_STD_TH.DTA_FND_FundCryGroup" }

type DTAFNDFundCryGroupEdit struct{ DTAFNDFundCryGroup }

func (DTAFNDFundCryGroupEdit) TableName() string { return "TA_STD_TH.DTA_FND_FundCryGroup_Edit" }

// ═══════════════════════════════════════════════════════════════════
// Child: DTA_FND_FundDetail (Sub-fund / share class details)
// ═══════════════════════════════════════════════════════════════════

type DTAFNDFundDetail struct {
	SysCoID           string          `gorm:"column:SysCoID;type:varchar(20);primaryKey;not null"`
	PrtFundCode       string          `gorm:"column:PrtFundCode;type:varchar(10);primaryKey;not null"`
	FundCode          string          `gorm:"column:FundCode;type:varchar(10);primaryKey;not null"`
	FundCry           string          `gorm:"column:FundCry;type:varchar(6);not null;default:''"`
	FundInShName      string          `gorm:"column:FundInShName;type:varchar(100);not null;default:''"`
	FundShMName       string          `gorm:"column:FundShMName;type:varchar(150);not null;default:''"`
	FundShSName       string          `gorm:"column:FundShSName;type:varchar(150);not null;default:''"`
	FHFundShare       string          `gorm:"column:FHFundShare;type:varchar(6);not null;default:''"`
	ShareFirSaleDate  time.Time       `gorm:"column:ShareFirSaleDate;type:timestamp with time zone;not null;default:'1900-01-01 00:00:00+08'"`
	SMARTFundCode     string          `gorm:"column:SMARTFundCode;type:varchar(10);not null;default:''"`
	CSDFundCode       string          `gorm:"column:CSDFundCode;type:varchar(10);not null;default:''"`
	SITCAFundCode     string          `gorm:"column:SITCAFundCode;type:varchar(20);not null;default:''"`
	FAFundCode        string          `gorm:"column:FAFundCode;type:varchar(20);not null;default:''"`
	FISFeeCode        string          `gorm:"column:FISFeeCode;type:varchar(20);not null;default:''"`
	ISINCode          string          `gorm:"column:ISINCode;type:varchar(20);not null;default:''"`
	IsTISA            string          `gorm:"column:IsTISA;type:varchar(1);not null;default:''"`
	TISAHoldM         int16           `gorm:"column:TISAHoldM;type:smallint;not null;default:0"`
	IsIShare          string          `gorm:"column:IsIShare;type:varchar(1);not null;default:''"`
	DivCalcType       string          `gorm:"column:DivCalcType;type:varchar(6);not null;default:''"`
	SubsFeeType       string          `gorm:"column:SubsFeeType;type:varchar(6);not null;default:''"`
	SubsMinAmt        decimal.Decimal `gorm:"column:SubsMinAmt;type:numeric;not null;default:0"`
	SubsAmtApc        decimal.Decimal `gorm:"column:SubsAmtApc;type:numeric;not null;default:0"`
	SwMinAmt          decimal.Decimal `gorm:"column:SwMinAmt;type:numeric;not null;default:0"`
	IsDistributionFee string          `gorm:"column:IsDistributionFee;type:varchar(1);not null;default:''"`
	MatureYear        int16           `gorm:"column:MatureYear;type:smallint;not null;default:0"`
	MatureFund        string          `gorm:"column:MatureFund;type:varchar(10);not null;default:''"`
	CSDSubsSetlDate   time.Time       `gorm:"column:CSDSubsSetlDate;type:timestamp with time zone;not null;default:'1900-01-01 00:00:00+08'"`
	RdmMinUnitSur     decimal.Decimal `gorm:"column:RdmMinUnitSur;type:numeric;not null;default:0"`
	RdmMinUnit        decimal.Decimal `gorm:"column:RdmMinUnit;type:numeric;not null;default:0"`
	RdmMinAmt         decimal.Decimal `gorm:"column:RdmMinAmt;type:numeric;not null;default:0"`
	CSDRdmSetlDate    time.Time       `gorm:"column:CSDRdmSetlDate;type:timestamp with time zone;not null;default:'1900-01-01 00:00:00+08'"`
	IsShareRSP        string          `gorm:"column:IsShareRSP;type:varchar(1);not null;default:''"`
	RSPMinAmt         decimal.Decimal `gorm:"column:RSPMinAmt;type:numeric;not null;default:0"`
	RSPAmtApc         decimal.Decimal `gorm:"column:RSPAmtApc;type:numeric;not null;default:0"`
	RSPMaxAmt         decimal.Decimal `gorm:"column:RSPMaxAmt;type:numeric;not null;default:0"`
	SubsMinAmtTxC     decimal.Decimal `gorm:"column:SubsMinAmtTxC;type:numeric;not null;default:0"`
	SubsAmtApcTxC     decimal.Decimal `gorm:"column:SubsAmtApcTxC;type:numeric;not null;default:0"`
	SwMinAmtTxC       decimal.Decimal `gorm:"column:SwMinAmtTxC;type:numeric;not null;default:0"`
	RSPMinAmtTxC      decimal.Decimal `gorm:"column:RSPMinAmtTxC;type:numeric;not null;default:0"`
	RSPAmtApcTxC      decimal.Decimal `gorm:"column:RSPAmtApcTxC;type:numeric;not null;default:0"`
	RSPMaxAmtTxC      decimal.Decimal `gorm:"column:RSPMaxAmtTxC;type:numeric;not null;default:0"`
	IsAdditional      string          `gorm:"column:IsAdditional;type:varchar(1);not null;default:''"`
	IsOld             string          `gorm:"column:IsOld;type:varchar(1);not null;default:''"`
	ShareMemo         string          `gorm:"column:ShareMemo;type:varchar(500);not null;default:''"`
	AuditFields
}

func (DTAFNDFundDetail) TableName() string { return "TA_STD_TH.DTA_FND_FundDetail" }

type DTAFNDFundDetailEdit struct{ DTAFNDFundDetail }

func (DTAFNDFundDetailEdit) TableName() string { return "TA_STD_TH.DTA_FND_FundDetail_Edit" }

// ═══════════════════════════════════════════════════════════════════
// Child: DTA_FND_FundAccount (Sub-fund bank accounts)
// ═══════════════════════════════════════════════════════════════════

type DTAFNDFundAccount struct {
	SysCoID       string `gorm:"column:SysCoID;type:varchar(20);primaryKey;not null"`
	PrtFundCode   string `gorm:"column:PrtFundCode;type:varchar(10);primaryKey;not null"`
	FundCode      string `gorm:"column:FundCode;type:varchar(10);primaryKey;not null"`
	DataSeq       int16  `gorm:"column:DataSeq;type:smallint;primaryKey;not null;default:0"`
	CusBankBrh    string `gorm:"column:CusBankBrh;type:varchar(10);not null;default:''"`
	FundAccount   string `gorm:"column:FundAccount;type:varchar(30);not null;default:''"`
	SRemitVCode   string `gorm:"column:SRemitVCode;type:varchar(20);not null;default:''"`
	AccountUseSet string `gorm:"column:AccountUseSet;type:varchar(200);not null;default:''"`
	Remark        string `gorm:"column:Remark;type:varchar(500);not null;default:''"`
	AuditFields
}

func (DTAFNDFundAccount) TableName() string { return "TA_STD_TH.DTA_FND_FundAccount" }

type DTAFNDFundAccountEdit struct{ DTAFNDFundAccount }

func (DTAFNDFundAccountEdit) TableName() string { return "TA_STD_TH.DTA_FND_FundAccount_Edit" }

// ═══════════════════════════════════════════════════════════════════
// Child: DTA_FND_FundRelGroup (Sub-fund related groups)
// ═══════════════════════════════════════════════════════════════════

type DTAFNDFundRelGroup struct {
	SysCoID             string          `gorm:"column:SysCoID;type:varchar(20);primaryKey;not null"`
	PrtFundCode         string          `gorm:"column:PrtFundCode;type:varchar(10);primaryKey;not null"`
	FundCode            string          `gorm:"column:FundCode;type:varchar(10);primaryKey;not null"`
	FundGroupNo         string          `gorm:"column:FundGroupNo;type:varchar(10);primaryKey;not null"`
	GroupTypeNo         string          `gorm:"column:GroupTypeNo;type:varchar(10);primaryKey;not null"`
	EffDate             time.Time       `gorm:"column:EffDate;type:timestamp with time zone;not null;default:'1900-01-01 00:00:00+08'"`
	AnnManagerRate      decimal.Decimal `gorm:"column:AnnManagerRate;type:numeric;not null;default:0"`
	AnnCustodyRate      decimal.Decimal `gorm:"column:AnnCustodyRate;type:numeric;not null;default:0"`
	AnnDistributionRate decimal.Decimal `gorm:"column:AnnDistributionRate;type:numeric;not null;default:0"`
	AnnOtherFee         decimal.Decimal `gorm:"column:AnnOtherFee;type:numeric;not null;default:0"`
	AuditFields
}

func (DTAFNDFundRelGroup) TableName() string { return "TA_STD_TH.DTA_FND_FundRelGroup" }

type DTAFNDFundRelGroupEdit struct{ DTAFNDFundRelGroup }

func (DTAFNDFundRelGroupEdit) TableName() string { return "TA_STD_TH.DTA_FND_FundRelGroup_Edit" }

// ═══════════════════════════════════════════════════════════════════
// Child: DTA_FND_FundDisclosureFee (Disclosure fees)
// ═══════════════════════════════════════════════════════════════════

type DTAFNDFundDisclosureFee struct {
	SysCoID     string `gorm:"column:SysCoID;type:varchar(20);primaryKey;not null"`
	PrtFundCode string `gorm:"column:PrtFundCode;type:varchar(10);primaryKey;not null"`
	FundCode    string `gorm:"column:FundCode;type:varchar(10);primaryKey;not null"`
	DataSeq     int16  `gorm:"column:DataSeq;type:smallint;primaryKey;not null;default:0"`
	AuditFields
}

func (DTAFNDFundDisclosureFee) TableName() string { return "TA_STD_TH.DTA_FND_FundDisclosureFee" }

type DTAFNDFundDisclosureFeeEdit struct{ DTAFNDFundDisclosureFee }

func (DTAFNDFundDisclosureFeeEdit) TableName() string {
	return "TA_STD_TH.DTA_FND_FundDisclosureFee_Edit"
}

// ═══════════════════════════════════════════════════════════════════
// Child: DTA_FND_FundTxCry (Transaction currencies per sub-fund)
// ═══════════════════════════════════════════════════════════════════

type DTAFNDFundTxCry struct {
	SysCoID     string `gorm:"column:SysCoID;type:varchar(20);primaryKey;not null"`
	PrtFundCode string `gorm:"column:PrtFundCode;type:varchar(10);primaryKey;not null"`
	FundCode    string `gorm:"column:FundCode;type:varchar(10);primaryKey;not null"`
	CryID       string `gorm:"column:CryID;type:varchar(6);primaryKey;not null"`
	CryDataSrc  string `gorm:"column:CryDataSrc;type:varchar(6);not null;default:''"`
	AuditFields
}

func (DTAFNDFundTxCry) TableName() string { return "TA_STD_TH.DTA_FND_FundTxCry" }

type DTAFNDFundTxCryEdit struct{ DTAFNDFundTxCry }

func (DTAFNDFundTxCryEdit) TableName() string { return "TA_STD_TH.DTA_FND_FundTxCry_Edit" }

// ═══════════════════════════════════════════════════════════════════
// Child: DTA_FND_FundShortInf (Short-term trading fee info)
// ═══════════════════════════════════════════════════════════════════

type DTAFNDFundShortInf struct {
	SysCoID       string `gorm:"column:SysCoID;type:varchar(20);primaryKey;not null"`
	PrtFundCode   string `gorm:"column:PrtFundCode;type:varchar(10);primaryKey;not null"`
	RdmRangeType  string `gorm:"column:RdmRangeType;type:varchar(6);not null;default:''"`
	ShortDateType string `gorm:"column:ShortDateType;type:varchar(6);not null;default:''"`
	RdmBaseID     string `gorm:"column:RdmBaseID;type:varchar(10);not null;default:''"`
	SubsBaseID    string `gorm:"column:SubsBaseID;type:varchar(10);not null;default:''"`
	ShortCalcID   string `gorm:"column:ShortCalcID;type:varchar(10);not null;default:''"`
	AuditFields
}

func (DTAFNDFundShortInf) TableName() string { return "TA_STD_TH.DTA_FND_FundShortInf" }

type DTAFNDFundShortInfEdit struct{ DTAFNDFundShortInf }

func (DTAFNDFundShortInfEdit) TableName() string { return "TA_STD_TH.DTA_FND_FundShortInf_Edit" }

// ═══════════════════════════════════════════════════════════════════
// Child: DTA_FND_FundShortDtl (Short-term trading fee detail tiers)
// ═══════════════════════════════════════════════════════════════════

type DTAFNDFundShortDtl struct {
	SysCoID     string          `gorm:"column:SysCoID;type:varchar(20);primaryKey;not null"`
	PrtFundCode string          `gorm:"column:PrtFundCode;type:varchar(10);primaryKey;not null"`
	DataSeq     int16           `gorm:"column:DataSeq;type:smallint;primaryKey;not null;default:0"`
	RdmRange    int16           `gorm:"column:RdmRange;type:smallint;not null;default:0"`
	ShortRate   decimal.Decimal `gorm:"column:ShortRate;type:numeric;not null;default:0"`
	AuditFields
}

func (DTAFNDFundShortDtl) TableName() string { return "TA_STD_TH.DTA_FND_FundShortDtl" }

type DTAFNDFundShortDtlEdit struct{ DTAFNDFundShortDtl }

func (DTAFNDFundShortDtlEdit) TableName() string { return "TA_STD_TH.DTA_FND_FundShortDtl_Edit" }

// ═══════════════════════════════════════════════════════════════════
// Child: DTA_FND_FundAntiDilution (Anti-dilution fee settings)
// ═══════════════════════════════════════════════════════════════════

type DTAFNDFundAntiDilution struct {
	SysCoID        string          `gorm:"column:SysCoID;type:varchar(20);primaryKey;not null"`
	PrtFundCode    string          `gorm:"column:PrtFundCode;type:varchar(10);primaryKey;not null"`
	AntiDilSetType string          `gorm:"column:AntiDilSetType;type:varchar(6);primaryKey;not null"`
	EffDate        time.Time       `gorm:"column:EffDate;type:timestamp with time zone;primaryKey;not null;default:'1900-01-01 00:00:00+08'"`
	TermDate       time.Time       `gorm:"column:TermDate;type:timestamp with time zone;not null;default:'9999-12-31 23:59:59.99'"`
	AntiDilTrigger decimal.Decimal `gorm:"column:AntiDilTrigger;type:numeric;not null;default:0"`
	AntiDilFeeRate decimal.Decimal `gorm:"column:AntiDilFeeRate;type:numeric;not null;default:0"`
	AdjRsn         string          `gorm:"column:AdjRsn;type:varchar(200);not null;default:''"`
	Remark         string          `gorm:"column:Remark;type:varchar(500);not null;default:''"`
	AuditFields
}

func (DTAFNDFundAntiDilution) TableName() string { return "TA_STD_TH.DTA_FND_FundAntiDilution" }

type DTAFNDFundAntiDilutionEdit struct{ DTAFNDFundAntiDilution }

func (DTAFNDFundAntiDilutionEdit) TableName() string {
	return "TA_STD_TH.DTA_FND_FundAntiDilution_Edit"
}

// ═══════════════════════════════════════════════════════════════════
// Child: DTA_FND_FundMGTFeeInf (Management fee info)
// ═══════════════════════════════════════════════════════════════════

type DTAFNDFundMGTFeeInf struct {
	SysCoID     string    `gorm:"column:SysCoID;type:varchar(20);primaryKey;not null"`
	PrtFundCode string    `gorm:"column:PrtFundCode;type:varchar(10);primaryKey;not null"`
	FundCode    string    `gorm:"column:FundCode;type:varchar(10);primaryKey;not null"`
	FHFundShare string    `gorm:"column:FHFundShare;type:varchar(6);not null;default:''"`
	EffDate     time.Time `gorm:"column:EffDate;type:timestamp with time zone;primaryKey;not null;default:'1900-01-01 00:00:00+08'"`
	AuditFields
}

func (DTAFNDFundMGTFeeInf) TableName() string { return "TA_STD_TH.DTA_FND_FundMGTFeeInf" }

type DTAFNDFundMGTFeeInfEdit struct{ DTAFNDFundMGTFeeInf }

func (DTAFNDFundMGTFeeInfEdit) TableName() string { return "TA_STD_TH.DTA_FND_FundMGTFeeInf_Edit" }

// ═══════════════════════════════════════════════════════════════════
// Child: DTA_FND_FundMGTFeeDtl (Management fee detail tiers)
// ═══════════════════════════════════════════════════════════════════

type DTAFNDFundMGTFeeDtl struct {
	SysCoID     string    `gorm:"column:SysCoID;type:varchar(20);primaryKey;not null"`
	PrtFundCode string    `gorm:"column:PrtFundCode;type:varchar(10);primaryKey;not null"`
	FundCode    string    `gorm:"column:FundCode;type:varchar(10);primaryKey;not null"`
	EffDate     time.Time `gorm:"column:EffDate;type:timestamp with time zone;primaryKey;not null;default:'1900-01-01 00:00:00+08'"`
	DataSeq     int16     `gorm:"column:DataSeq;type:smallint;primaryKey;not null;default:0"`
	AuditFields
}

func (DTAFNDFundMGTFeeDtl) TableName() string { return "TA_STD_TH.DTA_FND_FundMGTFeeDtl" }

type DTAFNDFundMGTFeeDtlEdit struct{ DTAFNDFundMGTFeeDtl }

func (DTAFNDFundMGTFeeDtlEdit) TableName() string { return "TA_STD_TH.DTA_FND_FundMGTFeeDtl_Edit" }
