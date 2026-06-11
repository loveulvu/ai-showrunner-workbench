package main

import (
	"log"
	"net/http"
	"net/url"
	"novel-to-screenplay-ai/internal/ai"
	"novel-to-screenplay-ai/internal/handlers"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	loadedEnvFiles := ai.LoadEnv()
	aiStatus := ai.RuntimeStatusFromEnv()

	r := gin.Default()
	r.Use(withCORS())
	r.GET("/api/health", handlers.Health)
	r.POST("/api/generate", handlers.Generate)
	addr := ":8080"

	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}

	if len(loadedEnvFiles) == 0 {
		log.Printf("No environment file loaded (checked .env and ../.env)")
	} else {
		for _, path := range loadedEnvFiles {
			log.Printf("Loaded environment file: %s", path)
		}
	}
	log.Printf("AI_PROVIDER=%s", aiStatus.AIProvider)
	log.Printf("AI_MODEL=%s", aiStatus.AIModel)
	log.Printf("AI_BASE_URL=%s", safeURLForLog(aiStatus.AIBaseURL))
	log.Printf("AI_API_KEY set: %t", aiStatus.AIAPIKeyConfigured)
	log.Printf("HTTP_PROXY set: %t", aiStatus.HTTPProxyConfigured)
	log.Printf("HTTPS_PROXY set: %t", aiStatus.HTTPSProxyConfigured)
	log.Printf("AI_TIMEOUT_SECONDS=%d", aiStatus.AITimeoutSeconds)
	log.Printf("novel-to-screenplay-ai backend listening on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}

func safeURLForLog(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}

	parsed, err := url.Parse(value)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return "<invalid URL configured>"
	}

	if parsed.User != nil {
		parsed.User = url.User("REDACTED")
	}
	parsed.RawQuery = ""
	parsed.Fragment = ""
	return parsed.String()
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
