package services

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"tpt-titan/backend/config"
	"tpt-titan/backend/models"

	"github.com/google/uuid"
)

// WorkflowAIService provides AI-powered workflow suggestions and analysis
type WorkflowAIService struct {
	aiService *AIService
}

// NewWorkflowAIService creates a new AI workflow service
func NewWorkflowAIService(aiService *AIService) *WorkflowAIService {
	return &WorkflowAIService{
		aiService: aiService,
	}
}

// AnalyzeUsagePatterns analyzes user workflow usage patterns
func (s *WorkflowAIService) AnalyzeUsagePatterns(userID uuid.UUID) (*WorkflowUsageAnalysis, error) {
	analysis := &WorkflowUsageAnalysis{
		UserID:         userID,
		TimeRange:      "30d",
		Patterns:       []UsagePattern{},
		Recommendations: []WorkflowRecommendation{},
	}

	// Get workflow execution data for the last 30 days
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	var executions []models.WorkflowExecution

	err := config.DB.Where("user_id = ? AND started_at >= ?", userID, thirtyDaysAgo).
		Order("started_at DESC").Find(&executions).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve execution data: %w", err)
	}

	if len(executions) == 0 {
		return analysis, nil
	}

	// Analyze execution patterns
	patterns := s.analyzeExecutionPatterns(executions)
	analysis.Patterns = patterns

	// Analyze workflow performance
	performance := s.analyzeWorkflowPerformance(executions)
	analysis.PerformanceMetrics = performance

	// Generate recommendations based on patterns
	recommendations := s.generateRecommendations(patterns, performance, userID)
	analysis.Recommendations = recommendations

	// Analyze peak usage times
	peakTimes := s.analyzePeakUsageTimes(executions)
	analysis.PeakUsageTimes = peakTimes

	return analysis, nil
}

// GetSmartTemplateRecommendations provides personalized template suggestions
func (s *WorkflowAIService) GetSmartTemplateRecommendations(userID uuid.UUID, context map[string]interface{}) ([]TemplateRecommendation, error) {
	var recommendations []TemplateRecommendation

	// Get user's existing workflows and usage patterns
	patterns, err := s.AnalyzeUsagePatterns(userID)
	if err != nil {
		return recommendations, err
	}

	// Get available templates
	var templates []models.WorkflowTemplate
	err = config.DB.Where("is_public = ?", true).Find(&templates).Error
	if err != nil {
		return recommendations, fmt.Errorf("failed to retrieve templates: %w", err)
	}

	// Score templates based on user patterns and context
	for _, template := range templates {
		score := s.scoreTemplateForUser(template, patterns, context)
		if score > 0.3 { // Only recommend templates with reasonable relevance
			recommendations = append(recommendations, TemplateRecommendation{
				TemplateID:   template.ID,
				TemplateName: template.Name,
				Category:     template.Category,
				RelevanceScore: score,
				Reason:       s.generateRecommendationReason(template, patterns, context),
			})
		}
	}

	// Sort by relevance score
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].RelevanceScore > recommendations[j].RelevanceScore
	})

	// Return top 5 recommendations
	if len(recommendations) > 5 {
		recommendations = recommendations[:5]
	}

	return recommendations, nil
}

// OptimizeWorkflow suggests optimizations for an existing workflow
func (s *WorkflowAIService) OptimizeWorkflow(workflowID uuid.UUID, userID uuid.UUID) (*WorkflowOptimization, error) {
	var workflow models.Workflow
	err := config.DB.Where("id = ? AND user_id = ?", workflowID, userID).First(&workflow).Error
	if err != nil {
		return nil, fmt.Errorf("workflow not found: %w", err)
	}

	optimization := &WorkflowOptimization{
		WorkflowID:   workflowID,
		Suggestions:  []OptimizationSuggestion{},
		EstimatedImprovement: WorkflowImprovement{},
	}

	// Parse workflow canvas data
	var canvasData map[string]interface{}
	if err := json.Unmarshal([]byte(workflow.CanvasData), &canvasData); err != nil {
		return nil, fmt.Errorf("failed to parse workflow data: %w", err)
	}

	nodes := canvasData["nodes"].([]interface{})
	connections := canvasData["connections"].([]interface{})

	// Analyze workflow structure
	suggestions := s.analyzeWorkflowStructure(nodes, connections, workflowID)
	optimization.Suggestions = suggestions

	// Estimate performance improvements
	improvement := s.estimatePerformanceImprovement(suggestions)
	optimization.EstimatedImprovement = improvement

	return optimization, nil
}

