package ai

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

const defaultAITimeoutSeconds = 180

var (
	loadEnvOnce    sync.Once
	loadedEnvFiles []string
)

type RuntimeStatus struct {
	AIProvider                string `json:"ai_provider"`
	AIModel                   string `json:"ai_model"`
	AIBaseURL                 string `json:"ai_base_url"`
	AIBaseURLConfigured       bool   `json:"ai_base_url_configured"`
	AIAPIKeyConfigured        bool   `json:"ai_api_key_configured"`
	DashScopeAPIKeyConfigured bool   `json:"dashscope_api_key_configured"`
	EffectiveAPIKeyConfigured bool   `json:"effective_api_key_configured"`
	HTTPProxyConfigured       bool   `json:"http_proxy_configured"`
	HTTPSProxyConfigured      bool   `json:"https_proxy_configured"`
	NOProxyConfigured         bool   `json:"no_proxy_configured"`
	AITimeoutSeconds          int    `json:"ai_timeout_seconds"`
}

type EnvFileStatus struct {
	Path   string
	Loaded bool
}

func LoadEnv() []string {
	loadEnvOnce.Do(func() {
		for _, path := range dotenvCandidates() {
			if loadDotEnvIfPresent(path) {
				loadedEnvFiles = append(loadedEnvFiles, path)
			}
		}
	})
	return append([]string(nil), loadedEnvFiles...)
}

func EnvFileStatuses() []EnvFileStatus {
	loaded := LoadEnv()
	statuses := make([]EnvFileStatus, 0, len(dotenvCandidates()))
	for _, path := range dotenvCandidates() {
		statuses = append(statuses, EnvFileStatus{
			Path:   path,
			Loaded: containsString(loaded, path),
		})
	}
	return statuses
}

func NewClientFromEnv() (Client, error) {
	LoadEnv()

	provider := normalizedProvider(os.Getenv("AI_PROVIDER"))
	if provider == ProviderMock {
		return NewMockClient(), nil
	}

	if provider != ProviderReal && provider != ProviderQwen {
		return nil, fmt.Errorf("unknown AI_PROVIDER %q; expected mock, real, or qwen", provider)
	}

	timeoutSeconds, err := aiTimeoutSecondsFromEnv()
	if err != nil {
		return nil, err
	}

	cfg := configFromEnv(provider, timeoutSeconds)

	var missing []string
	if cfg.APIKey == "" {
		missing = append(missing, "AI_API_KEY or DASHSCOPE_API_KEY")
	}
	if cfg.BaseURL == "" {
		missing = append(missing, "AI_BASE_URL")
	}
	if cfg.Model == "" {
		missing = append(missing, "AI_MODEL")
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("AI_PROVIDER=%s requires %s", provider, strings.Join(missing, ", "))
	}

	return NewRealClient(cfg), nil
}

func CheckConnectivityFromEnv(ctx context.Context) error {
	client, err := NewClientFromEnv()
	if err != nil {
		return err
	}

	realClient, ok := client.(*RealClient)
	if !ok {
		return fmt.Errorf("AI_PROVIDER must be real or qwen for connectivity checks")
	}

	_, err = realClient.callChatContent(ctx, "connectivity check", "Reply with exactly OK.")
	return err
}

func RuntimeStatusFromEnv() RuntimeStatus {
	LoadEnv()
	timeoutSeconds, err := aiTimeoutSecondsFromEnv()
	if err != nil {
		timeoutSeconds = defaultAITimeoutSeconds
	}

	provider := normalizedProvider(os.Getenv("AI_PROVIDER"))
	cfg := configFromEnv(provider, timeoutSeconds)

	return RuntimeStatus{
		AIProvider:                provider,
		AIModel:                   cfg.Model,
		AIBaseURL:                 cfg.BaseURL,
		AIBaseURLConfigured:       envSet("AI_BASE_URL"),
		AIAPIKeyConfigured:        envSet("AI_API_KEY"),
		DashScopeAPIKeyConfigured: envSet("DASHSCOPE_API_KEY"),
		EffectiveAPIKeyConfigured: cfg.APIKey != "",
		HTTPProxyConfigured:       envSet("HTTP_PROXY"),
		HTTPSProxyConfigured:      envSet("HTTPS_PROXY"),
		NOProxyConfigured:         envSet("NO_PROXY"),
		AITimeoutSeconds:          timeoutSeconds,
	}
}

