package db

import (
	"github.com/shopspring/decimal"

	models "go-transfer-agent/common/platform/model"

	"time"
)

// ═══════════════════════════════════════════════════════════════════
// TA_STD_TH Schema Models for FNDM009
// ═══════════════════════════════════════════════════════════════════

// TAFNDFavDisc represents the TA_FND_FavDisc table.
type TAFNDFavDisc struct {
	SysCoID      string `gorm:"column:SysCoID;type:varchar(20);primaryKey;not null"`
	CusIDCode    string `gorm:"column:CusIDCode;type:varchar(10);primaryKey;not null"`
	DiscItemSet  string `gorm:"column:DiscItemSet;type:varchar(200);not null;default:''"`
	DiscFundType string `gorm:"column:DiscFundType;type:varchar(6);not null;default:''"`

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

func (TAFNDFavDisc) TableName() string {
	return "TA_STD_TH.TA_FND_FavDisc"
}

// TAFNDFavDiscEdit represents the TA_FND_FavDisc_Edit table.
type TAFNDFavDiscEdit struct {
	TAFNDFavDisc
}

func (TAFNDFavDiscEdit) TableName() string {
	return "TA_STD_TH.TA_FND_FavDisc_Edit"
}

// ═══════════════════════════════════════════════════════════════════

// TAFNDFavDiscFund represents the TA_FND_FavDiscFund table.
type TAFNDFavDiscFund struct {
	SysCoID     string `gorm:"column:SysCoID;type:varchar(20);primaryKey;not null"`
	CusIDCode   string `gorm:"column:CusIDCode;type:varchar(10);primaryKey;not null"`
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

func (TAFNDFavDiscFund) TableName() string {
	return "TA_STD_TH.TA_FND_FavDiscFund"
}

type TAFNDFavDiscFundEdit struct {
	TAFNDFavDiscFund
}

func (TAFNDFavDiscFundEdit) TableName() string {
	return "TA_STD_TH.TA_FND_FavDiscFund_Edit"
}

// ═══════════════════════════════════════════════════════════════════

// TAFNDFavDiscType represents the TA_FND_FavDiscType table.
type TAFNDFavDiscType struct {
	SysCoID   string `gorm:"column:SysCoID;type:varchar(20);primaryKey;not null"`
	CusIDCode string `gorm:"column:CusIDCode;type:varchar(10);primaryKey;not null"`
	DiscItem  string `gorm:"column:DiscItem;type:varchar(6);primaryKey;not null"`
	DiscType  string `gorm:"column:DiscType;type:varchar(6);not null;default:''"`

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

func (TAFNDFavDiscType) TableName() string {
	return "TA_STD_TH.TA_FND_FavDiscType"
}

type TAFNDFavDiscTypeEdit struct {
	TAFNDFavDiscType
}

func (TAFNDFavDiscTypeEdit) TableName() string {
	return "TA_STD_TH.TA_FND_FavDiscType_Edit"
}

// ═══════════════════════════════════════════════════════════════════

// TAFNDFavDiscTypeDtl represents the TA_FND_FavDiscTypeDtl table.
type TAFNDFavDiscTypeDtl struct {
	SysCoID       string          `gorm:"column:SysCoID;type:varchar(20);primaryKey;not null"`
	CusIDCode     string          `gorm:"column:CusIDCode;type:varchar(10);primaryKey;not null"`
	DiscItem      string          `gorm:"column:DiscItem;type:varchar(6);primaryKey;not null"`
	TxCry         string          `gorm:"column:TxCry;type:varchar(3);primaryKey;not null"`
	RangeAmtAbove decimal.Decimal `gorm:"column:RangeAmtAbove;type:numeric;primaryKey;not null;default:0"`
	RangeFeeRate  decimal.Decimal `gorm:"column:RangeFeeRate;type:numeric;not null;default:0"`

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

func (TAFNDFavDiscTypeDtl) TableName() string {
	return "TA_STD_TH.TA_FND_FavDiscTypeDtl"
}

type TAFNDFavDiscTypeDtlEdit struct {
	TAFNDFavDiscTypeDtl
}

func (TAFNDFavDiscTypeDtlEdit) TableName() string {
	return "TA_STD_TH.TA_FND_FavDiscTypeDtl_Edit"
}
