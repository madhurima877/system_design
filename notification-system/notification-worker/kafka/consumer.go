package kafka

import (
	"context"
	"encoding/json"
	"log"
	"notification-system/models"
	"sync"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer() *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:29092"},
		Topic:   "notifications",
		GroupID: "notification-workers",
	})
	return &Consumer{
		reader: reader,
	}
}
func (c *Consumer) Worker(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := c.reader.ReadMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				continue
			}
			var event models.Message
			err = json.Unmarshal(msg.Value, &event)
			if err != nil {
				continue
			}
			log.Println(event)
		}
	}
}
func (c *Consumer) Close() error {
	return c.reader.Close()
}
