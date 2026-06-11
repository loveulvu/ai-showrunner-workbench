package video

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

var ErrTaskNotFound = errors.New("video task not found")

type VideoTaskStore interface {
	Save(ctx context.Context, task *VideoTask) error
	Get(ctx context.Context, taskID string) (*VideoTask, error)
	Update(ctx context.Context, task *VideoTask) error
}

type MemoryVideoTaskStore struct {
	mu    sync.RWMutex
	tasks map[string]VideoTask
}

func NewMemoryVideoTaskStore() *MemoryVideoTaskStore {
	return &MemoryVideoTaskStore{tasks: map[string]VideoTask{}}
}

func (s *MemoryVideoTaskStore) Save(ctx context.Context, task *VideoTask) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if task == nil || task.TaskID == "" {
		return fmt.Errorf("task_id is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.tasks[task.TaskID]; exists {
		return fmt.Errorf("video task %q already exists", task.TaskID)
	}
	s.tasks[task.TaskID] = *task
	return nil
}

func (s *MemoryVideoTaskStore) Get(ctx context.Context, taskID string) (*VideoTask, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	s.mu.RLock()
	task, ok := s.tasks[taskID]
	s.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrTaskNotFound, taskID)
	}
	return &task, nil
}

func (s *MemoryVideoTaskStore) Update(ctx context.Context, task *VideoTask) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if task == nil || task.TaskID == "" {
		return fmt.Errorf("task_id is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.tasks[task.TaskID]; !exists {
		return fmt.Errorf("%w: %s", ErrTaskNotFound, task.TaskID)
	}
	s.tasks[task.TaskID] = *task
	return nil
}
