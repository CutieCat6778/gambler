package routine

import (
	"fmt"
	"gambler/backend/database/models/customTypes"
	"gambler/backend/handlers"
	"gambler/backend/handlers/websocket"
	"gambler/backend/tools"
	"strings"

	"github.com/gofiber/fiber/v2/log"
)

// ListenForExpiredKeys listens for expired keys in Redis and handles them
func ListenForExpiredKeys() {
	// Subscribe to the Redis expired events
	pubsub := handlers.Cache.Redis.Conn().Subscribe(handlers.Cache.Context, "__keyevent@0__:expired") // Replace 0 with your Redis DB index if different

	// Handle messages in a separate goroutine
	go func() {
		for {
			msg, err := pubsub.ReceiveMessage(handlers.Cache.Context)
			if err != nil {
				log.Error("Failed to receive message from Redis:", err)
				continue
			}

			log.Info("Received expired key event:", msg.Payload)

			// Handle the expired key event (msg.Payload contains the expired key name)
			HandleExpiredKey(msg.Payload)
		}
	}()
}

// HandleExpiredKey processes the expired key event
func HandleExpiredKey(key string) {
	// Add your logic to handle expired keys here
	if strings.HasPrefix(key, "b-") {
		log.Info("Bet expired:", key)
		// You can add additional logic to handle the expiration of a bet, e.g., update the database, notify users, etc.
		betID := tools.ConvertKeyToBetID(key)
		bet, err := handlers.DB.UpdateBetStatus(betID, customTypes.Pending)
		if err != -1 {
			log.Error("Failed to update bet status:", err)
			tools.SendWebHook(fmt.Sprintf("Failed to update bet status: %d", betID))
		}
		log.Info("Updated bet status to Pending:", bet.ID)
		err = handlers.Cache.UpdateBet(bet.ID)
		if err != -1 {
			log.Error("Failed to update bet in cache:", err)
			tools.SendWebHook(fmt.Sprintf("Failed to update bet in cache: %d", betID))
		}
		log.Info("Updated bet in cache:", bet.ID)
		websocket.WebSocket.UpdateBet(betID)
	}
}
