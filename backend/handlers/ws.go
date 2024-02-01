package handlers

import (
	"GoBaatcheet/auth"
	. "GoBaatcheet/client_manager"
	"GoBaatcheet/config"
	"GoBaatcheet/helpers"
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
		w.WriteHeader(helpers.HttpForbidden)
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
	go mq.ReadFromQueue(username, ConnectedUsers[username].Send)
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
