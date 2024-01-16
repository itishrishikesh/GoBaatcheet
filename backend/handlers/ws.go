package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"GoBaatcheet/auth"
	"GoBaatcheet/config"
	"GoBaatcheet/models"

	"github.com/goombaio/namegenerator"
	"github.com/gorilla/websocket"
)

var connectedUsers map[string]*websocket.Conn = make(map[string]*websocket.Conn)

func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	if !auth.Authenticate(r) {
		w.WriteHeader(401)
	}
	config.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := config.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalln("E#1PZU6V - Failed to upgrade request to websocket connection!", err)
	}
	username, err := readOrAssignUsername(ws)
	if err != nil {
		log.Fatalln("")
	}
	connectedUsers[username] = ws
	if err != nil {
		log.Println("E#1PZUJN - Error while writing message back to client!")
	}
	reader(ws)
}

func readOrAssignUsername(conn *websocket.Conn) (string, error) {
	var user models.User
	err := conn.ReadJSON(&user)
	if err != nil {
		return "", fmt.Errorf("E#1QENH1 - Unable to read username from websocket connection. %v", err)
	}
	return user.Username, nil
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
			log.Println("E#1PZUA7 - Error while reading message for user:", msg.Sender, err)
			return
		}
		msgToSend := models.Message{
			Sender:   msg.Sender,
			Receiver: msg.Receiver,
			Msg:      msg.Msg,
		}
		connectedUsers[msg.Receiver].WriteJSON(msgToSend)
	}
}
