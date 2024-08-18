package handlers

import (
	"encoding/json"
	"fmt"
	"gambler/backend/tools"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2/log"
)

type WebSocketHandler struct {
	Cache *CacheHandler
}

var (
	activeConnections = make(map[string]*websocket.Conn)
	WebSocket         WebSocketHandler
)

// NewWebSocketHandler initializes a new WebSocketHandler
func NewWebSocketHandler(cache *CacheHandler) *WebSocketHandler {
	WebSocket = WebSocketHandler{
		Cache: cache,
	}
	return &WebSocket
}

// ErrorMessage defines the format for error messages sent to clients
type ErrorMessage struct {
	Type    string `json:"type"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// sendErrorMessage sends an error message to the WebSocket client
func (wsh *WebSocketHandler) sendErrorMessage(c *websocket.Conn, code int, errMessage string) {
	errorMsg := ErrorMessage{
		Type:    "error",
		Code:    code,
		Message: errMessage,
	}
	msg, _ := json.Marshal(errorMsg)
	c.WriteMessage(websocket.TextMessage, msg)
}

// HandleWebSocketConnection manages the WebSocket connection for a specific user
func (wsh *WebSocketHandler) HandleWebSocketConnection(c *websocket.Conn) {
	// Get unique connection ID (UUID) from WebSocket connection (or use other unique ID)
	uuid := c.Params("id")

	// Store the connection in Redis
	if errCode := wsh.Cache.StoreUserConnection(uuid, uuid); errCode != -1 {
		wsh.sendErrorMessage(c, errCode, "Failed to store WebSocket connection in Redis")
		c.Close()
		return
	}

	// Store the connection in the activeConnections map for in-memory access
	activeConnections[uuid] = c

	// Ensure the connection is removed from Redis and the map when the user disconnects
	defer func() {
		if errCode := wsh.Cache.RemoveUserConnection(uuid); errCode != -1 {
			log.Error("Failed to remove WebSocket connection from Redis:", errCode)
		}
		delete(activeConnections, uuid)
		c.Close()
	}()

	log.Info(fmt.Sprintf("User %s connected to WebSocket", uuid))

	wsh.SendMessageToUser(uuid, []byte{0, tools.WEBSOCKET_VERSION, 0})

	// Main loop to handle incoming WebSocket messages
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			wsh.sendErrorMessage(c, tools.WS_INVALID_CONN, "Error reading WebSocket message")
			break
		}
		log.Info(fmt.Sprintf("Received message from user %s: %x", uuid, msg))
	}
}

// SendMessageToUser sends a message to a specific user based on their UUID
func (wsh *WebSocketHandler) SendMessageToUser(uuid string, message []byte) error {
	// Get the WebSocket connection from the activeConnections map
	conn, exists := activeConnections[uuid]
	if !exists {
		return fmt.Errorf("connection not found for UUID %s", uuid)
	}

	// Send the message over the WebSocket connection
	if err := conn.WriteMessage(websocket.BinaryMessage, message); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}
