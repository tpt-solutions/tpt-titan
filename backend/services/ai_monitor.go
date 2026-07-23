package services

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
)

// AIMonitor provides comprehensive monitoring and logging for AI services
type AIMonitor struct {
	mu             sync.RWMutex
	metrics        map[string]*AIMetric
	alerts         []*AIAlert
	alertCallbacks []AlertCallback
	logFile        *os.File
	logChan        chan *LogEntry
	stopChan       chan bool
}

// AIMetric represents a monitoring metric
type AIMetric struct {
	Name       string            `json:"name"`
	Type       MetricType        `json:"type"`
	Value      interface{}       `json:"value"`
	Labels     map[string]string `json:"labels"`
	Timestamp  time.Time         `json:"timestamp"`
	Samples    []MetricSample    `json:"samples,omitempty"`
	WindowSize time.Duration     `json:"window_size"`
}

// MetricType defines the type of metric
type MetricType int

const (
	MetricGauge MetricType = iota
	MetricCounter
	MetricHistogram
	MetricSummary
)

// MetricSample represents a single metric measurement
type MetricSample struct {
	Value     float64   `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}

// AIAlert represents an alert condition
type AIAlert struct {
	ID          uuid.UUID      `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Severity    AlertLevel     `json:"severity"`
	Condition   AlertCondition `json:"condition"`
	Status      AlertStatus    `json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	ResolvedAt  *time.Time     `json:"resolved_at,omitempty"`
	Value       interface{}    `json:"value"`
	Threshold   interface{}    `json:"threshold"`
}

// AlertLevel defines alert severity
type AlertLevel int

const (
	AlertInfo AlertLevel = iota
	AlertWarning
	AlertError
	AlertCritical
)

// AlertStatus represents alert state
type AlertStatus int

const (
	AlertActive AlertStatus = iota
	AlertResolved
	AlertAcknowledged
)

// AlertCondition defines when an alert should trigger
type AlertCondition struct {
	MetricName string        `json:"metric_name"`
	Operator   string        `json:"operator"` // ">", "<", "==", "!="
	Threshold  interface{}   `json:"threshold"`
	Duration   time.Duration `json:"duration"` // How long condition must persist
}

// AlertCallback is called when alerts are triggered
type AlertCallback func(alert *AIAlert)

// LogEntry represents a log entry
type LogEntry struct {
	Timestamp time.Time              `json:"timestamp"`
	Level     LogLevel               `json:"level"`
	Service   string                 `json:"service"`
	Operation string                 `json:"operation"`
	UserID    uuid.UUID              `json:"user_id,omitempty"`
	RequestID uuid.UUID              `json:"request_id,omitempty"`
	Message   string                 `json:"message"`
	Error     string                 `json:"error,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	Duration  time.Duration          `json:"duration,omitempty"`
}

// LogLevel defines logging levels
type LogLevel int

const (
	LogDebug LogLevel = iota
	LogInfo
	LogWarn
	LogError
	LogFatal
)

// AIMonitorConfig configures the AI monitor
type AIMonitorConfig struct {
	LogFile            string
	EnableMetrics      bool
	EnableAlerts       bool
	EnableLogging      bool
	MetricsInterval    time.Duration
	AlertCheckInterval time.Duration
	LogBufferSize      int
}

// DefaultAIMonitorConfig returns sensible defaults
func DefaultAIMonitorConfig() *AIMonitorConfig {
	return &AIMonitorConfig{
		LogFile:            "logs/ai_monitor.log",
		EnableMetrics:      true,
		EnableAlerts:       true,
		EnableLogging:      true,
		MetricsInterval:    time.Minute * 5,
		AlertCheckInterval: time.Minute * 1,
		LogBufferSize:      1000,
	}
}

