package db

import (
	"github.com/shopspring/decimal"

	models "go-transfer-agent/common/platform/model"

	"time"
)

// ═══════════════════════════════════════════════════════════════════
// TA_STD_TH Schema Models for FNDM004
// ═══════════════════════════════════════════════════════════════════

// DTAFNDFundAgent represents the DTA_FND_FundAgent table.
type DTAFNDFundAgent struct {
	SysCoID     string `gorm:"column:SysCoID;type:varchar(20);primaryKey;not null"`
	PrtFundCode string `gorm:"column:PrtFundCode;type:varchar(10);primaryKey;not null"`

	// Standard Audit Fields
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

func (DTAFNDFundAgent) TableName() string {
	return "TA_STD_TH.DTA_FND_FundAgent"
}

type DTAFNDFundAgentEdit struct {
	DTAFNDFundAgent
}

func (DTAFNDFundAgentEdit) TableName() string {
	return "TA_STD_TH.DTA_FND_FundAgent_Edit"
}

// ═══════════════════════════════════════════════════════════════════

// DTAFNDFundAgentDtl represents the DTA_FND_FundAgentDtl table.
type DTAFNDFundAgentDtl struct {
	SysCoID           string          `gorm:"column:SysCoID;type:varchar(20);primaryKey;not null"`
	PrtFundCode       string          `gorm:"column:PrtFundCode;type:varchar(10);primaryKey;not null"`
	AgentType         string          `gorm:"column:AgentType;type:varchar(6);primaryKey;not null"`
	AgentCode         string          `gorm:"column:AgentCode;type:varchar(20);primaryKey;not null"`
	IsMAgentOPType    bool            `gorm:"column:IsMAgentOPType;type:boolean;not null;default:false"`
	AgentOPType       string          `gorm:"column:AgentOPType;type:varchar(6);not null;default:''"`
	IsMRcvTxnType     bool            `gorm:"column:IsMRcvTxnType;type:boolean;not null;default:false"`
	RcvTxnType        string          `gorm:"column:RcvTxnType;type:varchar(6);not null;default:''"`
	IsMSubsFeePct     bool            `gorm:"column:IsMSubsFeePct;type:boolean;not null;default:false"`
	SubsFeePctAG      decimal.Decimal `gorm:"column:SubsFeePctAG;type:numeric;not null;default:0"`
	SubsFeePctFH      decimal.Decimal `gorm:"column:SubsFeePctFH;type:numeric;not null;default:0"`
	IsMSubsFeePctType bool            `gorm:"column:IsMSubsFeePctType;type:boolean;not null;default:false"`
	SubsFeePctType    string          `gorm:"column:SubsFeePctType;type:varchar(6);not null;default:''"`
	IsMFundCrySet     bool            `gorm:"column:IsMFundCrySet;type:boolean;not null;default:false"`
	FundCrySet        string          `gorm:"column:FundCrySet;type:varchar(200);not null;default:''"`
	NoSaleShareSet    string          `gorm:"column:NoSaleShareSet;type:varchar(200);not null;default:''"`
	IsMEndOfSale      bool            `gorm:"column:IsMEndOfSale;type:boolean;not null;default:false"`
	TermDate          time.Time       `gorm:"column:TermDate;type:timestamp with time zone;not null;default:'1900-01-01 00:00:00+08'"`

	// Standard Audit Fields
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

func (DTAFNDFundAgentDtl) TableName() string {
	return "TA_STD_TH.DTA_FND_FundAgentDtl"
}

type DTAFNDFundAgentDtlEdit struct {
	DTAFNDFundAgentDtl
}

func (DTAFNDFundAgentDtlEdit) TableName() string {
	return "TA_STD_TH.DTA_FND_FundAgentDtl_Edit"
}
