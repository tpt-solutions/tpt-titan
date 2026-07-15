package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"tpt-titan/backend/services"
)

// RecognizeHandwriting recognizes handwritten mathematical expressions
func RecognizeHandwriting(c *gin.Context) {
	var req struct {
		Strokes []services.HandwritingStroke `json:"strokes" binding:"required"`
		Width   int                          `json:"width,omitempty"`
		Height  int                          `json:"height,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create handwriting recognition service
	// In a real implementation, this would use actual ML service credentials
	hrs := services.NewHandwritingRecognitionService("", "https://api.example.com", "math-model")

	result, err := hrs.RecognizeHandwriting(req.Strokes, req.Width, req.Height)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// RecognizeEquationFromImage recognizes equations from uploaded images
func RecognizeEquationFromImage(c *gin.Context) {
	// Get uploaded image file
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get uploaded image"})
		return
	}
	defer file.Close()

	// Read image data
	imageData := make([]byte, header.Size)
	_, err = file.Read(imageData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read image data"})
		return
	}

	// Get image format (optional parameter)
	format := c.PostForm("format")
	if format == "" {
		format = "png" // default
	}

	// Create handwriting recognition service
	hrs := services.NewHandwritingRecognitionService("", "https://api.example.com", "math-model")

	result, err := hrs.RecognizeEquationFromImage(imageData, format)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
