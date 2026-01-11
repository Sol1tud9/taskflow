package kafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/Sol1tud9/taskflow/pkg/logger"
	"go.uber.org/zap"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(brokers []string, topic string, groupID string) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:           brokers,
		GroupID:           groupID,
		Topic:             topic,
		HeartbeatInterval: 3 * time.Second,
		SessionTimeout:    30 * time.Second,
		MinBytes:          1,
		MaxBytes:          10e6,
	})

	return &Consumer{
		reader: reader,
	}
}

func (c *Consumer) Read(ctx context.Context) (kafka.Message, error) {
	msg, err := c.reader.ReadMessage(ctx)
	if err != nil {
		logger.Error("failed to read message from kafka", zap.Error(err))
		return kafka.Message{}, err
	}
	return msg, nil
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}

