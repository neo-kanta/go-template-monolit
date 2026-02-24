package db

import (
	"time"
)

// ═══════════════════════════════════════════════════════════════════
// TA_STD_TH Schema Models for FNDM012
// ═══════════════════════════════════════════════════════════════════

// DTAFNDPGFundFeeRdm represents the DTA_FND_PGFundFeeRdm table.
type DTAFNDPGFundFeeRdm struct {
	SysCoID        string    `gorm:"column:SysCoID;type:varchar(20);primaryKey;not null"`
	PrtFundCode    string    `gorm:"column:PrtFundCode;type:varchar(10);primaryKey;not null"`
	RdmCalcBegDate time.Time `gorm:"column:RdmCalcBegDate;type:timestamp with time zone;primaryKey;not null"`
	FeeRate        float64   `gorm:"column:FeeRate;type:numeric;not null;default:0"`

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
func (DTAFNDPGFundFeeRdm) TableName() string {
	return "TA_STD_TH.DTA_FND_PGFundFeeRdm"
}

// DTAFNDPGFundFeeRdmEdit represents the DTA_FND_PGFundFeeRdm_Edit table.
type DTAFNDPGFundFeeRdmEdit struct {
	DTAFNDPGFundFeeRdm
}

// TableName overrides the table name used by GORM.
func (DTAFNDPGFundFeeRdmEdit) TableName() string {
	return "TA_STD_TH.DTA_FND_PGFundFeeRdm_Edit"
}
