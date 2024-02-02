package main

import (
	"GoBaatcheet/models"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"
)

const Alice = "alice"
const Bob = "bob"

var wg sync.WaitGroup

func main() {
	//TestWebsocketOneToOneChat()
	TestIfMessagesArePushedToQueue()
}

func SetupDummyWebsocketConnections() (*websocket.Conn, *websocket.Conn) {
	// Establish Connection 1.
	ws1 := SetupWebsocketConnection(Alice)
	// Establish Connection 2.
	ws2 := SetupWebsocketConnection(Bob)
	return ws1, ws2
}

func SetupWebsocketConnection(username string) *websocket.Conn {
	url := "ws://localhost:8080/ws"
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		fmt.Println("E#1Q3MKH - Unable to connect to websocket server. E:", err)
	}
	_ = ws.WriteJSON(&models.User{Username: username})
	return ws
}

func ReadForSocket(ws *websocket.Conn) {
	for {
		var message models.Message
		err := ws.ReadJSON(&message)
		if err != nil {
			log.Println("error", err)
			break
		}
		log.Println(message.Msg)
		wg.Done()
	}
}

func WriteToSocket(sender *websocket.Conn, s, r, msg string) {
	err := sender.WriteJSON(models.Message{Sender: s, Receiver: r, Msg: msg})
	if err != nil {
		log.Println("error", err)
		return
	}
}

func TestWebsocketOneToOneChat() {
	ws1, ws2 := SetupDummyWebsocketConnections()
	wg.Add(2)
	go ReadForSocket(ws1)
	go ReadForSocket(ws2)
	WriteToSocket(ws1, Alice, Bob, "Alice say's hi!")
	WriteToSocket(ws2, Bob, Alice, "Bob say's hi in return!")
	wg.Wait()
}

func TestIfMessagesArePushedToQueue() {
	wsAlice := SetupWebsocketConnection(Alice)
	WriteToSocket(wsAlice, Alice, Bob, "Time "+time.Now().UTC().String())
	WriteToSocket(wsAlice, Alice, Bob, "Time "+time.Now().UTC().String())
	wsBob := SetupWebsocketConnection(Bob)
	wg.Add(2)
	go ReadForSocket(wsBob)
	wg.Wait()
}
