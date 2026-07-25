package kafka

import (
	"context"
	"encoding/json"
	"notification-system/models"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer() *Producer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:29092"},
		Topic:   "notifications",
	})
	return &Producer{writer: writer}
}
func (p *Producer) Publish(ctx context.Context, userID, message string) error {
	msg := models.Message{
		UserId:  userID,
		Message: message,
	}
	msgbyte, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Value: msgbyte,
	})

}
func (p *Producer) Close() error {
	return p.writer.Close()
}
