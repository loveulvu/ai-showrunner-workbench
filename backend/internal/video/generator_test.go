package video

import (
	"context"
	"testing"
)

func TestMockVideoGeneratorImplementsInterface(t *testing.T) {
	var generator VideoGenerator = NewMockVideoGenerator()
	if generator == nil {
		t.Fatal("VideoGenerator = nil")
	}
}

func TestMockVideoGeneratorCreateAndGetTask(t *testing.T) {
	generator := NewMockVideoGenerator()
	taskID, err := generator.CreateTask(context.Background(), VideoPrompt{
		ShotID:           "shot_001",
		Model:            "mock-video-model",
		Prompt:           "slow cinematic camera move",
		DurationSeconds:  5,
		AspectRatio:      "16:9",
		ExpectedClipName: "shot_001.mp4",
	})
	if err != nil {
		t.Fatalf("CreateTask() error = %v", err)
	}
	if taskID == "" {
		t.Fatal("CreateTask() taskID is empty")
	}

	task, err := generator.store.Get(context.Background(), taskID)
	if err != nil {
		t.Fatalf("store.Get() error = %v", err)
	}
	if task.Status != StatusPending {
		t.Fatalf("created task status = %q, want %q", task.Status, StatusPending)
	}
	if task.CreatedAt.IsZero() || task.UpdatedAt.IsZero() {
		t.Fatal("created task timestamps must be populated")
	}

	result, err := generator.GetTask(context.Background(), taskID)
	if err != nil {
		t.Fatalf("GetTask() error = %v", err)
	}
	if result.Status != StatusSucceeded {
		t.Fatalf("result status = %q, want %q", result.Status, StatusSucceeded)
	}
	if result.VideoURL == "" {
		t.Fatal("result video_url is empty")
	}
	if result.ShotID != "shot_001" {
		t.Fatalf("result shot_id = %q", result.ShotID)
	}
}

func TestMockVideoGeneratorValidationAndMissingTask(t *testing.T) {
	generator := NewMockVideoGenerator()
	if _, err := generator.CreateTask(context.Background(), VideoPrompt{}); err == nil {
		t.Fatal("CreateTask() error = nil, want validation error")
	}
	if _, err := generator.GetTask(context.Background(), "missing"); err == nil {
		t.Fatal("GetTask() error = nil, want not found error")
	}
}
