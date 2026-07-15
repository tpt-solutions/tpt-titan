package routes

import (
	"encoding/base64"
	"io"
	"net/http"
	"tpt-titan/backend/config"
	"tpt-titan/backend/models"
	"tpt-titan/backend/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)




// SpeechRequestPayload represents the payload for TTS requests
type SpeechRequestPayload struct {
	Text     string  `json:"text" binding:"required"`
	ModelID  string  `json:"model_id" binding:"required"`
	Voice    string  `json:"voice,omitempty"`
	Language string  `json:"language,omitempty"`
	Speed    float64 `json:"speed,omitempty"`
	Pitch    float64 `json:"pitch,omitempty"`
	Format   string  `json:"format,omitempty"`
}

// SpeechUploadPayload represents the payload for STT requests with audio upload
type SpeechUploadPayload struct {
	AudioData []byte `json:"-"` // Set from multipart form
	ModelID   string `json:"model_id" binding:"required"`
	Language  string `json:"language,omitempty"`
	Format    string `json:"format,omitempty"`
}

var speechService *services.SpeechService

// InitSpeechService initializes the speech service (called from main)
func InitSpeechService(cfg *config.SpeechConfig) {
	speechService = services.NewSpeechService(cfg)
}

// GetSpeechModels returns all available speech models for the user
func GetSpeechModels(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	modelType := c.Query("type") // "tts" or "stt"
	if modelType == "" {
		modelType = "tts" // Default to TTS
	}

	models, err := speechService.GetAvailableModels(userID.(uuid.UUID), modelType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve speech models"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"models": models})
}

// CreateSpeechModel creates a new custom speech model
func CreateSpeechModel(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var payload struct {
		Name     string `json:"name" binding:"required"`
		Provider string `json:"provider" binding:"required"`
		ModelID  string `json:"model_id" binding:"required"`
		Type     string `json:"type" binding:"required"`
		Language string `json:"language,omitempty"`
		Voice    string `json:"voice,omitempty"`
		APIKey   string `json:"api_key,omitempty"`
		Endpoint string `json:"endpoint,omitempty"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create the model
	model := &models.SpeechModel{
		UserID:   userID.(uuid.UUID),
		Name:     payload.Name,
		Provider: models.SpeechProvider(payload.Provider),
		ModelID:  payload.ModelID,
		Type:     payload.Type,
		Language: payload.Language,
		Voice:    payload.Voice,
		APIKey:   []byte(payload.APIKey), // TODO: Encrypt this
		Endpoint: payload.Endpoint,
		IsSystem: false,
		IsActive: true,
	}

	if err := config.DB.Create(model).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create speech model"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"model": model})
}

// TextToSpeech converts text to speech
func TextToSpeech(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var payload SpeechRequestPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse model ID
	modelID, err := uuid.Parse(payload.ModelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid model ID"})
		return
	}

	// Set default options
	options := services.SpeechOptions{
		Language:    payload.Language,
		Voice:       payload.Voice,
		AudioFormat: payload.Format,
		Speed:       payload.Speed,
		Pitch:       payload.Pitch,
	}

	if options.Language == "" {
		options.Language = "en"
	}
	if options.Voice == "" {
		options.Voice = "alloy"
	}
	if options.AudioFormat == "" {
		options.AudioFormat = "mp3"
	}
	if options.Speed == 0 {
		options.Speed = 1.0
	}
	if options.Pitch == 0 {
		options.Pitch = 1.0
	}

	// Process TTS
	request, err := speechService.TextToSpeech(userID.(uuid.UUID), payload.Text, modelID, options)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"request_id": request.ID,
		"status":     request.Status,
		"message":    "TTS request queued for processing",
	})
}

// SpeechToText converts speech to text
func SpeechToText(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Handle multipart form upload
	file, _, err := c.Request.FormFile("audio")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Audio file is required"})
		return
	}
	defer file.Close()

	// Read audio data
	audioData, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read audio file"})
		return
	}


	// Get form parameters
	modelIDStr := c.PostForm("model_id")
	language := c.PostForm("language")
	format := c.PostForm("format")

	if modelIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Model ID is required"})
		return
	}

	modelID, err := uuid.Parse(modelIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid model ID"})
		return
	}

	// Set default options
	options := services.SpeechOptions{
		Language:    language,
		AudioFormat: format,
	}

	if options.Language == "" {
		options.Language = "en"
	}
	if options.AudioFormat == "" {
		options.AudioFormat = "wav"
	}

	// Process STT
	request, err := speechService.SpeechToText(userID.(uuid.UUID), audioData, modelID, options)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"request_id": request.ID,
		"status":     request.Status,
		"message":    "STT request queued for processing",
	})
}

// GetSpeechRequestStatus gets the status of a speech request
func GetSpeechRequestStatus(c *gin.Context) {
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

	var request models.SpeechRequest
	if err := config.DB.Where("id = ? AND user_id = ?", reqID, userID).First(&request).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Request not found"})
		return
	}

	response := gin.H{
		"request":        request,
		"status":         request.Status,
		"error":          request.Error,
		"processing_time": request.ProcessingTime,
	}

	// If TTS completed, provide audio data as base64
	if request.RequestType == "tts" && request.Status == "completed" && len(request.OutputAudio) > 0 {
		response["audio_data"] = base64.StdEncoding.EncodeToString(request.OutputAudio)
		response["audio_format"] = request.AudioFormat
	}

	// If STT completed, provide transcript
	if request.RequestType == "stt" && request.Status == "completed" && request.OutputText != "" {
		response["transcript"] = request.OutputText
	}

	c.JSON(http.StatusOK, response)
}

// GetSpeechSettings returns user's speech preferences
func GetSpeechSettings(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	settings, err := speechService.GetSpeechSettings(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve speech settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"settings": settings})
}

// UpdateSpeechSettings updates user's speech preferences
func UpdateSpeechSettings(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var settings models.SpeechSettings
	if err := c.ShouldBindJSON(&settings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := speechService.UpdateSpeechSettings(userID.(uuid.UUID), &settings); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update speech settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Speech settings updated successfully"})
}

// GetSpeechHistory returns user's speech processing history
func GetSpeechHistory(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var requests []models.SpeechRequest
	query := config.DB.Where("user_id = ?", userID).Order("created_at DESC")

	// Optional filtering
	requestType := c.Query("type")
	if requestType != "" {
		query = query.Where("request_type = ?", requestType)
	}

	status := c.Query("status")
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Pagination
	limit := 50 // Default limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if l := parseInt(limitStr); l > 0 && l <= 100 {
			limit = l
		}
	}

	if err := query.Limit(limit).Find(&requests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve speech history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"requests": requests})
}
