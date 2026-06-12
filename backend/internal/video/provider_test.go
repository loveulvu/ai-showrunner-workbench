package video

import (
	"bytes"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"
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
	if wanGenerator.httpClient.Timeout != 10*time.Second {
		t.Fatalf("Wan timeout = %s, want 10s", wanGenerator.httpClient.Timeout)
	}
}

func TestNewGeneratorFromConfigWanRequiresConfiguration(t *testing.T) {
	t.Setenv("VIDEO_API_KEY", "")
	t.Setenv("AI_API_KEY", "")
	_, err := NewGeneratorFromConfig(ProviderConfig{Provider: "wan"}, NewMemoryVideoTaskStore())
	if err == nil || !strings.Contains(err.Error(), "VIDEO_BASE_URL") {
		t.Fatalf("error = %v, want VIDEO_BASE_URL requirement", err)
	}
}

func TestEffectiveAPIKeyFallsBackToAIAPIKey(t *testing.T) {
	t.Setenv("VIDEO_API_KEY", "")
	t.Setenv("AI_API_KEY", "fallback-key")

	key, fallbackUsed := EffectiveAPIKey()
	if key != "fallback-key" || !fallbackUsed {
		t.Fatalf("EffectiveAPIKey() = %q/%t, want fallback-key/true", key, fallbackUsed)
	}

	config, err := ProviderConfigFromEnv()
	if err != nil {
		t.Fatalf("ProviderConfigFromEnv() error = %v", err)
	}
	if config.APIKeySet || !config.AIAPIKeyFallbackUsed {
		t.Fatalf("config key status = %#v", config)
	}
}

func TestProviderConfigFromEnvUsesVideoTimeoutSeconds(t *testing.T) {
	t.Setenv("VIDEO_TIMEOUT_SECONDS", "725")

	config, err := ProviderConfigFromEnv()
	if err != nil {
		t.Fatalf("ProviderConfigFromEnv() error = %v", err)
	}
	if config.TimeoutSeconds != 725 {
		t.Fatalf("TimeoutSeconds = %d, want 725", config.TimeoutSeconds)
	}
}

func TestLogProviderConfigIncludesTimeoutWithoutSecrets(t *testing.T) {
	var output bytes.Buffer
	LogProviderConfig(log.New(&output, "", 0), ProviderConfig{
		Provider:       "wan",
		Model:          "wan2.6-t2v",
		BaseURL:        "https://secret.example/api",
		APIKeySet:      true,
		TimeoutSeconds: 600,
	})

	logged := output.String()
	for _, expected := range []string{"VIDEO_PROVIDER=wan", "VIDEO_MODEL=wan2.6-t2v", "VIDEO_BASE_URL set: true", "VIDEO_API_KEY set: true", "VIDEO_TIMEOUT_SECONDS=600"} {
		if !strings.Contains(logged, expected) {
			t.Fatalf("log missing %q: %s", expected, logged)
		}
	}
	if strings.Contains(logged, "secret.example") {
		t.Fatalf("log contains sensitive URL: %s", logged)
	}
}
