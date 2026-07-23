package services

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/google/uuid"
	"tpt-titan/backend/models"
)

// FileSyncService handles file synchronization operations
type FileSyncService struct {
	db *sql.DB
}

// NewFileSyncService creates a new file sync service
func NewFileSyncService(db *sql.DB) *FileSyncService {
	return &FileSyncService{db: db}
}

// RegisterDevice registers a new device for synchronization
func (s *FileSyncService) RegisterDevice(userID uuid.UUID, req models.SyncDevice) (*models.SyncDevice, error) {
	// Check if device already exists
	var existingID uuid.UUID
	err := s.db.QueryRow(`
		SELECT id FROM sync_devices WHERE user_id = $1 AND device_id = $2
	`, userID, req.DeviceID).Scan(&existingID)

	if err == nil {
		// Device exists, update it
		req.ID = existingID
		return s.updateDevice(userID, req)
	}

	// Device doesn't exist, create it
	req.ID = uuid.New()
	req.UserID = userID
	req.LastSeen = time.Now()
	req.IsOnline = true

	query := `
		INSERT INTO sync_devices (id, user_id, device_id, device_name, device_type, public_key, last_seen, is_online, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err = s.db.Exec(query,
		req.ID, req.UserID, req.DeviceID, req.DeviceName, req.DeviceType,
		req.PublicKey, req.LastSeen, req.IsOnline, req.LastSeen,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to register device: %w", err)
	}

	return &req, nil
}

// updateDevice updates an existing device
func (s *FileSyncService) updateDevice(userID uuid.UUID, device models.SyncDevice) (*models.SyncDevice, error) {
	query := `
		UPDATE sync_devices
		SET device_name = $1, device_type = $2, public_key = $3, last_seen = $4, is_online = $5
		WHERE id = $6 AND user_id = $7
	`

	now := time.Now()
	_, err := s.db.Exec(query,
		device.DeviceName, device.DeviceType, device.PublicKey, now, device.IsOnline,
		device.ID, userID,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update device: %w", err)
	}

	device.LastSeen = now
	return &device, nil
}

// GetDevices retrieves all devices for a user
func (s *FileSyncService) GetDevices(userID uuid.UUID) ([]models.SyncDevice, error) {
	query := `
		SELECT id, user_id, device_id, device_name, device_type, public_key, last_seen, is_online, created_at
		FROM sync_devices
		WHERE user_id = $1
		ORDER BY last_seen DESC
	`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query devices: %w", err)
	}
	defer rows.Close()

	var devices []models.SyncDevice
	for rows.Next() {
		var device models.SyncDevice
		err := rows.Scan(
			&device.ID, &device.UserID, &device.DeviceID, &device.DeviceName,
			&device.DeviceType, &device.PublicKey, &device.LastSeen, &device.IsOnline, &device.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan device: %w", err)
		}
		devices = append(devices, device)
	}

	return devices, nil
}

// CreateSyncFolder creates a new synchronized folder
func (s *FileSyncService) CreateSyncFolder(userID uuid.UUID, req models.SyncFolder) (*models.SyncFolder, error) {
	// Validate the folder
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid sync folder: %w", err)
	}

	req.ID = uuid.New()
	req.UserID = userID
	req.IsActive = true
	req.CreatedAt = time.Now()
	req.UpdatedAt = req.CreatedAt

	// Set default remote path if not provided
	if req.RemotePath == "" {
		req.RemotePath = fmt.Sprintf("/sync/%s/%s", userID.String(), req.Name)
	}

	query := `
		INSERT INTO sync_folders (id, user_id, name, local_path, remote_path, is_active, sync_mode, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := s.db.Exec(query,
		req.ID, req.UserID, req.Name, req.LocalPath, req.RemotePath,
		req.IsActive, req.SyncMode, req.CreatedAt, req.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create sync folder: %w", err)
	}

	return &req, nil
}

// GetSyncFolders retrieves all sync folders for a user
func (s *FileSyncService) GetSyncFolders(userID uuid.UUID) ([]models.SyncFolder, error) {
	query := `
		SELECT id, user_id, name, local_path, remote_path, is_active, sync_mode, last_sync, created_at, updated_at
		FROM sync_folders
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query sync folders: %w", err)
	}
	defer rows.Close()

	var folders []models.SyncFolder
	for rows.Next() {
		var folder models.SyncFolder
		err := rows.Scan(
			&folder.ID, &folder.UserID, &folder.Name, &folder.LocalPath, &folder.RemotePath,
			&folder.IsActive, &folder.SyncMode, &folder.LastSync, &folder.CreatedAt, &folder.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan sync folder: %w", err)
		}
		folders = append(folders, folder)
	}

	return folders, nil
}

// CreateFileVersion creates a new version of a file
func (s *FileSyncService) CreateFileVersion(fileID uuid.UUID, deviceID string, userID uuid.UUID) (*models.FileVersion, error) {
	// Get current version number
	var currentVersion int
	err := s.db.QueryRow(`
		SELECT COALESCE(MAX(version), 0) FROM file_versions WHERE file_id = $1
	`, fileID).Scan(&currentVersion)

	if err != nil {
		return nil, fmt.Errorf("failed to get current version: %w", err)
	}

	newVersion := currentVersion + 1

	// Get file info
	var filePath string
	var fileSize int64
	err = s.db.QueryRow(`
		SELECT path, size FROM files WHERE id = $1
	`, fileID).Scan(&filePath, &fileSize)

	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	// Calculate file hash
	hash, err := s.calculateFileHash(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate file hash: %w", err)
	}

	version := models.FileVersion{
		ID:         uuid.New(),
		FileID:     fileID,
		Version:    newVersion,
		Size:       fileSize,
		Hash:       hash,
		DeviceID:   deviceID,
		ModifiedBy: userID,
		CreatedAt:  time.Now(),
	}

	// Create chunks for large files
	if fileSize > 1024*1024 { // > 1MB
		chunks, err := s.createFileChunks(version.ID, filePath, fileSize)
		if err != nil {
			return nil, fmt.Errorf("failed to create file chunks: %w", err)
		}
		version.ChunkHashes = chunks
	}

	query := `
		INSERT INTO file_versions (id, file_id, version, size, hash, chunk_hashes, device_id, modified_by, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err = s.db.Exec(query,
		version.ID, version.FileID, version.Version, version.Size, version.Hash,
		version.ChunkHashes, version.DeviceID, version.ModifiedBy, version.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create file version: %w", err)
	}

	return &version, nil
}

// calculateFileHash calculates SHA-256 hash of a file
func (s *FileSyncService) calculateFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", fmt.Errorf("failed to hash file: %w", err)
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// createFileChunks splits a large file into chunks
func (s *FileSyncService) createFileChunks(versionID uuid.UUID, filePath string, fileSize int64) ([]string, error) {
	chunkSize := models.GetChunkSize(fileSize)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file for chunking: %w", err)
	}
	defer file.Close()

	var chunkHashes []string
	buffer := make([]byte, chunkSize)
	chunkIndex := 0

	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("failed to read file chunk: %w", err)
		}

		if n == 0 {
			break
		}

		// Calculate chunk hash
		hasher := sha256.New()
		hasher.Write(buffer[:n])
		chunkHash := hex.EncodeToString(hasher.Sum(nil))

		// Store chunk (in production, this would be encrypted)
		chunk := models.FileChunk{
			ID:         uuid.New(),
			VersionID:  versionID,
			ChunkIndex: chunkIndex,
			Size:       n,
			Hash:       chunkHash,
			Data:       buffer[:n], // In production: encrypt this
			CreatedAt:  time.Now(),
		}

		query := `
			INSERT INTO file_chunks (id, version_id, chunk_index, size, hash, data, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`

		_, insertErr := s.db.Exec(query,
			chunk.ID, chunk.VersionID, chunk.ChunkIndex, chunk.Size,
			chunk.Hash, chunk.Data, chunk.CreatedAt,
		)

		if insertErr != nil {
			return nil, fmt.Errorf("failed to store chunk: %w", insertErr)
		}

		chunkHashes = append(chunkHashes, chunkHash)
		chunkIndex++

		if err == io.EOF {
			break
		}
	}

	return chunkHashes, nil
}

