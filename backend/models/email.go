package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// EmailAccount represents an external email account (IMAP/SMTP)
type EmailAccount struct {
	ID                uuid.UUID `json:"id" db:"id"`
	UserID            uuid.UUID `json:"user_id" db:"user_id"`
	Email             string    `json:"email" db:"email"`
	Provider          string    `json:"provider" db:"provider"` // 'imap', 'smtp', etc.
	Server            string    `json:"server" db:"server"`
	Port              int       `json:"port" db:"port"`
	Username          string    `json:"username" db:"username"`
	PasswordEncrypted []byte    `json:"-" db:"password_encrypted"` // Not exposed in JSON
	UseSSL            bool      `json:"use_ssl" db:"use_ssl"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
}

// EmailAccountRequest represents the request payload for creating/updating email accounts
type EmailAccountRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Provider string `json:"provider" binding:"required"`
	Server   string `json:"server" binding:"required"`
	Port     int    `json:"port" binding:"required,min=1,max=65535"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	UseSSL   bool   `json:"use_ssl"`
}

// EmailAccountResponse represents the response payload for email accounts
type EmailAccountResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Provider  string    `json:"provider"`
	Server    string    `json:"server"`
	Port      int       `json:"port"`
	Username  string    `json:"username"`
	UseSSL    bool      `json:"use_ssl"`
	CreatedAt time.Time `json:"created_at"`
}

// Email represents an email message
type Email struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	AccountID       uuid.UUID  `json:"account_id" db:"account_id"`
	MessageID       *string    `json:"message_id,omitempty" db:"message_id"`
	Subject         *string    `json:"subject,omitempty" db:"subject"`
	SenderName      *string    `json:"sender_name,omitempty" db:"sender_name"`
	SenderEmail     string     `json:"sender_email" db:"sender_email"`
	RecipientEmails []string   `json:"recipient_emails" db:"recipient_emails"`
	CCEmails        []string   `json:"cc_emails,omitempty" db:"cc_emails"`
	BCCEmails       []string   `json:"bcc_emails,omitempty" db:"bcc_emails"`
	Content         *string    `json:"content,omitempty" db:"content"`
	HTMLContent     *string    `json:"html_content,omitempty" db:"html_content"`
	ReceivedAt      *time.Time `json:"received_at,omitempty" db:"received_at"`
	SentAt          *time.Time `json:"sent_at,omitempty" db:"sent_at"`
	IsRead          bool       `json:"is_read" db:"is_read"`
	IsStarred       bool       `json:"is_starred" db:"is_starred"`
	Folder          string     `json:"folder" db:"folder"`
	Labels          []string   `json:"labels,omitempty" db:"labels"`
	HasAttachments  bool       `json:"-" db:"has_attachments"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
}

// EmailRequest represents the request payload for sending emails
type EmailRequest struct {
	AccountID       uuid.UUID `json:"account_id" binding:"required"`
	Subject         string    `json:"subject"`
	Content         string    `json:"content" binding:"required"`
	HTMLContent     *string   `json:"html_content,omitempty"`
	RecipientEmails []string  `json:"recipient_emails" binding:"required,min=1"`
	CCEmails        []string  `json:"cc_emails,omitempty"`
	BCCEmails       []string  `json:"bcc_emails,omitempty"`
}

// EmailResponse represents the response payload for emails
type EmailResponse struct {
	ID              uuid.UUID  `json:"id"`
	AccountID       uuid.UUID  `json:"account_id"`
	MessageID       *string    `json:"message_id,omitempty"`
	Subject         *string    `json:"subject,omitempty"`
	SenderName      *string    `json:"sender_name,omitempty"`
	SenderEmail     string     `json:"sender_email"`
	RecipientEmails []string   `json:"recipient_emails"`
	CCEmails        []string   `json:"cc_emails,omitempty"`
	BCCEmails       []string   `json:"bcc_emails,omitempty"`
	Content         *string    `json:"content,omitempty"`
	HTMLContent     *string    `json:"html_content,omitempty"`
	ReceivedAt      *time.Time `json:"received_at,omitempty"`
	SentAt          *time.Time `json:"sent_at,omitempty"`
	IsRead          bool       `json:"is_read"`
	IsStarred       bool       `json:"is_starred"`
	Folder          string     `json:"folder"`
	Labels          []string   `json:"labels,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}

