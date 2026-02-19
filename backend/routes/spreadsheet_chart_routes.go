package routes

import (
	"database/sql"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"tpt-titan/backend/config"
	"tpt-titan/backend/services"
)

// GenerateChartSuggestion analyzes spreadsheet data and suggests charts
func GenerateChartSuggestion(c *gin.Context) {
	var req struct {
		Data      map[string]interface{} `json:"data" binding:"required"`
		Range     string                 `json:"range,omitempty"`
		DataTypes map[string]string      `json:"data_types,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	suggestions := analyzeDataForCharts(req.Data, req.Range, req.DataTypes)

	c.JSON(http.StatusOK, gin.H{
		"suggestions": suggestions,
	})
}

// CreateChart creates a chart from spreadsheet data
func CreateChart(c *gin.Context) {
	var req struct {
		SpreadsheetID uuid.UUID              `json:"spreadsheet_id" binding:"required"`
		ChartType     string                 `json:"chart_type" binding:"required"`
		DataRange     string                 `json:"data_range" binding:"required"`
		Title         string                 `json:"title"`
		XAxisLabel    string                 `json:"x_axis_label,omitempty"`
		YAxisLabel    string                 `json:"y_axis_label,omitempty"`
		Data          map[string]interface{} `json:"data,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate chart type
	validTypes := []string{"bar", "line", "pie", "scatter", "area"}
	validType := false
	for _, t := range validTypes {
		if req.ChartType == t {
			validType = true
			break
		}
	}

	if !validType {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chart type"})
		return
	}

	chartID := uuid.New()

	chart := gin.H{
		"id":            chartID,
		"spreadsheet_id": req.SpreadsheetID,
		"type":          req.ChartType,
		"data_range":    req.DataRange,
		"title":         req.Title,
		"x_axis_label":  req.XAxisLabel,
		"y_axis_label":  req.YAxisLabel,
		"data":          req.Data,
		"created_at":    "now",
	}

	c.JSON(http.StatusCreated, chart)
}

// GetCharts retrieves charts for a spreadsheet
func GetCharts(c *gin.Context) {
	spreadsheetIDStr := c.Param("id")
	_, err := uuid.Parse(spreadsheetIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid spreadsheet ID"})
		return
	}

	// In a real implementation, fetch from database
	// For now, return mock data
	charts := []gin.H{
		{
			"id":            uuid.New(),
			"type":          "bar",
			"data_range":    "A1:B5",
			"title":         "Sales Data",
			"x_axis_label":  "Products",
			"y_axis_label":  "Revenue",
		},
	}

	c.JSON(http.StatusOK, gin.H{"charts": charts})
}

// analyzeDataForCharts analyzes spreadsheet data and suggests appropriate charts
func analyzeDataForCharts(data map[string]interface{}, dataRange string, dataTypes map[string]string) []gin.H {
	suggestions := []gin.H{}
