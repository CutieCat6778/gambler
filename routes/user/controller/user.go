package controller

import (
	"gambler/backend/handlers"
	"gambler/backend/middleware"
	"gambler/backend/routes/user/service"
	"time"

	"github.com/gofiber/fiber/v2"
)

func InitUserRoute(c *fiber.App) {
	group := c.Group("/user", middleware.JwtGuardHandler)
	group.Get("/@me", handlers.AddCache(time.Second*5), service.GetSelf)
	group.Get("/balance", service.GetUserBalance)
	group.Get("/bets", service.GetUserBets)
	group.Get("/:name", service.GetUserByID)
}
