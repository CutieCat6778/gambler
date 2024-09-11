package service

import (
	"gambler/backend/handlers"
	"gambler/backend/tools"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
)

type (
	CreateUserReq struct {
		Username string `json:"username" validate:"required,min=3,max=20,alphanum"`
		Password string `json:"password" validate:"required,min=8,ascii,excludes=:"`
		Email    string `json:"email" validate:"required,email"`
		Name     string `json:"name" validate:"required,min=3,max=50,ascii"`
	}
)

func GetUserByID(c *fiber.Ctx) error {
	userId := c.Params("id")
	user, err := handlers.DB.GetUserByID(tools.ParseUInt(userId))
	log.Info(user)
	if err != -1 {
		return tools.ReturnData(c, 500, nil, err)
	}
	user.Password = ""
	user.Email = ""
	return tools.ReturnData(c, 200, user, -1)
}

func GetSelf(c *fiber.Ctx) error {
	claims := c.Locals("claims").(jwt.Claims)
	if claims == nil {
		return tools.ReturnData(c, 500, nil, -1)
	}
	log.Info(claims)
	userId, jwtErr := claims.GetSubject()
	if jwtErr != nil {
		return tools.ReturnData(c, 401, nil, -1)
	}
	log.Info(userId)
	user, err := handlers.DB.GetUserByID(tools.ParseUInt(userId))
	if err != -1 {
		return tools.ReturnData(c, 500, nil, err)
	}

	activeBets, err := handlers.Cache.GetAllBet()
	if err != -1 {
		return tools.ReturnData(c, 500, nil, err)
	}

	return tools.ReturnData(c, 200, fiber.Map{
		"user": user,
		"bets": activeBets,
	}, -1)
}

func GetUserBalance(c *fiber.Ctx) error {
	userId, jwtErr := c.Locals("claims").(jwt.Claims).GetSubject()
	if jwtErr != nil {
		return tools.ReturnData(c, 401, nil, -1)
	}
	balance, err := handlers.DB.FindBalanceHistoryByUser(tools.ParseUInt(userId))
	if err != -1 {
		if err == tools.DB_REC_NOTFOUND {
			return tools.ReturnData(c, 404, nil, tools.DB_REC_NOTFOUND)
		} else {
			return tools.ReturnData(c, 500, nil, err)
		}
	}
	return tools.ReturnData(c, 200, balance, -1)
}

func GetUserBets(c *fiber.Ctx) error {
	userId, jwtErr := c.Locals("claims").(jwt.Claims).GetSubject()
	if jwtErr != nil {
		return tools.ReturnData(c, 401, nil, -1)
	}
	bets, err := handlers.DB.GetUserBet(tools.ParseUInt(userId))
	if err != -1 {
		if err == tools.DB_REC_NOTFOUND {
			return tools.ReturnData(c, 404, nil, tools.DB_REC_NOTFOUND)
		} else {
			return tools.ReturnData(c, 500, nil, err)
		}
	}
	return tools.ReturnData(c, 200, bets, -1)
}
