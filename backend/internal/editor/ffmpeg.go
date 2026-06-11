package editor

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func RunFFmpeg(ctx context.Context, plan EditingPlan, concatFile string) (EditResult, error) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return EditResult{}, fmt.Errorf("ffmpeg not found")
	}
	if concatFile == "" {
		return EditResult{}, fmt.Errorf("concat file is required")
	}
	if err := os.MkdirAll(filepath.Dir(plan.OutputFile), 0o755); err != nil {
		return EditResult{}, fmt.Errorf("create output directory: %w", err)
	}

	command := exec.CommandContext(ctx, "ffmpeg",
		"-y",
		"-f", "concat",
		"-safe", "0",
		"-i", concatFile,
		"-c", "copy",
		plan.OutputFile,
	)
	output, err := command.CombinedOutput()
	if err != nil {
		return EditResult{}, fmt.Errorf("ffmpeg concat failed: %w: %s", err, summarizeCommandOutput(output))
	}

	return EditResult{
		OutputFile: plan.OutputFile,
		ClipCount:  len(plan.Clips),
		Warnings:   []string{},
	}, nil
}

func summarizeCommandOutput(value []byte) string {
	const limit = 500
	if len(value) <= limit {
		return string(value)
	}
	return string(value[:limit]) + "..."
}
