package middleware

import (
	"fmt"
	"gambler/backend/handlers"
	"gambler/backend/tools"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
)

type Jwt struct {
	AccessToken         string    `json:"accessToken"`
	RefreshToken        string    `json:"refreshToken"`
	AccessTokenExpDate  time.Time `json:"accessTokenExpDate"`
	RefreshTokenExpDate time.Time `json:"refreshTokenExpDate"`
}

func Sign(username string) (*Jwt, int) {
	user, dbErr := handlers.DB.GetUserByUsername(username)
	if dbErr != -1 {
		return nil, dbErr
	}
	accessTokenExpDate := time.Minute * 15
	refreshTokenExpDate := 24 * 7 * time.Hour
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenExpDate)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer:    "Gambler Backend Service",
		Subject:   username,
	})
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenExpDate)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer:    fmt.Sprintf("%d Version", user.RefreshTokenVersion),
		Subject:   username,
	})
	AccessToken, err := accessToken.SignedString(tools.JWT_SECRET)
	if err != nil {
		return nil, tools.JWT_FAILED_TO_SIGN
	}
	RefreshToken, err := refreshToken.SignedString(tools.JWT_SECRET)
	if err != nil {
		return nil, tools.JWT_FAILED_TO_SIGN
	}
	return &Jwt{
		AccessToken:         AccessToken,
		RefreshToken:        RefreshToken,
		AccessTokenExpDate:  time.Now().Add(accessTokenExpDate),
		RefreshTokenExpDate: time.Now().Add(refreshTokenExpDate),
	}, -1
}

func Decode(token string, isRefresh bool) (jwt.Claims, int) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return tools.JWT_SECRET, nil
	})
	if err != nil {
		log.Info(err)
		return nil, tools.JWT_FAILED_TO_DECODE
	}
	if !t.Valid {
		rawExpireIn := t.Claims.(jwt.MapClaims)["exp"]
		expireIn := time.Now().Unix() - rawExpireIn.(int64)
		if expireIn > 0 {
			return nil, tools.JWT_EXPIRED
		} else {
			return nil, tools.JWT_INVALID
		}
	}

	if isRefresh {
		issuer, jwtErr := t.Claims.GetIssuer()
		if jwtErr != nil {
			return nil, tools.JWT_FAILED_TO_DECODE
		}
		if !strings.Contains(issuer, "Version") {
			return nil, tools.JWT_INVALID
		}
		ver, err := strconv.ParseInt(strings.Split(issuer, " ")[0], 10, 0)
		if err != nil {
			return nil, tools.JWT_INVALID
		}
		username, jwtErr := t.Claims.GetSubject()
		if jwtErr != nil {
			return nil, tools.JWT_FAILED_TO_DECODE
		}

		user, dbErr := handlers.DB.GetUserByUsername(username)
		if dbErr != -1 {
			return nil, dbErr
		}

		if user.RefreshTokenVersion != int(ver) {
			return nil, tools.JWT_INVALID
		}
	}

	return t.Claims.(jwt.Claims), -1
}

func JwtGuardHandler(c *fiber.Ctx) error {
	log.Info("Connected")
	// Check if the request is authorized
	// If not, return an error
	token := tools.HeaderParser(c)
	if token == "" {
		log.Info("Unauthorized, no authorization protocol used")
		return c.Status(400).JSON(tools.GlobalErrorHandlerResp{
			Success: false,
			Message: "Unauthorized, no authorization protocol used",
			Code:    400,
		})
	}
	claims, err := Decode(token, false)
	if err != -1 {
		log.Info("Failed to decode token ", err)
		if err == tools.JWT_FAILED_TO_DECODE {
			log.Info(token)
			return c.Status(400).JSON(tools.GlobalErrorHandlerResp{
				Success: false,
				Message: "Failed to decode token",
				Code:    400,
			})
		} else if err == tools.JWT_INVALID {
			return c.Status(401).JSON(tools.GlobalErrorHandlerResp{
				Success: false,
				Message: "Invalid token",
				Code:    401,
			})
		} else if err == tools.JWT_EXPIRED {
			return c.Status(408).JSON(tools.GlobalErrorHandlerResp{
				Success: false,
				Message: "Token expired",
				Code:    408,
			})
		} else {
			return c.Status(401).JSON(tools.GlobalErrorHandlerResp{
				Success: false,
				Message: "Unknown error",
				Code:    401,
			})
		}
	}

	log.Info(claims)

	c.Locals("claims", claims)
	c.Locals("isAuthorized", true)

	return c.Next()
}

func JwtGuardMasterHandler(c *fiber.Ctx) error {
	log.Info("Connected")
	// Check if the request is authorized
	// If not, return an error
	token := tools.HeaderParser(c)
	if token == "" {
		log.Info("Unauthorized, no authorization protocol used")
		return c.Status(400).JSON(tools.GlobalErrorHandlerResp{
			Success: false,
			Message: "Unauthorized, no authorization protocol used",
			Code:    400,
		})
	}
	claims, err := Decode(token, false)
	if err != -1 {
		log.Info("Failed to decode token")
		if err == tools.JWT_FAILED_TO_DECODE {
			return c.Status(400).JSON(tools.GlobalErrorHandlerResp{
				Success: false,
				Message: "Failed to decode token",
				Code:    400,
			})
		} else if err == tools.JWT_INVALID {
			return c.Status(401).JSON(tools.GlobalErrorHandlerResp{
				Success: false,
				Message: "Invalid token",
				Code:    401,
			})
		} else if err == tools.JWT_EXPIRED {
			return c.Status(408).JSON(tools.GlobalErrorHandlerResp{
				Success: false,
				Message: "Token expired",
				Code:    408,
			})
		} else {
			return c.Status(401).JSON(tools.GlobalErrorHandlerResp{
				Success: false,
				Message: "Unknown error",
				Code:    401,
			})
		}
	}

	userId, jwtErr := claims.GetSubject()
	if jwtErr != nil {
		return c.Status(401).JSON(tools.GlobalErrorHandlerResp{
			Success: false,
			Message: "Failed to get user id",
			Code:    401,
		})
	}

	if !strings.Contains(tools.MASTER_IDS, userId) {
		return c.Status(401).JSON(tools.GlobalErrorHandlerResp{
			Success: false,
			Message: "Unauthorized, not a master",
			Code:    401,
		})
	}

	c.Locals("claims", claims)
	c.Locals("isAuthorized", true)

	return c.Next()
}
