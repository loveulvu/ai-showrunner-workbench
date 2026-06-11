package analysis

import (
	"context"

	"ai-showrunner-workbench/internal/novel"
)

type ChapterAIClient interface {
	AnalyzeChapter(ctx context.Context, chapter novel.Chapter) (ChapterAnalysis, error)
}

type Analyzer struct {
	client ChapterAIClient
}

func NewAnalyzer(client ChapterAIClient) Analyzer {
	return Analyzer{client: client}
}

func (a Analyzer) AnalyzeChapters(ctx context.Context, chapters []novel.Chapter) ([]ChapterAnalysis, error) {
	results := make([]ChapterAnalysis, 0, len(chapters))
	for _, chapter := range chapters {
		analysis, err := a.client.AnalyzeChapter(ctx, chapter)
		if err != nil {
			return nil, err
		}
		results = append(results, analysis)
	}
	return results, nil
}
