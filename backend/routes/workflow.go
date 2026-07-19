package routes

import (
	"encoding/json"
	"net/http"
	"tpt-titan/backend/config"
	"tpt-titan/backend/models"
	"tpt-titan/backend/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var workflowService *services.WorkflowService

// InitWorkflowService initializes the workflow service (called from main)
func InitWorkflowService() {
	workflowService = services.NewWorkflowService()
}

// CreateWorkflow creates a new workflow
func CreateWorkflow(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var workflow models.Workflow
	if err := c.ShouldBindJSON(&workflow); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := workflowService.CreateWorkflow(userID.(uuid.UUID), &workflow); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"workflow": workflow})
}

// GetWorkflows returns all workflows for the user
func GetWorkflows(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var workflows []models.Workflow
	query := config.DB.Where("user_id = ?", userID)

	// Optional filtering
	category := c.Query("category")
	if category != "" {
		query = query.Where("category = ?", category)
	}

	isActive := c.Query("active")
	if isActive == "true" {
		query = query.Where("is_active = ?", true)
	} else if isActive == "false" {
		query = query.Where("is_active = ?", false)
	}

	if err := query.Order("updated_at DESC").Find(&workflows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve workflows"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"workflows": workflows})
}

// GetWorkflow returns a specific workflow
func GetWorkflow(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	workflowID := c.Param("id")
	if workflowID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Workflow ID is required"})
		return
	}

	wfID, err := uuid.Parse(workflowID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workflow ID"})
		return
	}

	var workflow models.Workflow
	if err := config.DB.Where("id = ? AND user_id = ?", wfID, userID).First(&workflow).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workflow not found"})
		return
	}

	// Get workflow nodes
	var nodes []models.WorkflowNode
	config.DB.Where("workflow_id = ?", wfID).Find(&nodes)

	// Get workflow connections
	var connections []models.WorkflowConnection
	config.DB.Where("workflow_id = ?", wfID).Find(&connections)

	c.JSON(http.StatusOK, gin.H{
		"workflow":    workflow,
		"nodes":       nodes,
		"connections": connections,
	})
}

// UpdateWorkflow updates an existing workflow
func UpdateWorkflow(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	workflowID := c.Param("id")
	if workflowID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Workflow ID is required"})
		return
	}

	wfID, err := uuid.Parse(workflowID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workflow ID"})
		return
	}

	var updates models.Workflow
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := workflowService.UpdateWorkflow(userID.(uuid.UUID), wfID, &updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Workflow updated successfully"})
}

// DeleteWorkflow deletes a workflow
func DeleteWorkflow(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	workflowID := c.Param("id")
	if workflowID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Workflow ID is required"})
		return
	}

	wfID, err := uuid.Parse(workflowID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workflow ID"})
		return
	}

	// Delete workflow (cascade will handle nodes and connections)
	if err := config.DB.Where("id = ? AND user_id = ?", wfID, userID).Delete(&models.Workflow{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete workflow"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Workflow deleted successfully"})
}

