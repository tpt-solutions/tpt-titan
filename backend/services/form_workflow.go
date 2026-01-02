package services

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// FormWorkflowService manages workflow automation for forms
type FormWorkflowService struct {
	db              *sql.DB
	emailService    *EmailService
	relationshipSvc *FormRelationshipService
}

// WorkflowDefinition represents a workflow definition
type WorkflowDefinition struct {
	ID          uuid.UUID              `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	FormID      uuid.UUID              `json:"form_id"`
	Trigger     string                 `json:"trigger"`     // "on_submit", "on_update", "on_approve", "scheduled"
	Steps       []WorkflowStep         `json:"steps"`
	IsActive    bool                   `json:"is_active"`
	CreatedBy   uuid.UUID              `json:"created_by"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// WorkflowStep represents a step in a workflow
type WorkflowStep struct {
	ID          uuid.UUID              `json:"id"`
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`        // "approval", "notification", "assignment", "condition", "action"
	Config      map[string]interface{} `json:"config"`
	Order       int                    `json:"order"`
	NextStepID  *uuid.UUID             `json:"next_step_id,omitempty"`
	AltStepID   *uuid.UUID             `json:"alt_step_id,omitempty"` // Alternative path (e.g., rejection)
}

// WorkflowInstance represents a running workflow instance
type WorkflowInstance struct {
	ID          uuid.UUID            `json:"id"`
	WorkflowID  uuid.UUID            `json:"workflow_id"`
	RecordID    uuid.UUID            `json:"record_id"`    // Form response ID
	CurrentStep *uuid.UUID           `json:"current_step"`
	Status      string               `json:"status"`       // "running", "completed", "failed", "waiting"
	Context     map[string]interface{} `json:"context"`      // Workflow execution context
	StartedAt   time.Time            `json:"started_at"`
	CompletedAt *time.Time           `json:"completed_at,omitempty"`
}

// ApprovalRequest represents an approval request in a workflow
type ApprovalRequest struct {
	ID         uuid.UUID `json:"id"`
	InstanceID uuid.UUID `json:"instance_id"`
	StepID     uuid.UUID `json:"step_id"`
	RequestedBy uuid.UUID `json:"requested_by"`
	AssignedTo uuid.UUID `json:"assigned_to"`
	Status     string    `json:"status"` // "pending", "approved", "rejected"
	Message    string    `json:"message,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	RespondedAt *time.Time `json:"responded_at,omitempty"`
}

// NotificationTemplate represents a template for workflow notifications
type NotificationTemplate struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Type    string    `json:"type"`    // "email", "sms", "in_app"
	Subject string    `json:"subject"`
	Body    string    `json:"body"`
	Variables []string `json:"variables"` // Available template variables
}

// NewFormWorkflowService creates a new workflow service
func NewFormWorkflowService(db *sql.DB, emailSvc *EmailService, relSvc *FormRelationshipService) *FormWorkflowService {
	return &FormWorkflowService{
		db:              db,
		emailService:    emailSvc,
		relationshipSvc: relSvc,
	}
}

