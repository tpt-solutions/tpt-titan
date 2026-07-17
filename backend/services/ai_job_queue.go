package services

import (
	"fmt"
	"log"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Job represents an AI processing job
type Job struct {
	ID          uuid.UUID              `json:"id"`
	UserID      uuid.UUID              `json:"user_id"`
	Type        string                 `json:"type"`        // "document_analysis", "email_categorization", etc.
	Priority    JobPriority            `json:"priority"`
	Status      JobStatus              `json:"status"`
	Input       map[string]interface{} `json:"input"`
	Output      map[string]interface{} `json:"output,omitempty"`
	Error       string                 `json:"error,omitempty"`
	Progress    float64                `json:"progress"`    // 0.0 to 1.0
	CreatedAt   time.Time              `json:"created_at"`
	StartedAt   *time.Time             `json:"started_at,omitempty"`
	CompletedAt *time.Time             `json:"completed_at,omitempty"`
	EstimatedDuration time.Duration    `json:"estimated_duration"`
	WorkerID    string                 `json:"worker_id,omitempty"`
	RetryCount  int                    `json:"retry_count"`
	MaxRetries  int                    `json:"max_retries"`
}

// JobPriority defines job execution priority
type JobPriority int

const (
	PriorityLow JobPriority = iota
	PriorityNormal
	PriorityHigh
	PriorityCritical
)

// JobStatus represents the current state of a job
type JobStatus int

const (
	StatusQueued JobStatus = iota
	StatusRunning
	StatusCompleted
	StatusFailed
	StatusCancelled
)

// JobQueue manages background AI processing jobs
type JobQueue struct {
	mu           sync.RWMutex
	jobs         map[uuid.UUID]*Job
	queue        []*Job // Priority queue
	workers      []*Worker
	maxWorkers   int
	maxRetries   int
	stopChan     chan bool
	jobChan      chan *Job
	resultChan   chan *JobResult
}

// Worker processes jobs
type Worker struct {
	ID       string
	Busy     bool
	JobChan  chan *Job
	QuitChan chan bool
}

// JobResult contains the result of job processing
type JobResult struct {
	JobID   uuid.UUID
	Success bool
	Output  map[string]interface{}
	Error   string
}

// JobQueueConfig configures the job queue
type JobQueueConfig struct {
	MaxWorkers      int
	MaxRetries      int
	WorkerTimeout   time.Duration
	QueueSize       int
	CleanupInterval time.Duration
}

// DefaultJobQueueConfig returns sensible defaults
func DefaultJobQueueConfig() *JobQueueConfig {
	return &JobQueueConfig{
		MaxWorkers:      runtime.NumCPU(),
		MaxRetries:      3,
		WorkerTimeout:   time.Minute * 30,
		QueueSize:       1000,
		CleanupInterval: time.Minute * 10,
	}
}

// NewJobQueue creates a new job queue
func NewJobQueue(config *JobQueueConfig) *JobQueue {
	if config == nil {
		config = DefaultJobQueueConfig()
	}

	jq := &JobQueue{
		jobs:       make(map[uuid.UUID]*Job),
		queue:      make([]*Job, 0),
		workers:    make([]*Worker, config.MaxWorkers),
		maxWorkers: config.MaxWorkers,
		maxRetries: config.MaxRetries,
		stopChan:   make(chan bool),
		jobChan:    make(chan *Job, config.QueueSize),
		resultChan: make(chan *JobResult, config.QueueSize),
	}

	// Start workers
	jq.startWorkers()

	// Start job dispatcher
	go jq.dispatchJobs()

	// Start result processor
	go jq.processResults()

	// Start cleanup routine
	go jq.cleanupRoutine(config.CleanupInterval)

	return jq
}

// SubmitJob adds a new job to the queue
func (jq *JobQueue) SubmitJob(userID uuid.UUID, jobType string, input map[string]interface{}) (*Job, error) {
	jq.mu.Lock()
	defer jq.mu.Unlock()

	job := &Job{
		ID:               uuid.New(),
		UserID:           userID,
		Type:             jobType,
		Priority:         PriorityNormal,
		Status:           StatusQueued,
		Input:            input,
		Progress:         0.0,
		CreatedAt:        time.Now(),
		EstimatedDuration: jq.estimateJobDuration(jobType, input),
		MaxRetries:       jq.maxRetries,
	}

	// Store job
	jq.jobs[job.ID] = job

	// Add to priority queue
	jq.insertIntoQueue(job)

	log.Printf("Job submitted: %s (type: %s, user: %s)", job.ID, jobType, userID)

	return job, nil
}

// GetJob retrieves a job by ID
func (jq *JobQueue) GetJob(jobID uuid.UUID) (*Job, bool) {
	jq.mu.RLock()
	defer jq.mu.RUnlock()

	job, exists := jq.jobs[jobID]
	return job, exists
}

// GetUserJobs retrieves all jobs for a user
func (jq *JobQueue) GetUserJobs(userID uuid.UUID, limit int) []*Job {
	jq.mu.RLock()
	defer jq.mu.RUnlock()

	jobs := make([]*Job, 0)
	for _, job := range jq.jobs {
		if job.UserID == userID {
			jobs = append(jobs, job)
			if limit > 0 && len(jobs) >= limit {
				break
			}
		}
	}

	return jobs
}

// CancelJob cancels a queued job
func (jq *JobQueue) CancelJob(jobID uuid.UUID) error {
	jq.mu.Lock()
	defer jq.mu.Unlock()

	job, exists := jq.jobs[jobID]
	if !exists {
		return fmt.Errorf("job not found: %s", jobID)
	}

	if job.Status == StatusRunning || job.Status == StatusQueued {
		job.Status = StatusCancelled
		log.Printf("Job cancelled: %s", jobID)
		return nil
	}

	return fmt.Errorf("cannot cancel job in status: %v", job.Status)
}

// GetQueueStats returns queue statistics
func (jq *JobQueue) GetQueueStats() map[string]interface{} {
	jq.mu.RLock()
	defer jq.mu.RUnlock()

	stats := make(map[string]interface{})
	statusCounts := make(map[JobStatus]int)
	priorityCounts := make(map[JobPriority]int)

	totalJobs := len(jq.jobs)
	queuedJobs := len(jq.queue)

	for _, job := range jq.jobs {
		statusCounts[job.Status]++
		priorityCounts[job.Priority]++
	}

	// Calculate average wait time for queued jobs
	var totalWaitTime time.Duration
	waitingJobs := 0
	for _, job := range jq.queue {
		if job.Status == StatusQueued {
			totalWaitTime += time.Since(job.CreatedAt)
			waitingJobs++
		}
	}
	avgWaitTime := time.Duration(0)
	if waitingJobs > 0 {
		avgWaitTime = totalWaitTime / time.Duration(waitingJobs)
	}

	stats["total_jobs"] = totalJobs
	stats["queued_jobs"] = queuedJobs
	stats["status_counts"] = statusCounts
	stats["priority_counts"] = priorityCounts
	stats["active_workers"] = jq.getActiveWorkerCount()
	stats["avg_wait_time"] = avgWaitTime.String()

	return stats
}

// Stop gracefully shuts down the job queue
func (jq *JobQueue) Stop() {
	log.Println("Stopping job queue...")

	// Signal stop to all workers
	close(jq.stopChan)

	// Wait for workers to finish
	jq.mu.Lock()
	for _, worker := range jq.workers {
		if worker != nil {
			close(worker.QuitChan)
		}
	}
	jq.mu.Unlock()

	log.Println("Job queue stopped")
}

// startWorkers initializes worker goroutines
func (jq *JobQueue) startWorkers() {
	for i := 0; i < jq.maxWorkers; i++ {
		worker := &Worker{
			ID:       fmt.Sprintf("worker-%d", i+1),
			Busy:     false,
			JobChan:  make(chan *Job),
			QuitChan: make(chan bool),
		}

		jq.workers[i] = worker
		go jq.workerRoutine(worker)
	}
}

// workerRoutine handles job processing for a worker
func (jq *JobQueue) workerRoutine(worker *Worker) {
	log.Printf("Worker %s started", worker.ID)

	for {
		select {
		case job := <-worker.JobChan:
			jq.processJob(worker, job)

		case <-worker.QuitChan:
			log.Printf("Worker %s stopped", worker.ID)
			return

		case <-jq.stopChan:
			log.Printf("Worker %s shutting down", worker.ID)
			return
		}
	}
}

// dispatchJobs distributes jobs to available workers
func (jq *JobQueue) dispatchJobs() {
	for {
		select {
		case <-jq.stopChan:
			return

		default:
			jq.mu.Lock()
			if len(jq.queue) > 0 {
				// Find available worker
				for _, worker := range jq.workers {
					if !worker.Busy {
						// Get highest priority job
						job := jq.queue[0]
						jq.queue = jq.queue[1:]

						worker.Busy = true
						job.Status = StatusRunning
						job.WorkerID = worker.ID
						now := time.Now()
						job.StartedAt = &now

						log.Printf("Dispatching job %s to worker %s", job.ID, worker.ID)
						worker.JobChan <- job
						break
					}
				}
			}
			jq.mu.Unlock()

			time.Sleep(time.Millisecond * 100) // Prevent busy waiting
		}
	}
}

// processJob executes a job
func (jq *JobQueue) processJob(worker *Worker, job *Job) {
	defer func() {
		worker.Busy = false
	}()

	log.Printf("Worker %s processing job %s", worker.ID, job.ID)

	// Execute job based on type
	result := jq.executeJob(job)

	// Send result
	select {
	case jq.resultChan <- result:
	case <-time.After(time.Second * 5):
		log.Printf("Failed to send result for job %s", job.ID)
	}
}

// executeJob performs the actual job processing
func (jq *JobQueue) executeJob(job *Job) *JobResult {
	result := &JobResult{
		JobID:   job.ID,
		Success: false,
	}

	defer func() {
		if r := recover(); r != nil {
			result.Error = fmt.Sprintf("Job panicked: %v", r)
			log.Printf("Job %s panicked: %v", job.ID, r)
		}
	}()

	// Update job progress
	jq.updateJobProgress(job.ID, 0.1)

	// Process based on job type
	switch job.Type {
	case "document_analysis":
		result = jq.processDocumentAnalysis(job)
	case "email_categorization":
		result = jq.processEmailCategorization(job)
	case "speech_synthesis":
		result = jq.processSpeechSynthesis(job)
	case "workflow_optimization":
		result = jq.processWorkflowOptimization(job)
	default:
		result.Error = fmt.Sprintf("Unknown job type: %s", job.Type)
	}

	jq.updateJobProgress(job.ID, 1.0)

	return result
}

// processResults handles job completion results
func (jq *JobQueue) processResults() {
	for {
		select {
		case result := <-jq.resultChan:
			jq.handleJobResult(result)

		case <-jq.stopChan:
			return
		}
	}
}

// handleJobResult processes completed job results
func (jq *JobQueue) handleJobResult(result *JobResult) {
	jq.mu.Lock()
	defer jq.mu.Unlock()

	job, exists := jq.jobs[result.JobID]
	if !exists {
		log.Printf("Received result for unknown job: %s", result.JobID)
		return
	}

	now := time.Now()
	job.CompletedAt = &now

	if result.Success {
		job.Status = StatusCompleted
		job.Output = result.Output
		log.Printf("Job completed: %s", job.ID)
	} else {
		job.Status = StatusFailed
		job.Error = result.Error
		job.RetryCount++

		// Check if we should retry
		if job.RetryCount < job.MaxRetries {
			job.Status = StatusQueued
			jq.insertIntoQueue(job)
			log.Printf("Job failed, retrying: %s (attempt %d/%d)", job.ID, job.RetryCount, job.MaxRetries)
		} else {
			log.Printf("Job failed permanently: %s (error: %s)", job.ID, result.Error)
		}
	}
}

// updateJobProgress updates job progress
func (jq *JobQueue) updateJobProgress(jobID uuid.UUID, progress float64) {
	jq.mu.Lock()
	defer jq.mu.Unlock()

	if job, exists := jq.jobs[jobID]; exists {
		job.Progress = progress
	}
}

// insertIntoQueue adds a job to the priority queue
func (jq *JobQueue) insertIntoQueue(job *Job) {
	// Simple insertion sort by priority (higher priority first)
	inserted := false
	for i, queuedJob := range jq.queue {
		if job.Priority > queuedJob.Priority {
			jq.queue = append(jq.queue[:i], append([]*Job{job}, jq.queue[i:]...)...)
			inserted = true
			break
		}
	}

	if !inserted {
		jq.queue = append(jq.queue, job)
	}
}

// estimateJobDuration estimates how long a job will take
func (jq *JobQueue) estimateJobDuration(jobType string, input map[string]interface{}) time.Duration {
	baseEstimates := map[string]time.Duration{
		"document_analysis":    time.Minute * 2,
		"email_categorization": time.Second * 30,
		"speech_synthesis":     time.Minute * 1,
		"workflow_optimization": time.Minute * 3,
	}

	estimate, exists := baseEstimates[jobType]
	if !exists {
		estimate = time.Minute * 1
	}

	// Adjust based on input size
	if content, ok := input["content"].(string); ok {
		contentLength := len(content)
		if contentLength > 10000 {
			estimate = time.Duration(float64(estimate) * 1.5)
		} else if contentLength > 50000 {
			estimate = time.Duration(float64(estimate) * 2.0)
		}
	}

	return estimate
}

// getActiveWorkerCount returns the number of busy workers
func (jq *JobQueue) getActiveWorkerCount() int {
	count := 0
	for _, worker := range jq.workers {
		if worker.Busy {
			count++
		}
	}
	return count
}

// cleanupRoutine periodically cleans up old completed jobs
func (jq *JobQueue) cleanupRoutine(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			jq.cleanupOldJobs()
		case <-jq.stopChan:
			return
		}
	}
}

