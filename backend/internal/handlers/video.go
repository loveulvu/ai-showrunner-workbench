package handlers

import (
	"net/http"
	"strings"

	"ai-showrunner-workbench/internal/video"

	"github.com/gin-gonic/gin"
)

var videoGenerator video.VideoGenerator = video.NewMockVideoGenerator()

func CreateVideoTask(c *gin.Context) {
	var prompt video.VideoPrompt
	if err := c.ShouldBindJSON(&prompt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json body"})
		return
	}

	taskID, err := videoGenerator.CreateTask(c.Request.Context(), prompt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"task_id": taskID})
}

func GetVideoTask(c *gin.Context) {
	taskID := strings.TrimSpace(c.Param("task_id"))
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "task_id is required"})
		return
	}

	result, err := videoGenerator.GetTask(c.Request.Context(), taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