// PredictWorkflowCreation suggests new workflows based on user behavior patterns
func (s *WorkflowAIService) PredictWorkflowCreation(userID uuid.UUID) ([]PredictedWorkflow, error) {
	var predictions []PredictedWorkflow

	// Analyze user activity across different modules
	userActivity := s.analyzeUserActivity(userID)

	// Identify potential workflow opportunities
	opportunities := s.identifyWorkflowOpportunities(userActivity)

	// Generate predicted workflows
	for _, opportunity := range opportunities {
		if opportunity.Confidence > 0.6 { // Only high-confidence predictions
			prediction := s.generatePredictedWorkflow(opportunity)
			if prediction != nil {
				predictions = append(predictions, *prediction)
			}
		}
	}

	// Sort by confidence
	sort.Slice(predictions, func(i, j int) bool {
		return predictions[i].Confidence > predictions[j].Confidence
	})

	return predictions, nil
}

// analyzeExecutionPatterns analyzes workflow execution patterns
func (s *WorkflowAIService) analyzeExecutionPatterns(executions []models.WorkflowExecution) []UsagePattern {
	patterns := []UsagePattern{}

	// Group executions by workflow
	workflowCounts := make(map[uuid.UUID]int)
	workflowSuccessRates := make(map[uuid.UUID]float64)
	workflowAvgDuration := make(map[uuid.UUID]float64)

	for _, exec := range executions {
		workflowCounts[exec.WorkflowID]++

		// Calculate success rate
		if exec.Status == "completed" {
			currentRate := workflowSuccessRates[exec.WorkflowID]
			total := workflowCounts[exec.WorkflowID]
			workflowSuccessRates[exec.WorkflowID] = (currentRate*float64(total-1) + 1) / float64(total)
		}

		// Calculate average duration
		if exec.Duration > 0 {
			currentAvg := workflowAvgDuration[exec.WorkflowID]
			total := workflowCounts[exec.WorkflowID]
			workflowAvgDuration[exec.WorkflowID] = (currentAvg*float64(total-1) + float64(exec.Duration)) / float64(total)
		}
	}

	// Convert to patterns
	for workflowID, count := range workflowCounts {
		var workflow models.Workflow
		config.DB.Where("id = ?", workflowID).First(&workflow)

		patterns = append(patterns, UsagePattern{
			WorkflowID:      workflowID,
			WorkflowName:    workflow.Name,
			ExecutionCount:  count,
			SuccessRate:     workflowSuccessRates[workflowID],
			AvgDuration:     workflowAvgDuration[workflowID],
			LastExecuted:    workflow.LastRunAt,
			Frequency:       s.calculateExecutionFrequency(executions, workflowID),
		})
	}

	return patterns
}

// analyzeWorkflowPerformance analyzes overall workflow performance
func (s *WorkflowAIService) analyzeWorkflowPerformance(executions []models.WorkflowExecution) WorkflowPerformance {
	performance := WorkflowPerformance{
		TotalExecutions:    len(executions),
		SuccessfulExecutions: 0,
		FailedExecutions:   0,
		AverageDuration:    0,
		PeakUsageHour:      0,
		MostUsedWorkflow:   uuid.Nil,
	}

	totalDuration := int64(0)
	hourlyUsage := make(map[int]int)
	maxUsage := 0

	for _, exec := range executions {
		if exec.Status == "completed" {
			performance.SuccessfulExecutions++
		} else if exec.Status == "failed" {
			performance.FailedExecutions++
		}

		totalDuration += int64(exec.Duration)

		// Track hourly usage
		hour := exec.StartedAt.Hour()
		hourlyUsage[hour]++
		if hourlyUsage[hour] > maxUsage {
			maxUsage = hourlyUsage[hour]
			performance.PeakUsageHour = hour
		}
	}

	if len(executions) > 0 {
		performance.AverageDuration = float64(totalDuration) / float64(len(executions))
	}

	// Find most used workflow
	workflowUsage := make(map[uuid.UUID]int)
	for _, exec := range executions {
		workflowUsage[exec.WorkflowID]++
	}

	maxCount := 0
	for workflowID, count := range workflowUsage {
		if count > maxCount {
			maxCount = count
			performance.MostUsedWorkflow = workflowID
		}
	}

	return performance
}

