package routes

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"tpt-titan/backend/config"
	"tpt-titan/backend/models"
	"tpt-titan/backend/services"
	"tpt-titan/backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// DocumentRequest represents the request payload for creating/updating documents
type DocumentRequest struct {
	Title       string                 `json:"title" binding:"required"`
	ContentType string                 `json:"content_type" binding:"required"` // text, spreadsheet, form
	Content     map[string]interface{} `json:"content" binding:"required"`     // JSON content of the document
}

// GetDocuments returns all documents for the authenticated user
func GetDocuments(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var documents []models.EncryptedDocument
	if err := config.DB.Where("user_id = ? AND is_active = ?", userID, true).Order("updated_at DESC").Find(&documents).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve documents"})
		return
	}

	// Decrypt document metadata (titles) but not content for list view
	var response []gin.H
	for _, doc := range documents {
		response = append(response, gin.H{
			"id":          doc.ID,
			"title":       doc.Title,
			"content_type": doc.ContentType,
			"version":     doc.Version,
			"created_at":  doc.CreatedAt,
			"updated_at":  doc.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"documents": response})
}

// GetDocument returns a specific document with decrypted content
func GetDocument(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	documentID := c.Param("id")
	if documentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Document ID is required"})
		return
	}

	docID, err := uuid.Parse(documentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	var document models.EncryptedDocument
	if err := config.DB.Where("id = ? AND user_id = ? AND is_active = ?", docID, userID, true).First(&document).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	// Decrypt the content
	userPassword := utils.DeriveUserDocumentKey(userID)
	km, err := utils.DeriveKeyFromPassword(userPassword, document.Salt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize decryption"})
		return
	}

	decryptedData, err := km.Decrypt(document.EncryptedData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decrypt content"})
		return
	}

	// Parse decrypted JSON content
	var content map[string]interface{}
	if err := json.Unmarshal(decryptedData, &content); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse decrypted content"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":           document.ID,
		"title":        document.Title,
		"content_type": document.ContentType,
		"version":      document.Version,
		"created_at":   document.CreatedAt,
		"updated_at":   document.UpdatedAt,
		"content":      content,
	})
}

// CreateDocument creates a new document
func CreateDocument(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req DocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userPassword := utils.DeriveUserDocumentKey(userID)

	// Convert content to JSON for encryption
	contentJSON, err := json.Marshal(req.Content)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content format"})
		return
	}

	// Create key manager and encrypt
	km, err := utils.NewKeyManager(userPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize encryption"})
		return
	}

	encryptedContent, err := km.Encrypt(contentJSON)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt content"})
		return
	}

	// Create new document
	document := models.EncryptedDocument{
		UserID:        userID.(uuid.UUID),
		Title:         req.Title,
		ContentType:   req.ContentType,
		EncryptedData: encryptedContent,
		Salt:          km.GetSalt(),
		Algorithm:     "AES-256-GCM",
		FileSize:      int64(len(encryptedContent)),
		Version:       1,
		IsActive:      true,
	}

	if err := config.DB.Create(&document).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create document"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":           document.ID,
		"title":        document.Title,
		"content_type": document.ContentType,
		"version":      document.Version,
		"created_at":   document.CreatedAt,
		"updated_at":   document.UpdatedAt,
	})
}

// UpdateDocument updates an existing document (creates new version)
func UpdateDocument(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	documentID := c.Param("id")
	if documentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Document ID is required"})
		return
	}

	docID, err := uuid.Parse(documentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	var req DocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find existing document
	var existingDoc models.EncryptedDocument
	if err := config.DB.Where("id = ? AND user_id = ? AND is_active = ?", docID, userID, true).First(&existingDoc).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	// Encrypt the new content
	contentJSON, err := json.Marshal(req.Content)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content format"})
		return
	}

	// Create key manager and encrypt
	km, err := utils.NewKeyManager(utils.DeriveUserDocumentKey(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize encryption"})
		return
	}

	encryptedContent, err := km.Encrypt(contentJSON)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt content"})
		return
	}

	// Create new version
	newVersion := existingDoc.Version + 1
	newDocument := models.EncryptedDocument{
		UserID:        userID.(uuid.UUID),
		Title:         req.Title,
		ContentType:   req.ContentType,
		EncryptedData: encryptedContent,
		Salt:          km.GetSalt(),
		Algorithm:     "AES-256-GCM",
		FileSize:      int64(len(encryptedContent)),
		Version:       newVersion,
		IsActive:      true,
	}

	if err := config.DB.Create(&newDocument).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update document"})
		return
	}

	// Mark old version as inactive (soft delete)
	if err := config.DB.Model(&existingDoc).Update("is_active", false).Error; err != nil {
		// Log error but don't fail the request
		c.Error(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"id":           newDocument.ID,
		"title":        newDocument.Title,
		"content_type": newDocument.ContentType,
		"version":      newDocument.Version,
		"created_at":   newDocument.CreatedAt,
		"updated_at":   newDocument.UpdatedAt,
	})
}

