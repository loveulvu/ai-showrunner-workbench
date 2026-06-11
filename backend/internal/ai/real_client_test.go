package ai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestRealClientCallChatContentSendsOpenAICompatibleRequest(t *testing.T) {
	var received chatCompletionRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/chat/completions" {
			t.Errorf("request path = %q, want /v1/chat/completions", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
			t.Errorf("Authorization = %q, want Bearer test-key", got)
		}
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", got)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode request: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"choices":[{"message":{"role":"assistant","content":"  generated text  "}}]}`))
	}))
	defer server.Close()

	client := NewRealClient(Config{
		APIKey:         "test-key",
		BaseURL:        server.URL + "/v1",
		Model:          defaultQwenModel,
		TimeoutSeconds: 2,
	})

	content, err := client.callChatContent(context.Background(), "test request", "hello")
	if err != nil {
		t.Fatalf("callChatContent() error = %v", err)
	}
	if content != "generated text" {
		t.Fatalf("content = %q, want generated text", content)
	}
	if received.Model != defaultQwenModel {
		t.Fatalf("model = %q, want %q", received.Model, defaultQwenModel)
	}
}

func TestRealClientCallChatContentErrors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantError  string
	}{
		{
			name:       "non-200 response",
			statusCode: http.StatusUnauthorized,
			body:       `{"error":{"message":"invalid key"}}`,
			wantError:  "chat completions returned status 401",
		},
		{
			name:       "non-200 success response",
			statusCode: http.StatusCreated,
			body:       `{"choices":[{"message":{"role":"assistant","content":"generated text"}}]}`,
			wantError:  "chat completions returned status 201",
		},
		{
			name:       "invalid JSON response",
			statusCode: http.StatusOK,
			body:       `not-json`,
			wantError:  "invalid JSON response from chat completions",
		},
		{
			name:       "empty choices",
			statusCode: http.StatusOK,
			body:       `{"choices":[]}`,
			wantError:  "chat completions response has empty choices",
		},
		{
			name:       "empty model output",
			statusCode: http.StatusOK,
			body:       `{"choices":[{"message":{"role":"assistant","content":"  "}}]}`,
			wantError:  "chat completions response has empty model output",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.body))
			}))
			defer server.Close()

			client := NewRealClient(Config{
				APIKey:         "test-key",
				BaseURL:        server.URL,
				Model:          defaultQwenModel,
				TimeoutSeconds: 2,
			})

			_, err := client.callChatContent(context.Background(), "test request", "hello")
			if err == nil {
				t.Fatal("callChatContent() error = nil, want error")
			}
			if !strings.Contains(err.Error(), tt.wantError) {
				t.Fatalf("callChatContent() error = %q, want substring %q", err, tt.wantError)
			}
		})
	}
}

func TestRealClientCallChatContentTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		_, _ = w.Write([]byte(`{"choices":[{"message":{"role":"assistant","content":"late output"}}]}`))
	}))
	defer server.Close()

	client := NewRealClient(Config{
		APIKey:         "test-key",
		BaseURL:        server.URL,
		Model:          defaultQwenModel,
		TimeoutSeconds: 1,
	})
	client.timeout = 20 * time.Millisecond
	client.httpClient.Timeout = client.timeout

	_, err := client.callChatContent(context.Background(), "test request", "hello")
	if err == nil {
		t.Fatal("callChatContent() error = nil, want timeout error")
	}
	if !strings.Contains(err.Error(), "timeout after") {
		t.Fatalf("callChatContent() error = %q, want timeout message", err)
	}
}
