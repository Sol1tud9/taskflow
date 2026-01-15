package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/Sol1tud9/taskflow/internal/domain"
	"github.com/Sol1tud9/taskflow/pkg/logger"
	"go.uber.org/zap"
)

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

type TeamEventPublisher interface {
	PublishTeamUpdated(ctx context.Context, event domain.TeamUpdatedEvent) error
}

type TeamUseCase struct {
	teamRepo       TeamRepository
	teamMemberRepo TeamMemberRepository
	publisher      TeamEventPublisher
}

func NewTeamUseCase(
	teamRepo TeamRepository,
	teamMemberRepo TeamMemberRepository,
	publisher TeamEventPublisher,
) *TeamUseCase {
	return &TeamUseCase{
		teamRepo:       teamRepo,
		teamMemberRepo: teamMemberRepo,
		publisher:      publisher,
	}
}

func (uc *TeamUseCase) CreateTeam(ctx context.Context, name, ownerID string) (*domain.Team, error) {
	now := time.Now()
	team := &domain.Team{
		ID:        uuid.New().String(),
		Name:      name,
		OwnerID:   ownerID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := uc.teamRepo.Create(ctx, team); err != nil {
		return nil, err
	}

	ownerMember := &domain.TeamMember{
		ID:       uuid.New().String(),
		TeamID:   team.ID,
		UserID:   ownerID,
		Role:     "owner",
		JoinedAt: now,
	}
	if err := uc.teamMemberRepo.Add(ctx, ownerMember); err != nil {
		return nil, err
	}

	event := domain.TeamUpdatedEvent{
		TeamID:    team.ID,
		Name:      team.Name,
		OwnerID:   team.OwnerID,
		UpdatedAt: team.UpdatedAt,
	}
	if err := uc.publisher.PublishTeamUpdated(ctx, event); err != nil {
		logger.Error("failed to publish team.updated event", zap.Error(err), zap.String("team_id", team.ID))
	}

	return team, nil
}

func (uc *TeamUseCase) GetTeam(ctx context.Context, id string) (*domain.Team, error) {
	return uc.teamRepo.GetByID(ctx, id)
}

func (uc *TeamUseCase) AddTeamMember(ctx context.Context, teamID, userID, role string) (*domain.TeamMember, error) {
	member := &domain.TeamMember{
		ID:       uuid.New().String(),
		TeamID:   teamID,
		UserID:   userID,
		Role:     role,
		JoinedAt: time.Now(),
	}

	if err := uc.teamMemberRepo.Add(ctx, member); err != nil {
		return nil, err
	}

	return member, nil
}

func (uc *TeamUseCase) GetTeamMembers(ctx context.Context, teamID string) ([]*domain.TeamMember, error) {
	return uc.teamMemberRepo.GetByTeamID(ctx, teamID)
}

