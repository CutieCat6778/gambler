package service

import (
	"fmt"
	"gambler/backend/database/models"
	"gambler/backend/handlers"
	"gambler/backend/middleware"
	"gambler/backend/tools"
	"strings"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type (
	LoginReq struct {
		Username string `json:"username" validate:"required,min=3,max=20,alphanum"`
		Password string `json:"password" validate:"required,min=8,ascii,excludes=:"`
	}

	RegisterReq struct {
		Username string `json:"username" validate:"required,min=3,max=20,alphanum"`
		Password string `json:"password" validate:"required,min=8,ascii,excludes=:"`
		Email    string `json:"email" validate:"required,email"`
		Name     string `json:"name" validate:"required,min=3,max=50,ascii"`
	}

	LoginRes struct {
		User *models.User  `json:"user"`
		Bets *[]models.Bet `json:"bets"`
	}
)

func Login(c *fiber.Ctx) error {
	req := new(LoginReq)

	if err := c.BodyParser(req); err != nil {
		tools.ReturnData(c, 400, nil, -1)
	}

	if errs := handlers.VHandler.Validate(req); len(errs) > 0 && errs[0].Error {
		return tools.ReturnData(c, 400, errs, -1)
	}

	user, err := handlers.DB.GetUserByUsername(req.Username)
	if err != -1 {
		if err == tools.DB_REC_NOTFOUND {
			return tools.ReturnData(c, 404, nil, tools.DB_REC_NOTFOUND)
		} else {
			return c.Status(500).JSON(tools.GlobalErrorHandlerResp{
				Success: false,
				Message: "Internal server error",
				Code:    500,
			})
		}
	}
	hashedPassword := user.Password
	valid := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(tools.HASH_SECRET+":"+req.Password))
	if valid != nil {
		fmt.Println(valid)
		return tools.ReturnData(c, 401, strings.Split(valid.Error(), ": ")[1], -1)
	}
	tokens, err := middleware.Sign(user.ID)
	if err != -1 {
		return tools.ReturnData(c, 500, nil, err)
	}

	bets := &[]models.Bet{}
	bets, err = handlers.Cache.GetAllBet()
	if err != -1 {
		return tools.ReturnData(c, 500, nil, err)
	}

	tools.CreateCookie(c, tokens.AccessToken, tokens.RefreshToken, user.ID)
	return tools.ReturnData(c, 200, LoginRes{
		User: user,
		Bets: bets,
	}, -1)
}

func RefreshToken(c *fiber.Ctx) error {
	header := c.Cookies("refresh_token")

	claims, err := middleware.Decode(header, true)
	if err != -1 {
		return tools.ReturnData(c, 401, nil, err)
	}

	userId, jwtErr := claims.GetSubject()
	if jwtErr != nil {
		return tools.ReturnData(c, 401, nil, tools.JWT_INVALID)
	}

	tokens, err := middleware.Sign(tools.ParseUInt(userId))
	if err != -1 {
		return tools.ReturnData(c, 500, nil, err)
	}

	return tools.ReturnData(c, 200, tokens, -1)
}

func Register(c *fiber.Ctx) error {
	req := new(RegisterReq)

	if err := c.BodyParser(req); err != nil {
		fmt.Println(err)
		return tools.ReturnData(c, 400, nil, -1)
	}

	if errs := handlers.VHandler.Validate(req); len(errs) > 0 && errs[0].Error {
		return tools.ReturnData(c, 400, errs, -1)
	}

	hashedPasssword, err := bcrypt.GenerateFromPassword([]byte(tools.HASH_SECRET+":"+req.Password), bcrypt.MinCost*2)
	if err != nil {
		fmt.Println(err)
		return tools.ReturnData(c, 500, nil, -1)
	}

	user := models.User{
		Username: req.Username,
		Password: string(hashedPasssword),
		Email:    req.Email,
		Name:     req.Name,
		UserBet:  []models.UserBet{},
	}

	dbErr := handlers.DB.CreateUser(user)
	if dbErr != -1 {
		return tools.ReturnData(c, 500, nil, dbErr)
	}

	bets := &[]models.Bet{}
	bets, dbErr = handlers.Cache.GetAllBet()
	if dbErr != -1 {
		return tools.ReturnData(c, 500, nil, dbErr)
	}

	return tools.ReturnData(c, 200, LoginRes{
		User: &user,
		Bets: bets,
	}, -1)
}

func Ping(c *fiber.Ctx) error {
	return tools.ReturnData(c, 200, "Pong!", -1)
}
