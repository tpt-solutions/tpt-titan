package services

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// FormRelationshipService manages relationships between forms and data
type FormRelationshipService struct {
	db *sql.DB
}

// Relationship represents a relationship between two forms/tables
type Relationship struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	SourceFormID    uuid.UUID `json:"source_form_id"`
	SourceField     string    `json:"source_field"`
	TargetFormID    uuid.UUID `json:"target_form_id"`
	TargetField     string    `json:"target_field"`
	RelationshipType string   `json:"relationship_type"` // "one-to-one", "one-to-many", "many-to-many"
	CascadeDelete   bool      `json:"cascade_delete"`
	CascadeUpdate   bool      `json:"cascade_update"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// LookupField represents a lookup field that references another form
type LookupField struct {
	ID             uuid.UUID `json:"id"`
	FormID         uuid.UUID `json:"form_id"`
	FieldID        uuid.UUID `json:"field_id"`
	RelatedFormID  uuid.UUID `json:"related_form_id"`
	RelatedFieldID uuid.UUID `json:"related_field_id"`
	DisplayField   string    `json:"display_field"`   // Field to display in dropdown
	FilterField    string    `json:"filter_field,omitempty"`    // Field to filter by
	FilterValue    string    `json:"filter_value,omitempty"`    // Value to filter by
	AllowMultiple  bool      `json:"allow_multiple"` // Allow multiple selections
	CreatedAt      time.Time `json:"created_at"`
}

// FormIntegration represents integration between forms and external systems
type FormIntegration struct {
	ID          uuid.UUID            `json:"id"`
	FormID      uuid.UUID            `json:"form_id"`
	Name        string               `json:"name"`
	Type        string               `json:"type"`        // "webhook", "api", "email", "sms"
	Trigger     string               `json:"trigger"`     // "on_submit", "on_update", "on_delete"
	Config      map[string]interface{} `json:"config"`      // Configuration specific to integration type
	IsActive    bool                 `json:"is_active"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
}

// NewFormRelationshipService creates a new form relationship service
func NewFormRelationshipService(db *sql.DB) *FormRelationshipService {
	return &FormRelationshipService{db: db}
}

