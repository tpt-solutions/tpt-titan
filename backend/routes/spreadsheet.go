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
	"tpt-titan/backend/services"
)

// SpreadsheetData represents spreadsheet data structure
type SpreadsheetData struct {
	ID       uuid.UUID              `json:"id"`
	Name     string                 `json:"name"`
	OwnerID  uuid.UUID              `json:"owner_id"`
	Data     map[string]interface{} `json:"data"` // Cell data as map[cellRef]interface{}
	Formulas map[string]string      `json:"formulas,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// EvaluateFormulaRequest represents a formula evaluation request
type EvaluateFormulaRequest struct {
	Formula     string                 `json:"formula" binding:"required"`
	CellContext map[string]interface{} `json:"cell_context,omitempty"` // Cell values for references
}

// EvaluateFormulaResponse represents formula evaluation response
type EvaluateFormulaResponse struct {
	Result     interface{}            `json:"result"`
	DataType   string                 `json:"data_type"`
	Error      string                 `json:"error,omitempty"`
	DependsOn  []string               `json:"depends_on,omitempty"`
}

// CreateSpreadsheet creates a new spreadsheet
func CreateSpreadsheet(c *gin.Context) {
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
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create spreadsheet in database (simplified - would use proper model)
	spreadsheetID := uuid.New()

	c.JSON(http.StatusCreated, gin.H{
		"id":       spreadsheetID,
		"name":     req.Name,
		"owner_id": userID,
		"data":     make(map[string]interface{}),
		"formulas": make(map[string]string),
		"metadata": gin.H{
			"created_at": "now",
			"version":    1,
		},
	})
}

// GetSpreadsheet retrieves a spreadsheet by ID
func GetSpreadsheet(c *gin.Context) {
	spreadsheetIDStr := c.Param("id")
	spreadsheetID, err := uuid.Parse(spreadsheetIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid spreadsheet ID"})
		return
	}

	// In a real implementation, fetch from database
	// For now, return mock data
	c.JSON(http.StatusOK, gin.H{
		"id":   spreadsheetID,
		"name": "Sample Spreadsheet",
		"data": gin.H{
			"A1": "Product",
			"B1": "Price",
			"C1": "Quantity",
			"D1": "Total",
			"A2": "Widget A",
			"B2": 10.99,
			"C2": 5,
			"D2": "=B2*C2",
			"A3": "Widget B",
			"B3": 15.50,
			"C3": 3,
			"D3": "=B3*C3",
			"A4": "Total",
			"B4": "",
			"C4": "",
			"D4": "=SUM(D2:D3)",
		},
		"formulas": gin.H{
			"D2": "=B2*C2",
			"D3": "=B3*C3",
			"D4": "=SUM(D2:D3)",
		},
	})
}

// UpdateSpreadsheetCell updates a specific cell in the spreadsheet
func UpdateSpreadsheetCell(c *gin.Context) {
	spreadsheetIDStr := c.Param("id")
	_, err := uuid.Parse(spreadsheetIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid spreadsheet ID"})
		return
	}

	var req struct {
		CellReference string      `json:"cell_reference" binding:"required"`
		Value         interface{} `json:"value"`
		Formula       string      `json:"formula,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate cell reference
	if !isValidCellReference(req.CellReference) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cell reference"})
		return
	}

	// In a real implementation, update database
	c.JSON(http.StatusOK, gin.H{
		"message":        "Cell updated successfully",
		"cell_reference": req.CellReference,
		"value":          req.Value,
		"formula":        req.Formula,
	})
}

