package client_manager

import (
	"GoBaatcheet/models"
	"GoBaatcheet/mq"
	"encoding/json"
	"log"
)

var ConnectedUsers = make(map[string]*Client)

func SendMessage(message models.Message) error {
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
