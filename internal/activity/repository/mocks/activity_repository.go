package mocks

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/Sol1tud9/taskflow/internal/activity/repository"
	"github.com/Sol1tud9/taskflow/internal/domain"
)

type ActivityRepository struct {
	mock.Mock
}

func NewActivityRepository(t testing.TB) *ActivityRepository {
	mock := &ActivityRepository{}
	mock.Mock.Test(t)
	return mock
}

func (m *ActivityRepository) Create(ctx context.Context, activity *domain.Activity) error {
	args := m.Called(ctx, activity)
	return args.Error(0)
}

func (m *ActivityRepository) GetByUserID(ctx context.Context, userID string, filter repository.ActivityFilter) ([]*domain.Activity, int, error) {
	args := m.Called(ctx, userID, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int), args.Error(2)
	}
	return args.Get(0).([]*domain.Activity), args.Get(1).(int), args.Error(2)
}

func (m *ActivityRepository) GetByEntity(ctx context.Context, entityType, entityID string, filter repository.ActivityFilter) ([]*domain.Activity, int, error) {
	args := m.Called(ctx, entityType, entityID, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int), args.Error(2)
	}
	return args.Get(0).([]*domain.Activity), args.Get(1).(int), args.Error(2)
}

func (m *ActivityRepository) GetAll(ctx context.Context, filter repository.ActivityFilter) ([]*domain.Activity, int, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int), args.Error(2)
	}
	return args.Get(0).([]*domain.Activity), args.Get(1).(int), args.Error(2)
}

