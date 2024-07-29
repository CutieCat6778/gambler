package main

import (
	"gambler/backend/handlers"
	authController "gambler/backend/routes/auth/controller"
	userController "gambler/backend/routes/user/controller"
	"gambler/backend/tools"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

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

	// app.Use(cors.New())
	// app.Use(csrf.New())

	// app.Use(cors.New(cors.Config{
	// 	AllowOrigins: "http://localhost:3000",
	// 	AllowHeaders: "Origin, Content-Type, Accept",
	// }))

	// app.Use(csrf.New(csrf.Config{
	// 	KeyLookup:      "header:X-Csrf-Token",
	// 	CookieName:     "csrf_",
	// 	CookieSameSite: "Lax",
	// 	Expiration:     1 * time.Hour,
	// 	KeyGenerator:   utils.UUIDv4,
	// }))

	_ = handlers.NewDB()
	_ = handlers.NewValidator()

	userController.InitUserRoute(app)
	authController.InitAuthRoute(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(tools.GlobalErrorHandlerResp{
			Success: true,
			Message: "Welcome to Gambler API",
			Code:    200,
		})
	})

	app.Listen(":3000")
}
