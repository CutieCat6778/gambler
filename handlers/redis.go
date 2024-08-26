package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"gambler/backend/database/models"
	"gambler/backend/tools"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/storage/redis/v3"
	r "github.com/redis/go-redis/v9"
)

type (
	CacheHandler struct {
		Redis   *redis.Storage
		Context context.Context
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
		Context: context.Background(),
	}
	Cache.LoadDatabaseBets()
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
	// Save the JSON string to Redis with a key prefix
	res := c.Redis.Conn().Set(c.Context, "b-"+fmt.Sprintf("%d", bet.ID), bet, time.Hour*6).Err()
	if res != nil {
		log.Error(res)
		return HandleRedisError(res)
	}
	return -1
}

func (c *CacheHandler) GetBetById(key string) (*models.Bet, int) {
	var bet models.Bet

	// Get the JSON string from Redis
	res := c.Redis.Conn().Get(c.Context, key)
	err := res.Err()
	if err != nil {
		log.Error(err, key)
		if err == r.Nil {
			return nil, tools.RD_KEY_NOT_FOUND
		}
		return nil, HandleRedisError(err)
	}

	// Unmarshal the JSON string into the models.Bet struct
	data, err := res.Bytes()
	if err != nil {
		log.Error("Failed to get bytes from Redis:", err)
		return nil, HandleRedisError(err)
	}

	err = json.Unmarshal(data, &bet)
	if err != nil {
		log.Error("Failed to unmarshal JSON into Bet struct:", err)
		return nil, HandleRedisError(err)
	}

	return &bet, -1
}

func (c *CacheHandler) GetAllBet() (*[]models.Bet, int) {
	// Retrieve all keys from Redis
	keys, err := c.Redis.Keys()
	if err != nil {
		return nil, HandleRedisError(err)
	}
	log.Info(keys)

	bets := []models.Bet{}
	for _, key := range keys {
		key := string(key)
		// Skip keys that don't have the "b-" prefix
		if !strings.HasPrefix(key, "b-") {
			continue
		}
		log.Info(key)
		// Retrieve the bet by ID
		bet, err := c.GetBetById(key)
		if err != -1 {
			if err == tools.RD_KEY_NOT_FOUND {
				continue
			}
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
		log.Info("Loaded bet", bet.ID)
		err := c.SetBet(bet)
		if err != -1 {
			log.Error(err)
			return err
		}
	}
	return -1
}

func HandleRedisError(e error) int {
	if e == r.Nil {
		return tools.RD_KEY_NOT_FOUND
	} else {
		return tools.RD_UNKNOWN
	}
}
