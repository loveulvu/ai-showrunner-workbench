package main

import (
	"ai-showrunner-workbench/internal/ai"
	"ai-showrunner-workbench/internal/handlers"
	"ai-showrunner-workbench/internal/video"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	ai.LoadEnv()

	r := gin.Default()
	r.Use(withCORS())
	r.GET("/api/health", handlers.Health)
	r.POST("/api/generate", handlers.Generate)
	r.POST("/api/showrunner/generate", handlers.GenerateShowrunner)
	r.POST("/api/video/tasks", handlers.CreateVideoTask)
	r.GET("/api/video/tasks/:task_id", handlers.GetVideoTask)
	r.POST("/api/editor/render", handlers.RenderEditorDemo)
	addr := ":8080"

	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}

	ai.LogRuntimeConfiguration(log.Default())
	videoConfig, err := handlers.ConfigureVideoGeneratorFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	video.LogProviderConfig(log.Default(), videoConfig)
	log.Printf("ai-showrunner-workbench backend listening on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}

func withCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}

}