// NewAIMonitor creates a new AI monitor
func NewAIMonitor(config *AIMonitorConfig) (*AIMonitor, error) {
	if config == nil {
		config = DefaultAIMonitorConfig()
	}

	monitor := &AIMonitor{
		metrics:        make(map[string]*AIMetric),
		alerts:         make([]*AIAlert, 0),
		alertCallbacks: make([]AlertCallback, 0),
		logChan:        make(chan *LogEntry, config.LogBufferSize),
		stopChan:       make(chan bool),
	}

	// Initialize log file
	if config.EnableLogging {
		if err := monitor.initLogFile(config.LogFile); err != nil {
			return nil, fmt.Errorf("failed to initialize log file: %w", err)
		}

		// Start log writer
		go monitor.logWriter()
	}

	// Initialize default metrics
	monitor.initDefaultMetrics()

	// Initialize default alerts
	monitor.initDefaultAlerts()

	// Start monitoring routines
	if config.EnableMetrics {
		go monitor.metricsRoutine(config.MetricsInterval)
	}

	if config.EnableAlerts {
		go monitor.alertsRoutine(config.AlertCheckInterval)
	}

	return monitor, nil
}

// RecordMetric records a metric value
func (m *AIMonitor) RecordMetric(name string, value interface{}, labels map[string]string, metricType MetricType) {
	m.mu.Lock()
	defer m.mu.Unlock()

	key := name
	if labels != nil {
		for k, v := range labels {
			key += fmt.Sprintf(":%s=%s", k, v)
		}
	}

	metric, exists := m.metrics[key]
	if !exists {
		windowSize := time.Hour // Default 1 hour window
		if name == "ai_response_time" {
			windowSize = time.Minute * 30
		}

		metric = &AIMetric{
			Name:       name,
			Type:       metricType,
			Labels:     labels,
			Samples:    make([]MetricSample, 0),
			WindowSize: windowSize,
		}
		m.metrics[key] = metric
	}

	// Update metric value
	metric.Value = value
	metric.Timestamp = time.Now()

	// Add sample for time-series data
	if metric.Type == MetricGauge || metric.Type == MetricCounter {
		sample := MetricSample{
			Value:     toFloat64(value),
			Timestamp: metric.Timestamp,
		}
		metric.Samples = append(metric.Samples, sample)

		// Keep only recent samples within window
		cutoff := metric.Timestamp.Add(-metric.WindowSize)
		filtered := make([]MetricSample, 0)
		for _, s := range metric.Samples {
			if s.Timestamp.After(cutoff) {
				filtered = append(filtered, s)
			}
		}
		metric.Samples = filtered
	}
}

// GetMetric retrieves a metric value
func (m *AIMonitor) GetMetric(name string, labels map[string]string) (*AIMetric, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	key := name
	if labels != nil {
		for k, v := range labels {
			key += fmt.Sprintf(":%s=%s", k, v)
		}
	}

	metric, exists := m.metrics[key]
	return metric, exists
}

// GetAllMetrics returns all metrics
func (m *AIMonitor) GetAllMetrics() map[string]*AIMetric {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Return a copy to prevent external modification
	result := make(map[string]*AIMetric)
	for k, v := range m.metrics {
		result[k] = v
	}
	return result
}

// AddAlertCallback registers a callback for alerts
func (m *AIMonitor) AddAlertCallback(callback AlertCallback) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.alertCallbacks = append(m.alertCallbacks, callback)
}

// Log records a log entry
func (m *AIMonitor) Log(level LogLevel, service, operation string, userID, requestID uuid.UUID, message string, err error, metadata map[string]interface{}, duration time.Duration) {
	entry := &LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Service:   service,
		Operation: operation,
		UserID:    userID,
		RequestID: requestID,
		Message:   message,
		Metadata:  metadata,
		Duration:  duration,
	}

	if err != nil {
		entry.Error = err.Error()
	}

	// Record metrics for the operation
	if duration > 0 {
		labels := map[string]string{
			"service":   service,
			"operation": operation,
		}
		m.RecordMetric("ai_response_time", duration.Milliseconds(), labels, MetricHistogram)
	}

	// Send to log channel (non-blocking)
	select {
	case m.logChan <- entry:
	default:
		// Log channel full, log to console as fallback
		log.Printf("[%s] %s/%s: %s", level.String(), service, operation, message)
	}
}