// cleanupOldJobs removes jobs older than retention period
func (jq *JobQueue) cleanupOldJobs() {
	jq.mu.Lock()
	defer jq.mu.Unlock()

	now := time.Now()
	retentionPeriod := time.Hour * 24 * 7 // Keep jobs for 7 days

	cleaned := 0
	for id, job := range jq.jobs {
		if job.Status == StatusCompleted || job.Status == StatusFailed {
			if now.Sub(job.CreatedAt) > retentionPeriod {
				delete(jq.jobs, id)
				cleaned++
			}
		}
	}

	if cleaned > 0 {
		log.Printf("Cleaned up %d old jobs", cleaned)
	}
}

// Job processing implementations.
//
// These derive results from the actual job input payload so that different
// inputs produce different (and meaningful) outputs instead of canned data.

// extractText pulls a textual body out of the job input under a few common
// keys ("content", "text", "body", "document"). It returns "" when absent.
func extractJobText(input map[string]interface{}) string {
	for _, key := range []string{"content", "text", "body", "document", "transcript"} {
		if v, ok := input[key]; ok {
			if s, ok := v.(string); ok {
				return s
			}
		}
	}
	return ""
}

var topicKeywords = []struct {
	topic string
	words []string
}{
	{"AI", []string{"ai", "model", "neural", "machine learning", "llm", "gpt"}},
	{"productivity", []string{"task", "workflow", "automation", "schedule", "goal", "productive"}},
	{"security", []string{"security", "encryption", "password", "auth", "vulnerability", "key"}},
	{"finance", []string{"budget", "invoice", "payment", "finance", "revenue", "cost"}},
	{"meetings", []string{"meeting", "agenda", "standup", "call", "sync", "review"}},
}

