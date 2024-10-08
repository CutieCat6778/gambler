package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
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
		AllowOrigins:     "http://localhost:4200, http://192.168.178.2",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
		AllowOriginsFunc: func(origin string) bool {
			log.Info(origin)
			return true
		},
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

func CreateCookie(c *fiber.Ctx, accessToken string, refreshToken string, userId uint) int {
	rtC := &fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(24 * time.Hour * 7),
		HTTPOnly: true,
		Secure:   false,
		SameSite: "lax",
	}
	atC := &fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(15 * time.Minute),
		Secure:   false,
		SameSite: "lax",
		HTTPOnly: true,
	}
	udC := &fiber.Cookie{
		Name:     "user_id",
		Value:    fmt.Sprintf("%d", userId),
		Expires:  time.Now().Add(24 * time.Hour * 7),
		Secure:   false,
		SameSite: "lax",
		HTTPOnly: true,
	}
	c.Cookie(rtC)
	c.Cookie(atC)
	c.Cookie(udC)
	return -1
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

func ParseTimestamp(timestamp string) time.Time {
	t, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		log.Error(err)
	}
	return t
}

type Payload struct {
	Content string `json:"content"`
}

func SendWebHook(err string) int {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		log.Info(fmt.Sprintf("Called from %s, line %d", file, line))
	}
	webhookURL := "https://discordapp.com/api/webhooks/1274463723960143883/__YvfQkphIcyetuB0VBtS3RysKraGv2LORHolSyfMXDvmWuMFwVcEqXB4Hj7A0ZM5Hh4"

	// Create the payload
	payload := Payload{Content: fmt.Sprintf("----\n**Error:** %s\n**File:** %s\n**Line:** %d", err, file, line)}

	// Encode payload into JSON
	payloadBytes, jErr := json.Marshal(payload)
	if jErr != nil {
		log.Error(jErr)
		return JSON_MARSHAL_ERROR
	}

	// Send the POST request
	resp, hErr := http.Post(webhookURL, "application/json", bytes.NewBuffer(payloadBytes))
	if hErr != nil {
		log.Errorf("failed to send webhook: %v", err)
		return WEBHOOK_ERROR
	}
	defer resp.Body.Close()

	// Check for non-200 HTTP status codes
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		log.Errorf("unexpected status code: %d", resp.StatusCode)
		return WEBHOOK_ERROR
	}

	return -1
}

func ConvertKeyToBetID(key string) uint {
	return ParseUInt(strings.TrimPrefix(key, "b-"))
}

func ReturnData(c *fiber.Ctx, code int, body interface{}, err int) error {
	if err != -1 {
		return c.Status(code).JSON(GlobalErrorHandlerResp{
			Success: false,
			Message: GetErrorString(err),
			Code:    err,
			Body:    body,
		})
	} else if code >= 400 {
		return c.Status(code).JSON(GlobalErrorHandlerResp{
			Success: false,
			Message: StatusText(code),
			Code:    code,
			Body:    body,
		})
	}
	return c.Status(code).JSON(GlobalErrorHandlerResp{
		Success: true,
		Message: "Success",
		Code:    200,
		Body:    body,
	})
}

func StatusText(code int) string {
	switch code {
	case http.StatusOK:
		return "OK"
	case http.StatusCreated:
		return "Created"
	case http.StatusAccepted:
		return "Accepted"
	case http.StatusNoContent:
		return "No Content"
	case http.StatusBadRequest:
		return "Bad Request"
	case http.StatusUnauthorized:
		return "Unauthorized"
	case http.StatusForbidden:
		return "Forbidden"
	case http.StatusNotFound:
		return "Not Found"
	case http.StatusInternalServerError:
		return "Internal Server Error"
	case http.StatusNotImplemented:
		return "Not Implemented"
	case http.StatusBadGateway:
		return "Bad Gateway"
	case http.StatusServiceUnavailable:
		return "Service Unavailable"
	// Add more status codes as needed
	default:
		return "Unknown Status Code"
	}
}

func Contains(slice []string, item string) bool {
	var res bool = false
	for _, str := range slice {
		log.Info(str, item)
		if str == item {
			log.Info("Found")
			res = true
			break
		}
	}
	return res
}
