package consumer

import (
	"context"
	"log"

	"github.com/amogh-sood/distributed-autoscaler/worker/internal/hashing"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(broker, topic, groupID string) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   topic,
		GroupID: groupID,
	})

	return &Consumer{reader: r}
}

func (c *Consumer) Start() {
	log.Println("Worker: listening for jobs...")

	for {
		msg, err := c.reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error reading message: %v\n", err)
			continue
		}

		input := string(msg.Value)
		result, _ := hashing.Mine(input)

		log.Printf("Worker: hashed '%s' → %s\n", input, result)
	}
}
