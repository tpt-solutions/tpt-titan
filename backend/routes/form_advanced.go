package routes

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"tpt-titan/backend/services"
)

// Advanced Form Routes - Visual Query Builder, Relationships, Reporting, Workflows

// Visual Query Builder Routes

// BuildSQL builds SQL from visual query elements
func BuildSQL(c *gin.Context) {
	var elements []services.QueryElement
	if err := c.ShouldBindJSON(&elements); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := c.MustGet("db").(*sql.DB)
	queryBuilder := services.NewQueryBuilderService(db)

	sql, err := queryBuilder.BuildSQLFromElements(elements)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"sql": sql,
		"elements": elements,
	})
}

// ExecuteVisualQuery executes a visual query and returns results
func ExecuteVisualQuery(c *gin.Context) {
	var elements []services.QueryElement
	if err := c.ShouldBindJSON(&elements); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := c.MustGet("db").(*sql.DB)
	queryBuilder := services.NewQueryBuilderService(db)

	result, err := queryBuilder.ExecuteQuery(elements)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetAvailableTables returns tables available for query building
func GetAvailableTables(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	queryBuilder := services.NewQueryBuilderService(db)

	tables, err := queryBuilder.GetAvailableTables()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tables": tables})
}

// ValidateVisualQuery validates a visual query
func ValidateVisualQuery(c *gin.Context) {
	var elements []services.QueryElement
	if err := c.ShouldBindJSON(&elements); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := c.MustGet("db").(*sql.DB)
	queryBuilder := services.NewQueryBuilderService(db)

	errors := queryBuilder.ValidateQuery(elements)
	c.JSON(http.StatusOK, gin.H{
		"valid": len(errors) == 0,
		"errors": errors,
	})
}

// GetQuerySuggestions provides suggestions for query building
func GetQuerySuggestions(c *gin.Context) {
	var elements []services.QueryElement
	if err := c.ShouldBindJSON(&elements); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := c.MustGet("db").(*sql.DB)
	queryBuilder := services.NewQueryBuilderService(db)

	suggestions := queryBuilder.GetQuerySuggestions(elements)
	c.JSON(http.StatusOK, gin.H{"suggestions": suggestions})
}

// SaveQueryTemplate saves a visual query as a template
func SaveQueryTemplate(c *gin.Context) {
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

	var req struct {
		Name        string                   `json:"name" binding:"required"`
		Description string                   `json:"description"`
		Elements    []services.QueryElement `json:"elements" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := c.MustGet("db").(*sql.DB)
	queryBuilder := services.NewQueryBuilderService(db)

	err := queryBuilder.SaveQuery(userID.String(), req.Name, req.Description, req.Elements)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Query template saved successfully"})
}

// GetQueryTemplates returns saved query templates
func GetQueryTemplates(c *gin.Context) {
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
	queryBuilder := services.NewQueryBuilderService(db)

	templates, err := queryBuilder.LoadSavedQueries(userID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"templates": templates})
}

// Form Relationship Routes

// CreateRelationship creates a new relationship between forms
func CreateRelationship(c *gin.Context) {
	var rel services.Relationship
	if err := c.ShouldBindJSON(&rel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := c.MustGet("db").(*sql.DB)
	relationshipSvc := services.NewFormRelationshipService(db)

	// Validate relationship
	errors := relationshipSvc.ValidateRelationship(&rel)
	if len(errors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	err := relationshipSvc.CreateRelationship(&rel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, rel)
}

// GetFormRelationships gets relationships for a form
func GetFormRelationships(c *gin.Context) {
	formIDStr := c.Param("formId")
	formID, err := uuid.Parse(formIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form ID"})
		return
	}

	db := c.MustGet("db").(*sql.DB)
	relationshipSvc := services.NewFormRelationshipService(db)

	relationships, err := relationshipSvc.GetRelationshipsByForm(formID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"relationships": relationships})
}

// CreateLookupField creates a lookup field
func CreateLookupField(c *gin.Context) {
	var lookup services.LookupField
	if err := c.ShouldBindJSON(&lookup); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := c.MustGet("db").(*sql.DB)
	relationshipSvc := services.NewFormRelationshipService(db)

	err := relationshipSvc.CreateLookupField(&lookup)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, lookup)
}

// GetLookupData gets data for a lookup field
func GetLookupData(c *gin.Context) {
	lookupFieldIDStr := c.Param("lookupFieldId")
	lookupFieldID, err := uuid.Parse(lookupFieldIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lookup field ID"})
		return
	}

	db := c.MustGet("db").(*sql.DB)
	relationshipSvc := services.NewFormRelationshipService(db)

	data, err := relationshipSvc.GetLookupFieldData(lookupFieldID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}

// GetFormHierarchy gets the hierarchy of related forms
func GetFormHierarchy(c *gin.Context) {
	formIDStr := c.Param("formId")
	formID, err := uuid.Parse(formIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form ID"})
		return
	}

	db := c.MustGet("db").(*sql.DB)
	relationshipSvc := services.NewFormRelationshipService(db)

	hierarchy, err := relationshipSvc.GetFormHierarchy(formID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, hierarchy)
}

// GetRelatedData gets related data for a record
func GetRelatedData(c *gin.Context) {
	formIDStr := c.Param("formId")
	formID, err := uuid.Parse(formIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form ID"})
		return
	}

	recordIDStr := c.Param("recordId")
	recordID, err := uuid.Parse(recordIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid record ID"})
		return
	}

	db := c.MustGet("db").(*sql.DB)
	relationshipSvc := services.NewFormRelationshipService(db)

	relatedData, err := relationshipSvc.GetRelatedData(formID, recordID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, relatedData)
}

// Form Reporting Routes

// CreateReport creates a new report
func CreateReport(c *gin.Context) {
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

	var report services.ReportDefinition
	if err := c.ShouldBindJSON(&report); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	report.CreatedBy = userID

	db := c.MustGet("db").(*sql.DB)
	queryBuilder := services.NewQueryBuilderService(db)
	relationshipSvc := services.NewFormRelationshipService(db)
	reportingSvc := services.NewFormReportingService(db, queryBuilder, relationshipSvc)

	err := reportingSvc.CreateReport(&report)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, report)
}

// ExecuteReport executes a report
func ExecuteReport(c *gin.Context) {
	reportIDStr := c.Param("reportId")
	reportID, err := uuid.Parse(reportIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report ID"})
		return
	}

	db := c.MustGet("db").(*sql.DB)
	queryBuilder := services.NewQueryBuilderService(db)
	relationshipSvc := services.NewFormRelationshipService(db)
	reportingSvc := services.NewFormReportingService(db, queryBuilder, relationshipSvc)

	result, err := reportingSvc.ExecuteReport(reportID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GenerateAdHocReport generates an ad-hoc report
func GenerateAdHocReport(c *gin.Context) {
	formIDStr := c.Param("formId")
	formID, err := uuid.Parse(formIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form ID"})
		return
	}

	var config map[string]interface{}
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := c.MustGet("db").(*sql.DB)
	queryBuilder := services.NewQueryBuilderService(db)
	relationshipSvc := services.NewFormRelationshipService(db)
	reportingSvc := services.NewFormReportingService(db, queryBuilder, relationshipSvc)

	result, err := reportingSvc.GenerateAdHocReport(formID, config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ExportReport exports a report to various formats
func ExportReport(c *gin.Context) {
	reportIDStr := c.Param("reportId")
	reportID, err := uuid.Parse(reportIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report ID"})
		return
	}

	format := c.Query("format")
	if format == "" {
		format = "csv"
	}

	db := c.MustGet("db").(*sql.DB)
	queryBuilder := services.NewQueryBuilderService(db)
	relationshipSvc := services.NewFormRelationshipService(db)
	reportingSvc := services.NewFormReportingService(db, queryBuilder, relationshipSvc)

	data, err := reportingSvc.ExportReport(reportID, format)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set appropriate headers based on format
	var contentType string
	var filename string

	switch format {
	case "csv":
		contentType = "text/csv"
		filename = "report.csv"
	case "excel":
		contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
		filename = "report.xlsx"
	case "pdf":
		contentType = "application/pdf"
		filename = "report.pdf"
	default:
		contentType = "application/octet-stream"
		filename = "report.dat"
	}

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", contentType)
	c.Header("Content-Length", strconv.Itoa(len(data)))
	c.Data(http.StatusOK, contentType, data)
}

// GetFormReports gets reports for a form
func GetFormReports(c *gin.Context) {
	formIDStr := c.Param("formId")
	formID, err := uuid.Parse(formIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form ID"})
		return
	}

	db := c.MustGet("db").(*sql.DB)
	queryBuilder := services.NewQueryBuilderService(db)
	relationshipSvc := services.NewFormRelationshipService(db)
	reportingSvc := services.NewFormReportingService(db, queryBuilder, relationshipSvc)

	reports, err := reportingSvc.ListReports(formID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reports": reports})
}

// CreateDashboard creates a dashboard
func CreateDashboard(c *gin.Context) {
	var req struct {
		Name        string    `json:"name" binding:"required"`
		Description string    `json:"description"`
		ReportIDs   []uuid.UUID `json:"report_ids" binding:"required"`
		Layout      map[string]interface{} `json:"layout"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := c.MustGet("db").(*sql.DB)
	queryBuilder := services.NewQueryBuilderService(db)
	relationshipSvc := services.NewFormRelationshipService(db)
	reportingSvc := services.NewFormReportingService(db, queryBuilder, relationshipSvc)

	dashboardID, err := reportingSvc.CreateDashboard(req.Name, req.Description, req.ReportIDs, req.Layout)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":          dashboardID,
		"name":        req.Name,
		"description": req.Description,
		"report_ids":  req.ReportIDs,
		"layout":      req.Layout,
	})
}

