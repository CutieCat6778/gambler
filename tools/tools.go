package tools

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
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
)

var (
	DATABASE      string
	JWT_SECRET    []byte
	HASH_SECRET   string
	COOKIE_SECRET string
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
	fmt.Println("[ENV] Loaded Enviroment Variables")
	fmt.Println(DATABASE)
}

func ParseUInt(s string) uint {
	var n uint
	fmt.Sscanf(s, "%d", &n)
	return n
}

func ConfigureApp(app *fiber.App) {
	app.Use(encryptcookie.New(encryptcookie.Config{
		Key: encryptcookie.GenerateKey(),
	}))
	app.Use(cors.New())
	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowCredentials: true,
		AllowOriginsFunc: func(origin string) bool {
			fmt.Println(origin)
			return true
		},
		AllowOrigins: "http://localhost:3001",
	}))
}

func ClearAllCookies(c *fiber.Ctx) {
	c.ClearCookie("accesstoken")
	c.ClearCookie("refreshtoken")
	c.ClearCookie("username")
}

func SetCookieAfterAuth(c *fiber.Ctx, accessToken string, refreshToken string, username string) {
	c.Cookie(&fiber.Cookie{
		Name:     "username",
		Value:    username,
		Path:     "/",
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteLaxMode,
	})
	c.Cookie(&fiber.Cookie{
		Name:     "accesstoken",
		Value:    accessToken,
		MaxAge:   86400,
		Path:     "/",
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteLaxMode,
	})
	c.Cookie(&fiber.Cookie{
		Name:     "refreshtoken",
		Value:    refreshToken,
		MaxAge:   86400 * 7,
		Path:     "/",
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteLaxMode,
	})
	c.Cookie(&fiber.Cookie{
		Name:     "lastlogin",
		Value:    strconv.FormatInt(time.Now().Unix(), 10),
		Path:     "/",
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteLaxMode,
	})
}
