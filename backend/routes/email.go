package routes

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"
	"tpt-titan/backend/models"
	"tpt-titan/backend/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var emailService *services.EmailService

// InitEmailService initializes the email service (called from main)
func InitEmailService(db *sql.DB) {
	emailService = services.NewEmailService(db)
}

// GetEmailAccounts returns all email accounts for the authenticated user
func GetEmailAccounts(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	accounts, err := emailService.GetEmailAccounts(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve email accounts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"accounts": accounts})
}

// GetEmailAccount returns a specific email account
func GetEmailAccount(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	accountID := c.Param("id")
	if accountID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account ID is required"})
		return
	}

	id, err := uuid.Parse(accountID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	account, err := emailService.GetEmailAccount(userID.(uuid.UUID), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve email account"})
		return
	}

	if account == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Email account not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"account": account})
}

// CreateEmailAccount creates a new email account
func CreateEmailAccount(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.EmailAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account, err := emailService.CreateEmailAccount(userID.(uuid.UUID), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create email account"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"account": account})
}

// UpdateEmailAccount updates an existing email account
func UpdateEmailAccount(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	accountID := c.Param("id")
	if accountID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account ID is required"})
		return
	}

	id, err := uuid.Parse(accountID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	var req models.EmailAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account, err := emailService.UpdateEmailAccount(userID.(uuid.UUID), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update email account"})
		return
	}

	if account == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Email account not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"account": account})
}

// DeleteEmailAccount deletes an email account
func DeleteEmailAccount(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	accountID := c.Param("id")
	if accountID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account ID is required"})
		return
	}

	id, err := uuid.Parse(accountID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	err = emailService.DeleteEmailAccount(userID.(uuid.UUID), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete email account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email account deleted successfully"})
}

// GetEmails returns emails with search and pagination
func GetEmails(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse query parameters
	searchReq := models.EmailSearchRequest{
		Query:     c.Query("q"),
		Folder:    c.Query("folder"),
		Limit:     50,  // Default limit
		Offset:    0,   // Default offset
	}

	// Parse optional parameters
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 && limit <= 100 {
			searchReq.Limit = limit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			searchReq.Offset = offset
		}
	}

	if isReadStr := c.Query("is_read"); isReadStr != "" {
		if isRead, err := strconv.ParseBool(isReadStr); err == nil {
			searchReq.IsRead = &isRead
		}
	}

	if isStarredStr := c.Query("is_starred"); isStarredStr != "" {
		if isStarred, err := strconv.ParseBool(isStarredStr); err == nil {
			searchReq.IsStarred = &isStarred
		}
	}

	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if startDate, err := time.Parse(time.RFC3339, startDateStr); err == nil {
			searchReq.StartDate = &startDate
		}
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if endDate, err := time.Parse(time.RFC3339, endDateStr); err == nil {
			searchReq.EndDate = &endDate
		}
	}

	emails, err := emailService.GetEmails(userID.(uuid.UUID), searchReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve emails"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"emails": emails})
}

// GetEmail returns a specific email
func GetEmail(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	emailID := c.Param("id")
	if emailID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email ID is required"})
		return
	}

	id, err := uuid.Parse(emailID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email ID"})
		return
	}

	email, err := emailService.GetEmail(userID.(uuid.UUID), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve email"})
		return
	}

	if email == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Email not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"email": email})
}

// SendEmail sends a new email
func SendEmail(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.EmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	email, err := emailService.SendEmail(userID.(uuid.UUID), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"email": email})
}

// MarkEmailAsRead marks an email as read/unread
func MarkEmailAsRead(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	emailID := c.Param("id")
	if emailID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email ID is required"})
		return
	}

	id, err := uuid.Parse(emailID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email ID"})
		return
	}

	var req struct {
		IsRead bool `json:"is_read" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = emailService.MarkEmailAsRead(userID.(uuid.UUID), id, req.IsRead)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark email as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email marked successfully"})
}

// StarEmail marks an email as starred/unstarred
func StarEmail(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	emailID := c.Param("id")
	if emailID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email ID is required"})
		return
	}

	id, err := uuid.Parse(emailID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email ID"})
		return
	}

	var req struct {
		IsStarred bool `json:"is_starred" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = emailService.StarEmail(userID.(uuid.UUID), id, req.IsStarred)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to star email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email starred successfully"})
}

// MoveEmailToFolder moves an email to a different folder
func MoveEmailToFolder(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	emailID := c.Param("id")
	if emailID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email ID is required"})
		return
	}

	id, err := uuid.Parse(emailID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email ID"})
		return
	}

	var req struct {
		Folder string `json:"folder" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = emailService.MoveEmailToFolder(userID.(uuid.UUID), id, req.Folder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to move email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email moved successfully"})
}

// SyncEmails syncs emails from external accounts
func SyncEmails(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	accountID := c.Param("accountId")
	if accountID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account ID is required"})
		return
	}

	id, err := uuid.Parse(accountID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	err = emailService.SyncEmails(userID.(uuid.UUID), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sync emails"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Emails synced successfully"})
}

// GetEmailStats returns email statistics
func GetEmailStats(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	stats, err := emailService.GetEmailStats(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get email stats"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"stats": stats})
}