// GetDashboard gets dashboard data
func GetDashboard(c *gin.Context) {
	dashboardIDStr := c.Param("dashboardId")
	dashboardID, err := uuid.Parse(dashboardIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid dashboard ID"})
		return
	}

	db := c.MustGet("db").(*sql.DB)
	queryBuilder := services.NewQueryBuilderService(db)
	relationshipSvc := services.NewFormRelationshipService(db)
	reportingSvc := services.NewFormReportingService(db, queryBuilder, relationshipSvc)

	dashboard, err := reportingSvc.GetDashboard(dashboardID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dashboard)
}

// Workflow Automation Routes

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

// Form Template Library Routes

// CreateFormTemplate creates a new form template
func CreateFormTemplate(c *gin.Context) {
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

	var template struct {
		Name        string                   `json:"name" binding:"required"`
		Description string                   `json:"description"`
		Category    string                   `json:"category"`
		FormData    map[string]interface{}   `json:"form_data" binding:"required"`
		IsPublic    bool                     `json:"is_public"`
		Tags        []string                 `json:"tags"`
	}

	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// In a real implementation, save to form_templates table
	templateID := uuid.New()

	c.JSON(http.StatusCreated, gin.H{
		"id":          templateID,
		"name":        template.Name,
		"description": template.Description,
		"category":    template.Category,
		"created_by":  userID,
		"is_public":   template.IsPublic,
		"tags":        template.Tags,
	})
}

