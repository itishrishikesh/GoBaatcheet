package main

import (
	"fmt"
	"log"
	"net/http"

	"GoBaatcheet/handlers"
)

const port = "8080" // Todo: Move this to a centralized config

func main() {
	// Setup routes
	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/ws", handlers.WsEndpoint)
	fmt.Println("Server is Running on Port:", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
