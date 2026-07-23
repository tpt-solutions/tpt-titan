package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

// CacheService provides Redis-based caching functionality
type CacheService struct {
	client *redis.Client
	ctx    context.Context
}

// CacheItem represents a cached item with metadata
type CacheItem struct {
	Key        string        `json:"key"`
	Value      interface{}   `json:"value"`
	TTL        time.Duration `json:"ttl"`
	CreatedAt  time.Time     `json:"created_at"`
	AccessedAt time.Time     `json:"accessed_at"`
	HitCount   int64         `json:"hit_count"`
}

// CacheStats represents cache statistics
type CacheStats struct {
	Hits        int64   `json:"hits"`
	Misses      int64   `json:"misses"`
	HitRate     float64 `json:"hit_rate"`
	Keys        int64   `json:"keys"`
	MemoryUsage int64   `json:"memory_usage"`
	Evictions   int64   `json:"evictions"`
}

// NewCacheService creates a new Redis cache service
func NewCacheService(redisURL string) (*CacheService, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	client := redis.NewClient(opt)
	ctx := context.Background()

	// Test connection
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &CacheService{
		client: client,
		ctx:    ctx,
	}, nil
}

// Set stores a value in cache with TTL
func (cs *CacheService) Set(key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return cs.client.Set(cs.ctx, key, data, ttl).Err()
}

// Get retrieves a value from cache
func (cs *CacheService) Get(key string, dest interface{}) error {
	data, err := cs.client.Get(cs.ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(data), dest)
}

// GetString retrieves a string value from cache
func (cs *CacheService) GetString(key string) (string, error) {
	return cs.client.Get(cs.ctx, key).Result()
}

// Delete removes a key from cache
func (cs *CacheService) Delete(key string) error {
	return cs.client.Del(cs.ctx, key).Err()
}

// Exists checks if a key exists in cache
func (cs *CacheService) Exists(key string) bool {
	count, err := cs.client.Exists(cs.ctx, key).Result()
	return err == nil && count > 0
}

// Expire sets expiration time for a key
func (cs *CacheService) Expire(key string, ttl time.Duration) error {
	return cs.client.Expire(cs.ctx, key, ttl).Err()
}

// TTL gets remaining TTL for a key
func (cs *CacheService) TTL(key string) (time.Duration, error) {
	return cs.client.TTL(cs.ctx, key).Result()
}

// Increment increments a numeric value
func (cs *CacheService) Increment(key string) (int64, error) {
	return cs.client.Incr(cs.ctx, key).Result()
}

// Decrement decrements a numeric value
func (cs *CacheService) Decrement(key string) (int64, error) {
	return cs.client.Decr(cs.ctx, key).Result()
}

// SetNX sets value only if key doesn't exist
func (cs *CacheService) SetNX(key string, value interface{}, ttl time.Duration) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, fmt.Errorf("failed to marshal value: %w", err)
	}

	return cs.client.SetNX(cs.ctx, key, data, ttl).Result()
}

// GetSet gets old value and sets new value
func (cs *CacheService) GetSet(key string, value interface{}) (string, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return "", fmt.Errorf("failed to marshal value: %w", err)
	}

	return cs.client.GetSet(cs.ctx, key, data).Result()
}

// Keys returns all keys matching pattern
func (cs *CacheService) Keys(pattern string) ([]string, error) {
	return cs.client.Keys(cs.ctx, pattern).Result()
}

// FlushDB clears all keys in current database
func (cs *CacheService) FlushDB() error {
	return cs.client.FlushDB(cs.ctx).Err()
}

// GetStats returns cache statistics
func (cs *CacheService) GetStats() (*CacheStats, error) {
	info, err := cs.client.Info(cs.ctx, "stats").Result()
	if err != nil {
		return nil, err
	}

	stats := &CacheStats{}

	// Parse Redis INFO stats
	lines := splitLines(info)
	for _, line := range lines {
		if len(line) == 0 || line[0] == '#' {
			continue
		}

		parts := splitKV(line)
		if len(parts) != 2 {
			continue
		}

		key, value := parts[0], parts[1]

		switch key {
		case "keyspace_hits":
			fmt.Sscanf(value, "%d", &stats.Hits)
		case "keyspace_misses":
			fmt.Sscanf(value, "%d", &stats.Misses)
		case "total_connections_received":
			// Can be used for connection stats
		case "evicted_keys":
			fmt.Sscanf(value, "%d", &stats.Evictions)
		}
	}

	// Calculate hit rate
	total := stats.Hits + stats.Misses
	if total > 0 {
		stats.HitRate = float64(stats.Hits) / float64(total)
	}

	// Get key count
	keys, err := cs.client.DBSize(cs.ctx).Result()
	if err == nil {
		stats.Keys = keys
	}

	// Get memory usage
	memInfo, err := cs.client.Info(cs.ctx, "memory").Result()
	if err == nil {
		memLines := splitLines(memInfo)
		for _, line := range memLines {
			if parts := splitKV(line); len(parts) == 2 && parts[0] == "used_memory" {
				fmt.Sscanf(parts[1], "%d", &stats.MemoryUsage)
				break
			}
		}
	}

	return stats, nil
}

