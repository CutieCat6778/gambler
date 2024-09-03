package service

import (
	"gambler/backend/database/models"
	"gambler/backend/database/models/customTypes"
	"gambler/backend/handlers"
	"gambler/backend/tools"
	"math/rand"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lib/pq"
)

type (
	CreateBetReq struct {
		Name        string   `json:"name" validate:"required,min=3,max=50,ascii"`
		Description string   `json:"description" validate:"required,min=3,max=50,ascii"`
		BetOptions  []string `json:"betOptions" validate:"required,min=2,dive,min=3,max=50,ascii"`
		InputBet    float64  `json:"inputBet" validate:"required,min=1"`
		InputOption string   `json:"inputOption" validate:"required"`
	}
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

	err := handlers.DB.CreateBet(bet, userId, req.InputOption, req.InputBet)
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

func GetBet(c *fiber.Ctx) error {
	paramsId := c.Params("id")

	id, pErr := strconv.ParseUint(paramsId, 10, 32)
	if pErr != nil {
		return c.Status(400).JSON(tools.GlobalErrorHandlerResp{
			Success: false,
			Message: "Invalid id",
			Code:    400,
		})
	}

	bet, err := handlers.DB.GetBetByID(uint(id))
	if err != -1 {
		log.Info(tools.GetErrorString(err))
		return c.Status(500).JSON(tools.GlobalErrorHandlerResp{
			Success: false,
			Message: "Internal server error",
			Code:    500,
		})
	}

	if bet == nil {
		return c.Status(404).JSON(tools.GlobalErrorHandlerResp{
			Success: false,
			Message: "Bet not found",
			Code:    404,
		})
	}

	return c.Status(200).JSON(tools.GlobalErrorHandlerResp{
		Success: true,
		Message: "Bet found",
		Code:    200,
		Body:    bet,
	})
}
