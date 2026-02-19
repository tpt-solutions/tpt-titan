package routes

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"tpt-titan/backend/services"
)

// CreateFormTemplate creates a new form template
func CreateFormTemplate(c *gin.Context) {
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

	var template struct {
		Name        string                   `json:"name" binding:"required"`
		Description string                   `json:"description"`
		Category    string                   `json:"category"`
		FormData    map[string]interface{}   `json:"form_data" binding:"required"`
		IsPublic    bool                     `json:"is_public"`
		Tags        []string                 `json:"tags"`
	}

	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// In a real implementation, save to form_templates table
	templateID := uuid.New()

	c.JSON(http.StatusCreated, gin.H{
		"id":          templateID,
		"name":        template.Name,
		"description": template.Description,
		"category":    template.Category,
		"created_by":  userID,
		"is_public":   template.IsPublic,
		"tags":        template.Tags,
	})
}

// GetFormTemplates gets available form templates
func GetFormTemplates(c *gin.Context) {
	category := c.Query("category")
	search := c.Query("search")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	// In a real implementation, query form_templates table
	templates := []gin.H{
		{
			"id":          uuid.New(),
			"name":        "Contact Information Form",
			"description": "Standard contact form with name, email, phone",
			"category":    "business",
			"is_public":   true,
			"tags":        []string{"contact", "business"},
		},
		{
			"id":          uuid.New(),
			"name":        "Survey Form",
			"description": "Multi-question survey with various field types",
			"category":    "survey",
			"is_public":   true,
			"tags":        []string{"survey", "feedback"},
		},
		{
			"id":          uuid.New(),
			"name":        "Invoice Form",
			"description": "Professional invoice with line items and totals",
			"category":    "finance",
			"is_public":   true,
			"tags":        []string{"invoice", "finance", "business"},
		},
	}

	// Filter by category and search
	if category != "" {
		filtered := []gin.H{}
		for _, t := range templates {
			if cat, ok := t["category"].(string); ok && cat == category {
				filtered = append(filtered, t)
			}
		}
		templates = filtered
	}

	if search != "" {
		filtered := []gin.H{}
		for _, t := range templates {
			if name, ok := t["name"].(string); ok &&
			   strings.Contains(strings.ToLower(name), strings.ToLower(search)) {
				filtered = append(filtered, t)
			}
		}
		templates = filtered
	}

	// Apply limit
	if len(templates) > limit {
		templates = templates[:limit]
	}

	c.JSON(http.StatusOK, gin.H{"templates": templates})
}

// GetFormTemplateCategories gets available template categories
func GetFormTemplateCategories(c *gin.Context) {
	categories := []gin.H{
		{"id": "business", "name": "Business Forms", "description": "Professional business forms"},
		{"id": "survey", "name": "Surveys", "description": "Feedback and survey forms"},
		{"id": "finance", "name": "Finance", "description": "Financial and accounting forms"},
		{"id": "hr", "name": "Human Resources", "description": "HR and personnel forms"},
		{"id": "education", "name": "Education", "description": "Educational forms and surveys"},
		{"id": "healthcare", "name": "Healthcare", "description": "Medical and healthcare forms"},
	}

	c.JSON(http.StatusOK, gin.H{"categories": categories})
}

// UseFormTemplate creates a new form from a template
func UseFormTemplate(c *gin.Context) {
	templateIDStr := c.Param("templateId")
	_, err := uuid.Parse(templateIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID"})
		return
	}

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
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// In a real implementation, create form from template
	formID := uuid.New()

	c.JSON(http.StatusCreated, gin.H{
		"id":          formID,
		"name":        req.Name,
		"description": req.Description,
		"template_id": templateIDStr,
		"created_by":  userID,
		"created_at":  "now",
	})
}
