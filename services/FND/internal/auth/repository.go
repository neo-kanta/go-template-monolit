package auth

import (
	"log/slog"

	"gorm.io/gorm"
)

// UserRepository provides access to the sys_users table.
type UserRepository struct {
	db  *gorm.DB
	log *slog.Logger
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository(db *gorm.DB, log *slog.Logger) *UserRepository {
	return &UserRepository{db: db, log: log}
}

// FindByUsername looks up an active user by username.
// Returns nil if the user is not found or is inactive.
func (r *UserRepository) FindByUsername(username string) (*SysUser, error) {
	var user SysUser
	err := r.db.Where("username = ? AND is_active = true", username).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.log.Error("failed to query sys_users", slog.String("error", err.Error()))
		return nil, err
	}
	return &user, nil
}

// LoginLogRepository provides access to the sys_login_log table.
type LoginLogRepository struct {
	db  *gorm.DB
	log *slog.Logger
}

// NewLoginLogRepository creates a new LoginLogRepository.
func NewLoginLogRepository(db *gorm.DB, log *slog.Logger) *LoginLogRepository {
	return &LoginLogRepository{db: db, log: log}
}

// Record writes a login attempt to the sys_login_log table.
func (r *LoginLogRepository) Record(entry *SysLoginLog) {
	if err := r.db.Create(entry).Error; err != nil {
		r.log.Error("failed to write login log",
			slog.String("username", entry.Username),
			slog.String("error", err.Error()),
		)
	}
}
