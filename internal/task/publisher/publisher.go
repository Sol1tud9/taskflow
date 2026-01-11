package publisher

import (
	"context"

	"github.com/Sol1tud9/taskflow/internal/domain"
	"github.com/Sol1tud9/taskflow/pkg/kafka"
)

type Publisher struct {
	taskCreatedProducer *kafka.Producer
	taskUpdatedProducer *kafka.Producer
}

func NewPublisher(brokers []string, topics map[string]string) *Publisher {
	return &Publisher{
		taskCreatedProducer: kafka.NewProducer(brokers, topics["task_created"]),
		taskUpdatedProducer: kafka.NewProducer(brokers, topics["task_updated"]),
	}
}

func (p *Publisher) PublishTaskCreated(ctx context.Context, event domain.TaskCreatedEvent) error {
	return p.taskCreatedProducer.Publish(ctx, event.TaskID, event)
}

func (p *Publisher) PublishTaskUpdated(ctx context.Context, event domain.TaskUpdatedEvent) error {
	return p.taskUpdatedProducer.Publish(ctx, event.TaskID, event)
}

func (p *Publisher) Close() error {
	_ = p.taskCreatedProducer.Close()
	_ = p.taskUpdatedProducer.Close()
	return nil
}

