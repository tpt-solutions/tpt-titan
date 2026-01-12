package services

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// AICache provides intelligent caching for AI service results
type AICache struct {
	mu          sync.RWMutex
	cache       map[string]*CacheEntry
	maxSize     int
	defaultTTL  time.Duration
	cleanupInterval time.Duration
	stopCleanup chan bool
}

// CacheEntry represents a cached AI result
type CacheEntry struct {
	Key        string
	UserID     uuid.UUID
	Service    string
	Operation  string
	InputHash  string
	Result     interface{}
	Timestamp  time.Time
	TTL        time.Duration
	AccessCount int
	LastAccess time.Time
	Size       int // Approximate size in bytes
}

// CacheConfig configures the AI cache behavior
type CacheConfig struct {
	MaxSize         int           // Maximum number of entries
	DefaultTTL      time.Duration // Default time-to-live
	CleanupInterval time.Duration // How often to clean up expired entries
	EnableStats     bool          // Whether to track usage statistics
}

// DefaultCacheConfig returns sensible defaults for AI caching
func DefaultCacheConfig() *CacheConfig {
	return &CacheConfig{
		MaxSize:         1000,
		DefaultTTL:      time.Hour * 24, // 24 hours
		CleanupInterval: time.Minute * 30,
		EnableStats:     true,
	}
}

// NewAICache creates a new AI cache with the specified configuration
func NewAICache(config *CacheConfig) *AICache {
	if config == nil {
		config = DefaultCacheConfig()
	}

	cache := &AICache{
		cache:           make(map[string]*CacheEntry),
		maxSize:         config.MaxSize,
		defaultTTL:      config.DefaultTTL,
		cleanupInterval: config.CleanupInterval,
		stopCleanup:     make(chan bool),
	}

	// Start cleanup goroutine
	go cache.cleanupRoutine()

	return cache
}

// Get retrieves a cached result if available and not expired
func (c *AICache) Get(userID uuid.UUID, service, operation, input string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	key := c.generateKey(userID, service, operation, input)

	entry, exists := c.cache[key]
	if !exists {
		return nil, false
	}

	// Check if expired
	if time.Since(entry.Timestamp) > entry.TTL {
		return nil, false
	}

	// Update access statistics
	entry.AccessCount++
	entry.LastAccess = time.Now()

	return entry.Result, true
}

// Set stores a result in the cache
func (c *AICache) Set(userID uuid.UUID, service, operation, input string, result interface{}, customTTL ...time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := c.generateKey(userID, service, operation, input)
	ttl := c.defaultTTL
	if len(customTTL) > 0 {
		ttl = customTTL[0]
	}

	// Estimate size (rough approximation)
	size := c.estimateSize(result)

	entry := &CacheEntry{
		Key:        key,
		UserID:     userID,
		Service:    service,
		Operation:  operation,
		InputHash:  c.hashInput(input),
		Result:     result,
		Timestamp:  time.Now(),
		TTL:        ttl,
		AccessCount: 0,
		LastAccess: time.Now(),
		Size:       size,
	}

	// Check if we need to evict entries
	if len(c.cache) >= c.maxSize {
		c.evictEntries(1) // Evict 1 entry to make room
	}

	c.cache[key] = entry
}

// Delete removes a specific entry from the cache
func (c *AICache) Delete(userID uuid.UUID, service, operation, input string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := c.generateKey(userID, service, operation, input)
	delete(c.cache, key)
}

// Clear removes all entries from the cache
func (c *AICache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache = make(map[string]*CacheEntry)
}

// ClearExpired removes all expired entries
func (c *AICache) ClearExpired() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	expired := 0
	for key, entry := range c.cache {
		if time.Since(entry.Timestamp) > entry.TTL {
			delete(c.cache, key)
			expired++
		}
	}
	return expired
}

