package controller

import (
	"gambler/backend/middleware"
	"gambler/backend/routes/user/service"

	"github.com/gofiber/fiber/v2"
)

func InitUserRoute(c *fiber.App) {
	group := c.Group("/user", middleware.JwtGuardHandler)
	group.Get("/@me", service.GetSelf)
	group.Get("/balance", service.GetUserBalance)
	group.Get("/bets", service.GetUserBets)
	group.Get("/:name", service.GetUserByID)
}
