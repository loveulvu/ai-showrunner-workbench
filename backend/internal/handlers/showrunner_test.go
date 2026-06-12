package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGenerateShowrunnerMockEndpoint(t *testing.T) {
	t.Setenv("AI_PROVIDER", "mock")
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Request = httptest.NewRequest(http.MethodPost, "/api/showrunner/generate", bytes.NewBufferString(`{
		"screenplay":{"title":"Test","source_chapters":[],"characters":[],"scenes":[]},
		"story_bible":{"title":"Test","global_characters":[],"timeline":[],"scene_plan":[]},
		"chapters":[],
		"style":"cinematic animation",
		"language":"English"
	}`))
	context.Request.Header.Set("Content-Type", "application/json")

	GenerateShowrunner(context)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
	for _, expected := range []string{`"characters"`, `"scenes"`, `"shots"`, `"asset_prompts"`, `"warnings"`} {
		if !bytes.Contains(recorder.Body.Bytes(), []byte(expected)) {
			t.Fatalf("body missing %s: %s", expected, recorder.Body.String())
		}
	}
	if !bytes.Contains(recorder.Body.Bytes(), []byte(`"mode":"demo"`)) {
		t.Fatalf("body does not default to demo mode: %s", recorder.Body.String())
	}
}

func TestGenerateShowrunnerReturnsStructuredServiceError(t *testing.T) {
	t.Setenv("AI_PROVIDER", "invalid-provider")
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Request = httptest.NewRequest(http.MethodPost, "/api/showrunner/generate", bytes.NewBufferString(`{}`))
	context.Request.Header.Set("Content-Type", "application/json")

	GenerateShowrunner(context)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
	var body map[string]string
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body["stage"] != "service" || body["message"] == "" {
		t.Fatalf("body = %#v", body)
	}
}
