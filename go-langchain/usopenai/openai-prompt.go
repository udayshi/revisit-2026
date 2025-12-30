// main.go
package usopenai

import (
	"context"
	"fmt"
	"log"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/prompts"
)

func DemoPrompt() {
	llm, err := openai.New()
	if err != nil {
		log.Fatal(err)
	}

	// Define a prompt template with variables.
	tmpl := prompts.NewPromptTemplate(
		"Answer as a senior Go engineer.\nQuestion: {{.question}}\nAnswer:",
		[]string{"question"},
	)

	// Build an LLM chain from the template.
	chain := chains.NewLLMChain(llm, tmpl)

	ctx := context.Background()

	out, err := chains.Run(ctx, chain, map[string]any{
		"question": "How would you structure a Go microservice that calls an LLM for summarization?",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Demo Use of langchain Open AI Prompt")
	fmt.Println(out)
}
