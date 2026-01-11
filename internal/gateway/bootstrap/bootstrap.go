package bootstrap

import (
	"context"

	activityStorage "github.com/Sol1tud9/taskflow/internal/activity/storage/sharded"
	activityUsecase "github.com/Sol1tud9/taskflow/internal/activity/usecase"
	"github.com/Sol1tud9/taskflow/internal/domain"
	"github.com/Sol1tud9/taskflow/internal/gateway/cache"
	"github.com/Sol1tud9/taskflow/internal/gateway/handler"
	taskPublisher "github.com/Sol1tud9/taskflow/internal/task/publisher"
	taskStorage "github.com/Sol1tud9/taskflow/internal/task/storage/postgres"
	taskUsecase "github.com/Sol1tud9/taskflow/internal/task/usecase"
	userPublisher "github.com/Sol1tud9/taskflow/internal/user/publisher"
	userStorage "github.com/Sol1tud9/taskflow/internal/user/storage/postgres"
	userUsecase "github.com/Sol1tud9/taskflow/internal/user/usecase"
	"github.com/Sol1tud9/taskflow/pkg/config"
	"github.com/Sol1tud9/taskflow/pkg/logger"
	"go.uber.org/zap"
)

type App struct {
	Config          *config.GatewayConfig
	Cache           *cache.RedisCache
	Handler         *handler.Handler
	UserStorage     *userStorage.Storage
	TaskStorage     *taskStorage.Storage
	ActivityStorage *activityStorage.ShardedStorage
	UserPublisher   *userPublisher.Publisher
	TaskPublisher   *taskPublisher.Publisher
}

func NewApp(cfg *config.GatewayConfig) (*App, error) {
	redisCache := cache.NewRedisCache(cfg.Redis)

	userStore, err := userStorage.NewStorage(cfg.UserDB)
	if err != nil {
		logger.Error("failed to init user storage", zap.Error(err))
		return nil, err
	}

	taskStore, err := taskStorage.NewStorage(cfg.TaskDB)
	if err != nil {
		logger.Error("failed to init task storage", zap.Error(err))
		return nil, err
	}

	activityStore, err := activityStorage.NewShardedStorage(cfg.ActivityDB)
	if err != nil {
		logger.Error("failed to init activity storage", zap.Error(err))
		return nil, err
	}

	userTopics := map[string]string{
		"user_created": cfg.Kafka.Topics["user_created"],
		"user_updated": cfg.Kafka.Topics["user_updated"],
		"team_updated": cfg.Kafka.Topics["team_updated"],
	}
	taskTopics := map[string]string{
		"task_created": cfg.Kafka.Topics["task_created"],
		"task_updated": cfg.Kafka.Topics["task_updated"],
	}

	userPub := userPublisher.NewPublisher(cfg.Kafka.Brokers, userTopics)
	taskPub := taskPublisher.NewPublisher(cfg.Kafka.Brokers, taskTopics)

	userUC := userUsecase.NewUserUseCase(userStore, userPub)
	teamRepoAdapter := &teamRepoAdapter{storage: userStore}
	teamMemberRepoAdapter := &teamMemberRepoAdapter{storage: userStore}
	teamUC := userUsecase.NewTeamUseCase(teamRepoAdapter, teamMemberRepoAdapter, userPub)

	historyRepoAdapter := &historyRepoAdapter{storage: taskStore}
	taskUC := taskUsecase.NewTaskUseCase(taskStore, historyRepoAdapter, taskPub)

	activityUC := activityUsecase.NewActivityUseCase(activityStore)

	h := handler.NewHandler(redisCache, userUC, teamUC, taskUC, activityUC, userStore, userStore)

	return &App{
		Config:          cfg,
		Cache:           redisCache,
		Handler:         h,
		UserStorage:     userStore,
		TaskStorage:     taskStore,
		ActivityStorage: activityStore,
		UserPublisher:   userPub,
		TaskPublisher:   taskPub,
	}, nil
}

func (a *App) Close() {
	_ = a.Cache.Close()
	_ = a.UserPublisher.Close()
	_ = a.TaskPublisher.Close()
	a.UserStorage.Close()
	a.TaskStorage.Close()
	a.ActivityStorage.Close()
}

type teamRepoAdapter struct {
	storage *userStorage.Storage
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
	storage *userStorage.Storage
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

type historyRepoAdapter struct {
	storage *taskStorage.Storage
}

func (a *historyRepoAdapter) Create(ctx context.Context, history *domain.TaskHistory) error {
	return a.storage.CreateHistory(ctx, history)
}

func (a *historyRepoAdapter) GetByTaskID(ctx context.Context, taskID string) ([]*domain.TaskHistory, error) {
	return a.storage.GetHistoryByTaskID(ctx, taskID)
}
