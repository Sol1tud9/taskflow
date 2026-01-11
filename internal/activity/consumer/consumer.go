package consumer

import (
	"context"
	"encoding/json"

	"github.com/Sol1tud9/taskflow/internal/domain"
	"github.com/Sol1tud9/taskflow/pkg/kafka"
	"github.com/Sol1tud9/taskflow/pkg/logger"
	"go.uber.org/zap"
)

type ActivityRecorder interface {
	RecordUserCreated(ctx context.Context, event domain.UserCreatedEvent) error
	RecordUserUpdated(ctx context.Context, event domain.UserUpdatedEvent) error
	RecordTaskCreated(ctx context.Context, event domain.TaskCreatedEvent) error
	RecordTaskUpdated(ctx context.Context, event domain.TaskUpdatedEvent) error
}

type EventConsumer struct {
	userCreatedConsumer *kafka.Consumer
	userUpdatedConsumer *kafka.Consumer
	taskCreatedConsumer *kafka.Consumer
	taskUpdatedConsumer *kafka.Consumer
	recorder            ActivityRecorder
}

func NewEventConsumer(brokers []string, topics map[string]string, groupID string, recorder ActivityRecorder) *EventConsumer {
	return &EventConsumer{
		userCreatedConsumer: kafka.NewConsumer(brokers, topics["user_created"], groupID),
		userUpdatedConsumer: kafka.NewConsumer(brokers, topics["user_updated"], groupID),
		taskCreatedConsumer: kafka.NewConsumer(brokers, topics["task_created"], groupID),
		taskUpdatedConsumer: kafka.NewConsumer(brokers, topics["task_updated"], groupID),
		recorder:            recorder,
	}
}

func (c *EventConsumer) Start(ctx context.Context) {
	go c.consumeUserCreated(ctx)
	go c.consumeUserUpdated(ctx)
	go c.consumeTaskCreated(ctx)
	go c.consumeTaskUpdated(ctx)
}

func (c *EventConsumer) consumeUserCreated(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := c.userCreatedConsumer.Read(ctx)
			if err != nil {
				continue
			}

			var event domain.UserCreatedEvent
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				logger.Error("failed to unmarshal user created event", zap.Error(err))
				continue
			}

			if err := c.recorder.RecordUserCreated(ctx, event); err != nil {
				logger.Error("failed to record user created activity", zap.Error(err))
			}
		}
	}
}

func (c *EventConsumer) consumeUserUpdated(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := c.userUpdatedConsumer.Read(ctx)
			if err != nil {
				continue
			}

			var event domain.UserUpdatedEvent
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				logger.Error("failed to unmarshal user updated event", zap.Error(err))
				continue
			}

			if err := c.recorder.RecordUserUpdated(ctx, event); err != nil {
				logger.Error("failed to record user updated activity", zap.Error(err))
			}
		}
	}
}

func (c *EventConsumer) consumeTaskCreated(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := c.taskCreatedConsumer.Read(ctx)
			if err != nil {
				continue
			}

			var event domain.TaskCreatedEvent
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				logger.Error("failed to unmarshal task created event", zap.Error(err))
				continue
			}

			if err := c.recorder.RecordTaskCreated(ctx, event); err != nil {
				logger.Error("failed to record task created activity", zap.Error(err))
			}
		}
	}
}

func (c *EventConsumer) consumeTaskUpdated(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := c.taskUpdatedConsumer.Read(ctx)
			if err != nil {
				continue
			}

			var event domain.TaskUpdatedEvent
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				logger.Error("failed to unmarshal task updated event", zap.Error(err))
				continue
			}

			if err := c.recorder.RecordTaskUpdated(ctx, event); err != nil {
				logger.Error("failed to record task updated activity", zap.Error(err))
			}
		}
	}
}

func (c *EventConsumer) Close() error {
	_ = c.userCreatedConsumer.Close()
	_ = c.userUpdatedConsumer.Close()
	_ = c.taskCreatedConsumer.Close()
	_ = c.taskUpdatedConsumer.Close()
	return nil
}

