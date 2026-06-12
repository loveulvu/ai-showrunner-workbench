package showrunner

import (
	"ai-showrunner-workbench/internal/analysis"
)

func PrepareInput(input GenerateInput) GenerateInput {
	input.Mode = NormalizeMode(input.Mode)
	if input.Mode != ShowrunnerModeDemo {
		return input
	}

	input.Screenplay.SourceChapters = take(input.Screenplay.SourceChapters, 2)
	input.Screenplay.Characters = take(input.Screenplay.Characters, 4)
	input.Screenplay.Scenes = take(input.Screenplay.Scenes, 2)

	input.StoryBible.GlobalCharacters = take(input.StoryBible.GlobalCharacters, 4)
	input.StoryBible.Timeline = take(input.StoryBible.Timeline, 5)
	input.StoryBible.ScenePlan = take(input.StoryBible.ScenePlan, 2)

	input.Chapters = take(input.Chapters, 2)
	for index := range input.Chapters {
		input.Chapters[index] = slimChapter(input.Chapters[index])
	}
	return input
}

func NormalizeMode(mode ShowrunnerMode) ShowrunnerMode {
	if mode == ShowrunnerModeFull {
		return ShowrunnerModeFull
	}
	return ShowrunnerModeDemo
}

func LimitResultForMode(result ShowrunnerResult, mode ShowrunnerMode) ShowrunnerResult {
	result.Mode = NormalizeMode(mode)
	if result.Mode != ShowrunnerModeDemo {
		return result
	}
	result.Characters = take(result.Characters, 4)
	result.Scenes = take(result.Scenes, 2)
	result.Chapters = take(result.Chapters, 1)
	result.Shots = take(result.Shots, 3)
	result.Warnings = FlexibleStringList(take([]string(result.Warnings), 5))
	return result
}

func slimChapter(chapter analysis.ChapterAnalysis) analysis.ChapterAnalysis {
	chapter.Characters = take(chapter.Characters, 4)
	chapter.Locations = take(chapter.Locations, 2)
	chapter.KeyEvents = take(chapter.KeyEvents, 5)
	chapter.Conflicts = take(chapter.Conflicts, 3)
	chapter.SceneCandidates = take(chapter.SceneCandidates, 2)
	chapter.FactualAnchors = take(chapter.FactualAnchors, 5)
	return chapter
}

func take[T any](values []T, limit int) []T {
	if len(values) <= limit {
		return values
	}
	return values[:limit]
}
