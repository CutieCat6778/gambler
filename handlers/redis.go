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
)

var Cache CacheHandler

func NewCache(app *fiber.App) *CacheHandler {
	if tools.HOST_REDIS == "" || tools.PSW_REDIS == "" {
		log.Fatal("[CACHE] Redis connection details not found in .env")
		tools.SendWebHook("[CACHE] Redis connection details not found in .env")
		panic("Redis connection details not found in .env")
	}
	Cache = CacheHandler{
		Redis: redis.New(redis.Config{
			Host:     tools.HOST_REDIS,
			Port:     6379,
			Password: tools.PSW_REDIS,
			Database: 0,
		}),
		Context: context.Background(),
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

// ListenForExpiredKeys listens for expired keys in Redis and handles them
func (c *CacheHandler) ListenForExpiredKeys() {
	// Subscribe to the Redis expired events
	pubsub := c.Redis.Conn().Subscribe(c.Context, "__keyevent@0__:expired") // Replace 0 with your Redis DB index if different

	// Handle messages in a separate goroutine
	go func() {
		for {
			msg, err := pubsub.ReceiveMessage(c.Context)
			if err != nil {
				log.Error("Failed to receive message from Redis:", err)
				continue
			}

			log.Info("Received expired key event:", msg.Payload)

			// Handle the expired key event (msg.Payload contains the expired key name)
			c.HandleExpiredKey(msg.Payload)
		}
	}()
}

// HandleExpiredKey processes the expired key event
func (c *CacheHandler) HandleExpiredKey(key string) {
	// Add your logic to handle expired keys here
	if strings.HasPrefix(key, "b-") {
		log.Info("Bet expired:", key)
		// You can add additional logic to handle the expiration of a bet, e.g., update the database, notify users, etc.
	}
}

// Bets
func (c *CacheHandler) SetBet(bet models.Bet) int {
	betData, err := json.Marshal(bet)
	if err != nil {
		return HandleRedisError(err)
	}
	// Save the JSON string to Redis with a key prefix
	res := c.Redis.Conn().Set(c.Context, "b-"+fmt.Sprintf("%d", bet.ID), betData, time.Until(bet.EndsAt)).Err()
	if res != nil {
		log.Error(res)
		return HandleRedisError(res)
	}
	return -1
}

func (c *CacheHandler) RemoveBet(betID uint) int {
	// Remove the bet from Redis
	res := c.Redis.Conn().Del(c.Context, "b-"+fmt.Sprintf("%d", betID)).Err()
	if res != nil {
		return HandleRedisError(res)
	}
	return -1
}

func (c *CacheHandler) GetBetById(betID uint) (*models.Bet, int) {
	var bet models.Bet

	key := fmt.Sprintf("b-%d", betID)

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
		bet, err := c.GetBetById(tools.ConvertKeyToBetID(key))
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

func (c *CacheHandler) GetAllBetByAmount(amount int) (*[]models.Bet, int) {
	bets, err := c.GetAllBet()
	if err != -1 {
		return nil, err
	}
	filteredBets := []models.Bet{}
	for _, bet := range *bets {
		if len(filteredBets) <= amount {
			filteredBets = append(filteredBets, bet)
		}
	}
	return &filteredBets, -1
}

func (c *CacheHandler) UpdateBet(betID uint) int {
	bet, err := DB.GetBetByID(betID)
	if err != -1 {
		return err
	}

	err = c.SetBet(*bet)
	if err != -1 {
		return err
	}
	return -1
}

func (c *CacheHandler) LoadDatabaseBets() int {
	bets, err := DB.GetAllActiveBets()
	if err != -1 {
		return err
	}
	for _, bet := range *bets {
		log.Info("Loaded bet", bet.ID)
		log.Info(fmt.Sprintf("%v", bet))
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
