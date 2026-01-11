package publisher

import (
	"context"

	"github.com/Sol1tud9/taskflow/internal/domain"
	"github.com/Sol1tud9/taskflow/pkg/kafka"
)

type Publisher struct {
	userCreatedProducer *kafka.Producer
	userUpdatedProducer *kafka.Producer
	teamUpdatedProducer *kafka.Producer
}

func NewPublisher(brokers []string, topics map[string]string) *Publisher {
	return &Publisher{
		userCreatedProducer: kafka.NewProducer(brokers, topics["user_created"]),
		userUpdatedProducer: kafka.NewProducer(brokers, topics["user_updated"]),
		teamUpdatedProducer: kafka.NewProducer(brokers, topics["team_updated"]),
	}
}

func (p *Publisher) PublishUserCreated(ctx context.Context, event domain.UserCreatedEvent) error {
	return p.userCreatedProducer.Publish(ctx, event.UserID, event)
}

func (p *Publisher) PublishUserUpdated(ctx context.Context, event domain.UserUpdatedEvent) error {
	return p.userUpdatedProducer.Publish(ctx, event.UserID, event)
}

func (p *Publisher) PublishTeamUpdated(ctx context.Context, event domain.TeamUpdatedEvent) error {
	return p.teamUpdatedProducer.Publish(ctx, event.TeamID, event)
}

func (p *Publisher) Close() error {
	_ = p.userCreatedProducer.Close()
	_ = p.userUpdatedProducer.Close()
	_ = p.teamUpdatedProducer.Close()
	return nil
}

