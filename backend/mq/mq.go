package mq

import (
	"GoBaatcheet/models"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

const url = "localhost:9093" // Todo: Move this to a configuration server

func PushToQueue(message models.Message) error {
	conn, err := kafka.DialLeader(context.Background(), "tcp", url, message.Receiver, 0)
	if err != nil {
		log.Println("Didn't connect to kafka!")
		return fmt.Errorf("failed to push to kafka. %v", err)
	}
	tmp, _ := json.Marshal(message) // Todo: Handle error
	_, _ = conn.Write(tmp)          // Todo: Handle error
	return nil
}

func ReadFromQueue(topic string) ([]models.Message, error) {
	conn, err := kafka.DialLeader(context.Background(), "tcp", url, topic, 0)
	if err != nil {
		log.Println("Didn't connect to kafka!")
		return nil, fmt.Errorf("failed to push to kafka. %v", err)
	}
	messages := []models.Message{}
	buffer := make([]byte, 10e3)
	for _, err := conn.Read(buffer); err != nil; {
		var msg models.Message
		_ = json.Unmarshal(buffer, &msg) // Todo: Handle error
		messages = append(messages, msg)
	}
	return messages, nil
}
