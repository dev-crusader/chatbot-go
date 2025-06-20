package main

import (
	"os"

	"context"
	"log"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/shared"
)

func main() {
	apiKey := os.Getenv("API_KEY")
	baseURL := os.Getenv("BASE_URL")

	if apiKey == "" {
		log.Fatal("API_KEY not set in environment variable")
	}

	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}

	client := openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithBaseURL(baseURL),
	)

	prompt := `Tell me a joke`
	ctx := context.Background()
	chatCompletion, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		},
		Model: shared.ChatModelGPT4,
	})
	if err != nil {
		log.Fatal("Error: ", err)
	}

	if len(chatCompletion.Choices) == 0 {
		log.Println("No response choices returned by OpenAI API.")
		return
	}

	reply := chatCompletion.Choices[0].Message.Content

	log.Printf("Prompt: %s\n Response: %s", prompt, reply) 
}