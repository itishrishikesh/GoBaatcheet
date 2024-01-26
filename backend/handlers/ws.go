package handlers

import (
	"GoBaatcheet/auth"
	"GoBaatcheet/config"
	"GoBaatcheet/constants"
	"GoBaatcheet/models"
	"GoBaatcheet/mq"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var ConnectedUsers = make(map[string]*Client)

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
	if err != nil {
		log.Fatalln(err)
	}
	ConnectedUsers[username] = &Client{
		Conn:     ws,
		Send:     make(chan []byte),
		Username: username,
	}

	go ConnectedUsers[username].ReadPump()
	go ConnectedUsers[username].WritePump()

	messages, err := mq.ReadFromQueue(username)
	if err != nil {
		log.Println("E#1R2MKV - Failed to read from queue", err)
	}
	for _, v := range messages {
		_ = sendMessage(v) // Todo: handle error
	}
	if err != nil {
		log.Println("E#1PZUJN - Error while writing message back to client!")
	}
	//messageListener(ws)
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
		err = sendMessage(msgToSend)
		if err != nil {
			log.Println("E#1QZ3SJ - Error while sending the message", err)
		}
	}
}

func sendMessage(message models.Message) error {
	if ConnectedUsers[message.Receiver] != nil {
		ConnectedUsers[message.Receiver] = nil
	}
	if ConnectedUsers[message.Receiver] != nil {
		// err := ConnectedUsers[message.Receiver].Conn.WriteJSON(message)
		b, err := json.Marshal(message)
		ConnectedUsers[message.Receiver].Send <- b
		if err != nil {
			log.Println("E#1R2MTS - Failed to write JSON to websocket.", err)
		}
	} else {
		// Todo: check if message is for a valid user
		_ = mq.PushToQueue(message) // Todo: handle error
	}
	return nil
}