func (jq *JobQueue) processDocumentAnalysis(job *Job) *JobResult {
	text := extractJobText(job.Input)

	words := strings.Fields(text)
	wordCount := len(words)

	topicHits := map[string]int{}
	for _, w := range words {
		lw := strings.ToLower(strings.TrimRight(w, ",.;:!?"))
		for _, tk := range topicKeywords {
			for _, kw := range tk.words {
				if lw == kw || strings.Contains(lw, kw) {
					topicHits[tk.topic]++
				}
			}
		}
	}

	// Rank topics by hit count and keep up to 5 with at least one match.
	topics := make([]string, 0, len(topicHits))
	type scored struct {
		topic string
		score int
	}
	ranked := make([]scored, 0, len(topicHits))
	for t, s := range topicHits {
		ranked = append(ranked, scored{t, s})
	}
	sort.Slice(ranked, func(i, j int) bool { return ranked[i].score > ranked[j].score })
	for _, r := range ranked {
		topics = append(topics, r.topic)
		if len(topics) >= 5 {
			break
		}
	}
	if len(topics) == 0 {
		topics = []string{"general"}
	}

	charCount := len(text)
	readMinutes := 0
	if wordCount > 0 {
		readMinutes = (wordCount + 199) / 200
	}

	return &JobResult{
		JobID:   job.ID,
		Success: true,
		Output: map[string]interface{}{
			"analysis":        "Document analyzed successfully",
			"word_count":      wordCount,
			"character_count": charCount,
			"estimated_read_minutes": readMinutes,
			"key_topics":      topics,
		},
	}
}

