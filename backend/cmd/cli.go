package cmd

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"tpt-titan/backend/config"
	"tpt-titan/backend/models"
	"tpt-titan/backend/services"
	"tpt-titan/backend/utils"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Run parses backend CLI subcommands. It returns true if a subcommand was
// handled (in which case the caller should exit), or false to start the server.
func Run(args []string) bool {
	if len(args) == 0 {
		return false
	}

	sub := args[0]
	rest := args[1:]

	switch sub {
	case "admin":
		return runAdmin(rest)
	case "seed":
		return runSeed(rest)
	case "backup":
		return runBackup(rest)
	case "restore":
		return runRestore(rest)
	case "migrate":
		return runMigrate(rest)
	case "help", "--help", "-h":
		printHelp()
		return true
	default:
		// Unknown subcommand: let the server handle it (flag parsing in main).
		return false
	}
}

func printHelp() {
	fmt.Print(`TPT Titan — management CLI

Usage:
  tpt-titan [serve]            Start the API server (default)
  tpt-titan admin <command>    Admin user/role management
  tpt-titan seed [--force]     Insert demo/seed data for evaluation
  tpt-titan backup [name]      Create a full database backup
  tpt-titan restore <id>       Restore a backup by ID
  tpt-titan migrate [status]   Apply pending schema migrations (idempotent)

Admin commands:
  admin create-user --email E --username U --password P [--admin] [--first F] [--last L]
  admin list-users
  admin set-admin --username U [--admin=true|false]

Environment:
  Reads configuration from .env (see .env.example). DB_TYPE/DB_PATH or the
  PostgreSQL settings select the target database.
`)
}

func openDB() (*gorm.DB, *config.Config) {
	cfg := config.Load()
	if err := config.ConnectDatabase(&cfg.Database); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	db := config.GetDatabase()
	return db, cfg
}

func hashPassword(password string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}
	return string(hashed)
}

func persistEncryptionSalt(db *gorm.DB, userID uuid.UUID, password string) {
	km, err := utils.NewKeyManager(password)
	if err != nil {
		log.Printf("Warning: could not initialize encryption key manager: %v", err)
		return
	}
	salt := base64.StdEncoding.EncodeToString(km.GetSalt())
	if err := db.Model(&models.User{}).Where("id = ?", userID).
		Update("encryption_salt", salt).Error; err != nil {
		log.Printf("Warning: failed to persist encryption salt: %v", err)
	}
}

func userExists(db *gorm.DB, email, username string) bool {
	var count int64
	db.Model(&models.User{}).Where("email = ? OR username = ?", email, username).Count(&count)
	return count > 0
}

func runAdmin(args []string) bool {
	if len(args) == 0 {
		fmt.Println("usage: tpt-titan admin <create-user|list-users|set-admin> [flags]")
		return true
	}

	db, _ := openDB()
	defer config.CloseDatabase()

	switch args[0] {
	case "create-user":
		return adminCreateUser(db, args[1:])
	case "list-users":
		return adminListUsers(db)
	case "set-admin":
		return adminSetAdmin(db, args[1:])
	default:
		fmt.Printf("unknown admin command: %s\n", args[0])
		return true
	}
}

func adminCreateUser(db *gorm.DB, args []string) bool {
	var email, username, password, first, last string
	isAdmin := false
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--email":
			email = args[i+1]
			i++
		case "--username":
			username = args[i+1]
			i++
		case "--password":
			password = args[i+1]
			i++
		case "--first":
			first = args[i+1]
			i++
		case "--last":
			last = args[i+1]
			i++
		case "--admin":
			isAdmin = true
		}
	}

	if email == "" || username == "" || password == "" {
		fmt.Println("error: --email, --username and --password are required")
		return true
	}
	if len(password) < 8 {
		fmt.Println("error: password must be at least 8 characters")
		return true
	}
	if userExists(db, email, username) {
		fmt.Printf("error: a user with email %q or username %q already exists\n", email, username)
		return true
	}

	user := models.User{
		Email:        email,
		Username:     username,
		PasswordHash: hashPassword(password),
		FirstName:    first,
		LastName:     last,
		IsActive:     true,
		IsAdmin:      isAdmin,
		IsVerified:   true,
	}
	if err := db.Create(&user).Error; err != nil {
		fmt.Printf("error: failed to create user: %v\n", err)
		return true
	}
	persistEncryptionSalt(db, user.ID, password)

	fmt.Printf("Created user %q (id=%s, admin=%v)\n", username, user.ID, isAdmin)
	return true
}

