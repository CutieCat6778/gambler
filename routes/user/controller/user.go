package controller

import (
	"gambler/backend/middleware"
	"gambler/backend/routes/user/service"
	"gambler/backend/tools"

	"github.com/gofiber/fiber/v2"
)

func InitUserRoute(c *fiber.App) {
	group := c.Group("/user")
	group.Get("/:id<int>", middleware.JwtGuardHandler, service.GetUserByID)
	group.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(tools.GlobalErrorHandlerResp{
			Success: true,
			Message: "User",
			Code:    200,
		})
	})
}
