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
