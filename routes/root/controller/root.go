package controller

import (
	"gambler/backend/middleware"
	"gambler/backend/routes/root/service"

	"github.com/gofiber/fiber/v2"
)

func InitRootRoute(c *fiber.App) {
	group := c.Group("/", middleware.JwtGuardMasterHandler)
	group.Put("/bets/create", service.CreateBet)
}