// ExecuteWorkflow executes a workflow
func ExecuteWorkflow(c *gin.Context) {
	if _, exists := c.Get("user_id"); !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	workflowID := c.Param("id")
	if workflowID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Workflow ID is required"})
		return
	}

	wfID, err := uuid.Parse(workflowID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workflow ID"})
		return
	}

	var payload struct {
		TriggerData map[string]interface{} `json:"trigger_data,omitempty"`
		DryRun      bool                    `json:"dry_run,omitempty"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	execution, err := workflowService.ExecuteWorkflow(wfID, payload.TriggerData, payload.DryRun)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"execution_id": execution.ID,
		"status":       execution.Status,
		"message":      "Workflow execution started",
	})
}

// GetWorkflowExecution gets the status of a workflow execution
func GetWorkflowExecution(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	executionID := c.Param("executionId")
	if executionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Execution ID is required"})
		return
	}

	execID, err := uuid.Parse(executionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid execution ID"})
		return
	}

	var execution models.WorkflowExecution
	if err := config.DB.Where("id = ? AND user_id = ?", execID, userID).First(&execution).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Execution not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"execution": execution})
}

// GetWorkflowExecutions returns execution history for a workflow
func GetWorkflowExecutions(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	workflowID := c.Param("id")
	if workflowID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Workflow ID is required"})
		return
	}

	wfID, err := uuid.Parse(workflowID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workflow ID"})
		return
	}

	var executions []models.WorkflowExecution
	query := config.DB.Where("workflow_id = ? AND user_id = ?", wfID, userID).Order("started_at DESC")

	// Pagination
	limit := 50
	limitStr := c.Query("limit")
	if limitStr != "" {
		if l := parseInt(limitStr); l > 0 && l <= 100 {
			limit = l
		}
	}

	if err := query.Limit(limit).Find(&executions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve executions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"executions": executions})
}

// UpdateWorkflowNodes updates the nodes for a workflow
func UpdateWorkflowNodes(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	workflowID := c.Param("id")
	if workflowID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Workflow ID is required"})
		return
	}

	wfID, err := uuid.Parse(workflowID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workflow ID"})
		return
	}

	// Verify workflow ownership
	var workflow models.Workflow
	if err := config.DB.Where("id = ? AND user_id = ?", wfID, userID).First(&workflow).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workflow not found"})
		return
	}

	var payload struct {
		Nodes []models.WorkflowNode `json:"nodes"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Delete existing nodes
	config.DB.Where("workflow_id = ?", wfID).Delete(&models.WorkflowNode{})

	// Insert new nodes
	for i := range payload.Nodes {
		payload.Nodes[i].WorkflowID = wfID
		payload.Nodes[i].ID = uuid.New() // Ensure new ID
	}

	if err := config.DB.Create(&payload.Nodes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update nodes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Nodes updated successfully"})
}

// UpdateWorkflowConnections updates the connections for a workflow
func UpdateWorkflowConnections(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	workflowID := c.Param("id")
	if workflowID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Workflow ID is required"})
		return
	}

	wfID, err := uuid.Parse(workflowID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workflow ID"})
		return
	}

	// Verify workflow ownership
	var workflow models.Workflow
	if err := config.DB.Where("id = ? AND user_id = ?", wfID, userID).First(&workflow).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workflow not found"})
		return
	}

	var payload struct {
		Connections []models.WorkflowConnection `json:"connections"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Delete existing connections
	config.DB.Where("workflow_id = ?", wfID).Delete(&models.WorkflowConnection{})

	// Insert new connections
	for i := range payload.Connections {
		payload.Connections[i].WorkflowID = wfID
		payload.Connections[i].ID = uuid.New() // Ensure new ID
	}

	if err := config.DB.Create(&payload.Connections).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update connections"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Connections updated successfully"})
}

// GetWorkflowTemplates returns available workflow templates
func GetWorkflowTemplates(c *gin.Context) {
	category := c.Query("category")

	templates, err := workflowService.GetWorkflowTemplates(category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve templates"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"templates": templates})
}

// CreateWorkflowFromTemplate creates a workflow from a template
func CreateWorkflowFromTemplate(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	templateID := c.Param("templateId")
	if templateID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Template ID is required"})
		return
	}

	tmplID, err := uuid.Parse(templateID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID"})
		return
	}

	var template models.WorkflowTemplate
	if err := config.DB.Where("id = ? AND is_public = ?", tmplID, true).First(&template).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
		return
	}

	// Create workflow from template
	var templateData map[string]interface{}
	json.Unmarshal([]byte(template.TemplateData), &templateData)

	workflow := &models.Workflow{
		UserID:      userID.(uuid.UUID),
		Name:        template.Name + " (Copy)",
		Description: template.Description,
		Category:    template.Category,
		TriggerType: "manual",
		// Ships inactive: instantiating a preset is opt-in, but running it (and
		// producing any real side effects) requires the user to explicitly
		// activate it — ideally after a dry run.
		IsActive:   false,
		CanvasData: template.TemplateData,
	}

	if err := workflowService.CreateWorkflow(userID.(uuid.UUID), workflow); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"workflow": workflow})
}

// GetIntegrationConnectors returns available integration connectors
func GetIntegrationConnectors(c *gin.Context) {
	connectors, err := workflowService.GetIntegrationConnectors()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve connectors"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"connectors": connectors})
}

// GetNamedConnectorTemplates returns browse-only named connector presets
// (Slack/Discord/GitHub) that wrap the http.request connector with the right
// URL/header conventions pre-filled. The frontend builder offers these as
// starting points; selecting one yields an http.request node pre-filled from
// each template's "config".
func GetNamedConnectorTemplates(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"templates": workflowService.GetNamedConnectorTemplates()})
}
