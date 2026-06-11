package ai

import (
	"context"

	"novel-to-screenplay-ai/internal/analysis"
	"novel-to-screenplay-ai/internal/fidelity"
	"novel-to-screenplay-ai/internal/novel"
	"novel-to-screenplay-ai/internal/screenplay"
	"novel-to-screenplay-ai/internal/story"
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
	CheckFidelity(ctx context.Context, current screenplay.Screenplay, bible story.StoryBible, analyses []analysis.ChapterAnalysis) (fidelity.FidelityResult, error)
	RepairFidelity(ctx context.Context, current screenplay.Screenplay, bible story.StoryBible, analyses []analysis.ChapterAnalysis, result fidelity.FidelityResult) (screenplay.Screenplay, error)
}
