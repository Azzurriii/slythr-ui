package scripts

import (
	"log"

	config "github.com/Azzurriii/slythr-go-backend/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Define models for migration
type User struct {
	ID       uint   `gorm:"primaryKey"`
	Email    string `gorm:"uniqueIndex"`
	Password string
}

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

	// AutoMigrate runs the migration
	if err := db.AutoMigrate(&User{}); err != nil {
		log.Fatalf("Migration failed: %v", err)
		return err
	}

	log.Println("Migration completed successfully!")
	return nil
}
