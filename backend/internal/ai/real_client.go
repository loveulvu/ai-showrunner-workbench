package ai

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
	"strings"
	"time"

	"ai-showrunner-workbench/internal/analysis"
	"ai-showrunner-workbench/internal/fidelity"
	"ai-showrunner-workbench/internal/novel"
	"ai-showrunner-workbench/internal/screenplay"
	"ai-showrunner-workbench/internal/showrunner"
	"ai-showrunner-workbench/internal/story"
)

type RealClient struct {
	apiKey      string
	model       string
	endpoint    string
	httpClient  *http.Client
	timeout     time.Duration
	retryDelays []time.Duration
}

func NewRealClient(cfg Config) *RealClient {
	timeoutSeconds := cfg.TimeoutSeconds
	if timeoutSeconds <= 0 {
		timeoutSeconds = defaultAITimeoutSeconds
	}
	timeout := time.Duration(timeoutSeconds) * time.Second
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.Proxy = http.ProxyFromEnvironment

	return &RealClient{
		apiKey:   cfg.APIKey,
		model:    cfg.Model,
		endpoint: chatCompletionsEndpoint(cfg.BaseURL),
		httpClient: &http.Client{
			Timeout:   timeout,
			Transport: transport,
		},
		timeout:     timeout,
		retryDelays: []time.Duration{time.Second, 2 * time.Second},
	}
}

func (c *RealClient) AnalyzeChapter(ctx context.Context, chapter novel.Chapter) (analysis.ChapterAnalysis, error) {
	var result analysis.ChapterAnalysis
	step := fmt.Sprintf("analyze chapter %d", chapter.Number)
	if err := c.callJSON(ctx, step, BuildChapterAnalysisPrompt(chapter), chapterAnalysisSchemaDescription(), &result); err != nil {
		return result, err
	}
	result.ChapterNumber = chapter.Number
	result.ChapterTitle = chapter.Title
	return result, nil
}

func (c *RealClient) MergeStoryBible(ctx context.Context, analyses []analysis.ChapterAnalysis) (story.StoryBible, error) {
	var result story.StoryBible
	if err := c.callJSON(ctx, "merge story bible", BuildStoryBiblePrompt(analyses), storyBibleSchemaDescription(), &result); err != nil {
		return result, err
	}
	return result, nil
}

func (c *RealClient) GenerateScreenplay(ctx context.Context, bible story.StoryBible, analyses []analysis.ChapterAnalysis) (screenplay.Screenplay, error) {
	var result screenplay.Screenplay
	prompt := BuildScreenplayPrompt(bible, analyses)
	if err := c.callJSON(ctx, "generate screenplay", prompt, screenplaySchemaDescription(), &result); err != nil {
		return result, err
	}

	validation := screenplay.Validate(result)
	if validation.Passed {
		return result, nil
	}

	log.Printf("LLM generate screenplay validation failed; starting one validation repair retry: %s", strings.Join(validation.Errors, "; "))
	repairPrompt := BuildRepairScreenplayPrompt(prompt, result, validation.Errors)
	var repaired screenplay.Screenplay
	if err := c.callJSONNoRepair(ctx, "generate screenplay validation repair", repairPrompt, &repaired); err != nil {
		return result, fmt.Errorf("generate screenplay: validation failed and repair parse failed: %w", err)
	}

	repairedValidation := screenplay.Validate(repaired)
	if !repairedValidation.Passed {
		return result, fmt.Errorf("generate screenplay: validation failed after repair: %s", strings.Join(repairedValidation.Errors, "; "))
	}

	return repaired, nil
}