// CreateRelationship creates a new relationship between forms
func (frs *FormRelationshipService) CreateRelationship(rel *Relationship) error {
	rel.ID = uuid.New()
	rel.CreatedAt = time.Now()
	rel.UpdatedAt = time.Now()

	query := `
		INSERT INTO form_relationships (id, name, description, source_form_id, source_field,
			target_form_id, target_field, relationship_type, cascade_delete, cascade_update,
			created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err := frs.db.Exec(query,
		rel.ID, rel.Name, rel.Description, rel.SourceFormID, rel.SourceField,
		rel.TargetFormID, rel.TargetField, rel.RelationshipType, rel.CascadeDelete,
		rel.CascadeUpdate, rel.CreatedAt, rel.UpdatedAt,
	)

	return err
}

// GetRelationshipsByForm gets all relationships for a specific form
func (frs *FormRelationshipService) GetRelationshipsByForm(formID uuid.UUID) ([]Relationship, error) {
	query := `
		SELECT id, name, description, source_form_id, source_field, target_form_id,
		       target_field, relationship_type, cascade_delete, cascade_update,
		       created_at, updated_at
		FROM form_relationships
		WHERE source_form_id = $1 OR target_form_id = $1
		ORDER BY created_at
	`

	rows, err := frs.db.Query(query, formID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var relationships []Relationship
	for rows.Next() {
		var rel Relationship
		err := rows.Scan(
			&rel.ID, &rel.Name, &rel.Description, &rel.SourceFormID, &rel.SourceField,
			&rel.TargetFormID, &rel.TargetField, &rel.RelationshipType, &rel.CascadeDelete,
			&rel.CascadeUpdate, &rel.CreatedAt, &rel.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		relationships = append(relationships, rel)
	}

	return relationships, nil
}

// DeleteRelationship deletes a relationship
func (frs *FormRelationshipService) DeleteRelationship(relationshipID uuid.UUID) error {
	query := `DELETE FROM form_relationships WHERE id = $1`
	_, err := frs.db.Exec(query, relationshipID)
	return err
}

// ValidateRelationship checks if a relationship is valid
func (frs *FormRelationshipService) ValidateRelationship(rel *Relationship) []string {
	var errors []string

	// Check if forms exist
	if !frs.formExists(rel.SourceFormID) {
		errors = append(errors, "Source form does not exist")
	}
	if !frs.formExists(rel.TargetFormID) {
		errors = append(errors, "Target form does not exist")
	}

	// Check if fields exist
	if !frs.fieldExists(rel.SourceFormID, rel.SourceField) {
		errors = append(errors, "Source field does not exist")
	}
	if !frs.fieldExists(rel.TargetFormID, rel.TargetField) {
		errors = append(errors, "Target field does not exist")
	}

	// Validate relationship type
	validTypes := []string{"one-to-one", "one-to-many", "many-to-many"}
	validType := false
	for _, t := range validTypes {
		if rel.RelationshipType == t {
			validType = true
			break
		}
	}
	if !validType {
		errors = append(errors, "Invalid relationship type")
	}

	return errors
}

// CreateLookupField creates a lookup field that references another form
func (frs *FormRelationshipService) CreateLookupField(lookup *LookupField) error {
	lookup.ID = uuid.New()
	lookup.CreatedAt = time.Now()

	query := `
		INSERT INTO form_lookup_fields (id, form_id, field_id, related_form_id,
			related_field_id, display_field, filter_field, filter_value,
			allow_multiple, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := frs.db.Exec(query,
		lookup.ID, lookup.FormID, lookup.FieldID, lookup.RelatedFormID,
		lookup.RelatedFieldID, lookup.DisplayField, lookup.FilterField,
		lookup.FilterValue, lookup.AllowMultiple, lookup.CreatedAt,
	)

	return err
}

// GetLookupFieldData gets the data for a lookup field dropdown
func (frs *FormRelationshipService) GetLookupFieldData(lookupFieldID uuid.UUID) ([]map[string]interface{}, error) {
	// Get lookup field configuration
	var lookup LookupField
	query := `
		SELECT related_form_id, related_field_id, display_field, filter_field, filter_value
		FROM form_lookup_fields
		WHERE id = $1
	`

	err := frs.db.QueryRow(query, lookupFieldID).Scan(
		&lookup.RelatedFormID, &lookup.RelatedFieldID, &lookup.DisplayField,
		&lookup.FilterField, &lookup.FilterValue,
	)
	if err != nil {
		return nil, err
	}

	// Build query to get lookup data
	selectQuery := fmt.Sprintf(`
		SELECT %s, %s
		FROM form_responses
		WHERE form_id = $1
	`, lookup.RelatedFieldID, lookup.DisplayField)

	args := []interface{}{lookup.RelatedFormID}
	argCount := 1

	// Add filter if specified
	if lookup.FilterField != "" && lookup.FilterValue != "" {
		selectQuery += fmt.Sprintf(" AND %s = $%d", lookup.FilterField, argCount+1)
		args = append(args, lookup.FilterValue)
	}

	selectQuery += " ORDER BY " + lookup.DisplayField

	rows, err := frs.db.Query(selectQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var id, display interface{}
		err := rows.Scan(&id, &display)
		if err != nil {
			continue
		}

		results = append(results, map[string]interface{}{
			"id":      id,
			"display": display,
		})
	}

	return results, nil
}

// CreateFormIntegration creates a new form integration
func (frs *FormRelationshipService) CreateFormIntegration(integration *FormIntegration) error {
	integration.ID = uuid.New()
	integration.CreatedAt = time.Now()
	integration.UpdatedAt = time.Now()

	query := `
		INSERT INTO form_integrations (id, form_id, name, type, trigger, config,
			is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := frs.db.Exec(query,
		integration.ID, integration.FormID, integration.Name, integration.Type,
		integration.Trigger, integration.Config, integration.IsActive,
		integration.CreatedAt, integration.UpdatedAt,
	)

	return err
}

// ExecuteFormIntegration executes a form integration
func (frs *FormRelationshipService) ExecuteFormIntegration(integrationID uuid.UUID, formData map[string]interface{}) error {
	// Get integration configuration
	var integration FormIntegration
	query := `
		SELECT form_id, type, trigger, config
		FROM form_integrations
		WHERE id = $1 AND is_active = true
	`

	err := frs.db.QueryRow(query, integrationID).Scan(
		&integration.FormID, &integration.Type, &integration.Trigger, &integration.Config,
	)
	if err != nil {
		return err
	}

	// Execute based on integration type
	switch integration.Type {
	case "webhook":
		return frs.executeWebhookIntegration(integration, formData)
	case "email":
		return frs.executeEmailIntegration(integration, formData)
	case "api":
		return frs.executeAPIIntegration(integration, formData)
	default:
		return fmt.Errorf("unsupported integration type: %s", integration.Type)
	}
}

// GetRelatedData gets related data based on relationships
func (frs *FormRelationshipService) GetRelatedData(formID uuid.UUID, recordID uuid.UUID) (map[string]interface{}, error) {
	relationships, err := frs.GetRelationshipsByForm(formID)
	if err != nil {
		return nil, err
	}

	relatedData := make(map[string]interface{})

	for _, rel := range relationships {
		var query string
		var args []interface{}

		if rel.SourceFormID == formID {
			// This form is the source, get related target records
			query = fmt.Sprintf(`
				SELECT fr.* FROM form_responses fr
				WHERE fr.form_id = $1 AND fr.%s = $2
			`, rel.TargetField)
			args = []interface{}{rel.TargetFormID, recordID}
		} else {
			// This form is the target, get related source records
			query = fmt.Sprintf(`
				SELECT fr.* FROM form_responses fr
				WHERE fr.form_id = $1 AND fr.%s = $2
			`, rel.SourceField)
			args = []interface{}{rel.SourceFormID, recordID}
		}

		rows, err := frs.db.Query(query, args...)
		if err != nil {
			continue // Skip this relationship if query fails
		}
		defer rows.Close()

		var records []map[string]interface{}
		columns, _ := rows.Columns()

		for rows.Next() {
			values := make([]interface{}, len(columns))
			valuePtrs := make([]interface{}, len(columns))
			for i := range values {
				valuePtrs[i] = &values[i]
			}

			err := rows.Scan(valuePtrs...)
			if err != nil {
				continue
			}

			record := make(map[string]interface{})
			for i, col := range columns {
				record[col] = values[i]
			}
			records = append(records, record)
		}

		relatedData[rel.Name] = records
	}

	return relatedData, nil
}

// Helper methods

func (frs *FormRelationshipService) formExists(formID uuid.UUID) bool {
	var exists bool
	frs.db.QueryRow("SELECT EXISTS(SELECT 1 FROM forms WHERE id = $1)", formID).Scan(&exists)
	return exists
}

func (frs *FormRelationshipService) fieldExists(formID uuid.UUID, fieldName string) bool {
	var exists bool
	query := `
		SELECT EXISTS(
			SELECT 1 FROM form_fields ff
			JOIN forms f ON ff.form_id = f.id
			WHERE f.id = $1 AND ff.field_name = $2
		)
	`
	frs.db.QueryRow(query, formID, fieldName).Scan(&exists)
	return exists
}

func (frs *FormRelationshipService) executeWebhookIntegration(integration FormIntegration, formData map[string]interface{}) error {
	webhookURL, ok := integration.Config["url"].(string)
	if !ok {
		return fmt.Errorf("webhook URL not configured")
	}

	// In a real implementation, make HTTP request to webhook URL
	// For now, just log
	fmt.Printf("Executing webhook integration to %s with data: %v\n", webhookURL, formData)
	return nil
}

func (frs *FormRelationshipService) executeEmailIntegration(integration FormIntegration, formData map[string]interface{}) error {
	// In a real implementation, send email using email service
	// For now, just log
	fmt.Printf("Executing email integration with data: %v\n", formData)
	return nil
}

func (frs *FormRelationshipService) executeAPIIntegration(integration FormIntegration, formData map[string]interface{}) error {
	apiURL, ok := integration.Config["url"].(string)
	if !ok {
		return fmt.Errorf("API URL not configured")
	}

	// In a real implementation, make API call
	// For now, just log
	fmt.Printf("Executing API integration to %s with data: %v\n", apiURL, formData)
	return nil
}

// GetFormHierarchy gets the hierarchy of related forms
func (frs *FormRelationshipService) GetFormHierarchy(formID uuid.UUID) (map[string]interface{}, error) {
	hierarchy := map[string]interface{}{
		"form_id":    formID,
		"children":   []map[string]interface{}{},
		"parents":    []map[string]interface{}{},
	}

	// Get child relationships (where this form is the source)
	childrenQuery := `
		SELECT id, name, target_form_id
		FROM form_relationships
		WHERE source_form_id = $1
	`

	childrenRows, err := frs.db.Query(childrenQuery, formID)
	if err == nil {
		for childrenRows.Next() {
			var relID uuid.UUID
			var relName string
			var targetFormID uuid.UUID

			childrenRows.Scan(&relID, &relName, &targetFormID)
			hierarchy["children"] = append(hierarchy["children"].([]map[string]interface{}), map[string]interface{}{
				"relationship_id": relID,
				"name":           relName,
				"form_id":        targetFormID,
			})
		}
		childrenRows.Close()
	}

	// Get parent relationships (where this form is the target)
	parentsQuery := `
		SELECT id, name, source_form_id
		FROM form_relationships
		WHERE target_form_id = $1
	`

	parentsRows, err := frs.db.Query(parentsQuery, formID)
	if err == nil {
		for parentsRows.Next() {
			var relID uuid.UUID
			var relName string
			var sourceFormID uuid.UUID

			parentsRows.Scan(&relID, &relName, &sourceFormID)
			hierarchy["parents"] = append(hierarchy["parents"].([]map[string]interface{}), map[string]interface{}{
				"relationship_id": relID,
				"name":           relName,
				"form_id":        sourceFormID,
			})
		}
		parentsRows.Close()
	}

	return hierarchy, nil
}

// ValidateFormDataWithRelationships validates form data considering relationships
func (frs *FormRelationshipService) ValidateFormDataWithRelationships(formID uuid.UUID, formData map[string]interface{}) []string {
	var errors []string

	relationships, err := frs.GetRelationshipsByForm(formID)
	if err != nil {
		errors = append(errors, "Failed to check relationships")
		return errors
	}

	for _, rel := range relationships {
		if rel.SourceFormID == formID {
			// Check if referenced record exists
			if value, exists := formData[rel.SourceField]; exists && value != nil {
				var exists bool
				query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM form_responses WHERE form_id = $1 AND %s = $2)", rel.TargetField)
				frs.db.QueryRow(query, rel.TargetFormID, value).Scan(&exists)
				if !exists {
					errors = append(errors, fmt.Sprintf("Referenced record in %s does not exist", rel.Name))
				}
			}
		}
	}

	return errors
}
