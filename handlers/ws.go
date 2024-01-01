package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/itishrishikesh/GoBaatcheet/config"
)

func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	config.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := config.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("E#1PZU6V - Failed to upgrade request to websocket connection!")
	}
	fmt.Println("I#1PZUGS - Websocket Connection is Established with Client!")

	err = ws.WriteMessage(1, []byte("Hello from GoBaatCheet!"))

	if err != nil {
		fmt.Println("E#1PZUJN - Error while writing message back to client!")
	}

	reader(ws)
}

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("E#1PZUA7 - Error while reading message!")
			return
		}
		fmt.Println(string(p))
		if err := conn.WriteMessage(messageType, p); err != nil {
			fmt.Println("E#1PZUDM - Error while writing message!")
			return
		}
	}
}
