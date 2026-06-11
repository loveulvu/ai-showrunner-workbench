package video

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestMemoryVideoTaskStoreSaveGetUpdate(t *testing.T) {
	store := NewMemoryVideoTaskStore()
	now := time.Now().UTC()
	task := &VideoTask{
		TaskID:    "task_001",
		ShotID:    "shot_001",
		Provider:  "mock",
		Model:     "mock-model",
		Prompt:    "camera move",
		Status:    StatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := store.Save(context.Background(), task); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	got, err := store.Get(context.Background(), task.TaskID)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if got.Status != StatusPending || got.Prompt != "camera move" {
		t.Fatalf("Get() task = %#v", got)
	}

	got.Status = StatusRunning
	got.UpdatedAt = now.Add(time.Second)
	if err := store.Update(context.Background(), got); err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	updated, err := store.Get(context.Background(), task.TaskID)
	if err != nil {
		t.Fatalf("Get() after update error = %v", err)
	}
	if updated.Status != StatusRunning {
		t.Fatalf("updated status = %q, want running", updated.Status)
	}
}

func TestMemoryVideoTaskStoreReturnsCopiesAndNotFound(t *testing.T) {
	store := NewMemoryVideoTaskStore()
	task := &VideoTask{TaskID: "task_001", Status: StatusPending}
	if err := store.Save(context.Background(), task); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	task.Status = StatusFailed
	stored, err := store.Get(context.Background(), task.TaskID)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if stored.Status != StatusPending {
		t.Fatalf("stored status = %q, want pending copy", stored.Status)
	}

	_, err = store.Get(context.Background(), "missing")
	if !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("Get() error = %v, want ErrTaskNotFound", err)
	}
}