// generateRecommendations generates workflow recommendations based on analysis
func (s *WorkflowAIService) generateRecommendations(patterns []UsagePattern, performance WorkflowPerformance, userID uuid.UUID) []WorkflowRecommendation {
	recommendations := []WorkflowRecommendation{}

	// Check for workflows with low success rates
	for _, pattern := range patterns {
		if pattern.SuccessRate < 0.8 && pattern.ExecutionCount > 3 {
			recommendations = append(recommendations, WorkflowRecommendation{
				Type:        "optimization",
				WorkflowID:  pattern.WorkflowID,
				Title:       "Improve Workflow Reliability",
				Description: fmt.Sprintf("'%s' has a %.1f%% success rate. Consider reviewing error patterns and adding error handling.", pattern.WorkflowName, pattern.SuccessRate*100),
				Priority:    "high",
				Confidence:  0.85,
			})
		}
	}

	// Check for workflows with long execution times
	for _, pattern := range patterns {
		if pattern.AvgDuration > 30000 && pattern.ExecutionCount > 2 { // > 30 seconds
			recommendations = append(recommendations, WorkflowRecommendation{
				Type:        "optimization",
				WorkflowID:  pattern.WorkflowID,
				Title:       "Optimize Workflow Performance",
				Description: fmt.Sprintf("'%s' takes an average of %.1f seconds to execute. Consider optimizing slow steps or adding parallel processing.", pattern.WorkflowName, pattern.AvgDuration/1000),
				Priority:    "medium",
				Confidence:  0.75,
			})
		}
	}

	// Suggest workflow consolidation if user has many similar workflows
	if len(patterns) > 5 {
		recommendations = append(recommendations, WorkflowRecommendation{
			Type:        "consolidation",
			Title:       "Consolidate Similar Workflows",
			Description: "You have multiple workflows that might be performing similar tasks. Consider combining them for better maintainability.",
			Priority:    "low",
			Confidence:  0.6,
		})
	}

	// Suggest scheduling for frequently manual workflows
	for _, pattern := range patterns {
		if pattern.Frequency == "daily" && pattern.ExecutionCount > 7 {
			recommendations = append(recommendations, WorkflowRecommendation{
				Type:        "automation",
				WorkflowID:  pattern.WorkflowID,
				Title:       "Consider Scheduling",
				Description: fmt.Sprintf("'%s' is executed %s. Consider setting up automatic scheduling to reduce manual effort.", pattern.WorkflowName, pattern.Frequency),
				Priority:    "medium",
				Confidence:  0.7,
			})
		}
	}

	return recommendations
}

