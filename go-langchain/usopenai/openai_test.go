package usopenai

import (
	"context"
	"os"
	"testing"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

func TestOpenAIIntegration(t *testing.T) {
	if os.Getenv("OPENAI_API_KEY") == "" {
		t.Skip("OPENAI_API_KEY not set; skipping integration test")
	}

	llm, err := openai.New()
	if err != nil {
		t.Fatalf("openai.New() error = %v", err)
	}

	ctx := context.Background()
	out, err := llms.GenerateFromSinglePrompt(ctx, llm, "Say 'ok'.")
	if err != nil {
		t.Fatalf("GenerateFromSinglePrompt error = %v", err)
	}
	if out == "" {
		t.Fatalf("expected non-empty output")
	}
}
