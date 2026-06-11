package video

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBuildWanCreateRequest(t *testing.T) {
	request := buildWanCreateRequest(
		ProviderConfig{Model: "wan2.6-t2v"},
		VideoPrompt{
			Model:           "request-model-must-not-override-config",
			Prompt:          "a short animated shot",
			NegativePrompt:  "blur",
			DurationSeconds: 0,
			AspectRatio:     "unsupported",
		},
	)

	if request.Model != "wan2.6-t2v" {
		t.Fatalf("Model = %q", request.Model)
	}
	if request.Parameters.Size != "1280*720" {
		t.Fatalf("Size = %q, want 1280*720", request.Parameters.Size)
	}
	if request.Parameters.Duration != 5 {
		t.Fatalf("Duration = %d, want 5", request.Parameters.Duration)
	}
	if !request.Parameters.PromptExtend || request.Parameters.Watermark {
		t.Fatalf("Parameters = %#v", request.Parameters)
	}
}

func TestWanErrorsRedactAPIKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"message":"test-key is invalid"}`))
	}))
	defer server.Close()

	generator := NewWanVideoGenerator(
		ProviderConfig{Provider: "wan", Model: "wan2.6-t2v", BaseURL: server.URL},
		NewMemoryVideoTaskStore(),
		"test-key",
		server.Client(),
	)
	_, err := generator.CreateTask(context.Background(), VideoPrompt{ShotID: "shot_001", Prompt: "test"})
	if err == nil {
		t.Fatal("CreateTask() error = nil, want unauthorized error")
	}
	if strings.Contains(err.Error(), "test-key") {
		t.Fatalf("error contains API key: %q", err)
	}
}

func TestWanCreateTaskRequestAndResponse(t *testing.T) {
	var received wanCreateRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/services/aigc/video-generation/video-synthesis" {
			t.Errorf("path = %q", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Errorf("Authorization header is incorrect")
		}
		if r.Header.Get("X-DashScope-Async") != "enable" {
			t.Errorf("X-DashScope-Async = %q", r.Header.Get("X-DashScope-Async"))
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		_, _ = w.Write([]byte(`{"output":{"task_id":"wan-task-001","task_status":"PENDING"}}`))
	}))
	defer server.Close()

	store := NewMemoryVideoTaskStore()
	generator := NewWanVideoGenerator(
		ProviderConfig{Provider: "wan", Model: "wan2.6-t2v", BaseURL: server.URL + "/api/v1"},
		store,
		"test-key",
		server.Client(),
	)

	taskID, err := generator.CreateTask(context.Background(), VideoPrompt{
		ShotID:          "shot_001",
		Prompt:          "a short animated shot",
		NegativePrompt:  "blur",
		DurationSeconds: 5,
		AspectRatio:     "16:9",
	})
	if err != nil {
		t.Fatalf("CreateTask() error = %v", err)
	}
	if taskID != "wan-task-001" {
		t.Fatalf("taskID = %q", taskID)
	}
	if received.Parameters.Size != "1280*720" || received.Parameters.Duration != 5 {
		t.Fatalf("received parameters = %#v", received.Parameters)
	}
	task, err := store.Get(context.Background(), taskID)
	if err != nil {
		t.Fatalf("store.Get() error = %v", err)
	}
	if task.Status != StatusPending || task.Provider != "wan" {
		t.Fatalf("stored task = %#v", task)
	}
}

func TestWanGetTaskQueryResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/tasks/wan-task-001" {
			t.Errorf("path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"output":{"task_id":"wan-task-001","task_status":"SUCCEEDED","video_url":"https://example.com/video.mp4"}}`))
	}))
	defer server.Close()

	store := NewMemoryVideoTaskStore()
	if err := store.Save(context.Background(), &VideoTask{TaskID: "wan-task-001", ShotID: "shot_001", Status: StatusPending}); err != nil {
		t.Fatalf("store.Save() error = %v", err)
	}
	generator := NewWanVideoGenerator(
		ProviderConfig{Provider: "wan", BaseURL: server.URL + "/api/v1"},
		store,
		"test-key",
		server.Client(),
	)

	result, err := generator.GetTask(context.Background(), "wan-task-001")
	if err != nil {
		t.Fatalf("GetTask() error = %v", err)
	}
	if result.Status != StatusSucceeded || result.VideoURL != "https://example.com/video.mp4" {
		t.Fatalf("result = %#v", result)
	}
}

func TestWanGetTaskCodeAndMessageBecomeFailedResult(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"output":{"task_id":"wan-task-001","task_status":"FAILED"},"code":"InvalidParameter","message":"bad request"}`))
	}))
	defer server.Close()

	store := NewMemoryVideoTaskStore()
	_ = store.Save(context.Background(), &VideoTask{TaskID: "wan-task-001", ShotID: "shot_001", Status: StatusPending})
	generator := NewWanVideoGenerator(ProviderConfig{BaseURL: server.URL}, store, "test-key", server.Client())

	result, err := generator.GetTask(context.Background(), "wan-task-001")
	if err != nil {
		t.Fatalf("GetTask() error = %v", err)
	}
	if result.Status != StatusFailed || !strings.Contains(result.ErrorMessage, "InvalidParameter") {
		t.Fatalf("result = %#v", result)
	}
}

func TestMapWanStatus(t *testing.T) {
	tests := map[string]Status{
		"PENDING":   StatusPending,
		"RUNNING":   StatusRunning,
		"SUCCEEDED": StatusSucceeded,
		"FAILED":    StatusFailed,
		"UNKNOWN":   StatusFailed,
	}
	for input, expected := range tests {
		status, message := mapWanStatus(input)
		if status != expected {
			t.Fatalf("mapWanStatus(%q) = %q, want %q", input, status, expected)
		}
		if input == "UNKNOWN" && message == "" {
			t.Fatal("UNKNOWN status must include an error message")
		}
	}
}
