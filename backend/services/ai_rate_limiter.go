package services

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// RateLimiter manages API call rates for AI services
type RateLimiter struct {
	mu             sync.RWMutex
	userLimits     map[string]*UserRateLimit
	serviceLimits  map[string]*ServiceRateLimit
	globalLimits   *GlobalRateLimit
	cleanupTicker  *time.Ticker
	stopCleanup    chan bool
}

// UserRateLimit tracks per-user rate limits
type UserRateLimit struct {
	UserID         uuid.UUID
	Service        string
	RequestCount   int
	WindowStart    time.Time
	WindowDuration time.Duration
	MaxRequests    int
}

// ServiceRateLimit tracks per-service rate limits
type ServiceRateLimit struct {
	Service        string
	RequestCount   int
	WindowStart    time.Time
	WindowDuration time.Duration
	MaxRequests    int
	BurstLimit     int // Allow bursts above the sustained rate
}

// GlobalRateLimit tracks system-wide rate limits
type GlobalRateLimit struct {
	RequestCount   int
	WindowStart    time.Time
	WindowDuration time.Duration
	MaxRequests    int
}

// RateLimitConfig configures rate limiting behavior
type RateLimitConfig struct {
	// Per-user limits
	UserMaxRequestsPerHour int
	UserMaxRequestsPerDay  int

	// Per-service limits
	ServiceMaxRequestsPerMinute int
	ServiceBurstLimit           int

	// Global limits
	GlobalMaxRequestsPerMinute int

	// Special limits for expensive operations
	ExpensiveOperationLimit int // Per hour for document analysis, etc.

	// Cleanup settings
	CleanupInterval time.Duration
}

// DefaultRateLimitConfig returns sensible defaults for rate limiting
func DefaultRateLimitConfig() *RateLimitConfig {
	return &RateLimitConfig{
		UserMaxRequestsPerHour:     1000, // 1000 requests per user per hour
		UserMaxRequestsPerDay:      5000, // 5000 requests per user per day
		ServiceMaxRequestsPerMinute: 500, // 500 requests per service per minute
		ServiceBurstLimit:           100, // Allow bursts up to 100 requests
		GlobalMaxRequestsPerMinute:  2000, // 2000 requests globally per minute
		ExpensiveOperationLimit:     50,  // 50 expensive operations per hour
		CleanupInterval:            time.Minute * 5,
	}
}

// NewRateLimiter creates a new rate limiter with the specified configuration
func NewRateLimiter(config *RateLimitConfig) *RateLimiter {
	if config == nil {
		config = DefaultRateLimitConfig()
	}

	limiter := &RateLimiter{
		userLimits:    make(map[string]*UserRateLimit),
		serviceLimits: make(map[string]*ServiceRateLimit),
		globalLimits: &GlobalRateLimit{
			RequestCount:   0,
			WindowStart:    time.Now(),
			WindowDuration: time.Minute,
			MaxRequests:    config.GlobalMaxRequestsPerMinute,
		},
		stopCleanup: make(chan bool),
	}

	// Initialize service-specific limits
	services := []string{"openrouter", "ollama", "elevenlabs", "assemblyai", "replicate"}
	for _, service := range services {
		limiter.serviceLimits[service] = &ServiceRateLimit{
			Service:        service,
			RequestCount:   0,
			WindowStart:    time.Now(),
			WindowDuration: time.Minute,
			MaxRequests:    config.ServiceMaxRequestsPerMinute,
			BurstLimit:     config.ServiceBurstLimit,
		}
	}

	// Start cleanup routine
	limiter.cleanupTicker = time.NewTicker(config.CleanupInterval)
	go limiter.cleanupRoutine()

	return limiter
}