// emailCategoryKeywords maps a category to the keywords that indicate it.
var emailCategoryKeywords = map[string][]string{
	"urgent":   {"urgent", "asap", "immediately", "critical", "deadline", "emergency"},
	"meeting":  {"meeting", "agenda", "invite", "call", "standup", "sync", "conference"},
	"work":     {"project", "task", "client", "report", "review", "sprint", "ticket"},
	"personal": {"family", "weekend", "friend", "vacation", "birthday", "holiday"},
	"finance":  {"invoice", "payment", "billing", "receipt", "subscription", "refund"},
}

func (jq *JobQueue) processEmailCategorization(job *Job) *JobResult {
	subject, _ := job.Input["subject"].(string)
	from, _ := job.Input["from"].(string)
	body := extractJobText(job.Input)

	corpus := strings.ToLower(subject + " " + from + " " + body)

	categories := make([]string, 0, len(emailCategoryKeywords))
	bestScore := 0
	for category, keywords := range emailCategoryKeywords {
		score := 0
		for _, kw := range keywords {
			if strings.Contains(corpus, kw) {
				score++
			}
		}
		if score > 0 {
			// "work" is a weaker default than more specific categories.
			if category == "work" {
				score = score - 0
			}
			categories = append(categories, category)
			if score > bestScore {
				bestScore = score
			}
		}
	}
	if len(categories) == 0 {
		categories = []string{"other"}
	}

	confidence := 0.5
	if bestScore > 0 {
		confidence = float64(bestScore) / float64(bestScore+1)
		if confidence > 0.99 {
			confidence = 0.99
		}
	}

	return &JobResult{
		JobID:   job.ID,
		Success: true,
		Output: map[string]interface{}{
			"categories": categories,
			"confidence": confidence,
		},
	}
}