// GetActiveAlerts returns all active alerts
func (m *AIMonitor) GetActiveAlerts() []*AIAlert {
	m.mu.RLock()
	defer m.mu.RUnlock()

	active := make([]*AIAlert, 0)
	for _, alert := range m.alerts {
		if alert.Status == AlertActive {
			active = append(active, alert)
		}
	}
	return active
}

// AcknowledgeAlert acknowledges an alert
func (m *AIMonitor) AcknowledgeAlert(alertID uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, alert := range m.alerts {
		if alert.ID == alertID && alert.Status == AlertActive {
			alert.Status = AlertAcknowledged
			now := time.Now()
			alert.ResolvedAt = &now
			return nil
		}
	}

	return fmt.Errorf("alert not found or not active: %s", alertID)
}

// Stop gracefully shuts down the monitor
func (m *AIMonitor) Stop() {
	close(m.stopChan)

	if m.logFile != nil {
		m.logFile.Close()
	}
}

// initLogFile initializes the log file
func (m *AIMonitor) initLogFile(filename string) error {
	// Create logs directory if it doesn't exist
	os.MkdirAll("logs", 0755)

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	m.logFile = file
	return nil
}

// initDefaultMetrics initializes default metrics
func (m *AIMonitor) initDefaultMetrics() {
	// AI service metrics
	m.RecordMetric("ai_requests_total", 0, nil, MetricCounter)
	m.RecordMetric("ai_errors_total", 0, nil, MetricCounter)
	m.RecordMetric("ai_response_time", 0, map[string]string{"service": "all"}, MetricHistogram)

	// Cache metrics
	m.RecordMetric("ai_cache_hits", 0, nil, MetricCounter)
	m.RecordMetric("ai_cache_misses", 0, nil, MetricCounter)
	m.RecordMetric("ai_cache_size", 0, nil, MetricGauge)

	// Rate limiting metrics
	m.RecordMetric("ai_rate_limit_hits", 0, nil, MetricCounter)
	m.RecordMetric("ai_rate_limit_exceeded", 0, nil, MetricCounter)

	// Job queue metrics
	m.RecordMetric("ai_jobs_queued", 0, nil, MetricGauge)
	m.RecordMetric("ai_jobs_processing", 0, nil, MetricGauge)
	m.RecordMetric("ai_jobs_completed", 0, nil, MetricCounter)
	m.RecordMetric("ai_jobs_failed", 0, nil, MetricCounter)
}

// initDefaultAlerts initializes default alert conditions
func (m *AIMonitor) initDefaultAlerts() {
	alerts := []*AIAlert{
		{
			ID:          uuid.New(),
			Name:        "High Error Rate",
			Description: "AI service error rate is above 5%",
			Severity:    AlertError,
			Condition: AlertCondition{
				MetricName: "ai_errors_total",
				Operator:   ">",
				Threshold:  0.05,
				Duration:   time.Minute * 5,
			},
			Status:    AlertResolved,
			CreatedAt: time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        "Slow Response Time",
			Description: "Average AI response time exceeds 30 seconds",
			Severity:    AlertWarning,
			Condition: AlertCondition{
				MetricName: "ai_response_time",
				Operator:   ">",
				Threshold:  30000, // 30 seconds in milliseconds
				Duration:   time.Minute * 10,
			},
			Status:    AlertResolved,
			CreatedAt: time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        "High Rate Limit Usage",
			Description: "Rate limit usage exceeds 80%",
			Severity:    AlertWarning,
			Condition: AlertCondition{
				MetricName: "ai_rate_limit_exceeded",
				Operator:   ">",
				Threshold:  0.8,
				Duration:   time.Minute * 15,
			},
			Status:    AlertResolved,
			CreatedAt: time.Now(),
		},
	}

	m.alerts = alerts
}

// logWriter writes log entries to file
func (m *AIMonitor) logWriter() {
	for {
		select {
		case entry := <-m.logChan:
			m.writeLogEntry(entry)
		case <-m.stopChan:
			return
		}
	}
}

