package controller

import (
	"gambler/backend/middleware"
	"gambler/backend/routes/bets/service"

	"github.com/gofiber/fiber/v2"
)

func InitBetsRoute(c *fiber.App) {
	group := c.Group("/bets", middleware.JwtGuardHandler)
	group.Get("/", service.GetAllActiveBets)
	group.Post("/create", service.CreateBet)
}
