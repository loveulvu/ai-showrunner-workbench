package showrunner

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func SaveFailedRaw(raw string, stage string) (string, error) {
	root := debugOutputsRoot()
	dir := filepath.Join(root, "debug")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}

	name := fmt.Sprintf("showrunner_failed_%s_%s.txt", stage, time.Now().UTC().Format("20060102T150405.000000000Z"))
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(raw), 0o600); err != nil {
		return "", err
	}
	return path, nil
}

func debugOutputsRoot() string {
	cwd, err := os.Getwd()
	if err != nil {
		return "outputs"
	}

	current := cwd
	for {
		if _, err := os.Stat(filepath.Join(current, "go.mod")); err == nil {
			return filepath.Join(filepath.Dir(current), "outputs")
		}
		parent := filepath.Dir(current)
		if parent == current {
			return filepath.Join(cwd, "outputs")
		}
		current = parent
	}
}
