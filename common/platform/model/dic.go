package models

import (
	"time"
)

// ═══════════════════════════════════════════════════════════════════
// MWDP_TA_STD Schema Models (Shared)
// ═══════════════════════════════════════════════════════════════════

// DICOption represents the DIC_Options table.
type DICOption struct {
	ValidFrom    time.Time `gorm:"column:ValidFrom;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	ValidTo      time.Time `gorm:"column:ValidTo;type:timestamp;not null;default:'9999-12-31 23:59:59.99'"`
	DataID       string    `gorm:"column:DataID;type:uuid;not null;default:gen_random_uuid()"`
	CreateID     string    `gorm:"column:CreateID;type:varchar(30);not null;default:''"`
	CreateDate   time.Time `gorm:"column:CreateDate;type:timestamp with time zone;not null;default:'1900-01-01 00:00:00+08'"`
	Inspect      int16     `gorm:"column:Inspect;type:smallint;not null;default:1"`
	FlowID       int64     `gorm:"column:FlowID;type:bigint;not null;default:0"`
	InspectID    string    `gorm:"column:InspectID;type:varchar(30);not null;default:''"`
	InspectDate  time.Time `gorm:"column:InspectDate;type:timestamp with time zone;not null;default:'1900-01-01 00:00:00+08'"`
	UpdateID     string    `gorm:"column:UpdateID;type:varchar(30);not null;default:''"`
	UpdateDate   time.Time `gorm:"column:UpdateDate;type:timestamp with time zone;not null;default:'1900-01-01 00:00:00+08'"`
	DataFlag     []byte    `gorm:"column:DataFlag;type:bytea"` // Equivalent to SQL Server timestamp/rowversion
	LangID       string    `gorm:"column:LangID;type:varchar(15);primaryKey;not null"`
	ProductID    string    `gorm:"column:ProductID;type:varchar(30);primaryKey;not null"`
	CustomerID   string    `gorm:"column:CustomerID;type:varchar(10);primaryKey;not null"`
	OptID        string    `gorm:"column:OptID;type:varchar(6);primaryKey;not null"`
	OptName      string    `gorm:"column:OptName;type:varchar(50);not null;default:''"`
	ItemID       string    `gorm:"column:ItemID;type:varchar(50);primaryKey;not null"`
	ItemName     string    `gorm:"column:ItemName;type:varchar(200);not null;default:''"`
	DisplayOrder int16     `gorm:"column:DisplayOrder;type:smallint;not null;default:0"`
	IsDefault    bool      `gorm:"column:IsDefault;type:boolean;not null;default:false"`
	IsShowID     bool      `gorm:"column:IsShowID;type:boolean;not null;default:false"`
	IsEnabled    bool      `gorm:"column:IsEnabled;type:boolean;not null;default:false"`
	AccessType   int16     `gorm:"column:AccessType;type:smallint;not null;default:0"`
	CustomUse    string    `gorm:"column:CustomUse;type:varchar(4);not null;default:''"`
	Memo         string    `gorm:"column:Memo;type:varchar(1000);not null;default:''"`
}

// TableName overrides the table name used by GORM.
func (DICOption) TableName() string {
	return "MWDP_TA_STD.DIC_Options"
}

// DICMsgInfo represents the DIC_MsgInfo table.
type DICMsgInfo struct {
	DataID     string `gorm:"column:dataid;type:char(38);not null;default:gen_random_uuid()"`
	LangID     string `gorm:"column:LangID;type:varchar(15);primaryKey;not null"`
	ProductID  string `gorm:"column:ProductID;type:varchar(50);primaryKey;not null"`
	CustomerID string `gorm:"column:CustomerID;type:varchar(50);primaryKey;not null"`
	MsgCode    string `gorm:"column:MsgCode;type:varchar(20);primaryKey;not null"`
	MsgContent string `gorm:"column:MsgContent;type:varchar(255);not null;default:''"`
	MsgType    string `gorm:"column:MsgType;type:varchar(6);not null;default:''"`
	MsgKind    string `gorm:"column:MsgKind;type:varchar(30);not null;default:''"`
	MsgMemo    string `gorm:"column:MsgMemo;type:varchar(255);not null;default:''"`
	MsgAplyTo  string `gorm:"column:MsgAplyTo;type:varchar(200);not null;default:''"`
	AccessType int16  `gorm:"column:AccessType;type:smallint;not null;default:0"`
}

// TableName overrides the table name used by GORM.
func (DICMsgInfo) TableName() string {
	return "MWDP_TA_STD.DIC_MsgInfo"
}
