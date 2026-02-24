package db

import (
	"fmt"
	"log/slog"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TA_DEV_TH Database Connection
type Database struct {
	DB  *gorm.DB
	Log *slog.Logger
}

// NewConnection initializes the PostgreSQL connection to TA_DEV_TH.
func NewConnection(dsn string, logger *slog.Logger) (*Database, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	logger.Info("Successfully connected to PostgreSQL database TA_DEV_TH")
	return &Database{DB: db, Log: logger}, nil
}

// AutoMigrate setups the schemas and migrating the tables.
func (d *Database) AutoMigrate() error {
	d.Log.Info("Starting schema creation and auto-migration")

	// 1. Ensure schema exists
	schema := "TA_STD_TH"
	query := fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %q;", schema)
	if err := d.DB.Exec(query).Error; err != nil {
		return fmt.Errorf("failed to create schema %s: %w", schema, err)
	}
	d.Log.Info(fmt.Sprintf("Schema %q verified", schema))

	// 2. AutoMigrate Tables
	err := d.DB.AutoMigrate(
		&DTAFNDSwitch{},
		&DTAFNDSwitchEdit{},
		&DTAFNDSwitchFund{},
		&DTAFNDSwitchFundEdit{},
		&DTAFNDFundFeeSwitch{},
		&DTAFNDFundFeeSwitchEdit{},
		&DTAFNDSwitchCry{},
		&DTAFNDSwitchCryEdit{},
	)
	if err != nil {
		return fmt.Errorf("gorm auto-migrate failed: %w", err)
	}

	d.Log.Info("Auto-migration completed successfully")
	return nil
}
