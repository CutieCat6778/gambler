package service

import (
	"encoding/json"
	"gambler/backend/handlers"
	"gambler/backend/tools"
	"strconv"
	"strings"

	"github.com/gofiber/contrib/socketio"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type Message struct {
	Data  string `json:"data"`
	From  string `json:"from"`
	Event string `json:"event"`
	To    string `json:"to"`
}

func CompatibleCheck(c *fiber.Ctx) error {
	if !strings.HasPrefix(c.Route().Path, "/ws") {
		return c.Next()
	}
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		log.Info("Allowed Connection!")
		return c.Next()
	}
	log.Error("Connection refused!")
	return fiber.ErrUpgradeRequired
}

func UserConnect(kws *socketio.Websocket) {
	userId := kws.Params("id")
	gameId := kws.Params("game_id")

	users, err := handlers.Cache.UserJoinGame(userId, kws.UUID, gameId)
	if err != -1 {
		kws.Fire(socketio.EventError, ParseErr(err))
		return
	}

	kws.SetAttribute("user_id", userId)
	kws.SetAttribute("game_id", gameId)

	kws.EmitToList(*users, []byte("[Server] user connected: "+userId), socketio.TextMessage)
}

func OnHandshake(e *socketio.EventPayload) {
	log.Info("Connection event 1 - User: %s", e.Kws.GetStringAttribute("user_id"))
}

func OnMessage(e *socketio.EventPayload) {
	log.Info("Message event - User: %s - Message: %s", e.Kws.GetStringAttribute("user_id"), string(e.Data))

	message := Message{}

	err := json.Unmarshal(e.Data, &message)
	if err != nil {
		log.Error("Error parsing message: %s", err)
		return
	}
}

func OnError(e *socketio.EventPayload) error {
	kws := e.Kws
	message := Message{}

	err := json.Unmarshal(e.Data, &message)
	if err != nil {
		log.Error("Error parsing message: %s", err)
		return err
	}

	errCode, err := strconv.ParseInt(message.Data, 10, 32)
	if err != nil || errCode == tools.WS_UNKNOWN_ERR {
		return kws.EmitTo(message.To, []byte("Unkown Error"), socketio.TextMessage)
	}

	return kws.EmitTo(message.To, ParseErr(int(errCode)), socketio.TextMessage)
}

func ParseErr(e int) []byte {
	switch e {
	case tools.WS_UUID_DUP:
		return []byte("UUID Duplicated")
	case tools.WS_UUID_NOTFOUND:
		return []byte("User not found")
	case tools.WS_GAMEID_NOTFOUND:
		return []byte("Game not found")
	case tools.WS_INVALID_CONN:
		return []byte("Invalid Connection")
	default:
		return []byte("Unkown Error")
	}
}
