package video

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

type TaskDiagnostic struct {
	HTTPStatus int
	Code       string
	Message    string
	Latency    time.Duration
}

func CheckTask(ctx context.Context, config ProviderConfig, apiKey string, client *http.Client, taskID string) (TaskDiagnostic, error) {
	start := time.Now()
	endpoint := strings.TrimRight(strings.TrimSpace(config.BaseURL), "/") + "/tasks/" + url.PathEscape(taskID)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return TaskDiagnostic{}, fmt.Errorf("create Wan diagnostic request: %w", err)
	}
	request.Header.Set("Authorization", "Bearer "+apiKey)

	response, err := client.Do(request)
	if err != nil {
		return TaskDiagnostic{Latency: time.Since(start)}, fmt.Errorf("call Wan diagnostic endpoint: %w", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(io.LimitReader(response.Body, wanResponseLimit))
	if err != nil {
		return TaskDiagnostic{HTTPStatus: response.StatusCode, Latency: time.Since(start)}, fmt.Errorf("read Wan diagnostic response: %w", err)
	}
	var payload wanResponse
	_ = json.Unmarshal(body, &payload)
	return TaskDiagnostic{
		HTTPStatus: response.StatusCode,
		Code:       payload.Code,
		Message:    payload.Message,
		Latency:    time.Since(start),
	}, nil
}

func CheckMediaURL(ctx context.Context, client *http.Client, value string) (int, time.Duration, error) {
	start := time.Now()
	request, err := http.NewRequestWithContext(ctx, http.MethodHead, value, nil)
	if err != nil {
		return 0, 0, err
	}
	response, err := client.Do(request)
	if err != nil {
		return 0, time.Since(start), err
	}
	defer response.Body.Close()
	return response.StatusCode, time.Since(start), nil
}

func SafeURL(value string) string {
	parsed, err := url.Parse(strings.TrimSpace(value))
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return "<invalid URL>"
	}
	return parsed.Host + parsed.EscapedPath()
}

func SafeDiagnosticText(value string, secrets ...string) string {
	for _, secret := range append(secrets, os.Getenv("HTTP_PROXY"), os.Getenv("HTTPS_PROXY")) {
		secret = strings.TrimSpace(secret)
		if secret == "" {
			continue
		}
		value = strings.ReplaceAll(value, secret, "<redacted>")
		if parsed, err := url.Parse(secret); err == nil && parsed.Host != "" {
			value = strings.ReplaceAll(value, parsed.Host, "<redacted>")
		}
	}
	return diagnosticURLPattern.ReplaceAllString(value, "<redacted-url>")
}

func ClassifyTaskDiagnostic(result TaskDiagnostic, err error) string {
	if err != nil {
		message := strings.ToLower(err.Error())
		if strings.Contains(message, "timeout") || strings.Contains(message, "connection reset") || strings.Contains(message, "forcibly closed") {
			return "Network or proxy connection is unstable."
		}
		return "Wan task query failed before receiving a usable response."
	}

	code := strings.ToLower(result.Code)
	message := strings.ToLower(result.Message)
	switch {
	case result.HTTPStatus >= 500:
		return "Wan upstream service or proxy chain returned a server error."
	case result.HTTPStatus == http.StatusUnauthorized && (strings.Contains(code, "invalidapikey") || strings.Contains(message, "invalid api")):
		return "API key, endpoint region, or permission does not match."
	case result.HTTPStatus == http.StatusUnauthorized:
		return "No API key was accepted; VIDEO_API_KEY and AI_API_KEY fallback should be checked."
	case result.HTTPStatus == http.StatusNotFound || strings.Contains(code, "notfound") || strings.Contains(message, "not found") || strings.Contains(message, "invalid task"):
		return "Endpoint, network, proxy, and authentication are reachable; the task id does not exist."
	default:
		return "Wan task query completed."
	}
}

var diagnosticURLPattern = regexp.MustCompile(`https?://[^\s"']+`)
