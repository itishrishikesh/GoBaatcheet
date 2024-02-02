package mq

import (
	"GoBaatcheet/helpers"
	"GoBaatcheet/models"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

const url = "localhost:9093" // Todo: Move this to a configuration server

func PushToQueue(message models.Message) error {
	conn, err := kafka.DialLeader(context.Background(), "tcp", url, emailToHash(message.Receiver), 0)
	if err != nil {
		log.Println("E#1QX6I2 - Didn't connect to kafka!", err)
		return fmt.Errorf("failed to push to kafka. %v", err)
	}
	tmp, _ := json.Marshal(message) // Todo: Handle error
	_, err = conn.Write(tmp)
	if err != nil {
		log.Println("E#1R2MS5 - Failed to write to Kafka.", err)
		return fmt.Errorf("failed to write message to kafka. %v", err)
	}
	_ = conn.Close()
	return nil
}

func ConnectToKafka(topic string) *kafka.Reader {
	hashedTopic := emailToHash(topic)
	_, err := kafka.DialLeader(context.Background(), "tcp", url, hashedTopic, 0)
	if err != nil {
		log.Println("E#1QX6IW - Didn't connect to kafka!", err)
	}
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{url},
		Topic:          hashedTopic,
		Partition:      0,
		MaxBytes:       helpers.MaxMessageSize,
		CommitInterval: time.Second,
		GroupID:        "message-reader",
	})

}

func ReadFromQueue(topic string, client chan []byte) {
	queue := ConnectToKafka(topic)
	for {
		message, err := queue.ReadMessage(context.Background())
		if err != nil {
			log.Println("Error while reading from Kafka for", topic)
			break
		}
		client <- message.Value
	}
}

func emailToHash(email string) string {
	emailBytes := []byte(email)
	hashBytes := sha256.Sum256(emailBytes)
	email = hex.EncodeToString(hashBytes[:])
	return email
}
