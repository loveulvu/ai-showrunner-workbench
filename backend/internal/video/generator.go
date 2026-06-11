package video

import "context"

type VideoGenerator interface {
	CreateTask(ctx context.Context, prompt VideoPrompt) (string, error)
	GetTask(ctx context.Context, taskID string) (*VideoResult, error)
}
