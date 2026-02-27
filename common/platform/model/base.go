// Package models provides shared GORM models and base types for the FND module.
package models

import (
	"time"

	"gorm.io/gorm"
)

// ApprovalStatus represents the 4-eyes principle status
type ApprovalStatus string

const (
	StatusPendingApproval ApprovalStatus = "PENDING_APPROVAL"
	StatusApproved        ApprovalStatus = "APPROVED"
	StatusRejected        ApprovalStatus = "REJECTED"
	StatusDraft           ApprovalStatus = "DRAFT"
)

// Base Model
var BangkokLocation *time.Location

func init() {
	var err error
	BangkokLocation, err = time.LoadLocation("Asia/Bangkok")
	if err != nil {
		// Fallback to UTC+7 if timezone data is not available
		BangkokLocation = time.FixedZone("ICT", 7*60*60)
	}
}

// BaseModel is the shared GORM base model with audit fields
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// MakerCheckerFields contains 4-eyes principle fields without an ID (for embedding in composite PK models)
type MakerCheckerFields struct {
	MakerID   string         `gorm:"column:MakerID;type:varchar(50);not null;default:''" json:"maker_id"`
	CheckerID *string        `gorm:"column:CheckerID;type:varchar(50)" json:"checker_id,omitempty"`
	Status    ApprovalStatus `gorm:"column:Status;type:varchar(20);not null;default:'DRAFT'" json:"status"`
	Remark    *string        `gorm:"column:Remark;type:text" json:"remark,omitempty"`
}

// MakerCheckerModel extends BaseModel with MakerCheckerFields
type MakerCheckerModel struct {
	BaseModel
	MakerCheckerFields
}

const THBPrecision = 4

const THBScale = 10000
