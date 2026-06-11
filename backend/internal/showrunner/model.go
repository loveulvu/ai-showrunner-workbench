package showrunner

import (
	"ai-showrunner-workbench/internal/analysis"
	"ai-showrunner-workbench/internal/screenplay"
	"ai-showrunner-workbench/internal/story"
)

type CharacterProfile struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	Role             string   `json:"role"`
	Personality      []string `json:"personality"`
	Appearance       string   `json:"appearance"`
	Costume          string   `json:"costume"`
	VoiceStyle       string   `json:"voice_style"`
	KeyMotivation    string   `json:"key_motivation"`
	ConsistencyNotes []string `json:"consistency_notes"`
}

type SceneProfile struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	Location         string   `json:"location"`
	TimeOfDay        string   `json:"time_of_day"`
	Atmosphere       string   `json:"atmosphere"`
	VisualStyle      string   `json:"visual_style"`
	KeyProps         []string `json:"key_props"`
	ConsistencyNotes []string `json:"consistency_notes"`
}

type ChapterBreakdown struct {
	ChapterNumber  int      `json:"chapter_number"`
	ChapterTitle   string   `json:"chapter_title"`
	Summary        string   `json:"summary"`
	MainCharacters []string `json:"main_characters"`
	MainScenes     []string `json:"main_scenes"`
	EmotionalArc   string   `json:"emotional_arc"`
	KeyEvents      []string `json:"key_events"`
}

type Shot struct {
	ID            string   `json:"id"`
	ChapterNumber int      `json:"chapter_number"`
	SceneID       string   `json:"scene_id"`
	Characters    []string `json:"characters"`
	Dialogue      string   `json:"dialogue"`
	Action        string   `json:"action"`
	Camera        string   `json:"camera"`
	Background    string   `json:"background"`
	DurationHint  string   `json:"duration_hint"`
	ImagePrompt   string   `json:"image_prompt"`
	VideoPrompt   string   `json:"video_prompt"`
	AudioPrompt   string   `json:"audio_prompt"`
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
	Warnings     []string           `json:"warnings"`
}

type GenerateInput struct {
	Screenplay screenplay.Screenplay      `json:"screenplay"`
	StoryBible story.StoryBible           `json:"story_bible"`
	Chapters   []analysis.ChapterAnalysis `json:"chapters"`
	Style      string                     `json:"style"`
	Language   string                     `json:"language"`
}
