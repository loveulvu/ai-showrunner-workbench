package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"ai-showrunner-workbench/internal/ai"
	"ai-showrunner-workbench/internal/video"
)

const fakeTaskID = "00000000-0000-0000-0000-000000000000"

type options struct {
	TaskID string
	Create bool
	URL    string
}

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	if err := run(os.Args[1:], logger, nil); err != nil {
		logger.Fatal(err)
	}
}

func run(args []string, logger *log.Logger, client *http.Client) error {
	opts, err := parseOptions(args, io.Discard)
	if err != nil {
		return err
	}

	ai.LoadEnv()
	for _, status := range ai.EnvFileStatuses() {
		logger.Printf("Environment file %s loaded: %t", status.Path, status.Loaded)
	}

	config, err := video.ProviderConfigFromEnv()
	if err != nil {
		return err
	}
	apiKey, _ := video.EffectiveAPIKey()
	video.LogProviderConfig(logger, config)
	logger.Printf("VIDEO_BASE_URL=%s", video.SafeURL(config.BaseURL))
	logger.Printf("HTTP_PROXY set: %t", strings.TrimSpace(os.Getenv("HTTP_PROXY")) != "")
	logger.Printf("HTTPS_PROXY set: %t", strings.TrimSpace(os.Getenv("HTTPS_PROXY")) != "")
	logger.Printf("NO_PROXY set: %t", strings.TrimSpace(os.Getenv("NO_PROXY")) != "")
	logger.Printf("http.ProxyFromEnvironment enabled: true")

	if strings.TrimSpace(config.BaseURL) == "" {
		return fmt.Errorf("VIDEO_BASE_URL is required")
	}
	if apiKey == "" {
		return fmt.Errorf("VIDEO_API_KEY is empty and AI_API_KEY fallback is unavailable")
	}
	if client == nil {
		client = video.NewHTTPClient(config)
	}

	createdTaskID := ""
	if opts.Create {
		logger.Print("WARNING: This will create a real Wan video task and may consume credits.")
		generator := video.NewWanVideoGenerator(config, video.NewMemoryVideoTaskStore(), apiKey, client)
		taskID, err := generator.CreateTask(context.Background(), video.VideoPrompt{
			ShotID:          "video-check-shot",
			Prompt:          "A simple static animation test shot with minimal motion.",
			DurationSeconds: 5,
			AspectRatio:     "16:9",
		})
		if err != nil {
			return fmt.Errorf("create real Wan test task: %s", video.SafeDiagnosticText(err.Error(), apiKey))
		}
		logger.Printf("Real Wan task created: task_id=%s", taskID)
		createdTaskID = taskID
	}

	taskID := strings.TrimSpace(opts.TaskID)
	if taskID == "" {
		taskID = createdTaskID
	}
	if taskID == "" {
		taskID = fakeTaskID
	}
	result, checkErr := video.CheckTask(context.Background(), config, apiKey, client, taskID)
	logger.Printf("GET task HTTP status=%d code=%s message=%s latency=%s",
		result.HTTPStatus,
		video.SafeDiagnosticText(result.Code, apiKey),
		video.SafeDiagnosticText(result.Message, apiKey),
		result.Latency,
	)
	logger.Printf("Diagnosis: %s", video.ClassifyTaskDiagnostic(result, checkErr))
	if checkErr != nil {
		logger.Printf("Safe error: %s", video.SafeDiagnosticText(checkErr.Error(), apiKey))
	}

	if strings.TrimSpace(opts.URL) != "" {
		status, latency, err := video.CheckMediaURL(context.Background(), client, opts.URL)
		logger.Printf("OSS URL=%s status=%d latency=%s", video.SafeURL(opts.URL), status, latency)
		if err != nil {
			logger.Printf("OSS safe error: %s", video.SafeDiagnosticText(err.Error(), apiKey, opts.URL))
		}
	}
	return nil
}

func parseOptions(args []string, output io.Writer) (options, error) {
	var opts options
	flags := flag.NewFlagSet("video-check", flag.ContinueOnError)
	flags.SetOutput(output)
	flags.StringVar(&opts.TaskID, "task-id", "", "query a specific Wan task id")
	flags.BoolVar(&opts.Create, "create", false, "create one real Wan test task")
	flags.StringVar(&opts.URL, "url", "", "HEAD an OSS video URL without downloading it")
	return opts, flags.Parse(args)
}
