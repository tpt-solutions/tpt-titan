package services

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// BackupService handles data backup and recovery operations
type BackupService struct {
	db         *sql.DB
	backupPath string
}

// BackupMetadata contains metadata about a backup
type BackupMetadata struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        string    `json:"type"` // "full", "incremental", "user_data"
	CreatedAt   time.Time `json:"created_at"`
	Size        int64     `json:"size"`
	Checksum    string    `json:"checksum"`
	UserID      *uuid.UUID `json:"user_id,omitempty"`
	Tables      []string   `json:"tables"`
	RecordCount int       `json:"record_count"`
}

// BackupConfig contains backup configuration
type BackupConfig struct {
	Path           string
	RetentionDays  int
	MaxBackups     int
	Compress       bool
	Encrypt        bool
	IncludeFiles   bool
	Schedule       string // cron expression
}

// NewBackupService creates a new backup service
func NewBackupService(db *sql.DB, backupPath string) *BackupService {
	return &BackupService{
		db:         db,
		backupPath: backupPath,
	}
}

// CreateFullBackup creates a complete database backup
func (bs *BackupService) CreateFullBackup(name, description string) (*BackupMetadata, error) {
	backupID := uuid.New()
	backupDir := filepath.Join(bs.backupPath, backupID.String())
	metadataPath := filepath.Join(backupDir, "metadata.json")
	dataPath := filepath.Join(backupDir, "data.sql")

	// Create backup directory
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Get all tables
	tables, err := bs.getAllTables()
	if err != nil {
		return nil, fmt.Errorf("failed to get tables: %w", err)
	}

	// Create metadata
	metadata := &BackupMetadata{
		ID:          backupID,
		Name:        name,
		Description: description,
		Type:        "full",
		CreatedAt:   time.Now(),
		Tables:      tables,
	}

	_ = backupID // Mark as used

	// Export data
	recordCount, err := bs.exportTables(dataPath, tables)
	if err != nil {
		return nil, fmt.Errorf("failed to export data: %w", err)
	}

	metadata.RecordCount = recordCount

	// Calculate file size
	if info, err := os.Stat(dataPath); err == nil {
		metadata.Size = info.Size()
	}

	// Calculate checksum
	checksum, err := bs.calculateFileChecksum(dataPath)
	if err != nil {
		log.Printf("Failed to calculate checksum: %v", err)
	}
	metadata.Checksum = checksum

	// Save metadata
	if err := bs.saveMetadata(metadataPath, metadata); err != nil {
		return nil, fmt.Errorf("failed to save metadata: %w", err)
	}

	// Compress backup
	if err := bs.compressBackup(backupDir, backupID.String()+".tar.gz"); err != nil {
		log.Printf("Failed to compress backup: %v", err)
	}

	log.Printf("Created full backup: %s with %d records", name, recordCount)
	return metadata, nil
}

// CreateUserBackup creates a backup for a specific user
func (bs *BackupService) CreateUserBackup(userID uuid.UUID, name, description string) (*BackupMetadata, error) {
	backupID := uuid.New()
	backupDir := filepath.Join(bs.backupPath, backupID.String())

	// Create backup directory
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Tables that contain user-specific data
	userTables := []string{
		"users", "user_sessions", "documents", "document_versions",
		"email_accounts", "emails", "chat_rooms", "chat_participants",
		"chat_messages", "message_reactions", "contacts", "calendars",
		"events", "tasks", "forms", "form_responses", "sync_devices",
		"sync_folders", "file_versions", "meetings", "meeting_participants",
		"meeting_chat_messages", "webrtc_connections",
	}

	// Create metadata
	metadata := &BackupMetadata{
		ID:          backupID,
		Name:        name,
		Description: description,
		Type:        "user_data",
		CreatedAt:   time.Now(),
		UserID:      &userID,
		Tables:      userTables,
	}

	// Export user data
	recordCount, err := bs.exportUserData(backupDir, userID, userTables)
	if err != nil {
		return nil, fmt.Errorf("failed to export user data: %w", err)
	}

	metadata.RecordCount = recordCount

	// Calculate total size
	var totalSize int64
	filepath.Walk(backupDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			totalSize += info.Size()
		}
		return nil
	})
	metadata.Size = totalSize

	log.Printf("Created user backup for user %s: %s with %d records", userID, name, recordCount)
	return metadata, nil
}

