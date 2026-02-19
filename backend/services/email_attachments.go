package services

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
)

// EmailAttachmentService handles email attachments
type EmailAttachmentService struct {
	db         *sql.DB
	storageSvc *StorageService // Would be implemented for file storage
	maxSize    int64           // Maximum attachment size in bytes
	allowedTypes []string      // Allowed MIME types
}

// EmailAttachment represents an email attachment
type EmailAttachment struct {
	ID          uuid.UUID `json:"id"`
	EmailID     uuid.UUID `json:"email_id"`
	Filename    string    `json:"filename"`
	ContentType string    `json:"content_type"`
	Size        int64     `json:"size"`        // Size in bytes
	Hash        string    `json:"hash"`        // SHA256 hash for integrity
	StoragePath string    `json:"storage_path"`
	IsInline    bool      `json:"is_inline"`   // True for inline images
	ContentID   string    `json:"content_id,omitempty"` // For inline attachments
	CreatedAt   time.Time `json:"created_at"`
}

// AttachmentUpload represents an attachment being uploaded
type AttachmentUpload struct {
	Filename    string
	ContentType string
	Size        int64
	Data        []byte
	IsInline    bool
	ContentID   string
}

// AttachmentDownload represents an attachment being downloaded
type AttachmentDownload struct {
	ID       uuid.UUID
	Data     []byte
	Filename string
	MimeType string
}

// AttachmentVirusScan represents virus scan results
type AttachmentVirusScan struct {
	AttachmentID uuid.UUID `json:"attachment_id"`
	IsSafe       bool      `json:"is_safe"`
	ScanEngine   string    `json:"scan_engine"`
	ScanResult   string    `json:"scan_result"`
	ScannedAt    time.Time `json:"scanned_at"`
}

// NewEmailAttachmentService creates a new email attachment service
func NewEmailAttachmentService(db *sql.DB, storageSvc *StorageService) *EmailAttachmentService {
	return &EmailAttachmentService{
		db:         db,
		storageSvc: storageSvc,
		maxSize:    25 * 1024 * 1024, // 25MB default
		allowedTypes: []string{
			"application/pdf",
			"application/msword",
			"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
			"application/vnd.ms-excel",
			"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
			"application/vnd.ms-powerpoint",
			"application/vnd.openxmlformats-officedocument.presentationml.presentation",
			"image/jpeg",
			"image/png",
			"image/gif",
			"image/webp",
			"text/plain",
			"text/csv",
			"application/zip",
			"application/x-zip-compressed",
		},
	}
}

// SaveAttachment saves an email attachment to storage
func (eas *EmailAttachmentService) SaveAttachment(emailID uuid.UUID, upload *AttachmentUpload) (*EmailAttachment, error) {
	// Validate attachment
	if err := eas.validateAttachment(upload); err != nil {
		return nil, err
	}

	// Generate hash for integrity checking
	hash := eas.generateHash(upload.Data)

	// Generate storage path
	storagePath := eas.generateStoragePath(emailID, upload.Filename)

	// Save to storage (would use actual storage service)
	if err := eas.storageSvc.SaveFile(storagePath, upload.Data); err != nil {
		return nil, fmt.Errorf("failed to save attachment: %w", err)
	}

	// Create attachment record
	attachment := &EmailAttachment{
		ID:          uuid.New(),
		EmailID:     emailID,
		Filename:    upload.Filename,
		ContentType: upload.ContentType,
		Size:        upload.Size,
		Hash:        hash,
		StoragePath: storagePath,
		IsInline:    upload.IsInline,
		ContentID:   upload.ContentID,
		CreatedAt:   time.Now(),
	}

	// Save to database
	if err := eas.saveAttachmentToDB(attachment); err != nil {
		// Clean up storage if DB save fails
		eas.storageSvc.DeleteFile(storagePath)
		return nil, err
	}

	// Schedule virus scan
	go eas.scanAttachmentForViruses(attachment.ID, upload.Data)

	return attachment, nil
}

