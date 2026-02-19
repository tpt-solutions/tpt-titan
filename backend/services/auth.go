package services

import (
	"crypto/rand"
	"database/sql"
	"encoding/base32"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"

	"tpt-titan/backend/models"
)

// AuthService handles authentication and authorization
type AuthService struct {
	db         *sql.DB
	jwtSecret  []byte
	cache      *CacheService
}

// UserCredentials represents login credentials
type UserCredentials struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	TOTPCode string `json:"totp_code,omitempty"` // For 2FA
}

// RegisterRequest represents user registration data
type RegisterRequest struct {
	Username  string `json:"username" binding:"required,min=3,max=50"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

// LoginResponse represents successful login response
type LoginResponse struct {
	User         models.UserResponse `json:"user"`
	AccessToken  string             `json:"access_token"`
	RefreshToken string             `json:"refresh_token"`
	TokenType    string             `json:"token_type"`
	ExpiresIn    int64              `json:"expires_in"`
	RequiresTOTP bool               `json:"requires_totp,omitempty"`
}

// TOTPSecret represents TOTP setup information
type TOTPSecret struct {
	Secret    string `json:"secret"`
	QRCodeURL string `json:"qr_code_url"`
}

// PasswordResetRequest represents password reset request
type PasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// PasswordResetConfirm represents password reset confirmation
type PasswordResetConfirm struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

// NewAuthService creates a new authentication service
func NewAuthService(db *sql.DB, jwtSecret string, cache *CacheService) *AuthService {
	return &AuthService{
		db:        db,
		jwtSecret: []byte(jwtSecret),
		cache:     cache,
	}
}

// Register creates a new user account
func (as *AuthService) Register(req RegisterRequest) (*models.UserResponse, error) {
	// Check if user already exists
	var existingID uuid.UUID
	err := as.db.QueryRow("SELECT id FROM users WHERE email = $1 OR username = $2",
		req.Email, req.Username).Scan(&existingID)
	if err == nil {
		return nil, fmt.Errorf("user with this email or username already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	userID := uuid.New()
	now := time.Now()

	user := models.User{
		ID:           userID,
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		IsActive:     true,
		IsAdmin:      false,
		IsVerified:   false,
		CreatedAt:    now,
		UpdatedAt:    now,
		LastLoginAt:  nil,
	}

	query := `
		INSERT INTO users (id, username, email, password, first_name, last_name,
		                  is_active, is_admin, is_verified, created_at, updated_at, last_login_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err = as.db.Exec(query,
		user.ID, user.Username, user.Email, user.PasswordHash,
		user.FirstName, user.LastName, user.IsActive, user.IsAdmin, user.IsVerified,
		user.CreatedAt, user.UpdatedAt, user.LastLoginAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Assign default role
	if err := as.assignDefaultRole(userID); err != nil {
		log.Printf("Failed to assign default role to user %s: %v", userID, err)
		// Don't fail registration for this
	}

	// Invalidate cache
	if as.cache != nil {
		as.cache.Delete(fmt.Sprintf("user:%s", userID))
	}

	response := user.ToResponse()
	return &response, nil
}

// Login authenticates a user
func (as *AuthService) Login(credentials UserCredentials) (*LoginResponse, error) {
	// Get user by email
	var user models.User
	var lockedUntil *time.Time

	query := `
		SELECT id, username, email, password_hash, first_name, last_name,
		       is_active, is_verified, two_factor_enabled, two_factor_secret,
		       failed_login_attempts, locked_until, last_login_at
		FROM users WHERE email = $1
	`

	err := as.db.QueryRow(query, credentials.Email).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.FirstName, &user.LastName, &user.IsActive, &user.IsVerified,
		&user.TwoFactorEnabled, &user.TwoFactorSecret,
		&user.FailedLoginAttempts, &lockedUntil, &user.LastLoginAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("invalid credentials")
	}
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Check if account is active
	if !user.IsActive {
		return nil, fmt.Errorf("account is disabled")
	}

	// Check if account is locked
	if lockedUntil != nil && time.Now().Before(*lockedUntil) {
		return nil, fmt.Errorf("account is temporarily locked due to too many failed attempts")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(credentials.Password))
	if err != nil {
		// Increment failed login attempts
		as.incrementFailedLoginAttempts(user.ID)
		return nil, fmt.Errorf("invalid credentials")
	}

	// Check 2FA if enabled
	if user.TwoFactorEnabled {
		if credentials.TOTPCode == "" {
			return &LoginResponse{RequiresTOTP: true}, nil
		}

		if !totp.Validate(credentials.TOTPCode, user.TwoFactorSecret) {
			as.incrementFailedLoginAttempts(user.ID)
			return nil, fmt.Errorf("invalid TOTP code")
		}
	}

	// Reset failed login attempts on successful login
	as.resetFailedLoginAttempts(user.ID)

	// Update last login
	as.updateLastLogin(user.ID)

	// Generate tokens
	accessToken, err := as.generateAccessToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := as.generateRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	userResponse := user.ToResponse()

	return &LoginResponse{
		User:         userResponse,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600, // 1 hour
		RequiresTOTP: false,
	}, nil
}

// VerifyToken verifies and decodes a JWT token
func (as *AuthService) VerifyToken(tokenString string) (*models.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return as.jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid token claims")
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return nil, fmt.Errorf("invalid user ID in token")
		}

		// Get user from database or cache
		user, err := as.GetUserByID(userID)
		if err != nil {
			return nil, fmt.Errorf("user not found: %w", err)
		}

		return user, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// GetUserByID retrieves a user by ID with caching
func (as *AuthService) GetUserByID(userID uuid.UUID) (*models.User, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("user:%s", userID)
	if as.cache != nil {
		var user models.User
		if err := as.cache.Get(cacheKey, &user); err == nil {
			return &user, nil
		}
	}

	// Get from database
	var user models.User
	query := `
		SELECT id, username, email, first_name, last_name, avatar_url, bio,
		       timezone, language, is_active, is_verified, email_verified_at,
		       last_login_at, created_at, updated_at, two_factor_enabled
		FROM users WHERE id = $1
	`

	err := as.db.QueryRow(query, userID).Scan(
		&user.ID, &user.Username, &user.Email, &user.FirstName, &user.LastName,
		&user.AvatarURL, &user.Bio, &user.Timezone, &user.Language,
		&user.IsActive, &user.IsVerified, &user.EmailVerifiedAt,
		&user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt, &user.TwoFactorEnabled,
	)

	if err != nil {
		return nil, err
	}

	// Cache user data
	if as.cache != nil {
		as.cache.Set(cacheKey, user, time.Hour)
	}

	return &user, nil
}

// EnableTOTP enables two-factor authentication for a user
func (as *AuthService) EnableTOTP(userID uuid.UUID) (*TOTPSecret, error) {
	// Generate TOTP secret
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "TPT Titan",
		AccountName: fmt.Sprintf("user-%s", userID.String()[:8]),
		SecretSize:  32,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to generate TOTP secret: %w", err)
	}

	// Store secret in database (temporary, until verified)
	query := "UPDATE users SET two_factor_secret = $1 WHERE id = $2"
	_, err = as.db.Exec(query, key.Secret(), userID)
	if err != nil {
		return nil, fmt.Errorf("failed to store TOTP secret: %w", err)
	}

	return &TOTPSecret{
		Secret:    key.Secret(),
		QRCodeURL: key.URL(),
	}, nil
}

