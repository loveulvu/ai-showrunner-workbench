package editor

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func PrepareLocalPaths(plan EditingPlan) (EditingPlan, error) {
	outputDir, err := planOutputDir(plan)
	if err != nil {
		return plan, err
	}
	for index := range plan.Clips {
		if strings.TrimSpace(plan.Clips[index].LocalPath) != "" {
			continue
		}
		name, err := ClipFileName(plan.Clips[index].ShotID)
		if err != nil {
			return plan, err
		}
		plan.Clips[index].LocalPath = filepath.Join(outputDir, "clips", name)
	}
	return plan, nil
}

func BuildConcatList(plan EditingPlan) (string, error) {
	outputDir, err := planOutputDir(plan)
	if err != nil {
		return "", err
	}
	editDir := filepath.Join(outputDir, "edit")
	if err := os.MkdirAll(editDir, 0o755); err != nil {
		return "", fmt.Errorf("create edit directory: %w", err)
	}

	var lines strings.Builder
	for _, clip := range plan.Clips {
		if strings.TrimSpace(clip.LocalPath) == "" {
			return "", fmt.Errorf("clip %q local_path is required", clip.ShotID)
		}
		absolute, err := filepath.Abs(clip.LocalPath)
		if err != nil {
			return "", fmt.Errorf("resolve clip %q path: %w", clip.ShotID, err)
		}
		lines.WriteString("file '")
		lines.WriteString(escapeConcatPath(filepath.ToSlash(absolute)))
		lines.WriteString("'\n")
	}

	path := filepath.Join(editDir, "concat.txt")
	if err := os.WriteFile(path, []byte(lines.String()), 0o600); err != nil {
		return "", fmt.Errorf("write concat list: %w", err)
	}
	return path, nil
}

func BuildSRT(plan EditingPlan) (string, error) {
	outputDir, err := planOutputDir(plan)
	if err != nil {
		return "", err
	}

	var content strings.Builder
	start := time.Duration(0)
	index := 1
	for _, clip := range plan.Clips {
		duration := clip.DurationSeconds
		if duration <= 0 {
			duration = 5
		}
		end := start + time.Duration(duration)*time.Second
		if strings.TrimSpace(clip.Subtitle) != "" {
			content.WriteString(strconv.Itoa(index))
			content.WriteString("\n")
			content.WriteString(formatSRTTime(start))
			content.WriteString(" --> ")
			content.WriteString(formatSRTTime(end))
			content.WriteString("\n")
			content.WriteString(strings.TrimSpace(clip.Subtitle))
			content.WriteString("\n\n")
			index++
		}
		start = end
	}
	if index == 1 {
		return "", nil
	}

	editDir := filepath.Join(outputDir, "edit")
	if err := os.MkdirAll(editDir, 0o755); err != nil {
		return "", fmt.Errorf("create edit directory: %w", err)
	}
	path := filepath.Join(editDir, "subtitles.srt")
	if err := os.WriteFile(path, []byte(content.String()), 0o600); err != nil {
		return "", fmt.Errorf("write subtitles: %w", err)
	}
	return path, nil
}

func planOutputDir(plan EditingPlan) (string, error) {
	if strings.TrimSpace(plan.OutputFile) == "" {
		return "", fmt.Errorf("output_file is required")
	}
	return filepath.Dir(plan.OutputFile), nil
}

func escapeConcatPath(value string) string {
	return strings.ReplaceAll(value, "'", "'\\''")
}

func formatSRTTime(value time.Duration) string {
	totalMilliseconds := value.Milliseconds()
	hours := totalMilliseconds / 3_600_000
	minutes := (totalMilliseconds % 3_600_000) / 60_000
	seconds := (totalMilliseconds % 60_000) / 1_000
	milliseconds := totalMilliseconds % 1_000
	return fmt.Sprintf("%02d:%02d:%02d,%03d", hours, minutes, seconds, milliseconds)
}
