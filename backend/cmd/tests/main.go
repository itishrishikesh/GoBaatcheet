package main

import (
	"GoBaatcheet/models"
	"GoBaatcheet/mq"
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
	_ = ws1.WriteJSON(&models.User{Username: Alice})
	_ = ws2.WriteJSON(&models.User{Username: Bob})
	return ws1, ws2
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
	wsAlice, wsBob := SetupDummyWebsocketConnections()
	_ = wsBob.Close()
	WriteToSocket(wsAlice, Alice, Bob, "Alice sends message to bob, when he's offline!")
	//queue, _ := mq.ConnectToKafka(Bob)
	//queue.SetDeadline(time.Now().Add(2 * time.Second))
	//message, err := queue.ReadMessage(constants.MaxMessageSize)
	//if err != nil {
	//	log.Println("error while reading from kafka", err)
	//	return
	//}
	//log.Println("Message read from kafka (should be for bob)", message)
	for fromQueue, err := mq.ReadFromQueue(Bob); len(fromQueue) == 0 && err == nil; {
		log.Println("fromQueue", fromQueue)
		time.Sleep(time.Second)
	}
}
