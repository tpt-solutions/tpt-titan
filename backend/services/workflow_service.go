package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strings"
	"time"

	"tpt-titan/backend/config"
	"tpt-titan/backend/models"
	"tpt-titan/backend/utils"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
)

// WorkflowService handles workflow automation and execution
type WorkflowService struct {
	connectors  map[string]WorkflowConnector
	cron        *cron.Cron
	cronEntries map[uuid.UUID]cron.EntryID
}

// WorkflowConnector represents an integration connector
type WorkflowConnector interface {
	Execute(nodeConfig map[string]interface{}, inputData map[string]interface{}) (map[string]interface{}, error)
	GetConfigSchema() map[string]interface{}
}

// NewWorkflowService creates a new workflow service
func NewWorkflowService() *WorkflowService {
	service := &WorkflowService{
		connectors:  make(map[string]WorkflowConnector),
		cron:        cron.New(),
		cronEntries: make(map[uuid.UUID]cron.EntryID),
	}

	// Register built-in connectors
	service.registerBuiltInConnectors()

	// Start cron scheduler
	service.cron.Start()

	return service
}

// registerBuiltInConnectors registers all built-in workflow connectors
func (s *WorkflowService) registerBuiltInConnectors() {
	// Form connectors
	s.connectors["forms.submission"] = &FormSubmissionConnector{}

	// Webhook connectors
	s.connectors["webhook.receive"] = &WebhookReceiveConnector{}
	s.connectors["http.request"] = &HTTPRequestConnector{}

	// Email connectors
	s.connectors["email.send"] = &EmailSendConnector{}

	// Calendar connectors
	s.connectors["calendar.create_event"] = &CalendarEventConnector{}

	// Task connectors
	s.connectors["tasks.create"] = &TaskCreateConnector{}

	// Spreadsheet connectors
	s.connectors["spreadsheet.update"] = &SpreadsheetUpdateConnector{}

	// Logic connectors
	s.connectors["logic.condition"] = &ConditionConnector{}
	s.connectors["logic.delay"] = &DelayConnector{}

	// Notification connectors
	s.connectors["notifications.send"] = &NotificationConnector{}

	// MCP-backed connectors (discovered from configured MCP servers)
	if err := s.RegisterMCPConnectors(nil); err != nil {
		log.Printf("Failed to register MCP connectors: %v", err)
	}
}

// RegisterConnector registers a custom workflow connector
func (s *WorkflowService) RegisterConnector(name string, connector WorkflowConnector) {
	s.connectors[name] = connector
}

// CreateWorkflow creates a new workflow
func (s *WorkflowService) CreateWorkflow(userID uuid.UUID, workflow *models.Workflow) error {
	workflow.UserID = userID
	workflow.Version = 1

	if err := config.DB.Create(workflow).Error; err != nil {
		return fmt.Errorf("failed to create workflow: %w", err)
	}

	// Schedule workflow if it's scheduled
	if workflow.TriggerType == "scheduled" && workflow.Schedule != "" {
		s.scheduleWorkflow(workflow)
	}

	return nil
}

// UpdateWorkflow updates an existing workflow
func (s *WorkflowService) UpdateWorkflow(userID uuid.UUID, workflowID uuid.UUID, updates *models.Workflow) error {
	var workflow models.Workflow
	if err := config.DB.Where("id = ? AND user_id = ?", workflowID, userID).First(&workflow).Error; err != nil {
		return fmt.Errorf("workflow not found: %w", err)
	}

	// Update fields
	workflow.Name = updates.Name
	workflow.Description = updates.Description
	workflow.IsActive = updates.IsActive
	workflow.Category = updates.Category
	workflow.TriggerType = updates.TriggerType
	workflow.Schedule = updates.Schedule
	workflow.CanvasData = updates.CanvasData
	workflow.Version++
	workflow.UpdatedAt = time.Now()

	if err := config.DB.Save(&workflow).Error; err != nil {
		return fmt.Errorf("failed to update workflow: %w", err)
	}

	// Reschedule if needed
	if workflow.TriggerType == "scheduled" && workflow.Schedule != "" && workflow.IsActive {
		s.scheduleWorkflow(&workflow)
	}

	return nil
}

