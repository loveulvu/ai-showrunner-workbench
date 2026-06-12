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
	if len(validation.Warnings) != 3 {
		t.Fatalf("warnings = %v, want video, negative, and audio warnings", validation.Warnings)
	}
}

func TestValidateAllowsPartialAssetsAndShotFieldsWithWarnings(t *testing.T) {
	result := ShowrunnerResult{
		Shots: []Shot{{}},
	}

	validation := Validate(result)
	if !validation.Passed {
		t.Fatalf("Validate() passed = false, errors = %v", validation.Errors)
	}
	for _, expected := range []string{"characters is empty", "scenes is empty", "chapters is empty", "image_prompt is empty"} {
		if !containsSubstring(validation.Warnings, expected) {
			t.Fatalf("warnings = %v, want substring %q", validation.Warnings, expected)
		}
	}
}

func TestValidateFailsOnlyWhenShotsAreEmpty(t *testing.T) {
	validation := Validate(ShowrunnerResult{})
	if validation.Passed {
		t.Fatal("Validate() passed = true, want false")
	}
	if len(validation.Errors) != 1 || !containsSubstring(validation.Errors, "shots must contain at least one shot") {
		t.Fatalf("errors = %v", validation.Errors)
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