// writeLogEntry writes a log entry to file
func (m *AIMonitor) writeLogEntry(entry *LogEntry) {
	if m.logFile == nil {
		return
	}

	// Format log entry as JSON
	jsonData, err := json.Marshal(entry)
	if err != nil {
		log.Printf("Failed to marshal log entry: %v", err)
		return
	}

	// Write to file
	if _, err := m.logFile.WriteString(string(jsonData) + "\n"); err != nil {
		log.Printf("Failed to write log entry: %v", err)
	}

	// Also log to console for important messages
	if entry.Level >= LogError {
		log.Printf("[%s] %s/%s: %s", entry.Level.String(), entry.Service, entry.Operation, entry.Message)
	}
}

// metricsRoutine periodically updates computed metrics
func (m *AIMonitor) metricsRoutine(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.updateComputedMetrics()
		case <-m.stopChan:
			return
		}
	}
}

// updateComputedMetrics calculates derived metrics
func (m *AIMonitor) updateComputedMetrics() {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Calculate error rate
	if requestsMetric, exists := m.metrics["ai_requests_total"]; exists {
		if errorsMetric, exists := m.metrics["ai_errors_total"]; exists {
			requests := toFloat64(requestsMetric.Value)
			errors := toFloat64(errorsMetric.Value)

			if requests > 0 {
				errorRate := errors / requests
				m.RecordMetric("ai_error_rate", errorRate, nil, MetricGauge)
			}
		}
	}

	// Calculate cache hit rate
	if hitsMetric, exists := m.metrics["ai_cache_hits"]; exists {
		if missesMetric, exists := m.metrics["ai_cache_misses"]; exists {
			hits := toFloat64(hitsMetric.Value)
			misses := toFloat64(missesMetric.Value)
			total := hits + misses

			if total > 0 {
				hitRate := hits / total
				m.RecordMetric("ai_cache_hit_rate", hitRate, nil, MetricGauge)
			}
		}
	}
}

// alertsRoutine checks for alert conditions
func (m *AIMonitor) alertsRoutine(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.checkAlerts()
		case <-m.stopChan:
			return
		}
	}
}

// checkAlerts evaluates alert conditions
func (m *AIMonitor) checkAlerts() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, alert := range m.alerts {
		if alert.Status == AlertActive || alert.Status == AlertAcknowledged {
			continue // Already active or acknowledged
		}

		metric, exists := m.metrics[alert.Condition.MetricName]
		if !exists {
			continue
		}

		// Check if condition is met
		conditionMet := m.evaluateCondition(metric.Value, alert.Condition)

		if conditionMet {
			// Check if condition has persisted long enough
			if time.Since(alert.CreatedAt) >= alert.Condition.Duration {
				alert.Status = AlertActive
				alert.Value = metric.Value
				alert.Threshold = alert.Condition.Threshold

				// Notify callbacks
				for _, callback := range m.alertCallbacks {
					go callback(alert)
				}

				m.Log(LogWarn, "monitor", "alert_triggered", uuid.Nil, uuid.Nil,
					fmt.Sprintf("Alert triggered: %s", alert.Name), nil,
					map[string]interface{}{
						"alert_id":  alert.ID,
						"severity":  alert.Severity,
						"value":     alert.Value,
						"threshold": alert.Threshold,
					}, 0)
			}
		} else {
			// Reset alert timer if condition no longer met
			alert.CreatedAt = time.Now()
		}
	}
}

// evaluateCondition checks if a metric value meets an alert condition
func (m *AIMonitor) evaluateCondition(value interface{}, condition AlertCondition) bool {
	val := toFloat64(value)
	threshold := toFloat64(condition.Threshold)

	switch condition.Operator {
	case ">":
		return val > threshold
	case "<":
		return val < threshold
	case ">=":
		return val >= threshold
	case "<=":
		return val <= threshold
	case "==":
		return val == threshold
	case "!=":
		return val != threshold
	default:
		return false
	}
}

// Helper functions

func toFloat64(value interface{}) float64 {
	switch v := value.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int64:
		return float64(v)
	case int32:
		return float64(v)
	default:
		return 0.0
	}
}

func (l LogLevel) String() string {
	switch l {
	case LogDebug:
		return "DEBUG"
	case LogInfo:
		return "INFO"
	case LogWarn:
		return "WARN"
	case LogError:
		return "ERROR"
	case LogFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}
