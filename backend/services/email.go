package services

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	gomessagemail "github.com/emersion/go-message/mail"
	"github.com/go-mail/mail/v2"
	"github.com/google/uuid"
	"tpt-titan/backend/models"
	"tpt-titan/backend/utils"
)

// EmailService handles email-related business logic
type EmailService struct {
	db *sql.DB
}

// NewEmailService creates a new email service
func NewEmailService(db *sql.DB) *EmailService {
	return &EmailService{db: db}
}

// GetEmailAccounts retrieves all email accounts for a user
func (s *EmailService) GetEmailAccounts(userID uuid.UUID) ([]models.EmailAccountResponse, error) {
	query := `SELECT id, email, provider, server, port, username, use_ssl, created_at FROM email_accounts WHERE user_id = $1 ORDER BY created_at DESC`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query email accounts: %w", err)
	}
	defer rows.Close()

	var accounts []models.EmailAccountResponse
	for rows.Next() {
		var account models.EmailAccount
		err := rows.Scan(
			&account.ID, &account.Email, &account.Provider, &account.Server,
			&account.Port, &account.Username, &account.UseSSL, &account.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan email account: %w", err)
		}
		accounts = append(accounts, account.ToResponse())
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating email accounts: %w", err)
	}

	return accounts, nil
}

// GetEmailAccount retrieves a single email account by ID
func (s *EmailService) GetEmailAccount(userID, accountID uuid.UUID) (*models.EmailAccountResponse, error) {
	query := `SELECT id, email, provider, server, port, username, use_ssl, created_at FROM email_accounts WHERE id = $1 AND user_id = $2`

	var account models.EmailAccount
	err := s.db.QueryRow(query, accountID, userID).Scan(
		&account.ID, &account.Email, &account.Provider, &account.Server,
		&account.Port, &account.Username, &account.UseSSL, &account.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get email account: %w", err)
	}

	response := account.ToResponse()
	return &response, nil
}