// DeleteDocument soft deletes a document
func DeleteDocument(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	documentID := c.Param("id")
	if documentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Document ID is required"})
		return
	}

	docID, err := uuid.Parse(documentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	// Soft delete by marking as inactive
	result := config.DB.Model(&models.EncryptedDocument{}).Where("id = ? AND user_id = ?", docID, userID).Update("is_active", false)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete document"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document deleted successfully"})
}

// GetDocumentVersions returns version history for a document
func GetDocumentVersions(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	documentID := c.Param("id")
	if documentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Document ID is required"})
		return
	}

	docID, err := uuid.Parse(documentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	var versions []models.EncryptedDocument
	if err := config.DB.Where("id = ? AND user_id = ?", docID, userID).Order("version DESC").Find(&versions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve document versions"})
		return
	}

	var response []gin.H
	for _, version := range versions {
		response = append(response, gin.H{
			"id":         version.ID,
			"version":    version.Version,
			"title":      version.Title,
			"created_at": version.CreatedAt,
			"updated_at": version.UpdatedAt,
			"is_active":  version.IsActive,
		})
	}

	c.JSON(http.StatusOK, gin.H{"versions": response})
}

// RestoreDocumentVersion restores a specific version as the active version
func RestoreDocumentVersion(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	documentID := c.Param("id")
	versionStr := c.Param("version")

	if documentID == "" || versionStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Document ID and version are required"})
		return
	}

	docID, err := uuid.Parse(documentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	version, err := strconv.Atoi(versionStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid version number"})
		return
	}

	// Find the specified version
	var targetVersion models.EncryptedDocument
	if err := config.DB.Where("id = ? AND user_id = ? AND version = ?", docID, userID, version).First(&targetVersion).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Version not found"})
		return
	}

	// Mark all other versions as inactive
	if err := config.DB.Model(&models.EncryptedDocument{}).Where("id = ? AND user_id = ?", docID, userID).Update("is_active", false).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update versions"})
		return
	}

	// Mark target version as active
	if err := config.DB.Model(&targetVersion).Update("is_active", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to restore version"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Version restored successfully",
		"id":      targetVersion.ID,
		"version": targetVersion.Version,
	})
}

// DocumentUploadRequest represents the request payload for document uploads with AI processing
type DocumentUploadRequest struct {
	Title         string `json:"title" binding:"required"`
	FileName      string `json:"file_name" binding:"required"`
	FileData      string `json:"file_data" binding:"required"` // Base64 encoded file data
	FileType      string `json:"file_type" binding:"required"` // "pdf", "image", etc.
	ProcessWithAI bool   `json:"process_with_ai"`             // Whether to process with AI
	AnalysisType  string `json:"analysis_type,omitempty"`     // "ocr", "invoice", "receipt", "business_card", "contract"
}

// UploadDocumentWithAI uploads a document and optionally processes it with AI
func UploadDocumentWithAI(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req DocumentUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Decode base64 file data
	fileData, err := base64.StdEncoding.DecodeString(req.FileData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file data"})
		return
	}

	// Create document record
	document := &models.EncryptedDocument{
		UserID:      userID.(uuid.UUID),
		Title:       req.Title,
		FileName:    req.FileName,
		ContentType: req.FileType,
		FileSize:    int64(len(fileData)),
		Version:     1,
		IsActive:    true,
	}

	// Encrypt and store file data
	userPassword := utils.DeriveUserDocumentKey(userID)
	km, err := utils.NewKeyManager(userPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize encryption"})
		return
	}

	encryptedData, err := km.Encrypt(fileData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt file data"})
		return
	}

	document.EncryptedData = encryptedData
	document.Salt = km.GetSalt()
	document.Algorithm = "AES-256-GCM"

	// Save document
	if err := config.DB.Create(document).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save document"})
		return
	}

	response := gin.H{
		"id":          document.ID,
		"title":       document.Title,
		"file_type":   req.FileType,
		"file_size":   document.FileSize,
		"created_at":  document.CreatedAt,
		"status":      "uploaded",
	}

	// If AI processing is requested, start background processing
	if req.ProcessWithAI {
		go processDocumentWithAI(document.ID, userID.(uuid.UUID), req.FileName, fileData, req.FileType, req.AnalysisType)

		response["processing_status"] = "queued"
		response["message"] = "Document uploaded and AI processing started"
	} else {
		response["message"] = "Document uploaded successfully"
	}

	c.JSON(http.StatusCreated, response)
}

