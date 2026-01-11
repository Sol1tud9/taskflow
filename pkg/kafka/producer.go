package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/Sol1tud9/taskflow/pkg/logger"
	"go.uber.org/zap"
)

type Producer struct {
	writer *kafka.Writer
	topic  string
}

func NewProducer(brokers []string, topic string) *Producer {
	if topic == "" {
		logger.Error("kafka topic is empty", zap.Strings("brokers", brokers))
	}

	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
		RequiredAcks: kafka.RequireOne,
		WriteTimeout: 10 * time.Second,
	}

	logger.Info("kafka producer created", zap.String("topic", topic), zap.Strings("brokers", brokers))

	return &Producer{
		writer: writer,
		topic:  topic,
	}
}

func (p *Producer) Publish(ctx context.Context, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		logger.Error("failed to marshal message", zap.Error(err))
		return err
	}

	msg := kafka.Message{
		Key:   []byte(key),
		Value: data,
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		logger.Error("failed to write message to kafka", zap.Error(err), zap.String("topic", p.topic))
		return err
	}

	logger.Info("message published to kafka", zap.String("topic", p.topic), zap.String("key", key))
	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}

