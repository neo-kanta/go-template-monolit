package db

import (
	"time"
)

// ═══════════════════════════════════════════════════════════════════
// TA_STD_TH Schema Models for FNDM003
// ═══════════════════════════════════════════════════════════════════

// DTAFNDSwitch represents the DTA_FND_Switch table.
type DTAFNDSwitch struct {
	SysCoID     string `gorm:"column:SysCoID;type:varchar(20);primaryKey;not null"`
	PrtFundCode string `gorm:"column:PrtFundCode;type:varchar(10);primaryKey;not null"`
	IsSwitchIn  string `gorm:"column:IsSwitchIn;type:varchar(6);not null;default:''"`

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
}

func (DTAFNDSwitch) TableName() string {
	return "TA_STD_TH.DTA_FND_Switch"
}

type DTAFNDSwitchEdit struct {
	DTAFNDSwitch
}

func (DTAFNDSwitchEdit) TableName() string {
	return "TA_STD_TH.DTA_FND_Switch_Edit"
}

// ═══════════════════════════════════════════════════════════════════

// DTAFNDSwitchFund represents the DTA_FND_SwitchFund table.
type DTAFNDSwitchFund struct {
	SysCoID        string `gorm:"column:SysCoID;type:varchar(20);primaryKey;not null"`
	PrtFundCode    string `gorm:"column:PrtFundCode;type:varchar(10);primaryKey;not null"`
	SwOPrtFundCode string `gorm:"column:SwOPrtFundCode;type:varchar(10);primaryKey;not null"`
	SwDateType     string `gorm:"column:SwDateType;type:varchar(6);not null;default:''"`

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
}

func (DTAFNDSwitchFund) TableName() string {
	return "TA_STD_TH.DTA_FND_SwitchFund"
}

type DTAFNDSwitchFundEdit struct {
	DTAFNDSwitchFund
}

func (DTAFNDSwitchFundEdit) TableName() string {
	return "TA_STD_TH.DTA_FND_SwitchFund_Edit"
}

// ═══════════════════════════════════════════════════════════════════

// DTAFNDFundFeeSwitch represents the DTA_FND_FundFeeSwitch table.
type DTAFNDFundFeeSwitch struct {
	SysCoID     string  `gorm:"column:SysCoID;type:varchar(20);primaryKey;not null"`
	PrtFundCode string  `gorm:"column:PrtFundCode;type:varchar(10);primaryKey;not null"`
	FundCode    string  `gorm:"column:FundCode;type:varchar(10);primaryKey;not null"`
	SwFundType  string  `gorm:"column:SwFundType;type:varchar(6);not null;default:''"`
	SwDiscType  string  `gorm:"column:SwDiscType;type:varchar(6);not null;default:''"`
	SwitchRate  float64 `gorm:"column:SwitchRate;type:numeric;not null;default:0"`

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
}

func (DTAFNDFundFeeSwitch) TableName() string {
	return "TA_STD_TH.DTA_FND_FundFeeSwitch"
}

type DTAFNDFundFeeSwitchEdit struct {
	DTAFNDFundFeeSwitch
}

func (DTAFNDFundFeeSwitchEdit) TableName() string {
	return "TA_STD_TH.DTA_FND_FundFeeSwitch_Edit"
}

// ═══════════════════════════════════════════════════════════════════

// DTAFNDSwitchCry represents the DTA_FND_SwitchCry table.
type DTAFNDSwitchCry struct {
	SysCoID       string `gorm:"column:SysCoID;type:varchar(20);primaryKey;not null"`
	PrtFundCode   string `gorm:"column:PrtFundCode;type:varchar(10);primaryKey;not null"`
	SwIFundCry    string `gorm:"column:SwIFundCry;type:varchar(6);primaryKey;not null"`
	SwOFundCrySet string `gorm:"column:SwOFundCrySet;type:varchar(200);not null;default:''"`

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
}

func (DTAFNDSwitchCry) TableName() string {
	return "TA_STD_TH.DTA_FND_SwitchCry"
}

type DTAFNDSwitchCryEdit struct {
	DTAFNDSwitchCry
}

func (DTAFNDSwitchCryEdit) TableName() string {
	return "TA_STD_TH.DTA_FND_SwitchCry_Edit"
}
