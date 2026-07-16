package routes

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"tpt-titan/backend/models"
	"tpt-titan/backend/services"
)

var fileSyncService *services.FileSyncService

// InitFileSyncService initializes the file sync service (called from server setup)
func InitFileSyncService(db *sql.DB) {
	fileSyncService = services.NewFileSyncService(db)
}

func getFileSyncUserID(c *gin.Context) (uuid.UUID, bool) {
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

// GetSyncFolders returns all sync folders for the current user
func GetSyncFolders(c *gin.Context) {
	userID, ok := getFileSyncUserID(c)
	if !ok {
		return
	}

	folders, err := fileSyncService.GetSyncFolders(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve sync folders"})
		return
	}

	resp := make([]map[string]interface{}, 0, len(folders))
	for _, f := range folders {
		resp = append(resp, map[string]interface{}{
			"id":         f.ID.String(),
			"name":       f.Name,
			"path":       f.LocalPath,
			"local_path": f.LocalPath,
			"remote_path": f.RemotePath,
			"is_active":  f.IsActive,
			"sync_mode":  f.SyncMode,
			"last_sync":  f.LastSync,
			"created_at": f.CreatedAt,
		})
	}
	c.JSON(http.StatusOK, gin.H{"folders": resp})
}

// CreateSyncFolder creates a new sync folder
func CreateSyncFolder(c *gin.Context) {
	userID, ok := getFileSyncUserID(c)
	if !ok {
		return
	}

	var req struct {
		Name string `json:"name"`
		Path string `json:"path"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "path is required"})
		return
	}
	if req.Name == "" {
		req.Name = req.Path
	}

	folder := models.SyncFolder{
		Name:      req.Name,
		LocalPath: req.Path,
		SyncMode:  "bidirectional",
	}

	created, err := fileSyncService.CreateSyncFolder(userID, folder)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"folder": map[string]interface{}{
		"id":         created.ID.String(),
		"name":       created.Name,
		"path":       created.LocalPath,
		"local_path": created.LocalPath,
		"remote_path": created.RemotePath,
		"is_active":  created.IsActive,
		"sync_mode":  created.SyncMode,
		"last_sync":  created.LastSync,
		"created_at": created.CreatedAt,
	}})
}

// GetFileSyncStatus returns sync status summary for the current user
func GetFileSyncStatus(c *gin.Context) {
	userID, ok := getFileSyncUserID(c)
	if !ok {
		return
	}

	status, err := fileSyncService.GetSyncStatus(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve sync status"})
		return
	}

	resp := map[string]interface{}{
		"total_folders": status["folder_count"],
		"total_files":   0,
		"total_size":    0,
	}
	if v, ok := status["device_count"]; ok {
		resp["device_count"] = v
	}
	if v, ok := status["pending_operations"]; ok {
		resp["pending_operations"] = v
	}
	if v, ok := status["conflicts"]; ok {
		resp["conflicts"] = v
	}

	c.JSON(http.StatusOK, gin.H{"status": resp})
}

// SyncFolderRoute triggers synchronization for a folder
func SyncFolderRoute(c *gin.Context) {
	userID, ok := getFileSyncUserID(c)
	if !ok {
		return
	}
	folderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid folder ID"})
		return
	}

	// The web client acts as a virtual device for sync purposes.
	response, err := fileSyncService.SyncFolder(userID, folderID, "web")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Sync completed",
		"changes":       response.Changes,
		"conflicts":     response.Conflicts,
		"next_sync_token": response.NextSyncToken,
	})
}
