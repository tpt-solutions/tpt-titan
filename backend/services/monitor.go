package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

// MonitoringService provides comprehensive application monitoring
type MonitoringService struct {
	startTime      time.Time
	requestCount   int64
	errorCount     int64
	requestMetrics sync.Map // map[string]*RequestMetrics
	systemMetrics  *SystemMetrics
	cacheService   *CacheService
	dbOptimizer    *DatabaseOptimizer
	db             *sql.DB
}

// RequestMetrics tracks metrics for specific endpoints
type RequestMetrics struct {
	Path            string        `json:"path"`
	Method          string        `json:"method"`
	RequestCount    int64         `json:"request_count"`
	ErrorCount      int64         `json:"error_count"`
	AvgResponseTime time.Duration `json:"avg_response_time"`
	MinResponseTime time.Duration `json:"min_response_time"`
	MaxResponseTime time.Duration `json:"max_response_time"`
	LastRequestTime time.Time     `json:"last_request_time"`
}

// SystemMetrics represents system resource usage
type SystemMetrics struct {
	CPUUsage       float64 `json:"cpu_usage_percent"`
	MemoryUsage    float64 `json:"memory_usage_percent"`
	MemoryUsed     uint64  `json:"memory_used_bytes"`
	MemoryTotal    uint64  `json:"memory_total_bytes"`
	DiskUsage      float64 `json:"disk_usage_percent"`
	DiskUsed       uint64  `json:"disk_used_bytes"`
	DiskTotal      uint64  `json:"disk_total_bytes"`
	NetworkRxBytes uint64  `json:"network_rx_bytes"`
	NetworkTxBytes uint64  `json:"network_tx_bytes"`
	Goroutines     int     `json:"goroutines"`
	CGOCalls       int64   `json:"cgo_calls"`
	HeapAlloc      uint64  `json:"heap_alloc_bytes"`
	HeapSys        uint64  `json:"heap_sys_bytes"`
	HeapObjects    uint64  `json:"heap_objects"`
	GCCycles       uint32  `json:"gc_cycles"`
}

// ApplicationMetrics represents application-level metrics
type ApplicationMetrics struct {
	Uptime              time.Duration `json:"uptime"`
	TotalRequests       int64         `json:"total_requests"`
	ErrorRate           float64       `json:"error_rate"`
	ActiveUsers         int64         `json:"active_users"`
	DatabaseConnections int           `json:"database_connections"`
	CacheHitRate        float64       `json:"cache_hit_rate"`
	CacheKeys           int64         `json:"cache_keys"`
	MemoryUsage         uint64        `json:"memory_usage_bytes"`
}

// HealthStatus represents the overall health of the system
type HealthStatus struct {
	Status       string            `json:"status"` // "healthy", "degraded", "unhealthy"
	Timestamp    time.Time         `json:"timestamp"`
	Services     map[string]string `json:"services"`
	ResponseTime time.Duration     `json:"response_time"`
}

// Alert represents a system alert
type Alert struct {
	ID         string     `json:"id"`
	Level      string     `json:"level"` // "info", "warning", "error", "critical"
	Title      string     `json:"title"`
	Message    string     `json:"message"`
	Service    string     `json:"service"`
	Timestamp  time.Time  `json:"timestamp"`
	Resolved   bool       `json:"resolved"`
	ResolvedAt *time.Time `json:"resolved_at,omitempty"`
}

// NewMonitoringService creates a new monitoring service
func NewMonitoringService(cacheService *CacheService, dbOptimizer *DatabaseOptimizer, db *sql.DB) *MonitoringService {
	ms := &MonitoringService{
		startTime:     time.Now(),
		systemMetrics: &SystemMetrics{},
		cacheService:  cacheService,
		dbOptimizer:   dbOptimizer,
		db:            db,
	}

	// Start background monitoring
	go ms.startMonitoring()

	return ms
}

// startMonitoring begins background monitoring tasks
func (ms *MonitoringService) startMonitoring() {
	// Update system metrics every 30 seconds
	systemTicker := time.NewTicker(30 * time.Second)
	defer systemTicker.Stop()

	// Check health every 60 seconds
	healthTicker := time.NewTicker(60 * time.Second)
	defer healthTicker.Stop()

	for {
		select {
		case <-systemTicker.C:
			ms.updateSystemMetrics()
		case <-healthTicker.C:
			ms.checkHealth()
		}
	}
}