// CheckUserLimit checks if a user is within their rate limits
func (rl *RateLimiter) CheckUserLimit(userID uuid.UUID, service string) (*RateLimitResult, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	userKey := fmt.Sprintf("%s:%s", userID.String(), service)

	// Get or create user limit tracker
	limit, exists := rl.userLimits[userKey]
	if !exists {
		limit = &UserRateLimit{
			UserID:         userID,
			Service:        service,
			RequestCount:   0,
			WindowStart:    time.Now(),
			WindowDuration: time.Hour,
			MaxRequests:    1000, // Default per hour
		}
		rl.userLimits[userKey] = limit
	}

	// Check hourly limit
	if time.Since(limit.WindowStart) >= limit.WindowDuration {
		// Reset window
		limit.WindowStart = time.Now()
		limit.RequestCount = 0
	}

	// Check if limit exceeded
	if limit.RequestCount >= limit.MaxRequests {
		resetTime := limit.WindowStart.Add(limit.WindowDuration)
		return &RateLimitResult{
			Allowed:      false,
			LimitType:    "user_hourly",
			Remaining:    0,
			ResetTime:    resetTime,
			RetryAfter:   time.Until(resetTime),
			ErrorMessage: fmt.Sprintf("User hourly limit exceeded for %s. Try again after %v.", service, time.Until(resetTime).Truncate(time.Second)),
		}, nil
	}

	return &RateLimitResult{
		Allowed:   true,
		LimitType: "user_hourly",
		Remaining: limit.MaxRequests - limit.RequestCount - 1,
		ResetTime: limit.WindowStart.Add(limit.WindowDuration),
	}, nil
}

// CheckServiceLimit checks if a service is within its rate limits
func (rl *RateLimiter) CheckServiceLimit(service string) (*RateLimitResult, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limit, exists := rl.serviceLimits[service]
	if !exists {
		return &RateLimitResult{
			Allowed:      false,
			ErrorMessage: fmt.Sprintf("Unknown service: %s", service),
		}, fmt.Errorf("unknown service: %s", service)
	}

	// Check if window needs reset
	if time.Since(limit.WindowStart) >= limit.WindowDuration {
		limit.WindowStart = time.Now()
		limit.RequestCount = 0
	}

	// Check burst limit first (allows short bursts above sustained rate)
	burstRemaining := limit.BurstLimit - (limit.RequestCount - limit.MaxRequests)
	if burstRemaining > 0 && limit.RequestCount < limit.MaxRequests+limit.BurstLimit {
		return &RateLimitResult{
			Allowed:   true,
			LimitType: "service_burst",
			Remaining: burstRemaining - 1,
			ResetTime: limit.WindowStart.Add(limit.WindowDuration),
		}, nil
	}

	// Check sustained limit
	if limit.RequestCount >= limit.MaxRequests {
		resetTime := limit.WindowStart.Add(limit.WindowDuration)
		return &RateLimitResult{
			Allowed:      false,
			LimitType:    "service_sustained",
			Remaining:    0,
			ResetTime:    resetTime,
			RetryAfter:   time.Until(resetTime),
			ErrorMessage: fmt.Sprintf("Service rate limit exceeded for %s. Try again after %v.", service, time.Until(resetTime).Truncate(time.Second)),
		}, nil
	}

	return &RateLimitResult{
		Allowed:   true,
		LimitType: "service_sustained",
		Remaining: limit.MaxRequests - limit.RequestCount - 1,
		ResetTime: limit.WindowStart.Add(limit.WindowDuration),
	}, nil
}

// CheckGlobalLimit checks if the global rate limit is exceeded
func (rl *RateLimiter) CheckGlobalLimit() (*RateLimitResult, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Check if window needs reset
	if time.Since(rl.globalLimits.WindowStart) >= rl.globalLimits.WindowDuration {
		rl.globalLimits.WindowStart = time.Now()
		rl.globalLimits.RequestCount = 0
	}

	// Check if limit exceeded
	if rl.globalLimits.RequestCount >= rl.globalLimits.MaxRequests {
		resetTime := rl.globalLimits.WindowStart.Add(rl.globalLimits.WindowDuration)
		return &RateLimitResult{
			Allowed:      false,
			LimitType:    "global",
			Remaining:    0,
			ResetTime:    resetTime,
			RetryAfter:   time.Until(resetTime),
			ErrorMessage: fmt.Sprintf("Global rate limit exceeded. Try again after %v.", time.Until(resetTime).Truncate(time.Second)),
		}, nil
	}

	return &RateLimitResult{
		Allowed:   true,
		LimitType: "global",
		Remaining: rl.globalLimits.MaxRequests - rl.globalLimits.RequestCount - 1,
		ResetTime: rl.globalLimits.WindowStart.Add(rl.globalLimits.WindowDuration),
	}, nil
}