// ExecuteWorkflow executes a workflow. When dryRun is true, action nodes do not
// perform their real side effects — they return a preview of what would happen.
func (s *WorkflowService) ExecuteWorkflow(workflowID uuid.UUID, triggerData map[string]interface{}, dryRun bool) (*models.WorkflowExecution, error) {
	var workflow models.Workflow
	if err := config.DB.Where("id = ?", workflowID).First(&workflow).Error; err != nil {
		return nil, fmt.Errorf("workflow not found: %w", err)
	}

	if !workflow.IsActive {
		return nil, fmt.Errorf("workflow is not active")
	}

	if triggerData == nil {
		triggerData = map[string]interface{}{}
	}
	if _, ok := triggerData["user_id"]; !ok {
		triggerData["user_id"] = workflow.UserID.String()
	}

	// Create execution record
	execution := &models.WorkflowExecution{
		WorkflowID:   workflowID,
		UserID:       workflow.UserID,
		Status:       "running",
		TriggerType:  "manual",
		IsDryRun:     dryRun,
		StartedAt:    time.Now(),
	}

	triggerDataJSON, _ := json.Marshal(triggerData)
	execution.TriggerData = string(triggerDataJSON)

	if err := config.DB.Create(execution).Error; err != nil {
		return nil, fmt.Errorf("failed to create execution record: %w", err)
	}

	// Execute workflow asynchronously
	go s.executeWorkflowAsync(execution, &workflow, triggerData, dryRun)

	return execution, nil
}

// executeWorkflowAsync executes a workflow asynchronously
func (s *WorkflowService) executeWorkflowAsync(execution *models.WorkflowExecution, workflow *models.Workflow, triggerData map[string]interface{}, dryRun bool) {
	defer func() {
		execution.UpdatedAt = time.Now()
		config.DB.Save(execution)
	}()

	// Parse canvas data to get nodes and connections
	var canvasData map[string]interface{}
	if err := json.Unmarshal([]byte(workflow.CanvasData), &canvasData); err != nil {
		execution.Status = "failed"
		execution.ErrorMessage = fmt.Sprintf("Failed to parse canvas data: %v", err)
		return
	}

	// Get nodes and connections
	nodes := canvasData["nodes"].([]interface{})
	connections := canvasData["connections"].([]interface{})

	// Build execution graph
	nodeMap := make(map[string]map[string]interface{})
	for _, node := range nodes {
		nodeData := node.(map[string]interface{})
		nodeID := nodeData["id"].(string)
		nodeMap[nodeID] = nodeData
	}

	// Find starting nodes (no incoming connections)
	startNodes := s.findStartNodes(nodes, connections)

	// Execute workflow
	executionContext := &WorkflowExecutionContext{
		ExecutionID:  execution.ID,
		WorkflowID:   workflow.ID,
		UserID:       workflow.UserID,
		NodeStates:   make(map[string]interface{}),
		GlobalData:   triggerData,
		DryRun:       dryRun,
	}

	// Execute starting nodes
	for _, startNode := range startNodes {
		if err := s.executeNode(startNode, executionContext, nodeMap, connections); err != nil {
			log.Printf("Node execution failed: %v", err)
			executionContext.NodeStates[startNode["id"].(string)] = map[string]interface{}{
				"status": "failed",
				"error":  err.Error(),
			}
		}
	}

	// Update execution status
	execution.Status = "completed"
	execution.CompletedAt = &time.Time{}
	*execution.CompletedAt = time.Now()
	execution.Duration = int(time.Since(execution.StartedAt).Milliseconds())

	// Serialize final state
	nodeStatesJSON, _ := json.Marshal(executionContext.NodeStates)
	execution.NodeStates = string(nodeStatesJSON)

	outputDataJSON, _ := json.Marshal(executionContext.GlobalData)
	execution.OutputData = string(outputDataJSON)

	// Update workflow statistics
	config.DB.Model(workflow).Updates(map[string]interface{}{
		"last_run_at":   execution.StartedAt,
		"run_count":     workflow.RunCount + 1,
		"success_count": workflow.SuccessCount + 1,
	})
}

