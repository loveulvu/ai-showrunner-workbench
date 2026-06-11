package screenplay

import (
	"context"

	"ai-showrunner-workbench/internal/analysis"
	"ai-showrunner-workbench/internal/story"
)

type ScreenplayAIClient interface {
	GenerateScreenplay(ctx context.Context, bible story.StoryBible, analyses []analysis.ChapterAnalysis) (Screenplay, error)
}

type Generator struct {
	client ScreenplayAIClient
}

func NewGenerator(client ScreenplayAIClient) Generator {
	return Generator{client: client}
}

func (g Generator) Generate(ctx context.Context, bible story.StoryBible, analyses []analysis.ChapterAnalysis) (Screenplay, error) {
	return g.client.GenerateScreenplay(ctx, bible, analyses)
}
