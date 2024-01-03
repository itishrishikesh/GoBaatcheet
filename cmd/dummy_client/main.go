package main

import (
	"fmt"

	"github.com/gorilla/websocket"
)

func main() {
	// Establish Connection
	var url string
	fmt.Println("Specify a server url or type 'local'")
	fmt.Scan(&url)
	if url == "local" {
		url = "ws://localhost:8080/ws"
	}
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		fmt.Println("E#1Q3MKH - Unable to connect to websocket server")
	}

	// Read Confirmation String
	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil || msg == nil {
				fmt.Println("E#1Q3MWU - Failed to read messge from server")
			}
			fmt.Println("From Server:", string(msg))
		}
	}()

	// Send whatever is being typed back in console back.
	for {
		var msg string
		fmt.Scan(&msg)
		conn.WriteJSON(string(msg))
	}
}
