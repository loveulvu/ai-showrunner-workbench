package showrunner

import (
	"strings"
	"testing"
)

func TestValidateRequiresCoreAssetsAndShotFields(t *testing.T) {
	result := ShowrunnerResult{
		Characters: []CharacterProfile{{ID: "char_1", Name: "Lead"}},
		Scenes:     []SceneProfile{{ID: "scene_1", Name: "Room"}},
		Chapters:   []ChapterBreakdown{{ChapterNumber: 1, ChapterTitle: "Opening"}},
		Shots: []Shot{{
			ID:            "shot_1",
			ChapterNumber: 1,
			Action:        "The lead enters the room.",
			ImagePrompt:   "animated lead entering a room",
		}},
	}

	validation := Validate(result)
	if !validation.Passed {
		t.Fatalf("Validate() passed = false, errors = %v", validation.Errors)
	}
	if len(validation.Warnings) != 2 {
		t.Fatalf("warnings = %v, want video and audio warnings", validation.Warnings)
	}
}

func TestValidateRejectsInvalidShot(t *testing.T) {
	result := ShowrunnerResult{
		Characters: []CharacterProfile{{ID: "char_1"}},
		Scenes:     []SceneProfile{{ID: "scene_1"}},
		Chapters:   []ChapterBreakdown{{ChapterNumber: 1}},
		Shots:      []Shot{{}},
	}

	validation := Validate(result)
	if validation.Passed {
		t.Fatal("Validate() passed = true, want false")
	}
	for _, expected := range []string{"id is required", "chapter_number must be positive", "must contain action or dialogue", "image_prompt is required"} {
		if !containsSubstring(validation.Errors, expected) {
			t.Fatalf("errors = %v, want substring %q", validation.Errors, expected)
		}
	}
}

func containsSubstring(values []string, target string) bool {
	for _, value := range values {
		if strings.Contains(value, target) {
			return true
		}
	}
	return false
}
