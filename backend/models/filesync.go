package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// SyncDevice represents a device participating in file synchronization
type SyncDevice struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	UserID      uuid.UUID  `json:"user_id" db:"user_id"`
	DeviceID    string     `json:"device_id" db:"device_id"`    // Unique device identifier
	DeviceName  string     `json:"device_name" db:"device_name"`
	DeviceType  string     `json:"device_type" db:"device_type"` // desktop, mobile, web
	PublicKey   []byte     `json:"-" db:"public_key"`           // For encryption
	LastSeen    time.Time  `json:"last_seen" db:"last_seen"`
	IsOnline    bool       `json:"is_online" db:"is_online"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
}

// SyncFolder represents a synchronized folder
type SyncFolder struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	UserID      uuid.UUID  `json:"user_id" db:"user_id"`
	Name        string     `json:"name" db:"name"`
	LocalPath   string     `json:"local_path" db:"local_path"`   // Local filesystem path
	RemotePath  string     `json:"remote_path" db:"remote_path"` // Remote/cloud path
	IsActive    bool       `json:"is_active" db:"is_active"`
	SyncMode    string     `json:"sync_mode" db:"sync_mode"`     // bidirectional, upload-only, download-only
	LastSync    *time.Time `json:"last_sync,omitempty" db:"last_sync"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

// FileVersion represents a version of a file
type FileVersion struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	FileID      uuid.UUID  `json:"file_id" db:"file_id"`
	Version     int        `json:"version" db:"version"`
	Size        int64      `json:"size" db:"size"`
	Hash        string     `json:"hash" db:"hash"`        // SHA-256 hash of file content
	ChunkHashes []string   `json:"chunk_hashes,omitempty" db:"chunk_hashes"` // For chunked transfer
	DeviceID    string     `json:"device_id" db:"device_id"`   // Device that created this version
	ModifiedBy  uuid.UUID  `json:"modified_by" db:"modified_by"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
}

// FileChunk represents a chunk of a file for efficient transfer
type FileChunk struct {
	ID         uuid.UUID `json:"id" db:"id"`
	VersionID  uuid.UUID `json:"version_id" db:"version_id"`
	ChunkIndex int       `json:"chunk_index" db:"chunk_index"`
	Size       int       `json:"size" db:"size"`
	Hash       string    `json:"hash" db:"hash"`        // SHA-256 of chunk
	Data       []byte    `json:"-" db:"data"`           // Encrypted chunk data
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// SyncConflict represents a file synchronization conflict
type SyncConflict struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	FileID        uuid.UUID  `json:"file_id" db:"file_id"`
	DeviceID      string     `json:"device_id" db:"device_id"`
	LocalVersion  int        `json:"local_version" db:"local_version"`
	RemoteVersion int        `json:"remote_version" db:"remote_version"`
	ConflictType  string     `json:"conflict_type" db:"conflict_type"`   // concurrent_edit, delete_conflict, etc.
	Resolution    *string    `json:"resolution,omitempty" db:"resolution"` // keep_local, keep_remote, merge
	ResolvedAt    *time.Time `json:"resolved_at,omitempty" db:"resolved_at"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
}

