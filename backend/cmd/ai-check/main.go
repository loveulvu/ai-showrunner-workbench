package main

import (
	"context"
	"log"

	"ai-showrunner-workbench/internal/ai"
)

func main() {
	ai.LoadEnv()
	ai.LogRuntimeConfiguration(log.Default())

	if err := ai.CheckConnectivityFromEnv(context.Background()); err != nil {
		log.Fatalf("Qwen connectivity check failed: %s", ai.RedactedDiagnostic(err))
	}
	log.Printf("Qwen connectivity check passed")
}
