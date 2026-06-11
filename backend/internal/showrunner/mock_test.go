package showrunner

import "testing"

func TestMockResultProvidesDisplayableAssets(t *testing.T) {
	result := MockResult(GenerateInput{Style: "cinematic animation"})
	validation := Validate(result)

	if !validation.Passed {
		t.Fatalf("mock validation errors = %v", validation.Errors)
	}
	if len(result.Characters) < 2 {
		t.Fatalf("characters = %d, want at least 2", len(result.Characters))
	}
	if len(result.Scenes) < 2 {
		t.Fatalf("scenes = %d, want at least 2", len(result.Scenes))
	}
	if len(result.Chapters) < 1 {
		t.Fatalf("chapters = %d, want at least 1", len(result.Chapters))
	}
	if len(result.Shots) < 3 {
		t.Fatalf("shots = %d, want at least 3", len(result.Shots))
	}

	withVideo := 0
	withAudio := 0
	for _, shot := range result.Shots {
		if shot.ImagePrompt == "" {
			t.Fatalf("shot %q image_prompt is empty", shot.ID)
		}
		if shot.VideoPrompt != "" {
			withVideo++
		}
		if shot.AudioPrompt != "" {
			withAudio++
		}
	}
	if withVideo == 0 || withAudio == 0 {
		t.Fatalf("mock prompts with video=%d audio=%d, want at least one each", withVideo, withAudio)
	}
}
