package db

import (
	models "go-transfer-agent/common/platform/model"
	"time"
)

// ═══════════════════════════════════════════════════════════════════
// TA_STD_TH Schema Models for FNDM006
// ═══════════════════════════════════════════════════════════════════

// DTAFNDRPFeeChgType represents the DTA_FND_RPFeeChgType table.
type DTAFNDRPFeeChgType struct {
	SysCoID     string `gorm:"column:SysCoID;type:varchar(20);primaryKey;not null"`
	PmtTxnType  string `gorm:"column:PmtTxnType;type:varchar(6);primaryKey;not null"`
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

func (DTAFNDRPFeeChgType) TableName() string {
	return "TA_STD_TH.DTA_FND_RPFeeChgType"
}

type DTAFNDRPFeeChgTypeEdit struct {
	DTAFNDRPFeeChgType
}

func (DTAFNDRPFeeChgTypeEdit) TableName() string {
	return "TA_STD_TH.DTA_FND_RPFeeChgType_Edit"
}

// ═══════════════════════════════════════════════════════════════════

// DTAFNDRPFeeChgTypeDtl represents the DTA_FND_RPFeeChgTypeDtl table.
type DTAFNDRPFeeChgTypeDtl struct {
	SysCoID         string `gorm:"column:SysCoID;type:varchar(20);primaryKey;not null"`
	PmtTxnType      string `gorm:"column:PmtTxnType;type:varchar(6);primaryKey;not null"`
	PrtFundCode     string `gorm:"column:PrtFundCode;type:varchar(10);primaryKey;not null"`
	CryKind         string `gorm:"column:CryKind;type:varchar(6);primaryKey;not null"`
	RemitFeeObj     string `gorm:"column:RemitFeeObj;type:varchar(6);not null;default:''"`
	RemitFeeChgType string `gorm:"column:RemitFeeChgType;type:varchar(6);not null;default:''"`
	PostFeeObj      string `gorm:"column:PostFeeObj;type:varchar(6);not null;default:''"`

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

func (DTAFNDRPFeeChgTypeDtl) TableName() string {
	return "TA_STD_TH.DTA_FND_RPFeeChgTypeDtl"
}

type DTAFNDRPFeeChgTypeDtlEdit struct {
	DTAFNDRPFeeChgTypeDtl
}

func (DTAFNDRPFeeChgTypeDtlEdit) TableName() string {
	return "TA_STD_TH.DTA_FND_RPFeeChgTypeDtl_Edit"
}
