package handlers

import (
	"fmt"
	"log"

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

	_, _, room := createOrGetRoom(uuid)
	log.Println(room)
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
	return "", "", RoomTyp{}
}

func RoomViewerWebSocket(c *websocket.Conn) {

}

func roomViewerConn(c *websocket.Conn, p *w.Peers) {

}

type websocketMessage struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}
