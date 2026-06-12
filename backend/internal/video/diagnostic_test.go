package video

import (
	"errors"
	"net/http"
	"strings"
	"testing"
)

func TestClassifyTaskDiagnostic(t *testing.T) {
	tests := []struct {
		name   string
		result TaskDiagnostic
		err    error
		want   string
	}{
		{name: "not found", result: TaskDiagnostic{HTTPStatus: http.StatusNotFound, Code: "TaskNotFound"}, want: "reachable"},
		{name: "no key", result: TaskDiagnostic{HTTPStatus: http.StatusUnauthorized, Message: "No API-key provided"}, want: "No API key"},
		{name: "invalid key", result: TaskDiagnostic{HTTPStatus: http.StatusUnauthorized, Code: "InvalidApiKey"}, want: "does not match"},
		{name: "timeout", err: errors.New("request timeout"), want: "unstable"},
		{name: "upstream", result: TaskDiagnostic{HTTPStatus: http.StatusBadGateway}, want: "server error"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := ClassifyTaskDiagnostic(test.result, test.err)
			if !strings.Contains(got, test.want) {
				t.Fatalf("ClassifyTaskDiagnostic() = %q, want substring %q", got, test.want)
			}
		})
	}
}
