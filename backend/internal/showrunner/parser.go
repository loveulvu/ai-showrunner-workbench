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
	for index := range result.Characters {
		normalizeCharacterProfile(&result.Characters[index])
	}
	if result.Scenes == nil {
		result.Scenes = []SceneProfile{}
	}
	for index := range result.Scenes {
		normalizeSceneProfile(&result.Scenes[index])
	}
	if result.Chapters == nil {
		result.Chapters = []ChapterBreakdown{}
	}
	for index := range result.Chapters {
		normalizeChapterBreakdown(&result.Chapters[index])
	}
	if result.Shots == nil {
		result.Shots = []Shot{}
	}
	for index := range result.Shots {
		normalizeShot(&result.Shots[index])
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
		result.Warnings = FlexibleStringList{}
	}
}

func normalizeCharacterProfile(profile *CharacterProfile) {
	normalizeFlexibleList(&profile.Personality)
	normalizeFlexibleList(&profile.Appearance)
	normalizeFlexibleList(&profile.Costume)
	normalizeFlexibleList(&profile.VoiceStyle)
	normalizeFlexibleList(&profile.KeyMotivation)
	normalizeFlexibleList(&profile.ConsistencyNotes)
}

func normalizeSceneProfile(profile *SceneProfile) {
	normalizeFlexibleList(&profile.KeyProps)
	normalizeFlexibleList(&profile.ConsistencyNotes)
}

func normalizeChapterBreakdown(chapter *ChapterBreakdown) {
	normalizeFlexibleList(&chapter.MainCharacters)
	normalizeFlexibleList(&chapter.MainScenes)
	normalizeFlexibleList(&chapter.KeyEvents)
}

func normalizeShot(shot *Shot) {
	normalizeFlexibleList(&shot.Characters)
	normalizeFlexibleList(&shot.Dialogue)
}

func normalizeFlexibleList(list *FlexibleStringList) {
	if *list == nil {
		*list = FlexibleStringList{}
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
