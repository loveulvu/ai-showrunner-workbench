package handlers

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"ai-showrunner-workbench/internal/video"

	"github.com/gin-gonic/gin"
)

var videoGenerator video.VideoGenerator = video.NewMockVideoGenerator()

func ConfigureVideoGeneratorFromEnv() (video.ProviderConfig, error) {
	config, err := video.ProviderConfigFromEnv()
	if err != nil {
		return config, err
	}
	generator, err := video.NewGeneratorFromConfig(config, video.NewMemoryVideoTaskStore())
	if err != nil {
		return config, err
	}
	videoGenerator = generator
	return config, nil
}

func CreateVideoTask(c *gin.Context) {
	var prompt video.VideoPrompt
	if err := c.ShouldBindJSON(&prompt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid json body", "error": "invalid json body"})
		return
	}
	if strings.TrimSpace(prompt.ShotID) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "shot_id is required", "error": "shot_id is required"})
		return
	}
	if strings.TrimSpace(prompt.Prompt) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "prompt is required", "error": "prompt is required"})
		return
	}

	taskID, err := videoGenerator.CreateTask(c.Request.Context(), prompt)
	if err != nil {
		if video.IsRequestError(err) {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "error": err.Error()})
			return
		}
		log.Printf("Create video task failed: %s", err.Error())
		status := http.StatusInternalServerError
		if video.IsUpstreamError(err) {
			status = http.StatusBadGateway
		}
		c.JSON(status, gin.H{"message": err.Error(), "error": err.Error()})
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
		status := http.StatusInternalServerError
		if errors.Is(err, video.ErrTaskNotFound) {
			status = http.StatusNotFound
		} else if video.IsUpstreamError(err) {
			status = http.StatusBadGateway
		}
		c.JSON(status, gin.H{"message": err.Error(), "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
