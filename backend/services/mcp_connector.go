package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"tpt-titan/backend/config"
	"tpt-titan/backend/models"
	"tpt-titan/backend/utils"

	"github.com/google/uuid"
)

// MCPWorkflowConnector is a workflow connector backed by a single tool exposed
// by a configured MCP server. Its connector name is "mcp.<serverName>.<toolName>"
// so it slots into the existing connector registry unchanged; the builder UI can
// select it like any other connector and pass the tool's arguments via
// nodeConfig["parameters"].
type MCPWorkflowConnector struct {
	ServerName string
	ServerURL  string
	ToolName   string
	AuthType   string
	Token      string // decrypted bearer token, if any
}

// jsonRPCRequest is a minimal JSON-RPC 2.0 envelope.
type jsonRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

// jsonRPCResponse is a minimal JSON-RPC 2.0 envelope.
type jsonRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// mcpTool is the shape of an entry in an MCP tools/list result.
type mcpTool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

// mcpClient talks JSON-RPC 2.0 over HTTP to a single MCP server.
type mcpClient struct {
	url      string
	authType string
	token    string
	http     *http.Client
}

func newMCPClient(server models.MCPServer) *mcpClient {
	return &mcpClient{
		url:      strings.TrimRight(server.URL, "/"),
		authType: server.AuthType,
		token:    server.Token,
		http:     &http.Client{Timeout: 15 * time.Second},
	}
}

// NewMCPClient constructs an MCP client for the given server configuration.
func NewMCPClient(server models.MCPServer) *mcpClient {
	return newMCPClient(server)
}

func (c *mcpClient) call(method string, params interface{}) (*jsonRPCResponse, error) {
	reqBody, err := json.Marshal(jsonRPCRequest{JSONRPC: "2.0", ID: 1, Method: method, Params: params})
	if err != nil {
		return nil, fmt.Errorf("mcp: failed to encode request: %w", err)
	}

	httpReq, err := http.NewRequest(http.MethodPost, c.url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("mcp: failed to build request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json, text/event-stream")
	if c.authType == "bearer" && c.token != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("mcp: server unreachable: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("mcp: failed to read response: %w", err)
	}

	// Some MCP servers reply with Server-Sent Events; extract the first data: line.
	body := string(raw)
	if strings.Contains(body, "event:") || strings.HasPrefix(strings.TrimSpace(body), "data:") {
		for _, line := range strings.Split(body, "\n") {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "data:") {
				body = strings.TrimSpace(strings.TrimPrefix(line, "data:"))
				break
			}
		}
	}

	var parsed jsonRPCResponse
	if err := json.Unmarshal([]byte(body), &parsed); err != nil {
		return nil, fmt.Errorf("mcp: invalid JSON-RPC response: %w", err)
	}
	if parsed.Error != nil {
		return nil, fmt.Errorf("mcp: server error %d: %s", parsed.Error.Code, parsed.Error.Message)
	}
	return &parsed, nil
}

// Initialize performs the MCP handshake (best-effort; some servers skip it).
func (c *mcpClient) Initialize() error {
	_, err := c.call("initialize", map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"clientInfo":      map[string]interface{}{"name": "tpt-titan", "version": "1.0.0"},
		"capabilities":    map[string]interface{}{},
	})
	return err
}

// ListTools returns the tools exposed by the server.
func (c *mcpClient) ListTools() ([]mcpTool, error) {
	resp, err := c.call("tools/list", map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	var out struct {
		Tools []mcpTool `json:"tools"`
	}
	if err := json.Unmarshal(resp.Result, &out); err != nil {
		return nil, fmt.Errorf("mcp: failed to parse tools/list: %w", err)
	}
	return out.Tools, nil
}

// CallTool invokes a single tool with the supplied arguments.
func (c *mcpClient) CallTool(name string, args map[string]interface{}) (map[string]interface{}, error) {
	resp, err := c.call("tools/call", map[string]interface{}{
		"name":      name,
		"arguments": args,
	})
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return nil, fmt.Errorf("mcp: failed to parse tools/call result: %w", err)
	}
	return result, nil
}

// Execute satisfies the WorkflowConnector interface. The node config's
// "parameters" map is passed verbatim as the tool's arguments.
func (m *MCPWorkflowConnector) Execute(nodeConfig map[string]interface{}, inputData map[string]interface{}) (map[string]interface{}, error) {
	client := newMCPClient(models.MCPServer{
		URL:      m.ServerURL,
		AuthType: m.AuthType,
		Token:    m.Token,
	})

	args, ok := nodeConfig["parameters"].(map[string]interface{})
	if !ok {
		// Allow args to be supplied inline on the node config as well.
		args = map[string]interface{}{}
		for k, v := range nodeConfig {
			if k == "connector" || k == "action" {
				continue
			}
			args[k] = v
		}
	}

	return client.CallTool(m.ToolName, args)
}

// GetConfigSchema returns a schema describing the connector's parameters.
func (m *MCPWorkflowConnector) GetConfigSchema() map[string]interface{} {
	return map[string]interface{}{
		"type":        "object",
		"title":       fmt.Sprintf("MCP: %s", m.ToolName),
		"description": fmt.Sprintf("Calls the %q tool on MCP server %q", m.ToolName, m.ServerName),
		"properties": map[string]interface{}{
			"parameters": map[string]interface{}{
				"type":        "object",
				"description": "Arguments passed to the MCP tool",
			},
		},
	}
}

// RegisterMCPConnectors loads every active MCP server for the given user (or all
// users when userID is nil) and registers each of its tools as a workflow
// connector named "mcp.<serverName>.<toolName>". It is safe to call repeatedly —
// it only adds, never removes, stale registrations (the connector registry is
// process-global and MCP servers are low-churn).
func (s *WorkflowService) RegisterMCPConnectors(userID *uuid.UUID) error {
	if config.DB == nil {
		return nil
	}

	query := config.DB.Where("is_active = ?", true)
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	var servers []models.MCPServer
	if err := query.Find(&servers).Error; err != nil {
		return fmt.Errorf("mcp: failed to list servers: %w", err)
	}

	for _, server := range servers {
		token := server.Token
		if server.AuthType == "bearer" && token != "" {
			if dec, err := utils.DecryptPassword([]byte(token)); err == nil {
				token = dec
			}
		}

		client := newMCPClient(server)
		// Best-effort handshake; tolerate servers that don't require it.
		_ = client.Initialize()

		tools, err := client.ListTools()
		if err != nil {
			log.Printf("mcp: skipping server %q: %v", server.Name, err)
			continue
		}

		for _, tool := range tools {
			name := mcpConnectorName(server.Name, tool.Name)
			s.RegisterConnector(name, &MCPWorkflowConnector{
				ServerName: server.Name,
				ServerURL:  server.URL,
				ToolName:   tool.Name,
				AuthType:   server.AuthType,
				Token:      token,
			})
		}
	}
	return nil
}

// mcpConnectorName builds the registry key for an MCP tool connector.
func mcpConnectorName(serverName, toolName string) string {
	serverName = strings.ToLower(strings.ReplaceAll(serverName, " ", "_"))
	toolName = strings.ToLower(strings.ReplaceAll(toolName, " ", "_"))
	return fmt.Sprintf("mcp.%s.%s", serverName, toolName)
}
