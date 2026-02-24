package db

import (
	"time"
)

// ═══════════════════════════════════════════════════════════════════
// TA_STD_TH Schema Models for FNDM008
// ═══════════════════════════════════════════════════════════════════

// DTAFNDPauseTxn represents the DTA_FND_PauseTxn table.
type DTAFNDPauseTxn struct {
	SysCoID     string    `gorm:"column:SysCoID;type:varchar(20);primaryKey;not null"`
	PrtFundCode string    `gorm:"column:PrtFundCode;type:varchar(10);primaryKey;not null"`
	PTxnBegDate time.Time `gorm:"column:PTxnBegDate;type:timestamp with time zone;primaryKey;not null"`
	PTxnEndDate time.Time `gorm:"column:PTxnEndDate;type:timestamp with time zone;not null;default:'1900-01-01 00:00:00+08'"`
	Remark      string    `gorm:"column:Remark;type:varchar(200);not null;default:''"`

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

func (DTAFNDPauseTxn) TableName() string {
	return "TA_STD_TH.DTA_FND_PauseTxn"
}

type DTAFNDPauseTxnEdit struct {
	DTAFNDPauseTxn
}

func (DTAFNDPauseTxnEdit) TableName() string {
	return "TA_STD_TH.DTA_FND_PauseTxn_Edit"
}

// ═══════════════════════════════════════════════════════════════════

// DTAFNDPauseTxnCry represents the DTA_FND_PauseTxnCry table.
type DTAFNDPauseTxnCry struct {
	SysCoID     string    `gorm:"column:SysCoID;type:varchar(20);primaryKey;not null"`
	PrtFundCode string    `gorm:"column:PrtFundCode;type:varchar(10);primaryKey;not null"`
	PTxnBegDate time.Time `gorm:"column:PTxnBegDate;type:timestamp with time zone;primaryKey;not null"`
	CryID       string    `gorm:"column:CryID;type:varchar(3);primaryKey;not null"`

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

func (DTAFNDPauseTxnCry) TableName() string {
	return "TA_STD_TH.DTA_FND_PauseTxnCry"
}

type DTAFNDPauseTxnCryEdit struct {
	DTAFNDPauseTxnCry
}

func (DTAFNDPauseTxnCryEdit) TableName() string {
	return "TA_STD_TH.DTA_FND_PauseTxnCry_Edit"
}

// ═══════════════════════════════════════════════════════════════════

// DTAFNDPauseTxnDtl represents the DTA_FND_PauseTxnDtl table.
type DTAFNDPauseTxnDtl struct {
	SysCoID      string    `gorm:"column:SysCoID;type:varchar(20);primaryKey;not null"`
	PrtFundCode  string    `gorm:"column:PrtFundCode;type:varchar(10);primaryKey;not null"`
	PTxnBegDate  time.Time `gorm:"column:PTxnBegDate;type:timestamp with time zone;primaryKey;not null"`
	PauseTxnType string    `gorm:"column:PauseTxnType;type:varchar(6);primaryKey;not null"`

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

func (DTAFNDPauseTxnDtl) TableName() string {
	return "TA_STD_TH.DTA_FND_PauseTxnDtl"
}

type DTAFNDPauseTxnDtlEdit struct {
	DTAFNDPauseTxnDtl
}

func (DTAFNDPauseTxnDtlEdit) TableName() string {
	return "TA_STD_TH.DTA_FND_PauseTxnDtl_Edit"
}
