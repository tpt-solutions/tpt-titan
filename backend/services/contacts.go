package services

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"tpt-titan-simple/backend/models"
)

// ContactService handles contact-related business logic
type ContactService struct {
	db *sql.DB
}

// NewContactService creates a new contact service
func NewContactService(db *sql.DB) *ContactService {
	return &ContactService{db: db}
}

// GetContacts retrieves all contacts for a user
func (s *ContactService) GetContacts(userID uuid.UUID) ([]models.ContactResponse, error) {
	query := `
		SELECT id, first_name, last_name, email, phone, company, position, notes, created_at, updated_at
		FROM contacts
		WHERE user_id = $1
		ORDER BY COALESCE(first_name, '') || ' ' || COALESCE(last_name, '') ASC
	`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query contacts: %w", err)
	}
	defer rows.Close()

	var contacts []models.ContactResponse
	for rows.Next() {
		var contact models.Contact
		err := rows.Scan(
			&contact.ID, &contact.FirstName, &contact.LastName, &contact.Email,
			&contact.Phone, &contact.Company, &contact.Position, &contact.Notes,
			&contact.CreatedAt, &contact.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan contact: %w", err)
		}
		contacts = append(contacts, contact.ToResponse())
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating contacts: %w", err)
	}

	return contacts, nil
}

// GetContact retrieves a single contact by ID
func (s *ContactService) GetContact(userID, contactID uuid.UUID) (*models.ContactResponse, error) {
	query := `
		SELECT id, first_name, last_name, email, phone, company, position, notes, created_at, updated_at
		FROM contacts
		WHERE id = $1 AND user_id = $2
	`

	var contact models.Contact
	err := s.db.QueryRow(query, contactID, userID).Scan(
		&contact.ID, &contact.FirstName, &contact.LastName, &contact.Email,
		&contact.Phone, &contact.Company, &contact.Position, &contact.Notes,
		&contact.CreatedAt, &contact.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get contact: %w", err)
	}

	response := contact.ToResponse()
	return &response, nil
}

// CreateContact creates a new contact
func (s *ContactService) CreateContact(userID uuid.UUID, req models.ContactRequest) (*models.ContactResponse, error) {
	contactID := uuid.New()

	query := `
		INSERT INTO contacts (id, user_id, first_name, last_name, email, phone, company, position, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	now := time.Now()
	_, err := s.db.Exec(query,
		contactID, userID, req.FirstName, req.LastName, req.Email,
		req.Phone, req.Company, req.Position, req.Notes, now, now,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create contact: %w", err)
	}

	// Return the created contact
	return s.GetContact(userID, contactID)
}

// UpdateContact updates an existing contact
func (s *ContactService) UpdateContact(userID, contactID uuid.UUID, req models.ContactRequest) (*models.ContactResponse, error) {
	query := `
		UPDATE contacts
		SET first_name = $1, last_name = $2, email = $3, phone = $4,
		    company = $5, position = $6, notes = $7, updated_at = $8
		WHERE id = $9 AND user_id = $10
	`

	result, err := s.db.Exec(query,
		req.FirstName, req.LastName, req.Email, req.Phone,
		req.Company, req.Position, req.Notes, time.Now(),
		contactID, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update contact: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil, nil // Contact not found or doesn't belong to user
	}

	// Return the updated contact
	return s.GetContact(userID, contactID)
}

// DeleteContact deletes a contact
func (s *ContactService) DeleteContact(userID, contactID uuid.UUID) error {
	query := `DELETE FROM contacts WHERE id = $1 AND user_id = $2`

	result, err := s.db.Exec(query, contactID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete contact: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("contact not found or doesn't belong to user")
	}

	return nil
}

// SearchContacts searches contacts by name or email
func (s *ContactService) SearchContacts(userID uuid.UUID, query string) ([]models.ContactResponse, error) {
	if strings.TrimSpace(query) == "" {
		return s.GetContacts(userID)
	}

	searchTerm := "%" + strings.ToLower(query) + "%"
	sqlQuery := `
		SELECT id, first_name, last_name, email, phone, company, position, notes, created_at, updated_at
		FROM contacts
		WHERE user_id = $1 AND (
			LOWER(COALESCE(first_name, '') || ' ' || COALESCE(last_name, '')) LIKE $2 OR
			LOWER(COALESCE(email, '')) LIKE $2 OR
			LOWER(COALESCE(company, '')) LIKE $2
		)
		ORDER BY COALESCE(first_name, '') || ' ' || COALESCE(last_name, '') ASC
	`

	rows, err := s.db.Query(sqlQuery, userID, searchTerm)
	if err != nil {
		return nil, fmt.Errorf("failed to search contacts: %w", err)
	}
	defer rows.Close()

	var contacts []models.ContactResponse
	for rows.Next() {
		var contact models.Contact
		err := rows.Scan(
			&contact.ID, &contact.FirstName, &contact.LastName, &contact.Email,
			&contact.Phone, &contact.Company, &contact.Position, &contact.Notes,
			&contact.CreatedAt, &contact.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan contact: %w", err)
		}
		contacts = append(contacts, contact.ToResponse())
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating search results: %w", err)
	}

	return contacts, nil
}
