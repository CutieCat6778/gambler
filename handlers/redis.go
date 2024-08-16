package handlers

import (
	"fmt"
	"gambler/backend/database/models"
	"gambler/backend/tools"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/storage/redis/v3"
)

type (
	CacheHandler struct {
		Redis *redis.Storage
	}
	CurrentGame struct {
		GameId string
		Users  []string
	}
)

var Cache CacheHandler

func NewCache(app *fiber.App) *CacheHandler {
	log.Info(tools.HOST_REDIS, tools.PSW_REDIS, "ABC")
	if tools.HOST_REDIS == "" || tools.PSW_REDIS == "" {
		log.Fatal("[CACHE] Redis connection details not found in .env")
		panic("Redis connection details not found in .env")
	}
	Cache = CacheHandler{
		Redis: redis.New(redis.Config{
			Host:     tools.HOST_REDIS,
			Port:     18254,
			Password: tools.PSW_REDIS,
			Database: 0,
		}),
	}
	log.Info("[CACHE] Connected to Redis")
	return &Cache
}

func AddCache(exp time.Duration) fiber.Handler {
	return cache.New(cache.Config{
		Expiration:   exp,
		CacheControl: true,
		Storage:      Cache.Redis,
	})
}

// StoreUserConnection stores a user's WebSocket connection ID in Redis
func (c *CacheHandler) StoreUserConnection(userID, connectionID string) int {
	cacheKey := fmt.Sprintf("user:%s:connection", userID)
	err := c.Redis.Set(cacheKey, []byte(connectionID), time.Hour*6)
	if err != nil {
		return HandleRedisError(err)
	}
	return -1
}

// GetUserConnection retrieves a user's WebSocket connection ID from Redis
func (c *CacheHandler) GetUserConnection(userID string) (string, int) {
	cacheKey := fmt.Sprintf("user:%s:connection", userID)
	connectionID, err := c.Redis.Get(cacheKey)
	if err != nil {
		return "", tools.RD_KEY_NOT_FOUND
	} else if err != nil {
		return "", HandleRedisError(err)
	}
	return string(connectionID), -1
}

// RemoveUserConnection removes a user's WebSocket connection ID from Redis
func (c *CacheHandler) RemoveUserConnection(userID string) int {
	cacheKey := fmt.Sprintf("user:%s:connection", userID)
	err := c.Redis.Delete(cacheKey)
	if err != nil {
		return HandleRedisError(err)
	}
	return -1
}

func (c *CacheHandler) SetBet(bet models.Bet) int {
	res := c.Redis.Set("b-"+fmt.Sprintf("%d", bet.ID), []byte(bet.Name), time.Hour*6)
	if res != nil {
		return HandleRedisError(res)
	}
	return -1
}

func (c *CacheHandler) GetBetById(key string) (*models.Bet, int) {
	betId, err := c.Redis.Get(key)
	if err != nil {
		return nil, HandleRedisError(err)
	}
	bet, dbErr := DB.GetBetByBetName(string(betId))
	if dbErr != -1 {
		return nil, dbErr
	}
	return bet, -1
}

func (c *CacheHandler) GetAllBet() (*[]models.Bet, int) {
	keys, err := c.Redis.Keys()
	if err != nil {
		return nil, HandleRedisError(err)
	}
	var bets []models.Bet
	for _, key := range keys {
		key := string(key)
		if strings.HasPrefix(key, "b-") {
			continue
		}
		bet, err := c.GetBetById(key)
		if err != -1 {
			return nil, err
		}
		bets = append(bets, *bet)
	}
	return &bets, -1
}

func (c *CacheHandler) LoadDatabaseBets() int {
	bets, err := DB.GetAllActiveBets()
	if err != -1 {
		return err
	}
	for _, bet := range *bets {
		c.SetBet(bet)
	}
	return -1
}

func HandleRedisError(e error) int {
	return tools.RD_UNKNOWN
}