func (c *RealClient) GenerateShowrunner(ctx context.Context, input showrunner.GenerateInput) (showrunner.ShowrunnerResult, error) {
	raw, err := c.callChatContent(ctx, "generate showrunner assets", showrunner.BuildPrompt(input))
	if err != nil {
		return showrunner.ShowrunnerResult{}, &showrunner.StageError{Stage: showrunner.StageService, Message: "Qwen showrunner request failed", Err: err}
	}

	result, stageErr := parseAndValidateShowrunner(raw)
	if stageErr == nil {
		return result, nil
	}
	c.logShowrunnerFailure(stageErr, raw)

	log.Printf("LLM generate showrunner assets %s failed; starting one repair request: %s", stageErr.Stage, stageErr.Message)
	repairRaw, err := c.callChatContent(ctx, "generate showrunner assets repair", showrunner.BuildRepairPrompt(raw, stageErr.Error()))
	if err != nil {
		return result, &showrunner.StageError{Stage: stageErr.Stage, Message: "showrunner repair request failed", Err: err}
	}

	repaired, repairErr := parseAndValidateShowrunner(repairRaw)
	if repairErr != nil {
		c.logShowrunnerFailure(repairErr, repairRaw)
		return repaired, &showrunner.StageError{Stage: repairErr.Stage, Message: "showrunner repair output failed: " + repairErr.Message, Err: repairErr.Err}
	}
	return repaired, nil
}

func parseAndValidateShowrunner(raw string) (showrunner.ShowrunnerResult, *showrunner.StageError) {
	result, err := showrunner.ParseJSON(raw)
	if err != nil {
		var stageErr *showrunner.StageError
		if errors.As(err, &stageErr) {
			return result, stageErr
		}
		return result, &showrunner.StageError{Stage: showrunner.StageParseJSON, Message: "could not parse showrunner output", Err: err}
	}

	validation := showrunner.Validate(result)
	result.Warnings = showrunner.FlexibleStringList(validation.Warnings)
	if !validation.Passed {
		return result, &showrunner.StageError{Stage: showrunner.StageValidate, Message: strings.Join(validation.Errors, "; ")}
	}
	return result, nil
}

func (c *RealClient) logShowrunnerFailure(stageErr *showrunner.StageError, raw string) {
	path, err := showrunner.SaveFailedRaw(raw, stageErr.Stage)
	if err != nil {
		log.Printf("Showrunner failed stage=%s message=%s; raw debug save failed: %v", stageErr.Stage, stageErr.Message, err)
		return
	}
	log.Printf("Showrunner failed stage=%s message=%s raw_debug=%s", stageErr.Stage, stageErr.Message, path)
}

func (c *RealClient) CheckFidelity(ctx context.Context, current screenplay.Screenplay, bible story.StoryBible, analyses []analysis.ChapterAnalysis) (fidelity.FidelityResult, error) {
	var result fidelity.FidelityResult
	if err := c.callJSON(ctx, "fidelity check", BuildFidelityCheckPrompt(current, bible, analyses), fidelityResultSchemaDescription(), &result); err != nil {
		return result, err
	}
	if result.Issues == nil {
		result.Issues = []fidelity.FidelityIssue{}
	}
	if len(result.Issues) == 0 {
		result.Passed = true
	}
	return result, nil
}

func (c *RealClient) RepairFidelity(ctx context.Context, current screenplay.Screenplay, bible story.StoryBible, analyses []analysis.ChapterAnalysis, result fidelity.FidelityResult) (screenplay.Screenplay, error) {
	var repaired screenplay.Screenplay
	if err := c.callJSON(ctx, "fidelity repair", BuildFidelityRepairPrompt(current, bible, analyses, result), screenplaySchemaDescription(), &repaired); err != nil {
		return current, err
	}
	return repaired, nil
}

func (c *RealClient) callJSON(ctx context.Context, step string, prompt string, schema string, target any) error {
	raw, err := c.callChatContent(ctx, step, prompt)
	if err != nil {
		return err
	}

	jsonText := extractJSON(raw)
	if err := json.Unmarshal([]byte(jsonText), target); err == nil {
		return nil
	} else {
		log.Printf("LLM %s parse failed; starting one JSON repair retry: %v", step, err)
		repairPrompt := BuildRepairJSONPrompt(prompt, raw, err, schema)
		repairRaw, repairErr := c.callChatContent(ctx, step+" JSON repair", repairPrompt)
		if repairErr != nil {
			return fmt.Errorf("%s: parse failed and repair request failed: %w", step, repairErr)
		}

		repairJSON := extractJSON(repairRaw)
		if repairUnmarshalErr := json.Unmarshal([]byte(repairJSON), target); repairUnmarshalErr != nil {
			return fmt.Errorf("%s: parse failed after repair: %w; output: %s", step, repairUnmarshalErr, summarize(repairJSON, 500))
		}

		return nil
	}
}

