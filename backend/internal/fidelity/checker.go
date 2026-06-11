package fidelity

import (
	"context"
	"strings"

	"ai-showrunner-workbench/internal/analysis"
	"ai-showrunner-workbench/internal/screenplay"
	"ai-showrunner-workbench/internal/story"
)

type AIClient interface {
	CheckFidelity(ctx context.Context, screenplay screenplay.Screenplay, bible story.StoryBible, analyses []analysis.ChapterAnalysis) (FidelityResult, error)
	RepairFidelity(ctx context.Context, current screenplay.Screenplay, bible story.StoryBible, analyses []analysis.ChapterAnalysis, result FidelityResult) (screenplay.Screenplay, error)
}

type Checker struct {
	client AIClient
}

func NewChecker(client AIClient) Checker {
	return Checker{client: client}
}

func (c Checker) CheckAndRepairOnce(ctx context.Context, current screenplay.Screenplay, bible story.StoryBible, analyses []analysis.ChapterAnalysis) (screenplay.Screenplay, FidelityResult, error) {
	result, err := c.client.CheckFidelity(ctx, current, bible, analyses)
	if err != nil {
		return current, result, err
	}
	if !needsRepair(result) {
		return current, result, nil
	}

	repaired, err := c.client.RepairFidelity(ctx, current, bible, analyses, result)
	if err != nil {
		return current, result, nil
	}
	repaired.SourceChapters = current.SourceChapters
	if validation := screenplay.Validate(repaired); !validation.Passed {
		return current, result, nil
	}

	recheck, err := c.client.CheckFidelity(ctx, repaired, bible, analyses)
	if err != nil {
		return repaired, result, nil
	}
	return repaired, recheck, nil
}

func needsRepair(result FidelityResult) bool {
	for _, issue := range result.Issues {
		severity := strings.ToLower(strings.TrimSpace(issue.Severity))
		if severity == "high" || severity == "medium" {
			return true
		}
	}
	return false
}
