package mocks

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/Sol1tud9/taskflow/internal/domain"
)

type TaskHistoryRepository struct {
	mock.Mock
}

func NewTaskHistoryRepository(t testing.TB) *TaskHistoryRepository {
	mock := &TaskHistoryRepository{}
	mock.Mock.Test(t)
	return mock
}

func (m *TaskHistoryRepository) Create(ctx context.Context, history *domain.TaskHistory) error {
	args := m.Called(ctx, history)
	return args.Error(0)
}

func (m *TaskHistoryRepository) GetByTaskID(ctx context.Context, taskID string) ([]*domain.TaskHistory, error) {
	args := m.Called(ctx, taskID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.TaskHistory), args.Error(1)
}