// analyzePeakUsageTimes analyzes when workflows are most frequently executed
func (s *WorkflowAIService) analyzePeakUsageTimes(executions []models.WorkflowExecution) []PeakUsageTime {
	peakTimes := []PeakUsageTime{}

	hourlyUsage := make(map[int]int)
	weekdayUsage := make(map[string]int)

	for _, exec := range executions {
		hour := exec.StartedAt.Hour()
		hourlyUsage[hour]++

		weekday := exec.StartedAt.Weekday().String()
		weekdayUsage[weekday]++
	}

	// Find peak hours
	maxHourlyUsage := 0
	peakHour := 0
	for hour, count := range hourlyUsage {
		if count > maxHourlyUsage {
			maxHourlyUsage = count
			peakHour = hour
		}
	}

	if maxHourlyUsage > 0 {
		peakTimes = append(peakTimes, PeakUsageTime{
			TimeSlot: fmt.Sprintf("%02d:00-%02d:59", peakHour, peakHour),
			UsageCount: maxHourlyUsage,
			Type: "hourly",
		})
	}

	// Find peak weekdays
	maxWeekdayUsage := 0
	peakWeekday := ""
	for weekday, count := range weekdayUsage {
		if count > maxWeekdayUsage {
			maxWeekdayUsage = count
			peakWeekday = weekday
		}
	}

	if maxWeekdayUsage > 0 {
		peakTimes = append(peakTimes, PeakUsageTime{
			TimeSlot: peakWeekday,
			UsageCount: maxWeekdayUsage,
			Type: "weekday",
		})
	}

	return peakTimes
}

// scoreTemplateForUser scores how relevant a template is for a user
func (s *WorkflowAIService) scoreTemplateForUser(template models.WorkflowTemplate, patterns *WorkflowUsageAnalysis, context map[string]interface{}) float64 {
	score := 0.0

	// Base score from category matching
	userCategories := make(map[string]bool)
	for _, pattern := range patterns.Patterns {
		// Extract category from workflow name or description (simplified)
		category := s.inferWorkflowCategory(pattern.WorkflowName)
		userCategories[category] = true
	}

	if userCategories[template.Category] {
		score += 0.4
	}

	// Score based on usage patterns
	for _, pattern := range patterns.Patterns {
		if pattern.ExecutionCount > 5 {
			score += 0.2 // User is active with workflows
		}
		if pattern.SuccessRate > 0.9 {
			score += 0.1 // User creates successful workflows
		}
	}

	// Context-based scoring
	if context != nil {
		if context["recent_form_submissions"] != nil {
			if template.Category == "business" {
				score += 0.2
			}
		}
		if context["email_volume"] != nil {
			if template.Category == "communication" {
				score += 0.2
			}
		}
	}

	// Normalize score
	return math.Min(score, 1.0)
}

// generateRecommendationReason generates a human-readable reason for template recommendation
func (s *WorkflowAIService) generateRecommendationReason(template models.WorkflowTemplate, patterns *WorkflowUsageAnalysis, context map[string]interface{}) string {
	reasons := []string{}

	if len(patterns.Patterns) > 0 {
		reasons = append(reasons, "Based on your workflow usage patterns")
	}

	if template.Category == "business" && context != nil && context["recent_form_submissions"] != nil {
		reasons = append(reasons, "You frequently work with forms and data collection")
	}

	if template.Category == "communication" && context != nil && context["email_volume"] != nil {
		reasons = append(reasons, "You handle significant email communication")
	}

	if len(reasons) == 0 {
		reasons = append(reasons, "This template matches common workflow automation needs")
	}

	return strings.Join(reasons, ". ")
}

