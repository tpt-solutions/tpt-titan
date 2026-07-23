package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"tpt-titan/backend/models"
)

// TaskService handles task and project persistence operations
type TaskService struct {
	db *sql.DB
}

// NewTaskService creates a new task service
func NewTaskService(db *sql.DB) *TaskService {
	return &TaskService{db: db}
}

// CreateTask creates a new task (with optional subtasks) for a user
func (s *TaskService) CreateTask(userID uuid.UUID, req models.TaskRequest) (*models.Task, error) {
	if req.Title == "" {
		return nil, fmt.Errorf("task title is required")
	}

	if req.Status == "" {
		req.Status = "todo"
	}
	if req.Priority == "" {
		req.Priority = "medium"
	}

	taskID := uuid.New()
	now := time.Now()

	query := `
		INSERT INTO tasks (id, user_id, project_id, title, description, status, priority, assigned_to, due_date, tags, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	_, err := s.db.Exec(query,
		taskID, userID, req.ProjectID, req.Title, req.Description,
		req.Status, req.Priority, req.AssignedTo, req.DueDate,
		encodeTags(req.Tags), now, now,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	if err := s.replaceSubtasks(taskID, req.Subtasks); err != nil {
		return nil, err
	}

	return s.GetTask(userID, taskID)
}

// GetTask retrieves a single task with its subtasks
func (s *TaskService) GetTask(userID, taskID uuid.UUID) (*models.Task, error) {
	var task models.Task
	var dueDate sql.NullTime
	var projectID uuid.NullUUID
	var tagsJSON sql.NullString

	query := `
		SELECT id, user_id, project_id, title, description, status, priority, assigned_to, due_date, tags, created_at, updated_at
		FROM tasks WHERE id = $1 AND user_id = $2
	`
	err := s.db.QueryRow(query, taskID, userID).Scan(
		&task.ID, &task.UserID, &projectID, &task.Title, &task.Description,
		&task.Status, &task.Priority, &task.AssignedTo, &dueDate, &tagsJSON,
		&task.CreatedAt, &task.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	if projectID.Valid {
		task.ProjectID = &projectID.UUID
	}
	if dueDate.Valid {
		task.DueDate = &dueDate.Time
	}
	task.TagsList = decodeTags(tagsJSON)

	subtasks, err := s.getSubtasks(taskID)
	if err != nil {
		return nil, err
	}
	task.Subtasks = subtasks

	return &task, nil
}

// GetTasks retrieves all tasks for a user, optionally filtered by status/project
func (s *TaskService) GetTasks(userID uuid.UUID, status, projectID string) ([]models.Task, error) {
	query := `
		SELECT id, user_id, project_id, title, description, status, priority, assigned_to, due_date, tags, created_at, updated_at
		FROM tasks WHERE user_id = $1
	`
	args := []interface{}{userID}
	argIdx := 2

	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, status)
		argIdx++
	}
	if projectID != "" {
		if pid, err := uuid.Parse(projectID); err == nil {
			query += fmt.Sprintf(" AND project_id = $%d", argIdx)
			args = append(args, pid)
			argIdx++
		}
	}

	query += " ORDER BY created_at DESC"

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		var dueDate sql.NullTime
		var projID uuid.NullUUID
		var tagsJSON sql.NullString

		if err := rows.Scan(
			&task.ID, &task.UserID, &projID, &task.Title, &task.Description,
			&task.Status, &task.Priority, &task.AssignedTo, &dueDate, &tagsJSON,
			&task.CreatedAt, &task.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}

		if projID.Valid {
			task.ProjectID = &projID.UUID
		}
		if dueDate.Valid {
			task.DueDate = &dueDate.Time
		}
		task.TagsList = decodeTags(tagsJSON)
		tasks = append(tasks, task)
	}

	for i := range tasks {
		subtasks, err := s.getSubtasks(tasks[i].ID)
		if err != nil {
			return nil, err
		}
		tasks[i].Subtasks = subtasks
	}

	return tasks, nil
}

// UpdateTask updates an existing task
func (s *TaskService) UpdateTask(userID, taskID uuid.UUID, req models.TaskRequest) (*models.Task, error) {
	existing, err := s.GetTask(userID, taskID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, fmt.Errorf("task not found")
	}

	if req.Status == "" {
		req.Status = existing.Status
	}
	if req.Priority == "" {
		req.Priority = existing.Priority
	}

	query := `
		UPDATE tasks
		SET project_id = $1, title = $2, description = $3, status = $4, priority = $5,
		    assigned_to = $6, due_date = $7, tags = $8, updated_at = $9
		WHERE id = $10 AND user_id = $11
	`
	now := time.Now()
	result, err := s.db.Exec(query,
		req.ProjectID, req.Title, req.Description, req.Status, req.Priority,
		req.AssignedTo, req.DueDate, encodeTags(req.Tags), now, taskID, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		return nil, fmt.Errorf("task not found")
	}

	if req.Subtasks != nil {
		if err := s.replaceSubtasks(taskID, req.Subtasks); err != nil {
			return nil, err
		}
	}

	return s.GetTask(userID, taskID)
}

// DeleteTask removes a task and its subtasks
func (s *TaskService) DeleteTask(userID, taskID uuid.UUID) error {
	_, err := s.db.Exec("DELETE FROM task_subtasks WHERE task_id = $1", taskID)
	if err != nil {
		return fmt.Errorf("failed to delete subtasks: %w", err)
	}
	result, err := s.db.Exec("DELETE FROM tasks WHERE id = $1 AND user_id = $2", taskID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}

// UpdateTaskStatus changes only the status of a task
func (s *TaskService) UpdateTaskStatus(userID, taskID uuid.UUID, status string) (*models.Task, error) {
	result, err := s.db.Exec(
		"UPDATE tasks SET status = $1, updated_at = $2 WHERE id = $3 AND user_id = $4",
		status, time.Now(), taskID, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update task status: %w", err)
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		return nil, fmt.Errorf("task not found")
	}
	return s.GetTask(userID, taskID)
}

// CreateProject creates a new project
func (s *TaskService) CreateProject(userID uuid.UUID, req models.ProjectRequest) (*models.Project, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("project name is required")
	}
	if req.Color == "" {
		req.Color = "blue"
	}

	project := models.Project{
		ID:        uuid.New(),
		UserID:    userID,
		Name:      req.Name,
		Color:     req.Color,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	query := `
		INSERT INTO projects (id, user_id, name, color, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := s.db.Exec(query, project.ID, project.UserID, project.Name, project.Color, project.CreatedAt, project.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}
	return &project, nil
}

// GetProjects retrieves all projects for a user
func (s *TaskService) GetProjects(userID uuid.UUID) ([]models.Project, error) {
	query := `
		SELECT id, user_id, name, color, created_at, updated_at
		FROM projects WHERE user_id = $1 ORDER BY created_at DESC
	`
	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query projects: %w", err)
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var project models.Project
		if err := rows.Scan(
			&project.ID, &project.UserID, &project.Name, &project.Color,
			&project.CreatedAt, &project.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan project: %w", err)
		}
		projects = append(projects, project)
	}
	return projects, nil
}

