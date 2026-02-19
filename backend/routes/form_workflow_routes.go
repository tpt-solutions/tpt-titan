package routes

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"tpt-titan/backend/services"
)

// CreateWorkflow creates a new workflow
func CreateWorkflow(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	var workflow services.WorkflowDefinition
	if err := c.ShouldBindJSON(&workflow); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	workflow.CreatedBy = userID

	db := c.MustGet("db").(*sql.DB)
	emailSvc := services.NewEmailService(db, nil) // Would need proper email config
	relationshipSvc := services.NewFormRelationshipService(db)
	workflowSvc := services.NewFormWorkflowService(db, emailSvc, relationshipSvc)

	err := workflowSvc.CreateWorkflow(&workflow)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, workflow)
}

// StartWorkflow starts a workflow instance
func StartWorkflow(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	workflowIDStr := c.Param("workflowId")
	workflowID, err := uuid.Parse(workflowIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workflow ID"})
		return
	}

	recordIDStr := c.Param("recordId")
	recordID, err := uuid.Parse(recordIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid record ID"})
		return
	}

	db := c.MustGet("db").(*sql.DB)
	emailSvc := services.NewEmailService(db, nil)
	relationshipSvc := services.NewFormRelationshipService(db)
	workflowSvc := services.NewFormWorkflowService(db, emailSvc, relationshipSvc)

	err = workflowSvc.StartWorkflow(workflowID, recordID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Workflow started successfully"})
}

// ProcessApproval processes an approval response
func ProcessApproval(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	approvalIDStr := c.Param("approvalId")
	approvalID, err := uuid.Parse(approvalIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid approval ID"})
		return
	}

	var req struct {
		Approved bool   `json:"approved"`
		Message  string `json:"message"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := c.MustGet("db").(*sql.DB)
	emailSvc := services.NewEmailService(db, nil)
	relationshipSvc := services.NewFormRelationshipService(db)
	workflowSvc := services.NewFormWorkflowService(db, emailSvc, relationshipSvc)

	err = workflowSvc.ProcessApproval(approvalID, req.Approved, req.Message, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Approval processed successfully"})
}

// GetPendingApprovals gets pending approvals for the current user
func GetPendingApprovals(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	db := c.MustGet("db").(*sql.DB)
	emailSvc := services.NewEmailService(db, nil)
	relationshipSvc := services.NewFormRelationshipService(db)
	workflowSvc := services.NewFormWorkflowService(db, emailSvc, relationshipSvc)

	approvals, err := workflowSvc.GetPendingApprovals(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"approvals": approvals})
}

// CreateNotificationTemplate creates a notification template
func CreateNotificationTemplate(c *gin.Context) {
	var template services.NotificationTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := c.MustGet("db").(*sql.DB)
	emailSvc := services.NewEmailService(db, nil)
	relationshipSvc := services.NewFormRelationshipService(db)
	workflowSvc := services.NewFormWorkflowService(db, emailSvc, relationshipSvc)

	err := workflowSvc.CreateNotificationTemplate(&template)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, template)
}
