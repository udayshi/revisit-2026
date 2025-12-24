// main.go
package usopenai

import (
	"context"
	"fmt"
	"log"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

func Demo() {
	// Requires OPENAI_API_KEY env var to be set.
	llm, err := openai.New()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	prompt := "Explain what Go is in one short paragraph."

	resp, err := llms.GenerateFromSinglePrompt(ctx, llm, prompt)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Demo Use of langchain Open AI")
	fmt.Println(resp)
}