// Cached executes a function with caching
func (cs *CacheService) Cached(key string, ttl time.Duration, fn func() (interface{}, error)) (interface{}, error) {
	// Try to get from cache first
	if cs.Exists(key) {
		var cached interface{}
		if err := cs.Get(key, &cached); err == nil {
			return cached, nil
		}
		// If unmarshal failed, continue to execute function
	}

	// Execute function
	result, err := fn()
	if err != nil {
		return nil, err
	}

	// Cache the result
	if cacheErr := cs.Set(key, result, ttl); cacheErr != nil {
		log.Printf("Failed to cache result for key %s: %v", key, cacheErr)
	}

	return result, nil
}

// PubSubSubscribe subscribes to a Redis pub/sub channel
func (cs *CacheService) PubSubSubscribe(channel string) *redis.PubSub {
	return cs.client.Subscribe(cs.ctx, channel)
}

// Publish publishes a message to a Redis pub/sub channel
func (cs *CacheService) Publish(channel string, message interface{}) error {
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	return cs.client.Publish(cs.ctx, channel, data).Err()
}

// Close closes the Redis connection
func (cs *CacheService) Close() error {
	return cs.client.Close()
}

// Helper functions

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i, r := range s {
		if r == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

func splitKV(s string) []string {
	for i, r := range s {
		if r == ':' {
			return []string{s[:i], s[i+1:]}
		}
	}
	return []string{s}
}

// Cache key generators for common patterns

func UserCacheKey(userID string) string {
	return fmt.Sprintf("user:%s", userID)
}

func DocumentCacheKey(docID string) string {
	return fmt.Sprintf("document:%s", docID)
}

func EmailCacheKey(emailID string) string {
	return fmt.Sprintf("email:%s", emailID)
}

func ChatRoomCacheKey(roomID string) string {
	return fmt.Sprintf("chat_room:%s", roomID)
}

func MeetingCacheKey(meetingID string) string {
	return fmt.Sprintf("meeting:%s", meetingID)
}

// Cache TTL constants
const (
	UserCacheTTL     = 30 * time.Minute
	DocumentCacheTTL = 15 * time.Minute
	EmailCacheTTL    = 10 * time.Minute
	ChatCacheTTL     = 5 * time.Minute
	MeetingCacheTTL  = 5 * time.Minute
	SessionCacheTTL  = 24 * time.Hour
	StatsCacheTTL    = 5 * time.Minute
)

// Cache warming functions

// WarmUserCache preloads user data into cache
func (cs *CacheService) WarmUserCache(userID string, userData interface{}) error {
	key := UserCacheKey(userID)
	return cs.Set(key, userData, UserCacheTTL)
}

// WarmDocumentCache preloads document data into cache
func (cs *CacheService) WarmDocumentCache(docID string, docData interface{}) error {
	key := DocumentCacheKey(docID)
	return cs.Set(key, docData, DocumentCacheTTL)
}

// Cache invalidation helpers

// InvalidateUserCache removes user data from cache
func (cs *CacheService) InvalidateUserCache(userID string) error {
	return cs.Delete(UserCacheKey(userID))
}

// InvalidateDocumentCache removes document data from cache
func (cs *CacheService) InvalidateDocumentCache(docID string) error {
	return cs.Delete(DocumentCacheKey(docID))
}

// InvalidateEmailCache removes email data from cache
func (cs *CacheService) InvalidateEmailCache(emailID string) error {
	return cs.Delete(EmailCacheKey(emailID))
}

// Batch operations

// MSet sets multiple key-value pairs
func (cs *CacheService) MSet(pairs map[string]interface{}) error {
	redisPairs := make(map[string]interface{})
	for key, value := range pairs {
		data, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value for key %s: %w", key, err)
		}
		redisPairs[key] = data
	}

	return cs.client.MSet(cs.ctx, redisPairs).Err()
}

// MGet gets multiple values by keys
func (cs *CacheService) MGet(keys ...string) ([]interface{}, error) {
	values, err := cs.client.MGet(cs.ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	results := make([]interface{}, len(values))
	for i, val := range values {
		if val == nil {
			results[i] = nil
			continue
		}

		var result interface{}
		if err := json.Unmarshal([]byte(val.(string)), &result); err != nil {
			return nil, fmt.Errorf("failed to unmarshal value for key %s: %w", keys[i], err)
		}
		results[i] = result
	}

	return results, nil
}

// Pipeline operations for better performance

// Pipeline creates a Redis pipeline for batch operations
func (cs *CacheService) Pipeline() redis.Pipeliner {
	return cs.client.Pipeline()
}

// Transaction creates a Redis transaction
func (cs *CacheService) Transaction() redis.Pipeliner {
	return cs.client.TxPipeline()
}
