package auth

import (
	"log/slog"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// pocUser defines a hardcoded POC user for seeding.
type pocUser struct {
	Username string
	Password string
	Role     string
}

// pocUsers matches the hardcoded users from api-gateway/internal/http/handlers/auth.go.
var pocUsers = []pocUser{
	{Username: "admin", Password: "admin123", Role: "admin"},
	{Username: "user", Password: "user123", Role: "user"},
	{Username: "manager1", Password: "manager123", Role: "manager"},
	{Username: "manager2", Password: "manager234", Role: "manager"},
	{Username: "maker1", Password: "maker123", Role: "maker"},
	{Username: "maker2", Password: "maker234", Role: "maker"},
	{Username: "checker1", Password: "checker123", Role: "checker"},
	{Username: "checker2", Password: "checker234", Role: "checker"},
	{Username: "viewer1", Password: "viewer123", Role: "viewer"},
	{Username: "viewer2", Password: "viewer234", Role: "viewer"},
	{Username: "auditor1", Password: "auditor123", Role: "auditor"},
	{Username: "auditor2", Password: "auditor234", Role: "auditor"},
}

// AutoMigrateAndSeed creates the sys_users and sys_login_log tables (if they don't exist)
// and seeds the POC users with bcrypt-hashed passwords.
// This is safe to call multiple times — existing users are skipped.
func AutoMigrateAndSeed(db *gorm.DB, log *slog.Logger) error {
	// Auto-migrate tables
	if err := db.AutoMigrate(&SysUser{}, &SysLoginLog{}); err != nil {
		log.Error("failed to auto-migrate auth tables", slog.String("error", err.Error()))
		return err
	}
	log.Info("✅ Auth tables migrated (sys_users, sys_login_log)")

	// Seed POC users (skip if already exist)
	seeded := 0
	for _, u := range pocUsers {
		var count int64
		db.Model(&SysUser{}).Where("username = ?", u.Username).Count(&count)
		if count > 0 {
			continue
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Error("failed to hash password", slog.String("user", u.Username), slog.String("error", err.Error()))
			continue
		}

		user := SysUser{
			Username:  u.Username,
			Password:  string(hash),
			Role:      u.Role,
			IsActive:  true,
			CreatedBy: "system",
			UpdatedBy: "system",
		}
		if err := db.Create(&user).Error; err != nil {
			log.Error("failed to seed user", slog.String("user", u.Username), slog.String("error", err.Error()))
			continue
		}
		seeded++
	}

	if seeded > 0 {
		log.Info("✅ Seeded POC users", slog.Int("count", seeded))
	} else {
		log.Info("ℹ️  POC users already exist — skipped seeding")
	}

	return nil
}