// ClearUserEntries removes all entries for a specific user
func (c *AICache) ClearUserEntries(userID uuid.UUID) int {
	c.mu.Lock()
	defer c.mu.Unlock()

	cleared := 0
	for key, entry := range c.cache {
		if entry.UserID == userID {
			delete(c.cache, key)
			cleared++
		}
	}
	return cleared
}

// GetStats returns cache statistics
func (c *AICache) GetStats() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	totalEntries := len(c.cache)
	totalSize := 0
	totalAccesses := 0
	oldestEntry := time.Now()
	newestEntry := time.Time{}

	// Calculate statistics
	for _, entry := range c.cache {
		totalSize += entry.Size
		totalAccesses += entry.AccessCount

		if entry.Timestamp.Before(oldestEntry) {
			oldestEntry = entry.Timestamp
		}
		if entry.Timestamp.After(newestEntry) {
			newestEntry = entry.Timestamp
		}
	}

	// Service breakdown
	serviceStats := make(map[string]int)
	for _, entry := range c.cache {
		serviceStats[entry.Service]++
	}

	return map[string]interface{}{
		"total_entries":    totalEntries,
		"max_entries":      c.maxSize,
		"total_size_kb":    totalSize / 1024,
		"total_accesses":   totalAccesses,
		"oldest_entry":     oldestEntry,
		"newest_entry":     newestEntry,
		"service_breakdown": serviceStats,
		"utilization_percent": (float64(totalEntries) / float64(c.maxSize)) * 100,
	}
}

// Stop gracefully shuts down the cache
func (c *AICache) Stop() {
	c.stopCleanup <- true
}

// cleanupRoutine runs in the background to clean up expired entries
func (c *AICache) cleanupRoutine() {
	ticker := time.NewTicker(c.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			expired := c.ClearExpired()
			if expired > 0 {
				// Log cleanup activity (would integrate with logging system)
				fmt.Printf("AI Cache: Cleaned up %d expired entries\n", expired)
			}
		case <-c.stopCleanup:
			return
		}
	}
}

// evictEntries removes least recently used entries when cache is full
func (c *AICache) evictEntries(count int) {
	type evictionCandidate struct {
		key   string
		score float64 // Lower score = more likely to be evicted
	}

	candidates := make([]evictionCandidate, 0, len(c.cache))

	for key, entry := range c.cache {
		// Calculate eviction score based on:
		// - Age (newer entries get higher score)
		// - Access count (more accessed entries get higher score)
		// - Size (smaller entries get slightly higher score)
		ageHours := time.Since(entry.Timestamp).Hours()
		accessScore := float64(entry.AccessCount) * 10
		sizePenalty := float64(entry.Size) / 1000 // Penalty for large entries

		score := accessScore - ageHours - sizePenalty

		candidates = append(candidates, evictionCandidate{key: key, score: score})
	}

	// Sort by score (ascending - lowest scores first)
	for i := 0; i < len(candidates)-1; i++ {
		for j := i + 1; j < len(candidates); j++ {
			if candidates[i].score > candidates[j].score {
				candidates[i], candidates[j] = candidates[j], candidates[i]
			}
		}
	}

	// Evict the lowest scoring entries
	for i := 0; i < count && i < len(candidates); i++ {
		delete(c.cache, candidates[i].key)
	}
}

// generateKey creates a unique cache key
func (c *AICache) generateKey(userID uuid.UUID, service, operation, input string) string {
	inputHash := c.hashInput(input)
	return fmt.Sprintf("%s:%s:%s:%s", userID.String(), service, operation, inputHash)
}

// hashInput creates an MD5 hash of the input for consistent caching
func (c *AICache) hashInput(input string) string {
	hasher := md5.New()
	hasher.Write([]byte(input))
	return hex.EncodeToString(hasher.Sum(nil))
}

