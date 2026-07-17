package cmd

import (
	"fmt"
	"log"
	"time"

	"tpt-titan/backend/config"
	"tpt-titan/backend/models"

	"gorm.io/gorm"
)

// migration represents a single versioned schema change. The first migration
// (0001) runs the existing GORM AutoMigrate baseline so the runner can record
// everything the app already relies on. Subsequent migrations hold explicit
// up-SQL for breaking schema changes and run only once.
type migration struct {
	Version string
	Name    string
	Up      func(db *gorm.DB) error
}

// migrations must be appended in ascending version order. Never reorder.
var migrations = []migration{
	{
		Version: "0001",
		Name:    "gorm_automigrate_baseline",
		Up: func(db *gorm.DB) error {
			return db.AutoMigrate(
				&models.User{},
				&models.EncryptedDocument{},
				&models.EncryptedForm{},
				&models.EncryptedFormResponse{},
				&models.EncryptedEmail{},
				&models.EncryptedTask{},
				&models.KeyBackup{},
				&models.RecoveryShare{},
				&models.HardwareKey{},
				&models.RecoveryAttempt{},
				&models.AIModel{},
				&models.AITask{},
				&models.AIRequest{},
				&models.AIUsage{},
				&models.AIUpgradeCheck{},
				&models.ChatRoom{},
				&models.ChatParticipant{},
				&models.ChatMessage{},
				&models.MessageReaction{},
				&models.UserStatus{},
				&models.VoiceNote{},
				&models.VoiceAnnotation{},
				&models.Calendar{},
				&models.Event{},
				&models.EventAttendee{},
				&models.Contact{},
				&models.EmailAccount{},
				&models.Email{},
				&models.SyncDevice{},
				&models.SyncFolder{},
				&models.FileVersion{},
				&models.Task{},
				&models.TaskSubtask{},
				&models.Project{},
				&models.TaskTag{},
				&models.Workflow{},
				&models.WorkflowNode{},
				&models.WorkflowConnection{},
				&models.WorkflowExecution{},
				&models.WorkflowTemplate{},
				&models.SystemSetting{},
				&models.SpreadsheetChart{},
			)
		},
	},
	{
		Version: "0002",
		Name:    "add_workflow_execution_is_dry_run",
		Up: func(db *gorm.DB) error {
			return db.AutoMigrate(&models.WorkflowExecution{})
		},
	},
	{
		Version: "0003",
		Name:    "seed_default_workflow_templates",
		Up: func(db *gorm.DB) error {
			return config.EnsureDefaultWorkflowTemplates(db)
		},
	},
}

// runMigrate applies all pending migrations idempotently.
func runMigrate(args []string) bool {
	db, _ := openDB()
	defer config.CloseDatabase()

	if err := ensureMigrationsTable(db); err != nil {
		log.Fatalf("Failed to initialize migrations table: %v", err)
	}

	applied, err := appliedVersions(db)
	if err != nil {
		log.Fatalf("Failed to read applied migrations: %v", err)
	}

	// Support `migrate status` to just report.
	if len(args) > 0 && args[0] == "status" {
		printMigrationStatus(applied)
		return true
	}

	pending := 0
	for _, m := range migrations {
		if _, done := applied[m.Version]; done {
			continue
		}
		log.Printf("Applying migration %s (%s)...", m.Version, m.Name)
		if err := m.Up(db); err != nil {
			log.Fatalf("Migration %s failed: %v", m.Version, err)
		}
		if err := recordMigration(db, m); err != nil {
			log.Fatalf("Failed recording migration %s: %v", m.Version, err)
		}
		fmt.Printf("Applied %s (%s)\n", m.Version, m.Name)
		pending++
	}

	if pending == 0 {
		fmt.Println("Database schema is up to date.")
	} else {
		fmt.Printf("Applied %d migration(s).\n", pending)
	}
	return true
}

func ensureMigrationsTable(db *gorm.DB) error {
	return db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
		version   VARCHAR(64) PRIMARY KEY,
		name      VARCHAR(255) NOT NULL,
		applied_at TIMESTAMP NOT NULL
	)`).Error
}

func appliedVersions(db *gorm.DB) (map[string]bool, error) {
	rows, err := db.Raw("SELECT version FROM schema_migrations").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := map[string]bool{}
	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		applied[v] = true
	}
	return applied, rows.Err()
}

func recordMigration(db *gorm.DB, m migration) error {
	return db.Exec(
		"INSERT INTO schema_migrations (version, name, applied_at) VALUES (?, ?, ?)",
		m.Version, m.Name, time.Now(),
	).Error
}

func printMigrationStatus(applied map[string]bool) {
	fmt.Println("Migration status:")
	for _, m := range migrations {
		state := "pending"
		if applied[m.Version] {
			state = "applied"
		}
		fmt.Printf("  %s  %-32s %s\n", m.Version, m.Name, state)
	}
}
