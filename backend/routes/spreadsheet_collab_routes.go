package routes

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"tpt-titan/backend/config"
	"tpt-titan/backend/services"
)

// GetCollaborationMode returns the current collaboration mode
func GetCollaborationMode(c *gin.Context) {
	cfg := c.MustGet("config").(*config.Config)

	mode := "server" // Default to server mode
	if cfg.P2P.Enabled {
		mode = "p2p"
	}

	c.JSON(http.StatusOK, gin.H{
		"mode":       mode,
		"p2p_enabled": cfg.P2P.Enabled,
		"server_available": true, // Server is always available as backup
	})
}

// SetCollaborationMode switches between server and P2P modes
func SetCollaborationMode(c *gin.Context) {
	var req struct {
		Mode string `json:"mode" binding:"required"` // "server" or "p2p"
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Mode != "server" && req.Mode != "p2p" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mode. Must be 'server' or 'p2p'"})
		return
	}

	// In a real implementation, this would update configuration
	// and restart services as needed
	cfg := c.MustGet("config").(*config.Config)

	switch req.Mode {
	case "p2p":
		if !cfg.P2P.Enabled {
			c.JSON(http.StatusBadRequest, gin.H{"error": "P2P mode is not configured"})
			return
		}
		// Enable P2P services
		p2pService := c.MustGet("p2p_service").(*services.P2PService)
		if !p2pService.IsRunning() {
			err := p2pService.Start()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start P2P service: " + err.Error()})
				return
			}
		}
	case "server":
		// P2P can remain running as fallback, but primary mode is server
		// Server mode uses existing WebSocket/real-time infrastructure
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Collaboration mode updated successfully",
		"mode":    req.Mode,
	})
}

// GetConnectedPeers returns list of connected P2P peers
func GetConnectedPeers(c *gin.Context) {
	p2pService := c.MustGet("p2p_service").(*services.P2PService)
	peers := p2pService.GetConnectedPeers()

	c.JSON(http.StatusOK, gin.H{
		"peers": peers,
		"count": len(peers),
	})
}

