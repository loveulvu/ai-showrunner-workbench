package video

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type MockVideoGenerator struct {
	mu      sync.RWMutex
	tasks   map[string]VideoTask
	prompts map[string]VideoPrompt
	nextID  atomic.Uint64
}

func NewMockVideoGenerator() *MockVideoGenerator {
	return &MockVideoGenerator{
		tasks:   map[string]VideoTask{},
		prompts: map[string]VideoPrompt{},
	}
}

func (g *MockVideoGenerator) CreateTask(ctx context.Context, prompt VideoPrompt) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	if strings.TrimSpace(prompt.ShotID) == "" {
		return "", fmt.Errorf("shot_id is required")
	}
	if strings.TrimSpace(prompt.Prompt) == "" {
		return "", fmt.Errorf("prompt is required")
	}

	taskID := fmt.Sprintf("mock-video-%06d", g.nextID.Add(1))
	now := time.Now().UTC()
	g.mu.Lock()
	g.tasks[taskID] = VideoTask{
		TaskID:    taskID,
		ShotID:    prompt.ShotID,
		Status:    StatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}
	g.prompts[taskID] = prompt
	g.mu.Unlock()
	return taskID, nil
}

func (g *MockVideoGenerator) GetTask(ctx context.Context, taskID string) (*VideoResult, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	g.mu.Lock()
	task, ok := g.tasks[taskID]
	if !ok {
		g.mu.Unlock()
		return nil, fmt.Errorf("video task %q not found", taskID)
	}
	prompt := g.prompts[taskID]
	task.Status = StatusSucceeded
	task.UpdatedAt = time.Now().UTC()
	g.tasks[taskID] = task
	g.mu.Unlock()

	clipName := strings.TrimSpace(prompt.ExpectedClipName)
	if clipName == "" {
		clipName = prompt.ShotID + ".mp4"
	}

	return &VideoResult{
		TaskID:   task.TaskID,
		ShotID:   task.ShotID,
		Status:   StatusSucceeded,
		VideoURL: "https://mock.video.local/clips/" + url.PathEscape(clipName),
	}, nil
}
