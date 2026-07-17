package routes

import (
	"net/http"

	"tpt-titan/backend/config"
	"tpt-titan/backend/models"
	"tpt-titan/backend/services"
	"tpt-titan/backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ListMCPServers returns the user's configured MCP servers.
func ListMCPServers(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var servers []models.MCPServer
	if err := config.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&servers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list MCP servers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"servers": servers})
}

// CreateMCPServer registers a new MCP server and registers its tools as workflow
// connectors. The auth token (when auth_type is "bearer") is encrypted at rest.
func CreateMCPServer(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var payload struct {
		Name      string `json:"name"`
		URL       string `json:"url"`
		Transport string `json:"transport"`
		AuthType  string `json:"auth_type"`
		Token     string `json:"token"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if payload.Name == "" || payload.URL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name and url are required"})
		return
	}
	if payload.Transport == "" {
		payload.Transport = "http"
	}
	if payload.AuthType == "" {
		payload.AuthType = "none"
	}

	server := models.MCPServer{
		UserID:    userID.(uuid.UUID),
		Name:      payload.Name,
		URL:       payload.URL,
		Transport: payload.Transport,
		AuthType:  payload.AuthType,
		IsActive:  true,
	}
	if payload.AuthType == "bearer" && payload.Token != "" {
		enc, err := utils.EncryptPassword(payload.Token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt token"})
			return
		}
		server.Token = string(enc)
	}

	if err := config.DB.Create(&server).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create MCP server"})
		return
	}

	if workflowService != nil {
		uid := userID.(uuid.UUID)
		if err := workflowService.RegisterMCPConnectors(&uid); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"server":  server,
				"warning": "MCP server saved, but tools could not be registered: " + err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"server": server})
}

// DeleteMCPServer removes a configured MCP server. (Registered connectors remain
// in the running process until restart; this is acceptable for low-churn config.)
func DeleteMCPServer(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid server ID"})
		return
	}

	if err := config.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.MCPServer{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete MCP server"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "MCP server deleted"})
}

// ListMCPTools lists the tools exposed by a configured MCP server (after a live
// handshake + tools/list). This lets the frontend present available tools.
func ListMCPTools(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid server ID"})
		return
	}

	var server models.MCPServer
	if err := config.DB.Where("id = ? AND user_id = ?", id, userID).First(&server).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "MCP server not found"})
		return
	}

	token := server.Token
	if server.AuthType == "bearer" && token != "" {
		if dec, derr := utils.DecryptPassword([]byte(token)); derr == nil {
			token = dec
		}
	}

	client := services.NewMCPClient(models.MCPServer{URL: server.URL, AuthType: server.AuthType, Token: token})
	if err := client.Initialize(); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "MCP handshake failed: " + err.Error()})
		return
	}
	tools, err := client.ListTools()
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to list MCP tools: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tools": tools})
}

// GetMCPConnectors returns the MCP tool connectors currently registered in the
// workflow engine, so the builder UI can offer them as selectable connectors.
func GetMCPConnectors(c *gin.Context) {
	if _, ok := c.Get("user_id"); !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	if workflowService == nil {
		c.JSON(http.StatusOK, gin.H{"connectors": []interface{}{}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"connectors": workflowService.GetMCPConnectorDescriptors()})
}

// TestMCPServer performs a handshake + tools/list against the server URL supplied
// in the request body (no persistence) so the user can validate a config.
func TestMCPServer(c *gin.Context) {
	var payload struct {
		URL       string `json:"url"`
		Transport string `json:"transport"`
		AuthType  string `json:"auth_type"`
		Token     string `json:"token"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if payload.URL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url is required"})
		return
	}

	client := services.NewMCPClient(models.MCPServer{
		URL:      payload.URL,
		AuthType: payload.AuthType,
		Token:    payload.Token,
	})
	if err := client.Initialize(); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"ok": false, "error": err.Error()})
		return
	}
	tools, err := client.ListTools()
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"ok": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true, "tool_count": len(tools), "tools": tools})
}
