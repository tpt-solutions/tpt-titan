package services

import (
	"log"

	"github.com/google/uuid"
	"tpt-titan/backend/config"
	"tpt-titan/backend/models"
)

// truncateBytes caps a string to max bytes, appending an ellipsis note when it
// was shortened, so a large payload can't blow up the delivery-log table.
func truncateBytes(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "…[truncated]"
}

// RecordDeliveryLog persists one inbound/outbound webhook call row. Failures are
// logged but never propagated — the delivery log is best-effort observability,
// not part of the request's success path.
func RecordDeliveryLog(entry *models.WebhookDeliveryLog) {
	if entry.ID == "" {
		entry.ID = uuid.NewString()
	}
	db := config.GetDatabase()
	if db == nil {
		return
	}
	if err := db.Create(entry).Error; err != nil {
		log.Printf("RecordDeliveryLog: failed to persist entry: %v", err)
	}
}

// GetWebhookDeliveryLogs returns recent delivery-log rows for the monitoring /
// admin webhook dashboard, newest first.
func GetWebhookDeliveryLogs(limit int) ([]models.WebhookDeliveryLog, error) {
	if limit <= 0 || limit > 500 {
		limit = 200
	}
	db := config.GetDatabase()
	if db == nil {
		return nil, nil
	}
	var logs []models.WebhookDeliveryLog
	err := db.Order("created_at DESC").Limit(limit).Find(&logs).Error
	return logs, err
}
