package service

import (
	"gambler/backend/database/models"
	"gambler/backend/database/models/customTypes"
	"gambler/backend/handlers"
	"gambler/backend/tools"
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lib/pq"
)

type (
	CreateBetReq struct {
		Name        string   `json:"name" validate:"required,min=3,max=50,ascii"`
		Description string   `json:"description" validate:"required,min=3,max=50,ascii"`
		BetOptions  []string `json:"betOptions" validate:"required,min=2,dive,min=3,max=50,ascii"`
	}
	AddBalanceReq struct {
		Amount float64 `json:"amount" validate:"required,min=1"`
		Reason string  `json:"reason" validate:"required,min=1,max=3,ascii"`
		UserId string  `json:"user_id" validate:"required,min=1"`
	}
)

func CreateBet(c *fiber.Ctx) error {
	req := new(CreateBetReq)

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

	userId, jwtErr := c.Locals("claims").(jwt.Claims).GetSubject()
	if jwtErr != nil {
		return c.Status(401).JSON(tools.GlobalErrorHandlerResp{
			Success: false,
			Message: "Unauthorized",
			Code:    401,
		})
	}

	rand.NewSource(time.Now().UnixNano())

	bet := models.Bet{
		Name:        req.Name,
		Description: req.Description,
		BetOptions:  pq.StringArray(req.BetOptions),
		Status:      customTypes.Open,
	}

	err := handlers.DB.CreateBet(bet, userId)
	if err != -1 {
		if err == tools.DB_DUP_KEY {
			return c.Status(409).JSON(tools.GlobalErrorHandlerResp{
				Success: false,
				Message: "User already exists",
				Code:    409,
			})
		} else {
			return c.Status(500).JSON(tools.GlobalErrorHandlerResp{
				Success: false,
				Message: "Internal server error",
				Code:    500,
			})
		}
	}

	return c.Status(201).JSON(tools.GlobalErrorHandlerResp{
		Success: true,
		Message: "Bet created",
		Code:    201,
		Body:    bet,
	})
}

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
