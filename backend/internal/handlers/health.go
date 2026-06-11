package handlers

import (
	"net/http"

	"novel-to-screenplay-ai/internal/ai"

	"github.com/gin-gonic/gin"
)

func Health(c *gin.Context) {
	status := ai.RuntimeStatusFromEnv()

	c.JSON(http.StatusOK, gin.H{
		"status":                 "ok",
		"service":                "novel-to-screenplay-ai",
		"ai_provider":            status.AIProvider,
		"ai_model":               status.AIModel,
		"ai_base_url_configured": status.AIBaseURLConfigured,
		"ai_api_key_configured":  status.AIAPIKeyConfigured,
		"ai_timeout_seconds":     status.AITimeoutSeconds,
	})
}
