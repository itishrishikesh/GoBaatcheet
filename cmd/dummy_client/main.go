package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/gorilla/websocket"
	"github.com/itishrishikesh/GoBaatcheet/models"
)

func main() {
	// Establish Connection.
	var url string
	fmt.Println("Specify a server url or type 'local'")
	fmt.Scan(&url)
	if url == "local" {
		url = "ws://localhost:8080/ws"
	}
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		fmt.Println("E#1Q3MKH - Unable to connect to websocket server. E:", err)
	}
	fmt.Println("Whatever you'll type on console will be sent to the chat server!")

	// Read Confirmation String.
	go func() {
		for {
			var message models.Message
			err := conn.ReadJSON(&message)
			if err != nil {
				fmt.Println("E#1Q3MWU - Failed to read messge from server. E:", err)
				return
			}
			fmt.Printf("\n%s: %s\n", message.Sender, message.Msg)
		}
	}()

	// Send whatever is being typed in console to everyone connected.
	for {
		line, _, _ := bufio.NewReader(os.Stdin).ReadLine()
		err := conn.WriteMessage(1, line)
		if err != nil {
			fmt.Println("E#1Q5F07 - Failed to write message to server. E:", err)
		}
	}
}
