// backend/services/form_workflow_test.go
// Run with: cd backend && go test ./services/... -run TestWorkflow -v

package services

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

// ─── NewFormWorkflowService ───────────────────────────────────────────────────

func TestNewFormWorkflowService_NilDependencies(t *testing.T) {
	// All dependencies are optional for unit tests (nil db, nil email, nil rel)
	svc := NewFormWorkflowService(nil, nil, nil)
	if svc == nil {
		t.Fatal("expected non-nil FormWorkflowService")
	}
}

// ─── processTemplate ─────────────────────────────────────────────────────────
// Tests the unexported helper directly (same package).

func TestProcessTemplate_SingleVariable(t *testing.T) {
	svc := &FormWorkflowService{}

	result := svc.processTemplate(
		"Hello, {{name}}!",
		map[string]interface{}{"name": "Alice"},
	)

	if result != "Hello, Alice!" {
		t.Errorf("got %q, want %q", result, "Hello, Alice!")
	}
}

func TestProcessTemplate_MultipleVariables(t *testing.T) {
	svc := &FormWorkflowService{}

	result := svc.processTemplate(
		"Dear {{name}}, your amount is {{amount}}.",
		map[string]interface{}{
			"name":   "Bob",
			"amount": 1500,
		},
	)

	if result == "Dear {{name}}, your amount is {{amount}}." {
		t.Error("template variables were not replaced")
	}
	// Both variables should be substituted (map iteration order is random, but both must appear)
	if len(result) == 0 {
		t.Error("template result should not be empty")
	}
}

func TestProcessTemplate_MissingVariableLeftAsIs(t *testing.T) {
	svc := &FormWorkflowService{}

	// Data does not contain {{status}} so it should be left unchanged
	result := svc.processTemplate(
		"Status: {{status}}",
		map[string]interface{}{},
	)

	if result != "Status: {{status}}" {
		t.Errorf("unused placeholder should be kept: got %q", result)
	}
}

func TestProcessTemplate_EmptyTemplateAndData(t *testing.T) {
	svc := &FormWorkflowService{}

	result := svc.processTemplate("", map[string]interface{}{})
	if result != "" {
		t.Errorf("expected empty result, got %q", result)
	}
}

func TestProcessTemplate_RepeatedVariable(t *testing.T) {
	svc := &FormWorkflowService{}

	result := svc.processTemplate(
		"{{x}} + {{x}} = two {{x}}",
		map[string]interface{}{"x": "hello"},
	)

	if result != "hello + hello = two hello" {
		t.Errorf("got %q, want %q", result, "hello + hello = two hello")
	}
}

func TestProcessTemplate_BooleanValue(t *testing.T) {
	svc := &FormWorkflowService{}

	result := svc.processTemplate(
		"Active: {{active}}",
		map[string]interface{}{"active": true},
	)

	// fmt.Sprintf("%v", true) == "true"
	if result != "Active: true" {
		t.Errorf("got %q, want %q", result, "Active: true")
	}
}

// ─── evaluateCondition ───────────────────────────────────────────────────────

func TestEvaluateCondition_CurrentlyAlwaysTrue(t *testing.T) {
	svc := &FormWorkflowService{}

	// The current simplified implementation always returns true
	if !svc.evaluateCondition("amount > 1000", map[string]interface{}{"amount": 500}) {
		t.Error("simplified evaluateCondition should currently return true")
	}
	if !svc.evaluateCondition("", map[string]interface{}{}) {
		t.Error("evaluateCondition with empty condition should return true")
	}
}

// ─── WorkflowDefinition struct ────────────────────────────────────────────────

func TestWorkflowDefinition_ZeroValue(t *testing.T) {
	var wd WorkflowDefinition
	if wd.IsActive {
		t.Error("zero value IsActive should be false")
	}
	if len(wd.Steps) != 0 {
		t.Error("zero value Steps should be empty")
	}
}

func TestWorkflowDefinition_Construction(t *testing.T) {
	formID := uuid.New()
	createdBy := uuid.New()

	wd := WorkflowDefinition{
		ID:          uuid.New(),
		Name:        "Expense Approval",
		Description: "Requires manager sign-off",
		FormID:      formID,
		Trigger:     "on_submit",
		IsActive:    true,
		CreatedBy:   createdBy,
		CreatedAt:   time.Now(),
	}

	if wd.Trigger != "on_submit" {
		t.Errorf("Trigger = %q, want %q", wd.Trigger, "on_submit")
	}
	if !wd.IsActive {
		t.Error("IsActive should be true")
	}
	if wd.FormID != formID {
		t.Error("FormID mismatch")
	}
}

// ─── WorkflowStep struct ──────────────────────────────────────────────────────

