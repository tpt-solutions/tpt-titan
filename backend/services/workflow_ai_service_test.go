// backend/services/workflow_ai_service_test.go
// Regression tests for the from/to connection-key fix in workflow_ai_service.go.
// Run with: cd backend && go test ./services/... -run TestFindLongSequentialChains -v

package services

import "testing"

func TestFindLongSequentialChains_WalksFromToEdges(t *testing.T) {
	s := &WorkflowAIService{}

	nodes := []interface{}{
		map[string]interface{}{"id": "a"},
		map[string]interface{}{"id": "b"},
		map[string]interface{}{"id": "c"},
	}
	connections := []interface{}{
		map[string]interface{}{"from": "a", "to": "b"},
		map[string]interface{}{"from": "b", "to": "c"},
	}

	chains := s.findLongSequentialChains(nodes, connections)
	if len(chains) != 1 {
		t.Fatalf("expected 1 chain, got %d — the from/to edges were not walked", len(chains))
	}
	if len(chains[0]) != 3 {
		t.Fatalf("expected chain of length 3 (a->b->c), got %d", len(chains[0]))
	}
}

func TestFindLongSequentialChains_IgnoresSourceTargetKeys(t *testing.T) {
	s := &WorkflowAIService{}

	nodes := []interface{}{
		map[string]interface{}{"id": "a"},
		map[string]interface{}{"id": "b"},
	}
	// Old (wrong) key names — must not be treated as an edge.
	connections := []interface{}{
		map[string]interface{}{"source": "a", "target": "b"},
	}

	chains := s.findLongSequentialChains(nodes, connections)
	if len(chains) != 0 {
		t.Fatalf("expected no chains from source/target-keyed connections, got %d", len(chains))
	}
}

func TestFindUnusedDataFlow_UsesFromToEdges(t *testing.T) {
	s := &WorkflowAIService{}

	nodes := []interface{}{
		map[string]interface{}{"id": "connected-a"},
		map[string]interface{}{"id": "connected-b"},
		map[string]interface{}{"id": "isolated"},
	}
	connections := []interface{}{
		map[string]interface{}{"from": "connected-a", "to": "connected-b"},
	}

	unused := s.findUnusedDataFlow(nodes, connections)
	if len(unused) != 1 {
		t.Fatalf("expected exactly 1 unused node, got %d", len(unused))
	}
	unusedNode, ok := unused[0].(map[string]interface{})
	if !ok || unusedNode["id"] != "isolated" {
		t.Fatalf("expected the isolated node to be reported as unused, got %v", unused[0])
	}
}
