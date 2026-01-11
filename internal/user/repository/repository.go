package repository

import (
	"context"

	"github.com/Sol1tud9/taskflow/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id string) error
}

type TeamRepository interface {
	Create(ctx context.Context, team *domain.Team) error
	GetByID(ctx context.Context, id string) (*domain.Team, error)
	Update(ctx context.Context, team *domain.Team) error
	Delete(ctx context.Context, id string) error
}

type TeamMemberRepository interface {
	Add(ctx context.Context, member *domain.TeamMember) error
	GetByTeamID(ctx context.Context, teamID string) ([]*domain.TeamMember, error)
	Remove(ctx context.Context, teamID, userID string) error
}

