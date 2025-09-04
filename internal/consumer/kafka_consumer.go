package consumer

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

//go:generate mockgen -source=kafka_consumer.go -destination=../../test/mocks/kafka_consumer_mock.go

type Consumer interface {
	StartConsuming(processFunc func(message []byte) error)
	Close() error
}

type ConsumerImpl struct {
	Reader *kafka.Reader
}

func NewConsumer() *ConsumerImpl {
	return &ConsumerImpl{
		Reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{"localhost:9092"},
			Topic:   "order",
			GroupID: "order-service-group",
		}),
	}
}

func (c *ConsumerImpl) StartConsuming(processFunc func(message []byte) error) {
	go func() {
		for {
			msg, err := c.Reader.ReadMessage(context.Background())
			if err != nil {
				log.Printf("kafka error: %v", err)
				continue
			}

			if err := processFunc(msg.Value); err != nil {
				log.Printf("processing error: %v", err)
			}
		}
	}()
}

func (c *ConsumerImpl) Close() error {
	if err := c.Reader.Close(); err != nil {
		return fmt.Errorf("failed to close Kafka consumer: %w", err)
	}
	return nil
}
