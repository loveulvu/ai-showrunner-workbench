package ai

import (
	"bytes"
	"errors"
	"log"
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
	t.Setenv("DASHSCOPE_API_KEY", " ")
	t.Setenv("AI_BASE_URL", " ")
	t.Setenv("AI_MODEL", " ")

	_, err := NewClientFromEnv()
	if err == nil {
		t.Fatal("NewClientFromEnv() error = nil, want missing API key error")
	}
	if !strings.Contains(err.Error(), "AI_PROVIDER=qwen requires AI_API_KEY or DASHSCOPE_API_KEY") {
		t.Fatalf("NewClientFromEnv() error = %q", err)
	}
}

func TestConfigFromEnvPrefersAIAPIKeyAndFallsBackToDashScope(t *testing.T) {
	t.Setenv("AI_API_KEY", "preferred-key")
	t.Setenv("DASHSCOPE_API_KEY", "fallback-key")
	if got := configFromEnv(ProviderQwen, 10).APIKey; got != "preferred-key" {
		t.Fatalf("APIKey = %q, want preferred-key", got)
	}

	t.Setenv("AI_API_KEY", " ")
	if got := configFromEnv(ProviderQwen, 10).APIKey; got != "fallback-key" {
		t.Fatalf("APIKey = %q, want fallback-key", got)
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
	t.Setenv("HTTPS_PROXY", "https://proxy.example")
	t.Setenv("NO_PROXY", "localhost,127.0.0.1")

	status := RuntimeStatusFromEnv()
	if status.AIBaseURL != "https://example.com/v1" {
		t.Fatalf("AIBaseURL = %q", status.AIBaseURL)
	}
	if !status.HTTPProxyConfigured {
		t.Fatal("HTTPProxyConfigured = false, want true")
	}
	if !status.HTTPSProxyConfigured {
		t.Fatal("HTTPSProxyConfigured = false, want true")
	}
	if !status.NOProxyConfigured {
		t.Fatal("NOProxyConfigured = false, want true")
	}
}

func TestLogRuntimeConfigurationDoesNotPrintSecrets(t *testing.T) {
	const (
		apiKey     = "secret-api-key"
		dashKey    = "secret-dashscope-key"
		httpProxy  = "http://user:secret@proxy.example:8080"
		httpsProxy = "https://user:secret@proxy.example:8443"
	)
	t.Setenv("AI_PROVIDER", ProviderQwen)
	t.Setenv("AI_API_KEY", apiKey)
	t.Setenv("DASHSCOPE_API_KEY", dashKey)
	t.Setenv("AI_BASE_URL", "https://user:secret@example.com/v1?token=secret")
	t.Setenv("HTTP_PROXY", httpProxy)
	t.Setenv("HTTPS_PROXY", httpsProxy)
	t.Setenv("NO_PROXY", "localhost,127.0.0.1")

	var output bytes.Buffer
	LogRuntimeConfiguration(log.New(&output, "", 0))
	logged := output.String()

	for _, secret := range []string{apiKey, dashKey, httpProxy, httpsProxy, "user:secret", "token=secret"} {
		if strings.Contains(logged, secret) {
			t.Fatalf("runtime log contains secret value %q", secret)
		}
	}
	for _, expected := range []string{
		"AI_API_KEY set: true",
		"DASHSCOPE_API_KEY set: true",
		"HTTP_PROXY set: true",
		"HTTPS_PROXY set: true",
		"NO_PROXY set: true",
	} {
		if !strings.Contains(logged, expected) {
			t.Fatalf("runtime log missing %q", expected)
		}
	}
}

func TestSafeURLForLogAndRedactedDiagnostic(t *testing.T) {
	const secretKey = "secret-api-key"
	const secretProxy = "http://user:secret@proxy.example:8080"
	t.Setenv("AI_API_KEY", secretKey)
	t.Setenv("HTTP_PROXY", secretProxy)

	if got := SafeURLForLog("https://user:secret@example.com/v1?token=secret#fragment"); got != "https://REDACTED@example.com/v1" {
		t.Fatalf("SafeURLForLog() = %q", got)
	}

	diagnostic := RedactedDiagnostic(errors.New("request failed for https://user:secret@example.com/v1?token=secret using " + secretKey + " via " + secretProxy))
	for _, secret := range []string{secretKey, secretProxy, "proxy.example:8080", "user:secret", "token=secret"} {
		if strings.Contains(diagnostic, secret) {
			t.Fatalf("RedactedDiagnostic() contains secret value %q", secret)
		}
	}
}
