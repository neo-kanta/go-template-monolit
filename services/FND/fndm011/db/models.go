package db

import (
	"github.com/shopspring/decimal"

	models "go-transfer-agent/common/platform/model"

	"time"
)

// ═══════════════════════════════════════════════════════════════════
// TA_STD_TH Schema Models for FNDM011
// ═══════════════════════════════════════════════════════════════════

// DTAFNDIShareFundFeeRdm represents the DTA_FND_IShareFundFeeRdm table.
type DTAFNDIShareFundFeeRdm struct {
	SysCoID        string          `gorm:"column:SysCoID;type:varchar(20);primaryKey;not null"`
	FundCode       string          `gorm:"column:FundCode;type:varchar(10);primaryKey;not null"`
	FeeName        string          `gorm:"column:FeeName;type:varchar(50);primaryKey;not null"`
	RdmRangeType   string          `gorm:"column:RdmRangeType;type:varchar(6);not null;default:''"`
	RdmDateType    string          `gorm:"column:RdmDateType;type:varchar(6);not null;default:''"`
	RdmBaseID      string          `gorm:"column:RdmBaseID;type:varchar(6);not null;default:''"`
	SubsBaseID     string          `gorm:"column:SubsBaseID;type:varchar(6);not null;default:''"`
	RdmCalcID      string          `gorm:"column:RdmCalcID;type:varchar(6);not null;default:''"`
	ShouldHoldDays int16           `gorm:"column:ShouldHoldDays;type:smallint;not null;default:0"`
	FeeRate        decimal.Decimal `gorm:"column:FeeRate;type:numeric;not null;default:0"`

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

// TableName overrides the table name used by GORM.
func (DTAFNDIShareFundFeeRdm) TableName() string {
	return "TA_STD_TH.DTA_FND_IShareFundFeeRdm"
}

// DTAFNDIShareFundFeeRdmEdit represents the DTA_FND_IShareFundFeeRdm_Edit table.
type DTAFNDIShareFundFeeRdmEdit struct {
	DTAFNDIShareFundFeeRdm
}

// TableName overrides the table name used by GORM.
func (DTAFNDIShareFundFeeRdmEdit) TableName() string {
	return "TA_STD_TH.DTA_FND_IShareFundFeeRdm_Edit"
}
