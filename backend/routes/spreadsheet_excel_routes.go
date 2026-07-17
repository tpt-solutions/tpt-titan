package routes

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"tpt-titan/backend/config"
	"tpt-titan/backend/services"
)

// ImportExcel imports data from an Excel file
func ImportExcel(c *gin.Context) {
	if _, exists := c.Get("user_id"); !exists {
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

	// Fetch the real spreadsheet cell data from the database.
	db, err := config.GetDatabase().DB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get database connection"})
		return
	}

	rows, err := db.Query(`
		SELECT cell_reference, value, formula, data_type
		FROM spreadsheet_cells WHERE spreadsheet_id = $1`, spreadsheetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cell data: " + err.Error()})
		return
	}
	defer rows.Close()

	data := map[string]interface{}{}
	formulas := map[string]string{}
	for rows.Next() {
		var cellRef, value, formula, dataType string
		if err := rows.Scan(&cellRef, &value, &formula, &dataType); err != nil {
			continue
		}
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

	sheets := []services.ExcelSheet{
		{
			Name:     "Sheet1",
			Data:     data,
			Formulas: formulas,
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
