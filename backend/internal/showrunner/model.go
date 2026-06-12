package showrunner

import (
	"ai-showrunner-workbench/internal/analysis"
	"ai-showrunner-workbench/internal/screenplay"
	"ai-showrunner-workbench/internal/story"
	"encoding/json"
	"fmt"
)

type CharacterProfile struct {
	ID               string             `json:"id"`
	Name             string             `json:"name"`
	Role             string             `json:"role"`
	Personality      FlexibleStringList `json:"personality"`
	Appearance       FlexibleStringList `json:"appearance"`
	Costume          FlexibleStringList `json:"costume"`
	VoiceStyle       FlexibleStringList `json:"voice_style"`
	KeyMotivation    FlexibleStringList `json:"key_motivation"`
	ConsistencyNotes FlexibleStringList `json:"consistency_notes"`
}

type SceneProfile struct {
	ID               string             `json:"id"`
	Name             string             `json:"name"`
	Location         string             `json:"location"`
	TimeOfDay        string             `json:"time_of_day"`
	Atmosphere       string             `json:"atmosphere"`
	VisualStyle      string             `json:"visual_style"`
	KeyProps         FlexibleStringList `json:"key_props"`
	ConsistencyNotes FlexibleStringList `json:"consistency_notes"`
}

type ChapterBreakdown struct {
	ChapterNumber  int                `json:"chapter_number"`
	ChapterTitle   string             `json:"chapter_title"`
	Summary        string             `json:"summary"`
	MainCharacters FlexibleStringList `json:"main_characters"`
	MainScenes     FlexibleStringList `json:"main_scenes"`
	EmotionalArc   string             `json:"emotional_arc"`
	KeyEvents      FlexibleStringList `json:"key_events"`
}

func (chapter *ChapterBreakdown) UnmarshalJSON(data []byte) error {
	type chapterAlias ChapterBreakdown
	var fields struct {
		ChapterNumber json.RawMessage `json:"chapter_number"`
		*chapterAlias
	}
	fields.chapterAlias = (*chapterAlias)(chapter)
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}

	value, err := unmarshalFlexibleInt(fields.ChapterNumber)
	if err != nil {
		return fmt.Errorf("chapter_number: %w", err)
	}
	chapter.ChapterNumber = value
	return nil
}

type Shot struct {
	ID            string             `json:"id"`
	ChapterNumber int                `json:"chapter_number"`
	SceneID       string             `json:"scene_id"`
	Characters    FlexibleStringList `json:"characters"`
	Dialogue      FlexibleStringList `json:"dialogue"`
	Action        string             `json:"action"`
	Camera        string             `json:"camera"`
	Background    string             `json:"background"`
	DurationHint  string             `json:"duration_hint"`
	ImagePrompt   string             `json:"image_prompt"`
	VideoPrompt   string             `json:"video_prompt"`
	AudioPrompt   string             `json:"audio_prompt"`
}

func (shot *Shot) UnmarshalJSON(data []byte) error {
	type shotFields struct {
		ID              json.RawMessage    `json:"id"`
		ShotID          json.RawMessage    `json:"shot_id"`
		ChapterNumber   json.RawMessage    `json:"chapter_number"`
		SceneID         json.RawMessage    `json:"scene_id"`
		Characters      FlexibleStringList `json:"characters"`
		Dialogue        FlexibleStringList `json:"dialogue"`
		Subtitle        FlexibleStringList `json:"subtitle"`
		Action          string             `json:"action"`
		Camera          string             `json:"camera"`
		Background      string             `json:"background"`
		DurationHint    json.RawMessage    `json:"duration_hint"`
		DurationSeconds json.RawMessage    `json:"duration_seconds"`
		ImagePrompt     string             `json:"image_prompt"`
		VisualPrompt    string             `json:"visual_prompt"`
		VideoPrompt     string             `json:"video_prompt"`
		AudioPrompt     string             `json:"audio_prompt"`
	}

	var fields shotFields
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}

	id := fields.ID
	if len(id) == 0 {
		id = fields.ShotID
	}
	var err error
	if shot.ID, err = unmarshalFlexibleString(id); err != nil {
		return fmt.Errorf("id/shot_id: %w", err)
	}
	if shot.ChapterNumber, err = unmarshalFlexibleInt(fields.ChapterNumber); err != nil {
		return fmt.Errorf("chapter_number: %w", err)
	}
	if shot.SceneID, err = unmarshalFlexibleString(fields.SceneID); err != nil {
		return fmt.Errorf("scene_id: %w", err)
	}

	duration := fields.DurationHint
	if len(duration) == 0 {
		duration = fields.DurationSeconds
	}
	if shot.DurationHint, err = unmarshalFlexibleString(duration); err != nil {
		return fmt.Errorf("duration_hint/duration_seconds: %w", err)
	}
	if len(fields.Dialogue) == 0 {
		fields.Dialogue = fields.Subtitle
	}
	if fields.ImagePrompt == "" {
		fields.ImagePrompt = fields.VisualPrompt
	}

	shot.Characters = fields.Characters
	shot.Dialogue = fields.Dialogue
	shot.Action = fields.Action
	shot.Camera = fields.Camera
	shot.Background = fields.Background
	shot.ImagePrompt = fields.ImagePrompt
	shot.VideoPrompt = fields.VideoPrompt
	shot.AudioPrompt = fields.AudioPrompt
	return nil
}

