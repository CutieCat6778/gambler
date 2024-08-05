package controller

import (
	"gambler/backend/routes/ws/service"

	"github.com/gofiber/contrib/socketio"
	"github.com/gofiber/fiber/v2"
)

func InitWsRoute(c *fiber.App) {
	c.Use(service.CompatibleCheck)
	socketio.On(socketio.EventConnect, service.OnHandshake)
	c.Get("/ws/:id/:game_id", socketio.New(service.UserConnect))
}
