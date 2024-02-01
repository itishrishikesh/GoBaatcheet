package client_manager

import (
	"GoBaatcheet/models"
	"GoBaatcheet/mq"
	"bytes"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	Conn *websocket.Conn
	// Buffered channel of outbound messages.
	Send     chan []byte
	Username string
	Stop     chan bool
}

func (c *Client) ReadPump() {
	defer func() {
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for ConnectedUsers[c.Username] != nil {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("Connection is closed.", err)
			ConnectedUsers[c.Username] = nil
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		var msg models.Message
		err = json.Unmarshal(message, &msg)
		if err != nil {
			log.Println("Marshalling error.", err)
		}
		if ConnectedUsers[msg.Receiver] != nil {
			ConnectedUsers[msg.Receiver].Send <- message
		} else {
			_ = mq.PushToQueue(msg)
		}
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for c != nil {
		select {
		case message, ok := <-c.Send:
			if ConnectedUsers[c.Username] == nil {
				_ = SendMessage(models.FromBytes(message))
				continue
			}
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, _ = w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.Send)
			for i := 0; i < n; i++ {
				_, _ = w.Write(newline)
				_, _ = w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				ConnectedUsers[c.Username] = nil
				//return
			}
		}
	}
}