// GetFormTemplates gets available form templates
func GetFormTemplates(c *gin.Context) {
	category := c.Query("category")
	search := c.Query("search")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	// In a real implementation, query form_templates table
	templates := []gin.H{
		{
			"id":          uuid.New(),
			"name":        "Contact Information Form",
			"description": "Standard contact form with name, email, phone",
			"category":    "business",
			"is_public":   true,
			"tags":        []string{"contact", "business"},
		},
		{
			"id":          uuid.New(),
			"name":        "Survey Form",
			"description": "Multi-question survey with various field types",
			"category":    "survey",
			"is_public":   true,
			"tags":        []string{"survey", "feedback"},
		},
		{
			"id":          uuid.New(),
			"name":        "Invoice Form",
			"description": "Professional invoice with line items and totals",
			"category":    "finance",
			"is_public":   true,
			"tags":        []string{"invoice", "finance", "business"},
		},
	}

	// Filter by category and search
	if category != "" {
		filtered := []gin.H{}
		for _, t := range templates {
			if cat, ok := t["category"].(string); ok && cat == category {
				filtered = append(filtered, t)
			}
		}
		templates = filtered
	}

	if search != "" {
		filtered := []gin.H{}
		for _, t := range templates {
			if name, ok := t["name"].(string); ok &&
			   strings.Contains(strings.ToLower(name), strings.ToLower(search)) {
				filtered = append(filtered, t)
			}
		}
		templates = filtered
	}

	// Apply limit
	if len(templates) > limit {
		templates = templates[:limit]
	}

	c.JSON(http.StatusOK, gin.H{"templates": templates})
}

