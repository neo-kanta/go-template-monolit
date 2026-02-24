package db

import (
	"time"
)

// ═══════════════════════════════════════════════════════════════════
// TA_STD_TH Schema Models for FNDM013
// ═══════════════════════════════════════════════════════════════════

// DTAFNDTMFundFeeRdm represents the DTA_FND_TMFundFeeRdm table.
type DTAFNDTMFundFeeRdm struct {
	SysCoID        string    `gorm:"column:SysCoID;type:varchar(20);primaryKey;not null"`
	PrtFundCode    string    `gorm:"column:PrtFundCode;type:varchar(10);primaryKey;not null"`
	RdmCalcBegDate time.Time `gorm:"column:RdmCalcBegDate;type:timestamp with time zone;primaryKey;not null"`

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
func (DTAFNDTMFundFeeRdm) TableName() string {
	return "TA_STD_TH.DTA_FND_TMFundFeeRdm"
}

// DTAFNDTMFundFeeRdmEdit represents the DTA_FND_TMFundFeeRdm_Edit table.
type DTAFNDTMFundFeeRdmEdit struct {
	DTAFNDTMFundFeeRdm
}

// TableName overrides the table name used by GORM.
func (DTAFNDTMFundFeeRdmEdit) TableName() string {
	return "TA_STD_TH.DTA_FND_TMFundFeeRdm_Edit"
}

// ═══════════════════════════════════════════════════════════════════

// DTAFNDTMFundFeeRdmDtl represents the DTA_FND_TMFundFeeRdmDtl table.
type DTAFNDTMFundFeeRdmDtl struct {
	SysCoID        string    `gorm:"column:SysCoID;type:varchar(20);primaryKey;not null"`
	PrtFundCode    string    `gorm:"column:PrtFundCode;type:varchar(10);primaryKey;not null"`
	RdmCalcBegDate time.Time `gorm:"column:RdmCalcBegDate;type:timestamp with time zone;primaryKey;not null"`
	RdmCalcEndDate time.Time `gorm:"column:RdmCalcEndDate;type:timestamp with time zone;primaryKey;not null"`
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
func (DTAFNDTMFundFeeRdmDtl) TableName() string {
	return "TA_STD_TH.DTA_FND_TMFundFeeRdmDtl"
}

// DTAFNDTMFundFeeRdmDtlEdit represents the DTA_FND_TMFundFeeRdmDtl_Edit table.
type DTAFNDTMFundFeeRdmDtlEdit struct {
	DTAFNDTMFundFeeRdmDtl
}

// TableName overrides the table name used by GORM.
func (DTAFNDTMFundFeeRdmDtlEdit) TableName() string {
	return "TA_STD_TH.DTA_FND_TMFundFeeRdmDtl_Edit"
}
