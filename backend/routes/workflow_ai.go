package routes

import (
	"net/http"
	"tpt-titan/backend/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var workflowAIService *services.WorkflowAIService

// InitWorkflowAIService initializes the workflow AI service (called from main)
func InitWorkflowAIService(aiService *services.AIService) {
	workflowAIService = services.NewWorkflowAIService(aiService)
}

// AnalyzeWorkflowUsage analyzes user workflow usage patterns
func AnalyzeWorkflowUsage(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	analysis, err := workflowAIService.AnalyzeUsagePatterns(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"analysis": analysis})
}

// GetSmartTemplateRecommendations provides personalized template suggestions
func GetSmartTemplateRecommendations(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get context from query parameters
	context := make(map[string]interface{})

	if recentForms := c.Query("recent_forms"); recentForms != "" {
		context["recent_form_submissions"] = recentForms
	}

	if emailVol := c.Query("email_volume"); emailVol != "" {
		context["email_volume"] = emailVol
	}

	if taskCount := c.Query("task_count"); taskCount != "" {
		context["task_creation"] = taskCount
	}

	recommendations, err := workflowAIService.GetSmartTemplateRecommendations(userID.(uuid.UUID), context)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"recommendations": recommendations})
}

// OptimizeWorkflow suggests optimizations for an existing workflow
func OptimizeWorkflow(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	workflowID := c.Param("workflowId")
	if workflowID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Workflow ID is required"})
		return
	}

	wfID, err := uuid.Parse(workflowID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workflow ID"})
		return
	}

	optimization, err := workflowAIService.OptimizeWorkflow(wfID, userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"optimization": optimization})
}

// PredictWorkflows suggests new workflows based on user behavior patterns
func PredictWorkflows(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	predictions, err := workflowAIService.PredictWorkflowCreation(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"predictions": predictions})
}

// GetWorkflowInsights provides comprehensive workflow analytics
func GetWorkflowInsights(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get usage analysis
	usageAnalysis, err := workflowAIService.AnalyzeUsagePatterns(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get template recommendations
	recommendations, err := workflowAIService.GetSmartTemplateRecommendations(userID.(uuid.UUID), nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get workflow predictions
	predictions, err := workflowAIService.PredictWorkflowCreation(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	insights := gin.H{
		"usage_analysis": usageAnalysis,
		"template_recommendations": recommendations,
		"workflow_predictions": predictions,
		"summary": gin.H{
			"total_workflows": len(usageAnalysis.Patterns),
			"active_workflows": func() int {
				count := 0
				for _, pattern := range usageAnalysis.Patterns {
					if pattern.ExecutionCount > 0 {
						count++
					}
				}
				return count
			}(),
			"total_executions": usageAnalysis.PerformanceMetrics.TotalExecutions,
			"success_rate": func() float64 {
				if usageAnalysis.PerformanceMetrics.TotalExecutions == 0 {
					return 0
				}
				return float64(usageAnalysis.PerformanceMetrics.SuccessfulExecutions) /
					   float64(usageAnalysis.PerformanceMetrics.TotalExecutions) * 100
			}(),
			"recommendations_count": len(usageAnalysis.Recommendations),
			"predictions_count": len(predictions),
		},
	}

	c.JSON(http.StatusOK, insights)
}
