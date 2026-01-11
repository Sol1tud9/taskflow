package mocks

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/Sol1tud9/taskflow/internal/domain"
)

type UserRepository struct {
	mock.Mock
}

func NewUserRepository(t testing.TB) *UserRepository {
	mock := &UserRepository{}
	mock.Mock.Test(t)
	return mock
}

func (m *UserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *UserRepository) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *UserRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

