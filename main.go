package main

import (
	"gambler/backend/handlers"
	authController "gambler/backend/routes/auth/controller"
	betsController "gambler/backend/routes/bets/controller"
	rootController "gambler/backend/routes/root/controller"
	userController "gambler/backend/routes/user/controller"
	wsController "gambler/backend/routes/ws/controller"
	"gambler/backend/tools"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func init() {
	tools.InitEnvVars()
}

func main() {
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(tools.GlobalErrorHandlerResp{
				Success: false,
				Message: err.Error(),
				Code:    code,
			})
		},
	})

	tools.ConfigureApp(app)

	_ = handlers.NewDB()
	_ = handlers.NewValidator()
	cache := handlers.NewCache(app)
	_ = handlers.NewWebSocketHandler(cache)

	log.SetLevel(log.LevelInfo)

	userController.InitUserRoute(app)
	authController.InitAuthRoute(app)
	wsController.InitWsRoute(app)
	betsController.InitBetsRoute(app)
	rootController.InitRootRoute(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(tools.GlobalErrorHandlerResp{
			Success: true,
			Message: "Welcome to Gambler API",
			Code:    200,
		})
	})

	app.Listen(":3000")
}
