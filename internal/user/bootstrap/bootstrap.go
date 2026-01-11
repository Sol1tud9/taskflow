package bootstrap

import (
	"context"

	"github.com/Sol1tud9/taskflow/internal/domain"
	"github.com/Sol1tud9/taskflow/internal/user/publisher"
	"github.com/Sol1tud9/taskflow/internal/user/storage/postgres"
	"github.com/Sol1tud9/taskflow/internal/user/usecase"
	"github.com/Sol1tud9/taskflow/pkg/config"
	"github.com/Sol1tud9/taskflow/pkg/logger"
	"go.uber.org/zap"
)

type App struct {
	Config    *config.UserServiceConfig
	Storage   *postgres.Storage
	Publisher *publisher.Publisher
	UserUC    *usecase.UserUseCase
	TeamUC    *usecase.TeamUseCase
}

func NewApp(cfg *config.UserServiceConfig) (*App, error) {
	storage, err := postgres.NewStorage(cfg.Database)
	if err != nil {
		logger.Error("failed to init storage", zap.Error(err))
		return nil, err
	}

	pub := publisher.NewPublisher(cfg.Kafka.Brokers, cfg.Kafka.Topics)

	userUC := usecase.NewUserUseCase(storage, pub)

	teamRepoAdapter := &teamRepoAdapter{storage: storage}
	teamMemberRepoAdapter := &teamMemberRepoAdapter{storage: storage}
	teamUC := usecase.NewTeamUseCase(teamRepoAdapter, teamMemberRepoAdapter, pub)

	return &App{
		Config:    cfg,
		Storage:   storage,
		Publisher: pub,
		UserUC:    userUC,
		TeamUC:    teamUC,
	}, nil
}

func (a *App) Close() {
	a.Storage.Close()
	_ = a.Publisher.Close()
}

type teamRepoAdapter struct {
	storage *postgres.Storage
}

func (a *teamRepoAdapter) Create(ctx context.Context, team *domain.Team) error {
	return a.storage.CreateTeam(ctx, team)
}

func (a *teamRepoAdapter) GetByID(ctx context.Context, id string) (*domain.Team, error) {
	return a.storage.GetTeamByID(ctx, id)
}

func (a *teamRepoAdapter) Update(ctx context.Context, team *domain.Team) error {
	return a.storage.UpdateTeam(ctx, team)
}

func (a *teamRepoAdapter) Delete(ctx context.Context, id string) error {
	return a.storage.DeleteTeam(ctx, id)
}

type teamMemberRepoAdapter struct {
	storage *postgres.Storage
}

func (a *teamMemberRepoAdapter) Add(ctx context.Context, member *domain.TeamMember) error {
	return a.storage.AddTeamMember(ctx, member)
}

func (a *teamMemberRepoAdapter) GetByTeamID(ctx context.Context, teamID string) ([]*domain.TeamMember, error) {
	return a.storage.GetTeamMembersByTeamID(ctx, teamID)
}

func (a *teamMemberRepoAdapter) Remove(ctx context.Context, teamID, userID string) error {
	return a.storage.RemoveTeamMember(ctx, teamID, userID)
}

