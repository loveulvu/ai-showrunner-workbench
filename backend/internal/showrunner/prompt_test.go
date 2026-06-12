package showrunner

import (
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
	} {
		if !strings.Contains(prompt, expected) {
			t.Fatalf("prompt missing %q", expected)
		}
	}
}

func TestBuildPromptDefaultsToCinematicShortDrama(t *testing.T) {
	prompt := BuildPrompt(GenerateInput{})
	if !strings.Contains(prompt, `"style": "cinematic xianxia short drama"`) {
		t.Fatalf("prompt missing default style: %s", prompt)
	}
}
