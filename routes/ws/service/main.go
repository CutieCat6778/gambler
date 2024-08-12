package service

import (
	"strings"

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
