package video

import "time"

type Status string

const (
	StatusPending   Status = "pending"
	StatusRunning   Status = "running"
	StatusSucceeded Status = "succeeded"
	StatusFailed    Status = "failed"
)

type VideoPrompt struct {
	ShotID           string `json:"shot_id"`
	Model            string `json:"model"`
	Prompt           string `json:"prompt"`
	NegativePrompt   string `json:"negative_prompt"`
	DurationSeconds  int    `json:"duration_seconds"`
	AspectRatio      string `json:"aspect_ratio"`
	Subtitle         string `json:"subtitle"`
	ExpectedClipName string `json:"expected_clip_name"`
}

type VideoTask struct {
	TaskID    string    `json:"task_id"`
	ShotID    string    `json:"shot_id"`
	Status    Status    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type VideoResult struct {
	TaskID       string `json:"task_id"`
	ShotID       string `json:"shot_id"`
	Status       Status `json:"status"`
	VideoURL     string `json:"video_url"`
	ErrorMessage string `json:"error_message"`
}
