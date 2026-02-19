package services

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
)

// AIError represents different types of AI service errors
type AIErrorType int

const (
	AIErrorNetwork AIErrorType = iota
	AIErrorAuthentication
	AIErrorRateLimit
	AIErrorQuotaExceeded
	AIErrorInvalidInput
	AIErrorServiceUnavailable
	AIErrorTimeout
	AIErrorFallbackFailed
	AIErrorUnknown
)

// AIError provides detailed error information for AI operations
type AIError struct {
	Type        AIErrorType
	Service     string
	Operation   string
	UserID      uuid.UUID
	RequestID   uuid.UUID
	Message     string
	Details     map[string]interface{}
	Timestamp   time.Time
	Retryable   bool
	SuggestedAction string
}

func (e AIError) Error() string {
	return fmt.Sprintf("[%s] %s: %s", e.Service, e.Operation, e.Message)
}

// AIErrorHandler manages error handling and recovery for AI services
type AIErrorHandler struct {
	maxRetries       int
	baseRetryDelay   time.Duration
	maxRetryDelay    time.Duration
	circuitBreakers  map[string]*CircuitBreaker
	errorHistory     map[string][]AIError
}

// CircuitBreaker implements circuit breaker pattern for AI services
type CircuitBreaker struct {
	serviceName      string
	failureThreshold int
	resetTimeout     time.Duration
	failureCount     int
	lastFailureTime  time.Time
	state            CircuitState
}

type CircuitState int

const (
	CircuitClosed CircuitState = iota
	CircuitOpen
	CircuitHalfOpen
)

// NewAIErrorHandler creates a new AI error handler
func NewAIErrorHandler() *AIErrorHandler {
	return &AIErrorHandler{
		maxRetries:     3,
		baseRetryDelay: time.Second,
		maxRetryDelay:  time.Minute,
		circuitBreakers: make(map[string]*CircuitBreaker),
		errorHistory:   make(map[string][]AIError),
	}
}

// HandleError processes an AI error and determines recovery strategy
func (h *AIErrorHandler) HandleError(err error, service, operation string, userID, requestID uuid.UUID, details map[string]interface{}) *AIError {
	aiErr := h.classifyError(err, service, operation, userID, requestID, details)
	h.recordError(aiErr)

	// Update circuit breaker
	h.updateCircuitBreaker(service, aiErr)

	// Determine suggested action
	aiErr.SuggestedAction = h.getSuggestedAction(aiErr)

	return aiErr
}

// classifyError categorizes the error and determines if it's retryable
func (h *AIErrorHandler) classifyError(err error, service, operation string, userID, requestID uuid.UUID, details map[string]interface{}) *AIError {
	aiErr := &AIError{
		Service:   service,
		Operation: operation,
		UserID:    userID,
		RequestID: requestID,
		Message:   err.Error(),
		Details:   details,
		Timestamp: time.Now(),
		Type:      AIErrorUnknown,
		Retryable: false,
	}

	errMsg := err.Error()

	// Classify error type based on message patterns
	switch {
	case contains(errMsg, "connection", "timeout", "network"):
		aiErr.Type = AIErrorNetwork
		aiErr.Retryable = true
	case contains(errMsg, "unauthorized", "invalid api key", "authentication"):
		aiErr.Type = AIErrorAuthentication
		aiErr.Retryable = false
	case contains(errMsg, "rate limit", "too many requests"):
		aiErr.Type = AIErrorRateLimit
		aiErr.Retryable = true
	case contains(errMsg, "quota", "limit exceeded"):
		aiErr.Type = AIErrorQuotaExceeded
		aiErr.Retryable = false
	case contains(errMsg, "invalid input", "bad request"):
		aiErr.Type = AIErrorInvalidInput
		aiErr.Retryable = false
	case contains(errMsg, "service unavailable", "server error"):
		aiErr.Type = AIErrorServiceUnavailable
		aiErr.Retryable = true
	default:
		aiErr.Type = AIErrorUnknown
		aiErr.Retryable = false
	}

	return aiErr
}

// ShouldRetry determines if an operation should be retried
func (h *AIErrorHandler) ShouldRetry(aiErr *AIError, attempt int) bool {
	if !aiErr.Retryable || attempt >= h.maxRetries {
		return false
	}

	// Check circuit breaker
	if cb, exists := h.circuitBreakers[aiErr.Service]; exists && cb.state == CircuitOpen {
		return false
	}

	return true
}

// GetRetryDelay calculates the delay before retrying
func (h *AIErrorHandler) GetRetryDelay(attempt int, aiErr *AIError) time.Duration {
	// Exponential backoff with jitter
	delay := time.Duration(attempt) * h.baseRetryDelay
	if delay > h.maxRetryDelay {
		delay = h.maxRetryDelay
	}

	// Add jitter to prevent thundering herd
	jitter := time.Duration(rand.Int63n(int64(delay/10)))
	delay += jitter

	// Special handling for rate limits
	if aiErr.Type == AIErrorRateLimit {
		// Rate limits often specify retry-after header
		if retryAfter, ok := aiErr.Details["retry_after"]; ok {
			if ra, ok := retryAfter.(float64); ok {
				return time.Duration(ra) * time.Second
			}
		}
		// Default to longer delay for rate limits
		delay *= 2
	}

	return delay
}

