package models

import (
	"time"

	"gorm.io/gorm"
)

// WebhookDeliveryLog records inbound and outbound webhook/API calls for
// observability and debugging of misconfigured workflows, similar to a
// webhook dashboard. Inbound rows are created when the public webhook receiver
// is hit; outbound rows are created by the http.request connector (one per
// attempt, including retries).
type WebhookDeliveryLog struct {
	ID           string         `json:"id" gorm:"type:varchar(36);primaryKey"`
	Direction    string         `json:"direction" gorm:"type:varchar(16);index"` // "inbound" | "outbound"
	Connector    string         `json:"connector" gorm:"type:varchar(64);index"`
	WorkflowID   string         `json:"workflow_id" gorm:"type:varchar(36);index"`
	NodeID       string         `json:"node_id" gorm:"type:varchar(64);index"`
	// For inbound: the matched webhook token (kept for admin lookup, not exposed).
	Token       string `json:"-" gorm:"type:varchar(128);index"`
	Method      string `json:"method" gorm:"type:varchar(16)"`
	URL         string `json:"url" gorm:"type:text"`
	Host        string `json:"host" gorm:"type:varchar(255);index"`
	// RequestBody/ResponseBody are truncated to keep the table bounded.
	RequestBody  string `json:"request_body" gorm:"type:text"`
	StatusCode   int    `json:"status_code" gorm:"index"`
	ResponseBody string `json:"response_body" gorm:"type:text"`
	// Error captures transport/validation failures when no HTTP response exists.
	Error     string         `json:"error" gorm:"type:text"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

const webhookDeliveryLogBodyLimit = 4096