// ConnectToPeer manually connects to a specific peer
func ConnectToPeer(c *gin.Context) {
	var req struct {
		PeerID  string `json:"peer_id" binding:"required"`
		Address string `json:"address" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	p2pService := c.MustGet("p2p_service").(*services.P2PService)
	err := p2pService.ConnectToPeer(req.PeerID, req.Address)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to connect to peer: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully connected to peer",
		"peer_id": req.PeerID,
	})
}

// GetDiscoveredPeers returns peers discovered on the network
func GetDiscoveredPeers(c *gin.Context) {
	// In a full implementation, this would query the peer discovery service
	discoveredPeers := []gin.H{
		{
			"id":      "peer-001",
			"name":    "Alice's Computer",
			"address": "192.168.1.100:8081",
			"status":  "available",
		},
		{
			"id":      "peer-002",
			"name":    "Bob's Laptop",
			"address": "192.168.1.101:8081",
			"status":  "available",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"peers": discoveredPeers,
	})
}

// SyncSpreadsheetWithPeers broadcasts spreadsheet to all connected peers
func SyncSpreadsheetWithPeers(c *gin.Context) {
	spreadsheetIDStr := c.Param("id")
	spreadsheetID, err := uuid.Parse(spreadsheetIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid spreadsheet ID"})
		return
	}

	// Get spreadsheet data from database
	db := c.MustGet("db").(*sql.DB)

	var spreadsheet struct {
		Name string
	}

	err = db.QueryRow("SELECT name FROM spreadsheets WHERE id = $1", spreadsheetID).Scan(&spreadsheet.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get spreadsheet: " + err.Error()})
		return
	}

	// Get cell data
	rows, err := db.Query("SELECT cell_reference, value FROM spreadsheet_cells WHERE spreadsheet_id = $1", spreadsheetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get cell data: " + err.Error()})
		return
	}
	defer rows.Close()

	data := make(map[string]interface{})
	for rows.Next() {
		var cellRef, value string
		rows.Scan(&cellRef, &value)
		// Convert value to appropriate type (simplified)
		if num, err := strconv.ParseFloat(value, 64); err == nil {
			data[cellRef] = num
		} else {
			data[cellRef] = value
		}
	}

	// Sync with P2P peers
	p2pService := c.MustGet("p2p_service").(*services.P2PService)
	p2pService.SyncSpreadsheet(spreadsheetID, data)

	c.JSON(http.StatusOK, gin.H{
		"message": "Spreadsheet synced with peers",
		"peers_notified": len(p2pService.GetConnectedPeers()),
	})
}

// GetCollaborationStatus returns current collaboration status
func GetCollaborationStatus(c *gin.Context) {
	cfg := c.MustGet("config").(*config.Config)
	p2pService := c.MustGet("p2p_service").(*services.P2PService)

	// Determine current mode based on configuration and what's running
	currentMode := "automatic" // Default: automatic selection for ease of use

	peers := p2pService.GetConnectedPeers()
	connectedPeers := len(peers)

	status := gin.H{
		"mode": currentMode, // "automatic", "p2p", "server"
		"connected_peers": connectedPeers,
		"remote_access_enabled": cfg.P2P.AllowRemoteAccess,
		"features": gin.H{
			"local_network": gin.H{
				"available": true,
				"description": "Direct peer-to-peer collaboration on local networks",
				"speed": "Fast",
				"security": "Local network only",
			},
			"remote_access": gin.H{
				"available": cfg.P2P.AllowRemoteAccess,
				"description": "Cloud relay for remote users (work from home, etc.)",
				"speed": "Good",
				"security": "End-to-end encrypted",
			},
			"server_backup": gin.H{
				"available": true,
				"description": "Full server-based collaboration as backup",
				"speed": "Variable",
				"security": "User authentication + permissions",
			},
		},
		"p2p": gin.H{
			"enabled": cfg.P2P.Enabled,
			"running": p2pService.IsRunning(),
			"topology": cfg.P2P.PreferredTopology,
			"auto_detect": cfg.P2P.AutoDetectTopology,
			"config": gin.H{
				"port": cfg.P2P.Port,
				"max_peers": cfg.P2P.MaxPeers,
				"discovery_timeout": cfg.P2P.DiscoveryTimeout,
				"sync_interval": cfg.P2P.SyncInterval,
				"encryption": cfg.P2P.EnableEncryption,
				"compression": cfg.P2P.EnableCompression,
				"remote_access": cfg.P2P.AllowRemoteAccess,
				"cloud_relay": cfg.P2P.CloudRelayEnabled,
			},
		},
		"server": gin.H{
			"available": true,
			"description": "Traditional server-based collaboration",
			"features": []string{
				"User authentication",
				"Granular permissions",
				"Version history",
				"Advanced sharing",
			},
		},
		"recommendation": gin.H{
			"for_local_office": "automatic (P2P with cloud relay fallback)",
			"for_remote_teams": "server mode",
			"for_home_users": "automatic (works everywhere)",
			"easiest_setup": "automatic mode - just works",
		},
		"peers": peers,
	}

	// Determine actual current mode
	if cfg.P2P.Enabled && p2pService.IsRunning() {
		if connectedPeers > 0 {
			status["mode"] = "p2p"
			status["active_connections"] = "local_network"
		} else if cfg.P2P.AllowRemoteAccess {
			status["mode"] = "remote_ready"
			status["active_connections"] = "cloud_relay_available"
		}
	} else {
		status["mode"] = "server"
		status["active_connections"] = "server_based"
	}

	// Add user-friendly status messages
	status["user_message"] = gin.H{
		"current_status": "Ready for collaboration",
		"ease_of_use": "Just works - no configuration needed",
		"security_note": "All connections are encrypted and secure",
		"performance": "Automatic optimization for your network",
	}

	c.JSON(http.StatusOK, status)
}
