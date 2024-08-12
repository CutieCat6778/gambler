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

const (
	// DATABASE ERRORS
	DB_UNKOWN_ERR = iota
	DB_REC_NOTFOUND
	DB_DUP_KEY

	// JWT ERRORS
	JWT_FAILED_TO_SIGN
	JWT_FAILED_TO_DECODE
	JWT_INVALID
	JWT_EXPIRED

	// WS ERROS
	WS_UUID_DUP
	WS_UUID_NOTFOUND
	WS_GAMEID_NOTFOUND
	WS_INVALID_CONN
	WS_UNKNOWN_ERR

	// REDIS ERRORS
	RD_CONN_CLOSED
	RD_KEY_NOT_FOUND
	RD_TX_FAILED
	RD_UNKNOWN

	// GENERAL ERRORS
	JSON_UNMARSHAL_ERROR
)

var (
	DATABASE           string
	JWT_SECRET         []byte
	HASH_SECRET        string
	COOKIE_SECRET      string
	HOST_REDIS         string
	PSW_REDIS          string
	WEBSOCKET_VERSEION byte
)

func init() {
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
	ver, err := strconv.ParseInt(os.Getenv("WEBSOCKET_VERSION"), 10, 64)
	if err != nil {
		panic(err)
	}
	WEBSOCKET_VERSEION = byte(ver)
	fmt.Println("[ENV] Loaded Enviroment Variables")
	fmt.Println(DATABASE)
}

func ParseUInt(s string) uint {
	var n uint
	fmt.Sscanf(s, "%d", &n)
	return n
}

func ConfigureApp(app *fiber.App) {
	// app.Use(cors.New(cors.Config{
	// 	AllowHeaders: "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin,Authorization",
	// 	AllowMethods: "DELETE, POST, GET, PUT, OPTIONS",
	// 	AllowOriginsFunc: func(origin string) bool {
	// 		log.Info(origin)
	// 		return strings.Contains(origin, "localhost")
	// 	},
	// }))
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
