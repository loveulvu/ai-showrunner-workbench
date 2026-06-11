package video

import (
	"net/http"
	"strings"
	"testing"
)

func TestProviderConfigFromEnvDefaults(t *testing.T) {
	for _, name := range []string{
		"VIDEO_PROVIDER",
		"VIDEO_MODEL",
		"VIDEO_BASE_URL",
		"VIDEO_API_KEY",
		"VIDEO_TIMEOUT_SECONDS",
		"VIDEO_POLL_INTERVAL_SECONDS",
	} {
		t.Setenv(name, "")
	}

	config, err := ProviderConfigFromEnv()
	if err != nil {
		t.Fatalf("ProviderConfigFromEnv() error = %v", err)
	}
	if config.Provider != "mock" {
		t.Fatalf("Provider = %q, want mock", config.Provider)
	}
	if config.Model != "wan2.6-t2v" {
		t.Fatalf("Model = %q, want wan2.6-t2v", config.Model)
	}
	if config.TimeoutSeconds != 600 || config.PollIntervalSeconds != 10 {
		t.Fatalf("timeouts = %d/%d, want 600/10", config.TimeoutSeconds, config.PollIntervalSeconds)
	}
	if config.APIKeySet {
		t.Fatal("APIKeySet = true, want false")
	}
}

func TestNewGeneratorFromConfigMock(t *testing.T) {
	generator, err := NewGeneratorFromConfig(
		ProviderConfig{Provider: "mock", Model: "mock-model"},
		NewMemoryVideoTaskStore(),
	)
	if err != nil {
		t.Fatalf("NewGeneratorFromConfig() error = %v", err)
	}
	if _, ok := generator.(*MockVideoGenerator); !ok {
		t.Fatalf("generator type = %T, want *MockVideoGenerator", generator)
	}
}

func TestNewGeneratorFromConfigWan(t *testing.T) {
	t.Setenv("VIDEO_API_KEY", "test-key")
	generator, err := NewGeneratorFromConfig(
		ProviderConfig{Provider: "wan", Model: "wan2.6-t2v", BaseURL: "https://example.com/api/v1", TimeoutSeconds: 10},
		NewMemoryVideoTaskStore(),
	)
	if err != nil {
		t.Fatalf("NewGeneratorFromConfig() error = %v", err)
	}
	if _, ok := generator.(*WanVideoGenerator); !ok {
		t.Fatalf("generator type = %T, want *WanVideoGenerator", generator)
	}
	wanGenerator := generator.(*WanVideoGenerator)
	transport, ok := wanGenerator.httpClient.Transport.(*http.Transport)
	if !ok || transport.Proxy == nil {
		t.Fatalf("Wan transport = %T, want proxy-aware *http.Transport", wanGenerator.httpClient.Transport)
	}
}

func TestNewGeneratorFromConfigWanRequiresConfiguration(t *testing.T) {
	t.Setenv("VIDEO_API_KEY", "")
	_, err := NewGeneratorFromConfig(ProviderConfig{Provider: "wan"}, NewMemoryVideoTaskStore())
	if err == nil || !strings.Contains(err.Error(), "VIDEO_BASE_URL") {
		t.Fatalf("error = %v, want VIDEO_BASE_URL requirement", err)
	}
}
