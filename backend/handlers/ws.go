package handlers

import (
	"fmt"
	"net/http"
	"time"

	"GoBaatcheet/config"
	"GoBaatcheet/models"

	"github.com/goombaio/namegenerator"
	"github.com/gorilla/websocket"
)

var users map[string]*websocket.Conn = make(map[string]*websocket.Conn)

func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	config.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := config.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("E#1PZU6V - Failed to upgrade request to websocket connection!")
	}
	users[readOrAssignUsername(ws)] = ws
	if err != nil {
		fmt.Println("E#1PZUJN - Error while writing message back to client!")
	}
	reader(ws)
}

func readOrAssignUsername(conn *websocket.Conn) string {
	var user models.User
	err := conn.ReadJSON(&user)
	if err != nil {
		fmt.Println("E#1QENH1 - Unable to read username from websocket connection.", err)
		return getRandomUsername()
	}
	fmt.Println("TMP log:", user)
	return user.Username
}

func getRandomUsername() string {
	seed := time.Now().UTC().UnixNano()
	nameGenerator := namegenerator.NewNameGenerator(seed)
	return nameGenerator.Generate()
}

func reader(conn *websocket.Conn) {
	for {
		var msg models.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println("E#1PZUA7 - Error while reading message for user:", msg.Sender)
			return
		}
		msgToSend := models.Message{
			Sender:   msg.Sender,
			Receiver: msg.Receiver,
			Msg:      msg.Msg,
		}
		users[msg.Sender].WriteJSON(msgToSend)
	}
}
