package routes

import (
	"encoding/json"
	"net/http"

	"tpt-titan/backend/config"
	"tpt-titan/backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	dataJSON, _ := json.Marshal(req.Data)

	chart := models.SpreadsheetChart{
		ID:            chartID,
		SpreadsheetID: req.SpreadsheetID,
		ChartType:     req.ChartType,
		DataRange:     req.DataRange,
		Title:         req.Title,
		XAxisLabel:    req.XAxisLabel,
		YAxisLabel:    req.YAxisLabel,
		Data:          string(dataJSON),
	}

	if err := config.DB.Create(&chart).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save chart"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":             chart.ID,
		"spreadsheet_id": chart.SpreadsheetID,
		"type":           chart.ChartType,
		"data_range":     chart.DataRange,
		"title":          chart.Title,
		"x_axis_label":   chart.XAxisLabel,
		"y_axis_label":   chart.YAxisLabel,
		"data":           req.Data,
		"created_at":     chart.CreatedAt,
	})
}

// GetCharts retrieves charts for a spreadsheet
func GetCharts(c *gin.Context) {
	spreadsheetIDStr := c.Param("id")
	spreadsheetID, err := uuid.Parse(spreadsheetIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid spreadsheet ID"})
		return
	}

	var charts []models.SpreadsheetChart
	if err := config.DB.Where("spreadsheet_id = ?", spreadsheetID).Order("created_at DESC").Find(&charts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve charts"})
		return
	}

	result := make([]gin.H, 0, len(charts))
	for _, ch := range charts {
		var data map[string]interface{}
		_ = json.Unmarshal([]byte(ch.Data), &data)
		result = append(result, gin.H{
			"id":             ch.ID,
			"spreadsheet_id": ch.SpreadsheetID,
			"type":           ch.ChartType,
			"data_range":     ch.DataRange,
			"title":          ch.Title,
			"x_axis_label":   ch.XAxisLabel,
			"y_axis_label":   ch.YAxisLabel,
			"data":           data,
			"created_at":     ch.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"charts": result})
}

// analyzeDataForCharts analyzes spreadsheet data and suggests appropriate charts
func analyzeDataForCharts(data map[string]interface{}, dataRange string, dataTypes map[string]string) []gin.H {
	suggestions := []gin.H{}

	if len(data) == 0 {
		return suggestions
	}

	numericCols := 0
	categoricalCols := 0
	for _, t := range dataTypes {
		switch t {
		case "number", "float", "integer":
			numericCols++
		case "string", "category":
			categoricalCols++
		}
	}

	if categoricalCols > 0 && numericCols > 0 {
		suggestions = append(suggestions, gin.H{
			"type":        "bar",
			"title":       "Category Comparison",
			"data_range":  dataRange,
			"confidence":  0.9,
			"description": "Compare values across categories using a bar chart.",
		})
		suggestions = append(suggestions, gin.H{
			"type":        "pie",
			"title":       "Category Distribution",
			"data_range":  dataRange,
			"confidence":  0.7,
			"description": "Show the proportion of each category with a pie chart.",
		})
	}

	if numericCols >= 2 {
		suggestions = append(suggestions, gin.H{
			"type":        "line",
			"title":       "Trend Analysis",
			"data_range":  dataRange,
			"confidence":  0.8,
			"description": "Plot numeric series over an axis using a line chart.",
		})
		suggestions = append(suggestions, gin.H{
			"type":        "scatter",
			"title":       "Correlation",
			"data_range":  dataRange,
			"confidence":  0.6,
			"description": "Inspect correlation between two numeric series with a scatter plot.",
		})
	}

	return suggestions
}