// GetAttachment retrieves an attachment by ID
func (eas *EmailAttachmentService) GetAttachment(attachmentID uuid.UUID) (*EmailAttachment, error) {
	var attachment EmailAttachment
	query := `
		SELECT id, email_id, filename, content_type, size, hash, storage_path,
		       is_inline, content_id, created_at
		FROM email_attachments WHERE id = $1
	`

	err := eas.db.QueryRow(query, attachmentID).Scan(
		&attachment.ID, &attachment.EmailID, &attachment.Filename, &attachment.ContentType,
		&attachment.Size, &attachment.Hash, &attachment.StoragePath, &attachment.IsInline,
		&attachment.ContentID, &attachment.CreatedAt,
	)

	return &attachment, err
}

// DownloadAttachment downloads attachment data
func (eas *EmailAttachmentService) DownloadAttachment(attachmentID uuid.UUID) (*AttachmentDownload, error) {
	attachment, err := eas.GetAttachment(attachmentID)
	if err != nil {
		return nil, err
	}

	// Get data from storage
	data, err := eas.storageSvc.GetFile(attachment.StoragePath)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve attachment data: %w", err)
	}

	// Verify integrity
	if eas.generateHash(data) != attachment.Hash {
		return nil, fmt.Errorf("attachment integrity check failed")
	}

	return &AttachmentDownload{
		ID:       attachment.ID,
		Data:     data,
		Filename: attachment.Filename,
		MimeType: attachment.ContentType,
	}, nil
}

// DeleteAttachment deletes an attachment
func (eas *EmailAttachmentService) DeleteAttachment(attachmentID uuid.UUID) error {
	attachment, err := eas.GetAttachment(attachmentID)
	if err != nil {
		return err
	}

	// Delete from storage
	if err := eas.storageSvc.DeleteFile(attachment.StoragePath); err != nil {
		// Log error but continue with DB deletion
	}

	// Delete from database
	query := `DELETE FROM email_attachments WHERE id = $1`
	_, err = eas.db.Exec(query, attachmentID)

	return err
}

