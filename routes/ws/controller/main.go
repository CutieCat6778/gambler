package controller

import (
	"gambler/backend/routes/ws/service"

	W "gambler/backend/handlers/websocket"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func InitWsRoute(c *fiber.App) {
	c.Use(service.CompatibleCheck)
	c.Get("/ws/:id", websocket.New(W.WebSocket.HandleWebSocketConnection))
}
