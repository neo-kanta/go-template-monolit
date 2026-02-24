package db

import (
	"time"
)

// ═══════════════════════════════════════════════════════════════════
// TA_STD_TH Schema Models for FNDM010
// ═══════════════════════════════════════════════════════════════════

// DTAFNDIShareFundFee represents the DTA_FND_IShareFundFee table.
type DTAFNDIShareFundFee struct {
	SysCoID  string `gorm:"column:SysCoID;type:varchar(20);primaryKey;not null"`
	FundCode string `gorm:"column:FundCode;type:varchar(10);primaryKey;not null"`

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

// TableName overrides the table name used by GORM.
func (DTAFNDIShareFundFee) TableName() string {
	return "TA_STD_TH.DTA_FND_IShareFundFee"
}

// DTAFNDIShareFundFeeEdit represents the DTA_FND_IShareFundFee_Edit table.
type DTAFNDIShareFundFeeEdit struct {
	DTAFNDIShareFundFee
}

// TableName overrides the table name used by GORM.
func (DTAFNDIShareFundFeeEdit) TableName() string {
	return "TA_STD_TH.DTA_FND_IShareFundFee_Edit"
}

// ═══════════════════════════════════════════════════════════════════

// DTAFNDIShareFundFeeSub represents the DTA_FND_IShareFundFeeSub table.
type DTAFNDIShareFundFeeSub struct {
	SysCoID       string  `gorm:"column:SysCoID;type:varchar(20);primaryKey;not null"`
	FundCode      string  `gorm:"column:FundCode;type:varchar(10);primaryKey;not null"`
	RangeAmtAbove float64 `gorm:"column:RangeAmtAbove;type:numeric;primaryKey;not null;default:0"`
	SubsFeeRate   float64 `gorm:"column:SubsFeeRate;type:numeric;not null;default:0"`

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

// TableName overrides the table name used by GORM.
func (DTAFNDIShareFundFeeSub) TableName() string {
	return "TA_STD_TH.DTA_FND_IShareFundFeeSub"
}

// DTAFNDIShareFundFeeSubEdit represents the DTA_FND_IShareFundFeeSub_Edit table.
type DTAFNDIShareFundFeeSubEdit struct {
	DTAFNDIShareFundFeeSub
}

// TableName overrides the table name used by GORM.
func (DTAFNDIShareFundFeeSubEdit) TableName() string {
	return "TA_STD_TH.DTA_FND_IShareFundFeeSub_Edit"
}
