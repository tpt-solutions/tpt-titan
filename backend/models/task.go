package models

import (
	"time"

	"github.com/google/uuid"
)

// Task represents a task in the task management system
type Task struct {
	ID          uuid.UUID     `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID      uuid.UUID     `gorm:"type:uuid;not null;index" json:"user_id"`
	ProjectID   *uuid.UUID    `gorm:"type:uuid;index" json:"project_id,omitempty"`
	Title       string        `gorm:"size:255;not null" json:"title"`
	Description string        `gorm:"type:text" json:"description"`
	Status      string        `gorm:"size:20;not null;default:'todo'" json:"status"` // todo, in-progress, review, done
	Priority    string        `gorm:"size:10;not null;default:'medium'" json:"priority"`
	AssignedTo  string        `gorm:"size:255" json:"assigned_to"`
	DueDate     *time.Time    `json:"due_date,omitempty"`
	Tags        string        `gorm:"type:text" json:"-"` // JSON-encoded array
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	TagsList    []string      `gorm:"-" json:"tags,omitempty"`
	Subtasks    []TaskSubtask ` gorm:"foreignKey:TaskID" json:"subtasks,omitempty"`
}

// TaskSubtask represents a subtask belonging to a task
type TaskSubtask struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TaskID    uuid.UUID `gorm:"type:uuid;not null;index" json:"task_id"`
	Title     string    `gorm:"size:255;not null" json:"title"`
	Completed bool      `gorm:"default:false" json:"completed"`
	OrderIdx  int       `gorm:"default:0" json:"order"`
	CreatedAt time.Time `json:"created_at"`
}

// Project represents a project grouping tasks together
type Project struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Name      string    `gorm:"size:255;not null" json:"name"`
	Color     string    `gorm:"size:20;default:'blue'" json:"color"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TaskTag represents a tag associated with a task (join table)
type TaskTag struct {
	TaskID uuid.UUID `gorm:"type:uuid;primary_key" json:"-"`
	Tag    string    `gorm:"size:64;primary_key" json:"tag"`
}

// TableName overrides for GORM
func (Task) TableName() string        { return "tasks" }
func (TaskSubtask) TableName() string { return "task_subtasks" }
func (Project) TableName() string     { return "projects" }
func (TaskTag) TableName() string     { return "task_tags" }

// TaskRequest is the payload for creating/updating a task
type TaskRequest struct {
	ProjectID   *uuid.UUID       `json:"project_id"`
	Title       string           `json:"title"`
	Description string           `json:"description"`
	Status      string           `json:"status"`
	Priority    string           `json:"priority"`
	AssignedTo  string           `json:"assigned_to"`
	DueDate     *time.Time       `json:"due_date"`
	Tags        []string         `json:"tags"`
	Subtasks    []SubtaskRequest `json:"subtasks"`
}

// SubtaskRequest is the payload for a subtask
type SubtaskRequest struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// ProjectRequest is the payload for creating/updating a project
type ProjectRequest struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}
