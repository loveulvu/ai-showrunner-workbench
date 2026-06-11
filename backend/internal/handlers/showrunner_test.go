package handlers

import (
	"bytes"
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
}