// analyzeWorkflowStructure analyzes workflow structure for optimization opportunities
func (s *WorkflowAIService) analyzeWorkflowStructure(nodes []interface{}, connections []interface{}, workflowID uuid.UUID) []OptimizationSuggestion {
	suggestions := []OptimizationSuggestion{}

	// Check for sequential bottlenecks
	nodeMap := make(map[string]interface{})
	for _, node := range nodes {
		nodeData := node.(map[string]interface{})
		nodeID := nodeData["id"].(string)
		nodeMap[nodeID] = nodeData
	}

	// Find long sequential chains
	longChains := s.findLongSequentialChains(nodes, connections)
	for _, chain := range longChains {
		if len(chain) > 5 {
			suggestions = append(suggestions, OptimizationSuggestion{
				Type: "parallelization",
				Title: "Consider Parallel Processing",
				Description: fmt.Sprintf("Workflow has a long sequential chain of %d steps. Consider running independent steps in parallel.", len(chain)),
				Impact: "high",
				Effort: "medium",
			})
		}
	}

	// Check for redundant operations
	redundantOps := s.findRedundantOperations(nodes)
	for _, op := range redundantOps {
		suggestions = append(suggestions, OptimizationSuggestion{
			Type: "consolidation",
			Title: "Remove Redundant Operations",
			Description: fmt.Sprintf("Found %d similar %s operations that could be consolidated.", op.Count, op.OperationType),
			Impact: "medium",
			Effort: "low",
		})
	}

	// Check for error handling
	errorHandling := s.analyzeErrorHandling(nodes, connections)
	if !errorHandling.HasErrorHandlers {
		suggestions = append(suggestions, OptimizationSuggestion{
			Type: "reliability",
			Title: "Add Error Handling",
			Description: "Workflow lacks error handling for failed operations. Consider adding error paths and retry logic.",
			Impact: "high",
			Effort: "medium",
		})
	}

	// Check for unused variables/data
	unusedData := s.findUnusedDataFlow(nodes, connections)
	if len(unusedData) > 0 {
		suggestions = append(suggestions, OptimizationSuggestion{
			Type: "cleanup",
			Title: "Remove Unused Data Flows",
			Description: fmt.Sprintf("Found %d unused data connections that can be removed to simplify the workflow.", len(unusedData)),
			Impact: "low",
			Effort: "low",
		})
	}

	return suggestions
}

// Helper functions for analysis
func (s *WorkflowAIService) calculateExecutionFrequency(executions []models.WorkflowExecution, workflowID uuid.UUID) string {
	workflowExecs := []models.WorkflowExecution{}
	for _, exec := range executions {
		if exec.WorkflowID == workflowID {
			workflowExecs = append(workflowExecs, exec)
		}
	}

	if len(workflowExecs) < 2 {
		return "occasional"
	}

	// Calculate average interval
	sort.Slice(workflowExecs, func(i, j int) bool {
		return workflowExecs[i].StartedAt.Before(workflowExecs[j].StartedAt)
	})

	totalInterval := time.Duration(0)
	for i := 1; i < len(workflowExecs); i++ {
		interval := workflowExecs[i].StartedAt.Sub(workflowExecs[i-1].StartedAt)
		totalInterval += interval
	}

	avgInterval := totalInterval / time.Duration(len(workflowExecs)-1)
	days := avgInterval.Hours() / 24

	switch {
	case days < 1:
		return "multiple times daily"
	case days < 7:
		return "daily"
	case days < 30:
		return "weekly"
	default:
		return "monthly"
	}
}

func (s *WorkflowAIService) inferWorkflowCategory(workflowName string) string {
	name := strings.ToLower(workflowName)

	if strings.Contains(name, "invoice") || strings.Contains(name, "payment") || strings.Contains(name, "expense") {
		return "business"
	}
	if strings.Contains(name, "email") || strings.Contains(name, "notification") {
		return "communication"
	}
	if strings.Contains(name, "form") || strings.Contains(name, "survey") {
		return "data"
	}
	if strings.Contains(name, "calendar") || strings.Contains(name, "meeting") {
		return "scheduling"
	}

	return "general"
}

// Data structures for AI analysis results
type WorkflowUsageAnalysis struct {
	UserID            uuid.UUID               `json:"user_id"`
	TimeRange         string                  `json:"time_range"`
	Patterns          []UsagePattern          `json:"patterns"`
	PerformanceMetrics WorkflowPerformance    `json:"performance_metrics"`
	Recommendations   []WorkflowRecommendation `json:"recommendations"`
	PeakUsageTimes    []PeakUsageTime         `json:"peak_usage_times"`
}

type UsagePattern struct {
	WorkflowID     uuid.UUID `json:"workflow_id"`
	WorkflowName   string    `json:"workflow_name"`
	ExecutionCount int       `json:"execution_count"`
	SuccessRate    float64   `json:"success_rate"`
	AvgDuration    float64   `json:"avg_duration_ms"`
	LastExecuted   *time.Time `json:"last_executed"`
	Frequency      string    `json:"frequency"`
}

