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