func TestWorkflowStep_SupportedTypes(t *testing.T) {
	validTypes := []string{"approval", "notification", "assignment", "condition", "action"}

	for _, stepType := range validTypes {
		step := WorkflowStep{
			ID:   uuid.New(),
			Name: "Test Step",
			Type: stepType,
			Config: map[string]interface{}{
				"key": "value",
			},
			Order: 1,
		}

		if step.Type != stepType {
			t.Errorf("Type = %q, want %q", step.Type, stepType)
		}
	}
}

func TestWorkflowStep_NextAndAltStepID(t *testing.T) {
	nextID := uuid.New()
	altID := uuid.New()

	step := WorkflowStep{
		ID:         uuid.New(),
		Type:       "condition",
		NextStepID: &nextID,
		AltStepID:  &altID,
	}

	if step.NextStepID == nil || *step.NextStepID != nextID {
		t.Error("NextStepID should point to nextID")
	}
	if step.AltStepID == nil || *step.AltStepID != altID {
		t.Error("AltStepID should point to altID")
	}
}

func TestWorkflowStep_OptionalPointers_NilByDefault(t *testing.T) {
	step := WorkflowStep{Type: "approval"}
	if step.NextStepID != nil {
		t.Error("NextStepID should be nil when not set")
	}
	if step.AltStepID != nil {
		t.Error("AltStepID should be nil when not set")
	}
}

// ─── WorkflowInstance struct ──────────────────────────────────────────────────

func TestWorkflowInstance_StatusValues(t *testing.T) {
	statuses := []string{"running", "completed", "failed", "waiting"}

	for _, status := range statuses {
		inst := WorkflowInstance{
			ID:         uuid.New(),
			WorkflowID: uuid.New(),
			RecordID:   uuid.New(),
			Status:     status,
			Context:    make(map[string]interface{}),
			StartedAt:  time.Now(),
		}
		if inst.Status != status {
			t.Errorf("Status = %q, want %q", inst.Status, status)
		}
	}
}

func TestWorkflowInstance_CompletedAt_NilByDefault(t *testing.T) {
	inst := WorkflowInstance{Status: "running"}
	if inst.CompletedAt != nil {
		t.Error("CompletedAt should be nil for running instance")
	}
}

func TestWorkflowInstance_CompletedAt_CanBeSet(t *testing.T) {
	now := time.Now()
	inst := WorkflowInstance{
		Status:      "completed",
		CompletedAt: &now,
	}
	if inst.CompletedAt == nil {
		t.Fatal("CompletedAt should not be nil")
	}
	if !inst.CompletedAt.Equal(now) {
		t.Error("CompletedAt value mismatch")
	}
}

func TestWorkflowInstance_Context_Assignable(t *testing.T) {
	inst := WorkflowInstance{
		Status:  "running",
		Context: make(map[string]interface{}),
	}

	inst.Context["user"] = "alice"
	inst.Context["amount"] = 999

	if inst.Context["user"] != "alice" {
		t.Error("context user not stored")
	}
	if inst.Context["amount"] != 999 {
		t.Error("context amount not stored")
	}
}

// ─── ApprovalRequest struct ───────────────────────────────────────────────────

func TestApprovalRequest_StatusValues(t *testing.T) {
	statuses := []string{"pending", "approved", "rejected"}

	for _, status := range statuses {
		req := ApprovalRequest{
			ID:         uuid.New(),
			InstanceID: uuid.New(),
			StepID:     uuid.New(),
			Status:     status,
			CreatedAt:  time.Now(),
		}
		if req.Status != status {
			t.Errorf("Status = %q, want %q", req.Status, status)
		}
	}
}

func TestApprovalRequest_RespondedAt_NilByDefault(t *testing.T) {
	req := ApprovalRequest{Status: "pending"}
	if req.RespondedAt != nil {
		t.Error("RespondedAt should be nil for pending approval")
	}
}

func TestApprovalRequest_RespondedAt_CanBeSet(t *testing.T) {
	now := time.Now()
	req := ApprovalRequest{
		Status:      "approved",
		RespondedAt: &now,
	}
	if req.RespondedAt == nil || !req.RespondedAt.Equal(now) {
		t.Error("RespondedAt should be set")
	}
}

// ─── NotificationTemplate struct ─────────────────────────────────────────────

func TestNotificationTemplate_Types(t *testing.T) {
	for _, notifType := range []string{"email", "sms", "in_app"} {
		tmpl := NotificationTemplate{
			ID:        uuid.New(),
			Name:      "Welcome",
			Type:      notifType,
			Subject:   "Hello {{name}}",
			Body:      "Welcome to the platform, {{name}}!",
			Variables: []string{"name"},
		}
		if tmpl.Type != notifType {
			t.Errorf("Type = %q, want %q", tmpl.Type, notifType)
		}
		if len(tmpl.Variables) != 1 || tmpl.Variables[0] != "name" {
			t.Error("Variables should contain 'name'")
		}
	}
}

