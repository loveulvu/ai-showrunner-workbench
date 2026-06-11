package ai

import (
	"context"

	"ai-showrunner-workbench/internal/analysis"
	"ai-showrunner-workbench/internal/fidelity"
	"ai-showrunner-workbench/internal/novel"
	"ai-showrunner-workbench/internal/screenplay"
	"ai-showrunner-workbench/internal/showrunner"
	"ai-showrunner-workbench/internal/story"
)

type Config struct {
	Provider       string
	APIKey         string
	BaseURL        string
	Model          string
	TimeoutSeconds int
}

const ProviderMock = "mock"
const ProviderReal = "real"
const ProviderQwen = "qwen"

const defaultQwenBaseURL = "https://dashscope-intl.aliyuncs.com/compatible-mode/v1"
const defaultQwenModel = "qwen-plus"

type Client interface {
	AnalyzeChapter(ctx context.Context, chapter novel.Chapter) (analysis.ChapterAnalysis, error)
	MergeStoryBible(ctx context.Context, analyses []analysis.ChapterAnalysis) (story.StoryBible, error)
	GenerateScreenplay(ctx context.Context, bible story.StoryBible, analyses []analysis.ChapterAnalysis) (screenplay.Screenplay, error)
	GenerateShowrunner(ctx context.Context, input showrunner.GenerateInput) (showrunner.ShowrunnerResult, error)
	CheckFidelity(ctx context.Context, current screenplay.Screenplay, bible story.StoryBible, analyses []analysis.ChapterAnalysis) (fidelity.FidelityResult, error)
	RepairFidelity(ctx context.Context, current screenplay.Screenplay, bible story.StoryBible, analyses []analysis.ChapterAnalysis, result fidelity.FidelityResult) (screenplay.Screenplay, error)
}