// CreateEmailAccount creates a new email account
func (s *EmailService) CreateEmailAccount(userID uuid.UUID, req models.EmailAccountRequest) (*models.EmailAccountResponse, error) {
	// Validate the request
	account := models.EmailAccount{
		UserID:   userID,
		Email:    req.Email,
		Provider: req.Provider,
		Server:   req.Server,
		Port:     req.Port,
		Username: req.Username,
		UseSSL:   req.UseSSL,
	}
	if err := account.Validate(); err != nil {
		return nil, fmt.Errorf("invalid email account: %w", err)
	}

	// Encrypt the password
	encryptedPassword, err := utils.EncryptPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt password: %w", err)
	}

	accountID := uuid.New()

	query := `
		INSERT INTO email_accounts (id, user_id, email, provider, server, port, username, password_encrypted, use_ssl, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	now := time.Now()
	_, err = s.db.Exec(query,
		accountID, userID, req.Email, req.Provider, req.Server, req.Port,
		req.Username, encryptedPassword, req.UseSSL, now,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create email account: %w", err)
	}

	// Return the created account
	return s.GetEmailAccount(userID, accountID)
}

// UpdateEmailAccount updates an existing email account
func (s *EmailService) UpdateEmailAccount(userID, accountID uuid.UUID, req models.EmailAccountRequest) (*models.EmailAccountResponse, error) {
	// Validate the request
	account := models.EmailAccount{
		UserID:   userID,
		Email:    req.Email,
		Provider: req.Provider,
		Server:   req.Server,
		Port:     req.Port,
		Username: req.Username,
		UseSSL:   req.UseSSL,
	}
	if err := account.Validate(); err != nil {
		return nil, fmt.Errorf("invalid email account: %w", err)
	}

	// Check if account exists and belongs to user
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM email_accounts WHERE id = $1 AND user_id = $2)", accountID, userID).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("failed to check account existence: %w", err)
	}
	if !exists {
		return nil, nil
	}

	// Encrypt the password if provided
	var encryptedPassword []byte
	if req.Password != "" {
		encryptedPassword, err = utils.EncryptPassword(req.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt password: %w", err)
		}
	}

	query := `
		UPDATE email_accounts
		SET email = $1, provider = $2, server = $3, port = $4, username = $5, use_ssl = $6
	`
	args := []interface{}{req.Email, req.Provider, req.Server, req.Port, req.Username, req.UseSSL}

	if len(encryptedPassword) > 0 {
		query += `, password_encrypted = $7`
		args = append(args, encryptedPassword)
	}

	query += ` WHERE id = $8 AND user_id = $9`
	args = append(args, accountID, userID)

	_, err = s.db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update email account: %w", err)
	}

	// Return the updated account
	return s.GetEmailAccount(userID, accountID)
}

// DeleteEmailAccount deletes an email account
func (s *EmailService) DeleteEmailAccount(userID, accountID uuid.UUID) error {
	query := `DELETE FROM email_accounts WHERE id = $1 AND user_id = $2`

	result, err := s.db.Exec(query, accountID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete email account: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("email account not found or doesn't belong to user")
	}

	return nil
}

// GetEmails retrieves emails with pagination and filtering
func (s *EmailService) GetEmails(userID uuid.UUID, searchReq models.EmailSearchRequest) ([]models.EmailSummary, error) {
	query := `
		SELECT e.id, e.subject, e.sender_name, e.sender_email, e.received_at, e.sent_at,
			   e.is_read, e.is_starred, e.folder
		FROM emails e
		JOIN email_accounts ea ON e.account_id = ea.id
		WHERE ea.user_id = $1
	`
	args := []interface{}{userID}
	argCount := 1

	// Add filters
	if searchReq.Folder != "" {
		argCount++
		query += fmt.Sprintf(" AND e.folder = $%d", argCount)
		args = append(args, searchReq.Folder)
	}

	if searchReq.IsRead != nil {
		argCount++
		query += fmt.Sprintf(" AND e.is_read = $%d", argCount)
		args = append(args, *searchReq.IsRead)
	}

	if searchReq.IsStarred != nil {
		argCount++
		query += fmt.Sprintf(" AND e.is_starred = $%d", argCount)
		args = append(args, *searchReq.IsStarred)
	}

	if searchReq.StartDate != nil {
		argCount++
		query += fmt.Sprintf(" AND COALESCE(e.received_at, e.sent_at) >= $%d", argCount)
		args = append(args, *searchReq.StartDate)
	}

	if searchReq.EndDate != nil {
		argCount++
		query += fmt.Sprintf(" AND COALESCE(e.received_at, e.sent_at) <= $%d", argCount)
		args = append(args, *searchReq.EndDate)
	}

	if searchReq.Query != "" {
		query += fmt.Sprintf(" AND (e.subject ILIKE $%d OR e.sender_name ILIKE $%d OR e.sender_email ILIKE $%d OR e.content ILIKE $%d)", argCount+1, argCount+2, argCount+3, argCount+4)
		searchTerm := "%" + searchReq.Query + "%"
		args = append(args, searchTerm, searchTerm, searchTerm, searchTerm)
		argCount += 4
	}

	// Add ordering and pagination
	query += " ORDER BY COALESCE(e.received_at, e.sent_at) DESC"

	if searchReq.Limit > 0 {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, searchReq.Limit)
	}

	if searchReq.Offset > 0 {
		argCount++
		query += fmt.Sprintf(" OFFSET $%d", argCount)
		args = append(args, searchReq.Offset)
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query emails: %w", err)
	}
	defer rows.Close()

	var emails []models.EmailSummary
	for rows.Next() {
		var email models.Email
		err := rows.Scan(
			&email.ID, &email.Subject, &email.SenderName, &email.SenderEmail,
			&email.ReceivedAt, &email.SentAt, &email.IsRead, &email.IsStarred, &email.Folder,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan email: %w", err)
		}
		emails = append(emails, email.ToSummary())
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating emails: %w", err)
	}

	return emails, nil
}

// GetEmail retrieves a single email by ID
func (s *EmailService) GetEmail(userID, emailID uuid.UUID) (*models.EmailResponse, error) {
	query := `
		SELECT e.id, e.account_id, e.message_id, e.subject, e.sender_name, e.sender_email,
			   e.recipient_emails, e.cc_emails, e.bcc_emails, e.content, e.html_content,
			   e.received_at, e.sent_at, e.is_read, e.is_starred, e.folder, e.labels, e.created_at
		FROM emails e
		JOIN email_accounts ea ON e.account_id = ea.id
		WHERE e.id = $1 AND ea.user_id = $2
	`

	var email models.Email
	err := s.db.QueryRow(query, emailID, userID).Scan(
		&email.ID, &email.AccountID, &email.MessageID, &email.Subject, &email.SenderName, &email.SenderEmail,
		&email.RecipientEmails, &email.CCEmails, &email.BCCEmails, &email.Content, &email.HTMLContent,
		&email.ReceivedAt, &email.SentAt, &email.IsRead, &email.IsStarred, &email.Folder, &email.Labels, &email.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get email: %w", err)
	}

	response := email.ToResponse()
	return &response, nil
}

// SendEmail sends an email via SMTP
func (s *EmailService) SendEmail(userID uuid.UUID, req models.EmailRequest) (*models.EmailResponse, error) {
	// Get account details
	var account models.EmailAccount
	var encryptedPassword []byte
	err := s.db.QueryRow(`
		SELECT id, email, server, port, username, password_encrypted, use_ssl
		FROM email_accounts WHERE id = $1 AND user_id = $2
	`, req.AccountID, userID).Scan(
		&account.ID, &account.Email, &account.Server, &account.Port,
		&account.Username, &encryptedPassword, &account.UseSSL,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("email account not found or doesn't belong to user")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get account details: %w", err)
	}

	// Decrypt the password
	password, err := utils.DecryptPassword(encryptedPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt password: %w", err)
	}

	// Create the email message
	m := mail.NewMessage()
	m.SetHeader("From", account.Email)
	m.SetHeader("To", req.RecipientEmails...)
	if len(req.CCEmails) > 0 {
		m.SetHeader("Cc", req.CCEmails...)
	}
	if len(req.BCCEmails) > 0 {
		m.SetHeader("Bcc", req.BCCEmails...)
	}
	m.SetHeader("Subject", req.Subject)

	// Set body
	if req.HTMLContent != nil && *req.HTMLContent != "" {
		m.SetBody("text/html", *req.HTMLContent)
		if req.Content != "" {
			m.AddAlternative("text/plain", req.Content)
		}
	} else {
		m.SetBody("text/plain", req.Content)
	}

	// Create SMTP dialer
	d := mail.NewDialer(account.Server, account.Port, account.Username, password)
	d.SSL = account.UseSSL

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		return nil, fmt.Errorf("failed to send email: %w", err)
	}

	emailID := uuid.New()
	now := time.Now()
	messageID := fmt.Sprintf("<%s@%s>", emailID.String(), strings.Split(account.Email, "@")[1])

	// Store the sent email in the database
	query := `
		INSERT INTO emails (id, account_id, message_id, subject, sender_name, sender_email, recipient_emails, cc_emails, bcc_emails, content, html_content, sent_at, folder, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`

	// Get sender name from users table
	var senderName *string
	err = s.db.QueryRow("SELECT username FROM users WHERE id = $1", userID).Scan(&senderName)
	if err != nil {
		senderName = nil // Continue without name if not found
	}

	_, err = s.db.Exec(query,
		emailID, req.AccountID, messageID, req.Subject, senderName, account.Email,
		req.RecipientEmails, req.CCEmails, req.BCCEmails, req.Content, req.HTMLContent,
		now, "sent", now,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to store sent email: %w", err)
	}

	// Return the stored email
	return s.GetEmail(userID, emailID)
}

// MarkEmailAsRead marks an email as read/unread
func (s *EmailService) MarkEmailAsRead(userID, emailID uuid.UUID, isRead bool) error {
	query := `
		UPDATE emails
		SET is_read = $1
		WHERE id = $2 AND account_id IN (SELECT id FROM email_accounts WHERE user_id = $3)
	`

	result, err := s.db.Exec(query, isRead, emailID, userID)
	if err != nil {
		return fmt.Errorf("failed to mark email as read: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("email not found or doesn't belong to user")
	}

	return nil
}

// StarEmail marks an email as starred/unstarred
func (s *EmailService) StarEmail(userID, emailID uuid.UUID, isStarred bool) error {
	query := `
		UPDATE emails
		SET is_starred = $1
		WHERE id = $2 AND account_id IN (SELECT id FROM email_accounts WHERE user_id = $3)
	`

	result, err := s.db.Exec(query, isStarred, emailID, userID)
	if err != nil {
		return fmt.Errorf("failed to star email: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("email not found or doesn't belong to user")
	}

	return nil
}

// MoveEmailToFolder moves an email to a different folder
func (s *EmailService) MoveEmailToFolder(userID, emailID uuid.UUID, folder string) error {
	query := `
		UPDATE emails
		SET folder = $1
		WHERE id = $2 AND account_id IN (SELECT id FROM email_accounts WHERE user_id = $3)
	`

	result, err := s.db.Exec(query, folder, emailID, userID)
	if err != nil {
		return fmt.Errorf("failed to move email: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("email not found or doesn't belong to user")
	}

	return nil
}

// SyncEmails syncs emails from external IMAP accounts
func (s *EmailService) SyncEmails(userID, accountID uuid.UUID) error {
	// Get account details
	var account models.EmailAccount
	var encryptedPassword []byte
	err := s.db.QueryRow(`
		SELECT id, email, server, port, username, password_encrypted, use_ssl
		FROM email_accounts WHERE id = $1 AND user_id = $2
	`, accountID, userID).Scan(
		&account.ID, &account.Email, &account.Server, &account.Port,
		&account.Username, &encryptedPassword, &account.UseSSL,
	)
	if err == sql.ErrNoRows {
		return fmt.Errorf("email account not found or doesn't belong to user")
	}
	if err != nil {
		return fmt.Errorf("failed to get account details: %w", err)
	}

	// Decrypt the password
	password, err := utils.DecryptPassword(encryptedPassword)
	if err != nil {
		return fmt.Errorf("failed to decrypt password: %w", err)
	}

	// Connect to IMAP server
	var c *client.Client
	if account.UseSSL {
		c, err = client.DialTLS(account.Server+":"+strconv.Itoa(account.Port), nil)
	} else {
		c, err = client.Dial(account.Server + ":" + strconv.Itoa(account.Port))
	}
	if err != nil {
		return fmt.Errorf("failed to connect to IMAP server: %w", err)
	}
	defer c.Logout()

	// Login
	if err := c.Login(account.Username, password); err != nil {
		return fmt.Errorf("failed to login to IMAP server: %w", err)
	}

	// Get INBOX mailbox
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		return fmt.Errorf("failed to select INBOX: %w", err)
	}

	// Get the last synced UID
	var lastUID uint32
	err = s.db.QueryRow(`
		SELECT COALESCE(MAX(CAST(message_id AS INTEGER)), 0)
		FROM emails
		WHERE account_id = $1 AND folder = 'inbox' AND message_id ~ '^\d+$'
	`, accountID).Scan(&lastUID)
	if err != nil {
		log.Printf("Failed to get last UID, starting from beginning: %v", err)
		lastUID = 0
	}

	// Fetch new messages
	if mbox.Messages > 0 {
		from := uint32(1)
		if lastUID > 0 {
			from = lastUID + 1
		}
		to := mbox.Messages

		if from <= to {
			seqset := new(imap.SeqSet)
			seqset.AddRange(from, to)

			messages := make(chan *imap.Message, 10)
			done := make(chan error, 1)

			go func() {
				done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope, imap.FetchBody}, messages)
			}()

			for msg := range messages {
				if err := s.processIMAPMessage(accountID, msg); err != nil {
					log.Printf("Failed to process message %d: %v", msg.SeqNum, err)
				}
			}

			if err := <-done; err != nil {
				return fmt.Errorf("failed to fetch messages: %w", err)
			}
		}
	}

	return nil
}

// processIMAPMessage processes a single IMAP message and stores it in the database
func (s *EmailService) processIMAPMessage(accountID uuid.UUID, msg *imap.Message) error {
	envelope := msg.Envelope

	// Extract sender information
	senderName := ""
	senderEmail := ""
	if len(envelope.From) > 0 {
		senderName = envelope.From[0].PersonalName
		senderEmail = envelope.From[0].Address()
	}

	// Extract recipients
	var toEmails []string
	for _, addr := range envelope.To {
		toEmails = append(toEmails, addr.Address())
	}

	var ccEmails []string
	for _, addr := range envelope.Cc {
		ccEmails = append(ccEmails, addr.Address())
	}

	var bccEmails []string
	for _, addr := range envelope.Bcc {
		bccEmails = append(bccEmails, addr.Address())
	}

	// Get message body
	var content, htmlContent *string
	r := msg.GetBody(&imap.BodySectionName{})
	if r != nil {
		mr, err := gomessagemail.CreateReader(r)
		if err != nil {
			return fmt.Errorf("failed to create mail reader: %w", err)
		}

		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				continue
			}

			switch h := p.Header.(type) {
			case *gomessagemail.InlineHeader:
				ct, _, _ := h.ContentType()
				b, _ := io.ReadAll(p.Body)
				text := string(b)

				if strings.Contains(ct, "text/html") {
					htmlContent = &text
				} else if strings.Contains(ct, "text/plain") {
					content = &text
				}
			}
		}
	}

	emailID := uuid.New()
	now := time.Now()

	query := `
		INSERT INTO emails (id, account_id, message_id, subject, sender_name, sender_email,
						   recipient_emails, cc_emails, bcc_emails, content, html_content,
						   received_at, folder, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		ON CONFLICT (account_id, message_id) DO NOTHING
	`

	_, err := s.db.Exec(query,
		emailID, accountID, strconv.Itoa(int(msg.Uid)), envelope.Subject,
		senderName, senderEmail, toEmails, ccEmails, bccEmails,
		content, htmlContent, envelope.Date, "inbox", now,
	)

	if err != nil {
		return fmt.Errorf("failed to store email: %w", err)
	}

	return nil
}

// GetEmailStats returns email statistics for a user
func (s *EmailService) GetEmailStats(userID uuid.UUID) (map[string]int, error) {
	query := `
		SELECT folder, COUNT(*) as count
		FROM emails e
		JOIN email_accounts ea ON e.account_id = ea.id
		WHERE ea.user_id = $1
		GROUP BY folder
	`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get email stats: %w", err)
	}
	defer rows.Close()

	stats := make(map[string]int)
	for rows.Next() {
		var folder string
		var count int
		err := rows.Scan(&folder, &count)
		if err != nil {
			return nil, fmt.Errorf("failed to scan email stats: %w", err)
		}
		stats[folder] = count
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating email stats: %w", err)
	}

	return stats, nil
}