// processDocumentWithAI processes a document with AI in the background
func processDocumentWithAI(documentID, userID uuid.UUID, fileName string, fileData []byte, fileType, analysisType string) {
	startTime := time.Now()

	// Send initial processing notification
	if wsHub != nil {
		wsMessage := models.WebSocketMessage{
			Type: models.WSDocumentProcessing,
			Data: gin.H{
				"document_id": documentID,
				"status":      "started",
				"message":     "AI document processing has started",
			},
		}
		wsHub.BroadcastToUser(userID, wsMessage)
	}

	// Create analysis record
	analysis := &models.DocumentAnalysis{
		UserID:     userID,
		DocumentID: documentID,
		FileName:   fileName,
		FileType:   fileType,
		Status:     "processing",
		Pages:      1, // Default, will be updated after analysis
		Language:   "en",
	}

	if err := config.DB.Create(analysis).Error; err != nil {
		// Log error and update status
		analysis.Status = "failed"
		analysis.Error = "Failed to create analysis record: " + err.Error()
		config.DB.Save(analysis)

		// Send WebSocket notification for failure
		if wsHub != nil {
			wsMessage := models.WebSocketMessage{
				Type: models.WSDocumentFailed,
				Data: gin.H{
					"document_id": documentID,
					"error":       analysis.Error,
				},
			}
			wsHub.BroadcastToUser(userID, wsMessage)
		}
		return
	}

	// Find appropriate vision model
	var visionModel models.AIModel
	if err := config.DB.Where("capabilities @> ? AND is_active = ?", `["vision"]`, true).First(&visionModel).Error; err != nil {
		analysis.Status = "failed"
		analysis.Error = "No suitable vision model available"
		config.DB.Save(analysis)

		// Send WebSocket notification for failure
		if wsHub != nil {
			wsMessage := models.WebSocketMessage{
				Type: models.WSDocumentFailed,
				Data: gin.H{
					"document_id": documentID,
					"error":       analysis.Error,
				},
			}
			wsHub.BroadcastToUser(userID, wsMessage)
		}
		return
	}

	// Perform AI analysis
	ollamaService := services.NewOllamaService("localhost", "11434")
	result, err := ollamaService.AnalyzeDocument(visionModel.ModelID, fileData, analysisType)
	if err != nil {
		analysis.Status = "failed"
		analysis.Error = err.Error()
		config.DB.Save(analysis)

		// Send WebSocket notification for failure
		if wsHub != nil {
			wsMessage := models.WebSocketMessage{
				Type: models.WSDocumentFailed,
				Data: gin.H{
					"document_id": documentID,
					"error":       analysis.Error,
				},
			}
			wsHub.BroadcastToUser(userID, wsMessage)
		}
		return
	}

	// Update analysis record with results
	analysis.Status = "completed"
	analysis.TextContent = result.TextContent
	analysis.Confidence = result.Confidence
	analysis.Pages = result.Pages
	analysis.ProcessingTime = int(time.Since(startTime).Milliseconds()) // Track actual processing time

	// Serialize structured data
	if len(result.Fields) > 0 {
		fieldsJSON, _ := json.Marshal(result.Fields)
		analysis.Fields = string(fieldsJSON)
	}

	if len(result.Tables) > 0 {
		tablesJSON, _ := json.Marshal(result.Tables)
		analysis.Tables = string(tablesJSON)
	}

	if err := config.DB.Save(analysis).Error; err != nil {
		// Log error but analysis was successful
		fmt.Printf("Warning: Failed to save analysis results: %v\n", err)
	}

	// Send WebSocket notification for completion
	if wsHub != nil {
		wsMessage := models.WebSocketMessage{
			Type: models.WSDocumentProcessed,
			Data: gin.H{
				"document_id":    documentID,
				"analysis_id":    analysis.ID,
				"status":         "completed",
				"text_content":   result.TextContent,
				"confidence":     result.Confidence,
				"pages":          result.Pages,
				"fields_count":   len(result.Fields),
				"tables_count":   len(result.Tables),
				"message":        "AI document processing completed successfully",
			},
		}
		wsHub.BroadcastToUser(userID, wsMessage)
	}
}