func (jq *JobQueue) processSpeechSynthesis(job *Job) *JobResult {
	text, _ := job.Input["text"].(string)
	if text == "" {
		return &JobResult{
			JobID:   job.ID,
			Success: false,
			Error:   "speech synthesis requires a 'text' input field",
		}
	}

	// Estimate duration at ~150 words per minute.
	words := len(strings.Fields(text))
	durationSec := (words * 60) / 150
	if durationSec < 1 {
		durationSec = 1
	}

	// No real TTS engine is wired here; expose the requested text and the
	// estimated duration honestly instead of fabricating an audio URL.
	return &JobResult{
		JobID:   job.ID,
		Success: true,
		Output: map[string]interface{}{
			"text":           text,
			"estimated_duration_seconds": durationSec,
			"synthesized":    false,
			"note":           "no TTS provider configured; duration is an estimate",
		},
	}
}

func (jq *JobQueue) processWorkflowOptimization(job *Job) *JobResult {
	nodesRaw, hasNodes := job.Input["nodes"]
	connectionsRaw, hasConns := job.Input["connections"]

	if !hasNodes || !hasConns {
		return &JobResult{
			JobID:   job.ID,
			Success: false,
			Error:   "workflow optimization requires 'nodes' and 'connections' input",
		}
	}

	nodes, okN := nodesRaw.([]interface{})
	connections, okC := connectionsRaw.([]interface{})
	if !okN || !okC {
		return &JobResult{
			JobID:   job.ID,
			Success: false,
			Error:   "workflow optimization input has invalid node/connection types",
		}
	}

	nodeByID := map[string]map[string]interface{}{}
	nodeTypeCount := map[string]int{}
	for _, n := range nodes {
		nm, ok := n.(map[string]interface{})
		if !ok {
			continue
		}
		id, _ := nm["id"].(string)
		if id == "" {
			continue
		}
		nodeByID[id] = nm
		nt, _ := nm["type"].(string)
		if nt != "" {
			nodeTypeCount[nt]++
		}
	}

	// Build adjacency from connections (source -> target).
	adj := map[string][]string{}
	inDeg := map[string]int{}
	for _, c := range connections {
		cm, ok := c.(map[string]interface{})
		if !ok {
			continue
		}
		src, _ := cm["source"].(string)
		tgt, _ := cm["target"].(string)
		if src == "" || tgt == "" {
			continue
		}
		adj[src] = append(adj[src], tgt)
		inDeg[tgt]++
	}

	// Long chains: walk from each node that has no inbound edge.
	var suggestions []string
	longest := 0
	for id := range nodeByID {
		if inDeg[id] > 0 {
			continue
		}
		chain := 0
		cur := id
		visited := map[string]bool{}
		for cur != "" && !visited[cur] {
			visited[cur] = true
			chain++
			next := ""
			for _, t := range adj[cur] {
				next = t
				break
			}
			cur = next
		}
		if chain > longest {
			longest = chain
		}
		if chain > 5 {
			suggestions = append(suggestions, fmt.Sprintf("Flatten the sequential chain of %d steps starting at node %s", chain, id))
		}
	}

	// Redundant operations: node types appearing more than once.
	for nt, count := range nodeTypeCount {
		if count > 1 {
			suggestions = append(suggestions, fmt.Sprintf("Merge %d redundant '%s' nodes", count, nt))
		}
	}

	if len(suggestions) == 0 {
		suggestions = append(suggestions, "Workflow looks well structured; no optimizations detected")
	}

	return &JobResult{
		JobID:   job.ID,
		Success: true,
		Output: map[string]interface{}{
			"suggestions":          suggestions,
			"node_count":           len(nodeByID),
			"longest_chain_length": longest,
		},
	}
}
