package auth

import "time"

// SysUser maps to the sys_users table.
type SysUser struct {
	ID        uint      `gorm:"primaryKey"`
	Username  string    `gorm:"uniqueIndex;size:50;not null"`
	Password  string    `gorm:"size:255;not null"` // bcrypt hash
	Role      string    `gorm:"size:50;not null"`
	IsActive  bool      `gorm:"default:true"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	CreatedBy string    `gorm:"size:50;default:system"`
	UpdatedBy string    `gorm:"size:50;default:system"`
}

// TableName overrides the default table name.
func (SysUser) TableName() string { return "sys_users" }

// SysLoginLog maps to the sys_login_log table.
type SysLoginLog struct {
	ID        uint      `gorm:"primaryKey"`
	Username  string    `gorm:"size:50;not null"`
	IPAddress string    `gorm:"size:45"`
	UserAgent string    `gorm:"size:255"`
	Success   bool      `gorm:"not null"`
	Message   string    `gorm:"size:255"`
	LoggedAt  time.Time `gorm:"autoCreateTime"`
}

// TableName overrides the default table name.
func (SysLoginLog) TableName() string { return "sys_login_log" }
