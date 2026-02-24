package db

import (
	"time"
)

// ═══════════════════════════════════════════════════════════════════
// TA_STD_TH Schema Models for FNDM007
// ═══════════════════════════════════════════════════════════════════

// TAFNDCustGroup represents the TA_FND_CustGroup table.
type TAFNDCustGroup struct {
	SysCoID      string `gorm:"column:SysCoID;type:varchar(20);primaryKey;not null"`
	CustGrpCode  string `gorm:"column:CustGrpCode;type:varchar(6);primaryKey;not null"`
	CustGrpMName string `gorm:"column:CustGrpMName;type:varchar(200);not null;default:''"`
	CustGrpSName string `gorm:"column:CustGrpSName;type:varchar(200);not null;default:''"`

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

func (TAFNDCustGroup) TableName() string {
	return "TA_STD_TH.TA_FND_CustGroup"
}

type TAFNDCustGroupEdit struct {
	TAFNDCustGroup
}

func (TAFNDCustGroupEdit) TableName() string {
	return "TA_STD_TH.TA_FND_CustGroup_Edit"
}

// ═══════════════════════════════════════════════════════════════════

// TAFNDCustGroupDtl represents the TA_FND_CustGroupDtl table.
type TAFNDCustGroupDtl struct {
	SysCoID     string `gorm:"column:SysCoID;type:varchar(20);primaryKey;not null"`
	CustGrpCode string `gorm:"column:CustGrpCode;type:varchar(6);primaryKey;not null"`
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
}

func (TAFNDCustGroupDtl) TableName() string {
	return "TA_STD_TH.TA_FND_CustGroupDtl"
}

type TAFNDCustGroupDtlEdit struct {
	TAFNDCustGroupDtl
}

func (TAFNDCustGroupDtlEdit) TableName() string {
	return "TA_STD_TH.TA_FND_CustGroupDtl_Edit"
}
