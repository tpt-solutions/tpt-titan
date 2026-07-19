package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"tpt-titan/backend/config"
	"tpt-titan/backend/models"
	"tpt-titan/backend/services"
)

// GetWebhookDeliveryLogs returns recent inbound/outbound webhook call records
// for the monitoring / admin webhook dashboard.
func GetWebhookDeliveryLogs(c *gin.Context) {
	logs, err := services.GetWebhookDeliveryLogs(0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load delivery log"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"logs": logs})
}

// GetOutboundDomainAllowlist returns the admin-configured outbound domain
// allowlist (empty means "allow all public destinations").
func GetOutboundDomainAllowlist(c *gin.Context) {
	db := config.GetDatabase()
	var setting models.SystemSetting
	if err := db.Where("key = ?", "outbound_domains").First(&setting).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"domains": []string{}})
		return
	}
	var domains []string
	if err := json.Unmarshal([]byte(setting.Value), &domains); err != nil {
		c.JSON(http.StatusOK, gin.H{"domains": []string{}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"domains": domains})
}

// UpdateOutboundDomainAllowlist persists the admin-configured outbound domain
// allowlist.
func UpdateOutboundDomainAllowlist(c *gin.Context) {
	var payload struct {
		Domains []string `json:"domains"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.GetDatabase()
	encoded, err := json.Marshal(payload.Domains)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encode allowlist"})
		return
	}

	var setting models.SystemSetting
	if err := db.Where("key = ?", "outbound_domains").First(&setting).Error; err != nil {
		setting = models.SystemSetting{Key: "outbound_domains", Value: string(encoded)}
		if err := db.Create(&setting).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create allowlist"})
			return
		}
	} else {
		setting.Value = string(encoded)
		if err := db.Save(&setting).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update allowlist"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Outbound domain allowlist updated", "domains": payload.Domains})
}
