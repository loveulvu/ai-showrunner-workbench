package handlers

import (
	"net/http"

	"ai-showrunner-workbench/internal/editor"

	"github.com/gin-gonic/gin"
)

var renderEditorPlan = editor.Render

func RenderEditorDemo(c *gin.Context) {
	var plan editor.EditingPlan
	if err := c.ShouldBindJSON(&plan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json body"})
		return
	}
	if len(plan.Clips) == 0 || len(plan.Clips) > 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "clips must contain between 1 and 3 items"})
		return
	}

	result, err := renderEditorPlan(c.Request.Context(), plan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
