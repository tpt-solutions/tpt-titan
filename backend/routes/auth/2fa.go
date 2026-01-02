package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"tpt-titan/backend/models"
	"tpt-titan/backend/services"
)

// EnableTOTP enables two-factor authentication for the current user
func EnableTOTP(c *gin.Context) {
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

	// Get auth service from context
	authServiceInterface, exists := c.Get("auth_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Auth service not available"})
		return
	}

	authService, ok := authServiceInterface.(*services.AuthService)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid auth service"})
		return
	}

	// Enable TOTP
	totpSecret, err := authService.EnableTOTP(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "TOTP setup initiated",
		"data":    totpSecret,
	})
}

// VerifyAndEnableTOTP verifies TOTP code and enables 2FA
func VerifyAndEnableTOTP(c *gin.Context) {
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

	var req struct {
		Code string `json:"code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get auth service from context
	authServiceInterface, exists := c.Get("auth_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Auth service not available"})
		return
	}

	authService, ok := authServiceInterface.(*services.AuthService)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid auth service"})
		return
	}

	// Verify and enable TOTP
	err := authService.VerifyAndEnableTOTP(userID, req.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Two-factor authentication enabled successfully",
	})
}

// DisableTOTP disables two-factor authentication
func DisableTOTP(c *gin.Context) {
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

	var req struct {
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get auth service from context
	authServiceInterface, exists := c.Get("auth_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Auth service not available"})
		return
	}

	authService, ok := authServiceInterface.(*services.AuthService)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid auth service"})
		return
	}

	// Disable TOTP
	err := authService.DisableTOTP(userID, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Two-factor authentication disabled successfully",
	})
}

// RequestPasswordReset initiates password reset process
func RequestPasswordReset(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get auth service from context
	authServiceInterface, exists := c.Get("auth_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Auth service not available"})
		return
	}

	authService, ok := authServiceInterface.(*services.AuthService)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid auth service"})
		return
	}

	// Request password reset
	err := authService.RequestPasswordReset(req.Email)
	if err != nil {
		// Don't reveal if email exists or not for security
		c.JSON(http.StatusOK, gin.H{
			"message": "If an account with this email exists, a password reset link has been sent.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "If an account with this email exists, a password reset link has been sent.",
	})
}

// ResetPassword resets user password using reset token
func ResetPassword(c *gin.Context) {
	var req struct {
		Token    string `json:"token" binding:"required"`
		Password string `json:"password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get auth service from context
	authServiceInterface, exists := c.Get("auth_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Auth service not available"})
		return
	}

	authService, ok := authServiceInterface.(*services.AuthService)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid auth service"})
		return
	}

	// Reset password
	err := authService.ResetPassword(req.Token, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset successfully",
	})
}

// ChangePassword changes user password (authenticated)
func ChangePassword(c *gin.Context) {
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

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get auth service from context
	authServiceInterface, exists := c.Get("auth_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Auth service not available"})
		return
	}

	authService, ok := authServiceInterface.(*services.AuthService)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid auth service"})
		return
	}

	// Change password
	err := authService.ChangePassword(userID, req.OldPassword, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password changed successfully",
	})
}

// GetUserProfile returns the current user's profile
func GetUserProfile(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	user, ok := userInterface.(*models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user data"})
		return
	}

	response := user.ToResponse()
	c.JSON(http.StatusOK, gin.H{"user": response})
}

// UpdateUserProfile updates the current user's profile
func UpdateUserProfile(c *gin.Context) {
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

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get auth service from context
	authServiceInterface, exists := c.Get("auth_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Auth service not available"})
		return
	}

	authService, ok := authServiceInterface.(*services.AuthService)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid auth service"})
		return
	}

	// Update profile
	err := authService.UpdateProfile(userID, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
	})
}

// Logout logs out the current user
func Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) == 2 && tokenParts[0] == "Bearer" {
			// Get auth service from context
			authServiceInterface, exists := c.Get("auth_service")
			if exists {
				authService, ok := authServiceInterface.(*services.AuthService)
				if ok {
					authService.Logout(tokenParts[1])
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}
