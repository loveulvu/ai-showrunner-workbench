package video

import (
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

func TestNewGeneratorFromConfigWanNotImplemented(t *testing.T) {
	_, err := NewGeneratorFromConfig(
		ProviderConfig{Provider: "wan", Model: "wan2.6-t2v"},
		NewMemoryVideoTaskStore(),
	)
	if err == nil {
		t.Fatal("NewGeneratorFromConfig() error = nil, want not implemented error")
	}
	if !strings.Contains(err.Error(), "wan provider not implemented yet") {
		t.Fatalf("error = %q", err)
	}
}
