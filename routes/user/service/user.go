package service

import (
	"gambler/backend/handlers"
	"gambler/backend/tools"

	"github.com/gofiber/fiber/v2"
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
	return c.Status(200).JSON(tools.GlobalErrorHandlerResp{
		Success: true,
		Message: "User found",
		Code:    200,
		Body:    user,
	})
}