// updateSystemMetrics updates current system resource usage
func (ms *MonitoringService) updateSystemMetrics() {
	// CPU usage
	if cpuPercent, err := cpu.Percent(time.Second, false); err == nil && len(cpuPercent) > 0 {
		ms.systemMetrics.CPUUsage = cpuPercent[0]
	}

	// Memory usage
	if memStats, err := mem.VirtualMemory(); err == nil {
		ms.systemMetrics.MemoryUsage = memStats.UsedPercent
		ms.systemMetrics.MemoryUsed = memStats.Used
		ms.systemMetrics.MemoryTotal = memStats.Total
	}

	// Disk usage
	if diskStats, err := disk.Usage("/"); err == nil {
		ms.systemMetrics.DiskUsage = diskStats.UsedPercent
		ms.systemMetrics.DiskUsed = diskStats.Used
		ms.systemMetrics.DiskTotal = diskStats.Total
	}

	// Network stats
	if netStats, err := net.IOCounters(false); err == nil && len(netStats) > 0 {
		ms.systemMetrics.NetworkRxBytes = netStats[0].BytesRecv
		ms.systemMetrics.NetworkTxBytes = netStats[0].BytesSent
	}

	// Go runtime stats
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	ms.systemMetrics.Goroutines = runtime.NumGoroutine()
	ms.systemMetrics.CGOCalls = runtime.NumCgoCall()
	ms.systemMetrics.HeapAlloc = m.HeapAlloc
	ms.systemMetrics.HeapSys = m.HeapSys
	ms.systemMetrics.HeapObjects = m.HeapObjects
	ms.systemMetrics.GCCycles = m.NumGC
}

// RecordRequest records a request for monitoring
func (ms *MonitoringService) RecordRequest(path, method string, duration time.Duration, statusCode int) {
	ms.requestCount++

	key := method + ":" + path
	metricsInterface, exists := ms.requestMetrics.Load(key)

	var metrics *RequestMetrics
	if !exists {
		metrics = &RequestMetrics{
			Path:   path,
			Method: method,
		}
	} else {
		metrics = metricsInterface.(*RequestMetrics)
	}

	metrics.RequestCount++
	metrics.LastRequestTime = time.Now()

	if statusCode >= 400 {
		metrics.ErrorCount++
		ms.errorCount++
	}

	// Update response time statistics
	if metrics.RequestCount == 1 {
		metrics.AvgResponseTime = duration
		metrics.MinResponseTime = duration
		metrics.MaxResponseTime = duration
	} else {
		// Simple moving average
		metrics.AvgResponseTime = (metrics.AvgResponseTime + duration) / 2

		if duration < metrics.MinResponseTime {
			metrics.MinResponseTime = duration
		}
		if duration > metrics.MaxResponseTime {
			metrics.MaxResponseTime = duration
		}
	}

	ms.requestMetrics.Store(key, metrics)
}

// GetApplicationMetrics returns comprehensive application metrics
func (ms *MonitoringService) GetApplicationMetrics() (*ApplicationMetrics, error) {
	metrics := &ApplicationMetrics{
		Uptime:        time.Since(ms.startTime),
		TotalRequests: ms.requestCount,
		ErrorRate:     0,
		MemoryUsage:   ms.systemMetrics.HeapAlloc,
	}

	if ms.requestCount > 0 {
		metrics.ErrorRate = float64(ms.errorCount) / float64(ms.requestCount)
	}

	// Get database connections
	if ms.dbOptimizer != nil {
		if dbMetrics, err := ms.dbOptimizer.GetDatabaseMetrics(); err == nil {
			metrics.DatabaseConnections = dbMetrics.TotalConnections
		}
	}

	// Get cache statistics
	if ms.cacheService != nil {
		if cacheStats, err := ms.cacheService.GetStats(); err == nil {
			metrics.CacheHitRate = cacheStats.HitRate
			metrics.CacheKeys = cacheStats.Keys
		}
	}

	// Count active users (logged in within the last 15 minutes)
	metrics.ActiveUsers = ms.countActiveUsers()

	return metrics, nil
}

// countActiveUsers counts users who logged in within the last 15 minutes.
func (ms *MonitoringService) countActiveUsers() int64 {
	if ms.db == nil {
		return 0
	}
	var count int64
	cutoff := time.Now().Add(-15 * time.Minute)
	err := ms.db.QueryRow(
		"SELECT COUNT(*) FROM users WHERE last_login_at IS NOT NULL AND last_login_at >= $1",
		cutoff,
	).Scan(&count)
	if err != nil {
		log.Printf("Failed to count active users: %v", err)
		return 0
	}
	return count
}

// GetSystemMetrics returns current system metrics
func (ms *MonitoringService) GetSystemMetrics() *SystemMetrics {
	return ms.systemMetrics
}

// GetRequestMetrics returns metrics for all endpoints
func (ms *MonitoringService) GetRequestMetrics() map[string]*RequestMetrics {
	metrics := make(map[string]*RequestMetrics)
	ms.requestMetrics.Range(func(key, value interface{}) bool {
		metrics[key.(string)] = value.(*RequestMetrics)
		return true
	})
	return metrics
}

