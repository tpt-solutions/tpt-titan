package services

import (
	"fmt"
	"log"
	"runtime"
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

// Job processing implementations (simplified placeholders)

func (jq *JobQueue) processDocumentAnalysis(job *Job) *JobResult {
	// Placeholder implementation
	time.Sleep(time.Second * 2) // Simulate processing time

	return &JobResult{
		JobID:   job.ID,
		Success: true,
		Output: map[string]interface{}{
			"analysis": "Document analyzed successfully",
			"word_count": 1250,
			"key_topics": []string{"AI", "productivity", "automation"},
		},
	}
}

func (jq *JobQueue) processEmailCategorization(job *Job) *JobResult {
	// Placeholder implementation
	time.Sleep(time.Second * 1)

	return &JobResult{
		JobID:   job.ID,
		Success: true,
		Output: map[string]interface{}{
			"categories": []string{"work", "meeting", "urgent"},
			"confidence": 0.95,
		},
	}
}

func (jq *JobQueue) processSpeechSynthesis(job *Job) *JobResult {
	// Placeholder implementation
	time.Sleep(time.Second * 3)

	return &JobResult{
		JobID:   job.ID,
		Success: true,
		Output: map[string]interface{}{
			"audio_url": "https://example.com/audio/generated.mp3",
			"duration": 45,
		},
	}
}

func (jq *JobQueue) processWorkflowOptimization(job *Job) *JobResult {
	// Placeholder implementation
	time.Sleep(time.Second * 2)

	return &JobResult{
		JobID:   job.ID,
		Success: true,
		Output: map[string]interface{}{
			"suggestions": []string{
				"Add error handling to step 3",
				"Parallelize steps 4 and 5",
				"Add validation before step 2",
			},
			"estimated_improvement": "35% faster execution",
		},
	}
}
