package services

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"tpt-titan/backend/models"
)

// TaskIntegrationService provides integration between tasks and other TPT Titan components
type TaskIntegrationService struct {
	db              *sql.DB
	taskService     *TaskService
	emailService    *EmailService
	formService     interface{} // Form service interface (placeholder)
	calendarService *CalendarService
	userID          uuid.UUID
}

// TaskIntegration represents an integration configuration
type TaskIntegration struct {
	ID              uuid.UUID              `json:"id"`
	TaskID          uuid.UUID              `json:"task_id"`
	SourceType      string                 `json:"source_type"`      // "form", "email", "calendar", "chat"
	SourceID        uuid.UUID              `json:"source_id"`        // ID of the source item
	IntegrationType string                 `json:"integration_type"` // "auto_create", "link", "sync"
	Config          map[string]interface{} `json:"config"`           // Integration-specific configuration
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

// FormTaskIntegration represents integration between forms and tasks
type FormTaskIntegration struct {
	ID                  uuid.UUID `json:"id"`
	FormID              uuid.UUID `json:"form_id"`
	TaskTemplateID      uuid.UUID `json:"task_template_id,omitempty"`
	TriggerCondition    string    `json:"trigger_condition"` // JSON condition for when to create task
	AutoAssign          bool      `json:"auto_assign"`
	AssigneeField       string    `json:"assignee_field,omitempty"` // Form field containing assignee
	PriorityField       string    `json:"priority_field,omitempty"` // Form field containing priority
	DueDateField        string    `json:"due_date_field,omitempty"` // Form field containing due date
	TitleTemplate       string    `json:"title_template"`           // Template for task title
	DescriptionTemplate string    `json:"description_template"`     // Template for task description
	IsActive            bool      `json:"is_active"`
	CreatedAt           time.Time `json:"created_at"`
}

// EmailTaskIntegration represents integration between emails and tasks
type EmailTaskIntegration struct {
	ID                  uuid.UUID              `json:"id"`
	EmailFilter         map[string]interface{} `json:"email_filter"` // Conditions for email processing
	TaskTemplateID      uuid.UUID              `json:"task_template_id,omitempty"`
	AutoCreate          bool                   `json:"auto_create"`
	PriorityMapping     map[string]string      `json:"priority_mapping"` // Email priority to task priority
	CategoryMapping     map[string]string      `json:"category_mapping"` // Email category to task category
	AssigneeMapping     map[string]uuid.UUID   `json:"assignee_mapping"` // Email sender/domain to assignee
	TitleTemplate       string                 `json:"title_template"`
	DescriptionTemplate string                 `json:"description_template"`
	IsActive            bool                   `json:"is_active"`
	CreatedAt           time.Time              `json:"created_at"`
}

// TaskTemplate represents a reusable task template
type TaskTemplate struct {
	ID                uuid.UUID              `json:"id"`
	Name              string                 `json:"name"`
	Description       string                 `json:"description"`
	Category          string                 `json:"category"`
	Priority          string                 `json:"priority"`
	EstimatedDuration *int                   `json:"estimated_duration,omitempty"` // in minutes
	Tags              []string               `json:"tags"`
	Checklist         []TaskChecklistItem    `json:"checklist"`
	Subtasks          []SubtaskTemplate      `json:"subtasks"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
	CreatedBy         uuid.UUID              `json:"created_by"`
	IsPublic          bool                   `json:"is_public"`
	UsageCount        int                    `json:"usage_count"`
	CreatedAt         time.Time              `json:"created_at"`
}

// TaskChecklistItem represents an item in a task checklist
type TaskChecklistItem struct {
	Text     string `json:"text"`
	Required bool   `json:"required"`
	Order    int    `json:"order"`
}

// SubtaskTemplate represents a template for subtasks
type SubtaskTemplate struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	Order       int    `json:"order"`
}

// NewTaskIntegrationService creates a new task integration service
func NewTaskIntegrationService(db *sql.DB, taskSvc *TaskService, emailSvc *EmailService, formSvc interface{}, calSvc *CalendarService, userID uuid.UUID) *TaskIntegrationService {
	return &TaskIntegrationService{
		db:              db,
		taskService:     taskSvc,
		emailService:    emailSvc,
		formService:     formSvc,
		calendarService: calSvc,
		userID:          userID,
	}
}

// CreateFormTaskIntegration creates integration between a form and task creation
func (tis *TaskIntegrationService) CreateFormTaskIntegration(integration *FormTaskIntegration) error {
	integration.ID = uuid.New()
	integration.CreatedAt = time.Now()

	query := `
		INSERT INTO form_task_integrations (id, form_id, task_template_id, trigger_condition,
			auto_assign, assignee_field, priority_field, due_date_field, title_template,
			description_template, is_active, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err := tis.db.Exec(query,
		integration.ID, integration.FormID, integration.TaskTemplateID, integration.TriggerCondition,
		integration.AutoAssign, integration.AssigneeField, integration.PriorityField,
		integration.DueDateField, integration.TitleTemplate, integration.DescriptionTemplate,
		integration.IsActive, integration.CreatedAt,
	)

	return err
}

// ProcessFormSubmission creates tasks based on form submissions
func (tis *TaskIntegrationService) ProcessFormSubmission(formID uuid.UUID, submissionData map[string]interface{}) ([]uuid.UUID, error) {
	// Get active integrations for this form
	integrations, err := tis.getFormTaskIntegrations(formID)
	if err != nil {
		return nil, err
	}

	var createdTaskIDs []uuid.UUID

	for _, integration := range integrations {
		if !integration.IsActive {
			continue
		}

		// Check if trigger condition is met
		if tis.evaluateTriggerCondition(integration.TriggerCondition, submissionData) {
			taskID, err := tis.createTaskFromFormSubmission(integration, submissionData)
			if err != nil {
				// Log error but continue with other integrations
				continue
			}
			createdTaskIDs = append(createdTaskIDs, taskID)

			// Create integration record
			tis.createTaskIntegrationRecord(taskID, "form", formID, "auto_create", nil)
		}
	}

	return createdTaskIDs, nil
}

// CreateEmailTaskIntegration creates integration between emails and tasks
func (tis *TaskIntegrationService) CreateEmailTaskIntegration(integration *EmailTaskIntegration) error {
	integration.ID = uuid.New()
	integration.CreatedAt = time.Now()

	query := `
		INSERT INTO email_task_integrations (id, email_filter, task_template_id, auto_create,
			priority_mapping, category_mapping, assignee_mapping, title_template,
			description_template, is_active, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := tis.db.Exec(query,
		integration.ID, integration.EmailFilter, integration.TaskTemplateID, integration.AutoCreate,
		integration.PriorityMapping, integration.CategoryMapping, integration.AssigneeMapping,
		integration.TitleTemplate, integration.DescriptionTemplate, integration.IsActive,
		integration.CreatedAt,
	)

	return err
}

// ProcessIncomingEmail creates tasks based on incoming emails
func (tis *TaskIntegrationService) ProcessIncomingEmail(emailID uuid.UUID, emailData map[string]interface{}) ([]uuid.UUID, error) {
	// Get active email integrations
	integrations, err := tis.getEmailTaskIntegrations()
	if err != nil {
		return nil, err
	}

	var createdTaskIDs []uuid.UUID

	for _, integration := range integrations {
		if !integration.IsActive || !integration.AutoCreate {
			continue
		}

		// Check if email matches filter criteria
		if tis.emailMatchesFilter(emailData, integration.EmailFilter) {
			taskID, err := tis.createTaskFromEmail(integration, emailData)
			if err != nil {
				// Log error but continue with other integrations
				continue
			}
			createdTaskIDs = append(createdTaskIDs, taskID)

			// Create integration record
			tis.createTaskIntegrationRecord(taskID, "email", emailID, "auto_create", nil)
		}
	}

	return createdTaskIDs, nil
}

// CreateTaskTemplate creates a new task template
func (tis *TaskIntegrationService) CreateTaskTemplate(template *TaskTemplate) error {
	template.ID = uuid.New()
	template.CreatedAt = time.Now()
	template.UsageCount = 0

	query := `
		INSERT INTO task_templates (id, name, description, category, priority, estimated_duration,
			tags, checklist, subtasks, metadata, created_by, is_public, usage_count, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`

	_, err := tis.db.Exec(query,
		template.ID, template.Name, template.Description, template.Category, template.Priority,
		template.EstimatedDuration, template.Tags, template.Checklist, template.Subtasks,
		template.Metadata, template.CreatedBy, template.IsPublic, template.UsageCount, template.CreatedAt,
	)

	return err
}

// CreateTaskFromTemplate creates a task from a template
func (tis *TaskIntegrationService) CreateTaskFromTemplate(templateID uuid.UUID, customizations map[string]interface{}) (uuid.UUID, error) {
	template, err := tis.getTaskTemplate(templateID)
	if err != nil {
		return uuid.Nil, err
	}

	// Increment usage count
	tis.incrementTemplateUsage(templateID)

	// Create the task via the real task service so it is persisted and a real
	// ID is returned (previously this fabricated a UUID, risking orphans).
	subtasks := make([]models.SubtaskRequest, 0, len(template.Subtasks))
	for _, subtask := range template.Subtasks {
		subtasks = append(subtasks, models.SubtaskRequest{
			Title: subtask.Title,
			ID:    uuid.New().String(),
		})
	}

	task, err := tis.taskService.CreateTask(tis.userID, models.TaskRequest{
		Title:       template.Name,
		Description: template.Description,
		Status:      "todo",
		Priority:    template.Priority,
		Tags:        template.Tags,
		Subtasks:    subtasks,
	})
	if err != nil {
		return uuid.Nil, err
	}

	return task.ID, nil
}

// LinkTaskToItem links a task to another TPT Titan item
func (tis *TaskIntegrationService) LinkTaskToItem(taskID uuid.UUID, itemType string, itemID uuid.UUID, linkType string) error {
	integration := &TaskIntegration{
		TaskID:          taskID,
		SourceType:      itemType,
		SourceID:        itemID,
		IntegrationType: linkType,
		Config:          make(map[string]interface{}),
	}

	return tis.createTaskIntegrationRecord(taskID, itemType, itemID, linkType, integration.Config)
}

// GetTaskIntegrations gets all integrations for a task
func (tis *TaskIntegrationService) GetTaskIntegrations(taskID uuid.UUID) ([]TaskIntegration, error) {
	query := `
		SELECT id, task_id, source_type, source_id, integration_type, config, created_at, updated_at
		FROM task_integrations
		WHERE task_id = $1
		ORDER BY created_at
	`

	rows, err := tis.db.Query(query, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var integrations []TaskIntegration
	for rows.Next() {
		var integration TaskIntegration
		err := rows.Scan(
			&integration.ID, &integration.TaskID, &integration.SourceType,
			&integration.SourceID, &integration.IntegrationType, &integration.Config,
			&integration.CreatedAt, &integration.UpdatedAt,
		)
		if err != nil {
			continue
		}
		integrations = append(integrations, integration)
	}

	return integrations, nil
}

// GetTaskTemplates gets available task templates
func (tis *TaskIntegrationService) GetTaskTemplates(category string, includePublic bool) ([]TaskTemplate, error) {
	query := `SELECT id, name, description, category, priority, estimated_duration, tags,
		checklist, subtasks, metadata, created_by, is_public, usage_count, created_at
		FROM task_templates`

	args := []interface{}{}
	argCount := 0

	conditions := []string{}
	if category != "" {
		argCount++
		conditions = append(conditions, fmt.Sprintf("category = $%d", argCount))
		args = append(args, category)
	}

	if !includePublic {
		argCount++
		conditions = append(conditions, fmt.Sprintf("is_public = true OR created_by = $%d", argCount))
		args = append(args, "current_user_id") // Would be replaced with actual user ID
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY usage_count DESC, created_at DESC"

	rows, err := tis.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []TaskTemplate
	for rows.Next() {
		var template TaskTemplate
		err := rows.Scan(
			&template.ID, &template.Name, &template.Description, &template.Category,
			&template.Priority, &template.EstimatedDuration, &template.Tags,
			&template.Checklist, &template.Subtasks, &template.Metadata,
			&template.CreatedBy, &template.IsPublic, &template.UsageCount, &template.CreatedAt,
		)
		if err != nil {
			continue
		}
		templates = append(templates, template)
	}

	return templates, nil
}

// Helper methods

func (tis *TaskIntegrationService) getFormTaskIntegrations(formID uuid.UUID) ([]FormTaskIntegration, error) {
	query := `
		SELECT id, form_id, task_template_id, trigger_condition, auto_assign,
		       assignee_field, priority_field, due_date_field, title_template,
		       description_template, is_active, created_at
		FROM form_task_integrations
		WHERE form_id = $1 AND is_active = true
	`

	rows, err := tis.db.Query(query, formID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var integrations []FormTaskIntegration
	for rows.Next() {
		var integration FormTaskIntegration
		err := rows.Scan(
			&integration.ID, &integration.FormID, &integration.TaskTemplateID,
			&integration.TriggerCondition, &integration.AutoAssign, &integration.AssigneeField,
			&integration.PriorityField, &integration.DueDateField, &integration.TitleTemplate,
			&integration.DescriptionTemplate, &integration.IsActive, &integration.CreatedAt,
		)
		if err != nil {
			continue
		}
		integrations = append(integrations, integration)
	}

	return integrations, nil
}

func (tis *TaskIntegrationService) getEmailTaskIntegrations() ([]EmailTaskIntegration, error) {
	query := `
		SELECT id, email_filter, task_template_id, auto_create, priority_mapping,
		       category_mapping, assignee_mapping, title_template, description_template,
		       is_active, created_at
		FROM email_task_integrations
		WHERE is_active = true
	`

	rows, err := tis.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var integrations []EmailTaskIntegration
	for rows.Next() {
		var integration EmailTaskIntegration
		err := rows.Scan(
			&integration.ID, &integration.EmailFilter, &integration.TaskTemplateID,
			&integration.AutoCreate, &integration.PriorityMapping, &integration.CategoryMapping,
			&integration.AssigneeMapping, &integration.TitleTemplate, &integration.DescriptionTemplate,
			&integration.IsActive, &integration.CreatedAt,
		)
		if err != nil {
			continue
		}
		integrations = append(integrations, integration)
	}

	return integrations, nil
}

func (tis *TaskIntegrationService) evaluateTriggerCondition(condition string, data map[string]interface{}) bool {
	// Simple condition evaluation (would need full expression parser)
	// For now, assume conditions are met
	return true
}

func (tis *TaskIntegrationService) createTaskFromFormSubmission(integration FormTaskIntegration, submissionData map[string]interface{}) (uuid.UUID, error) {
	taskID := uuid.New() // Reserved so the integration record can reference it even if persistence fails.

	if tis.taskService == nil {
		return taskID, fmt.Errorf("task service not configured")
	}

	req := models.TaskRequest{
		Title:       tis.processTemplate(integration.TitleTemplate, submissionData),
		Description: tis.processTemplate(integration.DescriptionTemplate, submissionData),
		Status:      "todo",
		Priority:    "medium",
	}

	// Set assignee if auto-assign is enabled
	if integration.AutoAssign && integration.AssigneeField != "" {
		if assigneeValue, exists := submissionData[integration.AssigneeField]; exists {
			req.AssignedTo = fmt.Sprintf("%v", assigneeValue)
		}
	}

	// Set priority if field is specified
	if integration.PriorityField != "" {
		if priorityValue, exists := submissionData[integration.PriorityField]; exists {
			req.Priority = fmt.Sprintf("%v", priorityValue)
		}
	}

	// Set due date if field is specified
	if integration.DueDateField != "" {
		if dueDateValue, exists := submissionData[integration.DueDateField]; exists {
			if t, err := time.Parse(time.RFC3339, fmt.Sprintf("%v", dueDateValue)); err == nil {
				req.DueDate = &t
			}
		}
	}

	task, err := tis.taskService.CreateTask(tis.userID, req)
	if err != nil {
		return taskID, err
	}
	return task.ID, nil
}

func (tis *TaskIntegrationService) createTaskFromEmail(integration EmailTaskIntegration, emailData map[string]interface{}) (uuid.UUID, error) {
	taskID := uuid.New() // Reserved so the integration record can reference it even if persistence fails.

	if tis.taskService == nil {
		return taskID, fmt.Errorf("task service not configured")
	}

	req := models.TaskRequest{
		Title:       tis.processTemplate(integration.TitleTemplate, emailData),
		Description: tis.processTemplate(integration.DescriptionTemplate, emailData),
		Status:      "todo",
		Priority:    "medium",
	}

	// Apply priority mapping
	if sender, exists := emailData["sender"].(string); exists {
		if priority, mapped := integration.PriorityMapping[sender]; mapped {
			req.Priority = priority
		}
		if assigneeID, mapped := integration.AssigneeMapping[sender]; mapped {
			req.AssignedTo = assigneeID.String()
		}
	}

	task, err := tis.taskService.CreateTask(tis.userID, req)
	if err != nil {
		return taskID, err
	}
	return task.ID, nil
}

func (tis *TaskIntegrationService) createTaskIntegrationRecord(taskID uuid.UUID, sourceType string, sourceID uuid.UUID, integrationType string, config map[string]interface{}) error {
	integration := &TaskIntegration{
		ID:              uuid.New(),
		TaskID:          taskID,
		SourceType:      sourceType,
		SourceID:        sourceID,
		IntegrationType: integrationType,
		Config:          config,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	query := `
		INSERT INTO task_integrations (id, task_id, source_type, source_id, integration_type, config, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := tis.db.Exec(query,
		integration.ID, integration.TaskID, integration.SourceType, integration.SourceID,
		integration.IntegrationType, integration.Config, integration.CreatedAt, integration.UpdatedAt,
	)

	return err
}

func (tis *TaskIntegrationService) emailMatchesFilter(emailData, filter map[string]interface{}) bool {
	// Simple filter matching (would need more sophisticated implementation)
	return true
}

func (tis *TaskIntegrationService) processTemplate(template string, data map[string]interface{}) string {
	result := template

	// Simple template replacement
	for key, value := range data {
		placeholder := "{{" + key + "}}"
		result = strings.ReplaceAll(result, placeholder, fmt.Sprintf("%v", value))
	}

	return result
}

func (tis *TaskIntegrationService) getTaskTemplate(templateID uuid.UUID) (*TaskTemplate, error) {
	var template TaskTemplate
	query := `
		SELECT id, name, description, category, priority, estimated_duration, tags,
		       checklist, subtasks, metadata, created_by, is_public, usage_count, created_at
		FROM task_templates WHERE id = $1
	`

	err := tis.db.QueryRow(query, templateID).Scan(
		&template.ID, &template.Name, &template.Description, &template.Category,
		&template.Priority, &template.EstimatedDuration, &template.Tags,
		&template.Checklist, &template.Subtasks, &template.Metadata,
		&template.CreatedBy, &template.IsPublic, &template.UsageCount, &template.CreatedAt,
	)

	return &template, err
}

func (tis *TaskIntegrationService) incrementTemplateUsage(templateID uuid.UUID) {
	query := `UPDATE task_templates SET usage_count = usage_count + 1 WHERE id = $1`
	tis.db.Exec(query, templateID)
}

func (tis *TaskIntegrationService) createSubtask(subtaskData map[string]interface{}) error {
	parentID, ok := subtaskData["parent_task_id"].(uuid.UUID)
	if !ok {
		if s, ok2 := subtaskData["parent_task_id"].(string); ok2 {
			parsed, err := uuid.Parse(s)
			if err != nil {
				return err
			}
			parentID = parsed
		} else {
			return fmt.Errorf("subtask missing parent_task_id")
		}
	}

	title, _ := subtaskData["title"].(string)
	if title == "" {
		return fmt.Errorf("subtask missing title")
	}

	id := uuid.New()
	now := time.Now()
	query := `
		INSERT INTO task_subtasks (id, task_id, title, completed, "order", created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	completed := false
	if c, ok := subtaskData["completed"].(bool); ok {
		completed = c
	}
	order := 0
	if o, ok := subtaskData["order"].(int); ok {
		order = o
	}

	_, err := tis.db.Exec(query, id, parentID, title, completed, order, now)
	return err
}