func adminListUsers(db *gorm.DB) bool {
	var users []models.User
	if err := db.Order("created_at asc").Find(&users).Error; err != nil {
		fmt.Printf("error: failed to list users: %v\n", err)
		return true
	}
	if len(users) == 0 {
		fmt.Println("no users found")
		return true
	}
	fmt.Printf("%-38s %-24s %-32s admin\n", "ID", "USERNAME", "EMAIL")
	for _, u := range users {
		fmt.Printf("%-38s %-24s %-32s %v\n", u.ID, u.Username, u.Email, u.IsAdmin)
	}
	return true
}

func adminSetAdmin(db *gorm.DB, args []string) bool {
	var username string
	isAdmin := true
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--username":
			username = args[i+1]
			i++
		case "--admin":
			v := strings.ToLower(args[i+1])
			isAdmin = v != "false" && v != "0" && v != "no"
			i++
		}
	}
	if username == "" {
		fmt.Println("error: --username is required")
		return true
	}
	res := db.Model(&models.User{}).Where("username = ?", username).Update("is_admin", isAdmin)
	if res.Error != nil {
		fmt.Printf("error: %v\n", res.Error)
		return true
	}
	if res.RowsAffected == 0 {
		fmt.Printf("error: no user found with username %q\n", username)
		return true
	}
	fmt.Printf("Updated admin flag for %q -> %v\n", username, isAdmin)
	return true
}

func runSeed(args []string) bool {
	force := false
	for _, a := range args {
		if a == "--force" {
			force = true
		}
	}

	db, _ := openDB()
	defer config.CloseDatabase()

	if !force && userExists(db, "demo@tpt-titan.local", "demo") {
		fmt.Println("Seed data already present (demo user exists). Pass --force to re-seed.")
		return true
	}

	seedData(db, force)
	fmt.Println("Seed data created. Login with username 'demo' / password 'demo12345'.")
	return true
}

func runBackup(args []string) bool {
	name := "manual-backup"
	if len(args) > 0 {
		name = args[0]
	}

	db, cfg := openDB()
	defer config.CloseDatabase()

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get SQL DB: %v", err)
	}

	backupRoot := filepath.Join(filepath.Dir(cfg.Database.Path), "backups")
	if err := os.MkdirAll(backupRoot, 0755); err != nil {
		log.Fatalf("Failed to create backup directory: %v", err)
	}

	bs := services.NewBackupService(sqlDB, backupRoot)
	meta, err := bs.CreateFullBackup(name, "Created via management CLI")
	if err != nil {
		fmt.Printf("error: backup failed: %v\n", err)
		return true
	}
	fmt.Printf("Backup created: %s (id=%s, records=%d)\n", meta.Name, meta.ID, meta.RecordCount)
	return true
}

func runRestore(args []string) bool {
	if len(args) == 0 {
		fmt.Println("usage: tpt-titan restore <backup-id>")
		return true
	}
	backupID, err := uuid.Parse(args[0])
	if err != nil {
		fmt.Printf("error: invalid backup id %q\n", args[0])
		return true
	}

	db, cfg := openDB()
	defer config.CloseDatabase()

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get SQL DB: %v", err)
	}

	backupRoot := filepath.Join(filepath.Dir(cfg.Database.Path), "backups")
	bs := services.NewBackupService(sqlDB, backupRoot)
	if err := bs.RestoreBackup(backupID, nil); err != nil {
		fmt.Printf("error: restore failed: %v\n", err)
		return true
	}
	fmt.Printf("Restore completed for backup %s\n", backupID)
	return true
}
