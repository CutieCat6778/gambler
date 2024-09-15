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
	group.Get("/", handlers.AddCache(time.Minute*3), service.GetAllBetsHandler)
	group.Post("/create", service.CreateBet)
	group.Get("/:id<int>", handlers.AddCache(time.Second*10), service.GetBet)
	group.Put("/place/:id<int>", service.PlaceBet)
	group.Put("/remove/:id<int>", service.PlaceBet)
}
