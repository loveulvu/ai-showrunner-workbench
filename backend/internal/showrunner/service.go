package showrunner

import (
	"context"
	"fmt"
)

type Client interface {
	GenerateShowrunner(ctx context.Context, input GenerateInput) (ShowrunnerResult, error)
}

type Service struct {
	client Client
}

func NewService(client Client) *Service {
	return &Service{client: client}
}

func (s *Service) Generate(ctx context.Context, input GenerateInput) (ShowrunnerResult, error) {
	input = PrepareInput(input)
	result, err := s.client.GenerateShowrunner(ctx, input)
	if err != nil {
		if _, ok := err.(*StageError); ok {
			return result, err
		}
		return result, &StageError{Stage: StageService, Message: "showrunner generation service failed", Err: err}
	}

	result = LimitResultForMode(result, input.Mode)
	validation := Validate(result)
	result.Warnings = FlexibleStringList(validation.Warnings)
	if !validation.Passed {
		return result, &StageError{Stage: StageValidate, Message: fmt.Sprintf("showrunner validation failed: %v", validation.Errors)}
	}
	return LimitResultForMode(result, input.Mode), nil
}
