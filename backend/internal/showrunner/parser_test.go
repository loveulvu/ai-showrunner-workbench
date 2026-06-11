package showrunner

import "testing"

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
