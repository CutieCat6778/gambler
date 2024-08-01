package controller

import (
	"gambler/backend/middleware"
	"gambler/backend/routes/auth/service"

	"github.com/gofiber/fiber/v2"
)

func InitAuthRoute(c *fiber.App) {
	group := c.Group("/auth")
	group.Post("/login", service.Login)
	group.Put("/register", service.Register)
	group.Get("/ping", middleware.JwtGuardHandler, service.Ping)
}
