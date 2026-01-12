package services

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFormWorkflowCreation tests creating form workflows
func TestFormWorkflowCreation(t *testing.T) {
	workflowService := &FormWorkflowService{}

	formID := uuid.New()
	creatorID := uuid.New()

	// Test creating a simple approval workflow
	workflow := &FormWorkflow{
		FormID:    formID,
		CreatorID: creatorID,
		Name:      "Simple Approval",
		Description: "Basic approval workflow",
		Steps: []WorkflowStep{
			{
				Type:   "approval",
				Name:   "Manager Approval",
				Config: map[string]interface{}{
					"approver": "manager@example.com",
					"condition": map[string]interface{}{
						"field": "amount",
						"op":    "gt",
						"value": 1000,
					},
				},
			},
			{
				Type:   "notification",
				Name:   "Approval Notification",
				Config: map[string]interface{}{
					"to":      "submitter",
					"message": "Your request has been approved",
				},
			},
		},
	}

	createdWorkflow, err := workflowService.CreateWorkflow(workflow)
	assert.NoError(t, err)
	assert.NotNil(t, createdWorkflow)
	assert.Equal(t, formID, createdWorkflow.FormID)
	assert.Equal(t, creatorID, createdWorkflow.CreatorID)
	assert.Len(t, createdWorkflow.Steps, 2)
	assert.Equal(t, "draft", createdWorkflow.Status)
}

// TestFormWorkflowExecution tests executing workflows
func TestFormWorkflowExecution(t *testing.T) {
	workflowService := &FormWorkflowService{}

	// Create a workflow
	workflow := &FormWorkflow{
		FormID:    uuid.New(),
		CreatorID: uuid.New(),
		Name:      "Approval Workflow",
		Steps: []WorkflowStep{
			{
				Type: "approval",
				Name: "Manager Review",
				Config: map[string]interface{}{
					"approver": "manager@example.com",
					"condition": map[string]interface{}{
						"field": "amount",
						"op":    "gt",
						"value": 1000,
					},
				},
			},
			{
				Type: "notification",
				Name: "Success Notification",
				Config: map[string]interface{}{
					"to":      "submitter",
					"message": "Request approved",
				},
			},
		},
	}

	createdWorkflow, err := workflowService.CreateWorkflow(workflow)
	require.NoError(t, err)

	// Test form submission data
	submission := &FormSubmission{
		ID:     uuid.New(),
		FormID: workflow.FormID,
		Data: map[string]interface{}{
			"name":   "John Doe",
			"amount": 1500,
			"email":  "john@example.com",
		},
		SubmitterID: uuid.New(),
	}

	// Execute workflow
	result, err := workflowService.ExecuteWorkflow(createdWorkflow.ID, submission)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "in_progress", result.Status)

	// Check that approval step was triggered
	assert.Len(t, result.CurrentSteps, 1)
	assert.Equal(t, "approval", result.CurrentSteps[0].Type)
}

