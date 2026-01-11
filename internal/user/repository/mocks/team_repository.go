package mocks

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/Sol1tud9/taskflow/internal/domain"
)

type TeamRepository struct {
	mock.Mock
}

func NewTeamRepository(t testing.TB) *TeamRepository {
	mock := &TeamRepository{}
	mock.Mock.Test(t)
	return mock
}

func (m *TeamRepository) Create(ctx context.Context, team *domain.Team) error {
	args := m.Called(ctx, team)
	return args.Error(0)
}

func (m *TeamRepository) GetByID(ctx context.Context, id string) (*domain.Team, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Team), args.Error(1)
}

func (m *TeamRepository) Update(ctx context.Context, team *domain.Team) error {
	args := m.Called(ctx, team)
	return args.Error(0)
}

func (m *TeamRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

