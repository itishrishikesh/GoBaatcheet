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

func ConnectToKafka(topic string) (queue, qConfirm *kafka.Conn) {
	// I'm creating two queues, one for storing offline messages and one which will store the last updated offset.
	// queue --> Stores messages pushed for an offline username.
	queue, err := kafka.DialLeader(context.Background(), "tcp", url, emailToHash(topic), 0)
	// qConfirm --> Stores offset of the last read message from queue above.
	qConfirm, err = kafka.DialPartition(context.Background(), "tcp", url, kafka.Partition{Topic: emailToHash(topic), ID: 1, Leader: queue.Broker()})
	if err != nil {
		log.Println("E#1QX6IW - Didn't connect to kafka!", err)
		//return nil, fmt.Errorf("failed to read from kafka. %v", err)
	}
	return
}

func ReadFromQueue(topic string) ([]models.Message, error) {
	queue, qConfirm := ConnectToKafka(topic)
	// I'm checking whether we have message in the qConfirm or not
	var mConfirm kafka.Message
	if i, _ := qConfirm.ReadLastOffset(); i > 0 {
		mConfirm, _ = qConfirm.ReadMessage(constants.MaxMessageSize) // Todo: This is incorrect, correct this.
	} else {
		mConfirm = kafka.Message{Offset: 0}
	}

	// Read messages from queue.
	// messages variable will store the final list of messages that will be sent to user.
	var messages = make([]models.Message, 0)
	// setting up timeout for queue
	_ = queue.SetDeadline(time.Now().Add(10 * time.Second))
	// read batch of messages from queue.
	batch := queue.ReadBatch(constants.MinMessageSize, constants.MaxMessageSize)
	// iterate over all messages received.
	for m, err := batch.ReadMessage(); err != nil && m.Offset > 0; {
		// check if the offset of current message is less than offset last stored in mConfirm queue.
		// this tells us whether the message is already read by the user.
		if m.Offset <= int64(binary.LittleEndian.Uint16(mConfirm.Value)) {
			continue
		}
		// temp variable to store the current message value
		var msg models.Message
		// read the message value
		err = json.Unmarshal(m.Value, &msg)
		if err != nil {
			log.Println("E#1R2O6A - Error while unmarshalling a message from queue", err)
		}
		// append the message value in the final list
		messages = append(messages, msg)
		// store the offset of current message and push it mConfirm indicating that current message has been read.
		// Review here for more details: https://stackoverflow.com/a/35371760/5621167
		offset := make([]byte, binary.MaxVarintLen64)
		binary.LittleEndian.PutUint64(offset, uint64(m.Offset))
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
