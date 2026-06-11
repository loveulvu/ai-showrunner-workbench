package showrunner

import (
	"encoding/json"
	"fmt"
	"strings"
)

func ParseJSON(raw string) (ShowrunnerResult, error) {
	var result ShowrunnerResult
	text := stripJSONFence(raw)
	if err := json.Unmarshal([]byte(text), &result); err != nil {
		return result, fmt.Errorf("parse showrunner JSON: %w", err)
	}
	normalizeResult(&result)
	return result, nil
}

func normalizeResult(result *ShowrunnerResult) {
	if result.Characters == nil {
		result.Characters = []CharacterProfile{}
	}
	if result.Scenes == nil {
		result.Scenes = []SceneProfile{}
	}
	if result.Chapters == nil {
		result.Chapters = []ChapterBreakdown{}
	}
	if result.Shots == nil {
		result.Shots = []Shot{}
	}
	if result.AssetPrompts.CharacterPrompts == nil {
		result.AssetPrompts.CharacterPrompts = map[string]string{}
	}
	if result.AssetPrompts.BackgroundPrompts == nil {
		result.AssetPrompts.BackgroundPrompts = map[string]string{}
	}
	if result.AssetPrompts.ShotPrompts == nil {
		result.AssetPrompts.ShotPrompts = map[string]string{}
	}
	if result.AssetPrompts.VoicePrompts == nil {
		result.AssetPrompts.VoicePrompts = map[string]string{}
	}
	if result.Warnings == nil {
		result.Warnings = []string{}
	}
}

func stripJSONFence(raw string) string {
	text := strings.TrimSpace(raw)
	if !strings.HasPrefix(text, "```") {
		return text
	}
	if newline := strings.IndexByte(text, '\n'); newline >= 0 {
		text = text[newline+1:]
	}
	if end := strings.LastIndex(text, "```"); end >= 0 {
		text = text[:end]
	}
	return strings.TrimSpace(text)
}
