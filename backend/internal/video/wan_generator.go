package video

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const wanResponseLimit = 1_000_000

type WanVideoGenerator struct {
	config     ProviderConfig
	store      VideoTaskStore
	apiKey     string
	httpClient *http.Client
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
		config:     config,
		store:      store,
		apiKey:     apiKey,
		httpClient: client,
	}
}

func (g *WanVideoGenerator) CreateTask(ctx context.Context, prompt VideoPrompt) (string, error) {
	if strings.TrimSpace(prompt.ShotID) == "" {
		return "", fmt.Errorf("shot_id is required")
	}
	if strings.TrimSpace(prompt.Prompt) == "" {
		return "", fmt.Errorf("prompt is required")
	}

	requestBody := buildWanCreateRequest(g.config, prompt)
	var response wanResponse
	if err := g.doJSON(ctx, http.MethodPost, "/services/aigc/video-generation/video-synthesis", requestBody, true, &response); err != nil {
		return "", err
	}
	if response.Code != "" {
		return "", g.redactError(fmt.Errorf("wan create task failed: %s", wanErrorMessage(response.Code, response.Message)))
	}
	if strings.TrimSpace(response.Output.TaskID) == "" {
		return "", fmt.Errorf("wan create task response missing output.task_id")
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

func (g *WanVideoGenerator) GetTask(ctx context.Context, taskID string) (*VideoResult, error) {
	task, err := g.store.Get(ctx, taskID)
	if err != nil {
		return nil, err
	}

	var response wanResponse
	if err := g.doJSON(ctx, http.MethodGet, "/tasks/"+taskID, nil, false, &response); err != nil {
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
	if err := g.store.Update(ctx, task); err != nil {
		return nil, fmt.Errorf("update wan video task: %w", err)
	}

	return &VideoResult{
		TaskID:       task.TaskID,
		ShotID:       task.ShotID,
		Status:       task.Status,
		VideoURL:     task.VideoURL,
		ErrorMessage: task.ErrorMessage,
	}, nil
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
		return g.redactError(fmt.Errorf("create Wan request: %w", err))
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
		return g.redactError(fmt.Errorf("call Wan API: %w", err))
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(io.LimitReader(response.Body, wanResponseLimit))
	if err != nil {
		return fmt.Errorf("read Wan response: %w", err)
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return g.redactError(fmt.Errorf("Wan API returned status %d: %s", response.StatusCode, summarizeWanResponse(responseBody)))
	}
	if err := json.Unmarshal(responseBody, target); err != nil {
		return fmt.Errorf("parse Wan response: %w", err)
	}
	return nil
}

func (g *WanVideoGenerator) redactError(err error) error {
	if err == nil {
		return err
	}
	return fmt.Errorf("%s", g.redactText(err.Error()))
}

func (g *WanVideoGenerator) redactText(value string) string {
	if g.apiKey == "" {
		return value
	}
	return strings.ReplaceAll(value, g.apiKey, "<redacted>")
}

func summarizeWanResponse(value []byte) string {
	text := strings.TrimSpace(string(value))
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "\r", " ")
	if len(text) > 500 {
		return text[:500] + "..."
	}
	return text
}
