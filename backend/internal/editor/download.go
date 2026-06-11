package editor

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const maxClipDownloadBytes = 500 * 1024 * 1024

var unsafeFileCharacters = regexp.MustCompile(`[^a-zA-Z0-9_-]+`)

func ClipFileName(shotID string) (string, error) {
	name := strings.Trim(unsafeFileCharacters.ReplaceAllString(strings.TrimSpace(shotID), "_"), "_")
	if name == "" {
		return "", fmt.Errorf("shot_id is required")
	}
	return "clip_" + name + ".mp4", nil
}

func DownloadClip(ctx context.Context, clip ClipAsset) (ClipAsset, error) {
	if strings.TrimSpace(clip.SourceURL) == "" {
		return clip, fmt.Errorf("download clip %q: source_url is required", clip.ShotID)
	}
	if strings.TrimSpace(clip.LocalPath) == "" {
		name, err := ClipFileName(clip.ShotID)
		if err != nil {
			return clip, err
		}
		clip.LocalPath = filepath.Join("outputs", "clips", name)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, clip.SourceURL, nil)
	if err != nil {
		return clip, fmt.Errorf("download clip %q: create request: %s", clip.ShotID, redactSourceURL(err.Error(), clip.SourceURL))
	}

	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.Proxy = http.ProxyFromEnvironment
	client := &http.Client{
		Transport: transport,
		Timeout:   2 * time.Minute,
	}
	response, err := client.Do(request)
	if err != nil {
		return clip, fmt.Errorf("download clip %q failed: %s", clip.ShotID, redactSourceURL(err.Error(), clip.SourceURL))
	}
	defer response.Body.Close()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return clip, fmt.Errorf("download clip %q failed with status %d", clip.ShotID, response.StatusCode)
	}

	if err := os.MkdirAll(filepath.Dir(clip.LocalPath), 0o755); err != nil {
		return clip, fmt.Errorf("create clip directory: %w", err)
	}
	file, err := os.Create(clip.LocalPath)
	if err != nil {
		return clip, fmt.Errorf("create local clip %q: %w", clip.ShotID, err)
	}
	defer file.Close()

	written, err := io.Copy(file, io.LimitReader(response.Body, maxClipDownloadBytes+1))
	if err != nil {
		return clip, fmt.Errorf("write local clip %q: %w", clip.ShotID, err)
	}
	if written > maxClipDownloadBytes {
		return clip, fmt.Errorf("download clip %q exceeded %d bytes", clip.ShotID, maxClipDownloadBytes)
	}
	return clip, nil
}

func redactSourceURL(message string, sourceURL string) string {
	if strings.TrimSpace(sourceURL) == "" {
		return message
	}
	return strings.ReplaceAll(message, sourceURL, "<redacted-url>")
}
