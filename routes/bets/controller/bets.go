package controller

import (
	"gambler/backend/handlers"
	"gambler/backend/middleware"
	"gambler/backend/routes/bets/service"
	"time"

	"github.com/gofiber/fiber/v2"
)

func InitBetsRoute(c *fiber.App) {
	group := c.Group("/bets", middleware.JwtGuardHandler)
	group.Get("/", service.GetAllActiveBets)
	group.Post("/create", service.CreateBet)
	group.Get("/:id<int>", handlers.AddCache(time.Minute*15), service.GetBet)
}
