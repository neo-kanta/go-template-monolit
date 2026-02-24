package db

import (
	"time"
)

// ═══════════════════════════════════════════════════════════════════
// TA_STD_TH Schema Models
// ═══════════════════════════════════════════════════════════════════

// TAFNDFundInfo represents the TA_FND_FundInfo table.
type TAFNDFundInfo struct {
	SysCoID        string    `gorm:"column:SysCoID;type:varchar(20);primaryKey;not null"`
	PrtFundCode    string    `gorm:"column:PrtFundCode;type:varchar(10);primaryKey;not null"`
	UniCode        string    `gorm:"column:UniCode;type:varchar(30);not null;default:''"`
	FundInShName   string    `gorm:"column:FundInShName;type:varchar(100);not null;default:''"`
	FundMName      string    `gorm:"column:FundMName;type:varchar(200);not null;default:''"`
	FundShMName    string    `gorm:"column:FundShMName;type:varchar(150);not null;default:''"`
	FundSName      string    `gorm:"column:FundSName;type:varchar(200);not null;default:''"`
	FundShSName    string    `gorm:"column:FundShSName;type:varchar(150);not null;default:''"`
	FundRiskLevel  string    `gorm:"column:FundRiskLevel;type:varchar(6);not null;default:''"`
	IssueBaseCry   string    `gorm:"column:IssueBaseCry;type:varchar(3);not null;default:''"`
	FundSetupDate  time.Time `gorm:"column:FundSetupDate;type:timestamp with time zone;not null;default:'1900-01-01 00:00:00+08'"`
	GIINNo         string    `gorm:"column:GIINNo;type:varchar(20);not null;default:''"`
	FundStatus     string    `gorm:"column:FundStatus;type:varchar(6);not null;default:''"`
	FundWarningMsg string    `gorm:"column:FundWarningMsg;type:varchar(200);not null;default:''"`
	UnitDotLen     int16     `gorm:"column:UnitDotLen;type:smallint;not null;default:0"`
	FundSubsWay    string    `gorm:"column:FundSubsWay;type:varchar(6);not null;default:''"`
	FundRdmWay     string    `gorm:"column:FundRdmWay;type:varchar(6);not null;default:''"`
	RdmFNavDay     int16     `gorm:"column:RdmFNavDay;type:smallint;not null;default:0"`
	RdmFNavWay     string    `gorm:"column:RdmFNavWay;type:varchar(6);not null;default:''"`
	RdmFNavDRate   float64   `gorm:"column:RdmFNavDRate;type:numeric(5,2);not null;default:0"`
	RdmFNavDVal    float64   `gorm:"column:RdmFNavDVal;type:numeric(20,6);not null;default:0"`
	IsShortFee     string    `gorm:"column:IsShortFee;type:varchar(6);not null;default:''"`
	IsAntiDilFee   string    `gorm:"column:IsAntiDilFee;type:varchar(1);not null;default:''"`
	IsETF          string    `gorm:"column:IsETF;type:varchar(1);not null;default:''"`
	ETFCode        string    `gorm:"column:ETFCode;type:varchar(10);not null;default:''"`
	IsDIM          string    `gorm:"column:IsDIM;type:varchar(1);not null;default:''"`
	ShoreID        string    `gorm:"column:ShoreID;type:varchar(6);not null;default:''"`
	OFDFHCode      string    `gorm:"column:OFDFHCode;type:varchar(10);not null;default:''"`
	ValidFrom      time.Time `gorm:"column:ValidFrom;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	ValidTo        time.Time `gorm:"column:ValidTo;type:timestamp;not null;default:'9999-12-31 23:59:59.99'"`
	DataID         string    `gorm:"column:DataID;type:uuid;not null;default:gen_random_uuid()"`
	CreateID       string    `gorm:"column:CreateID;type:varchar(30);not null;default:''"`
	CreateDate     time.Time `gorm:"column:CreateDate;type:timestamp with time zone;not null;default:'1900-01-01 00:00:00+08'"`
	Inspect        int16     `gorm:"column:Inspect;type:smallint;not null;default:1"`
	FlowID         int64     `gorm:"column:FlowID;type:bigint;not null;default:0"`
	InspectID      string    `gorm:"column:InspectID;type:varchar(30);not null;default:''"`
	InspectDate    time.Time `gorm:"column:InspectDate;type:timestamp with time zone;not null;default:'1900-01-01 00:00:00+08'"`
	UpdateID       string    `gorm:"column:UpdateID;type:varchar(30);not null;default:''"`
	UpdateDate     time.Time `gorm:"column:UpdateDate;type:timestamp with time zone;not null;default:'1900-01-01 00:00:00+08'"`
	DataFlag       []byte    `gorm:"column:DataFlag;type:bytea"`
	DiffColumns    string    `gorm:"column:DiffColumns;type:text;not null;default:''"`
}

// TableName overrides the table name used by GORM.
func (TAFNDFundInfo) TableName() string {
	return "TA_STD_TH.TA_FND_FundInfo"
}

// TAFNDFundInfoEdit represents the TA_FND_FundInfo_Edit table.
type TAFNDFundInfoEdit struct {
	TAFNDFundInfo
}

// TableName overrides the table name used by GORM.
func (TAFNDFundInfoEdit) TableName() string {
	return "TA_STD_TH.TA_FND_FundInfo_Edit"
}
