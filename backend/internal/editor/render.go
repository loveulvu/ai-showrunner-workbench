package editor

import (
	"context"
	"fmt"
)

func Render(ctx context.Context, plan EditingPlan) (EditResult, error) {
	if len(plan.Clips) == 0 {
		return EditResult{}, fmt.Errorf("at least one clip is required")
	}

	prepared, err := PrepareLocalPaths(plan)
	if err != nil {
		return EditResult{}, err
	}
	for index, clip := range prepared.Clips {
		downloaded, err := DownloadClip(ctx, clip)
		if err != nil {
			return EditResult{}, err
		}
		prepared.Clips[index] = downloaded
	}

	concatFile, err := BuildConcatList(prepared)
	if err != nil {
		return EditResult{}, err
	}
	subtitlesFile, err := BuildSRT(prepared)
	if err != nil {
		return EditResult{}, err
	}
	result, err := RunFFmpeg(ctx, prepared, concatFile)
	if err != nil {
		return EditResult{}, err
	}
	result.SubtitlesFile = subtitlesFile
	return result, nil
}
