package ai

import "testing"

func TestExtractJSONPlain(t *testing.T) {
	got := extractJSON(`{"title":"x"}`)
	if got != `{"title":"x"}` {
		t.Fatalf("unexpected JSON: %q", got)
	}
}

func TestExtractJSONFromJSONCodeBlock(t *testing.T) {
	got := extractJSON("```json\n{\"title\":\"x\"}\n```")
	if got != `{"title":"x"}` {
		t.Fatalf("unexpected JSON: %q", got)
	}
}

func TestExtractJSONFromCodeBlock(t *testing.T) {
	got := extractJSON("```\n{\"title\":\"x\"}\n```")
	if got != `{"title":"x"}` {
		t.Fatalf("unexpected JSON: %q", got)
	}
}

func TestExtractJSONTrimsWhitespace(t *testing.T) {
	got := extractJSON(" \n\t{\"title\":\"x\"}\n ")
	if got != `{"title":"x"}` {
		t.Fatalf("unexpected JSON: %q", got)
	}
}

func TestTruncateTextShortValue(t *testing.T) {
	got := truncateText("  abc  ", 10)
	if got != "abc" {
		t.Fatalf("unexpected truncated value: %q", got)
	}
}

func TestTruncateTextLongValue(t *testing.T) {
	got := truncateText("abcdef", 3)
	if got != "abc..." {
		t.Fatalf("unexpected truncated value: %q", got)
	}
}

func TestTruncateTextZeroLimit(t *testing.T) {
	got := truncateText("abcdef", 0)
	if got != "abcdef" {
		t.Fatalf("unexpected truncated value: %q", got)
	}
}
