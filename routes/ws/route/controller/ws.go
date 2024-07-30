package controller

import "github.com/gofiber/contrib/websocket"

func HandleNewConnection(c *websocket.Conn) {
	if c.Locals("allowed") == nil || c.Locals("allowed") == false {
		c.Close()
		return
	}
}