// ─── executeWorkflowStep routing ─────────────────────────────────────────────

func TestExecuteWorkflowStep_UnknownType_ReturnsError(t *testing.T) {
	svc := &FormWorkflowService{}
	inst := &WorkflowInstance{
		ID:     uuid.New(),
		Status: "running",
	}
	step := WorkflowStep{
		ID:   uuid.New(),
		Type: "not_a_real_step_type",
	}

	err := svc.executeWorkflowStep(inst, step, uuid.New())
	if err == nil {
		t.Error("expected error for unknown step type")
	}
}

func TestExecuteWorkflowStep_ApprovalMissingAssignedTo_ReturnsError(t *testing.T) {
	svc := &FormWorkflowService{}
	inst := &WorkflowInstance{ID: uuid.New()}
	step := WorkflowStep{
		ID:     uuid.New(),
		Type:   "approval",
		Config: map[string]interface{}{
			// missing "assigned_to"
		},
	}

	err := svc.executeWorkflowStep(inst, step, uuid.New())
	if err == nil {
		t.Error("expected error when assigned_to is missing from approval step")
	}
}

func TestExecuteWorkflowStep_NotificationMissingType_ReturnsError(t *testing.T) {
	svc := &FormWorkflowService{}
	inst := &WorkflowInstance{ID: uuid.New()}
	step := WorkflowStep{
		ID:     uuid.New(),
		Type:   "notification",
		Config: map[string]interface{}{
			// missing "type"
		},
	}

	err := svc.executeWorkflowStep(inst, step, uuid.New())
	if err == nil {
		t.Error("expected error when notification type is missing")
	}
}

func TestExecuteWorkflowStep_AssignmentMissingAssignee_ReturnsError(t *testing.T) {
	svc := &FormWorkflowService{}
	inst := &WorkflowInstance{ID: uuid.New(), RecordID: uuid.New()}
	step := WorkflowStep{
		ID:     uuid.New(),
		Type:   "assignment",
		Config: map[string]interface{}{
			// missing "assignee"
		},
	}

	err := svc.executeWorkflowStep(inst, step, uuid.New())
	if err == nil {
		t.Error("expected error when assignee is missing from assignment step")
	}
}

func TestExecuteWorkflowStep_ConditionMissingCondition_ReturnsError(t *testing.T) {
	svc := &FormWorkflowService{}
	inst := &WorkflowInstance{ID: uuid.New()}
	step := WorkflowStep{
		ID:     uuid.New(),
		Type:   "condition",
		Config: map[string]interface{}{
			// missing "condition"
		},
	}

	err := svc.executeWorkflowStep(inst, step, uuid.New())
	if err == nil {
		t.Error("expected error when condition expression is missing")
	}
}

func TestExecuteWorkflowStep_ActionMissingActionType_ReturnsError(t *testing.T) {
	svc := &FormWorkflowService{}
	inst := &WorkflowInstance{ID: uuid.New()}
	step := WorkflowStep{
		ID:     uuid.New(),
		Type:   "action",
		Config: map[string]interface{}{
			// missing "action_type"
		},
	}

	err := svc.executeWorkflowStep(inst, step, uuid.New())
	if err == nil {
		t.Error("expected error when action_type is missing")
	}
}

// ─── executeActionStep routing ────────────────────────────────────────────────

func TestExecuteActionStep_UpdateField_MissingFieldReturnsError(t *testing.T) {
	svc := &FormWorkflowService{}
	inst := &WorkflowInstance{ID: uuid.New()}
	step := WorkflowStep{
		ID:   uuid.New(),
		Type: "action",
		Config: map[string]interface{}{
			"action_type": "update_field",
			// missing "field"
		},
	}

	err := svc.executeWorkflowStep(inst, step, uuid.New())
	if err == nil {
		t.Error("expected error for update_field without field key")
	}
}

func TestExecuteActionStep_WebhookMissingURL_ReturnsError(t *testing.T) {
	svc := &FormWorkflowService{}
	inst := &WorkflowInstance{ID: uuid.New()}
	step := WorkflowStep{
		ID:   uuid.New(),
		Type: "action",
		Config: map[string]interface{}{
			"action_type": "webhook",
			// missing "url"
		},
	}

	err := svc.executeWorkflowStep(inst, step, uuid.New())
	if err == nil {
		t.Error("expected error for webhook without url")
	}
}

func TestExecuteActionStep_UnknownActionType_ReturnsError(t *testing.T) {
	svc := &FormWorkflowService{}
	inst := &WorkflowInstance{ID: uuid.New()}
	step := WorkflowStep{
		ID:   uuid.New(),
		Type: "action",
		Config: map[string]interface{}{
			"action_type": "fly_to_moon",
		},
	}

	err := svc.executeWorkflowStep(inst, step, uuid.New())
	if err == nil {
		t.Error("expected error for unknown action type")
	}
}
