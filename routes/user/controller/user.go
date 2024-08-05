package controller

import (
	"gambler/backend/middleware"
	"gambler/backend/routes/user/service"

	"github.com/gofiber/fiber/v2"
)

func InitUserRoute(c *fiber.App) {
	group := c.Group("/user", middleware.JwtGuardHandler)
	group.Get("/:id<int>", service.GetUserByID)
	group.Get("/@me", service.GetSelf)
}
