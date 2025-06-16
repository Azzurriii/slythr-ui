package main

import (
	"log"

	config "github.com/Azzurriii/slythr-go-backend/config"
	"github.com/Azzurriii/slythr-go-backend/internal/domain/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Migrate function to run database migrations
func Migrate() error {
	// Load from environment variables or configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
		return err
	}

	dsn := cfg.Database.DSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return err
	}

	// AutoMigrate runs the migration for all entities
	if err := db.AutoMigrate(
		&entities.Contract{},
		&entities.StaticAnalysis{},
		&entities.DynamicAnalysis{},
	); err != nil {
		log.Fatalf("Migration failed: %v", err)
		return err
	}

	log.Println("Migration completed successfully!")
	return nil
}
