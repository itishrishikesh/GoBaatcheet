package mq

import (
	"GoBaatcheet/constants"
	"GoBaatcheet/models"
	"context"
	"crypto/sha256"
	"encoding/binary"
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
	queue, err := kafka.DialLeader(context.Background(), "tcp", url, emailToHash(email), 0)
	qConfirm, err := kafka.DialLeader(context.Background(), "tcp", url, emailToHash(email), 1)
	if err != nil {
		log.Println("E#1QX6IW - Didn't connect to kafka!", err)
		return nil, fmt.Errorf("failed to read from kafka. %v", err)
	}
	mConfirm, _ := qConfirm.ReadMessage(constants.MaxMessageSize) // Todo: This is incorrect, correct this.
	var messages = make([]models.Message, 0)
	_ = queue.SetDeadline(time.Now().Add(10 * time.Second))
	batch := queue.ReadBatch(constants.MinMessageSize, constants.MaxMessageSize)
	for m, err := batch.ReadMessage(); err != nil; {
		if m.Offset <= int64(binary.LittleEndian.Uint16(mConfirm.Value)) {
			continue
		}
		var msg models.Message
		err = json.Unmarshal(m.Value, &msg)
		if err != nil {
			log.Println("E#1R2O6A - Error while unmarshalling a message from queue", err)
		}
		messages = append(messages, msg)
		// Review here for more details: https://stackoverflow.com/a/35371760/5621167
		offset := make([]byte, binary.MaxVarintLen64)
		binary.LittleEndian.PutUint64(offset, uint64(m.Offset))
		// convertOffSetBack = int64(binary.LittleEndian.Uint64(offset))
		_, err := qConfirm.Write(offset)
		if err != nil {
			return nil, err
		}
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
