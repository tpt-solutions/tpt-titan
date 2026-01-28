package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/time/rate"
)

// SecurityMiddleware provides comprehensive security middleware
type SecurityMiddleware struct {
	auditLogger *AuditLogger
	rateLimiter *RateLimiter
	corsConfig  *CORSConfig
}

// AuditLogger handles security audit logging
type AuditLogger struct {
	db interface{} // Would be database connection
}

// RateLimiter manages API rate limiting
type RateLimiter struct {
	limiters map[string]*rate.Limiter
}

// CORSConfig manages CORS settings
type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

// SecurityEvent represents a security audit event
type SecurityEvent struct {
	ID          uuid.UUID `json:"id"`
	Timestamp   time.Time `json:"timestamp"`
	EventType   string    `json:"event_type"`
	UserID      *uuid.UUID `json:"user_id,omitempty"`
	IPAddress   string    `json:"ip_address"`
	UserAgent   string    `json:"user_agent"`
	Resource    string    `json:"resource"`
	Action      string    `json:"action"`
	StatusCode  int       `json:"status_code"`
	Details     string    `json:"details,omitempty"`
	RequestID   string    `json:"request_id"`
}

// NewSecurityMiddleware creates a new security middleware instance
func NewSecurityMiddleware() *SecurityMiddleware {
	return &SecurityMiddleware{
		auditLogger: NewAuditLogger(),
		rateLimiter: NewRateLimiter(),
		corsConfig:  NewCORSConfig(),
	}
}

// NewAuditLogger creates a new audit logger
func NewAuditLogger() *AuditLogger {
	return &AuditLogger{}
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
	}
}

// NewCORSConfig creates default CORS configuration
func NewCORSConfig() *CORSConfig {
	// Get allowed origins from environment, default to localhost for development
	originsStr := getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:5173")
	origins := splitOrigins(originsStr)
	
	return &CORSConfig{
		AllowedOrigins:   origins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "Accept", "Cache-Control", "X-Requested-With"},
		AllowCredentials: true,
		MaxAge:           86400, // 24 hours
	}
}

// splitOrigins splits a comma-separated string of origins
func splitOrigins(origins string) []string {
	if origins == "" {
		return []string{"http://localhost:3000"}
	}
	
	parts := strings.Split(origins, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// CORSMiddleware handles CORS requests
func (sm *SecurityMiddleware) CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Check if origin is allowed
		allowed := false
		for _, allowedOrigin := range sm.corsConfig.AllowedOrigins {
			if allowedOrigin == "*" || allowedOrigin == origin {
				allowed = true
				break
			}
		}

		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		c.Header("Access-Control-Allow-Methods", strings.Join(sm.corsConfig.AllowedMethods, ", "))
		c.Header("Access-Control-Allow-Headers", strings.Join(sm.corsConfig.AllowedHeaders, ", "))
		c.Header("Access-Control-Allow-Credentials", strconv.FormatBool(sm.corsConfig.AllowCredentials))
		c.Header("Access-Control-Max-Age", strconv.Itoa(sm.corsConfig.MaxAge))

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// RateLimitMiddleware implements API rate limiting
func (sm *SecurityMiddleware) RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get client identifier (IP address)
		clientIP := c.ClientIP()

		// Check rate limit
		if !sm.rateLimiter.Allow(clientIP) {
			sm.auditLogger.LogSecurityEvent(&SecurityEvent{
				EventType: "RATE_LIMIT_EXCEEDED",
				IPAddress: clientIP,
				UserAgent: c.Request.UserAgent(),
				Resource:  c.Request.URL.Path,
				Action:    c.Request.Method,
				StatusCode: http.StatusTooManyRequests,
				RequestID: c.GetString("request_id"),
			})

			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
				"message": "Too many requests. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Allow checks if the client is within rate limits
func (rl *RateLimiter) Allow(clientIP string) bool {
	limiter, exists := rl.limiters[clientIP]
	if !exists {
		// 100 requests per minute per IP
		limiter = rate.NewLimiter(rate.Limit(100.0/60.0), 100)
		rl.limiters[clientIP] = limiter
	}
	return limiter.Allow()
}

// SecurityHeadersMiddleware adds security headers
func (sm *SecurityMiddleware) SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevent clickjacking
		c.Header("X-Frame-Options", "DENY")

		// Prevent MIME type sniffing
		c.Header("X-Content-Type-Options", "nosniff")

		// Enable XSS protection
		c.Header("X-XSS-Protection", "1; mode=block")

		// Referrer policy
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// Content Security Policy (basic)
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' data:;")

		// HTTP Strict Transport Security (only for HTTPS)
		if c.Request.TLS != nil {
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}

		c.Next()
	}
}

// RequestIDMiddleware generates unique request IDs
func (sm *SecurityMiddleware) RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := generateRequestID()
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

// AuditMiddleware logs all API requests
func (sm *SecurityMiddleware) AuditMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Get user ID if authenticated
		var userID *uuid.UUID
		if userIDStr, exists := c.Get("user_id"); exists {
			if uid, ok := userIDStr.(uuid.UUID); ok {
				userID = &uid
			}
		}

		c.Next()

		// Log the request
		sm.auditLogger.LogSecurityEvent(&SecurityEvent{
			EventType: "API_ACCESS",
			UserID:    userID,
			IPAddress: c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
			Resource:  c.Request.URL.Path,
			Action:    c.Request.Method,
			StatusCode: c.Writer.Status(),
			Details:   fmt.Sprintf("Duration: %v", time.Since(start)),
			RequestID: c.GetString("request_id"),
		})
	}
}