// RestoreBackup restores data from a backup
func (bs *BackupService) RestoreBackup(backupID uuid.UUID, tables []string) error {
	backupDir := filepath.Join(bs.backupPath, backupID.String())
	metadataPath := filepath.Join(backupDir, "metadata.json")

	// Load metadata
	metadata, err := bs.loadMetadata(metadataPath)
	if err != nil {
		return fmt.Errorf("failed to load backup metadata: %w", err)
	}

	// Validate backup type
	if metadata.Type != "full" {
		return fmt.Errorf("only full backups can be restored")
	}

	// Determine which tables to restore
	if len(tables) == 0 {
		tables = metadata.Tables
	}

	// Restore each table
	for _, table := range tables {
		dataPath := filepath.Join(backupDir, table+".sql")
		if _, err := os.Stat(dataPath); os.IsNotExist(err) {
			log.Printf("Table data not found for table: %s", table)
			continue
		}

		if err := bs.restoreTable(dataPath, table); err != nil {
			return fmt.Errorf("failed to restore table %s: %w", table, err)
		}
	}

	log.Printf("Restored backup: %s", metadata.Name)
	return nil
}

// RestoreUserData restores user data from a backup
func (bs *BackupService) RestoreUserData(backupID uuid.UUID, userID uuid.UUID) error {
	backupDir := filepath.Join(bs.backupPath, backupID.String())

	// Load metadata
	metadataPath := filepath.Join(backupDir, "metadata.json")
	metadata, err := bs.loadMetadata(metadataPath)
	if err != nil {
		return fmt.Errorf("failed to load backup metadata: %w", err)
	}

	// Validate user ownership
	if metadata.UserID == nil || *metadata.UserID != userID {
		return fmt.Errorf("backup does not belong to user")
	}

	// Restore user data
	for _, table := range metadata.Tables {
		dataPath := filepath.Join(backupDir, table+".json")
		if _, err := os.Stat(dataPath); os.IsNotExist(err) {
			continue
		}

		if err := bs.restoreUserTable(dataPath, table, userID); err != nil {
			return fmt.Errorf("failed to restore user table %s: %w", table, err)
		}
	}

	log.Printf("Restored user data for user %s from backup: %s", userID, metadata.Name)
	return nil
}

// ListBackups returns a list of available backups
func (bs *BackupService) ListBackups(userID *uuid.UUID) ([]BackupMetadata, error) {
	var backups []BackupMetadata

	entries, err := os.ReadDir(bs.backupPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read backup directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

			_, err := uuid.Parse(entry.Name())
		if err != nil {
			continue // Skip invalid backup directories
		}

		metadataPath := filepath.Join(bs.backupPath, entry.Name(), "metadata.json")
		metadata, err := bs.loadMetadata(metadataPath)
		if err != nil {
			log.Printf("Failed to load metadata for backup %s: %v", entry.Name(), err)
			continue
		}

		// Filter by user if specified
		if userID != nil && (metadata.UserID == nil || *metadata.UserID != *userID) {
			continue
		}

		backups = append(backups, *metadata)
	}

	return backups, nil
}

// DeleteBackup deletes a backup
func (bs *BackupService) DeleteBackup(backupID uuid.UUID, userID *uuid.UUID) error {
	backupDir := filepath.Join(bs.backupPath, backupID.String())

	// Verify ownership if user is specified
	if userID != nil {
		metadataPath := filepath.Join(backupDir, "metadata.json")
		metadata, err := bs.loadMetadata(metadataPath)
		if err != nil {
			return fmt.Errorf("failed to load backup metadata: %w", err)
		}

		if metadata.UserID == nil || *metadata.UserID != *userID {
			return fmt.Errorf("backup does not belong to user")
		}
	}

	// Delete backup directory
	if err := os.RemoveAll(backupDir); err != nil {
		return fmt.Errorf("failed to delete backup: %w", err)
	}

	log.Printf("Deleted backup: %s", backupID)
	return nil
}

