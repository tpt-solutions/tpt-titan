package models

import (
	"time"

	"github.com/google/uuid"
)

// ContactEmail represents an email address for a contact
type ContactEmail struct {
	Type  string `json:"type"` // "work", "personal", "other"
	Value string `json:"value"`
	Label string `json:"label,omitempty"` // Custom label for "other" type
}

// ContactPhone represents a phone number for a contact
type ContactPhone struct {
	Type  string `json:"type"` // "work", "mobile", "home", "other"
	Value string `json:"value"`
	Label string `json:"label,omitempty"` // Custom label for "other" type
}

// Contact represents a user's contact information
type Contact struct {
	ID        uuid.UUID      `json:"id" db:"id"`
	UserID    uuid.UUID      `json:"user_id" db:"user_id"`
	FirstName *string        `json:"first_name,omitempty" db:"first_name"`
	LastName  *string        `json:"last_name,omitempty" db:"last_name"`
	Email     *string        `json:"email,omitempty" db:"email"`
	Phone     *string        `json:"phone,omitempty" db:"phone"`
	Emails    []ContactEmail `json:"emails,omitempty" db:"emails"`
	Phones    []ContactPhone `json:"phones,omitempty" db:"phones"`
	Company   *string        `json:"company,omitempty" db:"company"`
	Position  *string        `json:"position,omitempty" db:"position"`
	Notes     *string        `json:"notes,omitempty" db:"notes"`
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" db:"updated_at"`
}

// ContactRequest represents the request payload for creating/updating contacts
type ContactRequest struct {
	FirstName *string        `json:"first_name,omitempty"`
	LastName  *string        `json:"last_name,omitempty"`
	Email     *string        `json:"email,omitempty"`
	Phone     *string        `json:"phone,omitempty"`
	Emails    []ContactEmail `json:"emails,omitempty"`
	Phones    []ContactPhone `json:"phones,omitempty"`
	Company   *string        `json:"company,omitempty"`
	Position  *string        `json:"position,omitempty"`
	Notes     *string        `json:"notes,omitempty"`
}

// ContactResponse represents the response payload for contacts
type ContactResponse struct {
	ID        uuid.UUID      `json:"id"`
	FirstName *string        `json:"first_name,omitempty"`
	LastName  *string        `json:"last_name,omitempty"`
	Emails    []ContactEmail `json:"emails,omitempty"`
	Phones    []ContactPhone `json:"phones,omitempty"`
	Company   *string        `json:"company,omitempty"`
	Position  *string        `json:"position,omitempty"`
	Notes     *string        `json:"notes,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// GetFullName returns the full name of the contact
func (c *Contact) GetFullName() string {
	if c.FirstName != nil && c.LastName != nil {
		return *c.FirstName + " " + *c.LastName
	} else if c.FirstName != nil {
		return *c.FirstName
	} else if c.LastName != nil {
		return *c.LastName
	}
	return "Unknown Contact"
}

// ToResponse converts a Contact to ContactResponse
func (c *Contact) ToResponse() ContactResponse {
	return ContactResponse{
		ID:        c.ID,
		FirstName: c.FirstName,
		LastName:  c.LastName,
		Emails:    c.Emails,
		Phones:    c.Phones,
		Company:   c.Company,
		Position:  c.Position,
		Notes:     c.Notes,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}