// CheckAllLimits checks all applicable rate limits
func (rl *RateLimiter) CheckAllLimits(userID uuid.UUID, service string) (*RateLimitResult, error) {
	// Check global limit first
	globalResult, err := rl.CheckGlobalLimit()
	if err != nil {
		return nil, err
	}
	if !globalResult.Allowed {
		return globalResult, nil
	}

	// Check service limit
	serviceResult, err := rl.CheckServiceLimit(service)
	if err != nil {
		return nil, err
	}
	if !serviceResult.Allowed {
		return serviceResult, nil
	}

	// Check user limit
	userResult, err := rl.CheckUserLimit(userID, service)
	if err != nil {
		return nil, err
	}
	if !userResult.Allowed {
		return userResult, nil
	}

	// All limits passed
	return &RateLimitResult{
		Allowed:   true,
		LimitType: "all",
		Remaining: min(userResult.Remaining, serviceResult.Remaining, globalResult.Remaining),
		ResetTime: earliestReset(userResult.ResetTime, serviceResult.ResetTime, globalResult.ResetTime),
	}, nil
}

// RecordRequest records a successful request for rate limiting purposes
func (rl *RateLimiter) RecordRequest(userID uuid.UUID, service string) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Record user request
	userKey := fmt.Sprintf("%s:%s", userID.String(), service)
	if limit, exists := rl.userLimits[userKey]; exists {
		limit.RequestCount++
	}

	// Record service request
	if limit, exists := rl.serviceLimits[service]; exists {
		limit.RequestCount++
	}

	// Record global request
	rl.globalLimits.RequestCount++
}

// IsExpensiveOperation checks if an operation is considered expensive
func (rl *RateLimiter) IsExpensiveOperation(operation string) bool {
	expensiveOps := []string{
		"document_analysis",
		"document_summarization",
		"email_categorization",
		"speech_synthesis",
		"workflow_optimization",
	}

	for _, op := range expensiveOps {
		if operation == op {
			return true
		}
	}
	return false
}

// CheckExpensiveOperationLimit checks limits for expensive operations
func (rl *RateLimiter) CheckExpensiveOperationLimit(userID uuid.UUID) (*RateLimitResult, error) {
	// Use a special key for expensive operations
	return rl.CheckUserLimit(userID, "expensive_operations")
}

// GetRateLimitStats returns rate limiting statistics
func (rl *RateLimiter) GetRateLimitStats() map[string]interface{} {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	stats := make(map[string]interface{})

	// User limits stats
	userStats := make(map[string]interface{})
	for key, limit := range rl.userLimits {
		userStats[key] = map[string]interface{}{
			"requests":        limit.RequestCount,
			"max_requests":    limit.MaxRequests,
			"window_start":    limit.WindowStart,
			"window_duration": limit.WindowDuration,
			"utilization":     float64(limit.RequestCount) / float64(limit.MaxRequests) * 100,
		}
	}
	stats["user_limits"] = userStats

	// Service limits stats
	serviceStats := make(map[string]interface{})
	for service, limit := range rl.serviceLimits {
		serviceStats[service] = map[string]interface{}{
			"requests":        limit.RequestCount,
			"max_requests":    limit.MaxRequests,
			"burst_limit":     limit.BurstLimit,
			"window_start":    limit.WindowStart,
			"window_duration": limit.WindowDuration,
			"utilization":     float64(limit.RequestCount) / float64(limit.MaxRequests) * 100,
		}
	}
	stats["service_limits"] = serviceStats

	// Global limits stats
	stats["global_limits"] = map[string]interface{}{
		"requests":        rl.globalLimits.RequestCount,
		"max_requests":    rl.globalLimits.MaxRequests,
		"window_start":    rl.globalLimits.WindowStart,
		"window_duration": rl.globalLimits.WindowDuration,
		"utilization":     float64(rl.globalLimits.RequestCount) / float64(rl.globalLimits.MaxRequests) * 100,
	}

	return stats
}

