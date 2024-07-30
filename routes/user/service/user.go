package service

import (
	"gambler/backend/handlers"
	"gambler/backend/tools"

	"github.com/gofiber/fiber/v2"
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
	rawUserId := c.Params("id")
	userId := tools.ParseUInt(rawUserId)
	user, err := handlers.DB.GetUserByID(userId)
	if err != -1 {
		if err == tools.DB_REC_NOTFOUND {
			return c.Status(404).JSON(tools.GlobalErrorHandlerResp{
				Success: false,
				Message: "User not found",
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
	user.Password = ""
	user.Email = ""
	return c.Status(200).JSON(tools.GlobalErrorHandlerResp{
		Success: true,
		Message: "User found",
		Code:    200,
		Body:    user,
	})
}

func GetSelf(c *fiber.Ctx) error {
	userId, jwtErr := c.Locals("user").(jwt.Claims).GetSubject()
	if jwtErr != nil {
		tools.ClearAllCookies(c)
		return c.Status(401).JSON(tools.GlobalErrorHandlerResp{
			Success: false,
			Message: "Unauthorized",
			Code:    401,
		})
	}
	user, err := handlers.DB.GetUserByUsername(userId)
	if err != -1 {
		if err == tools.DB_REC_NOTFOUND {
			tools.ClearAllCookies(c)
			return c.Status(404).JSON(tools.GlobalErrorHandlerResp{
				Success: false,
				Message: "User not found",
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
	return c.Status(200).JSON(tools.GlobalErrorHandlerResp{
		Success: true,
		Message: "User found",
		Code:    200,
		Body:    user,
	})
}
