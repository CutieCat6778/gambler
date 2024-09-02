package controller

import (
	"gambler/backend/middleware"
	"gambler/backend/routes/root/service"

	"github.com/gofiber/fiber/v2"
)

func InitRootRoute(c *fiber.App) {
	group := c.Group("/s", middleware.JwtGuardMasterHandler)
	group.Put("/user/balance", service.AddBalanceToUser)
}