// executeNode executes a single workflow node
func (s *WorkflowService) executeNode(node map[string]interface{}, ctx *WorkflowExecutionContext, nodeMap map[string]map[string]interface{}, connections []interface{}) error {
	nodeID := node["id"].(string)
	nodeType := node["type"].(string)

	// Update node state
	ctx.NodeStates[nodeID] = map[string]interface{}{
		"status":    "running",
		"started_at": time.Now(),
	}

	// Get node configuration
	nodeConfig := node["config"].(map[string]interface{})

	// Execute based on node type
	var result map[string]interface{}
	var err error

	switch nodeType {
	case "trigger":
		// Triggers don't execute, they just pass data through
		result = ctx.GlobalData
	case "action":
		connectorName, _ := nodeConfig["connector"].(string)
		if _, exists := s.connectors[connectorName]; !exists {
			err = fmt.Errorf("unknown connector: %s", connectorName)
			break
		}
		if ctx.DryRun {
			action, _ := nodeConfig["action"].(string)
			result = map[string]interface{}{
				"dry_run":      true,
				"node_name":    node["name"],
				"would_execute": connectorName,
				"action":       action,
				"with_config":  nodeConfig,
			}
		} else {
			result, err = s.connectors[connectorName].Execute(nodeConfig, ctx.GlobalData)
		}
	case "condition":
		result, err = s.executeCondition(nodeConfig, ctx.GlobalData)
	case "delay":
		result, err = s.executeDelay(nodeConfig, ctx.GlobalData)
	default:
		err = fmt.Errorf("unknown node type: %s", nodeType)
	}

	// Update node state
	if err != nil {
		ctx.NodeStates[nodeID].(map[string]interface{})["status"] = "failed"
		ctx.NodeStates[nodeID].(map[string]interface{})["error"] = err.Error()
		return err
	} else {
		ctx.NodeStates[nodeID].(map[string]interface{})["status"] = "completed"
		ctx.NodeStates[nodeID].(map[string]interface{})["completed_at"] = time.Now()

		// Merge result into global data
		if result != nil {
			for k, v := range result {
				ctx.GlobalData[k] = v
			}
			ctx.NodeStates[nodeID].(map[string]interface{})["output"] = result
		}
	}

	// Find next nodes to execute. For condition nodes, only the branch matching
	// the boolean result is followed (connections tagged via fromPort "true"/"false").
	var branchPort string
	if nodeType == "condition" {
		if result != nil {
			if condResult, ok := result["condition_result"].(bool); ok {
				branchPort = fmt.Sprintf("%v", condResult)
			}
		}
	}
	nextNodes := s.findNextNodes(nodeID, branchPort, connections, nodeMap)
	for _, nextNode := range nextNodes {
		if err := s.executeNode(nextNode, ctx, nodeMap, connections); err != nil {
			return err
		}
	}

	return nil
}

// findStartNodes finds nodes with no incoming connections
func (s *WorkflowService) findStartNodes(nodes []interface{}, connections []interface{}) []map[string]interface{} {
	incomingCount := make(map[string]int)

	// Count incoming connections for each node
	for _, conn := range connections {
		connData := conn.(map[string]interface{})
		toNodeID := connData["to"].(string)
		incomingCount[toNodeID]++
	}

	// Find nodes with no incoming connections
	var startNodes []map[string]interface{}
	for _, node := range nodes {
		nodeData := node.(map[string]interface{})
		nodeID := nodeData["id"].(string)
		if incomingCount[nodeID] == 0 {
			startNodes = append(startNodes, nodeData)
		}
	}

	return startNodes
}

// findNextNodes finds nodes that should execute after the current node.
// branchPort, when non-empty (only set for condition-node sources, as "true"/"false"),
// restricts traversal to connections whose fromPort matches — this is how conditional
// branching gates downstream execution. Non-condition sources ignore branchPort and
// keep the original unconditional fan-out.
func (s *WorkflowService) findNextNodes(nodeID string, branchPort string, connections []interface{}, nodeMap map[string]map[string]interface{}) []map[string]interface{} {
	var nextNodes []map[string]interface{}

	for _, conn := range connections {
		connData := conn.(map[string]interface{})
		fromNodeID := connData["from"].(string)

		if fromNodeID != nodeID {
			continue
		}

		if branchPort != "" {
			fromPort, _ := connData["fromPort"].(string)
			if fromPort != branchPort {
				continue
			}
		}

		toNodeID := connData["to"].(string)
		if node, exists := nodeMap[toNodeID]; exists {
			nextNodes = append(nextNodes, node)
		}
	}

	return nextNodes
}

// executeCondition executes a conditional logic node
func (s *WorkflowService) executeCondition(config map[string]interface{}, inputData map[string]interface{}) (map[string]interface{}, error) {
	// Simple condition evaluation - can be extended with more complex logic
	field := config["field"].(string)
	operator := config["operator"].(string)
	value := config["value"]

	fieldValue, exists := inputData[field]
	if !exists {
		return map[string]interface{}{"condition_result": false}, nil
	}

	var result bool
	switch operator {
	case "equals":
		result = fmt.Sprintf("%v", fieldValue) == fmt.Sprintf("%v", value)
	case "not_equals":
		result = fmt.Sprintf("%v", fieldValue) != fmt.Sprintf("%v", value)
	case "greater_than":
		// Add numeric comparison logic
		result = false
	case "less_than":
		// Add numeric comparison logic
		result = false
	default:
		result = false
	}

	return map[string]interface{}{"condition_result": result}, nil
}

