package service

import (
	"gambler/backend/handlers"
	"gambler/backend/tools"

	"github.com/gofiber/fiber/v2"
)

type (
	AddBalanceReq struct {
		Amount float64 `json:"amount" validate:"required,min=1"`
		Reason string  `json:"reason" validate:"required,min=1,max=3,ascii"`
		UserId string  `json:"user_id" validate:"required,min=1"`
	}
)

func AddBalanceToUser(c *fiber.Ctx) error {
	var req AddBalanceReq

	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(tools.GlobalErrorHandlerResp{
			Success: false,
			Message: "[Parser] Bad request",
			Code:    400,
		})
	}

	if errs := handlers.VHandler.Validate(req); len(errs) > 0 && errs[0].Error {
		return c.Status(400).JSON(tools.GlobalErrorHandlerResp{
			Success: false,
			Message: "[Validator] Bad request",
			Code:    400,
			Body:    errs,
		})
	}

	user, err := handlers.DB.GetUserByUsername(req.UserId)
	if err != -1 {
		if err == tools.DB_REC_NOTFOUND {
			return c.Status(404).JSON(tools.GlobalErrorHandlerResp{
				Success: false,
				Message: "Bet not found",
				Code:    404,
			})
		} else {
			return c.Status(500).JSON(tools.GlobalErrorHandlerResp{
				Success: false,
				Message: "Internal server error",
				Code:    500,
			})
		}
	}
	newAmount := user.Balance + req.Amount
	err = handlers.DB.UpdateUserBalance(newAmount, *user, req.Reason)
	if err != -1 {
		return c.Status(500).JSON(tools.GlobalErrorHandlerResp{
			Success: false,
			Message: "Internal server error",
			Code:    500,
		})
	}
	return c.Status(200).JSON(tools.GlobalErrorHandlerResp{
		Success: true,
		Message: "Balance added",
		Code:    200,
	})
}
