package cmd

import (
	"log"
	"time"

	"tpt-titan/backend/models"

	"gorm.io/gorm"
)

// seedData inserts demo content so evaluators can explore the product without
// manual setup. It is idempotent per-user (guarded by runSeed's userExists check
// unless --force is passed).
func seedData(db *gorm.DB, force bool) {
	demoPassword := "demo12345"

	demo := models.User{
		Email:        "demo@tpt-titan.local",
		Username:     "demo",
		PasswordHash: hashPassword(demoPassword),
		FirstName:    "Demo",
		LastName:     "User",
		IsActive:     true,
		IsAdmin:      true,
		IsVerified:   true,
	}
	if err := db.Create(&demo).Error; err != nil {
		log.Printf("seed: failed to create demo user: %v", err)
		return
	}
	persistEncryptionSalt(db, demo.ID, demoPassword)

	// A second, non-admin account to demonstrate multi-user features.
	member := models.User{
		Email:        "member@tpt-titan.local",
		Username:     "member",
		PasswordHash: hashPassword(demoPassword),
		FirstName:    "Team",
		LastName:     "Member",
		IsActive:     true,
		IsAdmin:      false,
		IsVerified:   true,
	}
	if err := db.Create(&member).Error; err != nil {
		log.Printf("seed: failed to create member user: %v", err)
	} else {
		persistEncryptionSalt(db, member.ID, demoPassword)
	}

	now := time.Now()

	// Sample project + tasks for the task board.
	project := models.Project{
		UserID: demo.ID,
		Name:   "Getting Started",
		Color:  "#4f46e5",
	}
	if err := db.Create(&project).Error; err != nil {
		log.Printf("seed: failed to create project: %v", err)
	} else {
		tasks := []models.Task{
			{
				UserID:      demo.ID,
				ProjectID:   &project.ID,
				Title:       "Welcome to TPT Titan",
				Description: "Explore the sidebar to try Tasks, Forms, Spreadsheets, Documents and more.",
				Status:      "todo",
				Priority:    "high",
				DueDate:     &now,
			},
			{
				UserID:      demo.ID,
				ProjectID:   &project.ID,
				Title:       "Create your first form",
				Description: "Use the Forms page to build a survey and collect responses.",
				Status:      "in_progress",
				Priority:    "medium",
			},
			{
				UserID:      demo.ID,
				ProjectID:   &project.ID,
				Title:       "Configure AI settings",
				Description: "Point the AI panel at a local Ollama instance or add an OpenRouter key.",
				Status:      "done",
				Priority:    "low",
			},
		}
		if err := db.Create(&tasks).Error; err != nil {
			log.Printf("seed: failed to create tasks: %v", err)
		}
	}

	// Ensure sane system defaults exist.
	ensureSetting(db, "registration_open", "true")
	ensureSetting(db, "default_theme", "\"system\"")
	ensureSetting(db, "seeded_at", time.Now().Format(time.RFC3339))

	log.Println("seed: demo data inserted")
}

func ensureSetting(db *gorm.DB, key, value string) {
	var count int64
	db.Model(&models.SystemSetting{}).Where("key = ?", key).Count(&count)
	if count > 0 {
		return
	}
	if err := db.Create(&models.SystemSetting{Key: key, Value: value}).Error; err != nil {
		log.Printf("seed: failed to set %s: %v", key, err)
	}
}
