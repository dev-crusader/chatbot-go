package main

import (
	"github.com/dev-crusader/chatbot-go/controller"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	engine := html.New("./", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("static/style.css", "./static/style.css")

	store := session.New()

	chatController := controller.NewController(store)

	app.Get("/", func(c *fiber.Ctx) error {
		chatController.StartChat(c)
		return c.Render("chat", fiber.Map{})
	})
	app.Post("/create", chatController.CreateChat)
	app.Post("/send", chatController.SendMessage)

	app.Listen(":3000")
}