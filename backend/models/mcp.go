package models

import (
	"time"

	"github.com/google/uuid"
)

// MCPServer represents a configured Model Context Protocol server whose tools
// are exposed to the workflow engine as connectors named "mcp.<server>.<tool>".
//
// MCP (https://modelcontextprotocol.io) is the bridge layer that lets this app
// talk to external systems — including the sibling tpt-free-erp project — at the
// API boundary without merging codebases. A server is reached over JSON-RPC 2.0
// over HTTP (the `transport` field selects "http" or "streamable-http", both
// POST a JSON-RPC body to `URL`). Auth is optional; when `AuthType` is "bearer"
// the `Token` is encrypted at rest (AES-256-GCM) and never serialized out.
type MCPServer struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Name      string    `gorm:"size:255;not null" json:"name"`
	URL       string    `gorm:"size:1024;not null" json:"url"`
	Transport string    `gorm:"size:32;default:'http'" json:"transport"` // "http" | "streamable-http"
	AuthType  string    `gorm:"size:32;default:'none'" json:"auth_type"` // "none" | "bearer"
	Token     string    `gorm:"type:text" json:"-"`                       // encrypted; never serialized out
	IsActive  bool      `gorm:"default:true" json:"is_active"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName scopes MCP servers to their own table.
func (MCPServer) TableName() string { return "mcp_servers" }
