package showrunner

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestBuildPromptRequiresArrayFieldsAndStrictJSON(t *testing.T) {
	prompt := BuildPrompt(GenerateInput{})
	for _, expected := range []string{
		"Every list field must be a JSON array",
		"Never output an array field as a string",
		"strict JSON only",
		"no Markdown fences",
		"cinematic xianxia short drama",
		"first three shots must form one continuous mini-scene",
		"character_visuals",
		"negative_prompt",
		"Do not use ink painting unless",
		"Mode: demo",
		"maximum: 3 shots",
		"at most 2 scenes",
		"at most 1 chapter breakdown",
	} {
		if !strings.Contains(prompt, expected) {
			t.Fatalf("prompt missing %q", expected)
		}
	}
}

func TestBuildPromptSlimsDemoInput(t *testing.T) {
	prompt := BuildPrompt(oversizedDemoInput())
	parts := strings.SplitN(prompt, "Input JSON:\n", 2)
	if len(parts) != 2 {
		t.Fatalf("prompt does not contain input JSON")
	}
	var input GenerateInput
	if err := json.Unmarshal([]byte(parts[1]), &input); err != nil {
		t.Fatalf("decode prompt input JSON: %v", err)
	}

	if len(input.Screenplay.Scenes) != 2 || len(input.Chapters) != 2 {
		t.Fatalf("prompt screenplay scenes/chapters = %d/%d, want 2/2", len(input.Screenplay.Scenes), len(input.Chapters))
	}
}

func TestBuildPromptDefaultsToCinematicShortDrama(t *testing.T) {
	prompt := BuildPrompt(GenerateInput{})
	if !strings.Contains(prompt, `"style":"cinematic xianxia short drama"`) {
		t.Fatalf("prompt missing default style: %s", prompt)
	}
}
