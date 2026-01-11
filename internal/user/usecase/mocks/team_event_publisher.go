package mocks

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/Sol1tud9/taskflow/internal/domain"
)

type TeamEventPublisher struct {
	mock.Mock
}

func NewTeamEventPublisher(t testing.TB) *TeamEventPublisher {
	mock := &TeamEventPublisher{}
	mock.Mock.Test(t)
	return mock
}

func (m *TeamEventPublisher) PublishTeamUpdated(ctx context.Context, event domain.TeamUpdatedEvent) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

