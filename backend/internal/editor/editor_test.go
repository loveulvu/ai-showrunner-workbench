package editor

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestClipFileName(t *testing.T) {
	name, err := ClipFileName("shot_001")
	if err != nil {
		t.Fatalf("ClipFileName() error = %v", err)
	}
	if name != "clip_shot_001.mp4" {
		t.Fatalf("name = %q, want clip_shot_001.mp4", name)
	}

	name, err = ClipFileName("../shot unsafe")
	if err != nil {
		t.Fatalf("ClipFileName() unsafe error = %v", err)
	}
	if strings.Contains(name, "..") || strings.ContainsAny(name, `/\`) {
		t.Fatalf("unsafe name = %q", name)
	}
}

func TestBuildConcatList(t *testing.T) {
	root := t.TempDir()
	clipOne := filepath.Join(root, "clips", "clip_shot_001.mp4")
	clipTwo := filepath.Join(root, "clips", "clip_shot_002.mp4")
	plan := EditingPlan{
		OutputFile: filepath.Join(root, "final_demo.mp4"),
		Clips: []ClipAsset{
			{ShotID: "shot_001", LocalPath: clipOne},
			{ShotID: "shot_002", LocalPath: clipTwo},
		},
	}

	path, err := BuildConcatList(plan)
	if err != nil {
		t.Fatalf("BuildConcatList() error = %v", err)
	}
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read concat list: %v", err)
	}
	text := string(content)
	if !strings.Contains(text, filepath.ToSlash(clipOne)) || !strings.Contains(text, filepath.ToSlash(clipTwo)) {
		t.Fatalf("concat list = %q", text)
	}
}

func TestBuildSRTTimeline(t *testing.T) {
	root := t.TempDir()
	plan := EditingPlan{
		OutputFile: filepath.Join(root, "final_demo.mp4"),
		Clips: []ClipAsset{
			{ShotID: "shot_001", DurationSeconds: 5, Subtitle: "First subtitle."},
			{ShotID: "shot_002", DurationSeconds: 3, Subtitle: ""},
			{ShotID: "shot_003", DurationSeconds: 2, Subtitle: "Second subtitle."},
		},
	}

	path, err := BuildSRT(plan)
	if err != nil {
		t.Fatalf("BuildSRT() error = %v", err)
	}
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read subtitles: %v", err)
	}
	text := string(content)
	for _, expected := range []string{
		"00:00:00,000 --> 00:00:05,000",
		"First subtitle.",
		"00:00:08,000 --> 00:00:10,000",
		"Second subtitle.",
	} {
		if !strings.Contains(text, expected) {
			t.Fatalf("subtitles missing %q: %s", expected, text)
		}
	}
}

func TestRunFFmpegNotFound(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	_, err := RunFFmpeg(context.Background(), EditingPlan{OutputFile: filepath.Join(t.TempDir(), "final_demo.mp4")}, "concat.txt")
	if err == nil {
		t.Fatal("RunFFmpeg() error = nil, want ffmpeg not found")
	}
	if err.Error() != "ffmpeg not found" {
		t.Fatalf("RunFFmpeg() error = %q, want ffmpeg not found", err)
	}
}

func TestRedactSourceURL(t *testing.T) {
	sourceURL := "https://example.com/video.mp4?signature=secret"
	got := redactSourceURL("Get "+sourceURL+": connection failed", sourceURL)
	if strings.Contains(got, sourceURL) || strings.Contains(got, "signature=secret") {
		t.Fatalf("redacted error contains source URL: %q", got)
	}
}
