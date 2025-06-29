package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type ChatService struct {
	chatManager *ChatManager
	openaiClient *openai.Client
	systemPrompt string
}

func NewChatService() *ChatService {
	return &ChatService{
		chatManager: NewChatManager(),
		openaiClient: initOpenAIClient(),
		systemPrompt: loadSystemPrompt("app/data/system_prompt.txt"),
	}
}

func initOpenAIClient() *openai.Client {
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
	return &client
}

func loadSystemPrompt(filename string) string {
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("Error loading system prompt: %v\n", err)
		return "You are a friendly and efficient customer service attendant eager to assist customers with their inquiries and concerns."
	}
	return string(content)
}

func(cs *ChatService) CreateChat(userID string) string {
	chatID := uuid.New().String()
	cs.chatManager.CreateChat(userID, chatID, cs.systemPrompt)
	return chatID
}

func(cs *ChatService) ProcessMessage(userID, chatID, message string) (string, error) {
    if  _, exist := cs.chatManager.GetChat(userID, chatID); !exist {
        return "", errors.New("chat not found")
    }
    cs.chatManager.AddMessage(userID, chatID, openai.UserMessage(message))
	conversation := cs.chatManager.GetConversation(userID, chatID)

	response, err := cs.openaiClient.Chat.Completions.New(
		context.TODO(),
		openai.ChatCompletionNewParams{
			Messages: conversation,
			Model: openai.ChatModelGPT4,
			Temperature: openai.Float(0.7),
			MaxCompletionTokens: openai.Int(500),
		},
	)

	if err != nil {
		return "", fmt.Errorf("error getting response: %v", err)
	}
	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no response choices were returned by the OpenAI API")
	}
	reply := response.Choices[0].Message.Content
	cs.chatManager.AddMessage(userID, chatID, openai.AssistantMessage(reply))
    return reply, nil
}