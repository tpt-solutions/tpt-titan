package routes

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"tpt-titan/backend/models"
	"tpt-titan/backend/services"
)

const webhookInboundBodyLimit = 4096

// ReceiveWebhook is a public (unauthenticated) endpoint that external systems
// call to trigger a workflow. The URL's :token is the only thing identifying
// which workflow to run — it must be treated as a secret, so on any lookup
// failure this returns a generic 404 rather than revealing why the token
// didn't match.
func ReceiveWebhook(c *gin.Context) {
	token := c.Param("token")

	body, _ := io.ReadAll(c.Request.Body)
	var payload map[string]interface{}
	// Tolerate an empty or absent body — not every webhook sender includes one.
	if len(body) > 0 {
		_ = json.Unmarshal(body, &payload)
	}
	if payload == nil {
		payload = map[string]interface{}{}
	}
	redactedBody := string(body)
	if len(redactedBody) > webhookInboundBodyLimit {
		redactedBody = redactedBody[:webhookInboundBodyLimit] + "…[truncated]"
	}

	logInbound := func(status int, errMsg string) {
		// Record the inbound call for the dashboard. The token itself is never
		// persisted — only the matched workflow id (when known).
		services.RecordDeliveryLog(&models.WebhookDeliveryLog{
			Direction:   "inbound",
			Connector:   "webhook.receive",
			URL:         c.Request.URL.Path,
			Host:        c.Request.Host,
			Method:      c.Request.Method,
			RequestBody: redactedBody,
			StatusCode:  status,
			Error:       errMsg,
		})
	}

	workflow, err := workflowService.FindWebhookTriggeredWorkflow(token)
	if err != nil || workflow == nil {
		logInbound(http.StatusNotFound, "no matching workflow for token")
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	triggerData := make(map[string]interface{}, len(payload)+1)
	for k, v := range payload {
		triggerData[k] = v
	}
	triggerData["user_id"] = workflow.UserID.String()

	if _, err := workflowService.ExecuteWorkflow(workflow.ID, triggerData, false); err != nil {
		logInbound(http.StatusInternalServerError, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to start workflow"})
		return
	}

	logInbound(http.StatusAccepted, "")
	c.JSON(http.StatusAccepted, gin.H{"status": "accepted"})
}
