package services

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"tpt-titan/backend/models"

	"github.com/google/uuid"
)

// WorkflowCanvas is the JSON structure the builder and execution engine use.
// It mirrors the shape produced/consumed by the workflow canvas UI.
type WorkflowCanvas struct {
	Nodes       []map[string]interface{} `json:"nodes"`
	Connections []map[string]interface{} `json:"connections"`
}

// knownConnectors is the allowlist of connector names the AI generator may
// emit. Anything outside this set is rejected during validation so a generated
// workflow can never reference a connector that doesn't exist (and thus can
// never silently no-op). This is the safety gate the brainstorm asked for.
var knownConnectors = map[string]bool{
	"forms.submission":        true,
	"webhook.receive":         true,
	"http.request":            true,
	"email.send":              true,
	"calendar.create_event":   true,
	"tasks.create":            true,
	"spreadsheet.update":      true,
	"logic.condition":         true,
	"logic.delay":             true,
	"notifications.send":      true,
}

// aiGenerateFunc is the model-calling hook. It is a package var so tests can
// swap in a deterministic stub without a real Ollama/OpenRouter backend.
var aiGenerateFunc = defaultAIGenerate

// aiListModelsFunc selects the model to generate with. It is a package var so
// tests can stub model selection without a database.
var aiListModelsFunc = func(ai *AIService, userID uuid.UUID) ([]models.AIModel, error) {
	return ai.GetAvailableModels(userID)
}

// GenerateWorkflowFromPrompt asks the configured AI to turn a natural-language
// description into a workflow canvas. It returns a *draft* canvas only — the
// caller must let the user review/approve and then persist it (the dry-run /
// approval flow the brainstorm required). No workflow is created here.
func (s *WorkflowAIService) GenerateWorkflowFromPrompt(userID uuid.UUID, prompt string) (*WorkflowCanvas, error) {
	prompt = strings.TrimSpace(prompt)
	if prompt == "" {
		return nil, fmt.Errorf("a description is required")
	}

	models, err := aiListModelsFunc(s.aiService, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to load AI models: %w", err)
	}
	if len(models) == 0 {
		return nil, fmt.Errorf("no AI model is configured; add one in Settings before generating workflows")
	}
	// Prefer an online model if available, else the first enabled model.
	chosen := models[0]
	for _, m := range models {
		if m.Provider == "openrouter" {
			chosen = m
			break
		}
	}

	raw, err := aiGenerateFunc(s.aiService, chosen, prompt)
	if err != nil {
		return nil, fmt.Errorf("AI generation failed: %w", err)
	}

	canvas, err := parseWorkflowCanvas(raw)
	if err != nil {
		return nil, fmt.Errorf("generated workflow was malformed: %w", err)
	}
	return canvas, nil
}

// defaultAIGenerate calls the real configured provider for the chosen model.
func defaultAIGenerate(ai *AIService, model models.AIModel, prompt string) (string, error) {
	full := buildWorkflowGenerationPrompt(prompt)
	switch model.Provider {
	case "ollama":
		if ai.config == nil || !ai.config.EnableLocalAI {
			return "", fmt.Errorf("local AI is disabled")
		}
		return ai.ollamaService.GenerateResponse(model.ModelID, full)
	case "openrouter":
		if ai.config == nil || !ai.config.EnableOnlineAI {
			return "", fmt.Errorf("online AI is disabled")
		}
		return ai.openRouterService.GenerateResponse(model.ModelID, full)
	default:
		return "", fmt.Errorf("unsupported AI provider %q", model.Provider)
	}
}

