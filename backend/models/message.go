package models

import (
	"encoding/json"
	"log"
)

type Message struct {
	Msg      string `json:"message"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
}

func (m *Message) ToBytes() []byte {
	b, err := json.Marshal(m)
	if err != nil {
		log.Panic("error occurred while marshalling message", err)
	}
	return b
}

func FromBytes(b []byte) Message {
	var result Message
	err := json.Unmarshal(b, &result)
	if err != nil {
		log.Panic("error occurred while unmarshalling message")
	}
	return result
}
