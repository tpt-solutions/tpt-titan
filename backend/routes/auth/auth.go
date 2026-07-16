package auth

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"tpt-titan/backend/models"
	"tpt-titan/backend/utils"
)

var jwtSecret []byte

// InitAuth initializes the auth package with the JWT secret from configuration
func InitAuth(secret string) {
	jwtSecret = []byte(secret)
}

// getJWTSecret returns the JWT secret, panics if not initialized
func getJWTSecret() []byte {
	if jwtSecret == nil {
		panic("JWT secret not initialized. Call InitAuth() before using auth functions")
	}
	return jwtSecret
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Username  string `json:"username" binding:"required,min=3,max=50"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type AuthResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

// AuthMiddleware validates JWT tokens and sets user context
func AuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return getJWTSecret(), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Extract user ID from token claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if userIDStr, ok := claims["user_id"].(string); ok {
				userID, err := uuid.Parse(userIDStr)
				if err != nil {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
					c.Abort()
					return
				}

				var user models.User
				if err := db.Where("id = ? AND is_active = ?", userID, true).First(&user).Error; err != nil {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found or inactive"})
					c.Abort()
					return
				}

				// Set user and user_id in context for handlers to access
				c.Set("user", user)
				c.Set("user_id", userID)
			}
		}

		c.Next()
	}
}

// Register creates a new user account
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := c.MustGet("db").(*gorm.DB)

	// Check if user already exists
	var existingUser models.User
	if err := db.Where("email = ? OR username = ?", req.Email, req.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User with this email or username already exists"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create user
	user := models.User{
		Email:     req.Email,
		Username:  req.Username,
		PasswordHash: string(hashedPassword),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		IsActive:  true,
	}

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Generate JWT token
	token, err := generateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Initialize encryption keys for user
	if err := initializeUserEncryption(db, user.ID, req.Password); err != nil {
		// Log error but don't fail registration
		c.Error(err) // This will be logged
	}

	response := AuthResponse{
		User: UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Username:  user.Username,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
		},
		Token: token,
	}

	c.JSON(http.StatusCreated, response)
}

// Login authenticates a user and returns a JWT token
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := c.MustGet("db").(*gorm.DB)

	// Find user
	var user models.User
	if err := db.Where("email = ? AND is_active = ?", req.Email, true).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Update last login
	db.Model(&user).Update("last_login", time.Now())

	// Generate JWT token
	token, err := generateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	response := AuthResponse{
		User: UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Username:  user.Username,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
		},
		Token: token,
	}

	c.JSON(http.StatusOK, response)
}

// GetProfile returns the current user's profile
func GetProfile(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	response := UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// UpdateProfile updates the current user's profile
func UpdateProfile(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	var req struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := c.MustGet("db").(*gorm.DB)

	updates := map[string]interface{}{
		"first_name": req.FirstName,
		"last_name":  req.LastName,
		"updated_at": time.Now(),
	}

	if err := db.Model(&user).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	// Fetch updated user
	db.First(&user, user.ID)

	response := UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// generateToken creates a JWT token for the user
func generateToken(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // 24 hours
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJWTSecret())
}

// initializeUserEncryption sets up encryption keys for a new user and persists
// the per-user derivation salt so encryption keys can be recovered later.
func initializeUserEncryption(db *gorm.DB, userID uuid.UUID, password string) error {
	// Create user's encryption key manager
	km, err := utils.NewKeyManager(password)
	if err != nil {
		return err
	}

	// Store the salt used for key derivation (safe to store). It is required to
	// deterministically re-derive the user's encryption key from their password.
	salt := km.GetSalt()
	if err := db.Model(&models.User{}).Where("id = ?", userID).
		Update("encryption_salt", base64.StdEncoding.EncodeToString(salt)).Error; err != nil {
		return fmt.Errorf("failed to persist encryption salt: %w", err)
	}

	return nil
}