// CreateWorkflow creates a new workflow definition
func (fws *FormWorkflowService) CreateWorkflow(workflow *WorkflowDefinition) error {
	workflow.ID = uuid.New()
	workflow.CreatedAt = time.Now()
	workflow.UpdatedAt = time.Now()

	query := `
		INSERT INTO form_workflows (id, name, description, form_id, trigger, steps,
			is_active, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := fws.db.Exec(query,
		workflow.ID, workflow.Name, workflow.Description, workflow.FormID,
		workflow.Trigger, workflow.Steps, workflow.IsActive, workflow.CreatedBy,
		workflow.CreatedAt, workflow.UpdatedAt,
	)

	return err
}

// StartWorkflow starts a workflow instance for a form submission
func (fws *FormWorkflowService) StartWorkflow(workflowID uuid.UUID, recordID uuid.UUID, userID uuid.UUID) error {
	// Get workflow definition
	workflow, err := fws.GetWorkflow(workflowID)
	if err != nil {
		return err
	}

	if !workflow.IsActive {
		return fmt.Errorf("workflow is not active")
	}

	// Create workflow instance
	instance := &WorkflowInstance{
		ID:         uuid.New(),
		WorkflowID: workflowID,
		RecordID:   recordID,
		Status:     "running",
		Context:    make(map[string]interface{}),
		StartedAt:  time.Now(),
	}

	// Set initial context
	instance.Context["initiator"] = userID
	instance.Context["started_at"] = instance.StartedAt

	// Find first step
	if len(workflow.Steps) > 0 {
		firstStep := workflow.Steps[0]
		instance.CurrentStep = &firstStep.ID

		// Execute first step
		err = fws.executeWorkflowStep(instance, firstStep, userID)
		if err != nil {
			instance.Status = "failed"
		}
	} else {
		instance.Status = "completed"
		instance.CompletedAt = &instance.StartedAt
	}

	// Save instance
	return fws.saveWorkflowInstance(instance)
}

// ExecuteWorkflowStep executes a workflow step
func (fws *FormWorkflowService) executeWorkflowStep(instance *WorkflowInstance, step WorkflowStep, userID uuid.UUID) error {
	switch step.Type {
	case "approval":
		return fws.executeApprovalStep(instance, step, userID)
	case "notification":
		return fws.executeNotificationStep(instance, step, userID)
	case "assignment":
		return fws.executeAssignmentStep(instance, step, userID)
	case "condition":
		return fws.executeConditionStep(instance, step, userID)
	case "action":
		return fws.executeActionStep(instance, step, userID)
	default:
		return fmt.Errorf("unknown step type: %s", step.Type)
	}
}

// Execute approval step
func (fws *FormWorkflowService) executeApprovalStep(instance *WorkflowInstance, step WorkflowStep, userID uuid.UUID) error {
	// Create approval request
	assignedTo, ok := step.Config["assigned_to"].(string)
	if !ok {
		return fmt.Errorf("approval step missing assigned_to")
	}

	assignedToUUID, err := uuid.Parse(assignedTo)
	if err != nil {
		return fmt.Errorf("invalid assigned_to UUID: %v", err)
	}

	approval := &ApprovalRequest{
		ID:         uuid.New(),
		InstanceID: instance.ID,
		StepID:     step.ID,
		RequestedBy: userID,
		AssignedTo: assignedToUUID,
		Status:     "pending",
		CreatedAt:  time.Now(),
	}

	// Save approval request
	query := `
		INSERT INTO workflow_approvals (id, instance_id, step_id, requested_by,
			assigned_to, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err = fws.db.Exec(query,
		approval.ID, approval.InstanceID, approval.StepID, approval.RequestedBy,
		approval.AssignedTo, approval.Status, approval.CreatedAt,
	)

	if err != nil {
		return err
	}

	// Send notification to approver
	return fws.notifyApprover(approval, step)
}

// Execute notification step
func (fws *FormWorkflowService) executeNotificationStep(instance *WorkflowInstance, step WorkflowStep, userID uuid.UUID) error {
	notificationType, ok := step.Config["type"].(string)
	if !ok {
		return fmt.Errorf("notification step missing type")
	}

	recipients, ok := step.Config["recipients"].([]interface{})
	if !ok {
		return fmt.Errorf("notification step missing recipients")
	}

	templateID, ok := step.Config["template_id"].(string)
	if !ok {
		return fmt.Errorf("notification step missing template_id")
	}

	templateUUID, err := uuid.Parse(templateID)
	if err != nil {
		return fmt.Errorf("invalid template_id: %v", err)
	}

	// Get notification template
	template, err := fws.getNotificationTemplate(templateUUID)
	if err != nil {
		return err
	}

	// Get form data for template variables
	formData, err := fws.getFormData(instance.RecordID)
	if err != nil {
		return err
	}

	// Process template
	subject := fws.processTemplate(template.Subject, formData)
	body := fws.processTemplate(template.Body, formData)

	// Send notifications
	for _, recipient := range recipients {
		recipientID, err := uuid.Parse(recipient.(string))
		if err != nil {
			continue
		}

		switch notificationType {
		case "email":
			err = fws.sendEmailNotification(recipientID, subject, body)
		case "in_app":
			err = fws.sendInAppNotification(recipientID, subject, body)
		}

		if err != nil {
			// Log error but continue with other recipients
			fmt.Printf("Failed to send notification to %s: %v\n", recipientID, err)
		}
	}

	return nil
}

// Execute assignment step
func (fws *FormWorkflowService) executeAssignmentStep(instance *WorkflowInstance, step WorkflowStep, userID uuid.UUID) error {
	// Assign the record to a user or group
	assignee, ok := step.Config["assignee"].(string)
	if !ok {
		return fmt.Errorf("assignment step missing assignee")
	}

	// Update record assignment in database
	// This would depend on your form response structure

	return nil
}

// Execute condition step
func (fws *FormWorkflowService) executeConditionStep(instance *WorkflowInstance, step WorkflowStep, userID uuid.UUID) error {
	// Evaluate condition and choose next step
	condition, ok := step.Config["condition"].(string)
	if !ok {
		return fmt.Errorf("condition step missing condition")
	}

	// Evaluate condition (simplified - would need full expression evaluation)
	result := fws.evaluateCondition(condition, instance.Context)

	if result {
		instance.CurrentStep = step.NextStepID
	} else {
		instance.CurrentStep = step.AltStepID
	}

	return nil
}

// Execute action step
func (fws *FormWorkflowService) executeActionStep(instance *WorkflowInstance, step WorkflowStep, userID uuid.UUID) error {
	// Execute custom action
	actionType, ok := step.Config["action_type"].(string)
	if !ok {
		return fmt.Errorf("action step missing action_type")
	}

	switch actionType {
	case "update_field":
		return fws.executeFieldUpdateAction(step, instance.RecordID)
	case "create_record":
		return fws.executeRecordCreationAction(step, instance)
	case "webhook":
		return fws.executeWebhookAction(step, instance)
	default:
		return fmt.Errorf("unknown action type: %s", actionType)
	}
}

// Process approval response
func (fws *FormWorkflowService) ProcessApproval(approvalID uuid.UUID, approved bool, message string, userID uuid.UUID) error {
	// Get approval request
	var approval ApprovalRequest
	query := `SELECT instance_id, step_id FROM workflow_approvals WHERE id = $1`
	err := fws.db.QueryRow(query, approvalID).Scan(&approval.InstanceID, &approval.StepID)
	if err != nil {
		return err
	}

	// Update approval status
	status := "rejected"
	if approved {
		status = "approved"
	}

	updateQuery := `
		UPDATE workflow_approvals
		SET status = $1, message = $2, responded_at = $3
		WHERE id = $4
	`

	_, err = fws.db.Exec(updateQuery, status, message, time.Now(), approvalID)
	if err != nil {
		return err
	}

	// Update workflow instance
	instance, err := fws.getWorkflowInstance(approval.InstanceID)
	if err != nil {
		return err
	}

	// Add to context
	instance.Context["approval_"+approval.StepID.String()] = status
	instance.Context["approval_message_"+approval.StepID.String()] = message

	// Move to next step based on approval result
	workflow, err := fws.GetWorkflow(instance.WorkflowID)
	if err != nil {
		return err
	}

	// Find current step and determine next step
	for _, step := range workflow.Steps {
		if step.ID == approval.StepID {
			if approved && step.NextStepID != nil {
				instance.CurrentStep = step.NextStepID
			} else if !approved && step.AltStepID != nil {
				instance.CurrentStep = step.AltStepID
			} else {
				// End of workflow
				instance.Status = "completed"
				now := time.Now()
				instance.CompletedAt = &now
			}
			break
		}
	}

	// Execute next step if workflow continues
	if instance.Status == "running" && instance.CurrentStep != nil {
		for _, step := range workflow.Steps {
			if step.ID == *instance.CurrentStep {
				err = fws.executeWorkflowStep(instance, step, userID)
				if err != nil {
					instance.Status = "failed"
				}
				break
			}
		}
	}

	return fws.saveWorkflowInstance(instance)
}

// Get pending approvals for a user
func (fws *FormWorkflowService) GetPendingApprovals(userID uuid.UUID) ([]map[string]interface{}, error) {
	query := `
		SELECT wa.id, wa.instance_id, wa.step_id, wa.requested_by, wa.created_at,
		       fw.name as workflow_name, fr.form_id
		FROM workflow_approvals wa
		JOIN workflow_instances wi ON wa.instance_id = wi.id
		JOIN form_workflows fw ON wi.workflow_id = fw.id
		JOIN form_responses fr ON wi.record_id = fr.id
		WHERE wa.assigned_to = $1 AND wa.status = 'pending'
		ORDER BY wa.created_at DESC
	`

	rows, err := fws.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var approvals []map[string]interface{}
	for rows.Next() {
		var approval struct {
			ID          uuid.UUID `json:"id"`
			InstanceID  uuid.UUID `json:"instance_id"`
			StepID      uuid.UUID `json:"step_id"`
			RequestedBy uuid.UUID `json:"requested_by"`
			CreatedAt   time.Time `json:"created_at"`
			WorkflowName string   `json:"workflow_name"`
			FormID      uuid.UUID `json:"form_id"`
		}

		err := rows.Scan(
			&approval.ID, &approval.InstanceID, &approval.StepID, &approval.RequestedBy,
			&approval.CreatedAt, &approval.WorkflowName, &approval.FormID,
		)
		if err != nil {
			continue
		}

		approvals = append(approvals, map[string]interface{}{
			"id":             approval.ID,
			"instance_id":    approval.InstanceID,
			"step_id":        approval.StepID,
			"requested_by":   approval.RequestedBy,
			"created_at":     approval.CreatedAt,
			"workflow_name":  approval.WorkflowName,
			"form_id":        approval.FormID,
		})
	}

	return approvals, nil
}

// Create notification template
func (fws *FormWorkflowService) CreateNotificationTemplate(template *NotificationTemplate) error {
	template.ID = uuid.New()

	query := `
		INSERT INTO notification_templates (id, name, type, subject, body, variables)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := fws.db.Exec(query,
		template.ID, template.Name, template.Type, template.Subject,
		template.Body, template.Variables,
	)

	return err
}

// Helper methods

func (fws *FormWorkflowService) GetWorkflow(workflowID uuid.UUID) (*WorkflowDefinition, error) {
	var workflow WorkflowDefinition
	query := `
		SELECT id, name, description, form_id, trigger, steps, is_active,
		       created_by, created_at, updated_at
		FROM form_workflows WHERE id = $1
	`

	err := fws.db.QueryRow(query, workflowID).Scan(
		&workflow.ID, &workflow.Name, &workflow.Description, &workflow.FormID,
		&workflow.Trigger, &workflow.Steps, &workflow.IsActive, &workflow.CreatedBy,
		&workflow.CreatedAt, &workflow.UpdatedAt,
	)

	return &workflow, err
}

func (fws *FormWorkflowService) saveWorkflowInstance(instance *WorkflowInstance) error {
	query := `
		INSERT INTO workflow_instances (id, workflow_id, record_id, current_step,
			status, context, started_at, completed_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (id) DO UPDATE SET
			current_step = EXCLUDED.current_step,
			status = EXCLUDED.status,
			context = EXCLUDED.context,
			completed_at = EXCLUDED.completed_at
	`

	_, err := fws.db.Exec(query,
		instance.ID, instance.WorkflowID, instance.RecordID, instance.CurrentStep,
		instance.Status, instance.Context, instance.StartedAt, instance.CompletedAt,
	)

	return err
}

func (fws *FormWorkflowService) getWorkflowInstance(instanceID uuid.UUID) (*WorkflowInstance, error) {
	var instance WorkflowInstance
	query := `
		SELECT id, workflow_id, record_id, current_step, status, context,
		       started_at, completed_at
		FROM workflow_instances WHERE id = $1
	`

	err := fws.db.QueryRow(query, instanceID).Scan(
		&instance.ID, &instance.WorkflowID, &instance.RecordID, &instance.CurrentStep,
		&instance.Status, &instance.Context, &instance.StartedAt, &instance.CompletedAt,
	)

	return &instance, err
}

func (fws *FormWorkflowService) notifyApprover(approval *ApprovalRequest, step WorkflowStep) error {
	// Send notification to approver (would integrate with notification system)
	return nil
}

func (fws *FormWorkflowService) getNotificationTemplate(templateID uuid.UUID) (*NotificationTemplate, error) {
	var template NotificationTemplate
	query := `SELECT id, name, type, subject, body, variables FROM notification_templates WHERE id = $1`

	err := fws.db.QueryRow(query, templateID).Scan(
		&template.ID, &template.Name, &template.Type, &template.Subject,
		&template.Body, &template.Variables,
	)

	return &template, err
}

func (fws *FormWorkflowService) getFormData(recordID uuid.UUID) (map[string]interface{}, error) {
	// Get form response data
	var data map[string]interface{}
	query := `SELECT response_data FROM form_responses WHERE id = $1`

	err := fws.db.QueryRow(query, recordID).Scan(&data)
	return data, err
}

func (fws *FormWorkflowService) processTemplate(template string, data map[string]interface{}) string {
	// Simple template processing - replace {{variable}} with values
	result := template

	for key, value := range data {
		placeholder := fmt.Sprintf("{{%s}}", key)
		result = strings.ReplaceAll(result, placeholder, fmt.Sprintf("%v", value))
	}

	return result
}

func (fws *FormWorkflowService) sendEmailNotification(recipientID uuid.UUID, subject, body string) error {
	// Get user email
	var email string
	query := `SELECT email FROM users WHERE id = $1`
	err := fws.db.QueryRow(query, recipientID).Scan(&email)
	if err != nil {
		return err
	}

	// Send email using email service
	return fws.emailService.SendEmail([]string{email}, subject, body, "")
}

func (fws *FormWorkflowService) sendInAppNotification(recipientID uuid.UUID, title, message string) error {
	// Create in-app notification
	notification := map[string]interface{}{
		"id":      uuid.New(),
		"user_id": recipientID,
		"title":   title,
		"message": message,
		"type":    "workflow",
		"read":    false,
		"created_at": time.Now(),
	}

	// Save to notifications table
	query := `
		INSERT INTO notifications (id, user_id, title, message, type, read, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := fws.db.Exec(query,
		notification["id"], notification["user_id"], notification["title"],
		notification["message"], notification["type"], notification["read"],
		notification["created_at"],
	)

	return err
}

func (fws *FormWorkflowService) evaluateCondition(condition string, context map[string]interface{}) bool {
	// Simple condition evaluation (would need full expression parser)
	// For now, just check if context contains certain values
	return true // Simplified
}

func (fws *FormWorkflowService) executeFieldUpdateAction(step WorkflowStep, recordID uuid.UUID) error {
	// Update a field in the form response
	field, ok := step.Config["field"].(string)
	if !ok {
		return fmt.Errorf("field update action missing field")
	}

	value, ok := step.Config["value"]
	if !ok {
		return fmt.Errorf("field update action missing value")
	}

	// Update the field (would need to handle JSON field updates in form_responses)
	return nil
}

func (fws *FormWorkflowService) executeRecordCreationAction(step WorkflowStep, instance *WorkflowInstance) error {
	// Create a new record in another form
	targetForm, ok := step.Config["target_form"].(string)
	if !ok {
		return fmt.Errorf("record creation action missing target_form")
	}

	// Create new record logic
	return nil
}

func (fws *FormWorkflowService) executeWebhookAction(step WorkflowStep, instance *WorkflowInstance) error {
	// Execute webhook
	webhookURL, ok := step.Config["url"].(string)
	if !ok {
		return fmt.Errorf("webhook action missing url")
	}

	// Make HTTP request to webhook URL
	fmt.Printf("Executing workflow webhook to %s\n", webhookURL)
	return nil
}