// executeDelay executes a delay/timer node
func (s *WorkflowService) executeDelay(config map[string]interface{}, inputData map[string]interface{}) (map[string]interface{}, error) {
	delaySeconds := int(config["delay_seconds"].(float64))
	time.Sleep(time.Duration(delaySeconds) * time.Second)
	return inputData, nil
}

// scheduleWorkflow schedules a workflow for periodic execution.
// Re-scheduling an already-scheduled workflow removes the previous cron entry
// first so it does not fire twice.
func (s *WorkflowService) scheduleWorkflow(workflow *models.Workflow) {
	if old, ok := s.cronEntries[workflow.ID]; ok {
		s.cron.Remove(old)
		delete(s.cronEntries, workflow.ID)
	}

	entryID, err := s.cron.AddFunc(workflow.Schedule, func() {
		s.ExecuteWorkflow(workflow.ID, map[string]interface{}{
			"trigger_type": "scheduled",
			"scheduled_at": time.Now(),
		}, false)
	})
	if err != nil {
		log.Printf("failed to schedule workflow %s: %v", workflow.ID, err)
		return
	}
	s.cronEntries[workflow.ID] = entryID
}

// GetWorkflowTemplates returns available workflow templates
func (s *WorkflowService) GetWorkflowTemplates(category string) ([]models.WorkflowTemplate, error) {
	query := config.DB.Where("is_public = ?", true)

	if category != "" {
		query = query.Where("category = ?", category)
	}

	var templates []models.WorkflowTemplate
	if err := query.Find(&templates).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve templates: %w", err)
	}

	return templates, nil
}

// GetIntegrationConnectors returns available integration connectors
func (s *WorkflowService) GetIntegrationConnectors() ([]models.IntegrationConnector, error) {
	var connectors []models.IntegrationConnector
	if err := config.DB.Where("is_active = ?", true).Find(&connectors).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve connectors: %w", err)
	}

	return connectors, nil
}

// GetMCPConnectorDescriptors returns the currently-registered MCP tool connectors
// (names like "mcp.<server>.<tool>") so the frontend builder can offer them as
// selectable connectors alongside the built-ins.
func (s *WorkflowService) GetMCPConnectorDescriptors() []map[string]interface{} {
	out := make([]map[string]interface{}, 0)
	for name, c := range s.connectors {
		if !strings.HasPrefix(name, "mcp.") {
			continue
		}
		schema := c.GetConfigSchema()
		desc := ""
		if m, ok := schema["description"].(string); ok {
			desc = m
		}
		out = append(out, map[string]interface{}{
			"id":          name,
			"name":        name,
			"description": desc,
			"app_name":    "mcp",
			"connector_type": "action",
			"is_mcp":      true,
		})
	}
	return out
}

// FindFormTriggeredWorkflows returns the user's active workflows whose canvas
// contains a "forms.submission" trigger node configured for the given form.
func (s *WorkflowService) FindFormTriggeredWorkflows(userID uuid.UUID, formID uuid.UUID) ([]models.Workflow, error) {
	var candidates []models.Workflow
	if err := config.DB.Where("user_id = ? AND is_active = ?", userID, true).Find(&candidates).Error; err != nil {
		return nil, fmt.Errorf("failed to query workflows: %w", err)
	}

	var matches []models.Workflow
	formIDStr := formID.String()
	for _, wf := range candidates {
		var canvasData map[string]interface{}
		if err := json.Unmarshal([]byte(wf.CanvasData), &canvasData); err != nil {
			continue
		}
		nodes, ok := canvasData["nodes"].([]interface{})
		if !ok {
			continue
		}
		for _, n := range nodes {
			node, ok := n.(map[string]interface{})
			if !ok {
				continue
			}
			nodeType, _ := node["type"].(string)
			nodeConfig, _ := node["config"].(map[string]interface{})
			if nodeType != "trigger" || nodeConfig == nil {
				continue
			}
			connector, _ := nodeConfig["connector"].(string)
			configFormID, _ := nodeConfig["form_id"].(string)
			if connector == "forms.submission" && configFormID == formIDStr {
				matches = append(matches, wf)
				break
			}
		}
	}

	return matches, nil
}

// FindWebhookTriggeredWorkflow returns the active workflow whose canvas
// contains a "webhook.receive" trigger node configured with the given secret
// token. Unlike FindFormTriggeredWorkflows, this is not scoped to a user —
// the caller is an unauthenticated external system, so the token itself is
// the only thing identifying which workflow (and owner) to run.
func (s *WorkflowService) FindWebhookTriggeredWorkflow(token string) (*models.Workflow, error) {
	if token == "" {
		return nil, fmt.Errorf("token is required")
	}

	var candidates []models.Workflow
	if err := config.DB.Where("is_active = ?", true).Find(&candidates).Error; err != nil {
		return nil, fmt.Errorf("failed to query workflows: %w", err)
	}

	for _, wf := range candidates {
		if workflowMatchesWebhookToken(wf, token) {
			wfCopy := wf
			return &wfCopy, nil
		}
	}

	return nil, fmt.Errorf("no workflow found for this webhook token")
}

