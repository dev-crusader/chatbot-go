package service

import (
	"github.com/openai/openai-go"
)

type ChatManager struct {
	chats map[string]map[string]openai.ChatCompletionNewParams
}

func NewChatManager() *ChatManager {
	return &ChatManager{
		chats:	make(map[string]map[string]openai.ChatCompletionNewParams),
	}
}

func (cm *ChatManager) CreateChat(userID, chatID, systemPrompt string) {
	if _, exists := cm.chats[userID]; !exists {
		cm.chats[userID] = make(map[string]openai.ChatCompletionNewParams)
	}
	cm.chats[userID][chatID] = openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemPrompt),
		},
	}
}

func (cm *ChatManager) GetChat(userID, chatID string) (openai.ChatCompletionNewParams, bool) {
	userChats, exist := cm.chats[userID]
	if !exist {
		return openai.ChatCompletionNewParams{}, false
	}
	chat, exist := userChats[chatID]
	return chat, exist
}

func(cm *ChatManager) AddMessage(userID, chatID string, message openai.ChatCompletionMessageParamUnion){
	chat, exist := cm.GetChat(userID, chatID)
	if !exist {
		return
	}
	chat.Messages = append(chat.Messages, message)
	cm.chats[userID][chatID] = chat
}

func (cm *ChatManager) GetConversation(userID, chatID string) []openai.ChatCompletionMessageParamUnion {
	chats, exist := cm.GetChat(userID, chatID)
	if !exist {
		return []openai.ChatCompletionMessageParamUnion{}
	}
	return chats.Messages
}

