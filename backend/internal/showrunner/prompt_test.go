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
	} {
		if !strings.Contains(prompt, expected) {
			t.Fatalf("prompt missing %q", expected)
		}
	}
}
