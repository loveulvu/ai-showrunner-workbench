package showrunner

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func ParseJSON(raw string) (ShowrunnerResult, error) {
	var result ShowrunnerResult
	text, err := ExtractJSONObject(raw)
	if err != nil {
		return result, &StageError{Stage: StageExtractJSON, Message: "could not extract a valid JSON object from showrunner output", Err: err}
	}
	compatible, err := normalizeTopLevelFields([]byte(text))
	if err != nil {
		return result, &StageError{Stage: StageParseJSON, Message: "could not normalize showrunner JSON fields", Err: err}
	}
	if err := json.Unmarshal(compatible, &result); err != nil {
		return result, &StageError{Stage: StageParseJSON, Message: "could not parse showrunner JSON fields", Err: err}
	}
	normalizeResult(&result)
	return result, nil
}

func normalizeTopLevelFields(data []byte) ([]byte, error) {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return nil, err
	}

	aliases := map[string][]string{
		"characters":    {"character_profiles"},
		"scenes":        {"scene_profiles"},
		"chapters":      {"chapter_breakdowns"},
		"shots":         {"shot_list", "storyboard"},
		"asset_prompts": {"assetPrompts"},
	}
	for canonical, candidates := range aliases {
		if _, exists := fields[canonical]; exists {
			continue
		}
		for _, alias := range candidates {
			if value, exists := fields[alias]; exists {
				fields[canonical] = value
				break
			}
		}
	}
	return json.Marshal(fields)
}

func ExtractJSONObject(raw string) (string, error) {
	data := []byte(raw)
	for index, value := range data {
		if value != '{' {
			continue
		}

		decoder := json.NewDecoder(bytes.NewReader(data[index:]))
		var object map[string]json.RawMessage
		if err := decoder.Decode(&object); err != nil || object == nil {
			continue
		}

		var rawObject json.RawMessage
		decoder = json.NewDecoder(bytes.NewReader(data[index:]))
		if err := decoder.Decode(&rawObject); err == nil {
			return string(rawObject), nil
		}
	}
	return "", fmt.Errorf("no valid JSON object found")
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
	normalizeFlexibleList(&profile.VisualIdentity.KeyProps)
}

func normalizeChapterBreakdown(chapter *ChapterBreakdown) {
	normalizeFlexibleList(&chapter.MainCharacters)
	normalizeFlexibleList(&chapter.MainScenes)
	normalizeFlexibleList(&chapter.KeyEvents)
}

func normalizeShot(shot *Shot) {
	normalizeFlexibleList(&shot.Characters)
	normalizeFlexibleList(&shot.Dialogue)
	normalizeFlexibleList(&shot.CharacterVisuals)
}

func normalizeFlexibleList(list *FlexibleStringList) {
	if *list == nil {
		*list = FlexibleStringList{}
	}
}
