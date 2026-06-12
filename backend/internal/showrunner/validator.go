package showrunner

import (
	"fmt"
	"strings"
)

type ValidationResult struct {
	Passed   bool     `json:"passed"`
	Errors   []string `json:"errors"`
	Warnings []string `json:"warnings"`
}

func Validate(result ShowrunnerResult) ValidationResult {
	validation := ValidationResult{
		Passed:   true,
		Errors:   []string{},
		Warnings: append([]string(nil), result.Warnings...),
	}

	if len(result.Characters) == 0 {
		validation.Warnings = append(validation.Warnings, "characters is empty")
	}
	if len(result.Scenes) == 0 {
		validation.Warnings = append(validation.Warnings, "scenes is empty")
	}
	if len(result.Chapters) == 0 {
		validation.Warnings = append(validation.Warnings, "chapters is empty")
	}
	if len(result.Shots) == 0 {
		validation.Errors = append(validation.Errors, "shots must contain at least one shot")
	}

	for index, shot := range result.Shots {
		prefix := fmt.Sprintf("shots[%d]", index)
		if strings.TrimSpace(shot.ID) == "" {
			validation.Warnings = append(validation.Warnings, prefix+".id is empty")
		}
		if shot.ChapterNumber <= 0 {
			validation.Warnings = append(validation.Warnings, prefix+".chapter_number is not positive")
		}
		if strings.TrimSpace(shot.Action) == "" && strings.TrimSpace(shot.Dialogue.Text()) == "" {
			validation.Warnings = append(validation.Warnings, prefix+" has no action, dialogue, or subtitle")
		}
		if strings.TrimSpace(shot.ImagePrompt) == "" {
			validation.Warnings = append(validation.Warnings, prefix+".image_prompt is empty")
		}
		if strings.TrimSpace(shot.VideoPrompt) == "" {
			validation.Warnings = append(validation.Warnings, prefix+".video_prompt is empty")
		}
		if strings.TrimSpace(shot.NegativePrompt) == "" {
			validation.Warnings = append(validation.Warnings, prefix+".negative_prompt is empty")
		}
		if strings.TrimSpace(shot.AudioPrompt) == "" {
			validation.Warnings = append(validation.Warnings, prefix+".audio_prompt is empty")
		}
	}

	validation.Errors = uniqueStrings(validation.Errors)
	validation.Warnings = uniqueStrings(validation.Warnings)
	validation.Passed = len(validation.Errors) == 0
	return validation
}

func uniqueStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}
