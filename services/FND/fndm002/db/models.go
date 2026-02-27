package db

import (
	"github.com/shopspring/decimal"

	models "go-transfer-agent/common/platform/model"

	"time"
)

// ═══════════════════════════════════════════════════════════════════
// TA_STD_TH Schema Models for FNDM002
// ═══════════════════════════════════════════════════════════════════

// DTAFNDFundFee represents the DTA_FND_FundFee table.
type DTAFNDFundFee struct {
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

func (DTAFNDFundFee) TableName() string {
	return "TA_STD_TH.DTA_FND_FundFee"
}

type DTAFNDFundFeeEdit struct {
	DTAFNDFundFee
}

func (DTAFNDFundFeeEdit) TableName() string {
	return "TA_STD_TH.DTA_FND_FundFee_Edit"
}

// ═══════════════════════════════════════════════════════════════════

// DTAFNDFundFeeSub represents the DTA_FND_FundFeeSub table.
type DTAFNDFundFeeSub struct {
	SysCoID     string `gorm:"column:SysCoID;type:varchar(20);primaryKey;not null"`
	PrtFundCode string `gorm:"column:PrtFundCode;type:varchar(10);primaryKey;not null"`
	FundCode    string `gorm:"column:FundCode;type:varchar(10);primaryKey;not null"`
	CryID       string `gorm:"column:CryID;type:varchar(6);primaryKey;not null"`

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

func (DTAFNDFundFeeSub) TableName() string {
	return "TA_STD_TH.DTA_FND_FundFeeSub"
}

type DTAFNDFundFeeSubEdit struct {
	DTAFNDFundFeeSub
}

func (DTAFNDFundFeeSubEdit) TableName() string {
	return "TA_STD_TH.DTA_FND_FundFeeSub_Edit"
}

// ═══════════════════════════════════════════════════════════════════

// DTAFNDFundFeeSubDtl represents the DTA_FND_FundFeeSubDtl table.
type DTAFNDFundFeeSubDtl struct {
	SysCoID       string          `gorm:"column:SysCoID;type:varchar(20);primaryKey;not null"`
	PrtFundCode   string          `gorm:"column:PrtFundCode;type:varchar(10);primaryKey;not null"`
	FundCode      string          `gorm:"column:FundCode;type:varchar(10);primaryKey;not null"`
	CryID         string          `gorm:"column:CryID;type:varchar(6);primaryKey;not null"`
	RangeAmtAbove decimal.Decimal `gorm:"column:RangeAmtAbove;type:numeric;primaryKey;not null;default:0"`
	SubsFeeRate   decimal.Decimal `gorm:"column:SubsFeeRate;type:numeric;not null;default:0"`

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

func (DTAFNDFundFeeSubDtl) TableName() string {
	return "TA_STD_TH.DTA_FND_FundFeeSubDtl"
}

type DTAFNDFundFeeSubDtlEdit struct {
	DTAFNDFundFeeSubDtl
}

func (DTAFNDFundFeeSubDtlEdit) TableName() string {
	return "TA_STD_TH.DTA_FND_FundFeeSubDtl_Edit"
}

// ═══════════════════════════════════════════════════════════════════

// DTAFNDFundFeeCDSC represents the DTA_FND_FundFeeCDSC table.
type DTAFNDFundFeeCDSC struct {
	SysCoID     string `gorm:"column:SysCoID;type:varchar(20);primaryKey;not null"`
	PrtFundCode string `gorm:"column:PrtFundCode;type:varchar(10);primaryKey;not null"`
	FundCode    string `gorm:"column:FundCode;type:varchar(10);primaryKey;not null"`
	MatureYear  string `gorm:"column:MatureYear;type:varchar(6);primaryKey;not null"`
	SubsCalcID  string `gorm:"column:SubsCalcID;type:varchar(6);not null;default:''"`

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

func (DTAFNDFundFeeCDSC) TableName() string {
	return "TA_STD_TH.DTA_FND_FundFeeCDSC"
}

type DTAFNDFundFeeCDSCEdit struct {
	DTAFNDFundFeeCDSC
}

func (DTAFNDFundFeeCDSCEdit) TableName() string {
	return "TA_STD_TH.DTA_FND_FundFeeCDSC_Edit"
}

// ═══════════════════════════════════════════════════════════════════

// DTAFNDFundFeeCDSCDtl represents the DTA_FND_FundFeeCDSCDtl table.
type DTAFNDFundFeeCDSCDtl struct {
	SysCoID     string          `gorm:"column:SysCoID;type:varchar(20);primaryKey;not null"`
	PrtFundCode string          `gorm:"column:PrtFundCode;type:varchar(10);primaryKey;not null"`
	FundCode    string          `gorm:"column:FundCode;type:varchar(10);primaryKey;not null"`
	MatureYear  string          `gorm:"column:MatureYear;type:varchar(6);primaryKey;not null"`
	HoldBegDay  int16           `gorm:"column:HoldBegDay;type:smallint;primaryKey;not null;default:0"`
	CDSCFeeRate decimal.Decimal `gorm:"column:CDSCFeeRate;type:numeric;not null;default:0"`

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

func (DTAFNDFundFeeCDSCDtl) TableName() string {
	return "TA_STD_TH.DTA_FND_FundFeeCDSCDtl"
}

type DTAFNDFundFeeCDSCDtlEdit struct {
	DTAFNDFundFeeCDSCDtl
}

func (DTAFNDFundFeeCDSCDtlEdit) TableName() string {
	return "TA_STD_TH.DTA_FND_FundFeeCDSCDtl_Edit"
}

// ═══════════════════════════════════════════════════════════════════

// DTAFNDFundFeeBack represents the DTA_FND_FundFeeBack table.
type DTAFNDFundFeeBack struct {
	SysCoID        string `gorm:"column:SysCoID;type:varchar(20);primaryKey;not null"`
	PrtFundCode    string `gorm:"column:PrtFundCode;type:varchar(10);primaryKey;not null"`
	FundCode       string `gorm:"column:FundCode;type:varchar(10);primaryKey;not null"`
	HoldPeriodYear string `gorm:"column:HoldPeriodYear;type:varchar(6);primaryKey;not null"`
	SubsCalcID     string `gorm:"column:SubsCalcID;type:varchar(6);not null;default:''"`

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

func (DTAFNDFundFeeBack) TableName() string {
	return "TA_STD_TH.DTA_FND_FundFeeBack"
}

type DTAFNDFundFeeBackEdit struct {
	DTAFNDFundFeeBack
}

func (DTAFNDFundFeeBackEdit) TableName() string {
	return "TA_STD_TH.DTA_FND_FundFeeBack_Edit"
}

// ═══════════════════════════════════════════════════════════════════

// DTAFNDFundFeeBackDtl represents the DTA_FND_FundFeeBackDtl table.
type DTAFNDFundFeeBackDtl struct {
	SysCoID        string          `gorm:"column:SysCoID;type:varchar(20);primaryKey;not null"`
	PrtFundCode    string          `gorm:"column:PrtFundCode;type:varchar(10);primaryKey;not null"`
	FundCode       string          `gorm:"column:FundCode;type:varchar(10);primaryKey;not null"`
	HoldPeriodYear string          `gorm:"column:HoldPeriodYear;type:varchar(6);primaryKey;not null"`
	HoldBegDay     int16           `gorm:"column:HoldBegDay;type:smallint;primaryKey;not null;default:0"`
	BackFeeRate    decimal.Decimal `gorm:"column:BackFeeRate;type:numeric;not null;default:0"`

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

func (DTAFNDFundFeeBackDtl) TableName() string {
	return "TA_STD_TH.DTA_FND_FundFeeBackDtl"
}

type DTAFNDFundFeeBackDtlEdit struct {
	DTAFNDFundFeeBackDtl
}

func (DTAFNDFundFeeBackDtlEdit) TableName() string {
	return "TA_STD_TH.DTA_FND_FundFeeBackDtl_Edit"
}
