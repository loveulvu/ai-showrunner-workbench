package main

import (
	"testing"

	"ai-showrunner-workbench/internal/ai"
)

func TestSafeURLForLog(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{
			name:  "normal URL",
			value: "https://dashscope-intl.aliyuncs.com/compatible-mode/v1",
			want:  "https://dashscope-intl.aliyuncs.com/compatible-mode/v1",
		},
		{
			name:  "redacts credentials and query",
			value: "https://user:secret@example.com/v1?api_key=secret#fragment",
			want:  "https://REDACTED@example.com/v1",
		},
		{
			name:  "invalid URL",
			value: "not-a-url",
			want:  "<invalid URL configured>",
		},
		{
			name:  "empty URL",
			value: " ",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ai.SafeURLForLog(tt.value); got != tt.want {
				t.Fatalf("SafeURLForLog(%q) = %q, want %q", tt.value, got, tt.want)
			}
		})
	}
}