func (c *RealClient) callJSONNoRepair(ctx context.Context, step string, prompt string, target any) error {
	raw, err := c.callChatContent(ctx, step, prompt)
	if err != nil {
		return err
	}

	jsonText := extractJSON(raw)
	if err := json.Unmarshal([]byte(jsonText), target); err != nil {
		return fmt.Errorf("%s: parse failed: %w; output: %s", step, err, summarize(jsonText, 500))
	}

	return nil
}

func (c *RealClient) callChatContent(ctx context.Context, step string, prompt string) (string, error) {
	start := time.Now()
	log.Printf("LLM %s started (model=%s, timeout=%ds)", step, c.model, int(c.timeout.Seconds()))

	reqBody := chatCompletionRequest{
		Model: c.model,
		Messages: []chatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: prompt},
		},
		Temperature: 0.3,
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return c.failLLMCall(step, "marshal request", start, err)
	}

	var resp *http.Response
	var body []byte
	for attempt := 0; ; attempt++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(payload))
		if err != nil {
			return c.failLLMCall(step, "create request", start, err)
		}
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
		req.Header.Set("Content-Type", "application/json")

		resp, err = c.httpClient.Do(req)
		if err != nil {
			if c.shouldRetryNetworkError(ctx, step, err, attempt) {
				if waitErr := c.waitBeforeRetry(ctx, step, attempt, err); waitErr != nil {
					return c.failLLMCall(step, "call chat completions", start, waitErr)
				}
				continue
			}
			return c.failLLMCall(step, "call chat completions", start, err)
		}

		body, err = io.ReadAll(io.LimitReader(resp.Body, 1_000_000))
		_ = resp.Body.Close()
		if err != nil {
			if !isClientErrorStatus(resp.StatusCode) && c.shouldRetryNetworkError(ctx, step, err, attempt) {
				if waitErr := c.waitBeforeRetry(ctx, step, attempt, err); waitErr != nil {
					return c.failLLMCall(step, "read response", start, waitErr)
				}
				continue
			}
			return c.failLLMCall(step, "read response", start, err)
		}
		if c.shouldRetryHTTPStatus(ctx, resp.StatusCode, attempt) {
			statusErr := fmt.Errorf("chat completions returned retryable status %d", resp.StatusCode)
			if waitErr := c.waitBeforeRetry(ctx, step, attempt, statusErr); waitErr != nil {
				return c.failLLMCall(step, "call chat completions", start, waitErr)
			}
			continue
		}
		break
	}

	if resp.StatusCode != http.StatusOK {
		elapsed := time.Since(start)
		log.Printf("LLM %s failed in %.1fs: status %d", step, elapsed.Seconds(), resp.StatusCode)
		return "", fmt.Errorf("%s: chat completions returned status %d: %s", step, resp.StatusCode, summarize(string(body), 500))
	}

	var chatResp chatCompletionResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		elapsed := time.Since(start)
		log.Printf("LLM %s failed in %.1fs: parse chat completions response", step, elapsed.Seconds())
		return "", fmt.Errorf("%s: invalid JSON response from chat completions after %.1fs: %w; response: %s", step, elapsed.Seconds(), err, summarize(string(body), 300))
	}
	if len(chatResp.Choices) == 0 {
		elapsed := time.Since(start)
		log.Printf("LLM %s failed in %.1fs: empty choices", step, elapsed.Seconds())
		return "", fmt.Errorf("%s: chat completions response has empty choices after %.1fs", step, elapsed.Seconds())
	}

	content := strings.TrimSpace(chatResp.Choices[0].Message.Content)
	if content == "" {
		elapsed := time.Since(start)
		log.Printf("LLM %s failed in %.1fs: empty model output", step, elapsed.Seconds())
		return "", fmt.Errorf("%s: chat completions response has empty model output after %.1fs", step, elapsed.Seconds())
	}

	elapsed := time.Since(start)
	log.Printf("LLM %s finished in %.1fs", step, elapsed.Seconds())
	return content, nil
}

