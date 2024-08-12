package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"gambler/backend/database/models"
	"gambler/backend/tools"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/redis/go-redis/v9"
)

type (
	CacheHandler struct {
		redis *redis.Client
		ctx   context.Context
	}
	CurrentGame struct {
		GameId string
		Users  []string
	}
)

var Cache CacheHandler

func NewCache() *CacheHandler {
	Cache = CacheHandler{
		redis: redis.NewClient(&redis.Options{
			Addr:     tools.HOST_REDIS,
			Password: tools.PSW_REDIS,
			DB:       0,
		}),
		ctx: context.Background(),
	}
	log.Info("[CACHE] Connected to Redis")
	return &Cache
}

// StoreUserConnection stores a user's WebSocket connection ID in Redis
func (c *CacheHandler) StoreUserConnection(userID, connectionID string) int {
	cacheKey := fmt.Sprintf("user:%s:connection", userID)
	err := c.redis.Set(c.ctx, cacheKey, connectionID, time.Hour*6).Err()
	if err != nil {
		return HandleRedisError(err)
	}
	return -1
}

// GetUserConnection retrieves a user's WebSocket connection ID from Redis
func (c *CacheHandler) GetUserConnection(userID string) (string, int) {
	cacheKey := fmt.Sprintf("user:%s:connection", userID)
	connectionID, err := c.redis.Get(c.ctx, cacheKey).Result()
	if err == redis.Nil {
		return "", tools.RD_KEY_NOT_FOUND
	} else if err != nil {
		return "", HandleRedisError(err)
	}
	return connectionID, -1
}

// RemoveUserConnection removes a user's WebSocket connection ID from Redis
func (c *CacheHandler) RemoveUserConnection(userID string) int {
	cacheKey := fmt.Sprintf("user:%s:connection", userID)
	err := c.redis.Del(c.ctx, cacheKey).Err()
	if err != nil {
		return HandleRedisError(err)
	}
	return -1
}

func (c *CacheHandler) SetBet(bet models.Bet) int {
	res := c.redis.Set(c.ctx, "b-"+fmt.Sprintf("%d", bet.ID), bet, time.Hour*6)
	if res.Err() != nil {
		return HandleRedisError(res.Err())
	}
	return -1
}

func (c *CacheHandler) GetBetById(key string) (*models.Bet, int) {
	var bet models.Bet
	betDataString, err := c.redis.Get(c.ctx, key).Result()
	if err == redis.Nil {
		return nil, -1
	} else if err != nil {
		return nil, HandleRedisError(err)
	}
	err = json.Unmarshal(
		[]byte(betDataString),
		&bet,
	)
	if err != nil {
		return nil, tools.JSON_UNMARSHAL_ERROR
	}
	return &bet, -1
}

func (c *CacheHandler) GetAllBet() (*[]models.Bet, int) {
	keys, err := c.redis.Keys(c.ctx, "b-*").Result()
	if err != nil {
		return nil, HandleRedisError(err)
	}
	var bets []models.Bet
	for _, key := range keys {
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
	if e.Error() == redis.ErrClosed.Error() {
		return tools.RD_CONN_CLOSED
	} else if e.Error() == redis.Nil.Error() {
		return tools.RD_KEY_NOT_FOUND
	} else if e.Error() == redis.TxFailedErr.Error() {
		return tools.RD_TX_FAILED
	} else {
		return tools.RD_UNKNOWN
	}
}
