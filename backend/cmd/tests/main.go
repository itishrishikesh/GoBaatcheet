package main

import (
	"GoBaatcheet/models"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

func main() {
	TestWebsocketOneToOneChat()
}

func TestWebsocketOneToOneChat() {
	// Establish Connection 1.
	url := "ws://localhost:8080/ws"
	ws1, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		fmt.Println("E#1Q3MKH - Unable to connect to websocket server. E:", err)
	}
	// Establish Connection 2.
	ws2, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		fmt.Println("E#1Q3MKH - Unable to connect to websocket server. E:", err)
	}
	_ = ws1.WriteJSON(&models.User{Username: "websocket1"})
	_ = ws2.WriteJSON(&models.User{Username: "websocket2"})
	// Read Messages for Connection 1.
	go func() {
		for {
			var message models.Message
			err = ws1.ReadJSON(&message)
			log.Println("error", err)
			log.Println("websocket1 received:", message.Msg)
		}
	}()
	// Read Messages for Connection 2.
	go func() {
		for {
			var message models.Message
			err = ws2.ReadJSON(&message)
			log.Println("error", err)
			log.Println("websocket2 received:", message.Msg)
		}
	}()
	// Send Message from connection 1 to 2.
	err = ws1.WriteJSON(models.Message{Sender: "websocket1", Receiver: "websocket2", Msg: "Hello from websocket1"})
	// Send Message from connection 2 to 1.
	err = ws2.WriteJSON(models.Message{Sender: "websocket2", Receiver: "websocket1", Msg: "Hello from websocket2"})
	time.Sleep(5 * time.Second)
	// Close Connection 2.
	ws2.Close()
	// Send Message to closed connection 2 from 1.
	err = ws1.WriteJSON(models.Message{Sender: "websocket1", Receiver: "websocket2", Msg: "This message was pushed to queue"})
	log.Println("error", err)
	// Reconnect connection 2.
	ws2, _, err = websocket.DefaultDialer.Dial(url, nil)
	_ = ws2.WriteJSON(models.User{Username: "websocket2"})
	// Read Messages for Connection 2.
	//go func() {
	//	for {
	//		var message models.Message
	//		err = ws2.ReadJSON(&message)
	//		log.Println("error", err)
	//		log.Println("websocket2 received:", message.Msg)
	//	}
	//}()
}
