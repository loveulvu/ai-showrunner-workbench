package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ai-showrunner-workbench/internal/video"

	"github.com/gin-gonic/gin"
)

type failingVideoGenerator struct {
	err error
}

func (g failingVideoGenerator) CreateTask(context.Context, video.VideoPrompt) (string, error) {
	return "", g.err
}

func (g failingVideoGenerator) GetTask(context.Context, string) (*video.VideoResult, error) {
	return nil, g.err
}

func TestCreateAndGetVideoTaskEndpoints(t *testing.T) {
	original := videoGenerator
	videoGenerator = video.NewMockVideoGenerator()
	t.Cleanup(func() { videoGenerator = original })
	gin.SetMode(gin.TestMode)

	createRecorder := httptest.NewRecorder()
	createContext, _ := gin.CreateTestContext(createRecorder)
	createContext.Request = httptest.NewRequest(http.MethodPost, "/api/video/tasks", bytes.NewBufferString(`{
		"shot_id":"shot_001",
		"model":"mock-video-model",
		"prompt":"slow cinematic camera move",
		"negative_prompt":"",
		"duration_seconds":5,
		"aspect_ratio":"16:9",
		"subtitle":"",
		"expected_clip_name":"shot_001.mp4"
	}`))
	createContext.Request.Header.Set("Content-Type", "application/json")

	CreateVideoTask(createContext)
	if createRecorder.Code != http.StatusAccepted {
		t.Fatalf("create status = %d, body = %s", createRecorder.Code, createRecorder.Body.String())
	}

	var created struct {
		TaskID string `json:"task_id"`
	}
	if err := json.Unmarshal(createRecorder.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode create response: %v", err)
	}
	if created.TaskID == "" {
		t.Fatal("create response task_id is empty")
	}

	getRecorder := httptest.NewRecorder()
	getContext, _ := gin.CreateTestContext(getRecorder)
	getContext.Request = httptest.NewRequest(http.MethodGet, "/api/video/tasks/"+created.TaskID, nil)
	getContext.AddParam("task_id", created.TaskID)

	GetVideoTask(getContext)
	if getRecorder.Code != http.StatusOK {
		t.Fatalf("get status = %d, body = %s", getRecorder.Code, getRecorder.Body.String())
	}
	if !bytes.Contains(getRecorder.Body.Bytes(), []byte(`"status":"succeeded"`)) {
		t.Fatalf("get body = %s, want succeeded status", getRecorder.Body.String())
	}
	if !bytes.Contains(getRecorder.Body.Bytes(), []byte(`"video_url"`)) {
		t.Fatalf("get body = %s, want video_url", getRecorder.Body.String())
	}
}

func TestGetVideoTaskNotFound(t *testing.T) {
	original := videoGenerator
	videoGenerator = video.NewMockVideoGenerator()
	t.Cleanup(func() { videoGenerator = original })
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Request = httptest.NewRequest(http.MethodGet, "/api/video/tasks/missing", nil)
	context.AddParam("task_id", "missing")

	GetVideoTask(context)
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
}

func TestConfigureVideoGeneratorFromEnvMockAndWan(t *testing.T) {
	original := videoGenerator
	t.Cleanup(func() { videoGenerator = original })

	t.Setenv("VIDEO_PROVIDER", "mock")
	config, err := ConfigureVideoGeneratorFromEnv()
	if err != nil {
		t.Fatalf("ConfigureVideoGeneratorFromEnv() mock error = %v", err)
	}
	if config.Provider != "mock" {
		t.Fatalf("Provider = %q, want mock", config.Provider)
	}

	t.Setenv("VIDEO_PROVIDER", "wan")
	t.Setenv("VIDEO_BASE_URL", "https://example.com/api/v1")
	t.Setenv("VIDEO_API_KEY", "test-key")
	if _, err := ConfigureVideoGeneratorFromEnv(); err != nil {
		t.Fatalf("ConfigureVideoGeneratorFromEnv() wan error = %v", err)
	}
}

func TestCreateVideoTaskReturnsBadGatewayForWanNetworkError(t *testing.T) {
	original := videoGenerator
	videoGenerator = failingVideoGenerator{err: &video.Error{
		Kind:    video.ErrorKindUpstream,
		Message: "Could not reach Wan video service",
	}}
	t.Cleanup(func() { videoGenerator = original })
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Request = httptest.NewRequest(http.MethodPost, "/api/video/tasks", bytes.NewBufferString(`{"shot_id":"shot_001","prompt":"test"}`))
	context.Request.Header.Set("Content-Type", "application/json")

	CreateVideoTask(context)

	if recorder.Code != http.StatusBadGateway {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
	if !bytes.Contains(recorder.Body.Bytes(), []byte(`"message":"Could not reach Wan video service"`)) {
		t.Fatalf("body = %s", recorder.Body.String())
	}
}

func TestCreateVideoTaskReturnsBadRequestForInvalidPrompt(t *testing.T) {
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Request = httptest.NewRequest(http.MethodPost, "/api/video/tasks", bytes.NewBufferString(`{"shot_id":"","prompt":""}`))
	context.Request.Header.Set("Content-Type", "application/json")

	CreateVideoTask(context)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
}

func TestGetVideoTaskReturnsBadGatewayForWanNetworkError(t *testing.T) {
	original := videoGenerator
	videoGenerator = failingVideoGenerator{err: &video.Error{
		Kind:    video.ErrorKindUpstream,
		Message: "Could not reach Wan video service",
	}}
	t.Cleanup(func() { videoGenerator = original })

	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Request = httptest.NewRequest(http.MethodGet, "/api/video/tasks/task-1", nil)
	context.AddParam("task_id", "task-1")

	GetVideoTask(context)

	if recorder.Code != http.StatusBadGateway {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
}
