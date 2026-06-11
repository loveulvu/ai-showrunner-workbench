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
	result, err := s.client.GenerateShowrunner(ctx, input)
	if err != nil {
		return result, err
	}

	validation := Validate(result)
	result.Warnings = FlexibleStringList(validation.Warnings)
	if !validation.Passed {
		return result, fmt.Errorf("showrunner validation failed: %v", validation.Errors)
	}
	return result, nil
}
