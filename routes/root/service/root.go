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
		return tools.ReturnData(c, 400, nil, -1)
	}

	if errs := handlers.VHandler.Validate(req); len(errs) > 0 && errs[0].Error {
		return tools.ReturnData(c, 400, errs, -1)
	}

	user, err := handlers.DB.GetUserByID(tools.ParseUInt(req.UserId))
	if err != -1 {
		return tools.ReturnData(c, 500, nil, err)
	}
	newAmount := user.Balance + req.Amount
	err = handlers.DB.UpdateUserBalance(newAmount, *user, req.Reason)
	if err != -1 {
		return tools.ReturnData(c, 500, nil, err)
	}
	return tools.ReturnData(c, 200, nil, -1)
}