// VerifyAndEnableTOTP verifies TOTP code and enables 2FA
func (as *AuthService) VerifyAndEnableTOTP(userID uuid.UUID, code string) error {
	// Get user's TOTP secret
	var secret string
	err := as.db.QueryRow("SELECT two_factor_secret FROM users WHERE id = $1", userID).Scan(&secret)
	if err != nil {
		return fmt.Errorf("TOTP not set up for user")
	}

	// Verify the code
	if !totp.Validate(code, secret) {
		return fmt.Errorf("invalid TOTP code")
	}

	// Enable 2FA
	query := "UPDATE users SET two_factor_enabled = true WHERE id = $1"
	_, err = as.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to enable 2FA: %w", err)
	}

	// Invalidate cache
	if as.cache != nil {
		as.cache.Delete(fmt.Sprintf("user:%s", userID))
	}

	return nil
}

// DisableTOTP disables two-factor authentication
func (as *AuthService) DisableTOTP(userID uuid.UUID, password string) error {
	// Verify password first
	var passwordHash string
	err := as.db.QueryRow("SELECT password_hash FROM users WHERE id = $1", userID).Scan(&passwordHash)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		return fmt.Errorf("invalid password")
	}

	// Disable 2FA
	query := "UPDATE users SET two_factor_enabled = false, two_factor_secret = NULL WHERE id = $1"
	_, err = as.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to disable 2FA: %w", err)
	}

	// Invalidate cache
	if as.cache != nil {
		as.cache.Delete(fmt.Sprintf("user:%s", userID))
	}

	return nil
}

