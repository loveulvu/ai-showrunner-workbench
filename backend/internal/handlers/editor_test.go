package handlers

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"ai-showrunner-workbench/internal/editor"

	"github.com/gin-gonic/gin"
)

func TestRenderEditorDemo(t *testing.T) {
	original := renderEditorPlan
	renderEditorPlan = func(ctx context.Context, plan editor.EditingPlan) (editor.EditResult, error) {
		if len(plan.Clips) != 1 || plan.Clips[0].ShotID != "shot_001" {
			t.Fatalf("plan = %#v", plan)
		}
		return editor.EditResult{
			OutputFile:    plan.OutputFile,
			ClipCount:     len(plan.Clips),
			SubtitlesFile: "../outputs/edit/subtitles.srt",
			Warnings:      []string{},
		}, nil
	}
	t.Cleanup(func() { renderEditorPlan = original })
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Request = httptest.NewRequest(http.MethodPost, "/api/editor/render", bytes.NewBufferString(`{
		"output_file":"../outputs/final_demo.mp4",
		"aspect_ratio":"16:9",
		"resolution":"1280x720",
		"fps":24,
		"clips":[{"shot_id":"shot_001","source_url":"https://example.invalid/signed","duration_seconds":5,"subtitle":"Test"}]
	}`))
	context.Request.Header.Set("Content-Type", "application/json")

	RenderEditorDemo(context)
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
	if !bytes.Contains(recorder.Body.Bytes(), []byte(`"output_file":"../outputs/final_demo.mp4"`)) {
		t.Fatalf("body = %s", recorder.Body.String())
	}
}

func TestRenderEditorDemoFFmpegMissing(t *testing.T) {
	original := renderEditorPlan
	renderEditorPlan = func(context.Context, editor.EditingPlan) (editor.EditResult, error) {
		return editor.EditResult{}, errors.New("ffmpeg not found")
	}
	t.Cleanup(func() { renderEditorPlan = original })
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Request = httptest.NewRequest(http.MethodPost, "/api/editor/render", bytes.NewBufferString(`{
		"output_file":"../outputs/final_demo.mp4",
		"clips":[{"shot_id":"shot_001","source_url":"https://example.invalid/signed"}]
	}`))
	context.Request.Header.Set("Content-Type", "application/json")

	RenderEditorDemo(context)
	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
	if !bytes.Contains(recorder.Body.Bytes(), []byte(`ffmpeg not found`)) {
		t.Fatalf("body = %s", recorder.Body.String())
	}
}
