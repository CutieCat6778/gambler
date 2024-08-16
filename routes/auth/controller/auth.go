package controller

import (
	"gambler/backend/handlers"
	"gambler/backend/middleware"
	"gambler/backend/routes/auth/service"
	"time"

	"github.com/gofiber/fiber/v2"
)

func InitAuthRoute(c *fiber.App) {
	group := c.Group("/auth", handlers.AddCache(time.Hour*6))
	group.Post("/login", service.Login)
	group.Put("/register", service.Register)
	group.Get("/refresh", service.RefreshToken)
	group.Get("/ping", middleware.JwtGuardHandler, service.Ping)
}
