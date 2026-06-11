package main

import (
	"context"
	"log"

	"ai-showrunner-workbench/internal/ai"
	"ai-showrunner-workbench/internal/video"
)

func main() {
	ai.LoadEnv()

	config, err := video.ProviderConfigFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	video.LogProviderConfig(log.Default(), config)

	generator, err := video.NewGeneratorFromConfig(config, video.NewMemoryVideoTaskStore())
	if err != nil {
		log.Fatal(err)
	}

	taskID, err := generator.CreateTask(context.Background(), video.VideoPrompt{
		ShotID:          "video-check-shot",
		Prompt:          "A simple static animation test shot with minimal motion.",
		DurationSeconds: 5,
		AspectRatio:     "16:9",
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Video task created: task_id=%s", taskID)

	result, err := generator.GetTask(context.Background(), taskID)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Video task current status: %s", result.Status)
}
