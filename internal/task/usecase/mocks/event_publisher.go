package mocks

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/Sol1tud9/taskflow/internal/domain"
)

type EventPublisher struct {
	mock.Mock
}

func NewEventPublisher(t testing.TB) *EventPublisher {
	mock := &EventPublisher{}
	mock.Mock.Test(t)
	return mock
}

func (m *EventPublisher) PublishTaskCreated(ctx context.Context, event domain.TaskCreatedEvent) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *EventPublisher) PublishTaskUpdated(ctx context.Context, event domain.TaskUpdatedEvent) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