// LogSecurityEvent logs a security event
func (al *AuditLogger) LogSecurityEvent(event *SecurityEvent) {
	event.ID = uuid.New()
	event.Timestamp = time.Now()

	// In production, this would write to database or secure log file
	log.Printf("[AUDIT] %s | User: %v | IP: %s | Resource: %s | Action: %s | Status: %d | Details: %s",
		event.EventType,
		event.UserID,
		event.IPAddress,
		event.Resource,
		event.Action,
		event.StatusCode,
		event.Details,
	)
}

// CSRFProtectionMiddleware provides CSRF protection
func (sm *SecurityMiddleware) CSRFProtectionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip CSRF check for safe methods
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" || c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		// Get CSRF token from header
		token := c.GetHeader("X-CSRF-Token")
		if token == "" {
			// Try to get from form data
			token = c.PostForm("csrf_token")
		}

		// In production, validate token against session/server-side storage
		if token == "" {
			sm.auditLogger.LogSecurityEvent(&SecurityEvent{
				EventType: "CSRF_TOKEN_MISSING",
				IPAddress: c.ClientIP(),
				UserAgent: c.Request.UserAgent(),
				Resource:  c.Request.URL.Path,
				Action:    c.Request.Method,
				StatusCode: http.StatusForbidden,
				RequestID: c.GetString("request_id"),
			})

			c.JSON(http.StatusForbidden, gin.H{
				"error": "CSRF token missing",
				"message": "Cross-Site Request Forgery token is required",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// InputValidationMiddleware provides basic input validation
func (sm *SecurityMiddleware) InputValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for suspicious patterns in request
		if sm.containsSuspiciousPatterns(c) {
			sm.auditLogger.LogSecurityEvent(&SecurityEvent{
				EventType: "SUSPICIOUS_INPUT",
				IPAddress: c.ClientIP(),
				UserAgent: c.Request.UserAgent(),
				Resource:  c.Request.URL.Path,
				Action:    c.Request.Method,
				StatusCode: http.StatusBadRequest,
				RequestID: c.GetString("request_id"),
			})

			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid input",
				"message": "Request contains suspicious content",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// containsSuspiciousPatterns checks for common attack patterns
func (sm *SecurityMiddleware) containsSuspiciousPatterns(c *gin.Context) bool {
	// Check URL for suspicious patterns
	suspiciousPatterns := []string{
		"<script", "</script>", "javascript:", "vbscript:",
		"onload=", "onerror=", "eval(", "alert(",
		"union select", "1=1", "--", "/*", "*/",
		"../", "..\\", "\\x", "%2e%2e", "%2f",
	}

	url := c.Request.URL.String()
	body := c.Request.Body

	for _, pattern := range suspiciousPatterns {
		if strings.Contains(strings.ToLower(url), strings.ToLower(pattern)) {
			return true
		}
	}

	// Note: In production, you'd also check request body
	// This is simplified for the example
	_ = body

	return false
}

// IPWhitelistMiddleware restricts access to whitelisted IPs
func (sm *SecurityMiddleware) IPWhitelistMiddleware(allowedIPs []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		// Check if IP is in whitelist
		allowed := false
		for _, allowedIP := range allowedIPs {
			if allowedIP == clientIP {
				allowed = true
				break
			}

			// Check CIDR ranges
			if strings.Contains(allowedIP, "/") {
				_, network, err := net.ParseCIDR(allowedIP)
				if err == nil && network.Contains(net.ParseIP(clientIP)) {
					allowed = true
					break
				}
			}
		}

		if !allowed {
			sm.auditLogger.LogSecurityEvent(&SecurityEvent{
				EventType: "IP_BLOCKED",
				IPAddress: clientIP,
				UserAgent: c.Request.UserAgent(),
				Resource:  c.Request.URL.Path,
				Action:    c.Request.Method,
				StatusCode: http.StatusForbidden,
				RequestID: c.GetString("request_id"),
			})

			c.JSON(http.StatusForbidden, gin.H{
				"error": "Access denied",
				"message": "IP address not allowed",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// generateRequestID generates a unique request identifier
func generateRequestID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// SQLInjectionProtectionMiddleware provides basic SQL injection protection
func (sm *SecurityMiddleware) SQLInjectionProtectionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// This is a basic implementation
		// In production, use prepared statements and ORM sanitization

		suspiciousSQLPatterns := []string{
			"union select", "union all select", "select * from",
			"drop table", "delete from", "update ", "insert into",
			"--", "/*", "*/", "xp_", "sp_", "exec",
		}

		for key, values := range c.Request.URL.Query() {
			for _, value := range values {
				for _, pattern := range suspiciousSQLPatterns {
					if strings.Contains(strings.ToLower(value), pattern) {
						sm.auditLogger.LogSecurityEvent(&SecurityEvent{
							EventType: "SQL_INJECTION_ATTEMPT",
							IPAddress: c.ClientIP(),
							UserAgent: c.Request.UserAgent(),
							Resource:  c.Request.URL.Path,
							Action:    c.Request.Method,
							StatusCode: http.StatusBadRequest,
							Details:   fmt.Sprintf("Parameter: %s", key),
							RequestID: c.GetString("request_id"),
						})

						c.JSON(http.StatusBadRequest, gin.H{
							"error": "Invalid input",
							"message": "Request blocked due to security policy",
						})
						c.Abort()
						return
					}
				}
			}
		}

		c.Next()
	}
}
