package video

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

const wanResponseLimit = 1_000_000

type WanVideoGenerator struct {
	config      ProviderConfig
	store       VideoTaskStore
	apiKey      string
	httpClient  *http.Client
	retryDelays []time.Duration
}

type wanCreateRequest struct {
	Model      string        `json:"model"`
	Input      wanInput      `json:"input"`
	Parameters wanParameters `json:"parameters"`
}

type wanInput struct {
	Prompt         string `json:"prompt"`
	NegativePrompt string `json:"negative_prompt,omitempty"`
}

type wanParameters struct {
	Size         string `json:"size"`
	Duration     int    `json:"duration"`
	PromptExtend bool   `json:"prompt_extend"`
	Watermark    bool   `json:"watermark"`
}

type wanResponse struct {
	Output  wanOutput `json:"output"`
	Code    string    `json:"code"`
	Message string    `json:"message"`
}

type wanOutput struct {
	TaskID     string `json:"task_id"`
	TaskStatus string `json:"task_status"`
	VideoURL   string `json:"video_url"`
}

func NewWanVideoGenerator(config ProviderConfig, store VideoTaskStore, apiKey string, client *http.Client) *WanVideoGenerator {
	return &WanVideoGenerator{
		config:      config,
		store:       store,
		apiKey:      apiKey,
		httpClient:  client,
		retryDelays: []time.Duration{time.Second, 2 * time.Second},
	}
}

func (g *WanVideoGenerator) CreateTask(ctx context.Context, prompt VideoPrompt) (string, error) {
	if strings.TrimSpace(prompt.ShotID) == "" {
		return "", &Error{Kind: ErrorKindRequest, Message: "shot_id is required"}
	}
	if strings.TrimSpace(prompt.Prompt) == "" {
		return "", &Error{Kind: ErrorKindRequest, Message: "prompt is required"}
	}

	requestBody := buildWanCreateRequest(g.config, prompt)
	var response wanResponse
	if err := g.doCreateJSONWithRetry(ctx, requestBody, &response); err != nil {
		return "", err
	}
	if response.Code != "" {
		return "", g.upstreamError("Wan rejected video task creation", fmt.Errorf("%s", wanErrorMessage(response.Code, response.Message)), false)
	}
	if strings.TrimSpace(response.Output.TaskID) == "" {
		return "", g.upstreamError("Wan response missing task_id", nil, false)
	}

	now := time.Now().UTC()
	status, statusMessage := mapWanStatus(response.Output.TaskStatus)
	task := &VideoTask{
		TaskID:           response.Output.TaskID,
		ShotID:           prompt.ShotID,
		Provider:         "wan",
		Model:            requestBody.Model,
		Prompt:           prompt.Prompt,
		NegativePrompt:   prompt.NegativePrompt,
		DurationSeconds:  requestBody.Parameters.Duration,
		AspectRatio:      prompt.AspectRatio,
		Subtitle:         prompt.Subtitle,
		ExpectedClipName: prompt.ExpectedClipName,
		Status:           status,
		ErrorMessage:     statusMessage,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	if err := g.store.Save(ctx, task); err != nil {
		return "", fmt.Errorf("save wan video task: %w", err)
	}
	return task.TaskID, nil
}

func (g *WanVideoGenerator) doCreateJSONWithRetry(ctx context.Context, body any, target *wanResponse) error {
	for attempt := 0; ; attempt++ {
		err := g.doJSON(ctx, http.MethodPost, "/services/aigc/video-generation/video-synthesis", body, true, target)
		if err == nil || strings.TrimSpace(target.Output.TaskID) != "" {
			return err
		}
		if attempt >= len(g.retryDelays) || !isTemporaryWanNetworkError(err) {
			return err
		}

		delay := g.retryDelays[attempt]
		log.Printf("Wan create task temporary network error; retrying attempt %d/%d in %s", attempt+1, len(g.retryDelays), delay)
		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return g.upstreamError("Wan task creation canceled while waiting to retry", ctx.Err(), false)
		case <-timer.C:
		}
	}
}

func (g *WanVideoGenerator) GetTask(ctx context.Context, taskID string) (*VideoResult, error) {
	task, err := g.store.Get(ctx, taskID)
	taskWasCached := err == nil
	if err != nil && !errors.Is(err, ErrTaskNotFound) {
		return nil, err
	}
	if errors.Is(err, ErrTaskNotFound) {
		task = &VideoTask{TaskID: taskID, Provider: "wan", Model: g.config.Model}
	}

	var response wanResponse
	if err := g.doJSON(ctx, http.MethodGet, "/tasks/"+taskID, nil, false, &response); err != nil {
		if taskWasCached && task.Status == StatusSucceeded && strings.TrimSpace(task.VideoURL) != "" {
			return videoResultFromTask(task), nil
		}
		return nil, err
	}

	status, statusMessage := mapWanStatus(response.Output.TaskStatus)
	task.Status = status
	task.VideoURL = response.Output.VideoURL
	task.ErrorMessage = statusMessage
	if response.Code != "" {
		task.Status = StatusFailed
		task.ErrorMessage = g.redactText(wanErrorMessage(response.Code, response.Message))
	}
	task.UpdatedAt = time.Now().UTC()
	if taskWasCached {
		err = g.store.Update(ctx, task)
	} else {
		task.CreatedAt = task.UpdatedAt
		err = g.store.Save(ctx, task)
	}
	if err != nil {
		return nil, fmt.Errorf("update wan video task: %w", err)
	}

	return videoResultFromTask(task), nil
}

