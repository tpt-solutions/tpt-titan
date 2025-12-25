package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// EncryptedDocument represents encrypted document storage
type EncryptedDocument struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Title       string    `gorm:"size:255" json:"title"`
	ContentType string    `gorm:"size:50;not null" json:"content_type"` // spreadsheet, document, form, etc.
	EncryptedData []byte  `gorm:"type:bytea;not null" json:"encrypted_data"`
	Salt        []byte    `gorm:"type:bytea;not null" json:"salt"`
	Algorithm   string    `gorm:"size:50;not null;default:'AES-256-GCM'" json:"algorithm"`
	FileSize    int64     `gorm:"default:0" json:"file_size"`
	Version     int       `gorm:"default:1" json:"version"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// EncryptedForm represents encrypted form data
type EncryptedForm struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID         uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Name           string    `gorm:"size:255;not null" json:"name"`
	Description    string    `gorm:"type:text" json:"description"`
	EncryptedSchema []byte `gorm:"type:bytea;not null" json:"encrypted_schema"` // Form structure
	Salt           []byte    `gorm:"type:bytea;not null" json:"salt"`
	Algorithm      string    `gorm:"size:50;not null;default:'AES-256-GCM'" json:"algorithm"`
	ResponseCount  int       `gorm:"default:0" json:"response_count"`
	Status         string    `gorm:"size:20;default:'draft'" json:"status"` // draft, active, archived
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// EncryptedFormResponse represents encrypted form response data
type EncryptedFormResponse struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	FormID    uuid.UUID `gorm:"type:uuid;not null" json:"form_id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"` // Respondent user ID
	EncryptedData []byte `gorm:"type:bytea;not null" json:"encrypted_data"`
	Salt      []byte     `gorm:"type:bytea;not null" json:"salt"`
	Algorithm string     `gorm:"size:50;not null;default:'AES-256-GCM'" json:"algorithm"`
	SubmittedAt time.Time `json:"submitted_at"`
}

// EncryptedEmail represents encrypted email storage
type EncryptedEmail struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID     uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Folder     string    `gorm:"size:50;default:'inbox'" json:"folder"`
	From       string    `gorm:"size:255;not null" json:"from"`
	To         string    `gorm:"type:text;not null" json:"to"`
	Subject    string    `gorm:"size:500" json:"subject"`
	EncryptedBody []byte `gorm:"type:bytea;not null" json:"encrypted_body"`
	Salt       []byte    `gorm:"type:bytea;not null" json:"salt"`
	Algorithm  string    `gorm:"size:50;not null;default:'AES-256-GCM'" json:"algorithm"`
	IsRead     bool      `gorm:"default:false" json:"is_read"`
	IsEncrypted bool     `gorm:"default:true" json:"is_encrypted"`
	SentAt     time.Time `json:"sent_at"`
	ReceivedAt time.Time `json:"received_at"`
}

// EncryptedTask represents encrypted task data
type EncryptedTask struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Title       string    `gorm:"size:255;not null" json:"title"`
	EncryptedDescription []byte `gorm:"type:bytea" json:"encrypted_description"`
	Salt        []byte    `gorm:"type:bytea;not null" json:"salt"`
	Algorithm   string    `gorm:"size:50;not null;default:'AES-256-GCM'" json:"algorithm"`
	Status      string    `gorm:"size:20;default:'todo'" json:"status"`
	Priority    string    `gorm:"size:10;default:'medium'" json:"priority"`
	DueDate     *time.Time `json:"due_date"`
	AssignedTo  *uuid.UUID `gorm:"type:uuid" json:"assigned_to"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// KeyBackup represents encrypted key backup information
type KeyBackup struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID        uuid.UUID `gorm:"type:uuid;not null;unique" json:"user_id"`
	BackupMethod  string    `gorm:"size:50;not null" json:"backup_method"` // shamir, hardware, recovery_codes
	EncryptedData []byte    `gorm:"type:bytea;not null" json:"encrypted_data"` // Encrypted backup data
	Salt          []byte    `gorm:"type:bytea;not null" json:"salt"`
	Algorithm     string    `gorm:"size:50;not null;default:'AES-256-GCM'" json:"algorithm"`
	CreatedAt     time.Time `json:"created_at"`
	LastUpdated   time.Time `json:"last_updated"`
}

// RecoveryShare represents individual Shamir shares for key recovery
type RecoveryShare struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID     uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	ShareIndex int       `gorm:"not null" json:"share_index"`
	TotalShares int      `gorm:"not null" json:"total_shares"`
	Threshold  int       `gorm:"not null" json:"threshold"`
	EncryptedShare []byte `gorm:"type:bytea;not null" json:"encrypted_share"`
	Salt       []byte     `gorm:"type:bytea;not null" json:"salt"`
	Algorithm  string     `gorm:"size:50;not null;default:'AES-256-GCM'" json:"algorithm"`
	GuardianName string   `gorm:"size:255" json:"guardian_name"` // Who holds this share
	GuardianEmail string  `gorm:"size:255" json:"guardian_email"`
	Status     string     `gorm:"size:20;default:'active'" json:"status"` // active, used, revoked
	CreatedAt  time.Time  `json:"created_at"`
}

// HardwareKey represents hardware security keys for recovery
type HardwareKey struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID       uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	DeviceID     string    `gorm:"size:255;not null;unique" json:"device_id"`
	DeviceType   string    `gorm:"size:50;not null" json:"device_type"` // usb, yubikey, etc.
	PublicKey    []byte    `gorm:"type:bytea;not null" json:"public_key"`
	EncryptedKey []byte    `gorm:"type:bytea;not null" json:"encrypted_key"`
	Salt         []byte    `gorm:"type:bytea;not null" json:"salt"`
	Algorithm    string    `gorm:"size:50;not null;default:'AES-256-GCM'" json:"algorithm"`
	FaceTemplate []byte    `gorm:"type:bytea" json:"face_template"` // For biometric recovery
	GPSLocation  string    `gorm:"size:255" json:"gps_location"` // Optional location lock
	TimeLock     *time.Time `json:"time_lock"` // Optional time-based access
	Status       string     `gorm:"size:20;default:'active'" json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	LastUsed     *time.Time `json:"last_used"`
}

// RecoveryAttempt logs recovery attempts for security monitoring
type RecoveryAttempt struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Method    string    `gorm:"size:50;not null" json:"method"` // shamir, hardware, recovery_codes
	IPAddress string    `gorm:"size:45" json:"ip_address"`
	UserAgent string    `gorm:"size:500" json:"user_agent"`
	Success   bool      `gorm:"default:false" json:"success"`
	ErrorMessage string `gorm:"type:text" json:"error_message"`
	AttemptedAt time.Time `json:"attempted_at"`
}

// TableName overrides for better naming
func (EncryptedDocument) TableName() string { return "encrypted_documents" }
func (EncryptedForm) TableName() string { return "encrypted_forms" }
func (EncryptedFormResponse) TableName() string { return "encrypted_form_responses" }
func (EncryptedEmail) TableName() string { return "encrypted_emails" }
func (EncryptedTask) TableName() string { return "encrypted_tasks" }
func (KeyBackup) TableName() string { return "key_backups" }
func (RecoveryShare) TableName() string { return "recovery_shares" }
func (HardwareKey) TableName() string { return "hardware_keys" }
func (RecoveryAttempt) TableName() string { return "recovery_attempts" }