// EvaluateFormula evaluates a spreadsheet formula
func EvaluateFormula(c *gin.Context) {
	var req EvaluateFormulaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get spreadsheet math service
	mathService := services.NewSpreadsheetMathService()

	// Create cell resolver function
	cellResolver := func(cellRef string) (interface{}, error) {
		if value, exists := req.CellContext[cellRef]; exists {
			return value, nil
		}
		// Try to parse as number
		if num, err := strconv.ParseFloat(cellRef, 64); err == nil {
			return num, nil
		}
		return cellRef, nil // Return as string
	}

	// Evaluate formula
	result, err := mathService.EvaluateFormula(req.Formula, cellResolver)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	response := EvaluateFormulaResponse{
		Result:   result.Value,
		DataType: result.DataType,
	}

	if result.Error != "" {
		response.Error = result.Error
	}

	if result.DependsOn != nil {
		response.DependsOn = result.DependsOn
	}

	c.JSON(http.StatusOK, response)
}

// GetAvailableFunctions returns list of available mathematical functions
func GetAvailableFunctions(c *gin.Context) {
	mathService := services.NewSpreadsheetMathService()
	functions := mathService.GetAvailableFunctions()

	// Convert to response format
	funcList := make([]gin.H, 0, len(functions))
	for name, fn := range functions {
		funcList = append(funcList, gin.H{
			"name":        name,
			"description": fn.Description,
			"min_args":    fn.MinArgs,
			"max_args":    fn.MaxArgs,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"functions": funcList,
	})
}

// ValidateFormula validates formula syntax
func ValidateFormula(c *gin.Context) {
	var req struct {
		Formula string `json:"formula" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mathService := services.NewSpreadsheetMathService()
	err := mathService.ValidateFormula(req.Formula)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"valid": false,
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid": true,
	})
}

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

// Helper functions

func isValidCellReference(ref string) bool {
	ref = strings.ToUpper(strings.TrimSpace(ref))

	// Basic validation for A1, B2, etc.
	if matched, _ := regexp.MatchString(`^[A-Z]+\d+$`, ref); matched {
		return true
	}

	// Handle ranges like A1:B3
	if strings.Contains(ref, ":") {
		parts := strings.Split(ref, ":")
		if len(parts) == 2 {
			return isValidCellReference(parts[0]) && isValidCellReference(parts[1])
		}
	}

	return false
}

// analyzeDataForCharts analyzes spreadsheet data and suggests appropriate charts
func analyzeDataForCharts(data map[string]interface{}, dataRange string, dataTypes map[string]string) []gin.H {
	suggestions := []gin.H{}

	// Analyze data patterns and suggest charts
	// This is a simplified implementation

	// Check for numeric columns that could be charted
	numericColumns := []string{}
	categoricalColumns := []string{}

	for cellRef, value := range data {
		if dataTypes != nil {
			if dataType, exists := dataTypes[cellRef]; exists && dataType == "number" {
				numericColumns = append(numericColumns, cellRef)
			} else if dataType == "string" {
				categoricalColumns = append(categoricalColumns, cellRef)
			}
		} else {
			// Try to infer type
			switch value.(type) {
			case int, int32, int64, float32, float64:
				numericColumns = append(numericColumns, cellRef)
			case string:
				categoricalColumns = append(categoricalColumns, cellRef)
			}
		}
	}

	// Suggest bar chart for categorical vs numeric data
	if len(categoricalColumns) > 0 && len(numericColumns) > 0 {
		suggestions = append(suggestions, gin.H{
			"type":        "bar",
			"description": "Bar chart showing " + strings.Join(categoricalColumns, ", ") + " vs " + strings.Join(numericColumns, ", "),
			"confidence":  0.8,
			"data_range":  dataRange,
		})
	}

	// Suggest line chart for time series data
	if len(numericColumns) >= 2 {
		suggestions = append(suggestions, gin.H{
			"type":        "line",
			"description": "Line chart showing trends in numeric data",
			"confidence":  0.7,
			"data_range":  dataRange,
		})
	}

	// Suggest pie chart for single series with categories
	if len(categoricalColumns) >= 1 && len(numericColumns) == 1 {
		suggestions = append(suggestions, gin.H{
			"type":        "pie",
			"description": "Pie chart showing distribution of " + numericColumns[0] + " by " + categoricalColumns[0],
			"confidence":  0.6,
			"data_range":  dataRange,
		})
	}

	// Suggest scatter plot for two numeric series
	if len(numericColumns) >= 2 {
		suggestions = append(suggestions, gin.H{
			"type":        "scatter",
			"description": "Scatter plot showing relationship between numeric variables",
			"confidence":  0.5,
			"data_range":  dataRange,
		})
	}

	return suggestions
}

// Collaborative editing functions

// GetSpreadsheetVersion gets current version of spreadsheet
func GetSpreadsheetVersion(c *gin.Context) {
	spreadsheetIDStr := c.Param("id")
	_, err := uuid.Parse(spreadsheetIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid spreadsheet ID"})
		return
	}

	// In a real implementation, get version from database
	version := gin.H{
		"version":    1,
		"updated_at": "2025-12-26T07:28:00Z",
		"updated_by": uuid.New(),
	}

	c.JSON(http.StatusOK, version)
}

// UpdateSpreadsheetBatch updates multiple cells in batch
func UpdateSpreadsheetBatch(c *gin.Context) {
	spreadsheetIDStr := c.Param("id")
	_, err := uuid.Parse(spreadsheetIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid spreadsheet ID"})
		return
	}

	var req struct {
		Updates []struct {
			CellReference string      `json:"cell_reference"`
			Value         interface{} `json:"value"`
			Formula       string      `json:"formula,omitempty"`
		} `json:"updates" binding:"required"`
		Version int `json:"version,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate all cell references
	for _, update := range req.Updates {
		if !isValidCellReference(update.CellReference) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cell reference: " + update.CellReference})
			return
		}
	}

	// In a real implementation, perform batch update with conflict resolution
	// For now, acknowledge the update
	newVersion := req.Version + 1

	c.JSON(http.StatusOK, gin.H{
		"message":      "Batch update successful",
		"updated_cells": len(req.Updates),
		"new_version": newVersion,
	})
}

// GetSpreadsheetChanges gets changes since a specific version
func GetSpreadsheetChanges(c *gin.Context) {
	spreadsheetIDStr := c.Param("id")
	_, err := uuid.Parse(spreadsheetIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid spreadsheet ID"})
		return
	}

	sinceVersionStr := c.Query("since_version")
	sinceVersion, _ := strconv.Atoi(sinceVersionStr)

	// In a real implementation, get changes from version history
	changes := []gin.H{
		{
			"version":       sinceVersion + 1,
			"cell_reference": "A1",
			"old_value":     "Old Product",
			"new_value":     "New Product",
			"changed_by":    uuid.New(),
			"changed_at":    "2025-12-26T07:28:00Z",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"changes": changes,
		"current_version": sinceVersion + 1,
	})
}

// LockSpreadsheetCells locks cells for editing (for collaboration)
func LockSpreadsheetCells(c *gin.Context) {
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

	spreadsheetIDStr := c.Param("id")
	_, err := uuid.Parse(spreadsheetIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid spreadsheet ID"})
		return
	}

	var req struct {
		CellReferences []string `json:"cell_references" binding:"required"`
		LockDuration   int      `json:"lock_duration,omitempty"` // seconds
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate cell references
	for _, ref := range req.CellReferences {
		if !isValidCellReference(ref) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cell reference: " + ref})
			return
		}
	}

	if req.LockDuration == 0 {
		req.LockDuration = 300 // 5 minutes default
	}

	// In a real implementation, acquire locks in database/cache
	c.JSON(http.StatusOK, gin.H{
		"message":         "Cells locked successfully",
		"locked_cells":    req.CellReferences,
		"locked_by":       userID,
		"lock_duration":   req.LockDuration,
		"expires_at":      "2025-12-26T07:33:00Z", // 5 minutes from now
	})
}

// UnlockSpreadsheetCells releases cell locks
func UnlockSpreadsheetCells(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	spreadsheetIDStr := c.Param("id")
	_, err := uuid.Parse(spreadsheetIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid spreadsheet ID"})
		return
	}

	var req struct {
		CellReferences []string `json:"cell_references" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// In a real implementation, release locks
	c.JSON(http.StatusOK, gin.H{
		"message":          "Cells unlocked successfully",
		"unlocked_cells":   req.CellReferences,
	})
}

// Excel Import/Export Functions

// ImportExcel imports data from an Excel file
func ImportExcel(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get uploaded file
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get uploaded file"})
		return
	}
	defer file.Close()

	// Create Excel service
	mathService := services.NewSpreadsheetMathService()
	excelService := services.NewExcelService(mathService)

	// Import Excel file
	result, err := excelService.ImportExcel(file, header.Filename)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ExportExcel exports spreadsheet data to Excel format
func ExportExcel(c *gin.Context) {
	spreadsheetIDStr := c.Param("id")
	spreadsheetID, err := uuid.Parse(spreadsheetIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid spreadsheet ID"})
		return
	}

	var req services.ExcelExportOptions
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// In a real implementation, get spreadsheet data from database
	// For now, use mock data
	sheets := []services.ExcelSheet{
		{
			Name: "Sheet1",
			Data: map[string]interface{}{
				"A1": "Product",
				"B1": "Price",
				"C1": "Quantity",
				"D1": "Total",
				"A2": "Widget A",
				"B2": 10.99,
				"C2": 5,
				"A3": "Widget B",
				"B3": 15.50,
				"C3": 3,
			},
			Formulas: map[string]string{
				"D2": "=B2*C2",
				"D3": "=B3*C3",
				"D4": "=SUM(D2:D3)",
			},
		},
	}

	// Create Excel service
	mathService := services.NewSpreadsheetMathService()
	excelService := services.NewExcelService(mathService)

	// Export to Excel
	excelData, err := excelService.ExportExcel(spreadsheetID, sheets, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set headers for file download
	filename := fmt.Sprintf("spreadsheet-%s.xlsx", spreadsheetID.String()[:8])
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Length", strconv.Itoa(len(excelData)))

	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", excelData)
}

// GetExcelTemplate returns a pre-built Excel template
func GetExcelTemplate(c *gin.Context) {
	templateType := c.Query("type")
	if templateType == "" {
		templateType = "basic"
	}

	// Create Excel service
	mathService := services.NewSpreadsheetMathService()
	excelService := services.NewExcelService(mathService)

	// Get template
	templateData, err := excelService.GetExcelTemplate(templateType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set headers for file download
	filename := fmt.Sprintf("%s-template.xlsx", templateType)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Length", strconv.Itoa(len(templateData)))

	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", templateData)
}

// GetExcelInfo returns information about an uploaded Excel file
func GetExcelInfo(c *gin.Context) {
	// Get uploaded file
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get uploaded file"})
		return
	}
	defer file.Close()

	// Create Excel service
	mathService := services.NewSpreadsheetMathService()
	excelService := services.NewExcelService(mathService)

	// Get file info
	info, err := excelService.GetExcelInfo(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, info)
}

// ValidateExcelFile validates if an uploaded file is a valid Excel file
func ValidateExcelFile(c *gin.Context) {
	// Get uploaded file
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get uploaded file"})
		return
	}
	defer file.Close()

	// Create Excel service
	mathService := services.NewSpreadsheetMathService()
	excelService := services.NewExcelService(mathService)

	// Validate file
	err = excelService.ValidateExcelFile(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"valid": false,
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"valid": true})
}

// GetSupportedExcelFormats returns supported Excel formats
func GetSupportedExcelFormats(c *gin.Context) {
	// Create Excel service
	mathService := services.NewSpreadsheetMathService()
	excelService := services.NewExcelService(mathService)

	formats := excelService.GetSupportedFormats()
	c.JSON(http.StatusOK, formats)
}
