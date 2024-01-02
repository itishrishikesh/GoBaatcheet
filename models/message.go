package models

type Message struct {
	Msg      string `json:"message"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
}
