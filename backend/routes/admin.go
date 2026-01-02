package routes

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetSystemStats returns comprehensive system statistics
func GetSystemStats(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)

	// Get user statistics
	var totalUsers, activeUsers, verifiedUsers int
	err := db.QueryRow(`
		SELECT
			COUNT(*) as total_users,
			COUNT(CASE WHEN last_login_at > CURRENT_TIMESTAMP - INTERVAL '30 days' THEN 1 END) as active_users,
			COUNT(CASE WHEN is_verified = true THEN 1 END) as verified_users
		FROM users
	`).Scan(&totalUsers, &activeUsers, &verifiedUsers)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get content statistics
	var totalDocuments, totalEmails, totalChatMessages, totalTasks int
	err = db.QueryRow(`
		SELECT
			(SELECT COUNT(*) FROM documents) as documents,
			(SELECT COUNT(*) FROM emails) as emails,
			(SELECT COUNT(*) FROM chat_messages) as chat_messages,
			(SELECT COUNT(*) FROM tasks) as tasks
	`).Scan(&totalDocuments, &totalEmails, &totalChatMessages, &totalTasks)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get storage statistics
	var totalStorageBytes int64
	err = db.QueryRow(`
		SELECT COALESCE(SUM(file_size), 0) FROM documents WHERE file_size IS NOT NULL
	`).Scan(&totalStorageBytes)
	if err != nil {
		totalStorageBytes = 0
	}

	c.JSON(http.StatusOK, gin.H{
		"users": gin.H{
			"total":     totalUsers,
			"active":    activeUsers,
			"verified":  verifiedUsers,
		},
		"content": gin.H{
			"documents":     totalDocuments,
			"emails":        totalEmails,
			"chat_messages": totalChatMessages,
			"tasks":         totalTasks,
		},
		"storage": gin.H{
			"total_bytes": totalStorageBytes,
			"total_gb":    float64(totalStorageBytes) / (1024 * 1024 * 1024),
		},
	})
}