// buildWorkflowGenerationPrompt wraps the user's description with a strict
// schema instruction so the model returns only JSON we can parse.
func buildWorkflowGenerationPrompt(userPrompt string) string {
	var b strings.Builder
	b.WriteString("You are a workflow automation builder for TPT Titan. ")
	b.WriteString("Given the user's description, output a workflow as a single JSON object and nothing else. ")
	b.WriteString("The JSON must have exactly two keys: \"nodes\" and \"connections\".\n")
	b.WriteString("Each node is an object: {\"id\": string, \"type\": \"trigger\"|\"action\"|\"condition\"|\"delay\", \"name\": string, \"config\": object}.\n")
	b.WriteString("Each connection is an object: {\"from\": nodeId, \"to\": nodeId, optional \"fromPort\": \"true\"|\"false\" (only from condition nodes)}.\n")
	b.WriteString("Allowed action connectors (set config.connector to one of these): forms.submission, webhook.receive (trigger), http.request, email.send, calendar.create_event, tasks.create, spreadsheet.update, logic.condition, logic.delay, notifications.send.\n")
	b.WriteString("Condition nodes must set config.field, config.operator (equals|not_equals|greater_than|less_than), config.value; their outgoing connections must set fromPort.\n")
	b.WriteString("Use only the connectors above. Return JSON only, no markdown fences.\n\nUser description: ")
	b.WriteString(userPrompt)
	return b.String()
}

// parseWorkflowCanvas extracts and validates a WorkflowCanvas from a model
// response. It tolerates ```json fences and surrounding prose by locating the
// outermost {...} block, then validates node types, connector names, and that
// every connection references real node ids.
func parseWorkflowCanvas(raw string) (*WorkflowCanvas, error) {
	jsonStr := extractJSONBlock(raw)
	if jsonStr == "" {
		return nil, fmt.Errorf("no JSON object found in model output")
	}

	var canvas WorkflowCanvas
	if err := json.Unmarshal([]byte(jsonStr), &canvas); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}
	if len(canvas.Nodes) == 0 {
		return nil, fmt.Errorf("generated workflow has no nodes")
	}

	ids := make(map[string]bool, len(canvas.Nodes))
	for _, n := range canvas.Nodes {
		id, _ := n["id"].(string)
		if id == "" {
			return nil, fmt.Errorf("a node is missing its id")
		}
		if ids[id] {
			return nil, fmt.Errorf("duplicate node id %q", id)
		}
		ids[id] = true

		ntype, _ := n["type"].(string)
		switch ntype {
		case "trigger", "action", "condition", "delay":
		default:
			return nil, fmt.Errorf("node %q has invalid type %q", id, ntype)
		}

		if ntype == "action" {
			cfg, _ := n["config"].(map[string]interface{})
			connector, _ := cfg["connector"].(string)
			if connector == "" {
				return nil, fmt.Errorf("action node %q is missing a connector", id)
			}
			if !knownConnectors[connector] {
				return nil, fmt.Errorf("action node %q references unknown connector %q", id, connector)
			}
		}
	}

	for _, c := range canvas.Connections {
		from, _ := c["from"].(string)
		to, _ := c["to"].(string)
		if from == "" || to == "" {
			return nil, fmt.Errorf("a connection is missing from/to")
		}
		if !ids[from] || !ids[to] {
			return nil, fmt.Errorf("connection references unknown node id (from=%q to=%q)", from, to)
		}
	}

	return &canvas, nil
}

// extractJSONBlock returns the first balanced {...} substring, tolerating
// markdown code fences and trailing prose from the model.
func extractJSONBlock(raw string) string {
	clean := strings.TrimSpace(raw)
	// Strip ```json ... ``` fences if present.
	if fence := regexp.MustCompile("(?s)```(?:json)?\\s*(.*?)```"); fence.MatchString(clean) {
		if m := fence.FindStringSubmatch(clean); len(m) == 2 {
			clean = strings.TrimSpace(m[1])
		}
	}

	start := strings.Index(clean, "{")
	if start < 0 {
		return ""
	}
	depth := 0
	for i := start; i < len(clean); i++ {
		switch clean[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return clean[start : i+1]
			}
		}
	}
	return ""
}
