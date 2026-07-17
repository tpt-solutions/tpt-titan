// backend/services/workflow_service_test.go
// Run with: cd backend && go test ./services/... -run TestFindNextNodes -v

package services

import "testing"

func nodeMapOf(ids ...string) map[string]map[string]interface{} {
	m := make(map[string]map[string]interface{})
	for _, id := range ids {
		m[id] = map[string]interface{}{"id": id}
	}
	return m
}

// ─── findNextNodes: unconditional fan-out for non-condition sources ──────────

func TestFindNextNodes_UnconditionalFanoutIgnoresPorts(t *testing.T) {
	s := &WorkflowService{}
	nodeMap := nodeMapOf("b", "c")
	connections := []interface{}{
		map[string]interface{}{"from": "a", "to": "b", "fromPort": "true"},
		map[string]interface{}{"from": "a", "to": "c", "fromPort": "false"},
	}

	// branchPort == "" means "not a condition source" — both edges should fire.
	next := s.findNextNodes("a", "", connections, nodeMap)
	if len(next) != 2 {
		t.Fatalf("expected 2 next nodes for unconditional fan-out, got %d", len(next))
	}
}

// ─── findNextNodes: condition branching gates downstream execution ──────────

func TestFindNextNodes_ConditionBranching_TrueBranchOnly(t *testing.T) {
	s := &WorkflowService{}
	nodeMap := nodeMapOf("on-true", "on-false")
	connections := []interface{}{
		map[string]interface{}{"from": "cond", "to": "on-true", "fromPort": "true"},
		map[string]interface{}{"from": "cond", "to": "on-false", "fromPort": "false"},
	}

	next := s.findNextNodes("cond", "true", connections, nodeMap)
	if len(next) != 1 || next[0]["id"] != "on-true" {
		t.Fatalf("expected only the true-branch node, got %v", next)
	}
}

func TestFindNextNodes_ConditionBranching_FalseBranchOnly(t *testing.T) {
	s := &WorkflowService{}
	nodeMap := nodeMapOf("on-true", "on-false")
	connections := []interface{}{
		map[string]interface{}{"from": "cond", "to": "on-true", "fromPort": "true"},
		map[string]interface{}{"from": "cond", "to": "on-false", "fromPort": "false"},
	}

	next := s.findNextNodes("cond", "false", connections, nodeMap)
	if len(next) != 1 || next[0]["id"] != "on-false" {
		t.Fatalf("expected only the false-branch node, got %v", next)
	}
}

// ─── executeCondition ────────────────────────────────────────────────────────

func TestExecuteCondition_Equals(t *testing.T) {
	s := &WorkflowService{}
	config := map[string]interface{}{"field": "priority", "operator": "equals", "value": "high"}

	result, err := s.executeCondition(config, map[string]interface{}{"priority": "high"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["condition_result"] != true {
		t.Fatalf("expected condition_result true, got %v", result["condition_result"])
	}
}

func TestExecuteCondition_NotEquals(t *testing.T) {
	s := &WorkflowService{}
	config := map[string]interface{}{"field": "priority", "operator": "equals", "value": "high"}

	result, err := s.executeCondition(config, map[string]interface{}{"priority": "low"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["condition_result"] != false {
		t.Fatalf("expected condition_result false, got %v", result["condition_result"])
	}
}

// ─── executeNode: dry-run mode never invokes the real connector ─────────────

func TestExecuteNode_DryRun_ActionNode_SkipsRealConnector(t *testing.T) {
	s := NewWorkflowService()
	ctx := &WorkflowExecutionContext{
		NodeStates: make(map[string]interface{}),
		GlobalData: map[string]interface{}{"user_id": "not-a-real-user"},
		DryRun:     true,
	}
	node := map[string]interface{}{
		"id":   "action-1",
		"type": "action",
		"config": map[string]interface{}{
			"connector": "tasks.create",
			"title":     "Should never be persisted",
		},
	}

	// If dry-run didn't skip the real connector, TaskCreateConnector.Execute would
	// attempt a real config.DB.Create call — config.DB is nil in this unit test,
	// so a real invocation would panic. Reaching the assertions below proves the
	// dry-run path never called the connector.
	err := s.executeNode(node, ctx, nodeMapOf("action-1"), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ctx.GlobalData["dry_run"] != true {
		t.Fatalf("expected dry_run preview data merged into GlobalData, got %v", ctx.GlobalData)
	}
	if ctx.GlobalData["would_execute"] != "tasks.create" {
		t.Fatalf("expected would_execute to name the connector, got %v", ctx.GlobalData["would_execute"])
	}
}

func TestExecuteNode_UnknownConnector_ReturnsError(t *testing.T) {
	s := NewWorkflowService()
	ctx := &WorkflowExecutionContext{
		NodeStates: make(map[string]interface{}),
		GlobalData: map[string]interface{}{},
		DryRun:     true,
	}
	node := map[string]interface{}{
		"id":   "action-1",
		"type": "action",
		"config": map[string]interface{}{
			"connector": "not.a.real.connector",
		},
	}

	err := s.executeNode(node, ctx, nodeMapOf("action-1"), nil)
	if err == nil {
		t.Fatal("expected an error for an unregistered connector")
	}
}