// --- helpers ---

func (s *TaskService) getSubtasks(taskID uuid.UUID) ([]models.TaskSubtask, error) {
	query := `
		SELECT id, task_id, title, completed, order_idx, created_at
		FROM task_subtasks WHERE task_id = $1 ORDER BY order_idx ASC, created_at ASC
	`
	rows, err := s.db.Query(query, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to query subtasks: %w", err)
	}
	defer rows.Close()

	var subtasks []models.TaskSubtask
	for rows.Next() {
		var st models.TaskSubtask
		if err := rows.Scan(&st.ID, &st.TaskID, &st.Title, &st.Completed, &st.OrderIdx, &st.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan subtask: %w", err)
		}
		subtasks = append(subtasks, st)
	}
	return subtasks, nil
}

func (s *TaskService) replaceSubtasks(taskID uuid.UUID, subtasks []models.SubtaskRequest) error {
	if _, err := s.db.Exec("DELETE FROM task_subtasks WHERE task_id = $1", taskID); err != nil {
		return fmt.Errorf("failed to clear subtasks: %w", err)
	}

	for i, st := range subtasks {
		if st.Title == "" {
			continue
		}
		id := uuid.New()
		query := `
			INSERT INTO task_subtasks (id, task_id, title, completed, order_idx, created_at)
			VALUES ($1, $2, $3, $4, $5, $6)
		`
		if _, err := s.db.Exec(query, id, taskID, st.Title, st.Completed, i, time.Now()); err != nil {
			return fmt.Errorf("failed to create subtask: %w", err)
		}
	}
	return nil
}

func encodeTags(tags []string) string {
	if len(tags) == 0 {
		return "[]"
	}
	b, err := json.Marshal(tags)
	if err != nil {
		return "[]"
	}
	return string(b)
}

func decodeTags(s sql.NullString) []string {
	if !s.Valid || strings.TrimSpace(s.String) == "" {
		return []string{}
	}
	var tags []string
	if err := json.Unmarshal([]byte(s.String), &tags); err != nil {
		return []string{}
	}
	return tags
}