type WorkflowPerformance struct {
	TotalExecutions      int       `json:"total_executions"`
	SuccessfulExecutions int       `json:"successful_executions"`
	FailedExecutions     int       `json:"failed_executions"`
	AverageDuration      float64   `json:"average_duration_ms"`
	PeakUsageHour        int       `json:"peak_usage_hour"`
	MostUsedWorkflow     uuid.UUID `json:"most_used_workflow"`
}

type WorkflowRecommendation struct {
	Type       string    `json:"type"`
	WorkflowID uuid.UUID `json:"workflow_id,omitempty"`
	Title      string    `json:"title"`
	Description string   `json:"description"`
	Priority   string    `json:"priority"`
	Confidence float64   `json:"confidence"`
}

type PeakUsageTime struct {
	TimeSlot   string `json:"time_slot"`
	UsageCount int    `json:"usage_count"`
	Type       string `json:"type"` // "hourly" or "weekday"
}

type TemplateRecommendation struct {
	TemplateID     uuid.UUID `json:"template_id"`
	TemplateName   string    `json:"template_name"`
	Category       string    `json:"category"`
	RelevanceScore float64   `json:"relevance_score"`
	Reason         string    `json:"reason"`
}

type WorkflowOptimization struct {
	WorkflowID           uuid.UUID                `json:"workflow_id"`
	Suggestions          []OptimizationSuggestion `json:"suggestions"`
	EstimatedImprovement WorkflowImprovement      `json:"estimated_improvement"`
}