// ExecuteWithRetry executes an AI operation with automatic retry logic
func (h *AIErrorHandler) ExecuteWithRetry(operation func() error, service, operationName string, userID, requestID uuid.UUID) error {
	var lastErr error

	for attempt := 0; attempt <= h.maxRetries; attempt++ {
		err := operation()
		if err == nil {
			return nil // Success
		}

		lastErr = err
		aiErr := h.HandleError(err, service, operationName, userID, requestID, nil)

		if !h.ShouldRetry(aiErr, attempt) {
			break
		}

		delay := h.GetRetryDelay(attempt, aiErr)
		log.Printf("AI Error Handler: Retrying %s/%s in %v (attempt %d/%d)",
			service, operationName, delay, attempt+1, h.maxRetries+1)

		time.Sleep(delay)
	}

	return lastErr
}

// recordError stores error for analysis and monitoring
func (h *AIErrorHandler) recordError(aiErr *AIError) {
	key := fmt.Sprintf("%s:%s", aiErr.Service, aiErr.Operation)

	if h.errorHistory[key] == nil {
		h.errorHistory[key] = make([]AIError, 0)
	}

	// Keep only last 100 errors per operation
	history := h.errorHistory[key]
	if len(history) >= 100 {
		history = history[1:]
	}
	history = append(history, *aiErr)
	h.errorHistory[key] = history

	// Log significant errors
	if aiErr.Type == AIErrorAuthentication || aiErr.Type == AIErrorQuotaExceeded {
		log.Printf("AI Error Handler: Significant error - %s", aiErr.Error())
	}
}

// updateCircuitBreaker updates circuit breaker state based on errors
func (h *AIErrorHandler) updateCircuitBreaker(service string, aiErr *AIError) {
	cb, exists := h.circuitBreakers[service]
	if !exists {
		cb = &CircuitBreaker{
			serviceName:      service,
			failureThreshold: 5,
			resetTimeout:     time.Minute * 5,
			state:           CircuitClosed,
		}
		h.circuitBreakers[service] = cb
	}

	if aiErr.Retryable {
		cb.failureCount++
		cb.lastFailureTime = time.Now()

		if cb.failureCount >= cb.failureThreshold {
			cb.state = CircuitOpen
			log.Printf("AI Error Handler: Circuit breaker OPENED for service %s", service)
		}
	} else if cb.state == CircuitHalfOpen {
		// Successful call in half-open state
		cb.failureCount = 0
		cb.state = CircuitClosed
		log.Printf("AI Error Handler: Circuit breaker CLOSED for service %s", service)
	}
}

// getSuggestedAction provides user-friendly guidance for error resolution
func (h *AIErrorHandler) getSuggestedAction(aiErr *AIError) string {
	switch aiErr.Type {
	case AIErrorNetwork:
		return "Check your internet connection and try again. If the problem persists, try switching to local AI."
	case AIErrorAuthentication:
		return "Check your API key in AI Settings. Visit the provider's website to regenerate if needed."
	case AIErrorRateLimit:
		return "You've hit the rate limit. Wait a few minutes before trying again, or upgrade your plan."
	case AIErrorQuotaExceeded:
		return "You've exceeded your usage quota. Check your account limits or upgrade your plan."
	case AIErrorInvalidInput:
		return "Check your input data. Make sure it's in the correct format and within size limits."
	case AIErrorServiceUnavailable:
		return "The AI service is temporarily unavailable. Try again later or switch to a different provider."
	case AIErrorTimeout:
		return "The request timed out. Try with shorter input or switch to a faster AI provider."
	case AIErrorFallbackFailed:
		return "Both primary and fallback AI services failed. Check your settings and try again later."
	default:
		return "An unexpected error occurred. Check the troubleshooting guide or contact support."
	}
}

// GetErrorStats returns error statistics for monitoring
func (h *AIErrorHandler) GetErrorStats() map[string]interface{} {
	stats := make(map[string]interface{})

	for serviceOp, errors := range h.errorHistory {
		if len(errors) > 0 {
			recentErrors := 0
			for _, err := range errors {
				if time.Since(err.Timestamp) < time.Hour {
					recentErrors++
				}
			}

			stats[serviceOp] = map[string]interface{}{
				"total_errors":    len(errors),
				"recent_errors":   recentErrors,
				"last_error":      errors[len(errors)-1].Timestamp,
				"last_error_type": errors[len(errors)-1].Type,
			}
		}
	}

	// Circuit breaker states
	cbStats := make(map[string]string)
	for service, cb := range h.circuitBreakers {
		stateStr := "closed"
		if cb.state == CircuitOpen {
			stateStr = "open"
		} else if cb.state == CircuitHalfOpen {
			stateStr = "half-open"
		}
		cbStats[service] = stateStr
	}
	stats["circuit_breakers"] = cbStats

	return stats
}

// Helper function to check if string contains any of the given substrings
func contains(s string, substrings ...string) bool {
	for _, substr := range substrings {
		if strings.Contains(strings.ToLower(s), strings.ToLower(substr)) {
			return true
		}
	}
	return false
}
