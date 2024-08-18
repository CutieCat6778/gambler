package controller

import (
	"gambler/backend/handlers"
	"gambler/backend/routes/ws/service"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func InitWsRoute(c *fiber.App) {
	c.Use(service.CompatibleCheck)
	c.Get("/ws/:id", websocket.New(handlers.WebSocket.HandleWebSocketConnection))
}
