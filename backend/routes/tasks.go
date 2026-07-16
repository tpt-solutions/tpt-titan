package routes

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"tpt-titan/backend/models"
	"tpt-titan/backend/services"
)

var taskService *services.TaskService

// InitTaskService initializes the task service (called from server setup)
func InitTaskService(db *sql.DB) {
	taskService = services.NewTaskService(db)
}

func getTaskUserID(c *gin.Context) (uuid.UUID, bool) {
	v, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return uuid.Nil, false
	}
	userID, ok := v.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user context"})
		return uuid.Nil, false
	}
	return userID, true
}

func taskToResponse(t *models.Task) map[string]interface{} {
	resp := map[string]interface{}{
		"id":          t.ID.String(),
		"title":       t.Title,
		"description": t.Description,
		"status":      t.Status,
		"priority":    t.Priority,
		"assignedTo":  t.AssignedTo,
		"projectId":   nil,
		"dueDate":     nil,
		"tags":        t.Tags,
		"subtasks":    []map[string]interface{}{},
		"createdAt":   t.CreatedAt,
	}
	if t.ProjectID != nil {
		resp["projectId"] = t.ProjectID.String()
	}
	if t.DueDate != nil {
		resp["dueDate"] = t.DueDate
	}
	if t.Subtasks != nil {
		subtasks := make([]map[string]interface{}, 0, len(t.Subtasks))
		for _, st := range t.Subtasks {
			subtasks = append(subtasks, map[string]interface{}{
				"id":        st.ID.String(),
				"title":     st.Title,
				"completed": st.Completed,
			})
		}
		resp["subtasks"] = subtasks
	}
	return resp
}

// GetTasks returns all tasks for the current user
func GetTasks(c *gin.Context) {
	userID, ok := getTaskUserID(c)
	if !ok {
		return
	}
	status := c.Query("status")
	projectID := c.Query("project_id")

	tasks, err := taskService.GetTasks(userID, status, projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}

	resp := make([]map[string]interface{}, 0, len(tasks))
	for i := range tasks {
		resp = append(resp, taskToResponse(&tasks[i]))
	}
	c.JSON(http.StatusOK, gin.H{"tasks": resp})
}

// GetTask returns a single task
func GetTask(c *gin.Context) {
	userID, ok := getTaskUserID(c)
	if !ok {
		return
	}
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	task, err := taskService.GetTask(userID, taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve task"})
		return
	}
	if task == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"task": taskToResponse(task)})
}

// CreateTask creates a new task
func CreateTask(c *gin.Context) {
	userID, ok := getTaskUserID(c)
	if !ok {
		return
	}
	var req models.TaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := taskService.CreateTask(userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"task": taskToResponse(task)})
}

// UpdateTask updates an existing task
func UpdateTask(c *gin.Context) {
	userID, ok := getTaskUserID(c)
	if !ok {
		return
	}
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}
	var req models.TaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := taskService.UpdateTask(userID, taskID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"task": taskToResponse(task)})
}

// DeleteTask removes a task
func DeleteTask(c *gin.Context) {
	userID, ok := getTaskUserID(c)
	if !ok {
		return
	}
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}
	if err := taskService.DeleteTask(userID, taskID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

// UpdateTaskStatus changes only the status of a task
func UpdateTaskStatus(c *gin.Context) {
	userID, ok := getTaskUserID(c)
	if !ok {
		return
	}
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}
	var body struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.Status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status is required"})
		return
	}
	task, err := taskService.UpdateTaskStatus(userID, taskID, body.Status)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"task": taskToResponse(task)})
}

// GetProjects returns all projects for the current user
func GetProjects(c *gin.Context) {
	userID, ok := getTaskUserID(c)
	if !ok {
		return
	}
	projects, err := taskService.GetProjects(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve projects"})
		return
	}
	resp := make([]map[string]interface{}, 0, len(projects))
	for _, p := range projects {
		resp = append(resp, map[string]interface{}{
			"id":         p.ID.String(),
			"name":       p.Name,
			"color":      p.Color,
			"createdAt":  p.CreatedAt,
		})
	}
	c.JSON(http.StatusOK, gin.H{"projects": resp})
}

// CreateProject creates a new project
func CreateProject(c *gin.Context) {
	userID, ok := getTaskUserID(c)
	if !ok {
		return
	}
	var req models.ProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	project, err := taskService.CreateProject(userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"project": map[string]interface{}{
		"id":         project.ID.String(),
		"name":       project.Name,
		"color":      project.Color,
		"createdAt":  project.CreatedAt,
	}})
}
