package routes

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"
	"tpt-titan/backend/config"
	"tpt-titan/backend/models"
	"tpt-titan/backend/services"
	"tpt-titan/backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var voiceService *services.VoiceService

// InitVoiceService initializes the voice service (called from main)
func InitVoiceService() {
	voiceService = services.NewVoiceService()
}

// VoiceNoteRequest represents the request payload for creating/updating voice notes
type VoiceNoteRequest struct {
	Title      string   `json:"title" binding:"required"`
	Tags       []string `json:"tags,omitempty"`
	IsFavorite bool     `json:"is_favorite"`
	IsPublic   bool     `json:"is_public"`
}

// VoiceNoteUploadRequest represents the request payload for uploading voice note audio
type VoiceNoteUploadRequest struct {
	AudioData   string `json:"audio_data" binding:"required"` // Base64 encoded audio
	AudioFormat string `json:"audio_format" binding:"required"`
	Duration    int    `json:"duration" binding:"required"` // Duration in seconds
	Title       string `json:"title,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	IsFavorite  bool   `json:"is_favorite"`
	IsPublic    bool   `json:"is_public"`
}

// VoiceAnnotationRequest represents the request payload for creating voice annotations
type VoiceAnnotationRequest struct {
	ContentType string    `json:"content_type" binding:"required"` // document, task, email, calendar, contact
	ContentID   string    `json:"content_id" binding:"required"`   // ID of the content being annotated
	Title       string    `json:"title" binding:"required"`
	Position    *string   `json:"position,omitempty"` // JSON position data for highlighting
	IsPublic    bool      `json:"is_public"`
}

// VoiceAnnotationUploadRequest represents the request payload for uploading annotation audio
type VoiceAnnotationUploadRequest struct {
	AudioData   string  `json:"audio_data" binding:"required"` // Base64 encoded audio
	AudioFormat string  `json:"audio_format" binding:"required"`
	Duration    int     `json:"duration" binding:"required"` // Duration in seconds
	ContentType string  `json:"content_type" binding:"required"`
	ContentID   string  `json:"content_id" binding:"required"`
	Title       string  `json:"title" binding:"required"`
	Position    *string `json:"position,omitempty"`
	IsPublic    bool    `json:"is_public"`
}

// GetVoiceNotes returns all voice notes for the authenticated user
func GetVoiceNotes(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse query parameters
	limit := 50 // Default limit
	offset := 0 // Default offset
	favoritesOnly := false
	tag := c.Query("tag")

	if limitStr := c.Query("limit"); limitStr != "" {
		if l := parseInt(limitStr); l > 0 && l <= 200 {
			limit = l
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if o := parseInt(offsetStr); o >= 0 {
			offset = o
		}
	}

	if favStr := c.Query("favorites"); favStr == "true" {
		favoritesOnly = true
	}

	notes, err := voiceService.GetVoiceNotes(userID.(uuid.UUID), limit, offset, favoritesOnly, tag)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve voice notes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"notes": notes})
}

// GetVoiceNote returns a specific voice note
func GetVoiceNote(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	noteID := c.Param("id")
	if noteID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Voice note ID is required"})
		return
	}

	id, err := uuid.Parse(noteID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid voice note ID"})
		return
	}

	note, err := voiceService.GetVoiceNote(userID.(uuid.UUID), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve voice note"})
		return
	}

	if note == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Voice note not found"})
		return
	}

	// Decrypt audio data for playback
	userPassword := utils.DeriveUserDocumentKey(userID)
	km, err := utils.DeriveKeyFromPassword(userPassword, note.Salt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize decryption"})
		return
	}

	audioData, err := km.Decrypt(note.EncryptedAudioData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decrypt audio data"})
		return
	}

	response := gin.H{
		"id":           note.ID,
		"title":        note.Title,
		"content":      note.Content,
		"audio_format": note.AudioFormat,
		"duration":     note.Duration,
		"file_size":    note.FileSize,
		"tags":         note.Tags,
		"is_favorite":  note.IsFavorite,
		"is_public":    note.IsPublic,
		"created_at":   note.CreatedAt,
		"updated_at":   note.UpdatedAt,
		"audio_data":   base64.StdEncoding.EncodeToString(audioData),
	}

	c.JSON(http.StatusOK, response)
}

// CreateVoiceNote creates a new voice note with audio upload
func CreateVoiceNote(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req VoiceNoteUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Decode base64 audio data
	audioData, err := base64.StdEncoding.DecodeString(req.AudioData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid audio data"})
		return
	}

	// Auto-transcribe the audio using speech service
	transcript := ""
	if speechService != nil {
		// Use default STT model for transcription
		models, err := speechService.GetAvailableModels(userID.(uuid.UUID), "stt")
		if err == nil && len(models) > 0 {
			// Create a temporary speech request for transcription
			request, err := speechService.SpeechToText(userID.(uuid.UUID), audioData, models[0].ID, services.SpeechOptions{
				Language:    "en",
				AudioFormat: req.AudioFormat,
			})
			if err == nil {
				// Wait for transcription to complete (simplified - in production, this would be async)
				// For now, we'll store empty transcript and update it later via WebSocket
				transcript = "Transcription in progress..."
			}
		}
	}

	// Create voice note
	note := &models.VoiceNote{
		UserID:      userID.(uuid.UUID),
		Title:       req.Title,
		Content:     transcript,
		AudioFormat: req.AudioFormat,
		Duration:    req.Duration,
		FileSize:    int64(len(audioData)),
		Tags:        req.Tags,
		IsFavorite:  req.IsFavorite,
		IsPublic:    req.IsPublic,
	}

	// Encrypt audio data
	userPassword := utils.DeriveUserDocumentKey(userID)
	km, err := utils.NewKeyManager(userPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize encryption"})
		return
	}

	encryptedAudio, err := km.Encrypt(audioData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt audio data"})
		return
	}

	note.EncryptedAudioData = encryptedAudio
	note.Salt = km.GetSalt()
	note.Algorithm = "AES-256-GCM"

	if err := config.DB.Create(note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create voice note"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":           note.ID,
		"title":        note.Title,
		"content":      note.Content,
		"audio_format": note.AudioFormat,
		"duration":     note.Duration,
		"tags":         note.Tags,
		"is_favorite":  note.IsFavorite,
		"is_public":    note.IsPublic,
		"created_at":   note.CreatedAt,
	})
}

// UpdateVoiceNote updates an existing voice note
func UpdateVoiceNote(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	noteID := c.Param("id")
	if noteID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Voice note ID is required"})
		return
	}

	id, err := uuid.Parse(noteID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid voice note ID"})
		return
	}

	var req VoiceNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = voiceService.UpdateVoiceNote(userID.(uuid.UUID), id, req.Title, req.Tags, req.IsFavorite, req.IsPublic)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update voice note"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Voice note updated successfully"})
}

// DeleteVoiceNote deletes a voice note
func DeleteVoiceNote(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	noteID := c.Param("id")
	if noteID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Voice note ID is required"})
		return
	}

	id, err := uuid.Parse(noteID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid voice note ID"})
		return
	}

	err = voiceService.DeleteVoiceNote(userID.(uuid.UUID), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete voice note"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Voice note deleted successfully"})
}

// GetVoiceAnnotations returns voice annotations for specific content
func GetVoiceAnnotations(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	contentType := c.Query("content_type")
	contentIDStr := c.Query("content_id")

	if contentType == "" || contentIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "content_type and content_id are required"})
		return
	}

	contentID, err := uuid.Parse(contentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content ID"})
		return
	}

	annotations, err := voiceService.GetVoiceAnnotations(userID.(uuid.UUID), contentType, contentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve voice annotations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"annotations": annotations})
}

// GetVoiceAnnotation returns a specific voice annotation
func GetVoiceAnnotation(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	annotationID := c.Param("id")
	if annotationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Voice annotation ID is required"})
		return
	}

	id, err := uuid.Parse(annotationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid voice annotation ID"})
		return
	}

	annotation, err := voiceService.GetVoiceAnnotation(userID.(uuid.UUID), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve voice annotation"})
		return
	}

	if annotation == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Voice annotation not found"})
		return
	}

	// Decrypt audio data for playback
	userPassword := utils.DeriveUserDocumentKey(userID)
	km, err := utils.DeriveKeyFromPassword(userPassword, annotation.Salt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize decryption"})
		return
	}

	audioData, err := km.Decrypt(annotation.EncryptedAudioData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decrypt audio data"})
		return
	}

	response := gin.H{
		"id":           annotation.ID,
		"content_type": annotation.ContentType,
		"content_id":   annotation.ContentID,
		"title":        annotation.Title,
		"content":      annotation.Content,
		"audio_format": annotation.AudioFormat,
		"duration":     annotation.Duration,
		"position":     annotation.Position,
		"is_public":    annotation.IsPublic,
		"created_at":   annotation.CreatedAt,
		"updated_at":   annotation.UpdatedAt,
		"audio_data":   base64.StdEncoding.EncodeToString(audioData),
	}

	c.JSON(http.StatusOK, response)
}

// CreateVoiceAnnotation creates a new voice annotation with audio upload
func CreateVoiceAnnotation(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req VoiceAnnotationUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contentID, err := uuid.Parse(req.ContentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content ID"})
		return
	}

	// Decode base64 audio data
	audioData, err := base64.StdEncoding.DecodeString(req.AudioData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid audio data"})
		return
	}

	// Auto-transcribe the audio using speech service
	transcript := ""
	if speechService != nil {
		// Use default STT model for transcription
		models, err := speechService.GetAvailableModels(userID.(uuid.UUID), "stt")
		if err == nil && len(models) > 0 {
			request, err := speechService.SpeechToText(userID.(uuid.UUID), audioData, models[0].ID, services.SpeechOptions{
				Language:    "en",
				AudioFormat: req.AudioFormat,
			})
			if err == nil {
				transcript = "Transcription in progress..."
			}
		}
	}

	// Create voice annotation
	annotation := &models.VoiceAnnotation{
		UserID:      userID.(uuid.UUID),
		ContentType: req.ContentType,
		ContentID:   contentID,
		Title:       req.Title,
		Content:     transcript,
		AudioFormat: req.AudioFormat,
		Duration:    req.Duration,
		FileSize:    int64(len(audioData)),
		Position:    req.Position,
		IsPublic:    req.IsPublic,
	}

	// Encrypt audio data
	userPassword := utils.DeriveUserDocumentKey(userID)
	km, err := utils.NewKeyManager(userPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize encryption"})
		return
	}

	encryptedAudio, err := km.Encrypt(audioData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt audio data"})
		return
	}

	annotation.EncryptedAudioData = encryptedAudio
	annotation.Salt = km.GetSalt()
	annotation.Algorithm = "AES-256-GCM"

	if err := config.DB.Create(annotation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create voice annotation"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":           annotation.ID,
		"content_type": annotation.ContentType,
		"content_id":   annotation.ContentID,
		"title":        annotation.Title,
		"content":      annotation.Content,
		"audio_format": annotation.AudioFormat,
		"duration":     annotation.Duration,
		"position":     annotation.Position,
		"is_public":    annotation.IsPublic,
		"created_at":   annotation.CreatedAt,
	})
}

// DeleteVoiceAnnotation deletes a voice annotation
func DeleteVoiceAnnotation(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	annotationID := c.Param("id")
	if annotationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Voice annotation ID is required"})
		return
	}

	id, err := uuid.Parse(annotationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid voice annotation ID"})
		return
	}

	err = voiceService.DeleteVoiceAnnotation(userID.(uuid.UUID), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete voice annotation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Voice annotation deleted successfully"})
}

// Helper function to parse int (duplicate, but keeping for clarity)
func parseInt(s string) int {
	var result int
	fmt.Sscanf(s, "%d", &result)
	return result
}
