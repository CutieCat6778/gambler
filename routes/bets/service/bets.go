package service

import (
	"gambler/backend/handlers"
	"gambler/backend/tools"

	"github.com/gofiber/fiber/v2"
)

func GetAllActiveBets(c *fiber.Ctx) error {
	bets, err := handlers.Cache.GetAllBet()
	if err != -1 {
		return c.Status(500).JSON(tools.GlobalErrorHandlerResp{
			Success: false,
			Message: "Internal server error",
			Code:    500,
		})
	}
	return c.Status(200).JSON(tools.GlobalErrorHandlerResp{
		Success: true,
		Message: "Active bets found",
		Code:    200,
		Body:    bets,
	})
}
