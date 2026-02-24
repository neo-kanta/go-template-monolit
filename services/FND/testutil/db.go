package testutil

import (
	"log/slog"
	"os"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SetupTestDB initializes a connection to a test PostgreSQL database.
// It skips the test if TEST_DB_DSN is not provided in the environment.
// Example DSN: postgres://user:password@localhost:5432/dbname?sslmode=disable
func SetupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		t.Skip("Skipping integration test; TEST_DB_DSN environment variable not set.")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Ensure the target schema exists
	if err := db.Exec(`CREATE SCHEMA IF NOT EXISTS "TA_STD_TH"`).Error; err != nil {
		t.Fatalf("Failed to create TA_STD_TH schema: %v", err)
	}

	return db
}

// NewTestLogger creates a simple structured logger for testing.
func NewTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
}
