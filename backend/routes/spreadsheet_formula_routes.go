package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"tpt-titan/backend/services"
)

// EvaluateFormulaRequest is the payload for evaluating a spreadsheet formula
type EvaluateFormulaRequest struct {
	Formula     string                 `json:"formula" binding:"required"`
	CellContext map[string]interface{} `json:"cell_context,omitempty"`
}

// EvaluateFormulaResponse is the result of evaluating a spreadsheet formula
type EvaluateFormulaResponse struct {
	Result      interface{} `json:"result"`
	DataType    string      `json:"data_type"`
	Error       string      `json:"error,omitempty"`
	DependsOn   []string    `json:"depends_on,omitempty"`
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
