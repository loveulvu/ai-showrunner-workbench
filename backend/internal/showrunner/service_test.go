package showrunner

import (
	"context"
	"testing"
)

type stubClient struct {
	result ShowrunnerResult
	err    error
}

func (s stubClient) GenerateShowrunner(context.Context, GenerateInput) (ShowrunnerResult, error) {
	return s.result, s.err
}

func TestServiceAddsValidationWarnings(t *testing.T) {
	result := MockResult(GenerateInput{})
	service := NewService(stubClient{result: result})

	got, err := service.Generate(context.Background(), GenerateInput{})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if len(got.Warnings) != 0 {
		t.Fatalf("warnings = %v, want no warnings for complete mock prompts", got.Warnings)
	}
}