// TestFormWorkflowConditionEvaluation tests condition evaluation
func TestFormWorkflowConditionEvaluation(t *testing.T) {
	workflowService := &FormWorkflowService{}

	tests := []struct {
		name      string
		condition map[string]interface{}
		data      map[string]interface{}
		expected  bool
	}{
		{
			name: "greater than - true",
			condition: map[string]interface{}{
				"field": "amount",
				"op":    "gt",
				"value": 1000,
			},
			data: map[string]interface{}{
				"amount": 1500,
			},
			expected: true,
		},
		{
			name: "greater than - false",
			condition: map[string]interface{}{
				"field": "amount",
				"op":    "gt",
				"value": 1000,
			},
			data: map[string]interface{}{
				"amount": 500,
			},
			expected: false,
		},
		{
			name: "equals - true",
			condition: map[string]interface{}{
				"field": "status",
				"op":    "eq",
				"value": "urgent",
			},
			data: map[string]interface{}{
				"status": "urgent",
			},
			expected: true,
		},
		{
			name: "equals - false",
			condition: map[string]interface{}{
				"field": "status",
				"op":    "eq",
				"value": "urgent",
			},
			data: map[string]interface{}{
				"status": "normal",
			},
			expected: false,
		},
		{
			name: "contains - true",
			condition: map[string]interface{}{
				"field": "tags",
				"op":    "contains",
				"value": "urgent",
			},
			data: map[string]interface{}{
				"tags": []string{"normal", "urgent", "important"},
			},
			expected: true,
		},
		{
			name: "contains - false",
			condition: map[string]interface{}{
				"field": "tags",
				"op":    "contains",
				"value": "critical",
			},
			data: map[string]interface{}{
				"tags": []string{"normal", "urgent"},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := workflowService.EvaluateCondition(tt.condition, tt.data)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestFormWorkflowStepExecution tests individual step execution
func TestFormWorkflowStepExecution(t *testing.T) {
	workflowService := &FormWorkflowService{}

	submission := &FormSubmission{
		ID:     uuid.New(),
		FormID: uuid.New(),
		Data: map[string]interface{}{
			"name":   "Jane Smith",
			"email":  "jane@example.com",
			"amount": 2500,
		},
		SubmitterID: uuid.New(),
	}

	// Test approval step
	approvalStep := WorkflowStep{
		Type: "approval",
		Name: "Manager Approval",
		Config: map[string]interface{}{
			"approver": "manager@example.com",
		},
	}

	approvalResult, err := workflowService.ExecuteStep(approvalStep, submission, nil)
	assert.NoError(t, err)
	assert.NotNil(t, approvalResult)
	assert.Equal(t, "pending", approvalResult.Status)
	assert.Equal(t, "manager@example.com", approvalResult.AssignedTo)

	// Test notification step
	notificationStep := WorkflowStep{
		Type: "notification",
		Name: "Success Email",
		Config: map[string]interface{}{
			"to":      "submitter",
			"message": "Your request has been processed",
		},
	}

	notificationResult, err := workflowService.ExecuteStep(notificationStep, submission, nil)
	assert.NoError(t, err)
	assert.NotNil(t, notificationResult)
	assert.Equal(t, "completed", notificationResult.Status)
}

// TestFormWorkflowParallelExecution tests parallel workflow paths
func TestFormWorkflowParallelExecution(t *testing.T) {
	workflowService := &FormWorkflowService{}

	// Create workflow with parallel approval steps
	workflow := &FormWorkflow{
		FormID:    uuid.New(),
		CreatorID: uuid.New(),
		Name:      "Parallel Approval",
		Steps: []WorkflowStep{
			{
				Type: "parallel",
				Name: "Parallel Approvals",
				Config: map[string]interface{}{
					"steps": []map[string]interface{}{
						{
							"type":   "approval",
							"name":   "Manager Approval",
							"config": map[string]interface{}{
								"approver": "manager@example.com",
							},
						},
						{
							"type":   "approval",
							"name":   "Finance Approval",
							"config": map[string]interface{}{
								"approver": "finance@example.com",
							},
						},
					},
				},
			},
			{
				Type: "notification",
				Name: "Final Notification",
				Config: map[string]interface{}{
					"to":      "submitter",
					"message": "All approvals received",
				},
			},
		},
	}

	createdWorkflow, err := workflowService.CreateWorkflow(workflow)
	require.NoError(t, err)

	submission := &FormSubmission{
		ID:     uuid.New(),
		FormID: workflow.FormID,
		Data: map[string]interface{}{
			"name":   "Alice Johnson",
			"amount": 5000,
		},
		SubmitterID: uuid.New(),
	}

	// Execute workflow
	result, err := workflowService.ExecuteWorkflow(createdWorkflow.ID, submission)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Should have two parallel approval steps
	assert.Len(t, result.CurrentSteps, 2)
	for _, step := range result.CurrentSteps {
		assert.Equal(t, "approval", step.Type)
	}
}

// TestFormWorkflowConditionalBranching tests conditional workflow branching
func TestFormWorkflowConditionalBranching(t *testing.T) {
	workflowService := &FormWorkflowService{}

	// Create workflow with conditional branching
	workflow := &FormWorkflow{
		FormID:    uuid.New(),
		CreatorID: uuid.New(),
		Name:      "Conditional Workflow",
		Steps: []WorkflowStep{
			{
				Type: "condition",
				Name: "Amount Check",
				Config: map[string]interface{}{
					"condition": map[string]interface{}{
						"field": "amount",
						"op":    "gt",
						"value": 5000,
					},
					"if_true": []map[string]interface{}{
						{
							"type":   "approval",
							"name":   "Senior Manager Approval",
							"config": map[string]interface{}{
								"approver": "senior@example.com",
							},
						},
					},
					"if_false": []map[string]interface{}{
						{
							"type":   "approval",
							"name":   "Manager Approval",
							"config": map[string]interface{}{
								"approver": "manager@example.com",
							},
						},
					},
				},
			},
		},
	}

	createdWorkflow, err := workflowService.CreateWorkflow(workflow)
	require.NoError(t, err)

	// Test high amount path
	highAmountSubmission := &FormSubmission{
		ID:     uuid.New(),
		FormID: workflow.FormID,
		Data: map[string]interface{}{
			"name":   "Bob Wilson",
			"amount": 7500,
		},
		SubmitterID: uuid.New(),
	}

	result, err := workflowService.ExecuteWorkflow(createdWorkflow.ID, highAmountSubmission)
	assert.NoError(t, err)
	assert.Len(t, result.CurrentSteps, 1)
	assert.Equal(t, "Senior Manager Approval", result.CurrentSteps[0].Name)

	// Test low amount path
	lowAmountSubmission := &FormSubmission{
		ID:     uuid.New(),
		FormID: workflow.FormID,
		Data: map[string]interface{}{
			"name":   "Charlie Brown",
			"amount": 1500,
		},
		SubmitterID: uuid.New(),
	}

	result, err = workflowService.ExecuteWorkflow(createdWorkflow.ID, lowAmountSubmission)
	assert.NoError(t, err)
	assert.Len(t, result.CurrentSteps, 1)
	assert.Equal(t, "Manager Approval", result.CurrentSteps[0].Name)
}

// TestFormWorkflowDeadlineHandling tests deadline and escalation
func TestFormWorkflowDeadlineHandling(t *testing.T) {
	workflowService := &FormWorkflowService{}

	workflow := &FormWorkflow{
		FormID:    uuid.New(),
		CreatorID: uuid.New(),
		Name:      "Deadline Workflow",
		Steps: []WorkflowStep{
			{
				Type: "approval",
				Name: "Urgent Approval",
				Config: map[string]interface{}{
					"approver":     "approver@example.com",
					"deadline":     "24h",
					"escalation":   "manager@example.com",
					"reminder":     "6h",
				},
			},
		},
	}

	createdWorkflow, err := workflowService.CreateWorkflow(workflow)
	require.NoError(t, err)

	submission := &FormSubmission{
		ID:     uuid.New(),
		FormID: workflow.FormID,
		Data: map[string]interface{}{
			"name":   "Emergency Request",
			"amount": 10000,
		},
		SubmitterID: uuid.New(),
	}

	result, err := workflowService.ExecuteWorkflow(createdWorkflow.ID, submission)
	assert.NoError(t, err)
	assert.NotNil(t, result.Deadline)
	assert.NotNil(t, result.EscalationTime)

	// Test deadline checking
	isOverdue, err := workflowService.CheckDeadline(result.ID)
	assert.NoError(t, err)
	assert.False(t, isOverdue) // Not overdue yet

	// Test escalation trigger
	escalated, err := workflowService.TriggerEscalation(result.ID)
	assert.NoError(t, err)
	assert.False(t, escalated) // Not escalated yet
}

// TestFormWorkflowAuditTrail tests workflow audit logging
func TestFormWorkflowAuditTrail(t *testing.T) {
	workflowService := &FormWorkflowService{}

	workflow := &FormWorkflow{
		FormID:    uuid.New(),
		CreatorID: uuid.New(),
		Name:      "Audit Workflow",
		Steps: []WorkflowStep{
			{
				Type: "approval",
				Name: "Simple Approval",
				Config: map[string]interface{}{
					"approver": "approver@example.com",
				},
			},
		},
	}

	createdWorkflow, err := workflowService.CreateWorkflow(workflow)
	require.NoError(t, err)

	submission := &FormSubmission{
		ID:     uuid.New(),
		FormID: workflow.FormID,
		Data: map[string]interface{}{
			"name": "Audit Test",
		},
		SubmitterID: uuid.New(),
	}

	// Execute workflow
	_, err = workflowService.ExecuteWorkflow(createdWorkflow.ID, submission)
	require.NoError(t, err)

	// Check audit trail
	auditLogs, err := workflowService.GetWorkflowAuditTrail(createdWorkflow.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, auditLogs)

	// Should have at least workflow started and step created events
	assert.GreaterOrEqual(t, len(auditLogs), 2)
	assert.Contains(t, []string{"workflow_started", "step_created"}, auditLogs[0].Action)
}

// TestFormWorkflowErrorHandling tests error conditions
func TestFormWorkflowErrorHandling(t *testing.T) {
	workflowService := &FormWorkflowService{}

	// Test executing non-existent workflow
	_, err := workflowService.ExecuteWorkflow(uuid.New(), &FormSubmission{})
	assert.Error(t, err)

	// Test creating workflow with invalid steps
	invalidWorkflow := &FormWorkflow{
		FormID:    uuid.New(),
		CreatorID: uuid.New(),
		Name:      "Invalid Workflow",
		Steps: []WorkflowStep{
			{
				Type:   "invalid_type",
				Name:   "Invalid Step",
				Config: map[string]interface{}{},
			},
		},
	}

	_, err = workflowService.CreateWorkflow(invalidWorkflow)
	assert.Error(t, err)

	// Test executing workflow with missing submission data
	validWorkflow := &FormWorkflow{
		FormID:    uuid.New(),
		CreatorID: uuid.New(),
		Name:      "Valid Workflow",
		Steps: []WorkflowStep{
			{
				Type: "condition",
				Name: "Check Field",
				Config: map[string]interface{}{
					"condition": map[string]interface{}{
						"field": "missing_field",
						"op":    "eq",
						"value": "test",
					},
				},
			},
		},
	}

	createdWorkflow, err := workflowService.CreateWorkflow(validWorkflow)
	require.NoError(t, err)

	submission := &FormSubmission{
		ID:     uuid.New(),
		FormID: validWorkflow.FormID,
		Data:   map[string]interface{}{}, // Missing required field
	}

	_, err = workflowService.ExecuteWorkflow(createdWorkflow.ID, submission)
	assert.Error(t, err) // Should fail due to missing field
}

// TestFormWorkflowBulkOperations tests bulk workflow operations
func TestFormWorkflowBulkOperations(t *testing.T) {
	workflowService := &FormWorkflowService{}

	formID := uuid.New()
	creatorID := uuid.New()

	// Create multiple workflows
	workflows := make([]*FormWorkflow, 0)
	workflowIDs := make([]uuid.UUID, 0)

	for i := 0; i < 3; i++ {
		workflow := &FormWorkflow{
			FormID:    formID,
			CreatorID: creatorID,
			Name:      fmt.Sprintf("Bulk Workflow %d", i),
			Steps: []WorkflowStep{
				{
					Type:   "notification",
					Name:   "Test Notification",
					Config: map[string]interface{}{
						"to":      "test@example.com",
						"message": "Bulk test",
					},
				},
			},
		}

		createdWorkflow, err := workflowService.CreateWorkflow(workflow)
		require.NoError(t, err)
		workflows = append(workflows, createdWorkflow)
		workflowIDs = append(workflowIDs, createdWorkflow.ID)
	}

	// Test bulk status update
	err := workflowService.BulkUpdateWorkflowStatus(workflowIDs, "active")
	assert.NoError(t, err)

	// Verify status updates
	for _, workflowID := range workflowIDs {
		workflow, err := workflowService.GetWorkflow(workflowID)
		require.NoError(t, err)
		assert.Equal(t, "active", workflow.Status)
	}

	// Test bulk deletion
	err = workflowService.BulkDeleteWorkflows(workflowIDs[:2])
	assert.NoError(t, err)

	// Verify deletions
	for _, workflowID := range workflowIDs[:2] {
		_, err := workflowService.GetWorkflow(workflowID)
		assert.Error(t, err) // Should not exist
	}

	// Last workflow should still exist
	_, err = workflowService.GetWorkflow(workflowIDs[2])
	assert.NoError(t, err)
}