type AssetPromptSet struct {
	CharacterPrompts  map[string]string `json:"character_prompts"`
	BackgroundPrompts map[string]string `json:"background_prompts"`
	ShotPrompts       map[string]string `json:"shot_prompts"`
	VoicePrompts      map[string]string `json:"voice_prompts"`
}

func (prompts *AssetPromptSet) UnmarshalJSON(data []byte) error {
	prompts.CharacterPrompts = map[string]string{}
	prompts.BackgroundPrompts = map[string]string{}
	prompts.ShotPrompts = map[string]string{}
	prompts.VoicePrompts = map[string]string{}

	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return nil
	}

	var err error
	if prompts.CharacterPrompts, err = unmarshalPromptMap(fields["character_prompts"]); err != nil {
		return fmt.Errorf("character_prompts: %w", err)
	}
	if prompts.BackgroundPrompts, err = unmarshalPromptMap(fields["background_prompts"]); err != nil {
		return fmt.Errorf("background_prompts: %w", err)
	}
	if prompts.ShotPrompts, err = unmarshalPromptMap(fields["shot_prompts"]); err != nil {
		return fmt.Errorf("shot_prompts: %w", err)
	}
	if prompts.VoicePrompts, err = unmarshalPromptMap(fields["voice_prompts"]); err != nil {
		return fmt.Errorf("voice_prompts: %w", err)
	}
	return nil
}

func unmarshalPromptMap(data []byte) (map[string]string, error) {
	result := map[string]string{}
	if len(data) == 0 || string(data) == "null" {
		return result, nil
	}

	var object map[string]json.RawMessage
	if err := json.Unmarshal(data, &object); err == nil {
		for key, value := range object {
			text, err := unmarshalFlexibleString(value)
			if err != nil {
				encoded, marshalErr := json.Marshal(value)
				if marshalErr != nil {
					continue
				}
				text = string(encoded)
			}
			result[key] = text
		}
		return result, nil
	}

	var list []json.RawMessage
	if err := json.Unmarshal(data, &list); err == nil {
		for index, value := range list {
			text, err := unmarshalFlexibleString(value)
			if err == nil {
				result[fmt.Sprintf("%d", index+1)] = text
			}
		}
		return result, nil
	}
	return result, fmt.Errorf("expected prompt object, array, or null")
}

type ShowrunnerResult struct {
	Characters   []CharacterProfile `json:"characters"`
	Scenes       []SceneProfile     `json:"scenes"`
	Chapters     []ChapterBreakdown `json:"chapters"`
	Shots        []Shot             `json:"shots"`
	AssetPrompts AssetPromptSet     `json:"asset_prompts"`
	Warnings     FlexibleStringList `json:"warnings"`
}

type GenerateInput struct {
	Screenplay screenplay.Screenplay      `json:"screenplay"`
	StoryBible story.StoryBible           `json:"story_bible"`
	Chapters   []analysis.ChapterAnalysis `json:"chapters"`
	Style      string                     `json:"style"`
	Language   string                     `json:"language"`
}