// Stop gracefully shuts down the rate limiter
func (rl *RateLimiter) Stop() {
	rl.cleanupTicker.Stop()
	rl.stopCleanup <- true
}

// cleanupRoutine periodically cleans up old rate limit entries
func (rl *RateLimiter) cleanupRoutine() {
	for {
		select {
		case <-rl.cleanupTicker.C:
			rl.cleanupOldEntries()
		case <-rl.stopCleanup:
			return
		}
	}
}

// cleanupOldEntries removes old rate limit entries that are no longer relevant
func (rl *RateLimiter) cleanupOldEntries() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	// Clean up user limits older than 24 hours
	for key, limit := range rl.userLimits {
		if now.Sub(limit.WindowStart) > time.Hour*24 {
			delete(rl.userLimits, key)
		}
	}

	// Reset service limits if window has passed
	for _, limit := range rl.serviceLimits {
		if now.Sub(limit.WindowStart) >= limit.WindowDuration {
			limit.RequestCount = 0
			limit.WindowStart = now
		}
	}

	// Reset global limits if window has passed
	if now.Sub(rl.globalLimits.WindowStart) >= rl.globalLimits.WindowDuration {
		rl.globalLimits.RequestCount = 0
		rl.globalLimits.WindowStart = now
	}
}

// RateLimitResult contains the result of a rate limit check
type RateLimitResult struct {
	Allowed      bool          `json:"allowed"`
	LimitType    string        `json:"limit_type"`
	Remaining    int           `json:"remaining"`
	ResetTime    time.Time     `json:"reset_time"`
	RetryAfter   time.Duration `json:"retry_after"`
	ErrorMessage string        `json:"error_message,omitempty"`
}

// Helper functions

func min(values ...int) int {
	if len(values) == 0 {
		return 0
	}
	min := values[0]
	for _, v := range values[1:] {
		if v < min {
			min = v
		}
	}
	return min
}

func earliestReset(times ...time.Time) time.Time {
	if len(times) == 0 {
		return time.Time{}
	}
	earliest := times[0]
	for _, t := range times[1:] {
		if t.Before(earliest) {
			earliest = t
		}
	}
	return earliest
}

// RateLimitMiddleware provides middleware for automatic rate limiting
type RateLimitMiddleware struct {
	limiter *RateLimiter
}

// NewRateLimitMiddleware creates rate limiting middleware
func NewRateLimitMiddleware(limiter *RateLimiter) *RateLimitMiddleware {
	return &RateLimitMiddleware{limiter: limiter}
}

// ExecuteWithRateLimit executes an operation with automatic rate limiting
func (m *RateLimitMiddleware) ExecuteWithRateLimit(
	userID uuid.UUID,
	service, operation string,
	operationFunc func() (interface{}, error),
) (interface{}, error) {

	// Check rate limits
	result, err := m.limiter.CheckAllLimits(userID, service)
	if err != nil {
		return nil, err
	}

	if !result.Allowed {
		return nil, fmt.Errorf("rate limit exceeded: %s", result.ErrorMessage)
	}

	// Check expensive operation limits if applicable
	if m.limiter.IsExpensiveOperation(operation) {
		expensiveResult, err := m.limiter.CheckExpensiveOperationLimit(userID)
		if err != nil {
			return nil, err
		}
		if !expensiveResult.Allowed {
			return nil, fmt.Errorf("expensive operation limit exceeded: %s", expensiveResult.ErrorMessage)
		}
		m.limiter.RecordRequest(userID, "expensive_operations")
	}

	// Execute the operation
	resultInterface, err := operationFunc()
	if err != nil {
		return nil, err
	}

	// Record the successful request
	m.limiter.RecordRequest(userID, service)

	return resultInterface, nil
}