// SyncFolder performs synchronization for a folder
func (s *FileSyncService) SyncFolder(userID uuid.UUID, folderID uuid.UUID, deviceID string) (*models.SyncResponse, error) {
	// Get folder info
	var folder models.SyncFolder
	err := s.db.QueryRow(`
		SELECT id, user_id, name, local_path, remote_path, sync_mode FROM sync_folders
		WHERE id = $1 AND user_id = $2
	`, folderID, userID).Scan(
		&folder.ID, &folder.UserID, &folder.Name, &folder.LocalPath,
		&folder.RemotePath, &folder.SyncMode,
	)

	if err != nil {
		return nil, fmt.Errorf("folder not found: %w", err)
	}

	// Get local file hashes (this would come from the device in real implementation)
	localHashes := make(map[string]string) // path -> hash

	// Get remote file changes since last sync
	changes, err := s.getFileChanges(folderID, folder.LastSync)
	if err != nil {
		return nil, fmt.Errorf("failed to get file changes: %w", err)
	}

	// Convert file changes to path->hash map for conflict detection
	remoteChangeMap := make(map[string]string)
	for _, c := range changes {
		remoteChangeMap[c.Path] = c.Hash
	}
	// Check for conflicts
	conflicts, err := s.detectConflicts(folderID, deviceID, localHashes, remoteChangeMap)
	if err != nil {
		return nil, fmt.Errorf("failed to detect conflicts: %w", err)
	}

	// Update last sync time
	now := time.Now()
	_, err = s.db.Exec(`UPDATE sync_folders SET last_sync = $1 WHERE id = $2`, now, folderID)
	if err != nil {
		return nil, fmt.Errorf("failed to update last sync time: %w", err)
	}

	response := &models.SyncResponse{
		DeviceID:      deviceID,
		Changes:       changes,
		Conflicts:     conflicts,
		NextSyncToken: fmt.Sprintf("%d", now.Unix()),
	}

	return response, nil
}

