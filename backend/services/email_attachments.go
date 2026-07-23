package services

import (
	"bytes"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"io"
	"math"
	"mime"
	"mime/multipart"
	"net/mail"
	"net/textproto"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// EmailAttachmentService handles email attachments
type EmailAttachmentService struct {
	db           *sql.DB
	storageSvc   *StorageService // Would be implemented for file storage
	maxSize      int64           // Maximum attachment size in bytes
	allowedTypes []string        // Allowed MIME types
}

// EmailAttachment represents an email attachment
type EmailAttachment struct {
	ID          uuid.UUID `json:"id"`
	EmailID     uuid.UUID `json:"email_id"`
	Filename    string    `json:"filename"`
	ContentType string    `json:"content_type"`
	Size        int64     `json:"size"` // Size in bytes
	Hash        string    `json:"hash"` // SHA256 hash for integrity
	StoragePath string    `json:"storage_path"`
	IsInline    bool      `json:"is_inline"`            // True for inline images
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

// ProcessEmailAttachments processes attachments from incoming email.
func (eas *EmailAttachmentService) ProcessEmailAttachments(emailID uuid.UUID, email *mail.Message) ([]uuid.UUID, error) {
	var attachmentIDs []uuid.UUID

	// Walk the MIME tree starting from the top-level message. Attachment parts
	// are detected by their Content-Disposition / MIME type and extracted.
	var walk func(headers textproto.MIMEHeader, body io.Reader) error
	walk = func(headers textproto.MIMEHeader, body io.Reader) error {
		mediaType, params, err := mime.ParseMediaType(headers.Get("Content-Type"))
		if err != nil {
			mediaType = "text/plain"
		}

		if !strings.HasPrefix(mediaType, "multipart/") {
			// Non-multipart body: no nested attachments here.
			return nil
		}

		mr := multipart.NewReader(body, params["boundary"])
		for {
			part, perr := mr.NextPart()
			if perr != nil {
				break
			}
			ct := part.Header.Get("Content-Type")
			pmt, _, _ := mime.ParseMediaType(ct)
			if strings.HasPrefix(pmt, "multipart/") {
				if werr := walk(part.Header, part); werr != nil {
					return werr
				}
				continue
			}
			if isAttachmentHeader(part.Header) {
				data, rerr := io.ReadAll(part)
				if rerr != nil {
					return rerr
				}
				upload := &AttachmentUpload{
					Filename:    decodeFilename(part.FileName()),
					ContentType: pmt,
					Data:        data,
					IsInline:    isInlineHeader(part.Header),
					ContentID:   strings.Trim(part.Header.Get("Content-ID"), "<>"),
				}
				attachment, aerr := eas.SaveAttachment(emailID, upload)
				if aerr != nil {
					// Log error but continue processing other attachments.
					continue
				}
				attachmentIDs = append(attachmentIDs, attachment.ID)
			}
		}
		return nil
	}

	if werr := walk(textproto.MIMEHeader(email.Header), email.Body); werr != nil {
		return attachmentIDs, werr
	}

	return attachmentIDs, nil
}

// decodeFilename decodes RFC 2047 / parameter-encoded filenames when possible
// and strips path components for safety.
func decodeFilename(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "attachment.bin"
	}
	// decoder.DecodeHeader handles RFC 2047 encoded words.
	if decoded, err := mimeWordDecoder.DecodeHeader(name); err == nil && decoded != "" {
		name = decoded
	}
	// Strip directory components.
	if idx := strings.LastIndexAny(name, `/\`); idx >= 0 {
		name = name[idx+1:]
	}
	if name == "" {
		return "attachment.bin"
	}
	return name
}

var mimeWordDecoder = mime.WordDecoder{}

func isAttachmentHeader(h textproto.MIMEHeader) bool {
	disp := h.Get("Content-Disposition")
	ct := h.Get("Content-Type")
	mt, _, _ := mime.ParseMediaType(ct)
	isText := mt == "" || strings.HasPrefix(mt, "text/") || strings.HasPrefix(mt, "multipart/")

	switch {
	case disp == "":
		// No disposition: treat non-text (and non-multipart) parts as attachments.
		return !isText
	case disp == "attachment":
		return true
	case disp == "inline":
		// Inline text is message body, not an attachment; inline images/blobs are.
		return !isText
	}
	return !isText
}

func isInlineHeader(h textproto.MIMEHeader) bool {
	dm, _, _ := mime.ParseMediaType(h.Get("Content-Disposition"))
	return dm == "inline"
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

	// Generate thumbnail by decoding the image and scaling it down.
	thumbnail, err := generateImageThumbnail(data.Data, maxWidth, maxHeight)
	if err != nil {
		return nil, err
	}

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

// generateImageThumbnail decodes an image and scales it down to fit within
// maxWidth x maxHeight (preserving aspect ratio), returning JPEG-encoded bytes.
func generateImageThumbnail(data []byte, maxWidth, maxHeight int) ([]byte, error) {
	if maxWidth <= 0 {
		maxWidth = 200
	}
	if maxHeight <= 0 {
		maxHeight = 200
	}

	src, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image for thumbnail: %w", err)
	}

	b := src.Bounds()
	w, h := b.Dx(), b.Dy()
	if w == 0 || h == 0 {
		return nil, fmt.Errorf("invalid image dimensions")
	}

	scale := math.Min(float64(maxWidth)/float64(w), float64(maxHeight)/float64(h))
	if scale > 1 {
		scale = 1
	}
	nw, nh := int(float64(w)*scale), int(float64(h)*scale)
	if nw < 1 {
		nw = 1
	}
	if nh < 1 {
		nh = 1
	}

	dst := image.NewRGBA(image.Rect(0, 0, nw, nh))
	for y := 0; y < nh; y++ {
		for x := 0; x < nw; x++ {
			srcX := b.Min.X + int(float64(x)/scale)
			srcY := b.Min.Y + int(float64(y)/scale)
			dst.Set(x, y, src.At(srcX, srcY))
		}
	}

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, dst, &jpeg.Options{Quality: 85}); err != nil {
		return nil, fmt.Errorf("failed to encode thumbnail: %w", err)
	}
	return buf.Bytes(), nil
}

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

// StorageService persists attachment bytes. The default implementation stores
// files on the local filesystem under a configurable base directory; swap in an
// S3/GCS-backed implementation by satisfying this interface.
type StorageService struct {
	baseDir string
}

// NewLocalStorageService creates a filesystem-backed storage service rooted at
// baseDir. The directory is created if it does not already exist.
func NewLocalStorageService(baseDir string) (*StorageService, error) {
	if baseDir == "" {
		baseDir = "data/attachments"
	}
	if err := os.MkdirAll(baseDir, 0o750); err != nil {
		return nil, fmt.Errorf("failed to create storage dir: %w", err)
	}
	return &StorageService{baseDir: baseDir}, nil
}

func (ss *StorageService) resolve(path string) string {
	// Prevent path traversal outside the base directory.
	clean := filepath.Clean(path)
	if strings.Contains(clean, "..") {
		clean = filepath.Base(clean)
	}
	return filepath.Join(ss.baseDir, clean)
}

func (ss *StorageService) SaveFile(path string, data []byte) error {
	full := ss.resolve(path)
	if err := os.MkdirAll(filepath.Dir(full), 0o750); err != nil {
		return err
	}
	return os.WriteFile(full, data, 0o600)
}

func (ss *StorageService) GetFile(path string) ([]byte, error) {
	full := ss.resolve(path)
	data, err := os.ReadFile(full)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (ss *StorageService) DeleteFile(path string) error {
	full := ss.resolve(path)
	err := os.Remove(full)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// scanAttachmentForViruses scans attachment data for malware. If the ClamAV
// daemon (clamdscan / clamd) is available on PATH it is used for a real scan;
// otherwise the attachment is marked scanned-but-unverified with an honest
// engine label rather than claimed to be clean.
func (eas *EmailAttachmentService) scanAttachmentForViruses(attachmentID uuid.UUID, data []byte) {
	engine := "unverified"
	scanResult := "unscanned"
	isSafe := false

	if eas.storageSvc != nil {
		// Best-effort real scan via clamdscan if present on the system.
		if path, err := eas.writeTempForScan(data); err == nil {
			defer os.Remove(path)
			if out, serr := exec.Command("clamdscan", "--no-summary", path).CombinedOutput(); serr == nil || strings.Contains(string(out), "FOUND") {
				engine = "clamav"
				scanResult = "clamav:" + strings.TrimSpace(string(out))
				isSafe = !strings.Contains(string(out), "FOUND")
			} else {
				scanResult = "clamav unavailable: " + strings.TrimSpace(string(out))
			}
		}
	}

	result := &AttachmentVirusScan{
		AttachmentID: attachmentID,
		IsSafe:       isSafe,
		ScanEngine:   engine,
		ScanResult:   scanResult,
		ScannedAt:    time.Now(),
	}

	query := `
		INSERT INTO attachment_virus_scans (attachment_id, is_safe, scan_engine, scan_result, scanned_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	eas.db.Exec(query, result.AttachmentID, result.IsSafe, result.ScanEngine,
		result.ScanResult, result.ScannedAt)
}

func (eas *EmailAttachmentService) writeTempForScan(data []byte) (string, error) {
	f, err := os.CreateTemp("", "attachment-scan-*.bin")
	if err != nil {
		return "", err
	}
	if _, err := f.Write(data); err != nil {
		f.Close()
		os.Remove(f.Name())
		return "", err
	}
	f.Close()
	return f.Name(), nil
}
