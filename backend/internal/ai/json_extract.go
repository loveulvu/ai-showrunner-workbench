package ai

import "strings"

const repairPromptLimit = 4000

func extractJSON(raw string) string {
	text := strings.TrimSpace(raw)
	if !strings.HasPrefix(text, "```") {
		return text
	}

	firstNewline := strings.IndexByte(text, '\n')
	if firstNewline < 0 {
		return text
	}

	rest := text[firstNewline+1:]
	end := strings.LastIndex(rest, "```")
	if end < 0 {
		return strings.TrimSpace(rest)
	}

	return strings.TrimSpace(rest[:end])
}

func truncateText(value string, limit int) string {
	text := strings.TrimSpace(value)
	if limit <= 0 || len(text) <= limit {
		return text
	}
	return text[:limit] + "..."
}
