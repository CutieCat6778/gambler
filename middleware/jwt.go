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

func Sign(userId uint) (*Jwt, int) {
	user, dbErr := handlers.DB.GetUserByID(userId)
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
		Subject:   fmt.Sprintf("%d", user.ID),
	})
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenExpDate)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer:    fmt.Sprintf("%d Version", user.RefreshTokenVersion),
		Subject:   fmt.Sprintf("%d", user.ID),
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
		userId, jwtErr := t.Claims.GetSubject()
		if jwtErr != nil {
			return nil, tools.JWT_FAILED_TO_DECODE
		}

		user, dbErr := handlers.DB.GetUserByID(tools.ParseUInt(userId))
		if dbErr != -1 {
			return nil, dbErr
		}

		if user.RefreshTokenVersion != int(ver) {
			return nil, tools.JWT_INVALID
		}
	}

	return t.Claims, -1
}

func JwtGuardHandler(c *fiber.Ctx) error {
	log.Info("Connected")
	// Check if the request is authorized
	// If not, return an error
	token := c.Cookies("access_token")
	log.Info("Token", token)
	if token == "" {
		refresh_token := c.Cookies("refresh_token")
		if refresh_token == "" {
			return tools.ReturnData(c, 401, nil, tools.JWT_NO_KEY)
		}
		return c.Redirect("/auth/refresh", 307)
	}
	claims, err := Decode(token, false)
	if err != -1 {
		log.Info("Failed to decode token ", err)
		return tools.ReturnData(c, 401, nil, err)
	}

	log.Info(claims)

	c.Locals("claims", claims)
	c.Locals("isAuthorized", true)

	return c.Next()
}
