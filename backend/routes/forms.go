package routes

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Form represents a basic form structure
type Form struct {
	ID          uuid.UUID `json:"id" db:"id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Fields      []Field   `json:"fields" db:"fields"`
	Status      string    `json:"status" db:"status"`
	CreatedAt   string    `json:"created_at" db:"created_at"`
	UpdatedAt   string    `json:"updated_at" db:"updated_at"`
}

// Field represents a form field
type Field struct {
	ID          uuid.UUID `json:"id" db:"id"`
	FormID      uuid.UUID `json:"form_id" db:"form_id"`
	Type        string    `json:"type" db:"type"`
	Label       string    `json:"label" db:"label"`
	Placeholder string    `json:"placeholder" db:"placeholder"`
	Required    bool      `json:"required" db:"required"`
	Options     []string  `json:"options" db:"options"`
	Order       int       `json:"order" db:"order"`
}

// GetForms returns all forms for the authenticated user
func GetForms(c *gin.Context) {
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

	db := c.MustGet("db").(*sql.DB)

	// For now, return mock data until proper database schema is set up
	mockForms := []Form{
		{
			ID:          uuid.New(),
			UserID:      userID,
			Name:        "Customer Feedback Survey",
			Description: "Collect customer satisfaction data",
			Status:      "active",
			CreatedAt:   "2024-01-15T10:00:00Z",
			Fields: []Field{
				{
					ID:       uuid.New(),
					Type:     "text",
					Label:    "Name",
					Required: true,
					Order:    1,
				},
				{
					ID:       uuid.New(),
					Type:     "email",
					Label:    "Email",
					Required: true,
					Order:    2,
				},
			},
		},
		{
			ID:          uuid.New(),
			UserID:      userID,
			Name:        "Event Registration",
			Description: "Register attendees for company events",
			Status:      "active",
			CreatedAt:   "2024-01-10T10:00:00Z",
			Fields: []Field{
				{
					ID:       uuid.New(),
					Type:     "text",
					Label:    "Full Name",
					Required: true,
					Order:    1,
				},
				{
					ID:          uuid.New(),
					Type:        "select",
					Label:       "Event Type",
					Required:    true,
					Options:     []string{"Conference", "Workshop", "Webinar"},
					Order:       2,
				},
			},
		},
	}

	c.JSON(http.StatusOK, gin.H{"forms": mockForms})
}

// GetForm returns a specific form
func GetForm(c *gin.Context) {
	formIDStr := c.Param("id")
	formID, err := uuid.Parse(formIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form ID"})
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

	// Mock form data
	form := Form{
		ID:          formID,
		UserID:      userID,
		Name:        "Sample Form",
		Description: "A sample form for testing",
		Status:      "active",
		CreatedAt:   "2024-01-15T10:00:00Z",
		Fields: []Field{
			{
				ID:       uuid.New(),
				FormID:   formID,
				Type:     "text",
				Label:    "Name",
				Required: true,
				Order:    1,
			},
			{
				ID:       uuid.New(),
				FormID:   formID,
				Type:     "email",
				Label:    "Email",
				Required: true,
				Order:    2,
			},
		},
	}

	c.JSON(http.StatusOK, form)
}

// CreateForm creates a new form
func CreateForm(c *gin.Context) {
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

	var formData struct {
		Name        string  `json:"name" binding:"required"`
		Description string  `json:"description"`
		Fields      []Field `json:"fields"`
	}

	if err := c.ShouldBindJSON(&formData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create new form
	form := Form{
		ID:          uuid.New(),
		UserID:      userID,
		Name:        formData.Name,
		Description: formData.Description,
		Fields:      formData.Fields,
		Status:      "draft",
		CreatedAt:   "now",
		UpdatedAt:   "now",
	}

	c.JSON(http.StatusCreated, form)
}

// UpdateForm updates an existing form
func UpdateForm(c *gin.Context) {
	formIDStr := c.Param("id")
	_, err := uuid.Parse(formIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form ID"})
		return
	}

	var formData struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Fields      []Field `json:"fields"`
		Status      string  `json:"status"`
	}

	if err := c.ShouldBindJSON(&formData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Mock updated form
	form := Form{
		ID:          uuid.MustParse(formIDStr),
		Name:        formData.Name,
		Description: formData.Description,
		Fields:      formData.Fields,
		Status:      formData.Status,
		UpdatedAt:   "now",
	}

	c.JSON(http.StatusOK, form)
}

// DeleteForm deletes a form
func DeleteForm(c *gin.Context) {
	formIDStr := c.Param("id")
	_, err := uuid.Parse(formIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form ID"})
		return
	}

	// In a real implementation, delete from database
	c.JSON(http.StatusOK, gin.H{"message": "Form deleted successfully"})
}

// GetFormResponses gets responses for a form
func GetFormResponses(c *gin.Context) {
	formIDStr := c.Param("id")
	_, err := uuid.Parse(formIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form ID"})
		return
	}

	// Mock responses data
	responses := []map[string]interface{}{
		{
			"id":         uuid.New(),
			"form_id":    formIDStr,
			"responses":  map[string]interface{}{"name": "John Doe", "email": "john@example.com"},
			"submitted_at": "2024-01-16T10:00:00Z",
		},
		{
			"id":         uuid.New(),
			"form_id":    formIDStr,
			"responses":  map[string]interface{}{"name": "Jane Smith", "email": "jane@example.com"},
			"submitted_at": "2024-01-17T10:00:00Z",
		},
	}

	c.JSON(http.StatusOK, gin.H{"responses": responses})
}

// SubmitFormResponse submits a response to a form
func SubmitFormResponse(c *gin.Context) {
	formIDStr := c.Param("id")
	_, err := uuid.Parse(formIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form ID"})
		return
	}

	var responseData map[string]interface{}
	if err := c.ShouldBindJSON(&responseData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Mock response submission
	response := map[string]interface{}{
		"id":          uuid.New(),
		"form_id":     formIDStr,
		"responses":   responseData,
		"submitted_at": "now",
	}

	c.JSON(http.StatusCreated, response)
}
