package services

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"tpt-titan/backend/models"
)

// CalendarService handles calendar-related business logic
type CalendarService struct {
	db *sql.DB
}

// NewCalendarService creates a new calendar service
func NewCalendarService(db *sql.DB) *CalendarService {
	return &CalendarService{db: db}
}

// GetCalendars retrieves all calendars for a user
func (s *CalendarService) GetCalendars(userID uuid.UUID) ([]models.CalendarResponse, error) {
	query := `
		SELECT c.id, c.name, c.description, c.color, c.is_default, c.is_shared, c.created_at, c.updated_at,
			   COUNT(e.id) as event_count
		FROM calendars c
		LEFT JOIN events e ON c.id = e.calendar_id AND e.user_id = c.user_id
		WHERE c.user_id = $1
		GROUP BY c.id, c.name, c.description, c.color, c.is_default, c.is_shared, c.created_at, c.updated_at
		ORDER BY c.is_default DESC, c.name ASC
	`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query calendars: %w", err)
	}
	defer rows.Close()

	var calendars []models.CalendarResponse
	for rows.Next() {
		var calendar models.Calendar
		var eventCount int
		err := rows.Scan(
			&calendar.ID, &calendar.Name, &calendar.Description, &calendar.Color,
			&calendar.IsDefault, &calendar.IsShared, &calendar.CreatedAt, &calendar.UpdatedAt,
			&eventCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan calendar: %w", err)
		}

		response := calendar.ToResponse()
		response.EventCount = eventCount
		calendars = append(calendars, response)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating calendars: %w", err)
	}

	return calendars, nil
}

// GetCalendar retrieves a single calendar by ID
func (s *CalendarService) GetCalendar(userID, calendarID uuid.UUID) (*models.CalendarResponse, error) {
	query := `
		SELECT c.id, c.name, c.description, c.color, c.is_default, c.is_shared, c.created_at, c.updated_at,
			   COUNT(e.id) as event_count
		FROM calendars c
		LEFT JOIN events e ON c.id = e.calendar_id AND e.user_id = c.user_id
		WHERE c.id = $1 AND c.user_id = $2
		GROUP BY c.id, c.name, c.description, c.color, c.is_default, c.is_shared, c.created_at, c.updated_at
	`

	var calendar models.Calendar
	var eventCount int
	err := s.db.QueryRow(query, calendarID, userID).Scan(
		&calendar.ID, &calendar.Name, &calendar.Description, &calendar.Color,
		&calendar.IsDefault, &calendar.IsShared, &calendar.CreatedAt, &calendar.UpdatedAt,
		&eventCount,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get calendar: %w", err)
	}

	response := calendar.ToResponse()
	response.EventCount = eventCount
	return &response, nil
}

