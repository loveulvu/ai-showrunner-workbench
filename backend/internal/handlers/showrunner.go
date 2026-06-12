package handlers

import (
	"log"
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
		log.Printf("Showrunner failed stage=%s message=%s", showrunner.StageService, ai.RedactedDiagnostic(err))
		c.JSON(http.StatusInternalServerError, gin.H{"stage": showrunner.StageService, "message": err.Error(), "error": err.Error()})
		return
	}

	result, err := showrunner.NewService(client).Generate(c.Request.Context(), input)
	if err != nil {
		stage, message := showrunner.ErrorDetails(err)
		log.Printf("Showrunner failed stage=%s message=%s", stage, ai.RedactedDiagnostic(err))
		c.JSON(http.StatusInternalServerError, gin.H{"stage": stage, "message": message, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
