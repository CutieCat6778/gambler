package controller

import (
	"gambler/backend/middleware"
	"gambler/backend/routes/user/service"

	"github.com/gofiber/fiber/v2"
)

func InitUserRoute(c *fiber.App) {
	group := c.Group("/user")
	group.Get("/:id<int>", middleware.JwtGuardHandler, service.GetUserByID)
	group.Get("/@me", middleware.JwtGuardHandler, service.GetSelf)
}
