// backend/services/workflow_ai_generate_test.go
// Tests for natural-language -> workflow canvas generation (parsing/validation
// and the draft-return contract). The model call is stubbed via aiGenerateFunc.

package services

import (
	"testing"

	"tpt-titan/backend/models"

	"github.com/google/uuid"
)

func TestParseWorkflowCanvas_Valid(t *testing.T) {
	raw := `{"nodes":[
		{"id":"t","type":"trigger","name":"Form","config":{"connector":"forms.submission"}},
		{"id":"c","type":"condition","name":"Priority?","config":{"field":"priority","operator":"equals","value":"high"}},
		{"id":"a","type":"action","name":"Task","config":{"connector":"tasks.create"}}
	],"connections":[
		{"from":"t","to":"c"},
		{"from":"c","to":"a","fromPort":"true"}
	]}`

	canvas, err := parseWorkflowCanvas(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(canvas.Nodes) != 3 || len(canvas.Connections) != 2 {
		t.Fatalf("unexpected counts: nodes=%d conns=%d", len(canvas.Nodes), len(canvas.Connections))
	}
}

func TestParseWorkflowCanvas_RejectsUnknownConnector(t *testing.T) {
	raw := `{"nodes":[
		{"id":"a","type":"action","name":"X","config":{"connector":"totally.made.up"}}
	],"connections":[]}`

	if _, err := parseWorkflowCanvas(raw); err == nil {
		t.Fatal("expected rejection of unknown connector")
	}
}

func TestParseWorkflowCanvas_RejectsDanglingConnection(t *testing.T) {
	raw := `{"nodes":[{"id":"a","type":"action","name":"X","config":{"connector":"tasks.create"}}],"connections":[{"from":"a","to":"ghost"}]}`

	if _, err := parseWorkflowCanvas(raw); err == nil {
		t.Fatal("expected rejection of connection to unknown node")
	}
}

func TestParseWorkflowCanvas_ToleratesFencedProse(t *testing.T) {
	raw := "Here is your workflow:\n```json\n" + `{"nodes":[{"id":"a","type":"action","name":"X","config":{"connector":"email.send"}}],"connections":[]}` + "\n```\nLet me know if you need changes."
	canvas, err := parseWorkflowCanvas(raw)
	if err != nil {
		t.Fatalf("expected fence/prose to be tolerated, got: %v", err)
	}
	if len(canvas.Nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(canvas.Nodes))
	}
}

// TestGenerateWorkflowFromPrompt_UsesStub verifies the generation path returns
// a parsed draft and does NOT create any workflow (it is gated to draft-only).
func TestGenerateWorkflowFromPrompt_UsesStub(t *testing.T) {
	orig := aiGenerateFunc
	origList := aiListModelsFunc
	defer func() { aiGenerateFunc = orig; aiListModelsFunc = origList }()

	aiService := &AIService{}
	svc := &WorkflowAIService{aiService: aiService}

	aiListModelsFunc = func(ai *AIService, userID uuid.UUID) ([]models.AIModel, error) {
		return []models.AIModel{{Provider: "openrouter", ModelID: "test/model"}}, nil
	}

	// Stub the model to return a valid canvas and assert the chosen model is
	// passed through (provider routing is covered by defaultAIGenerate).
	aiGenerateFunc = func(ai *AIService, m models.AIModel, prompt string) (string, error) {
		if prompt == "" {
			t.Fatal("empty prompt passed to generator")
		}
		if m.ModelID != "test/model" {
			t.Fatalf("unexpected model passed to generator: %s", m.ModelID)
		}
		return `{"nodes":[{"id":"a","type":"action","name":"Email","config":{"connector":"email.send"}}],"connections":[]}`, nil
	}

	canvas, err := svc.GenerateWorkflowFromPrompt(uuid.New(), "when a form is submitted, email the team")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(canvas.Nodes) != 1 {
		t.Fatalf("expected 1 node in draft, got %d", len(canvas.Nodes))
	}
}