// EmailSummary represents a condensed version of an email for list views
type EmailSummary struct {
	ID             uuid.UUID  `json:"id"`
	Subject        *string    `json:"subject,omitempty"`
	SenderName     *string    `json:"sender_name,omitempty"`
	SenderEmail    string     `json:"sender_email"`
	ReceivedAt     *time.Time `json:"received_at,omitempty"`
	SentAt         *time.Time `json:"sent_at,omitempty"`
	IsRead         bool       `json:"is_read"`
	IsStarred      bool       `json:"is_starred"`
	Folder         string     `json:"folder"`
	HasAttachments bool       `json:"has_attachments"`
}

// EmailSearchRequest represents search parameters for emails
type EmailSearchRequest struct {
	Query     string     `json:"query,omitempty"`
	Folder    string     `json:"folder,omitempty"`
	IsRead    *bool      `json:"is_read,omitempty"`
	IsStarred *bool      `json:"is_starred,omitempty"`
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	Limit     int        `json:"limit,omitempty"`
	Offset    int        `json:"offset,omitempty"`
}

// ToResponse converts an EmailAccount to EmailAccountResponse
func (ea *EmailAccount) ToResponse() EmailAccountResponse {
	return EmailAccountResponse{
		ID:        ea.ID,
		Email:     ea.Email,
		Provider:  ea.Provider,
		Server:    ea.Server,
		Port:      ea.Port,
		Username:  ea.Username,
		UseSSL:    ea.UseSSL,
		CreatedAt: ea.CreatedAt,
	}
}

// ToResponse converts an Email to EmailResponse
func (e *Email) ToResponse() EmailResponse {
	return EmailResponse{
		ID:              e.ID,
		AccountID:       e.AccountID,
		MessageID:       e.MessageID,
		Subject:         e.Subject,
		SenderName:      e.SenderName,
		SenderEmail:     e.SenderEmail,
		RecipientEmails: e.RecipientEmails,
		CCEmails:        e.CCEmails,
		BCCEmails:       e.BCCEmails,
		Content:         e.Content,
		HTMLContent:     e.HTMLContent,
		ReceivedAt:      e.ReceivedAt,
		SentAt:          e.SentAt,
		IsRead:          e.IsRead,
		IsStarred:       e.IsStarred,
		Folder:          e.Folder,
		Labels:          e.Labels,
		CreatedAt:       e.CreatedAt,
	}
}

// ToSummary converts an Email to EmailSummary
func (e *Email) ToSummary() EmailSummary {
	return EmailSummary{
		ID:             e.ID,
		Subject:        e.Subject,
		SenderName:     e.SenderName,
		SenderEmail:    e.SenderEmail,
		ReceivedAt:     e.ReceivedAt,
		SentAt:         e.SentAt,
		IsRead:         e.IsRead,
		IsStarred:      e.IsStarred,
		Folder:         e.Folder,
		HasAttachments: e.HasAttachments,
	}
}

// Validate checks if the email account has valid configuration
func (ea *EmailAccount) Validate() error {
	if ea.Email == "" {
		return fmt.Errorf("email address is required")
	}
	if ea.Server == "" {
		return fmt.Errorf("server address is required")
	}
	if ea.Port < 1 || ea.Port > 65535 {
		return fmt.Errorf("invalid port number")
	}
	if ea.Username == "" {
		return fmt.Errorf("username is required")
	}
	return nil
}

// GetDisplayName returns a display name for the email account
func (ea *EmailAccount) GetDisplayName() string {
	if ea.Username != "" && ea.Username != ea.Email {
		return ea.Username + " (" + ea.Email + ")"
	}
	return ea.Email
}