func (c *RealClient) shouldRetryNetworkError(ctx context.Context, step string, err error, attempt int) bool {
	if isShowrunnerStep(step) && isAwaitingHeadersTimeout(err) {
		return false
	}
	return attempt < len(c.retryDelays) && ctx.Err() == nil && isTemporaryNetworkError(err)
}

func (c *RealClient) shouldRetryHTTPStatus(ctx context.Context, statusCode int, attempt int) bool {
	if ctx.Err() != nil || attempt >= len(c.retryDelays) {
		return false
	}
	return statusCode == http.StatusTooManyRequests ||
		statusCode == http.StatusBadGateway ||
		statusCode == http.StatusServiceUnavailable ||
		statusCode == http.StatusGatewayTimeout
}

func (c *RealClient) waitBeforeRetry(ctx context.Context, step string, attempt int, err error) error {
	delay := c.retryDelays[attempt]
	log.Printf("LLM %s temporary network error; retrying attempt %d/%d in %s: %s", step, attempt+1, len(c.retryDelays), delay, RedactedDiagnostic(err))

	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func (c *RealClient) failLLMCall(step string, phase string, start time.Time, err error) (string, error) {
	elapsed := time.Since(start)
	log.Printf("LLM %s failed in %.1fs: %s", step, elapsed.Seconds(), phase)

	if isTimeoutError(err) {
		if isShowrunnerStep(step) {
			return "", fmt.Errorf("%s: %s timeout after %.1fs (timeout seconds=%d): %w; reduce Showrunner input size before increasing AI_TIMEOUT_SECONDS", step, phase, elapsed.Seconds(), int(c.timeout.Seconds()), err)
		}
		return "", fmt.Errorf("%s: %s timeout after %.1fs (timeout seconds=%d): %w; increase AI_TIMEOUT_SECONDS or reduce input size", step, phase, elapsed.Seconds(), int(c.timeout.Seconds()), err)
	}

	return "", fmt.Errorf("%s: %s failed after %.1fs: %w", step, phase, elapsed.Seconds(), err)
}

func isShowrunnerStep(step string) bool {
	return strings.HasPrefix(strings.ToLower(strings.TrimSpace(step)), "generate showrunner assets")
}

func isAwaitingHeadersTimeout(err error) bool {
	return isTimeoutError(err) && strings.Contains(strings.ToLower(err.Error()), "awaiting headers")
}

func isTimeoutError(err error) bool {
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}

	var netErr net.Error
	return errors.As(err, &netErr) && netErr.Timeout()
}

func isTemporaryNetworkError(err error) bool {
	if err == nil || errors.Is(err, context.Canceled) {
		return false
	}
	if isTimeoutError(err) || errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
		return true
	}

	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Temporary() {
		return true
	}

	message := strings.ToLower(err.Error())
	for _, fragment := range []string{"connection reset", "forcibly closed", "unexpected eof", "temporary network error"} {
		if strings.Contains(message, fragment) {
			return true
		}
	}
	return false
}

func isClientErrorStatus(statusCode int) bool {
	return statusCode >= http.StatusBadRequest && statusCode < http.StatusInternalServerError
}

func chatCompletionsEndpoint(baseURL string) string {
	base := strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if strings.HasSuffix(base, "/chat/completions") {
		return base
	}
	return base + "/chat/completions"
}

func summarize(value string, limit int) string {
	text := strings.TrimSpace(value)
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "\r", " ")
	if len(text) <= limit {
		return text
	}
	return text[:limit] + "..."
}

type chatCompletionRequest struct {
	Model       string        `json:"model"`
	Messages    []chatMessage `json:"messages"`
	Temperature float64       `json:"temperature"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatCompletionResponse struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
}