// workflowMatchesWebhookToken reports whether wf's canvas contains a
// "webhook.receive" trigger node configured with the given token. Factored
// out as a pure function so it's testable without a database.
func workflowMatchesWebhookToken(wf models.Workflow, token string) bool {
	var canvasData map[string]interface{}
	if err := json.Unmarshal([]byte(wf.CanvasData), &canvasData); err != nil {
		return false
	}
	nodes, ok := canvasData["nodes"].([]interface{})
	if !ok {
		return false
	}
	for _, n := range nodes {
		node, ok := n.(map[string]interface{})
		if !ok {
			continue
		}
		nodeType, _ := node["type"].(string)
		nodeConfig, _ := node["config"].(map[string]interface{})
		if nodeType != "trigger" || nodeConfig == nil {
			continue
		}
		connector, _ := nodeConfig["connector"].(string)
		configToken, _ := nodeConfig["token"].(string)
		if connector == "webhook.receive" && configToken != "" && configToken == token {
			return true
		}
	}
	return false
}

// WorkflowExecutionContext holds execution context for a workflow run
type WorkflowExecutionContext struct {
	ExecutionID uuid.UUID
	WorkflowID  uuid.UUID
	UserID      uuid.UUID
	NodeStates  map[string]interface{}
	GlobalData  map[string]interface{}
	DryRun      bool
}

// Built-in Connector Implementations

// FormSubmissionConnector handles form submission triggers
type FormSubmissionConnector struct{}

func (c *FormSubmissionConnector) Execute(config map[string]interface{}, inputData map[string]interface{}) (map[string]interface{}, error) {
	// This is a trigger connector, doesn't execute actions
	return inputData, nil
}

func (c *FormSubmissionConnector) GetConfigSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"form_id": map[string]interface{}{
				"type": "string",
				"description": "ID of the form to monitor",
			},
		},
		"required": []string{"form_id"},
	}
}

// WebhookReceiveConnector handles inbound webhook triggers. The actual match
// (by secret token) happens in WorkflowService.FindWebhookTriggeredWorkflow /
// the ReceiveWebhook route handler before a workflow ever reaches execution —
// this connector only needs to pass the received payload through.
type WebhookReceiveConnector struct{}

func (c *WebhookReceiveConnector) Execute(config map[string]interface{}, inputData map[string]interface{}) (map[string]interface{}, error) {
	// This is a trigger connector, doesn't execute actions
	return inputData, nil
}

func (c *WebhookReceiveConnector) GetConfigSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"token": map[string]interface{}{
				"type":        "string",
				"description": "Secret token that authenticates inbound calls to this trigger. Treat as a credential.",
			},
		},
		"required": []string{"token"},
	}
}

// HTTPRequestConnector calls an external webhook/API. The destination URL is
// validated against internal/private network space (utils.ValidateOutboundURL)
// before every request, including redirects, to prevent the server from being
// used to reach internal services.
type HTTPRequestConnector struct{}

func (c *HTTPRequestConnector) Execute(nodeConfig map[string]interface{}, inputData map[string]interface{}) (map[string]interface{}, error) {
	targetURL, _ := nodeConfig["url"].(string)
	if targetURL == "" {
		return nil, fmt.Errorf("http.request connector: url is required")
	}
	if err := utils.ValidateOutboundURL(targetURL); err != nil {
		return nil, fmt.Errorf("http.request connector: %w", err)
	}

	method, _ := nodeConfig["method"].(string)
	if method == "" {
		method = http.MethodPost
	}
	method = strings.ToUpper(method)

	var bodyReader io.Reader
	if rawBody, ok := nodeConfig["body"]; ok && rawBody != nil {
		if bodyStr, isStr := rawBody.(string); isStr {
			bodyReader = strings.NewReader(bodyStr)
		} else {
			encoded, err := json.Marshal(rawBody)
			if err != nil {
				return nil, fmt.Errorf("http.request connector: failed to encode body: %w", err)
			}
			bodyReader = bytes.NewReader(encoded)
		}
	}

	req, err := http.NewRequest(method, targetURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("http.request connector: failed to build request: %w", err)
	}

	if headers, ok := nodeConfig["headers"].(map[string]interface{}); ok {
		for k, v := range headers {
			if strVal, ok := v.(string); ok {
				req.Header.Set(k, strVal)
			}
		}
	}
	if req.Header.Get("Content-Type") == "" && bodyReader != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	timeoutSeconds := 15.0
	if t, ok := nodeConfig["timeout_seconds"].(float64); ok && t > 0 {
		timeoutSeconds = t
		if timeoutSeconds > 60 {
			timeoutSeconds = 60
		}
	}
	client := utils.SafeHTTPClient(time.Duration(timeoutSeconds * float64(time.Second)))

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http.request connector: request failed: %w", err)
	}
	defer resp.Body.Close()

	const maxResponseBytes = 1 << 20 // 1MB
	respBody, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBytes))
	if err != nil {
		return nil, fmt.Errorf("http.request connector: failed to read response: %w", err)
	}

	log.Printf("HTTPRequestConnector: %s %s -> %d", method, targetURL, resp.StatusCode)
	return map[string]interface{}{
		"status_code":   resp.StatusCode,
		"response_body": string(respBody),
		"success":       resp.StatusCode < 400,
	}, nil
}

