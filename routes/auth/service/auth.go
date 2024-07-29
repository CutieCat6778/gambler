package service

import (
	"fmt"
	"gambler/backend/database/models"
	"gambler/backend/handlers"
	"gambler/backend/middleware"
	"gambler/backend/tools"

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
)

func Login(c *fiber.Ctx) error {
	req := new(LoginReq)

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

	return c.Status(200).JSON(tools.GlobalErrorHandlerResp{
		Success: true,
		Message: "Login success",
		Code:    200,
	})
}

func Register(c *fiber.Ctx) error {
	req := new(RegisterReq)

	if err := c.BodyParser(req); err != nil {
		fmt.Println(err)
		return c.Status(400).JSON(tools.GlobalErrorHandlerResp{
			Success: false,
			Message: "[Parser] Bad request: " + err.Error(),
			Code:    400,
			Body:    err,
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

	hashedPasssword, err := bcrypt.GenerateFromPassword([]byte(tools.HASH_SECRET+":"+req.Password), bcrypt.MinCost*2)
	if err != nil {
		fmt.Println(err)
		return c.Status(500).JSON(tools.GlobalErrorHandlerResp{
			Success: false,
			Message: "[Hash] Internal server error",
			Code:    500,
		})
	}

	user := models.User{
		Username: req.Username,
		Password: hashedPasssword,
		Email:    req.Email,
		Name:     req.Name,
	}

	res, dbErr := handlers.DB.CreateUser(user)
	if dbErr != -1 {
		if dbErr == tools.DB_DUP_KEY {
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

	tokens, jwtErr := middleware.Sign(res.Username, res.ID)
	if jwtErr != -1 {
		if jwtErr == tools.JWT_FAILED_TO_SIGN {
			return c.Status(500).JSON(tools.GlobalErrorHandlerResp{
				Success: false,
				Message: "Internal server error, failed to sign key",
				Code:    500,
			})
		} else {
			return c.Status(500).JSON(tools.GlobalErrorHandlerResp{
				Success: false,
				Message: "Internal server error",
				Code:    500,
			})
		}
	}

	return c.Status(200).JSON(tools.GlobalErrorHandlerResp{
		Success: true,
		Message: "Register success",
		Code:    200,
		Body:    tokens,
	})
}
