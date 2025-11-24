package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(broker string, topic string) *Producer {
	w := &kafka.Writer{
		Addr:     kafka.TCP(broker),
		Topic:    topic,
		Balancer: &kafka.RoundRobin{},
	}

	return &Producer{writer: w}
}

// Publish sends a message to Kafka
func (p *Producer) Publish(msg string) error {
	err := p.writer.WriteMessages(
		context.Background(),
		kafka.Message{
			Value: []byte(msg),
		},
	)

	if err != nil {
		log.Println("Kafka publish error:", err)
	}
	return err
}