func LogRuntimeConfiguration(logger *log.Logger) {
	for _, status := range EnvFileStatuses() {
		logger.Printf("Environment file %s loaded: %t", status.Path, status.Loaded)
	}

	status := RuntimeStatusFromEnv()
	logger.Printf("AI_PROVIDER=%s", status.AIProvider)
	logger.Printf("AI_MODEL=%s", status.AIModel)
	logger.Printf("AI_BASE_URL=%s", SafeURLForLog(status.AIBaseURL))
	logger.Printf("AI_API_KEY set: %t", status.AIAPIKeyConfigured)
	logger.Printf("DASHSCOPE_API_KEY set: %t", status.DashScopeAPIKeyConfigured)
	logger.Printf("Effective API key set: %t", status.EffectiveAPIKeyConfigured)
	logger.Printf("HTTP_PROXY set: %t", status.HTTPProxyConfigured)
	logger.Printf("HTTPS_PROXY set: %t", status.HTTPSProxyConfigured)
	logger.Printf("NO_PROXY set: %t", status.NOProxyConfigured)
	logger.Printf("AI_TIMEOUT_SECONDS=%d", status.AITimeoutSeconds)
}

func SafeURLForLog(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}

	parsed, err := url.Parse(value)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return "<invalid URL configured>"
	}

	if parsed.User != nil {
		parsed.User = url.User("REDACTED")
	}
	parsed.RawQuery = ""
	parsed.Fragment = ""
	return parsed.String()
}

var urlPattern = regexp.MustCompile(`https?://[^\s"']+`)

func RedactedDiagnostic(err error) string {
	if err == nil {
		return "success"
	}

	diagnostic := err.Error()
	for _, name := range []string{"AI_API_KEY", "DASHSCOPE_API_KEY", "HTTP_PROXY", "HTTPS_PROXY"} {
		if value := strings.TrimSpace(os.Getenv(name)); value != "" {
			diagnostic = strings.ReplaceAll(diagnostic, value, "<redacted>")
			if parsed, parseErr := url.Parse(value); parseErr == nil && parsed.Host != "" {
				diagnostic = strings.ReplaceAll(diagnostic, parsed.Host, "<redacted-proxy>")
			}
		}
	}
	diagnostic = urlPattern.ReplaceAllStringFunc(diagnostic, SafeURLForLog)
	return summarize(diagnostic, 300)
}

func configFromEnv(provider string, timeoutSeconds int) Config {
	baseURL := strings.TrimSpace(os.Getenv("AI_BASE_URL"))
	model := strings.TrimSpace(os.Getenv("AI_MODEL"))
	if provider == ProviderQwen {
		if baseURL == "" {
			baseURL = defaultQwenBaseURL
		}
		if model == "" {
			model = defaultQwenModel
		}
	}

	return Config{
		Provider:       provider,
		APIKey:         effectiveAPIKey(),
		BaseURL:        baseURL,
		Model:          model,
		TimeoutSeconds: timeoutSeconds,
	}
}

func effectiveAPIKey() string {
	if value := strings.TrimSpace(os.Getenv("AI_API_KEY")); value != "" {
		return value
	}
	return strings.TrimSpace(os.Getenv("DASHSCOPE_API_KEY"))
}

func envSet(name string) bool {
	return strings.TrimSpace(os.Getenv(name)) != ""
}

func containsString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func aiTimeoutSecondsFromEnv() (int, error) {
	value := strings.TrimSpace(os.Getenv("AI_TIMEOUT_SECONDS"))
	if value == "" {
		return defaultAITimeoutSeconds, nil
	}

	seconds, err := strconv.Atoi(value)
	if err != nil || seconds <= 0 {
		return 0, fmt.Errorf("AI_TIMEOUT_SECONDS must be a positive integer number of seconds")
	}
	return seconds, nil
}

func normalizedProvider(value string) string {
	provider := strings.ToLower(strings.TrimSpace(value))
	if provider == "" {
		return ProviderMock
	}
	return provider
}

func dotenvCandidates() []string {
	return []string{
		".env",
		"../.env",
	}
}

func loadDotEnvIfPresent(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}

		key = strings.TrimSpace(key)
		if key == "" || os.Getenv(key) != "" {
			continue
		}

		value = strings.TrimSpace(value)
		value = strings.Trim(value, `"'`)
		_ = os.Setenv(key, value)
	}
	return true
}
