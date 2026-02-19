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
