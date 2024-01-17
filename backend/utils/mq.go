package mq

import (
	"GoBaatcheet/models"
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

const url = "localhost:9093" //Todo: Move this configuration

func PushToQueue(msg models.Message, topic string) error {
	conn, err := kafka.DialLeader(context.Background(), "tcp", url, topic, 0)
	if err != nil {
		log.Println("Didn't connect to kafka!")
	}
	return nil
}

func ReadFromQueue(topic string) (models.Message, error) {
	return models.Message{}, nil
}