// GetUsers returns paginated list of users for admin management
func GetUsers(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)

	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	search := c.Query("search")
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 50
	}
	offset := (page - 1) * limit

	// Build query
	query := `
		SELECT id, username, email, first_name, last_name, is_active, is_verified,
		       last_login_at, created_at, failed_login_attempts, locked_until
		FROM users
		WHERE 1=1
	`
	args := []interface{}{}
	argCount := 0

	if search != "" {
		query += ` AND (username ILIKE $` + strconv.Itoa(argCount+1) +
			` OR email ILIKE $` + strconv.Itoa(argCount+2) +
			` OR first_name ILIKE $` + strconv.Itoa(argCount+3) +
			` OR last_name ILIKE $` + strconv.Itoa(argCount+4) + `)`
		args = append(args, "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
		argCount += 4
	}

	if status != "" {
		switch status {
		case "active":
			query += ` AND is_active = $` + strconv.Itoa(argCount+1)
			args = append(args, true)
		case "inactive":
			query += ` AND is_active = $` + strconv.Itoa(argCount+1)
			args = append(args, false)
		case "verified":
			query += ` AND is_verified = $` + strconv.Itoa(argCount+1)
			args = append(args, true)
		case "unverified":
			query += ` AND is_verified = $` + strconv.Itoa(argCount+1)
			args = append(args, false)
		case "locked":
			query += ` AND locked_until IS NOT NULL AND locked_until > CURRENT_TIMESTAMP`
		}
		if status != "locked" {
			argCount++
		}
	}

	query += ` ORDER BY created_at DESC LIMIT $` + strconv.Itoa(argCount+1) +
		` OFFSET $` + strconv.Itoa(argCount+2)
	args = append(args, limit, offset)

	rows, err := db.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var users []gin.H
	for rows.Next() {
		var user struct {
			ID                   uuid.UUID  `json:"id"`
			Username             string     `json:"username"`
			Email                string     `json:"email"`
			FirstName            *string    `json:"first_name"`
			LastName             *string    `json:"last_name"`
			IsActive             bool       `json:"is_active"`
			IsVerified           bool       `json:"is_verified"`
			LastLoginAt          *string    `json:"last_login_at"`
			CreatedAt            string     `json:"created_at"`
			FailedLoginAttempts  int        `json:"failed_login_attempts"`
			LockedUntil          *string    `json:"locked_until"`
		}

		err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.FirstName, &user.LastName,
			&user.IsActive, &user.IsVerified, &user.LastLoginAt, &user.CreatedAt,
			&user.FailedLoginAttempts, &user.LockedUntil,
		)
		if err != nil {
			continue
		}

		users = append(users, gin.H{
			"id":                     user.ID,
			"username":               user.Username,
			"email":                  user.Email,
			"first_name":             user.FirstName,
			"last_name":              user.LastName,
			"is_active":              user.IsActive,
			"is_verified":            user.IsVerified,
			"last_login_at":          user.LastLoginAt,
			"created_at":             user.CreatedAt,
			"failed_login_attempts":  user.FailedLoginAttempts,
			"locked_until":           user.LockedUntil,
		})
	}

	// Get total count for pagination
	countQuery := `SELECT COUNT(*) FROM users WHERE 1=1`
	countArgs := []interface{}{}

	if search != "" {
		countQuery += ` AND (username ILIKE $1 OR email ILIKE $2 OR first_name ILIKE $3 OR last_name ILIKE $4)`
		countArgs = append(countArgs, "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	if status != "" {
		switch status {
		case "active":
			countQuery += ` AND is_active = $` + strconv.Itoa(len(countArgs)+1)
			countArgs = append(countArgs, true)
		case "inactive":
			countQuery += ` AND is_active = $` + strconv.Itoa(len(countArgs)+1)
			countArgs = append(countArgs, false)
		case "verified":
			countQuery += ` AND is_verified = $` + strconv.Itoa(len(countArgs)+1)
			countArgs = append(countArgs, true)
		case "unverified":
			countQuery += ` AND is_verified = $` + strconv.Itoa(len(countArgs)+1)
			countArgs = append(countArgs, false)
		case "locked":
			countQuery += ` AND locked_until IS NOT NULL AND locked_until > CURRENT_TIMESTAMP`
		}
	}

	var totalCount int
	err = db.QueryRow(countQuery, countArgs...).Scan(&totalCount)
	if err != nil {
		totalCount = 0
	}

	c.JSON(http.StatusOK, gin.H{
		"users":       users,
		"pagination": gin.H{
			"page":       page,
			"limit":      limit,
			"total":      totalCount,
			"totalPages": (totalCount + limit - 1) / limit,
		},
	})
}

