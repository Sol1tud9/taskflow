package mocks

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/Sol1tud9/taskflow/internal/domain"
)

type TeamMemberRepository struct {
	mock.Mock
}

func NewTeamMemberRepository(t testing.TB) *TeamMemberRepository {
	mock := &TeamMemberRepository{}
	mock.Mock.Test(t)
	return mock
}

func (m *TeamMemberRepository) Add(ctx context.Context, member *domain.TeamMember) error {
	args := m.Called(ctx, member)
	return args.Error(0)
}

func (m *TeamMemberRepository) GetByTeamID(ctx context.Context, teamID string) ([]*domain.TeamMember, error) {
	args := m.Called(ctx, teamID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.TeamMember), args.Error(1)
}

func (m *TeamMemberRepository) Remove(ctx context.Context, teamID, userID string) error {
	args := m.Called(ctx, teamID, userID)
	return args.Error(0)
}

