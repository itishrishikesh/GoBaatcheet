package handlers

import (
	"crypto/sha256"
	"fmt"
	"gortc/pkg/chat"
	w "gortc/pkg/webrtc"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket"
	guuid "github.com/google/uuid"
)

func RoomCreate(context *fiber.Ctx) error {
	return context.Redirect(fmt.Sprintf("/room/%s", guuid.New().String()))
}

func Room(context *fiber.Ctx) error {
	uuid := context.Params("uuid")
	if uuid == "" {
		context.Status(400)
		return nil
	}
	ws := "ws"
	if os.Getenv("ENV") == "PROD" {
		ws = "wss"
	}
	uuid, suuid, _ := createOrGetRoom(uuid)
	return context.Render("peer", fiber.Map{
		"RoomWebSocketAddr":   fmt.Sprintf("%s://%s/room%s/websocket", ws, context.Hostname(), uuid),
		"RoomLink":            fmt.Sprintf("%s://%s/room%s", context.Protocol(), context.Hostname(), uuid),
		"ChatWebsocketAddr":   fmt.Sprintf("%s://%s/room%s/chat/websocket", ws, context.Hostname(), uuid),
		"ViewerWebsocketAddr": fmt.Sprintf("%s://%s/room%s/viewer/websocket", ws, context.Hostname(), uuid),
		"StreamLink":          fmt.Sprintf("%s://%s/stream%s", context.Protocol(), context.Hostname(), suuid),
		"Type":                "room",
	}, "layouts/main")
}

func RoomWebsocket(connection *websocket.Conn) {
	uuid := connection.Params("uuid")

	if uuid == "" {
		return
	}
	_, _, room := createOrGetRoom(uuid)
	w.RoomConn(c, room.Peers)
	log.Println(room)
}

type RoomTyp struct {
}

func createOrGetRoom(uuid string) (string, string, RoomTyp) {
	w.RoomsLock.Lock()
	defer w.RoomsLock.Unlock()
	h := sha256.New()
	h.Write([]byte(uuid))
	suuid := fmt.Sprintf("%x", h.Sum(nil))

	if room := w.Rooms[uuid]; room != nil {
		if _, ok := w.Streams[suuid]; !ok {
			w.Streams[suuid] = room
		}
		return uuid, suuid, room
	}
	hub := chat.NewHub()
	p := &w.Peers{}
	p.TrackLocals = make(map[string]*webrtc.TrackLocalStaticRTP)
	room := &w.Room{
		Peers: p,
		Hub:   hub,
	}
	w.Rooms[uuid] = room
	w.Streams[suuid] = room
	go hub.Run()
	return uuid, suuid, room
}

func RoomViewerWebSocket(c *websocket.Conn) {
	uuid := c.Params("uuid")
	if uuid == "" {
		return
	}
	w.RoomLock.Lock()
	if peer, ok := w.Rooms[uuid]; ok {
		w.RoomsLock.Unlock()
		roomViewerConn(c, peer.Peers)
	}
	w.RoomLock.Unlock()
}

func roomViewerConn(c *websocket.Conn, p *w.Peers) {
	ticket := time.NewTicker(1 * time.Second)
	defer ticket.Stop()
	defer c.Close()

	for {
		select {
		case <-ticket.C:
			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write([]byte(fmt.Sprintf("%d", len(p.Connections))))
		}
	}
}

type websocketMessage struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}
