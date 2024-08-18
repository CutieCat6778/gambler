package tools

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/joho/godotenv"
)

type (
	GlobalErrorHandlerResp struct {
		Success bool        `json:"success"`
		Message string      `json:"message"`
		Code    int         `json:"code"`
		Body    interface{} `json:"body,omitempty"`
	}
)

var (
	DATABASE          string
	JWT_SECRET        []byte
	HASH_SECRET       string
	COOKIE_SECRET     string
	HOST_REDIS        string
	PSW_REDIS         string
	URL_REDIS         string
	WEBSOCKET_VERSION byte
	MASTER_IDS        string
)

func InitEnvVars() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	DATABASE = os.Getenv("POSTGRES_DB")
	JWT_SECRET = []byte(os.Getenv("JWT_SECRET"))
	HASH_SECRET = os.Getenv("HASH_SECRET")
	COOKIE_SECRET = os.Getenv("COOKIE_SECRET")
	HOST_REDIS = os.Getenv("REDIS_HOST")
	PSW_REDIS = os.Getenv("REDIS_PSW")
	URL_REDIS = os.Getenv("REDIS_URL")
	ver, err := strconv.ParseInt(os.Getenv("WEBSOCKET_VERSION"), 10, 64)
	if err != nil {
		panic(err)
	}
	WEBSOCKET_VERSION = byte(ver)
	MASTER_IDS = os.Getenv("MASTER_IDS")
	// Check for missing variables and log them
	missingVars := []string{}
	if DATABASE == "" {
		missingVars = append(missingVars, "POSTGRES_DB")
	}
	if len(JWT_SECRET) == 0 {
		missingVars = append(missingVars, "JWT_SECRET")
	}
	if HASH_SECRET == "" {
		missingVars = append(missingVars, "HASH_SECRET")
	}
	if COOKIE_SECRET == "" {
		missingVars = append(missingVars, "COOKIE_SECRET")
	}
	if HOST_REDIS == "" {
		missingVars = append(missingVars, "REDIS_HOST")
	}
	if PSW_REDIS == "" {
		missingVars = append(missingVars, "REDIS_PSW")
	}
	if URL_REDIS == "" {
		missingVars = append(missingVars, "REDIS_URL")
	}
	if WEBSOCKET_VERSION == 0 {
		missingVars = append(missingVars, "WEBSOCKET_VERSION")
	}
	if MASTER_IDS == "" {
		missingVars = append(missingVars, "MASTER_IDS")
	}

	// If there are any missing variables, log them and panic
	if len(missingVars) > 0 {
		log.Fatalf("[ENV] The following environment variables are missing: %v", missingVars)
		panic("Some Environment Variables are missing")
	}
	fmt.Println("[ENV] Loaded Enviroment Variables")
	fmt.Println(DATABASE)
}

func ParseUInt(s string) uint {
	var n uint
	fmt.Sscanf(s, "%d", &n)
	return n
}

func AddCacheTime(c *fiber.Ctx, duration time.Duration) {
	c.Response().Header.Add("Cache-Time", fmt.Sprintf("%d", int(duration.Seconds())))
}

func ConfigureApp(app *fiber.App) {
	app.Use(healthcheck.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3001",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	app.Use(limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.IP() == "127.0.0.1"
		},
		Max:        20,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Get("x-forwarded-for")
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(429).JSON(GlobalErrorHandlerResp{
				Success: false,
				Message: "Too many requests",
				Code:    429,
			})
		},
	}))
}

func HeaderParser(c *fiber.Ctx) string {
	headers := c.GetReqHeaders()
	log.Info(headers)
	if len(headers["Authorization"]) == 0 || headers["Authorization"] == nil {
		return ""
	}
	rawBearer := headers["Authorization"][0]
	if !strings.HasPrefix(rawBearer, "Bearer ") {
		return ""
	}

	token := strings.Split(rawBearer, "Bearer ")[1]
	return token
}
