package handlers

import (
	"net/http"

	"ai-showrunner-workbench/internal/ai"
	"ai-showrunner-workbench/internal/showrunner"

	"github.com/gin-gonic/gin"
)

func GenerateShowrunner(c *gin.Context) {
	var input showrunner.GenerateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json body"})
		return
	}

	client, err := ai.NewClientFromEnv()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result, err := showrunner.NewService(client).Generate(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
