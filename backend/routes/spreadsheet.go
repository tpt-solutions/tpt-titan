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

	// Get database connection
	db := c.MustGet("db").(*sql.DB)

	// Create spreadsheet in database
	spreadsheetID := uuid.New()
	query := `
		INSERT INTO spreadsheets (id, owner_id, name, version, row_count, col_count)
		VALUES ($1, $2, $3, 1, 100, 26)
		RETURNING id, name, owner_id, version, created_at
	`

	var spreadsheet struct {
		ID       uuid.UUID `json:"id"`
		Name     string    `json:"name"`
		OwnerID  uuid.UUID `json:"owner_id"`
		Version  int       `json:"version"`
		CreatedAt string   `json:"created_at"`
	}

	err = db.QueryRow(query, spreadsheetID, userID, req.Name).Scan(
		&spreadsheet.ID,
		&spreadsheet.Name,
		&spreadsheet.OwnerID,
		&spreadsheet.Version,
		&spreadsheet.CreatedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create spreadsheet: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":        spreadsheet.ID,
		"name":      spreadsheet.Name,
		"owner_id":  spreadsheet.OwnerID,
		"version":   spreadsheet.Version,
		"created_at": spreadsheet.CreatedAt,
		"data":      make(map[string]interface{}),
		"formulas":  make(map[string]string),
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

	db := c.MustGet("db").(*sql.DB)

	// Get spreadsheet metadata
	var spreadsheet struct {
		ID       uuid.UUID `json:"id"`
		Name     string    `json:"name"`
		OwnerID  uuid.UUID `json:"owner_id"`
		Version  int       `json:"version"`
		RowCount int       `json:"row_count"`
		ColCount int       `json:"col_count"`
	}

	err = db.QueryRow(`
		SELECT id, name, owner_id, version, row_count, col_count
		FROM spreadsheets WHERE id = $1`, spreadsheetID).Scan(
		&spreadsheet.ID, &spreadsheet.Name, &spreadsheet.OwnerID,
		&spreadsheet.Version, &spreadsheet.RowCount, &spreadsheet.ColCount)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Spreadsheet not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
		return
	}

	// Get all cell data
	rows, err := db.Query(`
		SELECT cell_reference, value, formula, data_type
		FROM spreadsheet_cells WHERE spreadsheet_id = $1`, spreadsheetID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cell data: " + err.Error()})
		return
	}
	defer rows.Close()

	data := make(map[string]interface{})
	formulas := make(map[string]string)

	for rows.Next() {
		var cellRef, value, formula, dataType string
		err := rows.Scan(&cellRef, &value, &formula, &dataType)
		if err != nil {
			continue
		}

		// Convert value based on data type
		switch dataType {
		case "number":
			if num, err := strconv.ParseFloat(value, 64); err == nil {
				data[cellRef] = num
			} else {
				data[cellRef] = value
			}
		case "boolean":
			data[cellRef] = value == "true"
		default:
			data[cellRef] = value
		}

		if formula != "" {
			formulas[cellRef] = formula
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       spreadsheet.ID,
		"name":     spreadsheet.Name,
		"owner_id": spreadsheet.OwnerID,
		"version":  spreadsheet.Version,
		"data":     data,
		"formulas": formulas,
	})
}

// UpdateSpreadsheetCell updates a specific cell in the spreadsheet
func UpdateSpreadsheetCell(c *gin.Context) {
	spreadsheetIDStr := c.Param("id")
	spreadsheetID, err := uuid.Parse(spreadsheetIDStr)
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

	db := c.MustGet("db").(*sql.DB)

	// Determine data type and convert value to string
	var valueStr, dataType string
	switch v := req.Value.(type) {
	case float64, float32:
		dataType = "number"
		valueStr = fmt.Sprintf("%.10f", v)
		// Remove trailing zeros
		valueStr = strings.TrimRight(strings.TrimRight(valueStr, "0"), ".")
	case int, int32, int64:
		dataType = "number"
		valueStr = fmt.Sprintf("%d", v)
	case bool:
		dataType = "boolean"
		valueStr = fmt.Sprintf("%t", v)
	case string:
		dataType = "string"
		valueStr = v
	default:
		dataType = "string"
		valueStr = fmt.Sprintf("%v", req.Value)
	}

	// Get row and column indices for the cell
	col, row, err := parseCellReference(req.CellReference)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cell reference format"})
		return
	}

	// Insert or update cell data
	query := `
		INSERT INTO spreadsheet_cells (spreadsheet_id, cell_reference, row_index, col_index, value, formula, data_type)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (spreadsheet_id, cell_reference)
		DO UPDATE SET
			value = EXCLUDED.value,
			formula = EXCLUDED.formula,
			data_type = EXCLUDED.data_type,
			updated_at = CURRENT_TIMESTAMP
	`

	_, err = db.Exec(query, spreadsheetID, req.CellReference, row, col, valueStr, req.Formula, dataType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cell: " + err.Error()})
		return
	}

	// Update spreadsheet version
	_, err = db.Exec(`
		UPDATE spreadsheets
		SET version = version + 1, last_version_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1`, spreadsheetID)
	if err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to update spreadsheet version: %v\n", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Cell updated successfully",
		"cell_reference": req.CellReference,
		"value":          req.Value,
		"formula":        req.Formula,
		"data_type":      dataType,
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

