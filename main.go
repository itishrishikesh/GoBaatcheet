package main

import (
	"log"
	"net/http"

	"github.com/itishrishikesh/GoBaatcheet/handlers"
)

const port = "8080"

func main() {
	// Setup routes
	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/ws", handlers.WsEndpoint)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
