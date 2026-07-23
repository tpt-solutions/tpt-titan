package routes

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"tpt-titan/backend/services"
)

// SaveMathCanvas saves a math canvas drawing
func SaveMathCanvas(c *gin.Context) {
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
		Name        string                       `json:"name" binding:"required"`
		Description string                       `json:"description"`
		Strokes     []services.HandwritingStroke `json:"strokes"`
		Recognized  *services.RecognitionResult  `json:"recognized,omitempty"`
		Width       int                          `json:"width"`
		Height      int                          `json:"height"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// In a real implementation, save to database
	canvasID := uuid.New()

	c.JSON(http.StatusCreated, gin.H{
		"id":          canvasID,
		"name":        req.Name,
		"description": req.Description,
		"user_id":     userID,
		"created_at":  "now",
	})
}

// GetMathCanvases gets user's saved math canvases
func GetMathCanvases(c *gin.Context) {
	if _, exists := c.Get("user_id"); !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// In a real implementation, retrieve from database
	canvases := []gin.H{
		{
			"id":          uuid.New(),
			"name":        "Sample Equation",
			"description": "A handwritten mathematical expression",
			"created_at":  "2025-12-26T09:00:00Z",
		},
	}

	c.JSON(http.StatusOK, gin.H{"canvases": canvases})
}

// GenerateEquationImage generates an image representation of an equation
func GenerateEquationImage(c *gin.Context) {
	var req struct {
		LaTeX  string `json:"latex" binding:"required"`
		Format string `json:"format"` // "png", "svg", "pdf"
		Size   string `json:"size"`   // "small", "medium", "large"
		Color  string `json:"color"`  // hex color
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Format == "" {
		req.Format = "png"
	}
	if req.Size == "" {
		req.Size = "medium"
	}

	// Render a real artifact of the requested format (SVG / PDF / PNG) containing
	// the equation text. A full TeX layout engine is out of scope, so the LaTeX
	// is rendered as text rather than faked typeset output.
	data, contentType, err := services.RenderEquationImage(req.LaTeX, req.Format, req.Size, req.Color)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filename := fmt.Sprintf("equation.%s", strings.ToLower(req.Format))
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", contentType)
	c.Header("Content-Length", strconv.Itoa(len(data)))
	c.Data(http.StatusOK, contentType, data)
}
