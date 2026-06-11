package video

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"sync/atomic"
	"time"
)

type MockVideoGenerator struct {
	config ProviderConfig
	store  VideoTaskStore
	nextID atomic.Uint64
}

func NewMockVideoGenerator() *MockVideoGenerator {
	config := ProviderConfig{
		Provider:            defaultVideoProvider,
		Model:               defaultVideoModel,
		TimeoutSeconds:      defaultVideoTimeoutSeconds,
		PollIntervalSeconds: defaultVideoPollIntervalSeconds,
	}
	return NewMockVideoGeneratorWithStore(config, NewMemoryVideoTaskStore())
}

func NewMockVideoGeneratorWithStore(config ProviderConfig, store VideoTaskStore) *MockVideoGenerator {
	return &MockVideoGenerator{
		config: config,
		store:  store,
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
	model := strings.TrimSpace(prompt.Model)
	if model == "" {
		model = g.config.Model
	}

	task := &VideoTask{
		TaskID:           taskID,
		ShotID:           prompt.ShotID,
		Provider:         g.config.Provider,
		Model:            model,
		Prompt:           prompt.Prompt,
		NegativePrompt:   prompt.NegativePrompt,
		DurationSeconds:  prompt.DurationSeconds,
		AspectRatio:      prompt.AspectRatio,
		Subtitle:         prompt.Subtitle,
		ExpectedClipName: prompt.ExpectedClipName,
		Status:           StatusPending,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	if err := g.store.Save(ctx, task); err != nil {
		return "", fmt.Errorf("save video task: %w", err)
	}
	return taskID, nil
}

func (g *MockVideoGenerator) GetTask(ctx context.Context, taskID string) (*VideoResult, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	task, err := g.store.Get(ctx, taskID)
	if err != nil {
		return nil, err
	}

	task.Status = StatusSucceeded
	task.UpdatedAt = time.Now().UTC()
	clipName := strings.TrimSpace(task.ExpectedClipName)
	if clipName == "" {
		clipName = task.ShotID + ".mp4"
	}
	task.VideoURL = "https://mock.video.local/clips/" + url.PathEscape(clipName)
	if err := g.store.Update(ctx, task); err != nil {
		return nil, fmt.Errorf("update video task: %w", err)
	}

	return &VideoResult{
		TaskID:   task.TaskID,
		ShotID:   task.ShotID,
		Status:   task.Status,
		VideoURL: task.VideoURL,
	}, nil
}
