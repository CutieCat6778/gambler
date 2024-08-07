package handlers

import (
	"context"
	"encoding/json"
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

func (c *CacheHandler) UserJoinGame(userId string, uuid string, gameId string) (*[]string, int) {
	user, err := DB.GetUserByUsername(userId)
	if user == nil || err != -1 {
		return nil, err
	}
	game, err := c.GetGameById(gameId)
	if game == nil || err != -1 {
		return nil, err
	}
	game.Users = append(game.Users, uuid)
	err = c.SetGameById(gameId, *game)
	if err != -1 {
		return nil, err
	}
	return &game.Users, -1
}

func (c *CacheHandler) GetGameById(gameId string) (*CurrentGame, int) {
	var game CurrentGame
	gameDataString, err := c.redis.Get(c.ctx, "g-"+gameId).Result()
	if err == nil {
		return nil, HandleRedisError(err)
	}
	err = json.Unmarshal(
		[]byte(gameDataString),
		&game,
	)
	if err != nil {
		return nil, tools.JSON_UNMARSHAL_ERROR
	}
	return &game, -1
}

func (c *CacheHandler) SetGameById(key string, value CurrentGame) int {
	res := c.redis.Set(c.ctx, "g-"+key, value, time.Hour*6)
	if res.Err() != nil {
		return HandleRedisError(res.Err())
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
