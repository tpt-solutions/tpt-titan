package routes

import (
	"encoding/json"
	"net/http"
	"tpt-titan/backend/config"
	"tpt-titan/backend/models"
	"tpt-titan/backend/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AIRequestPayload represents the payload for AI requests
type AIRequestPayload struct {
	TaskID    string `json:"task_id" binding:"required"`
	ModelID   string `json:"model_id" binding:"required"`
	Input     string `json:"input" binding:"required"`
	InputType string `json:"input_type" binding:"required"`
}

// AIModelPayload represents the payload for creating/updating AI models
type AIModelPayload struct {
	Name         string   `json:"name" binding:"required"`
	Type         string   `json:"type" binding:"required"`
	Provider     string   `json:"provider" binding:"required"`
	ModelID      string   `json:"model_id" binding:"required"`
	Capabilities []string `json:"capabilities"`
	APIKey       string   `json:"api_key,omitempty"`
	Endpoint     string   `json:"endpoint,omitempty"`
	Config       string   `json:"config,omitempty"`
}

var aiService *services.AIService
var aiConfig *config.AIConfig

// InitAIService initializes the AI service (called from main)
func InitAIService(cfg *config.AIConfig) {
	aiConfig = cfg
	aiService = services.NewAIService(cfg)
}

// GetAIService returns the initialized AI service instance
func GetAIService() *services.AIService {
	return aiService
}

// GetAIModels returns all available AI models for the user
func GetAIModels(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	models, err := aiService.GetAvailableModels(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve AI models"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"models": models})
}

// CreateAIModel creates a new custom AI model
func CreateAIModel(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var payload AIModelPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create the model
	model := &models.AIModel{
		UserID:       userID.(uuid.UUID),
		Name:         payload.Name,
		Type:         payload.Type,
		Provider:     payload.Provider,
		ModelID:      payload.ModelID,
		Capabilities: payload.Capabilities,
		APIKey:       []byte(payload.APIKey), // TODO: Encrypt this
		Endpoint:     payload.Endpoint,
		Config:       payload.Config,
		IsSystem:     false,
		IsActive:     true,
	}

	if err := config.DB.Create(model).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create AI model"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"model": model})
}

// ProcessAIRequest processes an AI request
func ProcessAIRequest(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var payload AIRequestPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse UUIDs
	taskID, err := uuid.Parse(payload.TaskID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	modelID, err := uuid.Parse(payload.ModelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid model ID"})
		return
	}

	// Process the request
	request, err := aiService.ProcessRequest(userID.(uuid.UUID), taskID, modelID, payload.Input, payload.InputType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"request_id": request.ID,
		"status":     request.Status,
		"message":    "AI request queued for processing",
	})
}

// GetAIRequestStatus gets the status of an AI request
func GetAIRequestStatus(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	requestID := c.Param("requestId")
	if requestID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request ID is required"})
		return
	}

	reqID, err := uuid.Parse(requestID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	var request models.AIRequest
	if err := config.DB.Where("id = ? AND user_id = ?", reqID, userID).First(&request).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Request not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"request": request,
		"status":  request.Status,
		"output":  request.Output,
		"error":   request.Error,
	})
}

// ListOllamaModels lists available models from Ollama
func ListOllamaModels(c *gin.Context) {
	client := services.NewOllamaService(aiConfig.OllamaHost, aiConfig.OllamaPort)
	models, err := client.ListModels()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list Ollama models: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"models": models})
}

// PullOllamaModel pulls a model from Ollama
func PullOllamaModel(c *gin.Context) {
	modelName := c.Param("modelName")
	if modelName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Model name is required"})
		return
	}

	client := services.NewOllamaService(aiConfig.OllamaHost, aiConfig.OllamaPort)
	err := client.PullModel(modelName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to pull model: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Model pulled successfully"})
}

// GetAIUsage returns usage statistics for the user
func GetAIUsage(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var usage []models.AIUsage
	if err := config.DB.Where("user_id = ?", userID).Order("date DESC").Limit(30).Find(&usage).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve usage statistics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"usage": usage})
}

// CheckForUpgrades initiates a manual upgrade check
func CheckForUpgrades(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	check, err := aiService.CheckForUpgrades(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check for upgrades: " + err.Error()})
		return
	}

	// Parse upgrade options for response
	var upgradeOptions []models.UpgradeOption
	json.Unmarshal([]byte(check.UpgradeOptions), &upgradeOptions)

	c.JSON(http.StatusOK, gin.H{
		"check_id":        check.ID,
		"checked_at":      check.CheckedAt,
		"hardware":        check.HardwareInfo,
		"upgrade_options": upgradeOptions,
	})
}

// GetUpgradeHistory returns past upgrade checks
func GetUpgradeHistory(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var checks []models.AIUpgradeCheck
	if err := config.DB.Where("user_id = ?", userID).Order("checked_at DESC").Limit(10).Find(&checks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve upgrade history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"checks": checks})
}

// ApplyUpgrade applies a selected upgrade option
func ApplyUpgrade(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var payload struct {
		UpgradeID string `json:"upgrade_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find the upgrade check
	var check models.AIUpgradeCheck
	if err := config.DB.Where("id = ? AND user_id = ?", payload.UpgradeID, userID).First(&check).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Upgrade check not found"})
		return
	}

	// Parse upgrade options
	var upgradeOptions []models.UpgradeOption
	json.Unmarshal([]byte(check.UpgradeOptions), &upgradeOptions)

	// For now, this is a placeholder - in production, this would:
	// 1. Pull new model from Ollama/OpenRouter
	// 2. Update task assignments
	// 3. Test the new model
	// 4. Provide rollback capability

	c.JSON(http.StatusAccepted, gin.H{
		"message": "Upgrade application initiated",
		"status":  "This is a placeholder - upgrade application not yet implemented",
	})
}

// DetectHardware detects and returns system hardware capabilities
func DetectHardware(c *gin.Context) {
	if _, exists := c.Get("user_id"); !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	hardware, err := aiService.DetectHardware()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to detect hardware: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"hardware": hardware})
}

// GetRecommendedModels returns model recommendations based on hardware
func GetRecommendedModels(c *gin.Context) {
	if _, exists := c.Get("user_id"); !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Detect hardware first
	hardware, err := aiService.DetectHardware()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to detect hardware: " + err.Error()})
		return
	}

	// Get recommendations
	recommendations := aiService.GetRecommendedModels(hardware)

	c.JSON(http.StatusOK, gin.H{
		"hardware":       hardware,
		"recommendations": recommendations,
	})
}

// SetupRecommendedModels sets up default tasks with recommended models
func SetupRecommendedModels(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Detect hardware
	hardware, err := aiService.DetectHardware()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to detect hardware: " + err.Error()})
		return
	}

	// Get recommendations
	recommendations := aiService.GetRecommendedModels(hardware)

	// Setup models and tasks
	if err := aiService.SetupRecommendedModels(userID.(uuid.UUID), recommendations); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to setup recommended models: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Recommended models and tasks have been set up",
		"recommendations": recommendations,
	})
}
