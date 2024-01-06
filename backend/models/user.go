package models

import "github.com/gorilla/websocket"

type User struct {
	username string
	socket   *websocket.Conn
}
