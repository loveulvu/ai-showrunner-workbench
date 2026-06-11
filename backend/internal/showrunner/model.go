package showrunner

import (
	"ai-showrunner-workbench/internal/analysis"
	"ai-showrunner-workbench/internal/screenplay"
	"ai-showrunner-workbench/internal/story"
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

type AssetPromptSet struct {
	CharacterPrompts  map[string]string `json:"character_prompts"`
	BackgroundPrompts map[string]string `json:"background_prompts"`
	ShotPrompts       map[string]string `json:"shot_prompts"`
	VoicePrompts      map[string]string `json:"voice_prompts"`
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
