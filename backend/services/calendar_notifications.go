package services

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"tpt-titan/backend/models"
)

// CalendarNotificationService handles calendar reminders and notifications
type CalendarNotificationService struct {
	db              *sql.DB
	emailService    *EmailService
	websocketSvc    interface{} // WebSocket service interface (placeholder)
	smsService      interface{} // SMS service interface (placeholder)
}

// NotificationType represents different types of notifications
type NotificationType string

const (
	NotificationTypeEmail     NotificationType = "email"
	NotificationTypeInApp     NotificationType = "in_app"
	NotificationTypeSMS       NotificationType = "sms"
	NotificationTypePush      NotificationType = "push"
)

// Reminder represents a calendar reminder configuration
type Reminder struct {
	ID          uuid.UUID       `json:"id"`
	EventID     uuid.UUID       `json:"event_id"`
	UserID      uuid.UUID       `json:"user_id"`
	Type        NotificationType `json:"type"`
	MinutesBefore int           `json:"minutes_before"` // Minutes before event
	Message     string          `json:"message,omitempty"`
	SentAt      *time.Time      `json:"sent_at,omitempty"`
	Status      string          `json:"status"` // "pending", "sent", "failed"
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// EventNotification represents a notification sent for an event
type EventNotification struct {
	ID        uuid.UUID       `json:"id"`
	EventID   uuid.UUID       `json:"event_id"`
	UserID    uuid.UUID       `json:"user_id"`
	Type      NotificationType `json:"type"`
	Title     string          `json:"title"`
	Message   string          `json:"message"`
	SentAt    time.Time       `json:"sent_at"`
	Status    string          `json:"status"`
}

// NotificationSettings represents user notification preferences
type NotificationSettings struct {
	UserID                    uuid.UUID `json:"user_id"`
	EmailReminders           bool      `json:"email_reminders"`
	InAppReminders           bool      `json:"in_app_reminders"`
	SMSReminders             bool      `json:"sms_reminders"`
	PushReminders            bool      `json:"push_reminders"`
	DefaultReminderMinutes   int       `json:"default_reminder_minutes"`
	EmailForAllDayEvents     bool      `json:"email_for_all_day_events"`
	EmailForTentativeEvents  bool      `json:"email_for_tentative_events"`
	EmailForCancelledEvents  bool      `json:"email_for_cancelled_events"`
	QuietHoursStart          *string   `json:"quiet_hours_start,omitempty"` // HH:MM format
	QuietHoursEnd            *string   `json:"quiet_hours_end,omitempty"`
	Timezone                 string    `json:"timezone"`
	UpdatedAt                time.Time `json:"updated_at"`
}

// NewCalendarNotificationService creates a new calendar notification service
func NewCalendarNotificationService(db *sql.DB, emailSvc *EmailService, wsSvc interface{}) *CalendarNotificationService {
	return &CalendarNotificationService{
		db:           db,
		emailService: emailSvc,
		websocketSvc: wsSvc,
	}
}

// CreateReminder creates a reminder for an event
func (cns *CalendarNotificationService) CreateReminder(eventID, userID uuid.UUID, notificationType NotificationType, minutesBefore int, customMessage string) (*Reminder, error) {
	reminder := &Reminder{
		ID:            uuid.New(),
		EventID:       eventID,
		UserID:        userID,
		Type:          notificationType,
		MinutesBefore: minutesBefore,
		Message:       customMessage,
		Status:        "pending",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	query := `
		INSERT INTO calendar_reminders (id, event_id, user_id, type, minutes_before, message, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := cns.db.Exec(query,
		reminder.ID, reminder.EventID, reminder.UserID, reminder.Type,
		reminder.MinutesBefore, reminder.Message, reminder.Status,
		reminder.CreatedAt, reminder.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return reminder, nil
}

// GetEventReminders gets all reminders for an event
func (cns *CalendarNotificationService) GetEventReminders(eventID uuid.UUID) ([]Reminder, error) {
	query := `
		SELECT id, event_id, user_id, type, minutes_before, message, sent_at, status, created_at, updated_at
		FROM calendar_reminders
		WHERE event_id = $1
		ORDER BY minutes_before DESC
	`

	rows, err := cns.db.Query(query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reminders []Reminder
	for rows.Next() {
		var reminder Reminder
		err := rows.Scan(
			&reminder.ID, &reminder.EventID, &reminder.UserID, &reminder.Type,
			&reminder.MinutesBefore, &reminder.Message, &reminder.SentAt,
			&reminder.Status, &reminder.CreatedAt, &reminder.UpdatedAt,
		)
		if err != nil {
			continue
		}
		reminders = append(reminders, reminder)
	}

	return reminders, nil
}

// ProcessPendingReminders processes all pending reminders that should be sent
func (cns *CalendarNotificationService) ProcessPendingReminders() error {
	// Find reminders that should be sent now
	now := time.Now()
	query := `
		SELECT r.id, r.event_id, r.user_id, r.type, r.minutes_before, r.message,
		       e.title, e.start_time, e.end_time, e.location, u.email, u.username
		FROM calendar_reminders r
		JOIN calendar_events e ON r.event_id = e.id
		JOIN users u ON r.user_id = u.id
		WHERE r.status = 'pending'
		AND e.start_time - INTERVAL '1 minute' * r.minutes_before <= $1
		AND e.start_time > $1
	`

	rows, err := cns.db.Query(query, now)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var reminderID, eventID, userID uuid.UUID
		var reminderType NotificationType
		var minutesBefore int
		var customMessage, eventTitle, location, userEmail, username sql.NullString
		var startTime, endTime time.Time

		err := rows.Scan(
			&reminderID, &eventID, &userID, &reminderType, &minutesBefore,
			&customMessage, &eventTitle, &startTime, &endTime, &location,
			&userEmail, &username,
		)
		if err != nil {
			continue
		}

		// Send the reminder
		err = cns.sendReminder(reminderID, userID, reminderType, ReminderContext{
			EventID:    eventID,
			EventTitle: eventTitle.String,
			StartTime:  startTime,
			EndTime:    endTime,
			Location:   location.String,
			UserEmail:  userEmail.String,
			Username:   username.String,
			MinutesBefore: minutesBefore,
			CustomMessage: customMessage.String,
		})

		if err != nil {
			// Mark as failed
			cns.updateReminderStatus(reminderID, "failed")
		} else {
			// Mark as sent
			cns.updateReminderStatus(reminderID, "sent")
		}
	}

	return nil
}

// ReminderContext contains information about the event for the reminder
type ReminderContext struct {
	EventID       uuid.UUID
	EventTitle    string
	StartTime     time.Time
	EndTime       time.Time
	Location      string
	UserEmail     string
	Username      string
	MinutesBefore int
	CustomMessage string
}

// sendReminder sends a reminder notification
func (cns *CalendarNotificationService) sendReminder(reminderID, userID uuid.UUID, notificationType NotificationType, ctx ReminderContext) error {
	// Generate reminder message
	message := cns.generateReminderMessage(ctx)

	switch notificationType {
	case NotificationTypeEmail:
		return cns.sendEmailReminder(ctx.UserEmail, ctx.EventTitle, message, ctx)
	case NotificationTypeInApp:
		return cns.sendInAppReminder(userID, ctx.EventTitle, message, ctx)
	case NotificationTypeSMS:
		return cns.sendSMSReminder(ctx.UserEmail, message) // Using email as phone number placeholder
	case NotificationTypePush:
		return cns.sendPushReminder(userID, ctx.EventTitle, message, ctx)
	default:
		return fmt.Errorf("unsupported notification type: %s", notificationType)
	}
}

// generateReminderMessage generates a reminder message
func (cns *CalendarNotificationService) generateReminderMessage(ctx ReminderContext) string {
	if ctx.CustomMessage != "" {
		return ctx.CustomMessage
	}

	message := fmt.Sprintf("Reminder: %s is starting in %d minutes", ctx.EventTitle, ctx.MinutesBefore)

	if !ctx.StartTime.IsZero() {
		message += fmt.Sprintf(" at %s", ctx.StartTime.Format("3:04 PM"))
	}

	if ctx.Location != "" {
		message += fmt.Sprintf(" (Location: %s)", ctx.Location)
	}

	return message
}

// Send notification methods
func (cns *CalendarNotificationService) sendEmailReminder(email, subject, message string, ctx ReminderContext) error {
	htmlMessage := fmt.Sprintf(`
		<h2>%s</h2>
		<p><strong>Time:</strong> %s - %s</p>
		<p><strong>Location:</strong> %s</p>
		<p>%s</p>
		<p>This is a reminder sent %d minutes before the event.</p>
	`, ctx.EventTitle,
		ctx.StartTime.Format("Mon, Jan 2, 2006 at 3:04 PM"),
		ctx.EndTime.Format("3:04 PM"),
		ctx.Location,
		message,
		ctx.MinutesBefore,
	)

	_, err := cns.emailService.SendEmail(uuid.UUID{}, models.EmailRequest{
		RecipientEmails: []string{email},
		Subject:         "Event Reminder: " + subject,
		Content:         htmlMessage,
	})
	return err
}

func (cns *CalendarNotificationService) sendInAppReminder(userID uuid.UUID, title, message string, ctx ReminderContext) error {
	// In a real implementation, broadcast this via WebSocket
	// notification payload would be:
	// { type: "calendar_reminder", user_id, title, message, event_id, timestamp }
	_ = userID // used in WebSocket routing
	_ = title
	_ = message
	_ = ctx
	return nil
}

func (cns *CalendarNotificationService) sendSMSReminder(phoneNumber, message string) error {
	// Placeholder for SMS service integration
	// In a real implementation, integrate with Twilio, AWS SNS, etc.
	return nil
}

func (cns *CalendarNotificationService) sendPushReminder(userID uuid.UUID, title, message string, ctx ReminderContext) error {
	// Placeholder for push notification service
	// In a real implementation, integrate with FCM, APNS, etc.
	return nil
}

// updateReminderStatus updates the status of a reminder
func (cns *CalendarNotificationService) updateReminderStatus(reminderID uuid.UUID, status string) error {
	query := `
		UPDATE calendar_reminders
		SET status = $1, sent_at = $2, updated_at = $3
		WHERE id = $4
	`

	sentAt := time.Now()
	if status != "sent" {
		sentAt = time.Time{}
	}

	_, err := cns.db.Exec(query, status, sentAt, time.Now(), reminderID)
	return err
}

// GetNotificationSettings gets user notification settings
func (cns *CalendarNotificationService) GetNotificationSettings(userID uuid.UUID) (*NotificationSettings, error) {
	var settings NotificationSettings
	query := `
		SELECT user_id, email_reminders, in_app_reminders, sms_reminders, push_reminders,
		       default_reminder_minutes, email_for_all_day_events, email_for_tentative_events,
		       email_for_cancelled_events, quiet_hours_start, quiet_hours_end, timezone, updated_at
		FROM calendar_notification_settings
		WHERE user_id = $1
	`

	err := cns.db.QueryRow(query, userID).Scan(
		&settings.UserID, &settings.EmailReminders, &settings.InAppReminders,
		&settings.SMSReminders, &settings.PushReminders, &settings.DefaultReminderMinutes,
		&settings.EmailForAllDayEvents, &settings.EmailForTentativeEvents,
		&settings.EmailForCancelledEvents, &settings.QuietHoursStart,
		&settings.QuietHoursEnd, &settings.Timezone, &settings.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		// Return default settings
		settings = NotificationSettings{
			UserID:                  userID,
			EmailReminders:          true,
			InAppReminders:          true,
			SMSReminders:            false,
			PushReminders:           false,
			DefaultReminderMinutes:  15,
			EmailForAllDayEvents:    false,
			EmailForTentativeEvents: true,
			EmailForCancelledEvents: true,
			Timezone:                "UTC",
			UpdatedAt:               time.Now(),
		}
		return &settings, nil
	}

	return &settings, err
}

// UpdateNotificationSettings updates user notification settings
func (cns *CalendarNotificationService) UpdateNotificationSettings(settings *NotificationSettings) error {
	settings.UpdatedAt = time.Now()

	query := `
		INSERT INTO calendar_notification_settings (
			user_id, email_reminders, in_app_reminders, sms_reminders, push_reminders,
			default_reminder_minutes, email_for_all_day_events, email_for_tentative_events,
			email_for_cancelled_events, quiet_hours_start, quiet_hours_end, timezone, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		ON CONFLICT (user_id) DO UPDATE SET
			email_reminders = EXCLUDED.email_reminders,
			in_app_reminders = EXCLUDED.in_app_reminders,
			sms_reminders = EXCLUDED.sms_reminders,
			push_reminders = EXCLUDED.push_reminders,
			default_reminder_minutes = EXCLUDED.default_reminder_minutes,
			email_for_all_day_events = EXCLUDED.email_for_all_day_events,
			email_for_tentative_events = EXCLUDED.email_for_tentative_events,
			email_for_cancelled_events = EXCLUDED.email_for_cancelled_events,
			quiet_hours_start = EXCLUDED.quiet_hours_start,
			quiet_hours_end = EXCLUDED.quiet_hours_end,
			timezone = EXCLUDED.timezone,
			updated_at = EXCLUDED.updated_at
	`

	_, err := cns.db.Exec(query,
		settings.UserID, settings.EmailReminders, settings.InAppReminders,
		settings.SMSReminders, settings.PushReminders, settings.DefaultReminderMinutes,
		settings.EmailForAllDayEvents, settings.EmailForTentativeEvents,
		settings.EmailForCancelledEvents, settings.QuietHoursStart,
		settings.QuietHoursEnd, settings.Timezone, settings.UpdatedAt,
	)

	return err
}

// SendEventNotification sends a notification for an event change
func (cns *CalendarNotificationService) SendEventNotification(eventID uuid.UUID, notificationType string, title, message string, recipientUsers []uuid.UUID) error {
	for _, userID := range recipientUsers {
		notification := &EventNotification{
			ID:      uuid.New(),
			EventID: eventID,
			UserID:  userID,
			Type:    NotificationType(notificationType),
			Title:   title,
			Message: message,
			SentAt:  time.Now(),
			Status:  "sent",
		}

		// Save notification record
		query := `
			INSERT INTO calendar_event_notifications (id, event_id, user_id, type, title, message, sent_at, status)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`

		_, err := cns.db.Exec(query,
			notification.ID, notification.EventID, notification.UserID, notification.Type,
			notification.Title, notification.Message, notification.SentAt, notification.Status,
		)

		if err != nil {
			continue // Continue with other recipients
		}

		// Send the notification
		switch notification.Type {
		case NotificationTypeEmail:
			// Get user email
			var email string
			cns.db.QueryRow("SELECT email FROM users WHERE id = $1", userID).Scan(&email)
			_, _ = cns.emailService.SendEmail(uuid.UUID{}, models.EmailRequest{
				RecipientEmails: []string{email},
				Subject:         title,
				Content:         message,
			})
		case NotificationTypeInApp:
			// Send WebSocket notification
			cns.sendInAppReminder(userID, title, message, ReminderContext{})
		}
	}

	return nil
}

// IsInQuietHours checks if current time is within user's quiet hours
func (cns *CalendarNotificationService) IsInQuietHours(userID uuid.UUID) (bool, error) {
	settings, err := cns.GetNotificationSettings(userID)
	if err != nil {
		return false, err
	}

	if settings.QuietHoursStart == nil || settings.QuietHoursEnd == nil {
		return false, nil
	}

	now := time.Now()
	startTime, err := time.Parse("15:04", *settings.QuietHoursStart)
	if err != nil {
		return false, err
	}

	endTime, err := time.Parse("15:04", *settings.QuietHoursEnd)
	if err != nil {
		return false, err
	}

	// Convert to today's date
	startToday := time.Date(now.Year(), now.Month(), now.Day(), startTime.Hour(), startTime.Minute(), 0, 0, now.Location())
	endToday := time.Date(now.Year(), now.Month(), now.Day(), endTime.Hour(), endTime.Minute(), 0, 0, now.Location())

	// Handle case where end time is next day
	if endTime.Before(startTime) {
		endToday = endToday.AddDate(0, 0, 1)
	}

	return now.After(startToday) && now.Before(endToday), nil
}

// GetUpcomingReminders gets reminders that will be sent in the next specified minutes
func (cns *CalendarNotificationService) GetUpcomingReminders(minutes int) ([]map[string]interface{}, error) {
	query := `
		SELECT r.id, r.event_id, r.user_id, r.type, r.minutes_before,
		       e.title, e.start_time, u.username
		FROM calendar_reminders r
		JOIN calendar_events e ON r.event_id = e.id
		JOIN users u ON r.user_id = u.id
		WHERE r.status = 'pending'
		AND e.start_time - INTERVAL '1 minute' * r.minutes_before <= $1
		AND e.start_time > $1
		ORDER BY e.start_time
	`

	rows, err := cns.db.Query(query, time.Now().Add(time.Duration(minutes)*time.Minute))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reminders []map[string]interface{}
	for rows.Next() {
		var id, eventID, userID uuid.UUID
		var reminderType NotificationType
		var minutesBefore int
		var eventTitle, username string
		var startTime time.Time

		err := rows.Scan(&id, &eventID, &userID, &reminderType, &minutesBefore, &eventTitle, &startTime, &username)
		if err != nil {
			continue
		}

		reminders = append(reminders, map[string]interface{}{
			"id":             id,
			"event_id":       eventID,
			"user_id":        userID,
			"type":           reminderType,
			"minutes_before": minutesBefore,
			"event_title":    eventTitle,
			"start_time":     startTime,
			"username":       username,
		})
	}

	return reminders, nil
}

// DeleteReminder deletes a reminder
func (cns *CalendarNotificationService) DeleteReminder(reminderID uuid.UUID) error {
	query := `DELETE FROM calendar_reminders WHERE id = $1`
	_, err := cns.db.Exec(query, reminderID)
	return err
}

// GetReminderStats returns reminder statistics
func (cns *CalendarNotificationService) GetReminderStats(userID uuid.UUID) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Count reminders by status
	statusQuery := `
		SELECT status, COUNT(*) as count
		FROM calendar_reminders
		WHERE user_id = $1
		GROUP BY status
	`

	rows, err := cns.db.Query(statusQuery, userID)
	if err == nil {
		for rows.Next() {
			var status string
			var count int
			rows.Scan(&status, &count)
			stats[fmt.Sprintf("reminders_%s", status)] = count
		}
		rows.Close()
	}

	// Count reminders by type
	typeQuery := `
		SELECT type, COUNT(*) as count
		FROM calendar_reminders
		WHERE user_id = $1
		GROUP BY type
	`

	rows, err = cns.db.Query(typeQuery, userID)
	if err == nil {
		for rows.Next() {
			var reminderType NotificationType
			var count int
			rows.Scan(&reminderType, &count)
			stats[fmt.Sprintf("reminders_type_%s", reminderType)] = count
		}
		rows.Close()
	}

	// Total reminders
	var total int
	cns.db.QueryRow("SELECT COUNT(*) FROM calendar_reminders WHERE user_id = $1", userID).Scan(&total)
	stats["total_reminders"] = total

	return stats, nil
}