func videoResultFromTask(task *VideoTask) *VideoResult {
	return &VideoResult{
		TaskID:       task.TaskID,
		ShotID:       task.ShotID,
		Status:       task.Status,
		VideoURL:     task.VideoURL,
		ErrorMessage: task.ErrorMessage,
	}
}

func buildWanCreateRequest(config ProviderConfig, prompt VideoPrompt) wanCreateRequest {
	model := strings.TrimSpace(config.Model)
	if model == "" {
		model = strings.TrimSpace(prompt.Model)
	}
	duration := prompt.DurationSeconds
	if duration <= 0 {
		duration = 5
	}
	return wanCreateRequest{
		Model: model,
		Input: wanInput{
			Prompt:         prompt.Prompt,
			NegativePrompt: prompt.NegativePrompt,
		},
		Parameters: wanParameters{
			Size:         wanSize(prompt.AspectRatio),
			Duration:     duration,
			PromptExtend: true,
			Watermark:    false,
		},
	}
}

func wanSize(aspectRatio string) string {
	switch strings.TrimSpace(aspectRatio) {
	case "16:9":
		return "1280*720"
	default:
		return "1280*720"
	}
}

func mapWanStatus(value string) (Status, string) {
	switch strings.ToUpper(strings.TrimSpace(value)) {
	case "PENDING":
		return StatusPending, ""
	case "RUNNING":
		return StatusRunning, ""
	case "SUCCEEDED":
		return StatusSucceeded, ""
	case "FAILED":
		return StatusFailed, ""
	default:
		return StatusFailed, fmt.Sprintf("unknown Wan task status %q", value)
	}
}

func wanErrorMessage(code string, message string) string {
	code = strings.TrimSpace(code)
	message = strings.TrimSpace(message)
	if code == "" {
		return message
	}
	if message == "" {
		return code
	}
	return code + ": " + message
}

func (g *WanVideoGenerator) doJSON(ctx context.Context, method string, path string, body any, async bool, target any) error {
	var reader io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal Wan request: %w", err)
		}
		reader = bytes.NewReader(payload)
	}

	endpoint := strings.TrimRight(strings.TrimSpace(g.config.BaseURL), "/") + path
	request, err := http.NewRequestWithContext(ctx, method, endpoint, reader)
	if err != nil {
		return g.upstreamError("Could not create Wan video service request", err, false)
	}
	request.Header.Set("Authorization", "Bearer "+g.apiKey)
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	if async {
		request.Header.Set("X-DashScope-Async", "enable")
	}

	response, err := g.httpClient.Do(request)
	if err != nil {
		return g.upstreamError("Could not reach Wan video service", err, true)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(io.LimitReader(response.Body, wanResponseLimit))
	if err != nil {
		if json.Unmarshal(responseBody, target) == nil {
			if wanTarget, ok := target.(*wanResponse); ok && strings.TrimSpace(wanTarget.Output.TaskID) != "" {
				return nil
			}
		}
		return g.upstreamError("Could not read Wan video service response", err, true)
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return g.upstreamError(fmt.Sprintf("Wan video service returned status %d", response.StatusCode), fmt.Errorf("%s", summarizeWanResponse(responseBody)), false)
	}
	if err := json.Unmarshal(responseBody, target); err != nil {
		return g.upstreamError("Could not parse Wan video service response", err, false)
	}
	return nil
}

func (g *WanVideoGenerator) upstreamError(message string, err error, retryable bool) error {
	if err != nil {
		log.Printf("Wan upstream error: %s", g.redactText(err.Error()))
	}
	return &Error{Kind: ErrorKindUpstream, Message: message, Err: err, Retryable: retryable}
}

func (g *WanVideoGenerator) redactText(value string) string {
	for _, sensitive := range []string{g.apiKey, g.config.BaseURL} {
		if strings.TrimSpace(sensitive) != "" {
			value = strings.ReplaceAll(value, sensitive, "<redacted>")
		}
	}
	for _, name := range []string{"HTTP_PROXY", "HTTPS_PROXY"} {
		proxy := strings.TrimSpace(os.Getenv(name))
		if proxy == "" {
			continue
		}
		value = strings.ReplaceAll(value, proxy, "<redacted-proxy>")
		if parsed, err := url.Parse(proxy); err == nil && parsed.Host != "" {
			value = strings.ReplaceAll(value, parsed.Host, "<redacted-proxy>")
		}
	}
	return wanURLPattern.ReplaceAllString(value, "<redacted-url>")
}

func isTemporaryWanNetworkError(err error) bool {
	if err == nil || errors.Is(err, context.Canceled) {
		return false
	}
	var videoErr *Error
	if errors.As(err, &videoErr) {
		return videoErr.Retryable
	}
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
		return true
	}

	var netErr net.Error
	if errors.As(err, &netErr) && (netErr.Timeout() || netErr.Temporary()) {
		return true
	}

	message := strings.ToLower(err.Error())
	for _, fragment := range []string{"connection reset", "forcibly closed", "timeout", "unexpected eof", "temporary network error"} {
		if strings.Contains(message, fragment) {
			return true
		}
	}
	return false
}

var wanURLPattern = regexp.MustCompile(`https?://[^\s"']+`)

func summarizeWanResponse(value []byte) string {
	text := strings.TrimSpace(string(value))
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "\r", " ")
	if len(text) > 500 {
		return text[:500] + "..."
	}
	return text
}
