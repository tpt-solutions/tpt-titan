package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateEmailDistribution creates an email distribution for form responses
func CreateEmailDistribution(c *gin.Context) {
	formIDStr := c.Param("formId")
	_, err := uuid.Parse(formIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form ID"})
		return
	}

	var req struct {
		Name         string   `json:"name" binding:"required"`
		Recipients   []string `json:"recipients" binding:"required"`
		Subject      string   `json:"subject" binding:"required"`
		Message      string   `json:"message"`
		Trigger      string   `json:"trigger"` // "immediate", "daily", "weekly"
		IsActive     bool     `json:"is_active"`
		IncludeData  bool     `json:"include_data"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// In a real implementation, save to form_email_distributions table
	distributionID := uuid.New()

	c.JSON(http.StatusCreated, gin.H{
		"id":            distributionID,
		"form_id":       formIDStr,
		"name":          req.Name,
		"recipients":    req.Recipients,
		"subject":       req.Subject,
		"trigger":       req.Trigger,
		"is_active":     req.IsActive,
		"include_data":  req.IncludeData,
	})
}

// GetEmailDistributions gets email distributions for a form
func GetEmailDistributions(c *gin.Context) {
	formIDStr := c.Param("formId")
	_, err := uuid.Parse(formIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form ID"})
		return
	}

	// In a real implementation, query form_email_distributions table
	distributions := []gin.H{
		{
			"id":         uuid.New(),
			"name":       "Daily Summary",
			"recipients": []string{"admin@example.com"},
			"trigger":    "daily",
			"is_active":  true,
		},
	}

	c.JSON(http.StatusOK, gin.H{"distributions": distributions})
}

// SendFormResponseEmail sends an email with form response data
func SendFormResponseEmail(c *gin.Context) {
	responseIDStr := c.Param("responseId")
	_, err := uuid.Parse(responseIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid response ID"})
		return
	}

	var req struct {
		Recipients []string `json:"recipients" binding:"required"`
		Subject    string   `json:"subject" binding:"required"`
		Message    string   `json:"message"`
		IncludeData bool    `json:"include_data"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// In a real implementation, send email with form data
	c.JSON(http.StatusOK, gin.H{
		"message":    "Email sent successfully",
		"recipients": req.Recipients,
		"subject":    req.Subject,
	})
}
