package config

import (
	"fmt"
	"log"
	"tpt-titan/backend/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is the global database connection
var DB *gorm.DB

// ConnectDatabase initializes the database connection
func ConnectDatabase(cfg *DatabaseConfig) error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test the connection
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connected successfully")

	// Store the connection globally
	DB = db

	// Auto migrate the schema
	if err := migrateDatabase(db); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	return nil
}

// migrateDatabase runs auto-migration for all models
func migrateDatabase(db *gorm.DB) error {
	log.Println("Running database migrations...")

	// Add all models here
	if err := db.AutoMigrate(
		&models.User{},
		// Add other models as they are created
	); err != nil {
		return fmt.Errorf("failed to migrate: %w", err)
	}

	log.Println("Database migrations completed")
	return nil
}

// CloseDatabase closes the database connection
func CloseDatabase() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