// RequestPasswordReset initiates password reset process
func (as *AuthService) RequestPasswordReset(email string) error {
	// Generate reset token
	token := generateSecureToken()
	expiresAt := time.Now().Add(1 * time.Hour)

	// Store reset token
	query := `
		INSERT INTO password_resets (email, token, expires_at, created_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (email) DO UPDATE SET
			token = EXCLUDED.token,
			expires_at = EXCLUDED.expires_at,
			created_at = EXCLUDED.created_at
	`

	_, err := as.db.Exec(query, email, token, expiresAt, time.Now())
	if err != nil {
		return fmt.Errorf("failed to create password reset token: %w", err)
	}

	// TODO: Send email with reset link
	log.Printf("Password reset requested for email: %s", email)

	return nil
}

// ResetPassword resets user password using reset token
func (as *AuthService) ResetPassword(token, newPassword string) error {
	// Verify token
	var email string
	var expiresAt time.Time

	query := `
		SELECT email, expires_at FROM password_resets
		WHERE token = $1 AND expires_at > $2
	`

	err := as.db.QueryRow(query, token, time.Now()).Scan(&email, &expiresAt)
	if err != nil {
		return fmt.Errorf("invalid or expired reset token")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password
	query = `
		UPDATE users SET password_hash = $1, password_changed_at = $2
		WHERE email = $3
	`

	_, err = as.db.Exec(query, string(hashedPassword), time.Now(), email)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Delete used token
	as.db.Exec("DELETE FROM password_resets WHERE email = $1", email)

	// Invalidate cache
	if as.cache != nil {
		as.cache.Delete(fmt.Sprintf("user:email:%s", email))
	}

	return nil
}

// ChangePassword changes user password (authenticated)
func (as *AuthService) ChangePassword(userID uuid.UUID, oldPassword, newPassword string) error {
	// Verify old password
	var passwordHash string
	err := as.db.QueryRow("SELECT password_hash FROM users WHERE id = $1", userID).Scan(&passwordHash)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(oldPassword))
	if err != nil {
		return fmt.Errorf("invalid current password")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password
	query := `
		UPDATE users SET password_hash = $1, password_changed_at = $2
		WHERE id = $3
	`

	_, err = as.db.Exec(query, string(hashedPassword), time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Invalidate cache
	if as.cache != nil {
		as.cache.Delete(fmt.Sprintf("user:%s", userID))
	}

	return nil
}

// UpdateProfile updates user profile information
func (as *AuthService) UpdateProfile(userID uuid.UUID, updates map[string]interface{}) error {
	// Build dynamic update query
	setParts := []string{}
	args := []interface{}{}
	argCount := 1

	if firstName, ok := updates["first_name"].(string); ok {
		setParts = append(setParts, fmt.Sprintf("first_name = $%d", argCount))
		args = append(args, firstName)
		argCount++
	}

	if lastName, ok := updates["last_name"].(string); ok {
		setParts = append(setParts, fmt.Sprintf("last_name = $%d", argCount))
		args = append(args, lastName)
		argCount++
	}

	if bio, ok := updates["bio"].(string); ok {
		setParts = append(setParts, fmt.Sprintf("bio = $%d", argCount))
		args = append(args, bio)
		argCount++
	}

	if timezone, ok := updates["timezone"].(string); ok {
		setParts = append(setParts, fmt.Sprintf("timezone = $%d", argCount))
		args = append(args, timezone)
		argCount++
	}

	if language, ok := updates["language"].(string); ok {
		setParts = append(setParts, fmt.Sprintf("language = $%d", argCount))
		args = append(args, language)
		argCount++
	}

	if len(setParts) == 0 {
		return fmt.Errorf("no valid updates provided")
	}

	setParts = append(setParts, "updated_at = CURRENT_TIMESTAMP")
	query := fmt.Sprintf("UPDATE users SET %s WHERE id = $%d",
		strings.Join(setParts, ", "), argCount)
	args = append(args, userID)

	_, err := as.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}

	// Invalidate cache
	if as.cache != nil {
		as.cache.Delete(fmt.Sprintf("user:%s", userID))
	}

	return nil
}

// Logout invalidates user session
func (as *AuthService) Logout(tokenString string) error {
	// For stateless JWT, we could add token to blacklist
	// For now, just log the logout
	log.Printf("User logged out with token: %s...", tokenString[:20])
	return nil
}

// Helper methods

func (as *AuthService) assignDefaultRole(userID uuid.UUID) error {
	// Assign 'user' role by default
	query := `
		INSERT INTO user_roles (id, user_id, role_id, assigned_at)
		SELECT $1, $2, id, $3 FROM roles WHERE name = 'user'
	`
	_, err := as.db.Exec(query, uuid.New(), userID, time.Now())
	return err
}

func (as *AuthService) generateAccessToken(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"type":    "access",
		"exp":     time.Now().Add(time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(as.jwtSecret)
}

func (as *AuthService) generateRefreshToken(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"type":    "refresh",
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(), // 7 days
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(as.jwtSecret)
}

func (as *AuthService) incrementFailedLoginAttempts(userID uuid.UUID) {
	query := `
		UPDATE users SET failed_login_attempts = failed_login_attempts + 1,
		locked_until = CASE
			WHEN failed_login_attempts >= 5 THEN CURRENT_TIMESTAMP + INTERVAL '15 minutes'
			ELSE NULL
		END
		WHERE id = $1
	`
	as.db.Exec(query, userID)
}

func (as *AuthService) resetFailedLoginAttempts(userID uuid.UUID) {
	query := "UPDATE users SET failed_login_attempts = 0, locked_until = NULL WHERE id = $1"
	as.db.Exec(query, userID)
}

func (as *AuthService) updateLastLogin(userID uuid.UUID) {
	query := "UPDATE users SET last_login_at = CURRENT_TIMESTAMP WHERE id = $1"
	as.db.Exec(query, userID)
}

func generateSecureToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return base32.StdEncoding.EncodeToString(bytes)
}

// Middleware

// AuthMiddleware validates JWT tokens
func (as *AuthService) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract token
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		// Verify token
		user, err := as.VerifyToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Set user in context
		c.Set("user_id", user.ID)
		c.Set("user", user)
		c.Next()
	}
}

// RoleMiddleware checks if user has required role
func (as *AuthService) RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		// Check if user has required role
		var hasRole bool
		query := `
			SELECT EXISTS(
				SELECT 1 FROM user_roles ur
				JOIN roles r ON ur.role_id = r.id
				WHERE ur.user_id = $1 AND r.name = $2
			)
		`

		err := as.db.QueryRow(query, userID, requiredRole).Scan(&hasRole)
		if err != nil || !hasRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// PermissionMiddleware checks if user has required permission
func (as *AuthService) PermissionMiddleware(resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		// Check if user has required permission
		var hasPermission bool
		query := `
			SELECT EXISTS(
				SELECT 1 FROM user_roles ur
				JOIN role_permissions rp ON ur.role_id = rp.role_id
				JOIN permissions p ON rp.permission_id = p.id
				WHERE ur.user_id = $1 AND p.resource = $2 AND p.action = $3
			)
		`

		err := as.db.QueryRow(query, userID, resource, action).Scan(&hasPermission)
		if err != nil || !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}