// GetEmailAttachments gets all attachments for an email
func (eas *EmailAttachmentService) GetEmailAttachments(emailID uuid.UUID) ([]EmailAttachment, error) {
	query := `
		SELECT id, email_id, filename, content_type, size, hash, storage_path,
		       is_inline, content_id, created_at
		FROM email_attachments
		WHERE email_id = $1
		ORDER BY created_at
	`

	rows, err := eas.db.Query(query, emailID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attachments []EmailAttachment
	for rows.Next() {
		var attachment EmailAttachment
		err := rows.Scan(
			&attachment.ID, &attachment.EmailID, &attachment.Filename, &attachment.ContentType,
			&attachment.Size, &attachment.Hash, &attachment.StoragePath, &attachment.IsInline,
			&attachment.ContentID, &attachment.CreatedAt,
		)
		if err != nil {
			continue
		}
		attachments = append(attachments, attachment)
	}

	return attachments, nil
}

// ProcessEmailAttachments processes attachments from incoming email
func (eas *EmailAttachmentService) ProcessEmailAttachments(emailID uuid.UUID, email *mail.Message) ([]uuid.UUID, error) {
	var attachmentIDs []uuid.UUID

	// Parse email with attachments
	// This would use a proper email parsing library
	parts := eas.parseEmailParts(email)

	for _, part := range parts {
		if eas.isAttachmentPart(part) {
			upload := &AttachmentUpload{
				Filename:    eas.extractFilename(part),
				ContentType: eas.extractContentType(part),
				Data:        eas.extractPartData(part),
				IsInline:    eas.isInlineAttachment(part),
				ContentID:   eas.extractContentID(part),
			}

			upload.Size = int64(len(upload.Data))

			attachment, err := eas.SaveAttachment(emailID, upload)
			if err != nil {
				// Log error but continue processing other attachments
				continue
			}

			attachmentIDs = append(attachmentIDs, attachment.ID)
		}
	}

	return attachmentIDs, nil
}

// CreateAttachmentFromUpload creates an attachment from file upload
func (eas *EmailAttachmentService) CreateAttachmentFromUpload(emailID uuid.UUID, filename string, data []byte, contentType string) (*EmailAttachment, error) {
	upload := &AttachmentUpload{
		Filename:    filename,
		ContentType: contentType,
		Size:        int64(len(data)),
		Data:        data,
		IsInline:    false,
	}

	return eas.SaveAttachment(emailID, upload)
}

// GetAttachmentPreview generates a preview for an attachment
func (eas *EmailAttachmentService) GetAttachmentPreview(attachmentID uuid.UUID, maxWidth, maxHeight int) ([]byte, error) {
	attachment, err := eas.GetAttachment(attachmentID)
	if err != nil {
		return nil, err
	}

	// Only generate previews for images
	if !strings.HasPrefix(attachment.ContentType, "image/") {
		return nil, fmt.Errorf("preview not available for this file type")
	}

	data, err := eas.DownloadAttachment(attachmentID)
	if err != nil {
		return nil, err
	}

	// Generate thumbnail (would use image processing library)
	thumbnail := eas.generateThumbnail(data.Data, maxWidth, maxHeight, attachment.ContentType)

	return thumbnail, nil
}

// SearchAttachments searches for attachments by filename or content type
func (eas *EmailAttachmentService) SearchAttachments(query string, limit int) ([]EmailAttachment, error) {
	searchQuery := `
		SELECT id, email_id, filename, content_type, size, hash, storage_path,
		       is_inline, content_id, created_at
		FROM email_attachments
		WHERE filename ILIKE $1 OR content_type ILIKE $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := eas.db.Query(searchQuery, "%"+query+"%", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attachments []EmailAttachment
	for rows.Next() {
		var attachment EmailAttachment
		err := rows.Scan(
			&attachment.ID, &attachment.EmailID, &attachment.Filename, &attachment.ContentType,
			&attachment.Size, &attachment.Hash, &attachment.StoragePath, &attachment.IsInline,
			&attachment.ContentID, &attachment.CreatedAt,
		)
		if err != nil {
			continue
		}
		attachments = append(attachments, attachment)
	}

	return attachments, nil
}

// GetAttachmentStats returns attachment statistics
func (eas *EmailAttachmentService) GetAttachmentStats(userID uuid.UUID) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total attachments count
	var totalCount int
	eas.db.QueryRow("SELECT COUNT(*) FROM email_attachments").Scan(&totalCount)
	stats["total_attachments"] = totalCount

	// Total size
	var totalSize int64
	eas.db.QueryRow("SELECT COALESCE(SUM(size), 0) FROM email_attachments").Scan(&totalSize)
	stats["total_size_bytes"] = totalSize
	stats["total_size_mb"] = float64(totalSize) / (1024 * 1024)

	// Content type breakdown
	typeStats := []map[string]interface{}{}
	rows, err := eas.db.Query(`
		SELECT content_type, COUNT(*) as count, SUM(size) as total_size
		FROM email_attachments
		GROUP BY content_type
		ORDER BY count DESC
		LIMIT 10
	`)
	if err == nil {
		for rows.Next() {
			var contentType string
			var count int
			var size int64
			rows.Scan(&contentType, &count, &size)
			typeStats = append(typeStats, map[string]interface{}{
				"content_type": contentType,
				"count":        count,
				"total_size":   size,
			})
		}
		rows.Close()
	}
	stats["content_type_breakdown"] = typeStats

	// Recent attachments
	recent := []EmailAttachment{}
	rows, err = eas.db.Query(`
		SELECT id, email_id, filename, content_type, size, created_at
		FROM email_attachments
		ORDER BY created_at DESC
		LIMIT 5
	`)
	if err == nil {
		for rows.Next() {
			var attachment EmailAttachment
			rows.Scan(&attachment.ID, &attachment.EmailID, &attachment.Filename,
				&attachment.ContentType, &attachment.Size, &attachment.CreatedAt)
			recent = append(recent, attachment)
		}
		rows.Close()
	}
	stats["recent_attachments"] = recent

	return stats, nil
}

// Helper methods

func (eas *EmailAttachmentService) validateAttachment(upload *AttachmentUpload) error {
	// Check file size
	if upload.Size > eas.maxSize {
		return fmt.Errorf("attachment size %d bytes exceeds maximum allowed size %d bytes", upload.Size, eas.maxSize)
	}

	// Check content type
	allowed := false
	for _, allowedType := range eas.allowedTypes {
		if upload.ContentType == allowedType || strings.HasPrefix(upload.ContentType, allowedType) {
			allowed = true
			break
		}
	}
	if !allowed {
		return fmt.Errorf("content type %s is not allowed", upload.ContentType)
	}

	// Validate filename
	if upload.Filename == "" {
		return fmt.Errorf("filename is required")
	}

	// Check for dangerous filenames
	dangerousPatterns := []string{"../", "..\\", "<", ">", "|", "*", "?"}
	for _, pattern := range dangerousPatterns {
		if strings.Contains(upload.Filename, pattern) {
			return fmt.Errorf("filename contains dangerous characters")
		}
	}

	return nil
}

func (eas *EmailAttachmentService) generateHash(data []byte) string {
	hash := sha256.Sum256(data)
	return base64.StdEncoding.EncodeToString(hash[:])
}

func (eas *EmailAttachmentService) generateStoragePath(emailID uuid.UUID, filename string) string {
	timestamp := time.Now().Format("20060102_150405")
	return fmt.Sprintf("attachments/%s/%s_%s", emailID.String()[:8], timestamp, filename)
}

func (eas *EmailAttachmentService) saveAttachmentToDB(attachment *EmailAttachment) error {
	query := `
		INSERT INTO email_attachments (id, email_id, filename, content_type, size, hash,
			storage_path, is_inline, content_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := eas.db.Exec(query,
		attachment.ID, attachment.EmailID, attachment.Filename, attachment.ContentType,
		attachment.Size, attachment.Hash, attachment.StoragePath, attachment.IsInline,
		attachment.ContentID, attachment.CreatedAt,
	)

	return err
}

func (eas *EmailAttachmentService) scanAttachmentForViruses(attachmentID uuid.UUID, data []byte) {
	// Simulate virus scanning (would integrate with actual AV service)
	scanResult := &AttachmentVirusScan{
		AttachmentID: attachmentID,
		IsSafe:       true, // Assume safe for demo
		ScanEngine:   "simulated_scanner",
		ScanResult:   "clean",
		ScannedAt:    time.Now(),
	}

	// Save scan result
	query := `
		INSERT INTO attachment_virus_scans (attachment_id, is_safe, scan_engine, scan_result, scanned_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	eas.db.Exec(query, scanResult.AttachmentID, scanResult.IsSafe, scanResult.ScanEngine,
		scanResult.ScanResult, scanResult.ScannedAt)
}

func (eas *EmailAttachmentService) parseEmailParts(email *mail.Message) []interface{} {
	// Placeholder - would use proper email parsing
	return []interface{}{}
}

func (eas *EmailAttachmentService) isAttachmentPart(part interface{}) bool {
	// Placeholder - would check if email part is an attachment
	return false
}

func (eas *EmailAttachmentService) extractFilename(part interface{}) string {
	// Placeholder - would extract filename from email part
	return "attachment.bin"
}

func (eas *EmailAttachmentService) extractContentType(part interface{}) string {
	// Placeholder - would extract content type from email part
	return "application/octet-stream"
}

func (eas *EmailAttachmentService) extractPartData(part interface{}) []byte {
	// Placeholder - would extract data from email part
	return []byte{}
}

func (eas *EmailAttachmentService) isInlineAttachment(part interface{}) bool {
	// Placeholder - would check if attachment is inline
	return false
}

func (eas *EmailAttachmentService) extractContentID(part interface{}) string {
	// Placeholder - would extract content ID
	return ""
}

func (eas *EmailAttachmentService) generateThumbnail(data []byte, maxWidth, maxHeight int, contentType string) []byte {
	// Placeholder - would use image processing library to generate thumbnail
	// For now, return original data (not recommended for production)
	return data
}

// StorageService interface (would be implemented separately)
type StorageService struct{}

func (ss *StorageService) SaveFile(path string, data []byte) error {
	// Placeholder - would save to actual storage (S3, local filesystem, etc.)
	return nil
}

func (ss *StorageService) GetFile(path string) ([]byte, error) {
	// Placeholder - would retrieve from actual storage
	return []byte{}, nil
}

func (ss *StorageService) DeleteFile(path string) error {
	// Placeholder - would delete from actual storage
	return nil
}
