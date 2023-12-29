package handlers

import (
	"fmt"
	w "gortc/pkg/webrtc"
	"os"

	"github.com/gofiber/fiber"
	"github.com/gofiber/websocket"
)

func Stream(c *fiber.Ctx) error {
	suuid := c.Params("suuid")
	if suuid == "" {
		c.Status(400)
		return nil
	}
	ws := "ws"
	if os.Getenv("ENV") == "PROD" {
		ws = "wss"
	}
	w.RoomsLock.Lock()
	if _, ok := w.Streams[suuid]; ok {
		w.RoomsLock.Unlock()
		return c.Render("stream", fiber.Map{
			"StreamWebsocketAddr": fmt.Sprintf("%s://%s/stream/%s/websocket", ws, c.Hostname(), suuid),
			"ChatWebsocketAddr":   fmt.Sprintf("%s://%s/stream/%s/chat/websocket", ws, c.Hostname(), suuid),
			"ViewerWebsocketAddr": fmt.Sprintf("%s://%s/stream/%s/viewer/websocket", ws, c.Hostname(), suuid),
			"Type":                "stream",
		}, "layouts/main")
	}
	w.RoomsLock.Unlock()
	return c.Render("stream", fiber.Map{
		"NoStream": "true",
		"Leave":    "true",
	}, "layouts/main")
}

func StreamWebsocket(c *websocket.Conn) {

}

func StreamViewerWebsocket(c *websocket.Conn) {

}

func viewerConn(c *websocket.Conn, p *w.Peers) {

}