// getFileChanges gets file changes since last sync
func (s *FileSyncService) getFileChanges(folderID uuid.UUID, since *time.Time) ([]models.FileChange, error) {
	query := `
		SELECT f.name, f.path, f.size, fv.hash, fv.created_at, fv.version
		FROM files f
		JOIN file_versions fv ON f.id = fv.file_id
		WHERE f.parent_id IN (
			SELECT id FROM files WHERE name IN (
				SELECT name FROM sync_folders WHERE id = $1
			)
		)
	`

	args := []interface{}{folderID}
	if since != nil {
		query += " AND fv.created_at > $2"
		args = append(args, *since)
	}

	query += " ORDER BY fv.created_at ASC"

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query file changes: %w", err)
	}
	defer rows.Close()

	var changes []models.FileChange
	for rows.Next() {
		var change models.FileChange
		var fileName, filePath string
		var createdAt time.Time

		err := rows.Scan(
			&fileName, &filePath, &change.Size, &change.Hash,
			&createdAt, &change.Version,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan file change: %w", err)
		}

		change.Path = filePath
		change.ChangeType = "modify" // Simplified - in production, detect create/modify/delete
		change.ModifiedAt = createdAt

		changes = append(changes, change)
	}

	return changes, nil
}

// detectConflicts detects synchronization conflicts
func (s *FileSyncService) detectConflicts(folderID uuid.UUID, deviceID string, localHashes, remoteChanges map[string]string) ([]models.SyncConflictResponse, error) {
	var conflicts []models.SyncConflictResponse

	// This is a simplified conflict detection
	// In production, this would compare local and remote file states
	for path, remoteHash := range remoteChanges {
		if localHash, exists := localHashes[path]; exists && localHash != remoteHash {
			conflict := models.SyncConflictResponse{
				Path:          path,
				LocalVersion:  1, // Simplified
				RemoteVersion: 1, // Simplified
				ConflictType:  "concurrent_edit",
				Suggestions:   []string{"keep_local", "keep_remote", "merge"},
			}
			conflicts = append(conflicts, conflict)
		}
	}

	return conflicts, nil
}