// ProcessDocumentWithAI processes an existing document with AI
func ProcessDocumentWithAI(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	documentID := c.Param("id")
	if documentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Document ID is required"})
		return
	}

	docID, err := uuid.Parse(documentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	var req struct {
		AnalysisType string `json:"analysis_type" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find document
	var document models.EncryptedDocument
	if err := config.DB.Where("id = ? AND user_id = ? AND is_active = ?", docID, userID, true).First(&document).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	// Decrypt file data
	userPassword := utils.DeriveUserDocumentKey(userID)
	km, err := utils.DeriveKeyFromPassword(userPassword, document.Salt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize decryption"})
		return
	}

	fileData, err := km.Decrypt(document.EncryptedData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decrypt document"})
		return
	}

	// Start background processing
	go processDocumentWithAI(docID, userID.(uuid.UUID), document.FileName, fileData, document.ContentType, req.AnalysisType)

	c.JSON(http.StatusAccepted, gin.H{
		"message": "AI processing started",
		"status":  "queued",
	})
}

// GetDocumentAnalysis returns AI analysis results for a document
func GetDocumentAnalysis(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	documentID := c.Param("id")
	if documentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Document ID is required"})
		return
	}

	docID, err := uuid.Parse(documentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	// Find latest analysis for this document
	var analysis models.DocumentAnalysis
	if err := config.DB.Where("document_id = ? AND user_id = ?", docID, userID).
		Order("created_at DESC").First(&analysis).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No analysis found for this document"})
		return
	}

	// Parse structured data
	var fields []models.ExtractedField
	var tables []models.ExtractedTable

	if analysis.Fields != "" {
		json.Unmarshal([]byte(analysis.Fields), &fields)
	}

	if analysis.Tables != "" {
		json.Unmarshal([]byte(analysis.Tables), &tables)
	}

	c.JSON(http.StatusOK, gin.H{
		"analysis_id":    analysis.ID,
		"document_id":    analysis.DocumentID,
		"status":         analysis.Status,
		"text_content":   analysis.TextContent,
		"fields":         fields,
		"tables":         tables,
		"confidence":     analysis.Confidence,
		"pages":          analysis.Pages,
		"language":       analysis.Language,
		"processing_time": analysis.ProcessingTime,
		"created_at":     analysis.CreatedAt,
		"error":          analysis.Error,
	})
}

// GetDocumentProcessingStatus returns the processing status for a document
func GetDocumentProcessingStatus(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	documentID := c.Param("id")
	if documentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Document ID is required"})
		return
	}

	docID, err := uuid.Parse(documentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	// Find latest analysis for this document
	var analysis models.DocumentAnalysis
	err = config.DB.Where("document_id = ? AND user_id = ?", docID, userID).
		Order("created_at DESC").First(&analysis).Error

	if err != nil {
		// No analysis found - document not processed yet
		c.JSON(http.StatusOK, gin.H{
			"status": "not_processed",
			"message": "Document has not been processed with AI yet",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": analysis.Status,
		"analysis_id": analysis.ID,
		"error": analysis.Error,
		"processing_time": analysis.ProcessingTime,
		"created_at": analysis.CreatedAt,
		"updated_at": analysis.UpdatedAt,
	})
}

// GetDocumentAnalyses returns all analyses for a document
func GetDocumentAnalyses(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	documentID := c.Param("id")
	if documentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Document ID is required"})
		return
	}

	docID, err := uuid.Parse(documentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	var analyses []models.DocumentAnalysis
	if err := config.DB.Where("document_id = ? AND user_id = ?", docID, userID).
		Order("created_at DESC").Find(&analyses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve analyses"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"analyses": analyses})
}
