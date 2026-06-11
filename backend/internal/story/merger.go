package story

import (
	"context"

	"novel-to-screenplay-ai/internal/analysis"
)

type StoryAIClient interface {
	MergeStoryBible(ctx context.Context, analyses []analysis.ChapterAnalysis) (StoryBible, error)
}

type Merger struct {
	client StoryAIClient
}

func NewMerger(client StoryAIClient) Merger {
	return Merger{client: client}
}

func (m Merger) Merge(ctx context.Context, analyses []analysis.ChapterAnalysis) (StoryBible, error) {
	return m.client.MergeStoryBible(ctx, analyses)
}