// CheckHealth performs a comprehensive health check
func (ms *MonitoringService) CheckHealth() *HealthStatus {
	status := &HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now(),
		Services:  make(map[string]string),
	}

	start := time.Now()

	// Check database connectivity
	if ms.dbOptimizer != nil {
		if _, err := ms.dbOptimizer.GetDatabaseMetrics(); err != nil {
			status.Services["database"] = "unhealthy"
			status.Status = "degraded"
		} else {
			status.Services["database"] = "healthy"
		}
	}

	// Check cache connectivity
	if ms.cacheService != nil {
		if ms.cacheService.Exists("health_check") {
			status.Services["cache"] = "healthy"
		} else {
			status.Services["cache"] = "unhealthy"
			status.Status = "degraded"
		}
	}

	// Check system resources
	if ms.systemMetrics.CPUUsage > 90 {
		status.Services["cpu"] = "critical"
		status.Status = "unhealthy"
	} else if ms.systemMetrics.CPUUsage > 70 {
		status.Services["cpu"] = "warning"
		if status.Status == "healthy" {
			status.Status = "degraded"
		}
	} else {
		status.Services["cpu"] = "healthy"
	}

	if ms.systemMetrics.MemoryUsage > 90 {
		status.Services["memory"] = "critical"
		status.Status = "unhealthy"
	} else if ms.systemMetrics.MemoryUsage > 80 {
		status.Services["memory"] = "warning"
		if status.Status == "healthy" {
			status.Status = "degraded"
		}
	} else {
		status.Services["memory"] = "healthy"
	}

	// Check disk space
	if ms.systemMetrics.DiskUsage > 95 {
		status.Services["disk"] = "critical"
		status.Status = "unhealthy"
	} else if ms.systemMetrics.DiskUsage > 85 {
		status.Services["disk"] = "warning"
		if status.Status == "healthy" {
			status.Status = "degraded"
		}
	} else {
		status.Services["disk"] = "healthy"
	}

	status.ResponseTime = time.Since(start)

	// Overall status check
	allHealthy := true
	for _, serviceStatus := range status.Services {
		if serviceStatus != "healthy" {
			allHealthy = false
			break
		}
	}

	if !allHealthy && status.Status == "healthy" {
		status.Status = "degraded"
	}

	return status
}

// checkHealth is called periodically to update health status
func (ms *MonitoringService) checkHealth() {
	health := ms.CheckHealth()

	// Log warnings and errors
	if health.Status != "healthy" {
		log.Printf("Health check status: %s", health.Status)
		for service, status := range health.Services {
			if status != "healthy" {
				log.Printf("Service %s: %s", service, status)
			}
		}
	}

	// Cache health status
	if ms.cacheService != nil {
		healthData, _ := json.Marshal(health)
		ms.cacheService.Set("health_status", string(healthData), 5*time.Minute)
	}
}

// CreateAlert creates a new system alert
func (ms *MonitoringService) CreateAlert(level, title, message, service string) *Alert {
	alert := &Alert{
		ID:        fmt.Sprintf("alert_%d", time.Now().Unix()),
		Level:     level,
		Title:     title,
		Message:   message,
		Service:   service,
		Timestamp: time.Now(),
		Resolved:  false,
	}

	// In production, this would be stored in database and sent to alerting system
	log.Printf("[%s] %s: %s", alert.Level, alert.Title, alert.Message)

	return alert
}

// ResolveAlert marks an alert as resolved
func (ms *MonitoringService) ResolveAlert(alertID string) {
	// In production, update alert in database
	log.Printf("Alert %s resolved", alertID)
}

// GetAlerts returns recent alerts (simplified - would query database)
func (ms *MonitoringService) GetAlerts(limit int) []*Alert {
	// Mock alerts - in production, query database
	alerts := []*Alert{
		{
			ID:        "alert_1",
			Level:     "info",
			Title:     "System started",
			Message:   "TPT Titan backend started successfully",
			Service:   "system",
			Timestamp: ms.startTime,
			Resolved:  true,
		},
	}

	if limit > 0 && len(alerts) > limit {
		alerts = alerts[:limit]
	}

	return alerts
}

// Middleware returns Gin middleware for request monitoring
func (ms *MonitoringService) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Record metrics
		duration := time.Since(start)
		statusCode := c.Writer.Status()
		path := c.Request.URL.Path
		method := c.Request.Method

		ms.RecordRequest(path, method, duration, statusCode)

		// Add response time header
		c.Header("X-Response-Time", duration.String())
	}
}

