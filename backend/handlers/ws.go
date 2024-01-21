package handlers

import (
	"fmt"
	"log"
	"net/http"

	"GoBaatcheet/auth"
	"GoBaatcheet/config"
	"GoBaatcheet/constants"
	"GoBaatcheet/models"
	"GoBaatcheet/mq"

	"github.com/gorilla/websocket"
)

var ConnectedUsers map[string]*websocket.Conn = make(map[string]*websocket.Conn)

func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	if !auth.Authenticate(r) {
		w.WriteHeader(constants.HTTP_FORBIDDEN)
	}
	config.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := config.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalln("E#1PZU6V - Failed to upgrade request to websocket connection!", err)
	}
	username, err := readOrAssignUsername(ws)
	// if !mq.DoesTopicExists(username) {
	// 	_ = mq.CreateTopic(username) // Todo: handle error
	// }
	if err != nil {
		log.Fatalln(err)
	}
	ConnectedUsers[username] = ws
	messages, _ := mq.ReadFromQueue(username) // Todo: handle error
	for _, v := range messages {
		_ = sendMessage(v) // Todo: handle error
	}
	if err != nil {
		log.Println("E#1PZUJN - Error while writing message back to client!")
	}
	messageListener(ws)
}

func readOrAssignUsername(conn *websocket.Conn) (string, error) {
	var user models.User
	err := conn.ReadJSON(&user)
	if err != nil {
		return "", fmt.Errorf("E#1QENH1 - Unable to read username from websocket connection. %v", err)
	}
	return user.Username, nil
}

func messageListener(conn *websocket.Conn) {
	for {
		var msg models.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("E#1PZUA7 - Error while reading message for user:", err)
			return
		}
		msgToSend := models.Message{
			Sender:   msg.Sender,
			Receiver: msg.Receiver,
			Msg:      msg.Msg,
		}
		err = sendMessage(msgToSend) // Todo: handle error
		if err != nil {
			log.Println("E#1QZ3SJ - Error while sending the message", err)
		}
	}
}

func sendMessage(message models.Message) error {
	if ConnectedUsers[message.Receiver] != nil && !isConnectionAlive(ConnectedUsers[message.Receiver]) {
		ConnectedUsers[message.Receiver] = nil
	}
	if ConnectedUsers[message.Receiver] != nil {
		ConnectedUsers[message.Receiver].WriteJSON(message)
	} else {
		// Todo: check if message is for a valid user
		_ = mq.PushToQueue(message) // Todo: handle error
	}
	return nil
}

func isConnectionAlive(c *websocket.Conn) bool {
	if err := c.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
		return false
	}
	if mType, _, err := c.ReadMessage(); err != nil || mType != websocket.PongMessage {
		return false
	}
	return true
}
