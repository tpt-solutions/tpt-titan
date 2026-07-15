package routes

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"tpt-titan/backend/services"
)

// ExportEquation exports an equation to different formats
func ExportEquation(c *gin.Context) {
	var req struct {
		Expression services.MathExpression `json:"expression" binding:"required"`
		Format     string                  `json:"format" binding:"required"` // "latex", "mathml", "svg", "png", "pdf"
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hrs := services.NewHandwritingRecognitionService("", "", "")
	data, err := hrs.ExportEquation(&req.Expression, req.Format)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set appropriate headers based on format
	var contentType string
	var filename string

	switch req.Format {
	case "latex":
		contentType = "text/plain"
		filename = "equation.tex"
	case "mathml":
		contentType = "application/xml"
		filename = "equation.xml"
	case "svg":
		contentType = "image/svg+xml"
		filename = "equation.svg"
	case "png":
		contentType = "image/png"
		filename = "equation.png"
	case "pdf":
		contentType = "application/pdf"
		filename = "equation.pdf"
	default:
		contentType = "application/octet-stream"
		filename = "equation.dat"
	}

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", contentType)
	c.Header("Content-Length", strconv.Itoa(len(data)))
	c.Data(http.StatusOK, contentType, data)
}

// BatchExportEquations exports multiple equations
func BatchExportEquations(c *gin.Context) {
	var req struct {
		Expressions []services.MathExpression `json:"expressions" binding:"required"`
		Format      string                    `json:"format" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hrs := services.NewHandwritingRecognitionService("", "", "")

	results := make(map[string][]byte)
	for i, expr := range req.Expressions {
		data, err := hrs.ExportEquation(&expr, req.Format)
		if err != nil {
			// Skip failed exports
			continue
		}
		results[fmt.Sprintf("equation_%d", i+1)] = data
	}

	c.JSON(http.StatusOK, gin.H{
		"exports": results,
		"format":  req.Format,
		"count":   len(results),
	})
}
