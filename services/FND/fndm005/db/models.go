package db

import (
	models "go-transfer-agent/common/platform/model"
	"time"
)

// ═══════════════════════════════════════════════════════════════════
// TA_STD_TH Schema Models for FNDM005
// ═══════════════════════════════════════════════════════════════════

// TAFNDFundCal represents the TA_FND_FundCal table.
type TAFNDFundCal struct {
	SysCoID    string `gorm:"column:SysCoID;type:varchar(20);primaryKey;not null"`
	CalYear    string `gorm:"column:CalYear;type:varchar(4);primaryKey;not null"`
	FNDCalType string `gorm:"column:FNDCalType;type:varchar(6);primaryKey;not null"`
	FundCry    string `gorm:"column:FundCry;type:varchar(3);primaryKey;not null"`

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

func (TAFNDFundCal) TableName() string {
	return "TA_STD_TH.TA_FND_FundCal"
}

type TAFNDFundCalEdit struct {
	TAFNDFundCal
}

func (TAFNDFundCalEdit) TableName() string {
	return "TA_STD_TH.TA_FND_FundCal_Edit"
}

// ═══════════════════════════════════════════════════════════════════

// TAFNDFundCalMemo represents the TA_FND_FundCalMemo table.
type TAFNDFundCalMemo struct {
	SysCoID    string    `gorm:"column:SysCoID;type:varchar(20);primaryKey;not null"`
	CalYear    string    `gorm:"column:CalYear;type:varchar(4);primaryKey;not null"`
	FNDCalType string    `gorm:"column:FNDCalType;type:varchar(6);primaryKey;not null"`
	FundCry    string    `gorm:"column:FundCry;type:varchar(3);primaryKey;not null"`
	CalDate    time.Time `gorm:"column:CalDate;type:timestamp with time zone;primaryKey;not null"`

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

func (TAFNDFundCalMemo) TableName() string {
	return "TA_STD_TH.TA_FND_FundCalMemo"
}

type TAFNDFundCalMemoEdit struct {
	TAFNDFundCalMemo
}

func (TAFNDFundCalMemoEdit) TableName() string {
	return "TA_STD_TH.TA_FND_FundCalMemo_Edit"
}

// ═══════════════════════════════════════════════════════════════════

// TAFNDFundCalDtl represents the TA_FND_FundCalDtl table.
type TAFNDFundCalDtl struct {
	SysCoID     string    `gorm:"column:SysCoID;type:varchar(20);primaryKey;not null"`
	CalYear     string    `gorm:"column:CalYear;type:varchar(4);primaryKey;not null"`
	FNDCalType  string    `gorm:"column:FNDCalType;type:varchar(6);primaryKey;not null"`
	FundCry     string    `gorm:"column:FundCry;type:varchar(3);primaryKey;not null"`
	PrtFundCode string    `gorm:"column:PrtFundCode;type:varchar(10);primaryKey;not null"`
	CalDate     time.Time `gorm:"column:CalDate;type:timestamp with time zone;primaryKey;not null"`

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

func (TAFNDFundCalDtl) TableName() string {
	return "TA_STD_TH.TA_FND_FundCalDtl"
}

type TAFNDFundCalDtlEdit struct {
	TAFNDFundCalDtl
}

func (TAFNDFundCalDtlEdit) TableName() string {
	return "TA_STD_TH.TA_FND_FundCalDtl_Edit"
}