// estimateSize provides a rough estimate of the result size in bytes
func (c *AICache) estimateSize(result interface{}) int {
	if result == nil {
		return 0
	}

	// Try to serialize to JSON and measure length
	if jsonBytes, err := json.Marshal(result); err == nil {
		return len(jsonBytes)
	}

	// Fallback size estimation
	switch v := result.(type) {
	case string:
		return len(v) * 2 // Assume UTF-16 encoding
	case []byte:
		return len(v)
	case map[string]interface{}:
		size := 0
		for k, val := range v {
			size += len(k) + c.estimateSize(val)
		}
		return size
	case []interface{}:
		size := 0
		for _, item := range v {
			size += c.estimateSize(item)
		}
		return size
	default:
		// Default assumption for unknown types
		return 1000
	}
}

// CacheMiddleware provides middleware for automatic caching
type CacheMiddleware struct {
	cache *AICache
}

// NewCacheMiddleware creates caching middleware
func NewCacheMiddleware(cache *AICache) *CacheMiddleware {
	return &CacheMiddleware{cache: cache}
}

// ExecuteWithCache executes an operation with automatic caching
func (m *CacheMiddleware) ExecuteWithCache(
	userID uuid.UUID,
	service, operation, input string,
	operationFunc func() (interface{}, error),
	customTTL ...time.Duration,
) (interface{}, error) {

	// Try to get from cache first
	if cached, found := m.cache.Get(userID, service, operation, input); found {
		return cached, nil
	}

	// Execute the operation
	result, err := operationFunc()
	if err != nil {
		return nil, err
	}

	// Cache the result
	m.cache.Set(userID, service, operation, input, result, customTTL...)

	return result, nil
}

// Specialized caching methods for different AI operations

// CacheWritingSuggestions caches writing assistance results
func (m *CacheMiddleware) CacheWritingSuggestions(userID uuid.UUID, text, context string, fn func() (interface{}, error)) (interface{}, error) {
	// Writing suggestions can be cached for longer since they're relatively stable
	return m.ExecuteWithCache(userID, "ai", "writing_suggestions", text+context, fn, time.Hour*48)
}

// CacheDocumentAnalysis caches document analysis results
func (m *CacheMiddleware) CacheDocumentAnalysis(userID uuid.UUID, content string, fn func() (interface{}, error)) (interface{}, error) {
	// Document analysis can be cached for a moderate time
	return m.ExecuteWithCache(userID, "ai", "document_analysis", content, fn, time.Hour*12)
}

// CacheEmailCategorization caches email categorization results (shorter TTL due to dynamic nature)
func (m *CacheMiddleware) CacheEmailCategorization(userID uuid.UUID, emailsData string, fn func() (interface{}, error)) (interface{}, error) {
	// Email categorization changes frequently, so shorter cache time
	return m.ExecuteWithCache(userID, "ai", "email_categorization", emailsData, fn, time.Hour*2)
}

// CacheTaskPrioritization caches task priority predictions
func (m *CacheMiddleware) CacheTaskPrioritization(userID uuid.UUID, taskData string, fn func() (interface{}, error)) (interface{}, error) {
	// Task prioritization can be cached moderately
	return m.ExecuteWithCache(userID, "ai", "task_prioritization", taskData, fn, time.Hour*6)
}

// CacheSpeechSynthesis caches TTS results (longer TTL since audio generation is expensive)
func (m *CacheMiddleware) CacheSpeechSynthesis(userID uuid.UUID, text, voiceConfig string, fn func() (interface{}, error)) (interface{}, error) {
	// Speech synthesis results can be cached for a long time
	return m.ExecuteWithCache(userID, "speech", "synthesize", text+voiceConfig, fn, time.Hour*168) // 1 week
}

// CacheSpeechRecognition caches STT results (shorter TTL due to potential variations)
func (m *CacheMiddleware) CacheSpeechRecognition(userID uuid.UUID, audioFingerprint string, fn func() (interface{}, error)) (interface{}, error) {
	// Speech recognition results can be cached moderately
	return m.ExecuteWithCache(userID, "speech", "recognize", audioFingerprint, fn, time.Hour*24)
}
