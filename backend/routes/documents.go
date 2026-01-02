package routes

import (
	"net/http"
	"strconv"
	"tpt-titan/backend/config"
	"tpt-titan/backend/models"

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

	// For now, return document metadata. Full content decryption would require crypto implementation
	c.JSON(http.StatusOK, gin.H{
		"id":           document.ID,
		"title":        document.Title,
		"content_type": document.ContentType,
		"version":      document.Version,
		"created_at":   document.CreatedAt,
		"updated_at":   document.UpdatedAt,
		// TODO: Add decrypted content when crypto is implemented
		"content":      "Content decryption not yet implemented",
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

	// Create new document
	document := models.EncryptedDocument{
		UserID:      userID.(uuid.UUID),
		Title:       req.Title,
		ContentType: req.ContentType,
		// TODO: Encrypt content when crypto is implemented
		EncryptedData: []byte("Encrypted content placeholder"),
		Salt:          []byte("salt_placeholder"),
		Algorithm:     "AES-256-GCM",
		FileSize:      0, // TODO: Calculate actual size
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

	// Create new version
	newVersion := existingDoc.Version + 1
	newDocument := models.EncryptedDocument{
		UserID:        userID.(uuid.UUID),
		Title:         req.Title,
		ContentType:   req.ContentType,
		EncryptedData: []byte("Updated encrypted content placeholder"), // TODO: Encrypt actual content
		Salt:          []byte("salt_placeholder"),
		Algorithm:     "AES-256-GCM",
		FileSize:      0, // TODO: Calculate actual size
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