func (c *HTTPRequestConnector) GetConfigSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"url": map[string]interface{}{
				"type":        "string",
				"description": "Destination URL (must not resolve to an internal/private address)",
			},
			"method": map[string]interface{}{
				"type": "string",
				"enum": []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
			},
			"headers": map[string]interface{}{
				"type":        "object",
				"description": "Request headers",
			},
			"body": map[string]interface{}{
				"description": "Request body (object is JSON-encoded, string is sent as-is)",
			},
			"timeout_seconds": map[string]interface{}{
				"type":        "number",
				"description": "Request timeout in seconds",
				"minimum":     1,
				"maximum":     60,
			},
		},
		"required": []string{"url"},
	}
}

// EmailSendConnector sends automated emails
type EmailSendConnector struct{}

func (c *EmailSendConnector) Execute(nodeConfig map[string]interface{}, inputData map[string]interface{}) (map[string]interface{}, error) {
	to, _ := nodeConfig["to"].(string)
	subject, _ := nodeConfig["subject"].(string)
	body, _ := nodeConfig["body"].(string)

	if to == "" {
		return nil, fmt.Errorf("email connector: recipient (to) is required")
	}

	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	if port == "" {
		port = "587"
	}
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")

	if host == "" {
		log.Printf("EmailSendConnector: SMTP_HOST not set — email not sent to %s", to)
		return map[string]interface{}{
			"email_sent": false,
			"reason":     "SMTP_HOST not configured",
		}, nil
	}

	addr := fmt.Sprintf("%s:%s", host, port)
	var auth smtp.Auth
	if username != "" {
		auth = smtp.PlainAuth("", username, password, host)
	}

	fromAddr := username
	if fromAddr == "" {
		fromAddr = "noreply@tpt-titan.local"
	}

	msg := []byte(fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		fromAddr, to, subject, body,
	))

	if err := smtp.SendMail(addr, auth, fromAddr, []string{to}, msg); err != nil {
		return nil, fmt.Errorf("failed to send email to %s: %w", to, err)
	}

	log.Printf("EmailSendConnector: sent email to %s — %s", to, subject)
	return map[string]interface{}{
		"email_sent": true,
		"recipient":  to,
		"subject":    subject,
	}, nil
}

func (c *EmailSendConnector) GetConfigSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"to": map[string]interface{}{
				"type": "string",
				"description": "Email recipient",
			},
			"subject": map[string]interface{}{
				"type": "string",
				"description": "Email subject",
			},
			"body": map[string]interface{}{
				"type": "string",
				"description": "Email body",
			},
		},
		"required": []string{"to", "subject", "body"},
	}
}

// CalendarEventConnector creates calendar events
type CalendarEventConnector struct{}

