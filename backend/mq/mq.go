package mq

import (
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
	return nil
}

func ReadFromQueue(email string) ([]models.Message, error) {
	conn, err := kafka.DialLeader(context.Background(), "tcp", url, emailToHash(email), 0)
	if err != nil {
		log.Println("E#1QX6IW - Didn't connect to kafka!", err)
		return nil, fmt.Errorf("failed to push to kafka. %v", err)
	}
	messages := []models.Message{}
	buffer := make([]byte, 10e3)

	conn.SetDeadline(time.Now().Add(10 * time.Second))
	batch := conn.ReadBatch(10e3, 1e6)

	for {
		n, err := batch.Read(buffer)
		if err != nil {
			break
		}
		var msg models.Message
		err = json.Unmarshal(buffer[:n], &msg)
		if err != nil {
			log.Println("E#1R2O6A - Error while unmarshaling a message from queue", err)
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

// emailToHash converts an email to a SHA256 hash and returns it as a hexadecimal string
func emailToHash(email string) string {
	emailBytes := []byte(email)
	hashBytes := sha256.Sum256(emailBytes)
	hashString := hex.EncodeToString(hashBytes[:])
	return hashString
}