// ResolveConflict resolves a synchronization conflict
func (s *FileSyncService) ResolveConflict(userID uuid.UUID, conflictID uuid.UUID, resolution string) error {
	query := `
		UPDATE sync_conflicts
		SET resolution = $1, resolved_at = $2
		WHERE id = $3 AND file_id IN (
			SELECT id FROM files WHERE owner_id = $4
		)
	`

	now := time.Now()
	result, err := s.db.Exec(query, resolution, now, conflictID, userID)
	if err != nil {
		return fmt.Errorf("failed to resolve conflict: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("conflict not found or access denied")
	}

	return nil
}

// SetBandwidthLimit sets bandwidth throttling for a device
func (s *FileSyncService) SetBandwidthLimit(userID uuid.UUID, deviceID string, maxUpload, maxDownload int64) error {
	// Check if limit already exists
	var existingID uuid.UUID
	err := s.db.QueryRow(`
		SELECT id FROM bandwidth_limits WHERE user_id = $1 AND device_id = $2
	`, userID, deviceID).Scan(&existingID)

	now := time.Now()

	if err == sql.ErrNoRows {
		// Create new limit
		query := `
			INSERT INTO bandwidth_limits (id, user_id, device_id, max_upload, max_download, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`

		_, err = s.db.Exec(query, uuid.New(), userID, deviceID, maxUpload, maxDownload, now, now)
	} else {
		// Update existing limit
		query := `
			UPDATE bandwidth_limits
			SET max_upload = $1, max_download = $2, updated_at = $3
			WHERE id = $4
		`

		_, err = s.db.Exec(query, maxUpload, maxDownload, now, existingID)
	}

	if err != nil {
		return fmt.Errorf("failed to set bandwidth limit: %w", err)
	}

	return nil
}

// GetFileChunk retrieves a file chunk for transfer
func (s *FileSyncService) GetFileChunk(userID uuid.UUID, versionID uuid.UUID, chunkIndex int) (*models.FileChunk, error) {
	// Verify user has access to this file version
	var hasAccess bool
	err := s.db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM file_versions fv
			JOIN files f ON fv.file_id = f.id
			WHERE fv.id = $1 AND f.owner_id = $2
		)
	`, versionID, userID).Scan(&hasAccess)

	if err != nil || !hasAccess {
		return nil, fmt.Errorf("access denied to file chunk")
	}

	var chunk models.FileChunk
	err = s.db.QueryRow(`
		SELECT id, version_id, chunk_index, size, hash, data, created_at
		FROM file_chunks
		WHERE version_id = $1 AND chunk_index = $2
	`, versionID, chunkIndex).Scan(
		&chunk.ID, &chunk.VersionID, &chunk.ChunkIndex,
		&chunk.Size, &chunk.Hash, &chunk.Data, &chunk.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("chunk not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get chunk: %w", err)
	}

	return &chunk, nil
}

// WatchFolder starts watching a folder for changes using fsnotify
func (s *FileSyncService) WatchFolder(folderID uuid.UUID) error {
	// Get folder info
	var folder models.SyncFolder
	err := s.db.QueryRow(`
		SELECT id, user_id, name, local_path, remote_path, is_active FROM sync_folders
		WHERE id = $1
	`, folderID).Scan(
		&folder.ID, &folder.UserID, &folder.Name, &folder.LocalPath,
		&folder.RemotePath, &folder.IsActive,
	)

	if err != nil {
		return fmt.Errorf("folder not found: %w", err)
	}

	if !folder.IsActive {
		return fmt.Errorf("folder is not active")
	}

	// Check if folder exists
	if _, err := os.Stat(folder.LocalPath); os.IsNotExist(err) {
		return fmt.Errorf("folder does not exist: %s", folder.LocalPath)
	}

	// Create new watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create watcher: %w", err)
	}

	// Start watching in a goroutine
	go s.watchFolderLoop(watcher, folderID, folder.LocalPath)

	// Add the folder to watch
	err = watcher.Add(folder.LocalPath)
	if err != nil {
		watcher.Close()
		return fmt.Errorf("failed to watch folder: %w", err)
	}

	// Walk through subdirectories and add them to watch
	err = filepath.Walk(folder.LocalPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return watcher.Add(path)
		}
		return nil
	})

	if err != nil {
		watcher.Close()
		return fmt.Errorf("failed to walk folder: %w", err)
	}

	log.Printf("Started watching folder: %s", folder.LocalPath)
	return nil
}

// watchFolderLoop handles filesystem events
func (s *FileSyncService) watchFolderLoop(watcher *fsnotify.Watcher, folderID uuid.UUID, basePath string) {
	defer watcher.Close()

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			// Get relative path
			relPath, err := filepath.Rel(basePath, event.Name)
			if err != nil {
				log.Printf("Failed to get relative path: %v", err)
				continue
			}

			// Skip hidden files and directories
			if strings.HasPrefix(filepath.Base(relPath), ".") {
				continue
			}

			// Handle different event types
			eventType := ""
			if event.Has(fsnotify.Create) {
				eventType = "create"
			} else if event.Has(fsnotify.Write) {
				eventType = "modify"
			} else if event.Has(fsnotify.Remove) || event.Has(fsnotify.Rename) {
				eventType = "delete"
			}

			if eventType != "" {
				err := s.handleFileEvent(folderID, relPath, eventType, event.Name)
				if err != nil {
					log.Printf("Failed to handle file event: %v", err)
				}
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Watcher error: %v", err)
		}
	}
}

// handleFileEvent processes a filesystem event
func (s *FileSyncService) handleFileEvent(folderID uuid.UUID, relPath, eventType, fullPath string) error {
	// Get file info
	info, err := os.Stat(fullPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to stat file: %w", err)
	}

	size := int64(0)
	if info != nil {
		size = info.Size()
	}

	// Create file watch event
	event := models.FileWatchEvent{
		ID:         uuid.New(),
		FileID:     uuid.New(), // This would be looked up from database
		EventType:  eventType,
		Path:       relPath,
		Size:       size,
		ModifiedAt: time.Now(),
		DeviceID:   "local", // This would come from device registration
		Timestamp:  time.Now(),
	}

	// Store the event (in production, this would trigger sync operations)
	log.Printf("File event: %s %s (id=%s)", eventType, relPath, event.ID)

	// For now, just log. In production, this would:
	// 1. Check if file should be synced
	// 2. Calculate hash
	// 3. Create file version
	// 4. Queue for synchronization
	// 5. Notify other devices

	return nil
}

// CleanupOldVersions removes old file versions to save space
func (s *FileSyncService) CleanupOldVersions(userID uuid.UUID, keepVersions int) error {
	// Keep only the most recent N versions per file
	query := `
		DELETE FROM file_versions
		WHERE file_id IN (
			SELECT id FROM files WHERE owner_id = $1
		) AND id NOT IN (
			SELECT id FROM (
				SELECT id, ROW_NUMBER() OVER (PARTITION BY file_id ORDER BY version DESC) as rn
				FROM file_versions
				WHERE file_id IN (SELECT id FROM files WHERE owner_id = $1)
			) t WHERE rn <= $2
		)
	`

	_, err := s.db.Exec(query, userID, keepVersions)
	if err != nil {
		return fmt.Errorf("failed to cleanup old versions: %w", err)
	}

	return nil
}

// GetSyncStatus returns the current synchronization status
func (s *FileSyncService) GetSyncStatus(userID uuid.UUID) (map[string]interface{}, error) {
	status := make(map[string]interface{})

	// Get device count
	var deviceCount int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM sync_devices WHERE user_id = $1`, userID).Scan(&deviceCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get device count: %w", err)
	}
	status["device_count"] = deviceCount

	// Get folder count
	var folderCount int
	err = s.db.QueryRow(`SELECT COUNT(*) FROM sync_folders WHERE user_id = $1`, userID).Scan(&folderCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get folder count: %w", err)
	}
	status["folder_count"] = folderCount

	// Get pending operations
	var pendingOps int
	err = s.db.QueryRow(`
		SELECT COUNT(*) FROM sync_queue sq
		JOIN sync_devices sd ON sq.device_id = sd.device_id
		WHERE sd.user_id = $1 AND sq.status = 'pending'
	`, userID).Scan(&pendingOps)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending operations: %w", err)
	}
	status["pending_operations"] = pendingOps

	// Get conflicts
	var conflictCount int
	err = s.db.QueryRow(`
		SELECT COUNT(*) FROM sync_conflicts sc
		JOIN files f ON sc.file_id = f.id
		WHERE f.owner_id = $1 AND sc.resolution IS NULL
	`, userID).Scan(&conflictCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get conflict count: %w", err)
	}
	status["conflicts"] = conflictCount

	return status, nil
}