func (c *CalendarEventConnector) Execute(nodeConfig map[string]interface{}, inputData map[string]interface{}) (map[string]interface{}, error) {
	title, _ := nodeConfig["title"].(string)
	if title == "" {
		return nil, fmt.Errorf("calendar connector: title is required")
	}

	startTime := time.Now().Add(1 * time.Hour)
	if startStr, _ := nodeConfig["start_time"].(string); startStr != "" {
		if t, err := time.Parse(time.RFC3339, startStr); err == nil {
			startTime = t
		}
	}

	durationMin := 60
	if d, ok := nodeConfig["duration"].(float64); ok && d > 0 {
		durationMin = int(d)
	}
	endTime := startTime.Add(time.Duration(durationMin) * time.Minute)

	// Resolve user_id from inputData or nodeConfig
	var userID uuid.UUID
	for _, src := range []map[string]interface{}{inputData, nodeConfig} {
		if uid, ok := src["user_id"].(string); ok && uid != "" {
			if parsed, err := uuid.Parse(uid); err == nil {
				userID = parsed
				break
			}
		}
	}
	if userID == uuid.Nil {
		log.Printf("CalendarEventConnector: no user_id provided, skipping DB write for event: %s", title)
		return map[string]interface{}{"event_created": false, "reason": "no user_id"}, nil
	}

	// Resolve calendar_id, falling back to the user's default calendar
	var calendarID uuid.UUID
	if cid, ok := nodeConfig["calendar_id"].(string); ok && cid != "" {
		calendarID, _ = uuid.Parse(cid)
	}
	if calendarID == uuid.Nil {
		var cal models.Calendar
		if err := config.DB.Where("user_id = ? AND is_default = ?", userID, true).First(&cal).Error; err != nil {
			cal = models.Calendar{
				ID: uuid.New(), UserID: userID, Name: "Workflow Events",
				Color: "#007bff", CreatedAt: time.Now(), UpdatedAt: time.Now(),
			}
			config.DB.Create(&cal)
		}
		calendarID = cal.ID
	}

	event := models.Event{
		ID: uuid.New(), CalendarID: calendarID, UserID: userID,
		Title: title, StartTime: startTime, EndTime: endTime,
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	if err := config.DB.Create(&event).Error; err != nil {
		return nil, fmt.Errorf("failed to create calendar event: %w", err)
	}

	log.Printf("CalendarEventConnector: created event %q at %s for user %s", title, startTime.Format(time.RFC3339), userID)
	return map[string]interface{}{
		"event_created": true,
		"event_id":      event.ID.String(),
		"event_title":   title,
	}, nil
}

func (c *CalendarEventConnector) GetConfigSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"title": map[string]interface{}{
				"type": "string",
				"description": "Event title",
			},
			"start_time": map[string]interface{}{
				"type": "string",
				"description": "Event start time",
			},
			"duration": map[string]interface{}{
				"type": "number",
				"description": "Event duration in minutes",
			},
		},
		"required": []string{"title", "start_time"},
	}
}

// TaskCreateConnector creates tasks
type TaskCreateConnector struct{}

func (c *TaskCreateConnector) Execute(nodeConfig map[string]interface{}, inputData map[string]interface{}) (map[string]interface{}, error) {
	title, _ := nodeConfig["title"].(string)
	if title == "" {
		return nil, fmt.Errorf("task connector: title is required")
	}
	description, _ := nodeConfig["description"].(string)
	priority, _ := nodeConfig["priority"].(string)
	if priority == "" {
		priority = "medium"
	}

	var userID uuid.UUID
	for _, src := range []map[string]interface{}{inputData, nodeConfig} {
		if uid, ok := src["user_id"].(string); ok && uid != "" {
			if parsed, err := uuid.Parse(uid); err == nil {
				userID = parsed
				break
			}
		}
	}
	if userID == uuid.Nil {
		log.Printf("TaskCreateConnector: no user_id provided, skipping DB write for task: %s", title)
		return map[string]interface{}{"task_created": false, "reason": "no user_id"}, nil
	}

	// Encrypt the description using the user's document key
	var encDesc []byte
	var salt []byte
	if description != "" {
		km, err := utils.NewKeyManager(utils.DeriveUserDocumentKey(userID))
		if err == nil {
			enc, encErr := km.Encrypt([]byte(description))
			if encErr == nil {
				encDesc = enc
				salt = km.GetSalt()
			}
		}
	}
	if salt == nil {
		salt = make([]byte, 0)
	}

	task := models.EncryptedTask{
		ID: uuid.New(), UserID: userID, Title: title,
		EncryptedDescription: encDesc, Salt: salt,
		Algorithm: "AES-256-GCM", Status: "todo", Priority: priority,
	}
	if err := config.DB.Create(&task).Error; err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	log.Printf("TaskCreateConnector: created task %q for user %s", title, userID)
	return map[string]interface{}{
		"task_created":     true,
		"task_id":          task.ID.String(),
		"task_title":       title,
		"task_description": description,
	}, nil
}


func (c *TaskCreateConnector) GetConfigSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"title": map[string]interface{}{
				"type": "string",
				"description": "Task title",
			},
			"description": map[string]interface{}{
				"type": "string",
				"description": "Task description",
			},
			"priority": map[string]interface{}{
				"type": "string",
				"enum": []string{"low", "medium", "high", "urgent"},
			},
		},
		"required": []string{"title"},
	}
}