// CleanupOldBackups removes old backups based on retention policy
func (bs *BackupService) CleanupOldBackups(retentionDays int) error {
	entries, err := os.ReadDir(bs.backupPath)
	if err != nil {
		return fmt.Errorf("failed to read backup directory: %w", err)
	}

	cutoff := time.Now().AddDate(0, 0, -retentionDays)

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		metadataPath := filepath.Join(bs.backupPath, entry.Name(), "metadata.json")
		metadata, err := bs.loadMetadata(metadataPath)
		if err != nil {
			log.Printf("Failed to load metadata for backup %s: %v", entry.Name(), err)
			continue
		}

		if metadata.CreatedAt.Before(cutoff) {
			backupDir := filepath.Join(bs.backupPath, entry.Name())
			if err := os.RemoveAll(backupDir); err != nil {
				log.Printf("Failed to delete old backup %s: %v", entry.Name(), err)
			} else {
				log.Printf("Deleted old backup: %s", metadata.Name)
			}
		}
	}

	return nil
}

// getAllTables returns all tables in the database
func (bs *BackupService) getAllTables() ([]string, error) {
	rows, err := bs.db.Query(`
		SELECT tablename FROM pg_tables
		WHERE schemaname = 'public'
		ORDER BY tablename
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}

	return tables, nil
}

// exportTables exports data from specified tables
func (bs *BackupService) exportTables(filePath string, tables []string) (int, error) {
	file, err := os.Create(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	totalRecords := 0

	for _, table := range tables {
		// Get column names
		columns, err := bs.getTableColumns(table)
		if err != nil {
			log.Printf("Failed to get columns for table %s: %v", table, err)
			continue
		}

		if len(columns) == 0 {
			continue
		}

		// Export table data
		count, err := bs.exportTableData(file, table, columns)
		if err != nil {
			log.Printf("Failed to export table %s: %v", table, err)
			continue
		}

		totalRecords += count
	}

	return totalRecords, nil
}

// exportUserData exports data for a specific user
func (bs *BackupService) exportUserData(backupDir string, userID uuid.UUID, tables []string) (int, error) {
	totalRecords := 0

	for _, table := range tables {
		data, count, err := bs.getUserTableData(table, userID)
		if err != nil {
			log.Printf("Failed to get user data for table %s: %v", table, err)
			continue
		}

		if count == 0 {
			continue
		}

		// Save to JSON file
		filePath := filepath.Join(backupDir, table+".json")
		if err := bs.saveJSONData(filePath, data); err != nil {
			log.Printf("Failed to save data for table %s: %v", table, err)
			continue
		}

		totalRecords += count
	}

	return totalRecords, nil
}

// getTableColumns returns column names for a table
func (bs *BackupService) getTableColumns(table string) ([]string, error) {
	rows, err := bs.db.Query(`
		SELECT column_name FROM information_schema.columns
		WHERE table_name = $1 AND table_schema = 'public'
		ORDER BY ordinal_position
	`, table)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []string
	for rows.Next() {
		var column string
		if err := rows.Scan(&column); err != nil {
			return nil, err
		}
		columns = append(columns, column)
	}

	return columns, nil
}

// exportTableData exports data from a single table
func (bs *BackupService) exportTableData(file *os.File, table string, columns []string) (int, error) {
	columnList := strings.Join(columns, ", ")
	query := fmt.Sprintf("SELECT %s FROM %s", columnList, table)

	rows, err := bs.db.Query(query)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		// This is a simplified export - in production, you'd use proper SQL dumping
		values := make([]interface{}, len(columns))
		for i := range values {
			values[i] = new(interface{})
		}

		if err := rows.Scan(values...); err != nil {
			return 0, err
		}

		// Convert to JSON for simplicity
		data := make(map[string]interface{})
		for i, column := range columns {
			data[column] = *(values[i].(*interface{}))
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			return 0, err
		}

		if _, err := file.WriteString(fmt.Sprintf("INSERT INTO %s VALUES %s;\n", table, string(jsonData))); err != nil {
			return 0, err
		}

		count++
	}

	return count, nil
}

// getUserTableData gets user-specific data from a table
func (bs *BackupService) getUserTableData(table string, userID uuid.UUID) ([]interface{}, int, error) {
	var query string
	var userColumn string

	// Determine which column contains the user ID
	switch table {
	case "users":
		userColumn = "id"
	case "user_sessions":
		userColumn = "user_id"
	case "documents", "document_versions":
		userColumn = "owner_id"
	case "email_accounts", "emails":
		userColumn = "user_id"
	case "chat_rooms", "chat_participants", "chat_messages", "message_reactions":
		// These are more complex - get rooms where user is participant
		return bs.getUserChatData(table, userID)
	case "contacts":
		userColumn = "user_id"
	case "calendars", "events":
		userColumn = "user_id"
	case "tasks":
		userColumn = "user_id"
	case "forms", "form_responses":
		userColumn = "owner_id"
	case "sync_devices", "sync_folders", "file_versions":
		userColumn = "user_id"
	case "meetings", "meeting_participants", "meeting_chat_messages", "webrtc_connections":
		userColumn = "host_id"
	default:
		return nil, 0, fmt.Errorf("unknown table: %s", table)
	}

	if table == "users" {
		query = fmt.Sprintf("SELECT * FROM %s WHERE %s = $1", table, userColumn)
	} else {
		query = fmt.Sprintf("SELECT * FROM %s WHERE %s = $1", table, userColumn)
	}

	rows, err := bs.db.Query(query, userID)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var data []interface{}
	count := 0
	for rows.Next() {
		// Get column types
		columnTypes, err := rows.ColumnTypes()
		if err != nil {
			return nil, 0, err
		}

		values := make([]interface{}, len(columnTypes))
		for i := range values {
			values[i] = new(interface{})
		}

		if err := rows.Scan(values...); err != nil {
			return nil, 0, err
		}

		// Convert to map
		rowData := make(map[string]interface{})
		for i, col := range columnTypes {
			rowData[col.Name()] = *(values[i].(*interface{}))
		}

		data = append(data, rowData)
		count++
	}

	return data, count, nil
}

// getUserChatData gets user-specific chat data
func (bs *BackupService) getUserChatData(table string, userID uuid.UUID) ([]interface{}, int, error) {
	var query string

	switch table {
	case "chat_rooms":
		query = `
			SELECT cr.* FROM chat_rooms cr
			JOIN chat_participants cp ON cr.id = cp.room_id
			WHERE cp.user_id = $1
		`
	case "chat_participants":
		query = "SELECT * FROM chat_participants WHERE user_id = $1"
	case "chat_messages":
		query = `
			SELECT cm.* FROM chat_messages cm
			JOIN chat_participants cp ON cm.room_id = cp.room_id
			WHERE cp.user_id = $1
		`
	case "message_reactions":
		query = `
			SELECT mr.* FROM message_reactions mr
			JOIN chat_messages cm ON mr.message_id = cm.id
			JOIN chat_participants cp ON cm.room_id = cp.room_id
			WHERE cp.user_id = $1
		`
	default:
		return nil, 0, fmt.Errorf("unsupported chat table: %s", table)
	}

	rows, err := bs.db.Query(query, userID)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var data []interface{}
	count := 0

	// Simplified - in production, you'd handle this properly
	for rows.Next() {
		count++
	}

	return data, count, nil
}

// restoreTable restores data to a table
func (bs *BackupService) restoreTable(filePath, table string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Read and execute SQL statements
	content, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	statements := strings.Split(string(content), "\n")
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		if _, err := bs.db.Exec(stmt); err != nil {
			log.Printf("Failed to execute statement: %v", err)
			// Continue with other statements
		}
	}

	return nil
}

// restoreUserTable restores user data to a table
func (bs *BackupService) restoreUserTable(filePath, table string, userID uuid.UUID) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var data []interface{}
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return err
	}

	// Insert data (simplified - in production, handle conflicts properly)
	for range data {
		// This would need proper handling based on table structure
		log.Printf("Restoring item to table %s", table)
	}

	return nil
}

// saveMetadata saves backup metadata to file
func (bs *BackupService) saveMetadata(filePath string, metadata *BackupMetadata) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(metadata)
}

// loadMetadata loads backup metadata from file
func (bs *BackupService) loadMetadata(filePath string) (*BackupMetadata, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var metadata BackupMetadata
	if err := json.NewDecoder(file).Decode(&metadata); err != nil {
		return nil, err
	}

	return &metadata, nil
}

// saveJSONData saves data to JSON file
func (bs *BackupService) saveJSONData(filePath string, data interface{}) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(data)
}

// calculateFileChecksum calculates the SHA-256 checksum of a file
func (bs *BackupService) calculateFileChecksum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// compressBackup compresses a backup directory
func (bs *BackupService) compressBackup(sourceDir, targetFile string) error {
	targetPath := filepath.Join(bs.backupPath, targetFile)

	file, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer file.Close()

	gzipWriter := gzip.NewWriter(file)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}
		header.Name = relPath

		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(tarWriter, file)
			return err
		}

		return nil
	})
}
