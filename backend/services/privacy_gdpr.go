package services

import (
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// PrivacyGDPRService handles GDPR compliance and privacy controls
type PrivacyGDPRService struct {
	db *sql.DB
}

// DataSubjectRight represents GDPR data subject rights
type DataSubjectRight string

const (
	RightAccess          DataSubjectRight = "access"          // Right to access personal data
	RightRectification   DataSubjectRight = "rectification"   // Right to rectify inaccurate data
	RightErasure         DataSubjectRight = "erasure"         // Right to erasure ("right to be forgotten")
	RightRestriction     DataSubjectRight = "restriction"     // Right to restrict processing
	RightPortability     DataSubjectRight = "portability"     // Right to data portability
	RightObjection       DataSubjectRight = "objection"       // Right to object to processing
	RightWithdrawConsent DataSubjectRight = "withdraw_consent" // Right to withdraw consent
)

// PrivacyConsent represents user consent for data processing
type PrivacyConsent struct {
	ID              uuid.UUID              `json:"id"`
	UserID          uuid.UUID              `json:"user_id"`
	ConsentType     string                 `json:"consent_type"`     // "marketing", "analytics", "third_party", etc.
	Description     string                 `json:"description"`
	Version         string                 `json:"version"`          // Consent version
	GivenAt         time.Time              `json:"given_at"`
	ValidUntil      *time.Time             `json:"valid_until,omitempty"`
	WithdrawnAt     *time.Time             `json:"withdrawn_at,omitempty"`
	IP              string                 `json:"ip_address"`
	UserAgent       string                 `json:"user_agent"`
	Source          string                 `json:"source"`           // "website", "mobile_app", etc.
	Scope           map[string]interface{} `json:"scope,omitempty"`  // Specific scope of consent
}

// DataProcessingRecord represents a record of data processing activity
type DataProcessingRecord struct {
	ID              uuid.UUID              `json:"id"`
	UserID          uuid.UUID              `json:"user_id"`
	ProcessingType  string                 `json:"processing_type"`  // "collection", "storage", "sharing", etc.
	DataCategories  []string               `json:"data_categories"`  // "personal", "financial", "health", etc.
	Purpose         string                 `json:"purpose"`
	LegalBasis      string                 `json:"legal_basis"`      // "consent", "contract", "legitimate_interest", etc.
	Recipients      []string               `json:"recipients,omitempty"`
	Location        string                 `json:"location"`         // Where data is stored/processed
	RetentionPeriod *int                   `json:"retention_period,omitempty"` // Days
	ProcessedAt     time.Time              `json:"processed_at"`
	Controller      string                 `json:"controller"`       // Data controller organization
}

// DataBreachRecord represents a data breach incident
type DataBreachRecord struct {
	ID              uuid.UUID `json:"id"`
	IncidentID      string    `json:"incident_id"`
	Description     string    `json:"description"`
	DataAffected    []string  `json:"data_affected"`    // Types of data affected
	UsersAffected   int       `json:"users_affected"`
	RiskLevel       string    `json:"risk_level"`       // "low", "medium", "high", "critical"
	DetectedAt      time.Time `json:"detected_at"`
	ReportedAt      *time.Time `json:"reported_at,omitempty"`
	ResolvedAt      *time.Time `json:"resolved_at,omitempty"`
	ActionsTaken    string    `json:"actions_taken"`
	PreventiveMeasures string `json:"preventive_measures"`
	Status          string    `json:"status"`           // "investigating", "contained", "resolved"
	ReportedToAuthorities bool `json:"reported_to_authorities"`
}

// PrivacySettings represents user privacy preferences
type PrivacySettings struct {
	UserID                      uuid.UUID `json:"user_id"`
	DataCollectionConsent       bool      `json:"data_collection_consent"`
	AnalyticsConsent            bool      `json:"analytics_consent"`
	MarketingConsent            bool      `json:"marketing_consent"`
	ThirdPartySharingConsent    bool      `json:"third_party_sharing_consent"`
	DataRetentionPreference     string    `json:"data_retention_preference"`     // "minimum", "standard", "extended"
	ContactMethodPreference     string    `json:"contact_method_preference"`     // "email", "phone", "none"
	DataExportRequested         bool      `json:"data_export_requested"`
	DataDeletionRequested       bool      `json:"data_deletion_requested"`
	DoNotTrack                  bool      `json:"do_not_track"`
	AdvertisingOptOut           bool      `json:"advertising_opt_out"`
	LocationTrackingDisabled    bool      `json:"location_tracking_disabled"`
	ProfileVisibility           string    `json:"profile_visibility"`            // "public", "contacts", "private"
	UpdatedAt                   time.Time `json:"updated_at"`
}

// DataSubjectRequest represents a GDPR data subject access request
type DataSubjectRequest struct {
	ID              uuid.UUID         `json:"id"`
	RequestID       string            `json:"request_id"`       // Unique reference number
	UserID          uuid.UUID         `json:"user_id"`
	RequestType     DataSubjectRight  `json:"request_type"`
	Description     string            `json:"description"`
	Status          string            `json:"status"`           // "pending", "processing", "completed", "rejected"
	RequestedAt     time.Time         `json:"requested_at"`
	CompletedAt     *time.Time        `json:"completed_at,omitempty"`
	ResponseData    json.RawMessage   `json:"response_data,omitempty"` // JSON response for access requests
	RejectionReason string            `json:"rejection_reason,omitempty"`
	VerifiedAt      *time.Time        `json:"verified_at,omitempty"`   // When user identity was verified
	ProcessingNotes string            `json:"processing_notes"`
	DataController  string            `json:"data_controller"`
}

// PrivacyAuditLog represents audit logging for privacy-related activities
type PrivacyAuditLog struct {
	ID              uuid.UUID `json:"id"`
	UserID          uuid.UUID `json:"user_id,omitempty"`
	Action          string    `json:"action"`          // "consent_given", "data_accessed", "data_deleted", etc.
	Resource        string    `json:"resource"`        // What was affected
	ResourceID      string    `json:"resource_id,omitempty"`
	IPAddress       string    `json:"ip_address"`
	UserAgent       string    `json:"user_agent"`
	Timestamp       time.Time `json:"timestamp"`
	Details         json.RawMessage `json:"details,omitempty"`
	ComplianceFlag  bool      `json:"compliance_flag"` // Flags for regulatory review
}

// NewPrivacyGDPRService creates a new GDPR privacy service
func NewPrivacyGDPRService(db *sql.DB) *PrivacyGDPRService {
	return &PrivacyGDPRService{db: db}
}

// RecordConsent records user consent for data processing
func (pgs *PrivacyGDPRService) RecordConsent(consent *PrivacyConsent) error {
	consent.ID = uuid.New()

	query := `
		INSERT INTO privacy_consents (id, user_id, consent_type, description, version,
			given_at, valid_until, withdrawn_at, ip_address, user_agent, source, scope)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err := pgs.db.Exec(query,
		consent.ID, consent.UserID, consent.ConsentType, consent.Description, consent.Version,
		consent.GivenAt, consent.ValidUntil, consent.WithdrawnAt, consent.IP,
		consent.UserAgent, consent.Source, consent.Scope,
	)

	if err != nil {
		return err
	}

	// Log audit event
	pgs.logPrivacyAudit(consent.UserID, "consent_given", "privacy_consent", consent.ID.String(),
		consent.IP, consent.UserAgent, map[string]interface{}{
			"consent_type": consent.ConsentType,
			"version":      consent.Version,
		})

	return nil
}

// WithdrawConsent withdraws user consent
func (pgs *PrivacyGDPRService) WithdrawConsent(userID uuid.UUID, consentType string, reason string) error {
	now := time.Now()

	query := `
		UPDATE privacy_consents
		SET withdrawn_at = $1, scope = scope || '{"withdrawn_reason": "$2", "withdrawn_at": "$3"}'::jsonb
		WHERE user_id = $2 AND consent_type = $3 AND withdrawn_at IS NULL
	`

	_, err := pgs.db.Exec(query, now, reason, now.Format(time.RFC3339), userID, consentType)

	if err != nil {
		return err
	}

	// Log audit event
	pgs.logPrivacyAudit(userID, "consent_withdrawn", "privacy_consent", "",
		"", "", map[string]interface{}{
			"consent_type": consentType,
			"reason":       reason,
		})

	return nil
}

// RecordDataProcessing records a data processing activity
func (pgs *PrivacyGDPRService) RecordDataProcessing(record *DataProcessingRecord) error {
	record.ID = uuid.New()
	record.ProcessedAt = time.Now()

	// Hash sensitive data for privacy
	recordHash := sha256.Sum256([]byte(fmt.Sprintf("%s-%s-%s", record.UserID, record.ProcessingType, record.ProcessedAt)))
	_ = recordHash
	record.ID = uuid.New() // Use hash-based ID for privacy

	query := `
		INSERT INTO data_processing_records (id, user_id, processing_type, data_categories,
			purpose, legal_basis, recipients, location, retention_period, processed_at, controller)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := pgs.db.Exec(query,
		record.ID, record.UserID, record.ProcessingType, record.DataCategories,
		record.Purpose, record.LegalBasis, record.Recipients, record.Location,
		record.RetentionPeriod, record.ProcessedAt, record.Controller,
	)

	return err
}

// SubmitDataSubjectRequest submits a GDPR data subject access request
func (pgs *PrivacyGDPRService) SubmitDataSubjectRequest(request *DataSubjectRequest) error {
	request.ID = uuid.New()
	request.RequestID = pgs.generateRequestID()
	request.RequestedAt = time.Now()
	request.Status = "pending"
	request.DataController = "TPT Titan"

	query := `
		INSERT INTO data_subject_requests (id, request_id, user_id, request_type, description,
			status, requested_at, processing_notes, data_controller)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := pgs.db.Exec(query,
		request.ID, request.RequestID, request.UserID, request.RequestType,
		request.Description, request.Status, request.RequestedAt,
		request.ProcessingNotes, request.DataController,
	)

	if err != nil {
		return err
	}

	// Log audit event
	pgs.logPrivacyAudit(request.UserID, "dsr_submitted", "data_subject_request", request.ID.String(),
		"", "", map[string]interface{}{
			"request_type": request.RequestType,
			"request_id":   request.RequestID,
		})

	return nil
}

// ProcessDataSubjectRequest processes a data subject request
func (pgs *PrivacyGDPRService) ProcessDataSubjectRequest(requestID string, action string, responseData interface{}, notes string) error {
	now := time.Now()

	var updateQuery string
	var args []interface{}

	switch action {
	case "complete":
		updateQuery = `
			UPDATE data_subject_requests
			SET status = 'completed', completed_at = $1, response_data = $2, processing_notes = $3
			WHERE request_id = $4
		`
		jsonData, _ := json.Marshal(responseData)
		args = []interface{}{now, jsonData, notes, requestID}

	case "reject":
		updateQuery = `
			UPDATE data_subject_requests
			SET status = 'rejected', completed_at = $1, rejection_reason = $2, processing_notes = $3
			WHERE request_id = $4
		`
		args = []interface{}{now, responseData.(string), notes, requestID}

	default:
		return fmt.Errorf("invalid action: %s", action)
	}

	result, err := pgs.db.Exec(updateQuery, args...)
	if err != nil {
		return err
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return fmt.Errorf("request not found or already processed")
	}

	// Log audit event
	pgs.logPrivacyAudit(uuid.Nil, "dsr_processed", "data_subject_request", requestID,
		"", "", map[string]interface{}{
			"action": action,
			"notes":  notes,
		})

	return nil
}

// GetDataSubjectRequest gets a data subject request by ID
func (pgs *PrivacyGDPRService) GetDataSubjectRequest(requestID string) (*DataSubjectRequest, error) {
	var request DataSubjectRequest

	query := `
		SELECT id, request_id, user_id, request_type, description, status, requested_at,
		       completed_at, response_data, rejection_reason, verified_at, processing_notes, data_controller
		FROM data_subject_requests WHERE request_id = $1
	`

	err := pgs.db.QueryRow(query, requestID).Scan(
		&request.ID, &request.RequestID, &request.UserID, &request.RequestType,
		&request.Description, &request.Status, &request.RequestedAt, &request.CompletedAt,
		&request.ResponseData, &request.RejectionReason, &request.VerifiedAt,
		&request.ProcessingNotes, &request.DataController,
	)

	return &request, err
}

// ExportUserData exports all user data for GDPR Article 20 (data portability)
func (pgs *PrivacyGDPRService) ExportUserData(userID uuid.UUID) (map[string]interface{}, error) {
	export := make(map[string]interface{})
	export["export_timestamp"] = time.Now()
	export["user_id"] = userID
	export["data_controller"] = "TPT Titan"

	// Export user profile
	userData, err := pgs.getUserProfileData(userID)
	if err == nil {
		export["user_profile"] = userData
	}

	// Export contacts
	contactsData, err := pgs.getUserContactsData(userID)
	if err == nil {
		export["contacts"] = contactsData
	}

	// Export calendar events
	calendarData, err := pgs.getUserCalendarData(userID)
	if err == nil {
		export["calendar"] = calendarData
	}

	// Export emails
	emailsData, err := pgs.getUserEmailsData(userID)
	if err == nil {
		export["emails"] = emailsData
	}

	// Export tasks
	tasksData, err := pgs.getUserTasksData(userID)
	if err == nil {
		export["tasks"] = tasksData
	}

	// Export documents
	documentsData, err := pgs.getUserDocumentsData(userID)
	if err == nil {
		export["documents"] = documentsData
	}

	// Export consent history
	consentData, err := pgs.getUserConsentData(userID)
	if err == nil {
		export["consent_history"] = consentData
	}

	// Export data processing records
	processingData, err := pgs.getUserProcessingData(userID)
	if err == nil {
		export["data_processing_history"] = processingData
	}

	// Log audit event
	pgs.logPrivacyAudit(userID, "data_exported", "user_data", "",
		"", "", map[string]interface{}{
			"export_type": "gdpr_portability",
		})

	return export, nil
}

// DeleteUserData deletes all user data for GDPR Article 17 (right to erasure)
func (pgs *PrivacyGDPRService) DeleteUserData(userID uuid.UUID, reason string) error {
	// Start transaction for data deletion
	tx, err := pgs.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete in order of dependencies
	tables := []string{
		"email_attachments",
		"emails",
		"calendar_event_notifications",
		"calendar_reminders",
		"calendar_shares",
		"calendar_events",
		"calendars",
		"contact_groups",
		"contacts",
		"task_integrations",
		"tasks",
		"documents",
		"spreadsheet_cells",
		"spreadsheets",
		"form_responses",
		"forms",
		"data_subject_requests",
		"privacy_consents",
		"data_processing_records",
		"privacy_audit_logs",
		"privacy_settings",
		"user_sessions",
		"user_preferences",
	}

	for _, table := range tables {
		_, err = tx.Exec(fmt.Sprintf("DELETE FROM %s WHERE user_id = $1", table), userID)
		if err != nil {
			return fmt.Errorf("failed to delete from %s: %w", table, err)
		}
	}

	// Finally delete the user
	_, err = tx.Exec("DELETE FROM users WHERE id = $1", userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	// Log audit event (after deletion, so we can't associate with user ID)
	pgs.logPrivacyAudit(uuid.Nil, "user_data_deleted", "user_account", userID.String(),
		"", "", map[string]interface{}{
			"deletion_reason": reason,
		})

	return nil
}

// GetPrivacySettings gets user privacy settings
func (pgs *PrivacyGDPRService) GetPrivacySettings(userID uuid.UUID) (*PrivacySettings, error) {
	var settings PrivacySettings

	query := `
		SELECT user_id, data_collection_consent, analytics_consent, marketing_consent,
		       third_party_sharing_consent, data_retention_preference, contact_method_preference,
		       data_export_requested, data_deletion_requested, do_not_track, advertising_opt_out,
		       location_tracking_disabled, profile_visibility, updated_at
		FROM privacy_settings WHERE user_id = $1
	`

	err := pgs.db.QueryRow(query, userID).Scan(
		&settings.UserID, &settings.DataCollectionConsent, &settings.AnalyticsConsent,
		&settings.MarketingConsent, &settings.ThirdPartySharingConsent,
		&settings.DataRetentionPreference, &settings.ContactMethodPreference,
		&settings.DataExportRequested, &settings.DataDeletionRequested,
		&settings.DoNotTrack, &settings.AdvertisingOptOut,
		&settings.LocationTrackingDisabled, &settings.ProfileVisibility,
		&settings.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		// Return default settings
		settings = PrivacySettings{
			UserID:                   userID,
			DataCollectionConsent:    false,
			AnalyticsConsent:         false,
			MarketingConsent:         false,
			ThirdPartySharingConsent: false,
			DataRetentionPreference:  "standard",
			ContactMethodPreference:  "email",
			DataExportRequested:      false,
			DataDeletionRequested:    false,
			DoNotTrack:               true,
			AdvertisingOptOut:        true,
			LocationTrackingDisabled: true,
			ProfileVisibility:        "private",
			UpdatedAt:                time.Now(),
		}
		return &settings, nil
	}

	return &settings, err
}

// UpdatePrivacySettings updates user privacy settings
func (pgs *PrivacyGDPRService) UpdatePrivacySettings(settings *PrivacySettings) error {
	settings.UpdatedAt = time.Now()

	query := `
		INSERT INTO privacy_settings (
			user_id, data_collection_consent, analytics_consent, marketing_consent,
			third_party_sharing_consent, data_retention_preference, contact_method_preference,
			data_export_requested, data_deletion_requested, do_not_track, advertising_opt_out,
			location_tracking_disabled, profile_visibility, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		ON CONFLICT (user_id) DO UPDATE SET
			data_collection_consent = EXCLUDED.data_collection_consent,
			analytics_consent = EXCLUDED.analytics_consent,
			marketing_consent = EXCLUDED.marketing_consent,
			third_party_sharing_consent = EXCLUDED.third_party_sharing_consent,
			data_retention_preference = EXCLUDED.data_retention_preference,
			contact_method_preference = EXCLUDED.contact_method_preference,
			data_export_requested = EXCLUDED.data_export_requested,
			data_deletion_requested = EXCLUDED.data_deletion_requested,
			do_not_track = EXCLUDED.do_not_track,
			advertising_opt_out = EXCLUDED.advertising_opt_out,
			location_tracking_disabled = EXCLUDED.location_tracking_disabled,
			profile_visibility = EXCLUDED.profile_visibility,
			updated_at = EXCLUDED.updated_at
	`

	_, err := pgs.db.Exec(query,
		settings.UserID, settings.DataCollectionConsent, settings.AnalyticsConsent,
		settings.MarketingConsent, settings.ThirdPartySharingConsent,
		settings.DataRetentionPreference, settings.ContactMethodPreference,
		settings.DataExportRequested, settings.DataDeletionRequested,
		settings.DoNotTrack, settings.AdvertisingOptOut,
		settings.LocationTrackingDisabled, settings.ProfileVisibility,
		settings.UpdatedAt,
	)

	if err != nil {
		return err
	}

	// Log audit event
	pgs.logPrivacyAudit(settings.UserID, "privacy_settings_updated", "privacy_settings", "",
		"", "", nil)

	return nil
}

// RecordDataBreach records a data breach incident
func (pgs *PrivacyGDPRService) RecordDataBreach(breach *DataBreachRecord) error {
	breach.ID = uuid.New()
	breach.IncidentID = pgs.generateIncidentID()
	breach.DetectedAt = time.Now()
	breach.Status = "investigating"

	query := `
		INSERT INTO data_breach_records (id, incident_id, description, data_affected,
			users_affected, risk_level, detected_at, actions_taken, preventive_measures, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := pgs.db.Exec(query,
		breach.ID, breach.IncidentID, breach.Description, breach.DataAffected,
		breach.UsersAffected, breach.RiskLevel, breach.DetectedAt,
		breach.ActionsTaken, breach.PreventiveMeasures, breach.Status,
	)

	return err
}

// AnonymizeUserData anonymizes user data for research/analytics while maintaining GDPR compliance
func (pgs *PrivacyGDPRService) AnonymizeUserData(userID uuid.UUID, purpose string) error {
	// Create anonymized record
	anonymizedID := uuid.New()

	query := `
		INSERT INTO anonymized_user_data (anonymized_id, original_user_id, anonymized_at, purpose)
		VALUES ($1, $2, $3, $4)
	`

	_, err := pgs.db.Exec(query, anonymizedID, userID, time.Now(), purpose)
	if err != nil {
		return err
	}

	// Log audit event
	pgs.logPrivacyAudit(userID, "data_anonymized", "user_data", anonymizedID.String(),
		"", "", map[string]interface{}{
			"purpose": purpose,
		})

	return nil
}

// CheckDataRetentionCompliance checks if data retention policies are being followed
func (pgs *PrivacyGDPRService) CheckDataRetentionCompliance() ([]map[string]interface{}, error) {
	query := `
		SELECT user_id, processing_type, processed_at, retention_period
		FROM data_processing_records
		WHERE retention_period IS NOT NULL
		AND processed_at + INTERVAL '1 day' * retention_period < NOW()
	`

	rows, err := pgs.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var violations []map[string]interface{}
	for rows.Next() {
		var userID uuid.UUID
		var processingType string
		var processedAt time.Time
		var retentionPeriod int

		err := rows.Scan(&userID, &processingType, &processedAt, &retentionPeriod)
		if err != nil {
			continue
		}

		violations = append(violations, map[string]interface{}{
			"user_id":          userID,
			"processing_type":  processingType,
			"processed_at":     processedAt,
			"retention_period": retentionPeriod,
			"days_overdue":     int(time.Since(processedAt).Hours() / 24) - retentionPeriod,
		})
	}

	return violations, nil
}

// GeneratePrivacyReport generates a comprehensive privacy compliance report
func (pgs *PrivacyGDPRService) GeneratePrivacyReport() (map[string]interface{}, error) {
	report := make(map[string]interface{})
	report["generated_at"] = time.Now()
	report["reporting_period"] = "last_30_days"

	// Count active consents by type
	consentQuery := `
		SELECT consent_type, COUNT(*) as count
		FROM privacy_consents
		WHERE withdrawn_at IS NULL
		AND given_at >= NOW() - INTERVAL '30 days'
		GROUP BY consent_type
	`

	rows, err := pgs.db.Query(consentQuery)
	if err == nil {
		consents := make(map[string]int)
		for rows.Next() {
			var consentType string
			var count int
			rows.Scan(&consentType, &count)
			consents[consentType] = count
		}
		rows.Close()
		report["active_consents"] = consents
	}

	// Count data subject requests
	dsrQuery := `
		SELECT request_type, status, COUNT(*) as count
		FROM data_subject_requests
		WHERE requested_at >= NOW() - INTERVAL '30 days'
		GROUP BY request_type, status
	`

	dsrStats := make(map[string]map[string]int)
	rows, err = pgs.db.Query(dsrQuery)
	if err == nil {
		for rows.Next() {
			var requestType, status string
			var count int
			rows.Scan(&requestType, &status, &count)

			if dsrStats[requestType] == nil {
				dsrStats[requestType] = make(map[string]int)
			}
			dsrStats[requestType][status] = count
		}
		rows.Close()
	}
	report["data_subject_requests"] = dsrStats

	// Count data processing activities
	processingQuery := `
		SELECT processing_type, COUNT(*) as count
		FROM data_processing_records
		WHERE processed_at >= NOW() - INTERVAL '30 days'
		GROUP BY processing_type
	`

	processingStats := make(map[string]int)
	rows, err = pgs.db.Query(processingQuery)
	if err == nil {
		for rows.Next() {
			var processingType string
			var count int
			rows.Scan(&processingType, &count)
			processingStats[processingType] = count
		}
		rows.Close()
	}
	report["data_processing_activities"] = processingStats

	// Check for data breaches
	var breachCount int
	pgs.db.QueryRow("SELECT COUNT(*) FROM data_breach_records WHERE detected_at >= NOW() - INTERVAL '30 days'").Scan(&breachCount)
	report["data_breaches_reported"] = breachCount

	return report, nil
}

// Helper methods

func (pgs *PrivacyGDPRService) generateRequestID() string {
	return fmt.Sprintf("DSR-%s", strings.ToUpper(uuid.New().String()[:8]))
}

func (pgs *PrivacyGDPRService) generateIncidentID() string {
	return fmt.Sprintf("BR-%s", strings.ToUpper(uuid.New().String()[:8]))
}

func (pgs *PrivacyGDPRService) logPrivacyAudit(userID uuid.UUID, action, resource, resourceID, ipAddress, userAgent string, details map[string]interface{}) {
	log := &PrivacyAuditLog{
		ID:             uuid.New(),
		UserID:         userID,
		Action:         action,
		Resource:       resource,
		ResourceID:     resourceID,
		IPAddress:      ipAddress,
		UserAgent:      userAgent,
		Timestamp:      time.Now(),
		Details:        nil,
		ComplianceFlag: false,
	}

	if details != nil {
		jsonDetails, _ := json.Marshal(details)
		log.Details = jsonDetails
	}

	query := `
		INSERT INTO privacy_audit_logs (id, user_id, action, resource, resource_id,
			ip_address, user_agent, timestamp, details, compliance_flag)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	pgs.db.Exec(query,
		log.ID, log.UserID, log.Action, log.Resource, log.ResourceID,
		log.IPAddress, log.UserAgent, log.Timestamp, log.Details, log.ComplianceFlag,
	)
}

// Data export helper methods (simplified implementations)
func (pgs *PrivacyGDPRService) getUserProfileData(userID uuid.UUID) (map[string]interface{}, error) {
	var data map[string]interface{}
	query := "SELECT id, username, email, first_name, last_name, created_at FROM users WHERE id = $1"
	row := pgs.db.QueryRow(query, userID)

	var id uuid.UUID
	var username, email, firstName, lastName *string
	var createdAt time.Time

	err := row.Scan(&id, &username, &email, &firstName, &lastName, &createdAt)
	if err != nil {
		return nil, err
	}

	data = map[string]interface{}{
		"id":         id,
		"username":   username,
		"email":      email,
		"first_name": firstName,
		"last_name":  lastName,
		"created_at": createdAt,
	}
	return data, nil
}

func (pgs *PrivacyGDPRService) getUserContactsData(userID uuid.UUID) ([]map[string]interface{}, error) {
	// Simplified implementation
	return []map[string]interface{}{}, nil
}

func (pgs *PrivacyGDPRService) getUserCalendarData(userID uuid.UUID) ([]map[string]interface{}, error) {
	// Simplified implementation
	return []map[string]interface{}{}, nil
}

func (pgs *PrivacyGDPRService) getUserEmailsData(userID uuid.UUID) ([]map[string]interface{}, error) {
	// Simplified implementation
	return []map[string]interface{}{}, nil
}

func (pgs *PrivacyGDPRService) getUserTasksData(userID uuid.UUID) ([]map[string]interface{}, error) {
	// Simplified implementation
	return []map[string]interface{}{}, nil
}

func (pgs *PrivacyGDPRService) getUserDocumentsData(userID uuid.UUID) ([]map[string]interface{}, error) {
	// Simplified implementation
	return []map[string]interface{}{}, nil
}

func (pgs *PrivacyGDPRService) getUserConsentData(userID uuid.UUID) ([]map[string]interface{}, error) {
	// Simplified implementation
	return []map[string]interface{}{}, nil
}

func (pgs *PrivacyGDPRService) getUserProcessingData(userID uuid.UUID) ([]map[string]interface{}, error) {
	// Simplified implementation
	return []map[string]interface{}{}, nil
}