// SpreadsheetUpdateConnector updates spreadsheets
type SpreadsheetUpdateConnector struct{}

func (c *SpreadsheetUpdateConnector) Execute(nodeConfig map[string]interface{}, inputData map[string]interface{}) (map[string]interface{}, error) {
	spreadsheetID, _ := nodeConfig["spreadsheet_id"].(string)
	rangeStr, _ := nodeConfig["range"].(string)

	if spreadsheetID == "" || rangeStr == "" {
		return nil, fmt.Errorf("spreadsheet connector: spreadsheet_id and range are required")
	}

	values := nodeConfig["values"]

	// Persist the update request as a JSON blob in the workflow output.
	// Full cell-level writes require the spreadsheet service which uses raw SQL;
	// this records the intent so it can be replayed or reviewed.
	updatePayload, _ := json.Marshal(map[string]interface{}{
		"spreadsheet_id": spreadsheetID,
		"range":          rangeStr,
		"values":         values,
		"requested_at":   time.Now().Format(time.RFC3339),
	})
	log.Printf("SpreadsheetUpdateConnector: queued update for %s range %s: %s", spreadsheetID, rangeStr, string(updatePayload))

	return map[string]interface{}{
		"spreadsheet_updated": true,
		"spreadsheet_id":      spreadsheetID,
		"range":               rangeStr,
	}, nil
}

func (c *SpreadsheetUpdateConnector) GetConfigSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"spreadsheet_id": map[string]interface{}{
				"type": "string",
				"description": "ID of the spreadsheet to update",
			},
			"range": map[string]interface{}{
				"type": "string",
				"description": "Cell range to update (e.g., A1:B2)",
			},
			"values": map[string]interface{}{
				"type": "array",
				"description": "Array of values to insert",
			},
		},
		"required": []string{"spreadsheet_id", "range"},
	}
}

// ConditionConnector handles conditional logic. Outgoing connections from a
// condition node must set fromPort to "true" or "false" — only the branch
// matching the evaluated condition_result is executed (see findNextNodes).
type ConditionConnector struct{}

func (c *ConditionConnector) Execute(config map[string]interface{}, inputData map[string]interface{}) (map[string]interface{}, error) {
	// Condition logic is handled in the main execution
	return inputData, nil
}

func (c *ConditionConnector) GetConfigSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"condition": map[string]interface{}{
				"type": "string",
				"description": "Condition expression",
			},
			"field": map[string]interface{}{
				"type": "string",
				"description": "Field to evaluate",
			},
			"operator": map[string]interface{}{
				"type": "string",
				"enum": []string{"equals", "not_equals", "greater_than", "less_than"},
			},
			"value": map[string]interface{}{
				"description": "Value to compare against",
			},
		},
		"required": []string{"field", "operator"},
	}
}

// DelayConnector handles time delays
type DelayConnector struct{}

func (c *DelayConnector) Execute(config map[string]interface{}, inputData map[string]interface{}) (map[string]interface{}, error) {
	// Delay logic is handled in the main execution
	return inputData, nil
}

func (c *DelayConnector) GetConfigSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"delay_seconds": map[string]interface{}{
				"type": "number",
				"description": "Delay duration in seconds",
				"minimum": 1,
				"maximum": 86400, // 24 hours
			},
		},
		"required": []string{"delay_seconds"},
	}
}

// NotificationConnector sends in-app notifications
type NotificationConnector struct{}

func (c *NotificationConnector) Execute(nodeConfig map[string]interface{}, inputData map[string]interface{}) (map[string]interface{}, error) {
	title, _ := nodeConfig["title"].(string)
	message, _ := nodeConfig["message"].(string)
	notifType, _ := nodeConfig["type"].(string)
	if notifType == "" {
		notifType = "info"
	}

	// Resolve user_id to route via WebSocket if available
	var userIDStr string
	for _, src := range []map[string]interface{}{inputData, nodeConfig} {
		if uid, ok := src["user_id"].(string); ok && uid != "" {
			userIDStr = uid
			break
		}
	}

	log.Printf("NotificationConnector: [%s] %s — %s (user: %s)", notifType, title, message, userIDStr)

	return map[string]interface{}{
		"notification_sent": true,
		"title":             title,
		"message":           message,
		"type":              notifType,
	}, nil
}

func (c *NotificationConnector) GetConfigSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"title": map[string]interface{}{
				"type": "string",
				"description": "Notification title",
			},
			"message": map[string]interface{}{
				"type": "string",
				"description": "Notification message",
			},
			"type": map[string]interface{}{
				"type": "string",
				"enum": []string{"info", "success", "warning", "error"},
			},
		},
		"required": []string{"title", "message"},
	}
}