type OptimizationSuggestion struct {
	Type        string `json:"type"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Impact      string `json:"impact"`
	Effort      string `json:"effort"`
}

type WorkflowImprovement struct {
	TimeReduction     float64 `json:"time_reduction_percent"`
	ReliabilityIncrease float64 `json:"reliability_increase_percent"`
	ComplexityReduction int     `json:"complexity_reduction_steps"`
}

type PredictedWorkflow struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Confidence  float64   `json:"confidence"`
	Reason      string    `json:"reason"`
	TemplateData string   `json:"template_data,omitempty"`
}

// Placeholder implementations for analysis methods
func (s *WorkflowAIService) findLongSequentialChains(nodes []interface{}, connections []interface{}) [][]interface{} {
	// Build adjacency from connection "source" -> "target" edges.
	adj := make(map[string][]string)
	targets := make(map[string]bool)
	for _, c := range connections {
		conn, ok := c.(map[string]interface{})
		if !ok {
			continue
		}
		src, _ := conn["source"].(string)
		tgt, _ := conn["target"].(string)
		if src == "" || tgt == "" {
			continue
		}
		adj[src] = append(adj[src], tgt)
		targets[tgt] = true
	}

	// A node with no inbound edge is a chain start.
	var chains [][]interface{}
	for _, n := range nodes {
		node, ok := n.(map[string]interface{})
		if !ok {
			continue
		}
		id, _ := node["id"].(string)
		if id == "" || targets[id] {
			continue
		}

		// Walk the longest simple path from this start.
		chain := []interface{}{node}
		visited := map[string]bool{id: true}
		cur := id
		for {
			nexts := adj[cur]
			if len(nexts) == 0 {
				break
			}
			// Follow the first unvisited successor to stay simple.
			var picked string
			for _, nx := range nexts {
				if !visited[nx] {
					picked = nx
					break
				}
			}
			if picked == "" {
				break
			}
			visited[picked] = true
			chain = append(chain, picked)
			cur = picked
		}
		if len(chain) > 1 {
			chains = append(chains, chain)
		}
	}
	return chains
}

func (s *WorkflowAIService) findRedundantOperations(nodes []interface{}) []struct{OperationType string; Count int} {
	counts := make(map[string]int)
	for _, n := range nodes {
		node, ok := n.(map[string]interface{})
		if !ok {
			continue
		}
		op, _ := node["type"].(string)
		if op == "" {
			op, _ = node["operation_type"].(string)
		}
		if op == "" {
			continue
		}
		counts[op]++
	}

	result := []struct{OperationType string; Count int}{}
	for op, count := range counts {
		if count > 1 {
			result = append(result, struct{OperationType string; Count int}{OperationType: op, Count: count})
		}
	}
	return result
}

func (s *WorkflowAIService) analyzeErrorHandling(nodes []interface{}, connections []interface{}) struct{HasErrorHandlers bool} {
	// Heuristic: a workflow has error handling if any node declares an
	// on_error/error_handler property or any connection leads to an
	// error-handling node.
	errorNodeTypes := map[string]bool{"error": true, "error_handler": true, "notify_error": true}
	for _, n := range nodes {
		node, ok := n.(map[string]interface{})
		if !ok {
			continue
		}
		if onErr, _ := node["on_error"].(string); onErr != "" {
			return struct{HasErrorHandlers bool}{HasErrorHandlers: true}
		}
		if eh, _ := node["error_handler"].(string); eh != "" {
			return struct{HasErrorHandlers bool}{HasErrorHandlers: true}
		}
		op, _ := node["type"].(string)
		if errorNodeTypes[op] {
			return struct{HasErrorHandlers bool}{HasErrorHandlers: true}
		}
	}
	for _, c := range connections {
		conn, ok := c.(map[string]interface{})
		if !ok {
			continue
		}
		if label, _ := conn["label"].(string); strings.EqualFold(label, "error") {
			return struct{HasErrorHandlers bool}{HasErrorHandlers: true}
		}
	}
	return struct{HasErrorHandlers bool}{HasErrorHandlers: false}
}

func (s *WorkflowAIService) findUnusedDataFlow(nodes []interface{}, connections []interface{}) []interface{} {
	// A node is "unused" if it is never referenced as a connection target and
	// never acts as a connection source (i.e. it is disconnected from the flow).
	referenced := make(map[string]bool)
	for _, c := range connections {
		conn, ok := c.(map[string]interface{})
		if !ok {
			continue
		}
		if src, _ := conn["source"].(string); src != "" {
			referenced[src] = true
		}
		if tgt, _ := conn["target"].(string); tgt != "" {
			referenced[tgt] = true
		}
	}

	unused := []interface{}{}
	for _, n := range nodes {
		node, ok := n.(map[string]interface{})
		if !ok {
			continue
		}
		id, _ := node["id"].(string)
		if id == "" || referenced[id] {
			continue
		}
		unused = append(unused, n)
	}
	return unused
}

func (s *WorkflowAIService) estimatePerformanceImprovement(suggestions []OptimizationSuggestion) WorkflowImprovement {
	// Simplified estimation based on suggestion types
	improvement := WorkflowImprovement{}

	for _, suggestion := range suggestions {
		switch suggestion.Type {
		case "parallelization":
			improvement.TimeReduction += 30.0
			improvement.ComplexityReduction += 2
		case "consolidation":
			improvement.TimeReduction += 15.0
			improvement.ComplexityReduction += 1
		case "reliability":
			improvement.ReliabilityIncrease += 25.0
		}
	}

	return improvement
}

func (s *WorkflowAIService) analyzeUserActivity(userID uuid.UUID) map[string]interface{} {
	// Simplified user activity analysis
	return map[string]interface{}{
		"form_submissions": 0,
		"email_volume": 0,
		"task_creation": 0,
		"calendar_events": 0,
	}
}

func (s *WorkflowAIService) identifyWorkflowOpportunities(activity map[string]interface{}) []struct{Type string; Confidence float64} {
	// Simplified opportunity identification
	return []struct{Type string; Confidence float64}{
		{Type: "form_processing", Confidence: 0.8},
	}
}

func (s *WorkflowAIService) generatePredictedWorkflow(opportunity struct{Type string; Confidence float64}) *PredictedWorkflow {
	// Generate predicted workflow based on opportunity type
	switch opportunity.Type {
	case "form_processing":
		return &PredictedWorkflow{
			ID:          uuid.New(),
			Name:        "Automated Form Processing",
			Description: "Automatically process and route form submissions",
			Category:    "business",
			Confidence:  opportunity.Confidence,
			Reason:      "Based on your form submission patterns",
		}
	default:
		return nil
	}
}
