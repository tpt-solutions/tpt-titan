package routes

import (
	"database/sql"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"tpt-titan/backend/config"
	"tpt-titan/backend/services"
)

var pluginSystem *services.PluginSystem

// InitPluginService initializes the plugin system (called from server setup)
func InitPluginService(db *sql.DB, cfg *config.Config) {
	pluginDir := pluginDirectory(cfg)
	pluginSystem = services.NewPluginSystem(db, pluginDir)
	if err := pluginSystem.InitializePluginSystem(); err != nil {
		// Log but don't fail startup — plugins are optional.
		// (logger not imported to avoid extra deps; ignore)
		_ = err
	}
}

func pluginDirectory(cfg *config.Config) string {
	dataDir := filepath.Dir(cfg.Database.Path)
	return filepath.Join(dataDir, "plugins")
}

func getPluginUserID(c *gin.Context) (uuid.UUID, bool) {
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

// GetPlugins returns the currently loaded plugins
func GetPlugins(c *gin.Context) {
	_, ok := getPluginUserID(c)
	if !ok {
		return
	}

	loaded := pluginSystem.GetLoadedPlugins()
	plugins := make([]map[string]interface{}, 0, len(loaded))
	for _, p := range loaded {
		plugins = append(plugins, map[string]interface{}{
			"id":          p.ID,
			"name":        p.Name,
			"version":     p.Version,
			"description": p.Description,
			"author":      p.Author,
			"enabled":     p.Enabled,
			"loaded_at":   p.LoadedAt,
			"hooks":       p.Hooks,
			"apis":        p.APIs,
		})
	}

	c.JSON(http.StatusOK, gin.H{"plugins": plugins})
}

// GetPluginStats returns plugin system statistics
func GetPluginStats(c *gin.Context) {
	_, ok := getPluginUserID(c)
	if !ok {
		return
	}
	c.JSON(http.StatusOK, gin.H{"stats": pluginSystem.GetPluginStats()})
}

// EnablePluginRoute enables a plugin
func EnablePluginRoute(c *gin.Context) {
	_, ok := getPluginUserID(c)
	if !ok {
		return
	}
	id := c.Param("id")
	if err := pluginSystem.EnablePlugin(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Plugin enabled", "id": id})
}

// DisablePluginRoute disables a plugin
func DisablePluginRoute(c *gin.Context) {
	_, ok := getPluginUserID(c)
	if !ok {
		return
	}
	id := c.Param("id")
	if err := pluginSystem.DisablePlugin(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Plugin disabled", "id": id})
}

// UnloadPluginRoute unloads and removes a plugin
func UnloadPluginRoute(c *gin.Context) {
	_, ok := getPluginUserID(c)
	if !ok {
		return
	}
	id := c.Param("id")
	if err := pluginSystem.UninstallPlugin(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Plugin uninstalled", "id": id})
}

// GetPluginSettingsRoute returns settings for a plugin
func GetPluginSettingsRoute(c *gin.Context) {
	_, ok := getPluginUserID(c)
	if !ok {
		return
	}
	id := c.Param("id")
	settings, err := pluginSystem.GetPluginSettings(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"settings": settings})
}

// UpdatePluginSettingsRoute updates settings for a plugin
func UpdatePluginSettingsRoute(c *gin.Context) {
	_, ok := getPluginUserID(c)
	if !ok {
		return
	}
	id := c.Param("id")
	var settings map[string]interface{}
	if err := c.ShouldBindJSON(&settings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := pluginSystem.UpdatePluginSettings(id, settings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Plugin settings updated", "id": id})
}