// CreateCalendar creates a new calendar
func (s *CalendarService) CreateCalendar(userID uuid.UUID, req models.CalendarRequest) (*models.CalendarResponse, error) {
	// If this is set as default, unset other defaults first
	if req.IsDefault {
		if err := s.unsetDefaultCalendars(userID); err != nil {
			return nil, fmt.Errorf("failed to unset default calendars: %w", err)
		}
	}

	calendarID := uuid.New()
	color := req.Color
	if color == "" {
		color = "#3B82F6" // Default blue color
	}

	query := `
		INSERT INTO calendars (id, user_id, name, description, color, is_default, is_shared, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	now := time.Now()
	_, err := s.db.Exec(query,
		calendarID, userID, req.Name, req.Description, color, req.IsDefault, false, now, now,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create calendar: %w", err)
	}

	// Return the created calendar
	return s.GetCalendar(userID, calendarID)
}

// UpdateCalendar updates an existing calendar
func (s *CalendarService) UpdateCalendar(userID, calendarID uuid.UUID, req models.CalendarRequest) (*models.CalendarResponse, error) {
	// If this is set as default, unset other defaults first
	if req.IsDefault {
		if err := s.unsetDefaultCalendars(userID); err != nil {
			return nil, fmt.Errorf("failed to unset default calendars: %w", err)
		}
	}

	color := req.Color
	if color == "" {
		color = "#3B82F6" // Default blue color
	}

	query := `
		UPDATE calendars
		SET name = $1, description = $2, color = $3, is_default = $4, updated_at = $5
		WHERE id = $6 AND user_id = $7
	`

	result, err := s.db.Exec(query,
		req.Name, req.Description, color, req.IsDefault, time.Now(),
		calendarID, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update calendar: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil, nil // Calendar not found or doesn't belong to user
	}

	// Return the updated calendar
	return s.GetCalendar(userID, calendarID)
}

// DeleteCalendar deletes a calendar
func (s *CalendarService) DeleteCalendar(userID, calendarID uuid.UUID) error {
	query := `DELETE FROM calendars WHERE id = $1 AND user_id = $2`

	result, err := s.db.Exec(query, calendarID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete calendar: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("calendar not found or doesn't belong to user")
	}

	return nil
}

// unsetDefaultCalendars removes default flag from all calendars for a user
func (s *CalendarService) unsetDefaultCalendars(userID uuid.UUID) error {
	query := `UPDATE calendars SET is_default = false WHERE user_id = $1`
	_, err := s.db.Exec(query, userID)
	return err
}

// GetEvents retrieves events for a user within a date range
func (s *CalendarService) GetEvents(userID uuid.UUID, startDate, endDate time.Time) ([]models.EventResponse, error) {
	query := `
		SELECT e.id, e.calendar_id, e.title, e.description, e.location, e.start_time, e.end_time,
			   e.is_all_day, e.recurrence_rule, e.recurrence_id, e.reminder_minutes, e.is_completed,
			   e.created_at, e.updated_at, c.name as calendar_name, c.color as calendar_color
		FROM events e
		JOIN calendars c ON e.calendar_id = c.id
		WHERE e.user_id = $1 AND e.start_time <= $2 AND e.end_time >= $3
		ORDER BY e.start_time ASC
	`

	rows, err := s.db.Query(query, userID, endDate, startDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query events: %w", err)
	}
	defer rows.Close()

	var events []models.EventResponse
	for rows.Next() {
		var event models.Event
		var calendarName, calendarColor string
		err := rows.Scan(
			&event.ID, &event.CalendarID, &event.Title, &event.Description, &event.Location,
			&event.StartTime, &event.EndTime, &event.IsAllDay, &event.RecurrenceRule,
			&event.RecurrenceID, &event.ReminderMinutes, &event.IsCompleted,
			&event.CreatedAt, &event.UpdatedAt, &calendarName, &calendarColor,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}

		response := event.ToResponse(calendarName, calendarColor)

		// Get attendees for this event
		attendees, err := s.getEventAttendees(event.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get attendees for event %s: %w", event.ID, err)
		}
		response.Attendees = attendees

		events = append(events, response)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating events: %w", err)
	}

	return events, nil
}

// GetEvent retrieves a single event by ID
func (s *CalendarService) GetEvent(userID, eventID uuid.UUID) (*models.EventResponse, error) {
	query := `
		SELECT e.id, e.calendar_id, e.title, e.description, e.location, e.start_time, e.end_time,
			   e.is_all_day, e.recurrence_rule, e.recurrence_id, e.reminder_minutes, e.is_completed,
			   e.created_at, e.updated_at, c.name as calendar_name, c.color as calendar_color
		FROM events e
		JOIN calendars c ON e.calendar_id = c.id
		WHERE e.id = $1 AND e.user_id = $2
	`

	var event models.Event
	var calendarName, calendarColor string
	err := s.db.QueryRow(query, eventID, userID).Scan(
		&event.ID, &event.CalendarID, &event.Title, &event.Description, &event.Location,
		&event.StartTime, &event.EndTime, &event.IsAllDay, &event.RecurrenceRule,
		&event.RecurrenceID, &event.ReminderMinutes, &event.IsCompleted,
		&event.CreatedAt, &event.UpdatedAt, &calendarName, &calendarColor,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	response := event.ToResponse(calendarName, calendarColor)

	// Get attendees for this event
	attendees, err := s.getEventAttendees(event.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get attendees for event: %w", err)
	}
	response.Attendees = attendees

	return &response, nil
}

// CreateEvent creates a new event
func (s *CalendarService) CreateEvent(userID uuid.UUID, req models.EventRequest) (*models.EventResponse, error) {
	// Verify calendar belongs to user
	var calendarExists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM calendars WHERE id = $1 AND user_id = $2)",
		req.CalendarID, userID).Scan(&calendarExists)
	if err != nil {
		return nil, fmt.Errorf("failed to verify calendar: %w", err)
	}
	if !calendarExists {
		return nil, fmt.Errorf("calendar not found or doesn't belong to user")
	}

	// Validate event
	event := models.Event{
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}
	if err := event.Validate(); err != nil {
		return nil, fmt.Errorf("invalid event: %w", err)
	}

	eventID := uuid.New()

	query := `
		INSERT INTO events (id, calendar_id, user_id, title, description, location, start_time, end_time,
						   is_all_day, recurrence_rule, reminder_minutes, is_completed, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`

	now := time.Now()
	_, err = s.db.Exec(query,
		eventID, req.CalendarID, userID, req.Title, req.Description, req.Location,
		req.StartTime, req.EndTime, req.IsAllDay, req.RecurrenceRule, req.ReminderMinutes,
		false, now, now,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	// Add attendees if provided
	if len(req.AttendeeIDs) > 0 {
		if err := s.addEventAttendees(eventID, req.AttendeeIDs); err != nil {
			return nil, fmt.Errorf("failed to add attendees: %w", err)
		}
	}

	// Return the created event
	return s.GetEvent(userID, eventID)
}

// UpdateEvent updates an existing event
func (s *CalendarService) UpdateEvent(userID, eventID uuid.UUID, req models.EventRequest) (*models.EventResponse, error) {
	// Verify calendar belongs to user
	var calendarExists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM calendars WHERE id = $1 AND user_id = $2)",
		req.CalendarID, userID).Scan(&calendarExists)
	if err != nil {
		return nil, fmt.Errorf("failed to verify calendar: %w", err)
	}
	if !calendarExists {
		return nil, fmt.Errorf("calendar not found or doesn't belong to user")
	}

	// Validate event
	event := models.Event{
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}
	if err := event.Validate(); err != nil {
		return nil, fmt.Errorf("invalid event: %w", err)
	}

	query := `
		UPDATE events
		SET calendar_id = $1, title = $2, description = $3, location = $4, start_time = $5,
		    end_time = $6, is_all_day = $7, recurrence_rule = $8, reminder_minutes = $9, updated_at = $10
		WHERE id = $11 AND user_id = $12
	`

	result, err := s.db.Exec(query,
		req.CalendarID, req.Title, req.Description, req.Location, req.StartTime, req.EndTime,
		req.IsAllDay, req.RecurrenceRule, req.ReminderMinutes, time.Now(),
		eventID, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update event: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil, nil // Event not found or doesn't belong to user
	}

	// Update attendees
	if err := s.updateEventAttendees(eventID, req.AttendeeIDs); err != nil {
		return nil, fmt.Errorf("failed to update attendees: %w", err)
	}

	// Return the updated event
	return s.GetEvent(userID, eventID)
}

// DeleteEvent deletes an event
func (s *CalendarService) DeleteEvent(userID, eventID uuid.UUID) error {
	query := `DELETE FROM events WHERE id = $1 AND user_id = $2`

	result, err := s.db.Exec(query, eventID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("event not found or doesn't belong to user")
	}

	return nil
}

// getEventAttendees retrieves attendees for an event
func (s *CalendarService) getEventAttendees(eventID uuid.UUID) ([]models.EventAttendeeResponse, error) {
	query := `
		SELECT ea.id, ea.event_id, ea.contact_id, ea.status, ea.created_at,
			   COALESCE(c.first_name || ' ', '') || COALESCE(c.last_name, '') as name, c.email
		FROM event_attendees ea
		JOIN contacts c ON ea.contact_id = c.id
		WHERE ea.event_id = $1
		ORDER BY name ASC
	`

	rows, err := s.db.Query(query, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to query attendees: %w", err)
	}
	defer rows.Close()

	var attendees []models.EventAttendeeResponse
	for rows.Next() {
		var attendee models.EventAttendee
		var name string
		var email *string
		err := rows.Scan(
			&attendee.ID, &attendee.EventID, &attendee.ContactID, &attendee.Status, &attendee.CreatedAt,
			&name, &email,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan attendee: %w", err)
		}

		response := attendee.ToResponse(name, email)
		attendees = append(attendees, response)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating attendees: %w", err)
	}

	return attendees, nil
}

// addEventAttendees adds attendees to an event
func (s *CalendarService) addEventAttendees(eventID uuid.UUID, contactIDs []uuid.UUID) error {
	for _, contactID := range contactIDs {
		attendeeID := uuid.New()
		query := `
			INSERT INTO event_attendees (id, event_id, contact_id, status, created_at)
			VALUES ($1, $2, $3, $4, $5)
		`
		_, err := s.db.Exec(query, attendeeID, eventID, contactID, "pending", time.Now())
		if err != nil {
			return fmt.Errorf("failed to add attendee %s: %w", contactID, err)
		}
	}
	return nil
}

// updateEventAttendees replaces attendees for an event
func (s *CalendarService) updateEventAttendees(eventID uuid.UUID, contactIDs []uuid.UUID) error {
	// Remove existing attendees
	_, err := s.db.Exec("DELETE FROM event_attendees WHERE event_id = $1", eventID)
	if err != nil {
		return fmt.Errorf("failed to remove existing attendees: %w", err)
	}

	// Add new attendees
	if len(contactIDs) > 0 {
		return s.addEventAttendees(eventID, contactIDs)
	}

	return nil
}

// CreateDefaultCalendar creates a default calendar for a new user
func (s *CalendarService) CreateDefaultCalendar(userID uuid.UUID) error {
	req := models.CalendarRequest{
		Name:       "My Calendar",
		Color:      "#3B82F6",
		IsDefault:  true,
	}

	_, err := s.CreateCalendar(userID, req)
	return err
}