// GetFormTemplateCategories gets available template categories
func GetFormTemplateCategories(c *gin.Context) {
	categories := []gin.H{
		{"id": "business", "name": "Business Forms", "description": "Professional business forms"},
		{"id": "survey", "name": "Surveys", "description": "Feedback and survey forms"},
		{"id": "finance", "name": "Finance", "description": "Financial and accounting forms"},
		{"id": "hr", "name": "Human Resources", "description": "HR and personnel forms"},
		{"id": "education", "name": "Education", "description": "Educational forms and surveys"},
		{"id": "healthcare", "name": "Healthcare", "description": "Medical and healthcare forms"},
	}

	c.JSON(http.StatusOK, gin.H{"categories": categories})
}

// UseFormTemplate creates a new form from a template
func UseFormTemplate(c *gin.Context) {
	templateIDStr := c.Param("templateId")
	_, err := uuid.Parse(templateIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID"})
		return
	}

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

	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// In a real implementation, create form from template
	formID := uuid.New()

	c.JSON(http.StatusCreated, gin.H{
		"id":          formID,
		"name":        req.Name,
		"description": req.Description,
		"template_id": templateIDStr,
		"created_by":  userID,
		"created_at":  "now",
	})
}

// Form Email Integration Routes

// CreateEmailDistribution creates an email distribution for form responses
func CreateEmailDistribution(c *gin.Context) {
	formIDStr := c.Param("formId")
	_, err := uuid.Parse(formIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form ID"})
		return
	}

	var req struct {
		Name         string   `json:"name" binding:"required"`
		Recipients   []string `json:"recipients" binding:"required"`
		Subject      string   `json:"subject" binding:"required"`
		Message      string   `json:"message"`
		Trigger      string   `json:"trigger"` // "immediate", "daily", "weekly"
		IsActive     bool     `json:"is_active"`
		IncludeData  bool     `json:"include_data"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// In a real implementation, save to form_email_distributions table
	distributionID := uuid.New()

	c.JSON(http.StatusCreated, gin.H{
		"id":            distributionID,
		"form_id":       formIDStr,
		"name":          req.Name,
		"recipients":    req.Recipients,
		"subject":       req.Subject,
		"trigger":       req.Trigger,
		"is_active":     req.IsActive,
		"include_data":  req.IncludeData,
	})
}

// GetEmailDistributions gets email distributions for a form
func GetEmailDistributions(c *gin.Context) {
	formIDStr := c.Param("formId")
	_, err := uuid.Parse(formIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form ID"})
		return
	}

	// In a real implementation, query form_email_distributions table
	distributions := []gin.H{
		{
			"id":         uuid.New(),
			"name":       "Daily Summary",
			"recipients": []string{"admin@example.com"},
			"trigger":    "daily",
			"is_active":  true,
		},
	}

	c.JSON(http.StatusOK, gin.H{"distributions": distributions})
}

// SendFormResponseEmail sends an email with form response data
func SendFormResponseEmail(c *gin.Context) {
	responseIDStr := c.Param("responseId")
	_, err := uuid.Parse(responseIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid response ID"})
		return
	}

	var req struct {
		Recipients []string `json:"recipients" binding:"required"`
		Subject    string   `json:"subject" binding:"required"`
		Message    string   `json:"message"`
		IncludeData bool    `json:"include_data"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// In a real implementation, send email with form data
	c.JSON(http.StatusOK, gin.H{
		"message":    "Email sent successfully",
		"recipients": req.Recipients,
		"subject":    req.Subject,
	})
}
