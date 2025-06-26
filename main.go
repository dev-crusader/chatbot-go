package main

import (
	"fmt"
	"os"

	"context"
	"log"

	"github.com/google/uuid"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

var (
	client = openai.Client{}
	chatSessions map[string]*[]openai.ChatCompletionMessageParamUnion
	systemPrompt = openai.SystemMessage(loadSystemPrompt("app/data/system_prompt.txt"))
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

	client = openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithBaseURL(baseURL),
	)

	chatID1 := createUniqueID()
    _, err := sendMessage(chatID1, "I'm having trouble with my recent order. Can you help me track it?")
    if err != nil {
        log.Fatal(err)
    }

	_, err = sendMessage(chatID1, "It was supposed to arrive yesterday but hasn't. What should I do next?")
    if err != nil {
        log.Fatal(err)
    }
	printChatHistory(chatSessions[chatID1])
}

func sendMessage(id string, message string) (string, error) {
    if _, exists := chatSessions[id]; !exists {
        return "", fmt.Errorf("chat session Id: %s not found", id)
    }
    *chatSessions[id] = append(*chatSessions[id], openai.UserMessage(message))
    
    req := openai.ChatCompletionNewParams{
        Messages: *chatSessions[id],
        Model: openai.ChatModelGPT4,
    }
    
    response, err := client.Chat.Completions.New(context.TODO(), req)
    if err != nil {
        return "", err
    }
    if len(response.Choices) == 0 {
        return "", fmt.Errorf("no response choices were returned by the OpenAI API")
    }
    answer := response.Choices[0].Message.Content
    *chatSessions[id] = append(*chatSessions[id], openai.AssistantMessage(answer))
    return answer, nil
}

func printChatHistory(conversation *[]openai.ChatCompletionMessageParamUnion) {
	for _, message := range *conversation {
        if message.OfUser != nil {
            fmt.Println("User: ", message.OfUser.Content.OfString)
        }
        
        if message.OfAssistant != nil {
            fmt.Println("Assistant: ", message.OfAssistant.Content.OfString)
        }

		if message.OfSystem != nil {
			fmt.Println("System: ", message.OfSystem.Content.OfString)
		}
    }
}

func createUniqueID() string {
	chatId := uuid.New().String()
	chatSessions[chatId] = &[]openai.ChatCompletionMessageParamUnion{systemPrompt}
	return chatId
}

func loadSystemPrompt(filename string) string {
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("Error loading system prompt: %v\n", err)
		return "You are a friendly and efficient customer service attendant eager to assist customers with their inquiries and concerns."
	}
	return string(content)
}