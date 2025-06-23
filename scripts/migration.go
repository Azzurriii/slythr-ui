package main

import (
	"github.com/Azzurriii/slythr/pkg/logger"

	config "github.com/Azzurriii/slythr/config"
	"github.com/Azzurriii/slythr/internal/domain/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Migrate() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatalf("Failed to load config: %v", err)
		return err
	}

	dsn := cfg.Database.DSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
		return err
	}

	if err := db.AutoMigrate(
		&entities.Contract{},
		&entities.StaticAnalysis{},
		&entities.DynamicAnalysis{},
		&entities.GeneratedTestCases{},
	); err != nil {
		logger.Fatalf("Migration failed: %v", err)
		return err
	}

	logger.Info("Migration completed successfully!")
	return nil
}