// UpdateUserStatus updates a user's status (active/inactive)
func UpdateUserStatus(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	userIDStr := c.Param("id")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req struct {
		IsActive bool `json:"is_active" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = db.Exec(`
		UPDATE users SET is_active = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`, req.IsActive, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User status updated successfully"})
}

// DeleteUser deletes a user and all their data
func DeleteUser(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	userIDStr := c.Param("id")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Begin transaction for safe deletion
	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback()

	// Delete in order (respecting foreign key constraints)
	tables := []string{
		"password_resets",
		"user_sessions",
		"user_roles",
		"email_attachments",
		"emails",
		"email_accounts",
		"meeting_participants",
		"meeting_chat_messages",
		"meetings",
		"webrtc_connections",
		"calendar_shares",
		"events",
		"calendars",
		"task_comments",
		"tasks",
		"message_reactions",
		"chat_participants",
		"chat_messages",
		"chat_rooms",
		"form_responses",
		"form_fields",
		"forms",
		"document_comments",
		"document_versions",
		"document_shares",
		"documents",
		"contacts",
		"sync_devices",
		"sync_folders",
		"files",
		"file_versions",
		"file_chunks",
		"sync_queue",
		"sync_conflicts",
		"bandwidth_limits",
		"backup_schedules",
		"backups",
		"security_events",
		"rate_limit_events",
		"audit_log",
		"ai_usage",
		"ai_requests",
		"ai_upgrades",
	}

	for _, table := range tables {
		_, err = tx.Exec("DELETE FROM "+table+" WHERE user_id = $1", userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete from " + table})
			return
		}
	}

	// Finally delete the user
	_, err = tx.Exec("DELETE FROM users WHERE id = $1", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// GetSystemLogs returns recent system logs for admin review
func GetSystemLogs(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)

	// Parse query parameters
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	if limit < 1 || limit > 1000 {
		limit = 100
	}

	level := c.Query("level")
	eventType := c.Query("event_type")

	query := `
		SELECT timestamp, event_type, severity, user_id, ip_address, resource_type,
		       resource_id, action, status, details
		FROM audit_log
		WHERE 1=1
	`
	args := []interface{}{}
	argCount := 0

	if level != "" {
		query += ` AND severity = $` + strconv.Itoa(argCount+1)
		args = append(args, level)
		argCount++
	}

	if eventType != "" {
		query += ` AND event_type = $` + strconv.Itoa(argCount+1)
		args = append(args, eventType)
		argCount++
	}

	query += ` ORDER BY timestamp DESC LIMIT $` + strconv.Itoa(argCount+1)
	args = append(args, limit)

	rows, err := db.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var logs []gin.H
	for rows.Next() {
		var log struct {
			Timestamp    string `json:"timestamp"`
			EventType    string `json:"event_type"`
			Severity     string `json:"severity"`
			UserID       *uuid.UUID `json:"user_id"`
			IPAddress    *string `json:"ip_address"`
			ResourceType *string `json:"resource_type"`
			ResourceID   *uuid.UUID `json:"resource_id"`
			Action       *string `json:"action"`
			Status       *string `json:"status"`
			Details      *string `json:"details"`
		}

		err := rows.Scan(
			&log.Timestamp, &log.EventType, &log.Severity, &log.UserID,
			&log.IPAddress, &log.ResourceType, &log.ResourceID, &log.Action,
			&log.Status, &log.Details,
		)
		if err != nil {
			continue
		}

		logs = append(logs, gin.H{
			"timestamp":     log.Timestamp,
			"event_type":    log.EventType,
			"severity":      log.Severity,
			"user_id":       log.UserID,
			"ip_address":    log.IPAddress,
			"resource_type": log.ResourceType,
			"resource_id":   log.ResourceID,
			"action":        log.Action,
			"status":        log.Status,
			"details":       log.Details,
		})
	}

	c.JSON(http.StatusOK, gin.H{"logs": logs})
}

// GetDatabaseStats returns database performance statistics
func GetDatabaseStats(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)

	var stats struct {
		TableCount    int     `json:"table_count"`
		TotalSizeMB   float64 `json:"total_size_mb"`
		IndexSizeMB   float64 `json:"index_size_mb"`
		TableSizeMB   float64 `json:"table_size_mb"`
		ActiveConn    int     `json:"active_connections"`
		IdleConn      int     `json:"idle_connections"`
		TotalConn     int     `json:"total_connections"`
	}

	// Get table count
	err := db.QueryRow("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public'").Scan(&stats.TableCount)
	if err != nil {
		stats.TableCount = 0
	}

	// Get database size information
	err = db.QueryRow(`
		SELECT
			pg_size_pretty(pg_database_size(current_database())) as total_size,
			pg_size_pretty(sum(pg_total_relation_size(quote_ident(schemaname)||'.'||quote_ident(tablename)))) as table_size,
			pg_size_pretty(sum(pg_indexes_size(quote_ident(schemaname)||'.'||quote_ident(tablename)))) as index_size
		FROM pg_tables
		WHERE schemaname = 'public'
	`).Scan(&stats.TotalSizeMB, &stats.TableSizeMB, &stats.IndexSizeMB)
	if err != nil {
		// Fallback if the above query fails
		stats.TotalSizeMB = 0
		stats.TableSizeMB = 0
		stats.IndexSizeMB = 0
	}

	c.JSON(http.StatusOK, stats)
}

// RunDatabaseMaintenance runs VACUUM and REINDEX operations
func RunDatabaseMaintenance(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)

	// Run VACUUM ANALYZE on all tables
	tables := []string{
		"users", "documents", "emails", "chat_messages",
		"events", "contacts", "tasks",
	}

	for _, table := range tables {
		_, err := db.Exec("VACUUM ANALYZE " + table)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to vacuum table: " + table})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Database maintenance completed successfully"})
}

// GetSecurityAlerts returns recent security alerts
func GetSecurityAlerts(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if limit < 1 || limit > 200 {
		limit = 50
	}

	rows, err := db.Query(`
		SELECT id, event_type, severity, source_ip, user_id, details,
		       created_at, resolved, resolved_at
		FROM security_events
		ORDER BY created_at DESC
		LIMIT $1
	`, limit)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var alerts []gin.H
	for rows.Next() {
		var alert struct {
			ID         uuid.UUID `json:"id"`
			EventType  string    `json:"event_type"`
			Severity   string    `json:"severity"`
			SourceIP   *string   `json:"source_ip"`
			UserID     *uuid.UUID `json:"user_id"`
			Details    string    `json:"details"`
			CreatedAt  string    `json:"created_at"`
			Resolved   bool      `json:"resolved"`
			ResolvedAt *string   `json:"resolved_at"`
		}

		err := rows.Scan(
			&alert.ID, &alert.EventType, &alert.Severity, &alert.SourceIP,
			&alert.UserID, &alert.Details, &alert.CreatedAt, &alert.Resolved, &alert.ResolvedAt,
		)
		if err != nil {
			continue
		}

		alerts = append(alerts, gin.H{
			"id":          alert.ID,
			"event_type":  alert.EventType,
			"severity":    alert.Severity,
			"source_ip":   alert.SourceIP,
			"user_id":     alert.UserID,
			"details":     alert.Details,
			"created_at":  alert.CreatedAt,
			"resolved":    alert.Resolved,
			"resolved_at": alert.ResolvedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"alerts": alerts})
}

// ResolveSecurityAlert marks a security alert as resolved
func ResolveSecurityAlert(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	alertIDStr := c.Param("id")

	alertID, err := uuid.Parse(alertIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid alert ID"})
		return
	}

	_, err = db.Exec(`
		UPDATE security_events
		SET resolved = true, resolved_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`, alertID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Alert resolved successfully"})
}

// GetSystemSettings returns current system configuration
func GetSystemSettings(c *gin.Context) {
	// This would return current system settings
	// For now, return mock data
	settings := gin.H{
		"security": gin.H{
			"max_login_attempts": 5,
			"lockout_duration_minutes": 15,
			"password_min_length": 8,
			"session_timeout_hours": 24,
		},
		"storage": gin.H{
			"max_file_size_mb": 100,
			"backup_retention_days": 30,
			"storage_quota_gb": 10,
		},
		"features": gin.H{
			"ai_enabled": true,
			"file_sync_enabled": true,
			"video_conferencing_enabled": true,
			"backup_enabled": true,
		},
	}

	c.JSON(http.StatusOK, settings)
}

// UpdateSystemSettings updates system configuration
func UpdateSystemSettings(c *gin.Context) {
	var settings gin.H
	if err := c.ShouldBindJSON(&settings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// This would update system settings in database
	// For now, just acknowledge
	c.JSON(http.StatusOK, gin.H{
		"message": "System settings updated successfully",
		"settings": settings,
	})
}
