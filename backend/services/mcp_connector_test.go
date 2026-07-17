package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"tpt-titan/backend/models"
)

// TestMCPClientToolsList verifies the JSON-RPC client can handshake and list tools
// against a minimal MCP-compatible server.
func TestMCPClientToolsList(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req jsonRPCRequest
		_ = json.NewDecoder(r.Body).Decode(&req)
		w.Header().Set("Content-Type", "application/json")
		switch req.Method {
		case "initialize":
			_ = json.NewEncoder(w).Encode(jsonRPCResponse{JSONRPC: "2.0", ID: req.ID, Result: json.RawMessage(`{"protocolVersion":"2024-11-05"}`)})
		case "tools/list":
			_ = json.NewEncoder(w).Encode(jsonRPCResponse{JSONRPC: "2.0", ID: req.ID, Result: json.RawMessage(`{"tools":[{"name":"create_order","description":"Create an order"}]}`)})
		case "tools/call":
			_ = json.NewEncoder(w).Encode(jsonRPCResponse{JSONRPC: "2.0", ID: req.ID, Result: json.RawMessage(`{"content":[{"type":"text","text":"ok"}]}`)})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	client := NewMCPClient(models.MCPServer{URL: srv.URL})
	if err := client.Initialize(); err != nil {
		t.Fatalf("initialize failed: %v", err)
	}
	tools, err := client.ListTools()
	if err != nil {
		t.Fatalf("list tools failed: %v", err)
	}
	if len(tools) != 1 || tools[0].Name != "create_order" {
		t.Fatalf("unexpected tools: %+v", tools)
	}

	res, err := client.CallTool("create_order", map[string]interface{}{"id": 1})
	if err != nil {
		t.Fatalf("call tool failed: %v", err)
	}
	if res == nil {
		t.Fatal("expected non-nil result")
	}
}

// TestMCPConnectorName verifies the registry naming convention.
func TestMCPConnectorName(t *testing.T) {
	got := mcpConnectorName("My ERP", "Create Order")
	want := "mcp.my_erp.create_order"
	if got != want {
		t.Fatalf("mcpConnectorName = %q, want %q", got, want)
	}
}
