package main

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestVideoCheckDefaultOnlyQueriesFakeTask(t *testing.T) {
	var methods []string
	var paths []string
	var mu sync.Mutex
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		methods = append(methods, r.Method)
		paths = append(paths, r.URL.Path)
		mu.Unlock()
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"code":"TaskNotFound","message":"task not found"}`))
	}))
	defer server.Close()

	setVideoCheckEnv(t, server.URL)
	var output bytes.Buffer
	if err := run(nil, log.New(&output, "", 0), server.Client()); err != nil {
		t.Fatalf("run() error = %v", err)
	}

	if len(methods) != 1 || methods[0] != http.MethodGet {
		t.Fatalf("methods = %v, want one GET", methods)
	}
	if paths[0] != "/tasks/"+fakeTaskID {
		t.Fatalf("path = %q", paths[0])
	}
}

func TestVideoCheckPostsOnlyWithCreateFlag(t *testing.T) {
	var methods []string
	var mu sync.Mutex
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		methods = append(methods, r.Method)
		mu.Unlock()
		if r.Method == http.MethodPost {
			_, _ = w.Write([]byte(`{"output":{"task_id":"test-task","task_status":"PENDING"}}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"code":"TaskNotFound","message":"task not found"}`))
	}))
	defer server.Close()

	setVideoCheckEnv(t, server.URL)
	var output bytes.Buffer
	if err := run([]string{"--create"}, log.New(&output, "", 0), server.Client()); err != nil {
		t.Fatalf("run() error = %v", err)
	}

	if len(methods) != 2 || methods[0] != http.MethodPost || methods[1] != http.MethodGet {
		t.Fatalf("methods = %v, want POST then GET", methods)
	}
	if !bytes.Contains(output.Bytes(), []byte("This will create a real Wan video task and may consume credits.")) {
		t.Fatalf("output missing create warning: %s", output.String())
	}
}

func setVideoCheckEnv(t *testing.T, baseURL string) {
	t.Helper()
	t.Setenv("VIDEO_PROVIDER", "wan")
	t.Setenv("VIDEO_BASE_URL", baseURL)
	t.Setenv("VIDEO_API_KEY", "test-key")
	t.Setenv("AI_API_KEY", "")
	t.Setenv("HTTP_PROXY", "")
	t.Setenv("HTTPS_PROXY", "")
	t.Setenv("NO_PROXY", "")
}
