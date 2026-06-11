package ai

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestNewClientFromEnvQwenDefaults(t *testing.T) {
	t.Setenv("AI_PROVIDER", ProviderQwen)
	t.Setenv("AI_API_KEY", "test-key")
	t.Setenv("AI_BASE_URL", " ")
	t.Setenv("AI_MODEL", " ")
	t.Setenv("AI_TIMEOUT_SECONDS", "12")

	client, err := NewClientFromEnv()
	if err != nil {
		t.Fatalf("NewClientFromEnv() error = %v", err)
	}

	realClient, ok := client.(*RealClient)
	if !ok {
		t.Fatalf("NewClientFromEnv() type = %T, want *RealClient", client)
	}
	if realClient.endpoint != defaultQwenBaseURL+"/chat/completions" {
		t.Fatalf("endpoint = %q, want %q", realClient.endpoint, defaultQwenBaseURL+"/chat/completions")
	}
	if realClient.model != defaultQwenModel {
		t.Fatalf("model = %q, want %q", realClient.model, defaultQwenModel)
	}
}

func TestNewClientFromEnvQwenRequiresAPIKey(t *testing.T) {
	t.Setenv("AI_PROVIDER", ProviderQwen)
	t.Setenv("AI_API_KEY", " ")
	t.Setenv("AI_BASE_URL", " ")
	t.Setenv("AI_MODEL", " ")

	_, err := NewClientFromEnv()
	if err == nil {
		t.Fatal("NewClientFromEnv() error = nil, want missing API key error")
	}
	if !strings.Contains(err.Error(), "AI_PROVIDER=qwen requires AI_API_KEY") {
		t.Fatalf("NewClientFromEnv() error = %q", err)
	}
}

func TestRuntimeStatusFromEnvUsesQwenModelDefault(t *testing.T) {
	t.Setenv("AI_PROVIDER", ProviderQwen)
	t.Setenv("AI_API_KEY", "test-key")
	t.Setenv("AI_BASE_URL", " ")
	t.Setenv("AI_MODEL", " ")

	status := RuntimeStatusFromEnv()
	if status.AIProvider != ProviderQwen {
		t.Fatalf("AIProvider = %q, want %q", status.AIProvider, ProviderQwen)
	}
	if status.AIModel != defaultQwenModel {
		t.Fatalf("AIModel = %q, want %q", status.AIModel, defaultQwenModel)
	}
}

func TestNewClientFromEnvKeepsMockProvider(t *testing.T) {
	t.Setenv("AI_PROVIDER", ProviderMock)
	t.Setenv("AI_API_KEY", " ")
	t.Setenv("AI_BASE_URL", " ")
	t.Setenv("AI_MODEL", " ")

	client, err := NewClientFromEnv()
	if err != nil {
		t.Fatalf("NewClientFromEnv() error = %v", err)
	}
	if _, ok := client.(MockClient); !ok {
		t.Fatalf("NewClientFromEnv() type = %T, want MockClient", client)
	}
}

func TestDotenvCandidates(t *testing.T) {
	want := []string{".env", "../.env"}
	if got := dotenvCandidates(); !reflect.DeepEqual(got, want) {
		t.Fatalf("dotenvCandidates() = %#v, want %#v", got, want)
	}
}

func TestLoadDotEnvIfPresent(t *testing.T) {
	keyFromFile := "NOVEL_TO_SCREENPLAY_TEST_FROM_FILE"
	keyAlreadySet := "NOVEL_TO_SCREENPLAY_TEST_ALREADY_SET"
	t.Setenv(keyFromFile, "")
	t.Setenv(keyAlreadySet, "process-value")

	path := filepath.Join(t.TempDir(), ".env")
	content := keyFromFile + "=file-value\n" + keyAlreadySet + "=file-value\n"
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write test .env: %v", err)
	}

	if loaded := loadDotEnvIfPresent(path); !loaded {
		t.Fatal("loadDotEnvIfPresent() = false, want true")
	}
	if got := os.Getenv(keyFromFile); got != "file-value" {
		t.Fatalf("%s = %q, want file-value", keyFromFile, got)
	}
	if got := os.Getenv(keyAlreadySet); got != "process-value" {
		t.Fatalf("%s = %q, want process-value", keyAlreadySet, got)
	}
}

func TestLoadDotEnvIfPresentMissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.env")
	if loaded := loadDotEnvIfPresent(path); loaded {
		t.Fatal("loadDotEnvIfPresent() = true, want false")
	}
}

func TestRuntimeStatusFromEnvReportsBaseURLAndProxyStatus(t *testing.T) {
	t.Setenv("AI_PROVIDER", ProviderQwen)
	t.Setenv("AI_BASE_URL", "https://example.com/v1")
	t.Setenv("HTTP_PROXY", "http://proxy.example")
	t.Setenv("HTTPS_PROXY", "")

	status := RuntimeStatusFromEnv()
	if status.AIBaseURL != "https://example.com/v1" {
		t.Fatalf("AIBaseURL = %q", status.AIBaseURL)
	}
	if !status.HTTPProxyConfigured {
		t.Fatal("HTTPProxyConfigured = false, want true")
	}
	if status.HTTPSProxyConfigured {
		t.Fatal("HTTPSProxyConfigured = true, want false")
	}
}
