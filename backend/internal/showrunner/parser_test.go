package showrunner

import (
	"encoding/json"
	"testing"
)

func TestParseJSON(t *testing.T) {
	raw := "```json\n" + `{
		"characters":[{"id":"char_1","name":"Lead"}],
		"scenes":[{"id":"scene_1","name":"Room"}],
		"chapters":[{"chapter_number":1,"chapter_title":"Opening"}],
		"shots":[{"id":"shot_1","chapter_number":1,"action":"The lead enters.","image_prompt":"animated entrance"}],
		"asset_prompts":{},
		"warnings":[]
	}` + "\n```"

	result, err := ParseJSON(raw)
	if err != nil {
		t.Fatalf("ParseJSON() error = %v", err)
	}
	if len(result.Shots) != 1 || result.Shots[0].ID != "shot_1" {
		t.Fatalf("shots = %#v", result.Shots)
	}
	if result.AssetPrompts.CharacterPrompts == nil {
		t.Fatal("CharacterPrompts = nil, want normalized empty map")
	}
}

func TestParseJSONExtractsObjectWithSurroundingText(t *testing.T) {
	raw := `Here is the requested result:
	{"shots":[{"shot_id":7,"chapter_number":"2","duration_seconds":6,"visual_prompt":"wide shot","subtitle":"Hello"}]}
	This is ready for production.`

	result, err := ParseJSON(raw)
	if err != nil {
		t.Fatalf("ParseJSON() error = %v", err)
	}
	shot := result.Shots[0]
	if shot.ID != "7" || shot.ChapterNumber != 2 || shot.DurationHint != "6" {
		t.Fatalf("shot = %#v", shot)
	}
	if shot.ImagePrompt != "wide shot" || shot.Dialogue.Text() != "Hello" {
		t.Fatalf("shot compatibility fields = %#v", shot)
	}
}

func TestParseJSONAcceptsDurationAndShotIDStringOrNumber(t *testing.T) {
	raw := `{"shots":[
		{"shot_id":"shot-a","chapter_number":1,"duration_seconds":"5","image_prompt":"a"},
		{"shot_id":2,"chapter_number":"2","duration_seconds":6,"image_prompt":"b"}
	]}`

	result, err := ParseJSON(raw)
	if err != nil {
		t.Fatalf("ParseJSON() error = %v", err)
	}
	if result.Shots[0].ID != "shot-a" || result.Shots[0].DurationHint != "5" {
		t.Fatalf("first shot = %#v", result.Shots[0])
	}
	if result.Shots[1].ID != "2" || result.Shots[1].ChapterNumber != 2 || result.Shots[1].DurationHint != "6" {
		t.Fatalf("second shot = %#v", result.Shots[1])
	}
}

func TestParseJSONAcceptsContinuityAndVideoPromptFields(t *testing.T) {
	raw := `{
		"characters":[{"id":"lead","visual_identity":{"age":"24","face":"oval","hairstyle":"high ponytail","costume":"blue robe","color_palette":"blue silver","expression_baseline":"focused","body_type":"lean","consistency_prompt":"same lead"}}],
		"scenes":[{"id":"hall","visual_identity":{"architecture":"timber hall","lighting":"moonlight","color_palette":"blue amber","atmosphere":"tense","key_props":"bronze key","consistency_prompt":"same hall"}}],
		"shots":[{"id":"shot-1","character_visuals":"same lead","scene_visuals":"same hall","camera_angle":"eye level","camera_movement":"push in","composition":"centered","lighting":"moonlight","motion":"raises key","continuity_notes":"same position","video_prompt":"cinematic action","negative_prompt":"blurry, text"}]
	}`

	result, err := ParseJSON(raw)
	if err != nil {
		t.Fatalf("ParseJSON() error = %v", err)
	}
	if result.Characters[0].VisualIdentity.Costume != "blue robe" {
		t.Fatalf("character identity = %#v", result.Characters[0].VisualIdentity)
	}
	if result.Scenes[0].VisualIdentity.KeyProps.Text() != "bronze key" {
		t.Fatalf("scene identity = %#v", result.Scenes[0].VisualIdentity)
	}
	shot := result.Shots[0]
	if shot.CharacterVisuals.Text() != "same lead" || shot.CameraMovement != "push in" || shot.NegativePrompt != "blurry, text" {
		t.Fatalf("shot = %#v", shot)
	}
}

func TestParseJSONAcceptsTopLevelAliasesAndFlexibleAssetPrompts(t *testing.T) {
	raw := `{
		"shot_list":[{"id":"shot-1","chapter_number":1,"image_prompt":"frame"}],
		"assetPrompts":{"shot_prompts":["first prompt", "second prompt"]}
	}`

	result, err := ParseJSON(raw)
	if err != nil {
		t.Fatalf("ParseJSON() error = %v", err)
	}
	if len(result.Shots) != 1 || len(result.AssetPrompts.ShotPrompts) != 2 {
		t.Fatalf("result = %#v", result)
	}
}

func TestParseJSONDoesNotFailWhenAssetPromptsHasUnexpectedShape(t *testing.T) {
	result, err := ParseJSON(`{"shots":[{"id":"shot-1","image_prompt":"frame"}],"asset_prompts":["unexpected"]}`)
	if err != nil {
		t.Fatalf("ParseJSON() error = %v", err)
	}
	if len(result.Shots) != 1 || result.AssetPrompts.ShotPrompts == nil {
		t.Fatalf("result = %#v", result)
	}
}

func TestParseJSONRejectsInvalidJSON(t *testing.T) {
	if _, err := ParseJSON(`{"shots":`); err == nil {
		t.Fatal("ParseJSON() error = nil, want parse error")
	}
}

func TestParseJSONAcceptsRealLikeStringLists(t *testing.T) {
	raw := `{
		"characters":[{
			"id":"char_1",
			"name":"Lead",
			"personality":"calm, controlled",
			"appearance":"dark coat",
			"costume":["dark coat"],
			"voice_style":"quiet",
			"key_motivation":"find the truth",
			"consistency_notes":null
		}],
		"scenes":[{"id":"scene_1","name":"Room","key_props":"clock","consistency_notes":"same lighting"}],
		"chapters":[{"chapter_number":1,"chapter_title":"Opening","main_characters":"char_1","main_scenes":["scene_1"],"key_events":"arrival"}],
		"shots":[{"id":"shot_1","chapter_number":1,"characters":"char_1","dialogue":"Hello","image_prompt":"animated entrance"}],
		"asset_prompts":{},
		"warnings":"video prompt missing"
	}`

	result, err := ParseJSON(raw)
	if err != nil {
		t.Fatalf("ParseJSON() error = %v", err)
	}
	if len(result.Characters[0].Personality) != 1 || result.Characters[0].Personality[0] != "calm, controlled" {
		t.Fatalf("personality = %#v", result.Characters[0].Personality)
	}
	if result.Shots[0].Dialogue.Text() != "Hello" {
		t.Fatalf("dialogue = %#v", result.Shots[0].Dialogue)
	}
	if len(result.Warnings) != 1 {
		t.Fatalf("warnings = %#v", result.Warnings)
	}
}

func TestMockResultJSONStillParses(t *testing.T) {
	payload, err := json.Marshal(MockResult(GenerateInput{}))
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}
	if _, err := ParseJSON(string(payload)); err != nil {
		t.Fatalf("ParseJSON(mock result) error = %v", err)
	}
}
