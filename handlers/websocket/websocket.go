package websocket

import (
	"encoding/json"
	"fmt"
	"gambler/backend/handlers"
	"gambler/backend/tools"
	"runtime"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2/log"
)

type WebSocketHandler struct {
	Cache             *handlers.CacheHandler
	ActiveConnections map[string]*websocket.Conn
}

var (
	WebSocket WebSocketHandler
)

// NewWebSocketHandler initializes a new WebSocketHandler
func NewWebSocketHandler(cache *handlers.CacheHandler) *WebSocketHandler {
	WebSocket = WebSocketHandler{
		Cache:             cache,
		ActiveConnections: make(map[string]*websocket.Conn),
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
func (wsh *WebSocketHandler) SendErrorMessage(uuid string, code int, errMessage string) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		log.Info(fmt.Sprintf("Called from %s, line %d", file, line))
	}
	log.Info("Called from", file, line)
	log.Info("Sending error message to user:", uuid, code, errMessage)
	errorMsg := ErrorMessage{
		Type:    "error",
		Code:    code,
		Message: errMessage,
	}
	msg, _ := json.Marshal(errorMsg)
	msgAsByte := []byte(msg)
	headers := []byte{tools.WS_ERR, tools.WEBSOCKET_VERSION}
	headers = append(headers, msgAsByte...)
	err := wsh.SendMessageToUser(uuid, headers)
	if err != -1 {
		log.Error(tools.GetErrorString(err))
	}
}

// HandleWebSocketConnection manages the WebSocket connection for a specific user
func (wsh *WebSocketHandler) HandleWebSocketConnection(c *websocket.Conn) {
	// Get unique connection ID (UUID) from WebSocket connection (or use other unique ID)
	uuid := c.Params("id")

	// Store the connection in the activeConnections map for in-memory access
	wsh.ActiveConnections[uuid] = c

	// Ensure the connection is removed from Redis and the map when the user disconnects
	defer func() {
		delete(wsh.ActiveConnections, uuid)
		c.Close()
	}()

	log.Info(fmt.Sprintf("User %s connected to WebSocket", uuid))

	wsh.SendMessageToUser(uuid, []byte{0, tools.WEBSOCKET_VERSION, 0})

	// Main loop to handle incoming WebSocket messages
	// go func() {
	for {
		msgType, msg, err := c.ReadMessage()
		if err != nil {
			log.Info("Message type ", msgType, msg)
			log.Info(err.Error())
			wsh.SendErrorMessage(uuid, tools.WS_INVALID_CONN, "Error reading WebSocket message")
			break
		}
		log.Info(fmt.Sprintf("Received message from user %s: %v", uuid, msg))
		HandleMessageEvent(wsh, uuid, int(msg[0]), msg[2:])
	}
	// }()
}

// SendMessageToUser sends a message to a specific user based on their UUID
func (wsh *WebSocketHandler) SendMessageToUser(uuid string, message []byte) int {
	// Get the WebSocket connection from the activeConnections map
	conn, exists := wsh.ActiveConnections[uuid]
	if !exists {
		return tools.WS_UUID_NOTFOUND
	}

	// Send the message over the WebSocket connection
	if err := conn.WriteMessage(websocket.BinaryMessage, message); err != nil {
		return tools.WS_UNKNOWN_ERR
	}
	return -1
}

func (wsh *WebSocketHandler) SendMessageToAll(message []byte) int {
	for _, conn := range wsh.ActiveConnections {
		if err := conn.WriteMessage(websocket.BinaryMessage, message); err != nil {
			log.Info(err)
			continue
		}
		log.Info(fmt.Sprintf("Sent message to user %v", conn.Params("id")))
	}
	return -1
}

func (wsh *WebSocketHandler) UpdateBet(betID uint) int {
	result := []byte{tools.BET_UPDATE, tools.WEBSOCKET_VERSION}
	betIdChunks := tools.ChunkBigNumber(int(betID))
	result = append(result, betIdChunks...)
	err := wsh.SendMessageToAll(result)
	if err != -1 {
		return err
	}
	return -1
}

func (wsh *WebSocketHandler) UpdateUser(uuid string) int {
	err := wsh.SendMessageToUser(uuid, []byte{tools.USER_UPDATE, tools.WEBSOCKET_VERSION})
	if err != -1 {
		return err
	}
	return -1
}