// SyncQueue represents pending synchronization operations
type SyncQueue struct {
	ID         uuid.UUID  `json:"id" db:"id"`
	FileID     uuid.UUID  `json:"file_id" db:"file_id"`
	DeviceID   string     `json:"device_id" db:"device_id"`
	Operation  string     `json:"operation" db:"operation"`   // create, update, delete, rename
	Priority   int        `json:"priority" db:"priority"`     // 1=low, 5=high
	Status     string     `json:"status" db:"status"`         // pending, processing, completed, failed
	RetryCount int        `json:"retry_count" db:"retry_count"`
	ErrorMsg   *string    `json:"error_msg,omitempty" db:"error_msg"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	ProcessedAt *time.Time `json:"processed_at,omitempty" db:"processed_at"`
}

// P2PMessage represents messages exchanged between devices in P2P network
type P2PMessageType string

const (
	P2PMessageHandshake  P2PMessageType = "handshake"
	P2PMessageFileList   P2PMessageType = "file_list"
	P2PMessageFileReq    P2PMessageType = "file_request"
	P2PMessageFileChunk  P2PMessageType = "file_chunk"
	P2PMessageSyncStatus P2PMessageType = "sync_status"
	P2PMessageConflict   P2PMessageType = "conflict"
)

type P2PMessage struct {
	ID          uuid.UUID       `json:"id"`
	Type        P2PMessageType  `json:"type"`
	FromDevice  string          `json:"from_device"`
	ToDevice    string          `json:"to_device"`
	Payload     interface{}     `json:"payload"`
	Timestamp   time.Time       `json:"timestamp"`
	Signature   string          `json:"signature"` // For message authenticity
}

// BandwidthLimit represents bandwidth throttling settings
type BandwidthLimit struct {
	ID         uuid.UUID `json:"id" db:"id"`
	UserID     uuid.UUID `json:"user_id" db:"user_id"`
	DeviceID   string    `json:"device_id" db:"device_id"`
	MaxUpload  int64     `json:"max_upload" db:"max_upload"`   // bytes per second
	MaxDownload int64    `json:"max_download" db:"max_download"` // bytes per second
	IsActive   bool      `json:"is_active" db:"is_active"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// FileWatchEvent represents filesystem change events
type FileWatchEvent struct {
	ID         uuid.UUID `json:"id"`
	FileID     uuid.UUID `json:"file_id"`
	EventType  string    `json:"event_type"`  // create, modify, delete, rename
	Path       string    `json:"path"`
	OldPath    *string   `json:"old_path,omitempty"` // For renames
	Size       int64     `json:"size"`
	ModifiedAt time.Time `json:"modified_at"`
	DeviceID   string    `json:"device_id"`
	Timestamp  time.Time `json:"timestamp"`
}

// SyncRequest represents synchronization requests
type SyncRequest struct {
	DeviceID    string        `json:"device_id"`
	FolderID    uuid.UUID     `json:"folder_id"`
	FileHashes  map[string]string `json:"file_hashes"` // path -> hash mapping
	LastSync    *time.Time    `json:"last_sync,omitempty"`
	IncludeDeletes bool       `json:"include_deletes"`
}

type SyncResponse struct {
	DeviceID      string                   `json:"device_id"`
	Changes       []FileChange             `json:"changes"`
	Conflicts     []SyncConflictResponse   `json:"conflicts"`
	NextSyncToken string                   `json:"next_sync_token"`
}

type FileChange struct {
	Path         string     `json:"path"`
	ChangeType   string     `json:"change_type"` // add, modify, delete, rename
	NewPath      *string    `json:"new_path,omitempty"`
	Size         int64      `json:"size"`
	Hash         string     `json:"hash"`
	ModifiedAt   time.Time  `json:"modified_at"`
	Version      int        `json:"version"`
}

type SyncConflictResponse struct {
	FileID       uuid.UUID `json:"file_id"`
	Path         string    `json:"path"`
	LocalVersion int       `json:"local_version"`
	RemoteVersion int      `json:"remote_version"`
	ConflictType string    `json:"conflict_type"`
	Suggestions  []string  `json:"suggestions"` // keep_local, keep_remote, merge, rename
}

// ToResponse converts SyncDevice to response format
func (sd *SyncDevice) ToResponse() map[string]interface{} {
	return map[string]interface{}{
		"id":           sd.ID,
		"device_id":    sd.DeviceID,
		"device_name":  sd.DeviceName,
		"device_type":  sd.DeviceType,
		"last_seen":    sd.LastSeen,
		"is_online":    sd.IsOnline,
		"created_at":   sd.CreatedAt,
	}
}

// ToResponse converts SyncFolder to response format
func (sf *SyncFolder) ToResponse() map[string]interface{} {
	return map[string]interface{}{
		"id":          sf.ID,
		"name":        sf.Name,
		"local_path":  sf.LocalPath,
		"remote_path": sf.RemotePath,
		"is_active":   sf.IsActive,
		"sync_mode":   sf.SyncMode,
		"last_sync":   sf.LastSync,
		"created_at":  sf.CreatedAt,
		"updated_at":  sf.UpdatedAt,
	}
}

// GetMerkleRoot calculates Merkle root for a file's chunks
func (fv *FileVersion) GetMerkleRoot() string {
	if len(fv.ChunkHashes) == 0 {
		return fv.Hash
	}

	// Simplified Merkle tree calculation
	// In production, this would build a proper Merkle tree
	hash := fv.ChunkHashes[0]
	for i := 1; i < len(fv.ChunkHashes); i++ {
		hash = HashStrings(hash, fv.ChunkHashes[i])
	}
	return hash
}

// HashStrings combines two strings with SHA-256 (placeholder)
func HashStrings(s1, s2 string) string {
	// In production, use crypto/sha256
	return s1 + s2 // Placeholder
}

// Validate checks if SyncFolder has valid data
func (sf *SyncFolder) Validate() error {
	if sf.Name == "" {
		return fmt.Errorf("folder name is required")
	}
	if sf.LocalPath == "" {
		return fmt.Errorf("local path is required")
	}

	validModes := []string{"bidirectional", "upload-only", "download-only"}
	for _, mode := range validModes {
		if sf.SyncMode == mode {
			return nil
		}
	}
	return fmt.Errorf("invalid sync mode: %s", sf.SyncMode)
}

// IsResolved checks if a conflict has been resolved
func (sc *SyncConflict) IsResolved() bool {
	return sc.Resolution != nil && sc.ResolvedAt != nil
}

// ShouldRetry checks if a sync operation should be retried
func (sq *SyncQueue) ShouldRetry() bool {
	return sq.RetryCount < 3 && sq.Status == "failed"
}

// GetChunkSize returns optimal chunk size for file transfer
func GetChunkSize(fileSize int64) int {
	// Adaptive chunk sizing based on file size
	if fileSize < 1024*1024 { // < 1MB
		return 64 * 1024 // 64KB
	} else if fileSize < 100*1024*1024 { // < 100MB
		return 1024 * 1024 // 1MB
	} else {
		return 4 * 1024 * 1024 // 4MB
	}
}

// IsFileIgnored checks if a file should be ignored during sync
func IsFileIgnored(filename string) bool {
	ignoredPatterns := []string{
		".DS_Store",
		"Thumbs.db",
		"desktop.ini",
		".tmp",
		"~$", // Office temp files
	}

	for _, pattern := range ignoredPatterns {
		if strings.Contains(filename, pattern) {
			return true
		}
	}

	return false
}
