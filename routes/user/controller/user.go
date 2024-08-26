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
	group.Get("/:id<int>", handlers.AddCache(time.Minute*15), service.GetUserByID)
	group.Get("/@me", service.GetSelf)
	group.Get("/balance", service.GetUserBalance)
	group.Get("/bets", service.GetUserBets)
}
