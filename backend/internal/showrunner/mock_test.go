package showrunner

import "testing"

func TestMockResultProvidesDisplayableAssets(t *testing.T) {
	result := MockResult(GenerateInput{})
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
	if result.Mode != ShowrunnerModeDemo || len(result.Shots) != 3 {
		t.Fatalf("demo mock mode/shots = %q/%d, want demo/3", result.Mode, len(result.Shots))
	}

	withVideo := 0
	withAudio := 0
	sceneID := result.Shots[0].SceneID
	for _, shot := range result.Shots {
		if shot.ImagePrompt == "" {
			t.Fatalf("shot %q image_prompt is empty", shot.ID)
		}
		if shot.VideoPrompt != "" {
			withVideo++
		}
		if shot.NegativePrompt == "" || shot.CameraMovement == "" || shot.ContinuityNotes == "" {
			t.Fatalf("shot %q missing video continuity fields: %#v", shot.ID, shot)
		}
		if shot.SceneID != sceneID {
			t.Fatalf("first three shots are not continuous: %q != %q", shot.SceneID, sceneID)
		}
		if shot.AudioPrompt != "" {
			withAudio++
		}
	}
	if withVideo == 0 || withAudio == 0 {
		t.Fatalf("mock prompts with video=%d audio=%d, want at least one each", withVideo, withAudio)
	}
	if result.Characters[0].VisualIdentity.ConsistencyPrompt == "" || result.Scenes[0].VisualIdentity.ConsistencyPrompt == "" {
		t.Fatal("mock continuity bible is incomplete")
	}
	if result.Scenes[0].VisualStyle != defaultVideoStyle {
		t.Fatalf("default visual style = %q, want %q", result.Scenes[0].VisualStyle, defaultVideoStyle)
	}
}
