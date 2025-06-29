package controller

import (
	cs "github.com/dev-crusader/chatbot-go/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/google/uuid"
)

type ChatController struct {
	store *session.Store
	chatService *cs.ChatService
}

func NewController(store *session.Store) *ChatController {
	return &ChatController{
		store: store,
		chatService: cs.NewChatService(),
	}
}

func(cc *ChatController) StartChat(ctx *fiber.Ctx) error {
	sess, _ := cc.store.Get(ctx)
	if sess.Get("user_id") == nil {
		sess.Set("user_id", uuid.New().String())
		sess.Save()
	}
	return ctx.SendString("Welcome to Chat Bot Service!")
}

func (cc *ChatController) CreateChat(ctx *fiber.Ctx) error {
	sess, _ := cc.store.Get(ctx)
    userID := sess.Get("user_id")
	if userID == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error":"Session expired"})
	}
    chatID := cc.chatService.CreateChat(userID.(string))
    return ctx.JSON(fiber.Map{
		"chat_id": chatID,
		"message":"Chat created successfully!",
	})
}

func(cc *ChatController) SendMessage(ctx *fiber.Ctx) error {
	sess, _ := cc.store.Get(ctx)
	userID := sess.Get("user_id")
	if userID == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error":"Session expired"})
	}

	var data struct {
		ChatID		string `json:"chat_id"`
		UserMessage string `json:"message"`
	}

	if err := ctx.BodyParser(&data); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"Invalid request"})
	}
    if data.ChatID == "" || data.UserMessage == "" {
        return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"Missing ChatID or UserMessage"})
    }
    
    response, err := cc.chatService.ProcessMessage(userID.(string), data.ChatID, data.UserMessage)
    if err != nil {
        return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }
    return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": response})
}