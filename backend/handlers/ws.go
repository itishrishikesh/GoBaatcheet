package handlers

import (
	"GoBaatcheet/auth"
	. "GoBaatcheet/client_manager"
	"GoBaatcheet/config"
	"GoBaatcheet/constants"
	"GoBaatcheet/models"
	"GoBaatcheet/mq"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var bypassForTest = true

func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	if !auth.Authenticate(r) && !bypassForTest {
		w.WriteHeader(constants.HttpForbidden)
	}
	config.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := config.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalln("E#1PZU6V - Failed to upgrade request to websocket connection!", err)
	}
	username, err := readOrAssignUsername(ws)
	if err != nil {
		log.Fatalln(err)
	}
	ConnectedUsers[username] = &Client{
		Conn:     ws,
		Send:     make(chan []byte),
		Username: username,
	}
	log.Println("Client is created.", ConnectedUsers[username])
	messages, err := mq.ReadFromQueue(username)
	if err != nil {
		log.Println("E#1R2MKV - Failed to read from queue", err)
	}
	for _, v := range messages {
		_ = SendMessage(v) // Todo: handle error
	}
	if err != nil {
		log.Println("E#1PZUJN - Error while writing message back to client!")
	}
	go ConnectedUsers[username].ReadPump()
	go ConnectedUsers[username].WritePump()
}

func readOrAssignUsername(conn *websocket.Conn) (string, error) {
	var user models.User
	err := conn.ReadJSON(&user)
	if err != nil {
		return "", fmt.Errorf("E#1QENH1 - Unable to read username from websocket connection. %v", err)
	}
	return user.Username, nil
}
