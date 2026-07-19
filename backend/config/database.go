package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"tpt-titan/backend/models"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is the global database connection
var DB *gorm.DB

// ConnectDatabase initializes the database connection
func ConnectDatabase(cfg *DatabaseConfig) error {
	var db *gorm.DB
	var err error

	switch cfg.Type {
	case "sqlite":
		// Ensure the directory exists for SQLite database
		dir := filepath.Dir(cfg.Path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create database directory: %w", err)
		}

		db, err = gorm.Open(sqlite.Open(cfg.Path), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			return fmt.Errorf("failed to connect to SQLite database: %w", err)
		}

		log.Printf("SQLite database connected successfully at %s", cfg.Path)

	case "postgres":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode)

		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			return fmt.Errorf("failed to connect to PostgreSQL database: %w", err)
		}

		log.Println("PostgreSQL database connected successfully")

	default:
		return fmt.Errorf("unsupported database type: %s (supported: sqlite, postgres)", cfg.Type)
	}

	// Test the connection
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

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
		// Core auth
		&models.User{},
		// Encrypted document/form/task storage
		&models.EncryptedDocument{},
		&models.EncryptedForm{},
		&models.EncryptedFormResponse{},
		&models.EncryptedEmail{},
		&models.EncryptedTask{},
		// Encryption key management
		&models.KeyBackup{},
		&models.RecoveryShare{},
		&models.HardwareKey{},
		&models.RecoveryAttempt{},
		// AI models and usage
		&models.AIModel{},
		&models.AITask{},
		&models.AIRequest{},
		&models.AIUsage{},
		&models.AIUpgradeCheck{},
		// Chat
		&models.ChatRoom{},
		&models.ChatParticipant{},
		&models.ChatMessage{},
		&models.MessageReaction{},
		&models.UserStatus{},
		// Voice notes
		&models.VoiceNote{},
		&models.VoiceAnnotation{},
		// Calendar
		&models.Calendar{},
		&models.Event{},
		&models.EventAttendee{},
		// Contacts
		&models.Contact{},
		// Email
		&models.EmailAccount{},
		&models.Email{},
		// File sync
		&models.SyncDevice{},
		&models.SyncFolder{},
		&models.FileVersion{},
		// Tasks & projects
		&models.Task{},
		&models.TaskSubtask{},
		&models.Project{},
		&models.TaskTag{},
		// Workflows
		&models.Workflow{},
		&models.WorkflowNode{},
		&models.WorkflowConnection{},
		&models.WorkflowExecution{},
		&models.WorkflowTemplate{},
		// System settings
		&models.SystemSetting{},
		// MCP servers (Model Context Protocol bridges)
		&models.MCPServer{},
		// Webhook delivery log (inbound/outbound call audit trail)
		&models.WebhookDeliveryLog{},
	); err != nil {
		return fmt.Errorf("failed to migrate: %w", err)
	}

	// Initialize system AI models if they don't exist
	if err := initializeSystemAIModels(db); err != nil {
		log.Printf("Warning: Failed to initialize system AI models: %v", err)
	}

	// Initialize built-in preset workflow templates if they don't exist
	if err := EnsureDefaultWorkflowTemplates(db); err != nil {
		log.Printf("Warning: Failed to initialize default workflow templates: %v", err)
	}

	log.Println("Database migrations completed")
	return nil
}

// GetDatabase returns the current database instance
func GetDatabase() *gorm.DB {
	return DB
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

// initializeSystemAIModels sets up the default AI models
func initializeSystemAIModels(db *gorm.DB) error {
	for _, model := range models.DefaultAIModels {
		// Check if model already exists
		var count int64
		db.Model(&models.AIModel{}).Where("model_id = ? AND is_system = ?", model.ModelID, true).Count(&count)

		if count == 0 {
			if err := db.Create(&model).Error; err != nil {
				log.Printf("Failed to create system model %s: %v", model.Name, err)
				continue
			}
			log.Printf("Created system AI model: %s", model.Name)
		}
	}
	return nil
}

// EnsureDefaultWorkflowTemplates seeds the built-in preset workflow templates
// if they don't already exist, following the same idempotent pattern as
// initializeSystemAIModels above.
func EnsureDefaultWorkflowTemplates(db *gorm.DB) error {
	for _, tmpl := range models.DefaultWorkflowTemplates {
		var count int64
		db.Model(&models.WorkflowTemplate{}).Where("name = ? AND is_system = ?", tmpl.Name, true).Count(&count)

		if count == 0 {
			if err := db.Create(&tmpl).Error; err != nil {
				log.Printf("Failed to create workflow template %s: %v", tmpl.Name, err)
				continue
			}
			log.Printf("Created preset workflow template: %s", tmpl.Name)
		}
	}
	return nil
}
