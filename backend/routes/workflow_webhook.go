package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ReceiveWebhook is a public (unauthenticated) endpoint that external systems
// call to trigger a workflow. The URL's :token is the only thing identifying
// which workflow to run — it must be treated as a secret, so on any lookup
// failure this returns a generic 404 rather than revealing why the token
// didn't match.
func ReceiveWebhook(c *gin.Context) {
	token := c.Param("token")

	workflow, err := workflowService.FindWebhookTriggeredWorkflow(token)
	if err != nil || workflow == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	var payload map[string]interface{}
	// Tolerate an empty or absent body — not every webhook sender includes one.
	_ = c.ShouldBindJSON(&payload)
	if payload == nil {
		payload = map[string]interface{}{}
	}

	triggerData := make(map[string]interface{}, len(payload)+1)
	for k, v := range payload {
		triggerData[k] = v
	}
	triggerData["user_id"] = workflow.UserID.String()

	if _, err := workflowService.ExecuteWorkflow(workflow.ID, triggerData, false); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to start workflow"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"status": "accepted"})
}