// MetricsHandler returns an HTTP handler for metrics endpoint
func (ms *MonitoringService) MetricsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		metrics := make(map[string]interface{})

		// Application metrics
		appMetrics, err := ms.GetApplicationMetrics()
		if err == nil {
			metrics["application"] = appMetrics
		}

		// System metrics
		metrics["system"] = ms.GetSystemMetrics()

		// Request metrics
		metrics["requests"] = ms.GetRequestMetrics()

		// Health status
		metrics["health"] = ms.CheckHealth()

		// Database metrics
		if ms.dbOptimizer != nil {
			if dbMetrics, err := ms.dbOptimizer.GetDatabaseMetrics(); err == nil {
				metrics["database"] = dbMetrics
			}
		}

		// Cache metrics
		if ms.cacheService != nil {
			if cacheStats, err := ms.cacheService.GetStats(); err == nil {
				metrics["cache"] = cacheStats
			}
		}

		c.JSON(http.StatusOK, metrics)
	}
}

// HealthHandler returns an HTTP handler for health check endpoint
func (ms *MonitoringService) HealthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		health := ms.CheckHealth()

		statusCode := http.StatusOK
		if health.Status == "degraded" {
			statusCode = http.StatusOK // Still return 200 for degraded
		} else if health.Status == "unhealthy" {
			statusCode = http.StatusServiceUnavailable
		}

		c.JSON(statusCode, health)
	}
}

// PrometheusMetricsHandler returns Prometheus-compatible metrics
func (ms *MonitoringService) PrometheusMetricsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var output string

		// Application metrics
		appMetrics, _ := ms.GetApplicationMetrics()
		output += fmt.Sprintf("# HELP tpt_titan_uptime_seconds Time since application start\n")
		output += fmt.Sprintf("# TYPE tpt_titan_uptime_seconds gauge\n")
		output += fmt.Sprintf("tpt_titan_uptime_seconds %f\n", appMetrics.Uptime.Seconds())

		output += fmt.Sprintf("# HELP tpt_titan_requests_total Total number of requests\n")
		output += fmt.Sprintf("# TYPE tpt_titan_requests_total counter\n")
		output += fmt.Sprintf("tpt_titan_requests_total %d\n", appMetrics.TotalRequests)

		// System metrics
		sysMetrics := ms.GetSystemMetrics()
		output += fmt.Sprintf("# HELP tpt_titan_cpu_usage_percent CPU usage percentage\n")
		output += fmt.Sprintf("# TYPE tpt_titan_cpu_usage_percent gauge\n")
		output += fmt.Sprintf("tpt_titan_cpu_usage_percent %f\n", sysMetrics.CPUUsage)

		output += fmt.Sprintf("# HELP tpt_titan_memory_usage_percent Memory usage percentage\n")
		output += fmt.Sprintf("# TYPE tpt_titan_memory_usage_percent gauge\n")
		output += fmt.Sprintf("tpt_titan_memory_usage_percent %f\n", sysMetrics.MemoryUsage)

		output += fmt.Sprintf("# HELP tpt_titan_goroutines_number Number of goroutines\n")
		output += fmt.Sprintf("# TYPE tpt_titan_goroutines_number gauge\n")
		output += fmt.Sprintf("tpt_titan_goroutines_number %d\n", sysMetrics.Goroutines)

		c.Header("Content-Type", "text/plain; charset=utf-8")
		c.String(http.StatusOK, output)
	}
}

// LogSlowRequest logs requests that exceed a duration threshold
func (ms *MonitoringService) LogSlowRequest(path, method string, duration time.Duration, threshold time.Duration) {
	if duration > threshold {
		log.Printf("SLOW REQUEST: %s %s took %v", method, path, duration)
		ms.CreateAlert("warning", "Slow Request", fmt.Sprintf("%s %s took %v", method, path, duration), "api")
	}
}

// GetPerformanceReport generates a comprehensive performance report
func (ms *MonitoringService) GetPerformanceReport() (map[string]interface{}, error) {
	report := make(map[string]interface{})

	// Application metrics
	appMetrics, err := ms.GetApplicationMetrics()
	if err != nil {
		return nil, err
	}
	report["application"] = appMetrics

	// System metrics
	report["system"] = ms.GetSystemMetrics()

	// Request performance
	requestMetrics := ms.GetRequestMetrics()
	report["requests"] = requestMetrics

	// Database performance
	if ms.dbOptimizer != nil {
		if dbReport, err := ms.dbOptimizer.GenerateQueryReport(); err == nil {
			report["database"] = dbReport
		}
	}

	// Health status
	report["health"] = ms.CheckHealth()

	// Recent alerts
	report["alerts"] = ms.GetAlerts(10)

	report["generated_at"] = time.Now()

	return report, nil
}
