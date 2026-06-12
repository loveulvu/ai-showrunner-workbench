package showrunner

import (
	"fmt"
	"testing"

	"ai-showrunner-workbench/internal/analysis"
	"ai-showrunner-workbench/internal/screenplay"
	"ai-showrunner-workbench/internal/story"
)

func TestPrepareInputDefaultsToSlimDemoMode(t *testing.T) {
	input := PrepareInput(oversizedDemoInput())

	if input.Mode != ShowrunnerModeDemo {
		t.Fatalf("mode = %q, want demo", input.Mode)
	}
	if len(input.Screenplay.Characters) != 4 || len(input.Screenplay.Scenes) != 2 {
		t.Fatalf("screenplay characters/scenes = %d/%d, want 4/2", len(input.Screenplay.Characters), len(input.Screenplay.Scenes))
	}
	if len(input.StoryBible.GlobalCharacters) != 4 || len(input.StoryBible.ScenePlan) != 2 {
		t.Fatalf("story bible characters/scenes = %d/%d, want 4/2", len(input.StoryBible.GlobalCharacters), len(input.StoryBible.ScenePlan))
	}
	if len(input.Chapters) != 2 || len(input.Chapters[0].FactualAnchors) != 5 {
		t.Fatalf("chapter analyses/factual anchors = %d/%d, want 2/5", len(input.Chapters), len(input.Chapters[0].FactualAnchors))
	}
}

func TestLimitResultForDemoMode(t *testing.T) {
	result := ShowrunnerResult{
		Characters: make([]CharacterProfile, 6),
		Scenes:     make([]SceneProfile, 4),
		Chapters:   make([]ChapterBreakdown, 3),
		Shots:      make([]Shot, 6),
		Warnings:   FlexibleStringList{"1", "2", "3", "4", "5", "6"},
	}

	got := LimitResultForMode(result, "")
	if got.Mode != ShowrunnerModeDemo || len(got.Characters) != 4 || len(got.Scenes) != 2 || len(got.Chapters) != 1 || len(got.Shots) != 3 || len(got.Warnings) != 5 {
		t.Fatalf("limited demo result = mode:%q characters:%d scenes:%d chapters:%d shots:%d warnings:%d", got.Mode, len(got.Characters), len(got.Scenes), len(got.Chapters), len(got.Shots), len(got.Warnings))
	}
}

func oversizedDemoInput() GenerateInput {
	input := GenerateInput{}
	for index := 1; index <= 6; index++ {
		id := fmt.Sprintf("%d", index)
		input.Screenplay.SourceChapters = append(input.Screenplay.SourceChapters, screenplay.SourceChapter{Number: index})
		input.Screenplay.Characters = append(input.Screenplay.Characters, screenplay.Character{ID: "char_" + id})
		input.Screenplay.Scenes = append(input.Screenplay.Scenes, screenplay.Scene{ID: "scene_" + id})
		input.StoryBible.GlobalCharacters = append(input.StoryBible.GlobalCharacters, story.Character{ID: "char_" + id})
		input.StoryBible.Timeline = append(input.StoryBible.Timeline, story.TimelineEvent{ChapterNumber: index})
		input.StoryBible.ScenePlan = append(input.StoryBible.ScenePlan, story.ScenePlanItem{ID: "scene_" + id})
		input.Chapters = append(input.Chapters, analysis.ChapterAnalysis{
			ChapterNumber:  index,
			FactualAnchors: []string{"1", "2", "3", "4", "5", "6"},
		})
	}
	return input
}