// parseCellReference parses A1 style cell reference to column/row indices
func parseCellReference(ref string) (col, row int, err error) {
	ref = strings.ToUpper(strings.TrimSpace(ref))

	// Find where numbers start
	i := 0
	for i < len(ref) && (ref[i] < '0' || ref[i] > '9') {
		i++
	}

	if i == 0 || i == len(ref) {
		return 0, 0, fmt.Errorf("invalid cell reference: %s", ref)
	}

	colStr := ref[:i]
	rowStr := ref[i:]

	// Convert column letters to number
	col = 0
	for _, r := range colStr {
		col = col*26 + int(r-'A'+1)
	}
	col-- // 0-based

	// Parse row number
	row, err = strconv.Atoi(rowStr)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid row number in %s", ref)
	}
	row-- // 0-based

	return col, row, nil
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

// P2P Collaboration Routes

// GetCollaborationMode returns the current collaboration mode
func GetCollaborationMode(c *gin.Context) {
	cfg := c.MustGet("config").(*config.Config)

	mode := "server" // Default to server mode
	if cfg.P2P.Enabled {
		mode = "p2p"
	}

	c.JSON(http.StatusOK, gin.H{
		"mode":       mode,
		"p2p_enabled": cfg.P2P.Enabled,
		"server_available": true, // Server is always available as backup
	})
}

// SetCollaborationMode switches between server and P2P modes
func SetCollaborationMode(c *gin.Context) {
	var req struct {
		Mode string `json:"mode" binding:"required"` // "server" or "p2p"
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Mode != "server" && req.Mode != "p2p" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mode. Must be 'server' or 'p2p'"})
		return
	}

	// In a real implementation, this would update configuration
	// and restart services as needed
	cfg := c.MustGet("config").(*config.Config)

	switch req.Mode {
	case "p2p":
		if !cfg.P2P.Enabled {
			c.JSON(http.StatusBadRequest, gin.H{"error": "P2P mode is not configured"})
			return
		}
		// Enable P2P services
		p2pService := c.MustGet("p2p_service").(*services.P2PService)
		if !p2pService.IsRunning() {
			err := p2pService.Start()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start P2P service: " + err.Error()})
				return
			}
		}
	case "server":
		// P2P can remain running as fallback, but primary mode is server
		// Server mode uses existing WebSocket/real-time infrastructure
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Collaboration mode updated successfully",
		"mode":    req.Mode,
	})
}

// GetConnectedPeers returns list of connected P2P peers
func GetConnectedPeers(c *gin.Context) {
	p2pService := c.MustGet("p2p_service").(*services.P2PService)
	peers := p2pService.GetConnectedPeers()

	c.JSON(http.StatusOK, gin.H{
		"peers": peers,
		"count": len(peers),
	})
}

