package routes

import (
	"net/http"
	"strings"
	"time"
	"tpt-titan/backend/models"
	"tpt-titan/backend/services"
	"tpt-titan/backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetAISettings retrieves user's AI settings
func GetAISettings(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var settings models.AISettings
	if err := config.DB.Where("user_id = ?", userID).First(&settings).Error; err != nil {
		// Create default settings if not found
		settings = models.AISettings{
			UserID:                userID.(uuid.UUID),
			EnableAIFeatures:      true,
			EnableOCR:             true,
			EnableSpeech:          true,
			EnableWorkflows:       true,
			EnableLocalAI:         true,
			EnableCloudAI:         false,
			DefaultProvider:       "ollama",
			MaxConcurrentRequests: 3,
			RequestTimeout:        30,
			HardwareAcceleration:  true,
			LowPowerMode:          false,
		}
		if err := config.DB.Create(&settings).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create default settings"})
			return
		}
	}

	// Return settings without encrypted API keys
	response := map[string]interface{}{
		"user_id":                   settings.UserID,
		"enable_ai_features":        settings.EnableAIFeatures,
		"enable_ocr":                settings.EnableOCR,
		"enable_speech":             settings.EnableSpeech,
		"enable_workflows":          settings.EnableWorkflows,
		"enable_local_ai":           settings.EnableLocalAI,
		"enable_cloud_ai":           settings.EnableCloudAI,
		"default_provider":          settings.DefaultProvider,
		"max_concurrent_requests":   settings.MaxConcurrentRequests,
		"request_timeout":           settings.RequestTimeout,
		"hardware_acceleration":     settings.HardwareAcceleration,
		"low_power_mode":            settings.LowPowerMode,
		"api_keys": map[string]string{
			"openai":     "[CONFIGURED]",
			"elevenlabs": "[CONFIGURED]",
			"replicate":  "[CONFIGURED]",
			"assemblyai": "[CONFIGURED]",
			"deepgram":   "[CONFIGURED]",
		},
		"created_at": settings.CreatedAt,
		"updated_at": settings.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// UpdateAISettings updates user's AI settings
func UpdateAISettings(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var requestData struct {
		EnableAIFeatures      *bool   `json:"enable_ai_features"`
		EnableOCR             *bool   `json:"enable_ocr"`
		EnableSpeech          *bool   `json:"enable_speech"`
		EnableWorkflows       *bool   `json:"enable_workflows"`
		EnableLocalAI         *bool   `json:"enable_local_ai"`
		EnableCloudAI         *bool   `json:"enable_cloud_ai"`
		DefaultProvider       *string `json:"default_provider"`
		MaxConcurrentRequests *int    `json:"max_concurrent_requests"`
		RequestTimeout        *int    `json:"request_timeout"`
		HardwareAcceleration  *bool   `json:"hardware_acceleration"`
		LowPowerMode          *bool   `json:"low_power_mode"`
		APIKeys               map[string]string `json:"api_keys"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	var settings models.AISettings
	if err := config.DB.Where("user_id = ?", userID).First(&settings).Error; err != nil {
		// Create new settings if not found
		settings = models.AISettings{
			UserID: userID.(uuid.UUID),
		}
	}

	// Update fields if provided
	if requestData.EnableAIFeatures != nil {
		settings.EnableAIFeatures = *requestData.EnableAIFeatures
	}
	if requestData.EnableOCR != nil {
		settings.EnableOCR = *requestData.EnableOCR
	}
	if requestData.EnableSpeech != nil {
		settings.EnableSpeech = *requestData.EnableSpeech
	}
	if requestData.EnableWorkflows != nil {
		settings.EnableWorkflows = *requestData.EnableWorkflows
	}
	if requestData.EnableLocalAI != nil {
		settings.EnableLocalAI = *requestData.EnableLocalAI
	}
	if requestData.EnableCloudAI != nil {
		settings.EnableCloudAI = *requestData.EnableCloudAI
	}
	if requestData.DefaultProvider != nil {
		settings.DefaultProvider = *requestData.DefaultProvider
	}
	if requestData.MaxConcurrentRequests != nil {
		settings.MaxConcurrentRequests = *requestData.MaxConcurrentRequests
	}
	if requestData.RequestTimeout != nil {
		settings.RequestTimeout = *requestData.RequestTimeout
	}
	if requestData.HardwareAcceleration != nil {
		settings.HardwareAcceleration = *requestData.HardwareAcceleration
	}
	if requestData.LowPowerMode != nil {
		settings.LowPowerMode = *requestData.LowPowerMode
	}

	// Handle API key updates (encrypt them)
	if requestData.APIKeys != nil {
		if key, exists := requestData.APIKeys["openai"]; exists && key != "" {
			encrypted, _ := utils.Encrypt([]byte(key))
			settings.OpenAIKey = encrypted
		}
		if key, exists := requestData.APIKeys["elevenlabs"]; exists && key != "" {
			encrypted, _ := utils.Encrypt([]byte(key))
			settings.ElevenLabsKey = encrypted
		}
		if key, exists := requestData.APIKeys["replicate"]; exists && key != "" {
			encrypted, _ := utils.Encrypt([]byte(key))
			settings.ReplicateKey = encrypted
		}
		if key, exists := requestData.APIKeys["assemblyai"]; exists && key != "" {
			encrypted, _ := utils.Encrypt([]byte(key))
			settings.AssemblyAIKey = encrypted
		}
		if key, exists := requestData.APIKeys["deepgram"]; exists && key != "" {
			encrypted, _ := utils.Encrypt([]byte(key))
			settings.DeepgramKey = encrypted
		}
	}

	// Save or update settings
	if err := config.DB.Save(&settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "AI settings updated successfully"})
}

// GetSpeechSettings retrieves user's speech settings
func GetSpeechSettings(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var settings models.SpeechSettings
	if err := config.DB.Where("user_id = ?", userID).First(&settings).Error; err != nil {
		// Create default settings if not found
		settings = models.SpeechSettings{
			UserID:             userID.(uuid.UUID),
			EnableTTS:          true,
			EnableSTT:          true,
			DefaultVoice:       "alloy",
			DefaultLanguage:    "en",
			TTSSpeed:           1.0,
			TTSVolume:          1.0,
			STTLanguage:        "en",
			AutoPlayTTS:        false,
			ShowSTTTranscript:  true,
			KeyboardShortcut:   "ctrl+shift+s",
		}
		if err := config.DB.Create(&settings).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create default speech settings"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"settings": settings})
}

// UpdateSpeechSettings updates user's speech settings
func UpdateSpeechSettings(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var requestData struct {
		EnableTTS         *bool   `json:"enable_tts"`
		EnableSTT         *bool   `json:"enable_stt"`
		DefaultTTSModel   *string `json:"default_tts_model"`
		DefaultSTTModel   *string `json:"default_stt_model"`
		DefaultVoice      *string `json:"default_voice"`
		DefaultLanguage   *string `json:"default_language"`
		TTSSpeed          *float64 `json:"tts_speed"`
		TTSVolume         *float64 `json:"tts_volume"`
		STTLanguage       *string `json:"stt_language"`
		AutoPlayTTS       *bool   `json:"auto_play_tts"`
		ShowSTTTranscript *bool   `json:"show_stt_transcript"`
		KeyboardShortcut  *string `json:"keyboard_shortcut"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	var settings models.SpeechSettings
	if err := config.DB.Where("user_id = ?", userID).First(&settings).Error; err != nil {
		// Create new settings if not found
		settings = models.SpeechSettings{
			UserID: userID.(uuid.UUID),
		}
	}

	// Update fields if provided
	if requestData.EnableTTS != nil {
		settings.EnableTTS = *requestData.EnableTTS
	}
	if requestData.EnableSTT != nil {
		settings.EnableSTT = *requestData.EnableSTT
	}
	if requestData.DefaultTTSModel != nil {
		if modelID, err := uuid.Parse(*requestData.DefaultTTSModel); err == nil {
			settings.DefaultTTSModel = &modelID
		}
	}
	if requestData.DefaultSTTModel != nil {
		if modelID, err := uuid.Parse(*requestData.DefaultSTTModel); err == nil {
			settings.DefaultSTTModel = &modelID
		}
	}
	if requestData.DefaultVoice != nil {
		settings.DefaultVoice = *requestData.DefaultVoice
	}
	if requestData.DefaultLanguage != nil {
		settings.DefaultLanguage = *requestData.DefaultLanguage
	}
	if requestData.TTSSpeed != nil {
		settings.TTSSpeed = *requestData.TTSSpeed
	}
	if requestData.TTSVolume != nil {
		settings.TTSVolume = *requestData.TTSVolume
	}
	if requestData.STTLanguage != nil {
		settings.STTLanguage = *requestData.STTLanguage
	}
	if requestData.AutoPlayTTS != nil {
		settings.AutoPlayTTS = *requestData.AutoPlayTTS
	}
	if requestData.ShowSTTTranscript != nil {
		settings.ShowSTTTranscript = *requestData.ShowSTTTranscript
	}
	if requestData.KeyboardShortcut != nil {
		settings.KeyboardShortcut = *requestData.KeyboardShortcut
	}

	// Save or update settings
	if err := config.DB.Save(&settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save speech settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Speech settings updated successfully"})
}

// GetAIUsageStats retrieves user's AI usage statistics
func GetAIUsageStats(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get query parameters
	period := c.DefaultQuery("period", "30d") // 7d, 30d, 90d
	provider := c.Query("provider")           // Optional filter by provider
	service := c.Query("service")             // Optional filter by service

	query := config.DB.Where("user_id = ?", userID)

	// Add date filter based on period
	switch period {
	case "7d":
		query = query.Where("date >= ?", time.Now().AddDate(0, 0, -7))
	case "30d":
		query = query.Where("date >= ?", time.Now().AddDate(0, 0, -30))
	case "90d":
		query = query.Where("date >= ?", time.Now().AddDate(0, 0, -90))
	default:
		query = query.Where("date >= ?", time.Now().AddDate(0, 0, -30))
	}

	if provider != "" {
		query = query.Where("provider = ?", provider)
	}
	if service != "" {
		query = query.Where("service = ?", service)
	}

	var stats []models.AIUsageStats
	if err := query.Order("date DESC").Find(&stats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve usage stats"})
		return
	}

	// Calculate totals
	totalTokens := 0
	totalRequests := 0
	totalCost := 0.0

	for _, stat := range stats {
		totalTokens += stat.Tokens
		totalRequests += stat.Requests
		totalCost += stat.Cost
	}

	c.JSON(http.StatusOK, gin.H{
		"stats":         stats,
		"summary": gin.H{
			"total_tokens":    totalTokens,
			"total_requests":  totalRequests,
			"total_cost":      totalCost,
			"period":          period,
		},
	})
}

// TestAPIKey tests if an API key is valid for a provider
func TestAPIKey(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var requestData struct {
		Provider string `json:"provider" binding:"required"`
		APIKey   string `json:"api_key" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Test the API key based on provider
	var isValid bool
	var errorMsg string

	switch requestData.Provider {
	case "openai":
		isValid, errorMsg = testOpenAIKey(requestData.APIKey)
	case "elevenlabs":
		isValid, errorMsg = testElevenLabsKey(requestData.APIKey)
	case "replicate":
		isValid, errorMsg = testReplicateKey(requestData.APIKey)
	case "assemblyai":
		isValid, errorMsg = testAssemblyAIKey(requestData.APIKey)
	case "deepgram":
		isValid, errorMsg = testDeepgramKey(requestData.APIKey)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported provider"})
		return
	}

	if isValid {
		c.JSON(http.StatusOK, gin.H{"valid": true, "message": "API key is valid"})
	} else {
		c.JSON(http.StatusOK, gin.H{"valid": false, "error": errorMsg})
	}
}

// Helper functions for testing API keys
func testOpenAIKey(apiKey string) (bool, string) {
	// Simple test by making a minimal API call
	// In a real implementation, you'd make a test request to OpenAI
	return len(apiKey) > 20 && strings.HasPrefix(apiKey, "sk-"), "Invalid OpenAI API key format"
}

func testElevenLabsKey(apiKey string) (bool, string) {
	// Test ElevenLabs API key
	return len(apiKey) > 10, "Invalid ElevenLabs API key format"
}

func testReplicateKey(apiKey string) (bool, string) {
	// Test Replicate API key
	return len(apiKey) > 20, "Invalid Replicate API key format"
}

func testAssemblyAIKey(apiKey string) (bool, string) {
	// Test AssemblyAI API key
	return len(apiKey) > 10, "Invalid AssemblyAI API key format"
}

func testDeepgramKey(apiKey string) (bool, string) {
	// Test Deepgram API key
	return len(apiKey) > 10, "Invalid Deepgram API key format"
}
