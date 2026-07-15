package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"tpt-titan/backend/services"
)

// GetEquationTemplates returns available equation templates
func GetEquationTemplates(c *gin.Context) {
	category := c.Query("category")
	search := c.Query("search")

	// Create handwriting recognition service
	hrs := services.NewHandwritingRecognitionService("", "", "")

	templates, err := hrs.GetEquationTemplates(category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Filter by search query if provided
	if search != "" {
		searchResults, err := hrs.SearchEquations(search)
		if err == nil {
			templates = searchResults
		}
	}

	c.JSON(http.StatusOK, gin.H{"templates": templates})
}

// GetEquationTemplateCategories returns available template categories
func GetEquationTemplateCategories(c *gin.Context) {
	hrs := services.NewHandwritingRecognitionService("", "", "")
	categories := hrs.GetEquationTemplateCategories()

	c.JSON(http.StatusOK, gin.H{"categories": categories})
}

// SaveEquationTemplate saves a new equation template
func SaveEquationTemplate(c *gin.Context) {
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

	var template services.EquationTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	template.CreatedBy = userID

	hrs := services.NewHandwritingRecognitionService("", "", "")
	err := hrs.SaveEquationTemplate(&template)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, template)
}

// SearchEquations searches for equations
func SearchEquations(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	hrs := services.NewHandwritingRecognitionService("", "", "")
	results, err := hrs.SearchEquations(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"results": results})
}