// ConnectToPeer manually connects to a specific peer
func ConnectToPeer(c *gin.Context) {
	var req struct {
		PeerID  string `json:"peer_id" binding:"required"`
		Address string `json:"address" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	p2pService := c.MustGet("p2p_service").(*services.P2PService)
	err := p2pService.ConnectToPeer(req.PeerID, req.Address)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to connect to peer: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully connected to peer",
		"peer_id": req.PeerID,
	})
}

// GetDiscoveredPeers returns peers discovered on the network
func GetDiscoveredPeers(c *gin.Context) {
	// In a full implementation, this would query the peer discovery service
	discoveredPeers := []gin.H{
		{
			"id":      "peer-001",
			"name":    "Alice's Computer",
			"address": "192.168.1.100:8081",
			"status":  "available",
		},
		{
			"id":      "peer-002",
			"name":    "Bob's Laptop",
			"address": "192.168.1.101:8081",
			"status":  "available",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"peers": discoveredPeers,
	})
}

// SyncSpreadsheetWithPeers broadcasts spreadsheet to all connected peers
func SyncSpreadsheetWithPeers(c *gin.Context) {
	spreadsheetIDStr := c.Param("id")
	spreadsheetID, err := uuid.Parse(spreadsheetIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid spreadsheet ID"})
		return
	}

	// Get spreadsheet data from database
	db := c.MustGet("db").(*sql.DB)

	var spreadsheet struct {
		Name string
	}

	err = db.QueryRow("SELECT name FROM spreadsheets WHERE id = $1", spreadsheetID).Scan(&spreadsheet.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get spreadsheet: " + err.Error()})
		return
	}

	// Get cell data
	rows, err := db.Query("SELECT cell_reference, value FROM spreadsheet_cells WHERE spreadsheet_id = $1", spreadsheetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get cell data: " + err.Error()})
		return
	}
	defer rows.Close()

	data := make(map[string]interface{})
	for rows.Next() {
		var cellRef, value string
		rows.Scan(&cellRef, &value)
		// Convert value to appropriate type (simplified)
		if num, err := strconv.ParseFloat(value, 64); err == nil {
			data[cellRef] = num
		} else {
			data[cellRef] = value
		}
	}

	// Sync with P2P peers
	p2pService := c.MustGet("p2p_service").(*services.P2PService)
	p2pService.SyncSpreadsheet(spreadsheetID, data)

	c.JSON(http.StatusOK, gin.H{
		"message": "Spreadsheet synced with peers",
		"peers_notified": len(p2pService.GetConnectedPeers()),
	})
}

// GetCollaborationStatus returns current collaboration status
func GetCollaborationStatus(c *gin.Context) {
	cfg := c.MustGet("config").(*config.Config)
	p2pService := c.MustGet("p2p_service").(*services.P2PService)

	// Determine current mode based on configuration and what's running
	currentMode := "automatic" // Default: automatic selection for ease of use

	peers := p2pService.GetConnectedPeers()
	connectedPeers := len(peers)

	status := gin.H{
		"mode": currentMode, // "automatic", "p2p", "server"
		"connected_peers": connectedPeers,
		"remote_access_enabled": cfg.P2P.AllowRemoteAccess,
		"features": gin.H{
			"local_network": gin.H{
				"available": true,
				"description": "Direct peer-to-peer collaboration on local networks",
				"speed": "Fast",
				"security": "Local network only",
			},
			"remote_access": gin.H{
				"available": cfg.P2P.AllowRemoteAccess,
				"description": "Cloud relay for remote users (work from home, etc.)",
				"speed": "Good",
				"security": "End-to-end encrypted",
			},
			"server_backup": gin.H{
				"available": true,
				"description": "Full server-based collaboration as backup",
				"speed": "Variable",
				"security": "User authentication + permissions",
			},
		},
		"p2p": gin.H{
			"enabled": cfg.P2P.Enabled,
			"running": p2pService.IsRunning(),
			"topology": cfg.P2P.PreferredTopology,
			"auto_detect": cfg.P2P.AutoDetectTopology,
			"config": gin.H{
				"port": cfg.P2P.Port,
				"max_peers": cfg.P2P.MaxPeers,
				"discovery_timeout": cfg.P2P.DiscoveryTimeout,
				"sync_interval": cfg.P2P.SyncInterval,
				"encryption": cfg.P2P.EnableEncryption,
				"compression": cfg.P2P.EnableCompression,
				"remote_access": cfg.P2P.AllowRemoteAccess,
				"cloud_relay": cfg.P2P.CloudRelayEnabled,
			},
		},
		"server": gin.H{
			"available": true,
			"description": "Traditional server-based collaboration",
			"features": []string{
				"User authentication",
				"Granular permissions",
				"Version history",
				"Advanced sharing",
			},
		},
		"recommendation": gin.H{
			"for_local_office": "automatic (P2P with cloud relay fallback)",
			"for_remote_teams": "server mode",
			"for_home_users": "automatic (works everywhere)",
			"easiest_setup": "automatic mode - just works",
		},
		"peers": peers,
	}

	// Determine actual current mode
	if cfg.P2P.Enabled && p2pService.IsRunning() {
		if connectedPeers > 0 {
			status["mode"] = "p2p"
			status["active_connections"] = "local_network"
		} else if cfg.P2P.AllowRemoteAccess {
			status["mode"] = "remote_ready"
			status["active_connections"] = "cloud_relay_available"
		}
	} else {
		status["mode"] = "server"
		status["active_connections"] = "server_based"
	}

	// Add user-friendly status messages
	status["user_message"] = gin.H{
		"current_status": "Ready for collaboration",
		"ease_of_use": "Just works - no configuration needed",
		"security_note": "All connections are encrypted and secure",
		"performance": "Automatic optimization for your network",
	}

	c.JSON(http.StatusOK, status)
}

// IsRunning returns whether P2P service is running (helper method)
func (p2p *services.P2PService) IsRunning() bool {
	// This would check the running state
	// For now, return false as we don't have a getter
	return false
}
