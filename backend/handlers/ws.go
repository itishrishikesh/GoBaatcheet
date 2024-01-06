package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"GoBaatcheet/config"
	"GoBaatcheet/models"

	"github.com/goombaio/namegenerator"
	"github.com/gorilla/websocket"
)

var users map[*websocket.Conn]string = make(map[*websocket.Conn]string)

func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	config.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := config.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("E#1PZU6V - Failed to upgrade request to websocket connection!")
	}
	username := getRandomUsername()
	users[ws] = username
	fmt.Println("I#1PZUGS - Websocket Connection is assigned a random username:", username)
	if err != nil {
		fmt.Println("E#1PZUJN - Error while writing message back to client!")
	}
	reader(ws)
}

func getRandomUsername() string {
	seed := time.Now().UTC().UnixNano()
	nameGenerator := namegenerator.NewNameGenerator(seed)
	return nameGenerator.Generate()
}

func reader(conn *websocket.Conn) {
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("E#1PZUA7 - Error while reading message!")
			return
		}
		for receiver, username := range users {
			if strings.EqualFold(username, users[conn]) {
				continue
			}
			message := models.Message{
				Msg:      string(p),
				Sender:   users[conn],
				Receiver: username,
			}
			receiver.WriteJSON(message)
		}
	}
}
