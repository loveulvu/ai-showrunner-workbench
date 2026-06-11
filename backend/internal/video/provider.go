package video

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultVideoProvider            = "mock"
	defaultVideoModel               = "wan2.6-t2v"
	defaultVideoTimeoutSeconds      = 600
	defaultVideoPollIntervalSeconds = 10
)

type ProviderConfig struct {
	Provider            string
	Model               string
	BaseURL             string
	APIKeySet           bool
	TimeoutSeconds      int
	PollIntervalSeconds int
}

func ProviderConfigFromEnv() (ProviderConfig, error) {
	timeoutSeconds, err := positiveIntFromEnv("VIDEO_TIMEOUT_SECONDS", defaultVideoTimeoutSeconds)
	if err != nil {
		return ProviderConfig{}, err
	}
	pollIntervalSeconds, err := positiveIntFromEnv("VIDEO_POLL_INTERVAL_SECONDS", defaultVideoPollIntervalSeconds)
	if err != nil {
		return ProviderConfig{}, err
	}

	provider := strings.ToLower(strings.TrimSpace(os.Getenv("VIDEO_PROVIDER")))
	if provider == "" {
		provider = defaultVideoProvider
	}
	model := strings.TrimSpace(os.Getenv("VIDEO_MODEL"))
	if model == "" {
		model = defaultVideoModel
	}

	return ProviderConfig{
		Provider:            provider,
		Model:               model,
		BaseURL:             strings.TrimSpace(os.Getenv("VIDEO_BASE_URL")),
		APIKeySet:           strings.TrimSpace(os.Getenv("VIDEO_API_KEY")) != "",
		TimeoutSeconds:      timeoutSeconds,
		PollIntervalSeconds: pollIntervalSeconds,
	}, nil
}

func NewGeneratorFromConfig(config ProviderConfig, store VideoTaskStore) (VideoGenerator, error) {
	if store == nil {
		return nil, fmt.Errorf("video task store is required")
	}

	switch config.Provider {
	case "mock":
		return NewMockVideoGeneratorWithStore(config, store), nil
	case "wan":
		apiKey := strings.TrimSpace(os.Getenv("VIDEO_API_KEY"))
		if strings.TrimSpace(config.BaseURL) == "" {
			return nil, fmt.Errorf("VIDEO_PROVIDER=wan requires VIDEO_BASE_URL")
		}
		if apiKey == "" {
			return nil, fmt.Errorf("VIDEO_PROVIDER=wan requires VIDEO_API_KEY")
		}
		timeoutSeconds := config.TimeoutSeconds
		if timeoutSeconds <= 0 {
			timeoutSeconds = defaultVideoTimeoutSeconds
		}
		transport := http.DefaultTransport.(*http.Transport).Clone()
		transport.Proxy = http.ProxyFromEnvironment
		client := &http.Client{
			Transport: transport,
			Timeout:   time.Duration(timeoutSeconds) * time.Second,
		}
		return NewWanVideoGenerator(config, store, apiKey, client), nil
	default:
		return nil, fmt.Errorf("unknown VIDEO_PROVIDER %q", config.Provider)
	}
}

func NewGeneratorFromEnv(store VideoTaskStore) (VideoGenerator, error) {
	config, err := ProviderConfigFromEnv()
	if err != nil {
		return nil, err
	}
	return NewGeneratorFromConfig(config, store)
}

func LogProviderConfig(logger *log.Logger, config ProviderConfig) {
	logger.Printf("VIDEO_PROVIDER=%s", config.Provider)
	logger.Printf("VIDEO_MODEL=%s", config.Model)
	logger.Printf("VIDEO_BASE_URL set: %t", strings.TrimSpace(config.BaseURL) != "")
	logger.Printf("VIDEO_API_KEY set: %t", config.APIKeySet)
}

func positiveIntFromEnv(name string, fallback int) (int, error) {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return fallback, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return 0, fmt.Errorf("%s must be a positive integer number of seconds", name)
	}
	return parsed, nil
}
