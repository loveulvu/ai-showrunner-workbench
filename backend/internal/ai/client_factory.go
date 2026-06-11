package ai

import (
	"bufio"
	"fmt"
	"os"
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
	AIProvider           string `json:"ai_provider"`
	AIModel              string `json:"ai_model"`
	AIBaseURL            string `json:"ai_base_url"`
	AIBaseURLConfigured  bool   `json:"ai_base_url_configured"`
	AIAPIKeyConfigured   bool   `json:"ai_api_key_configured"`
	HTTPProxyConfigured  bool   `json:"http_proxy_configured"`
	HTTPSProxyConfigured bool   `json:"https_proxy_configured"`
	AITimeoutSeconds     int    `json:"ai_timeout_seconds"`
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
		missing = append(missing, "AI_API_KEY")
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

func RuntimeStatusFromEnv() RuntimeStatus {
	LoadEnv()
	timeoutSeconds, err := aiTimeoutSecondsFromEnv()
	if err != nil {
		timeoutSeconds = defaultAITimeoutSeconds
	}

	provider := normalizedProvider(os.Getenv("AI_PROVIDER"))
	cfg := configFromEnv(provider, timeoutSeconds)

	return RuntimeStatus{
		AIProvider:           provider,
		AIModel:              cfg.Model,
		AIBaseURL:            cfg.BaseURL,
		AIBaseURLConfigured:  strings.TrimSpace(os.Getenv("AI_BASE_URL")) != "",
		AIAPIKeyConfigured:   strings.TrimSpace(os.Getenv("AI_API_KEY")) != "",
		HTTPProxyConfigured:  strings.TrimSpace(os.Getenv("HTTP_PROXY")) != "",
		HTTPSProxyConfigured: strings.TrimSpace(os.Getenv("HTTPS_PROXY")) != "",
		AITimeoutSeconds:     timeoutSeconds,
	}
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
		APIKey:         strings.TrimSpace(os.Getenv("AI_API_KEY")),
		BaseURL:        baseURL,
		Model:          model,
		TimeoutSeconds: timeoutSeconds,
	}
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
