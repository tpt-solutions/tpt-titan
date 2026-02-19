package routes

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"tpt-titan/backend/services"
)

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
