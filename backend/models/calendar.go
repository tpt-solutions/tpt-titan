package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Calendar represents a user's calendar
type Calendar struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	UserID      uuid.UUID  `json:"user_id" db:"user_id"`
	Name        string     `json:"name" db:"name"`
	Description *string    `json:"description,omitempty" db:"description"`
	Color       string     `json:"color" db:"color"`
	IsDefault   bool       `json:"is_default" db:"is_default"`
	IsShared    bool       `json:"is_shared" db:"is_shared"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

// CalendarRequest represents the request payload for creating/updating calendars
type CalendarRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description,omitempty"`
	Color       string  `json:"color"`
	IsDefault   bool    `json:"is_default"`
}

// CalendarResponse represents the response payload for calendars
type CalendarResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	Color       string    `json:"color"`
	IsDefault   bool      `json:"is_default"`
	IsShared    bool      `json:"is_shared"`
	EventCount  int       `json:"event_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Event represents a calendar event
type Event struct {
	ID               uuid.UUID  `json:"id" db:"id"`
	CalendarID       uuid.UUID  `json:"calendar_id" db:"calendar_id"`
	UserID           uuid.UUID  `json:"user_id" db:"user_id"`
	Title            string     `json:"title" db:"title"`
	Description      *string    `json:"description,omitempty" db:"description"`
	Location         *string    `json:"location,omitempty" db:"location"`
	StartTime        time.Time  `json:"start_time" db:"start_time"`
	EndTime          time.Time  `json:"end_time" db:"end_time"`
	IsAllDay         bool       `json:"is_all_day" db:"is_all_day"`
	RecurrenceRule   *string    `json:"recurrence_rule,omitempty" db:"recurrence_rule"`
	RecurrenceID     *uuid.UUID `json:"recurrence_id,omitempty" db:"recurrence_id"`
	ReminderMinutes  int        `json:"reminder_minutes" db:"reminder_minutes"`
	IsCompleted      bool       `json:"is_completed" db:"is_completed"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
}

// EventRequest represents the request payload for creating/updating events
type EventRequest struct {
	CalendarID      uuid.UUID  `json:"calendar_id" binding:"required"`
	Title           string     `json:"title" binding:"required"`
	Description     *string    `json:"description,omitempty"`
	Location        *string    `json:"location,omitempty"`
	StartTime       time.Time  `json:"start_time" binding:"required"`
	EndTime         time.Time  `json:"end_time" binding:"required"`
	IsAllDay        bool       `json:"is_all_day"`
	RecurrenceRule  *string    `json:"recurrence_rule,omitempty"`
	ReminderMinutes int        `json:"reminder_minutes"`
	AttendeeIDs     []uuid.UUID `json:"attendee_ids,omitempty"`
}

// EventResponse represents the response payload for events
type EventResponse struct {
	ID               uuid.UUID             `json:"id"`
	CalendarID       uuid.UUID             `json:"calendar_id"`
	CalendarName     string                `json:"calendar_name"`
	CalendarColor    string                `json:"calendar_color"`
	Title            string                `json:"title"`
	Description      *string               `json:"description,omitempty"`
	Location         *string               `json:"location,omitempty"`
	StartTime        time.Time             `json:"start_time"`
	EndTime          time.Time             `json:"end_time"`
	IsAllDay         bool                  `json:"is_all_day"`
	RecurrenceRule   *string               `json:"recurrence_rule,omitempty"`
	RecurrenceID     *uuid.UUID            `json:"recurrence_id,omitempty"`
	ReminderMinutes  int                   `json:"reminder_minutes"`
	IsCompleted      bool                  `json:"is_completed"`
	Attendees        []EventAttendeeResponse `json:"attendees,omitempty"`
	CreatedAt        time.Time             `json:"created_at"`
	UpdatedAt        time.Time             `json:"updated_at"`
}

// EventAttendee represents an event attendee
type EventAttendee struct {
	ID        uuid.UUID `json:"id" db:"id"`
	EventID   uuid.UUID `json:"event_id" db:"event_id"`
	ContactID uuid.UUID `json:"contact_id" db:"contact_id"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// EventAttendeeResponse represents the response payload for event attendees
type EventAttendeeResponse struct {
	ID        uuid.UUID `json:"id"`
	EventID   uuid.UUID `json:"event_id"`
	ContactID uuid.UUID `json:"contact_id"`
	Name      string    `json:"name"`
	Email     *string   `json:"email,omitempty"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// ToResponse converts a Calendar to CalendarResponse
func (c *Calendar) ToResponse() CalendarResponse {
	return CalendarResponse{
		ID:         c.ID,
		Name:       c.Name,
		Description: c.Description,
		Color:      c.Color,
		IsDefault:  c.IsDefault,
		IsShared:   c.IsShared,
		EventCount: 0, // This will be populated by the service
		CreatedAt:  c.CreatedAt,
		UpdatedAt:  c.UpdatedAt,
	}
}

// ToResponse converts an Event to EventResponse
func (e *Event) ToResponse(calendarName, calendarColor string) EventResponse {
	return EventResponse{
		ID:              e.ID,
		CalendarID:      e.CalendarID,
		CalendarName:    calendarName,
		CalendarColor:   calendarColor,
		Title:           e.Title,
		Description:     e.Description,
		Location:        e.Location,
		StartTime:       e.StartTime,
		EndTime:         e.EndTime,
		IsAllDay:        e.IsAllDay,
		RecurrenceRule:  e.RecurrenceRule,
		RecurrenceID:    e.RecurrenceID,
		ReminderMinutes: e.ReminderMinutes,
		IsCompleted:     e.IsCompleted,
		Attendees:       []EventAttendeeResponse{}, // This will be populated by the service
		CreatedAt:       e.CreatedAt,
		UpdatedAt:       e.UpdatedAt,
	}
}

// ToResponse converts an EventAttendee to EventAttendeeResponse
func (ea *EventAttendee) ToResponse(contactName string, contactEmail *string) EventAttendeeResponse {
	return EventAttendeeResponse{
		ID:        ea.ID,
		EventID:   ea.EventID,
		ContactID: ea.ContactID,
		Name:      contactName,
		Email:     contactEmail,
		Status:    ea.Status,
		CreatedAt: ea.CreatedAt,
	}
}

// Validate checks if the event has valid times
func (e *Event) Validate() error {
	if e.EndTime.Before(e.StartTime) {
		return fmt.Errorf("end time cannot be before start time")
	}
	return nil
}
