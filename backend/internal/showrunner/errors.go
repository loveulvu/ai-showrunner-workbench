package showrunner

import (
	"errors"
	"fmt"
)

const (
	StageExtractJSON = "extract_json"
	StageParseJSON   = "parse_json"
	StageValidate    = "validate"
	StageService     = "service"
)

type StageError struct {
	Stage   string
	Message string
	Err     error
}

func (e *StageError) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("%s: %s", e.Stage, e.Message)
	}
	return fmt.Sprintf("%s: %s: %v", e.Stage, e.Message, e.Err)
}

func (e *StageError) Unwrap() error {
	return e.Err
}

func ErrorDetails(err error) (string, string) {
	var stageErr *StageError
	if errors.As(err, &stageErr) {
		message := stageErr.Message
		if stageErr.Err != nil {
			message += ": " + stageErr.Err.Error()
		}
		if len(message) > 500 {
			message = message[:500] + "..."
		}
		return stageErr.Stage, message
	}
	return StageService, err.Error()
}
