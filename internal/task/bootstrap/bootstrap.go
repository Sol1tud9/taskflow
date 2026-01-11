package bootstrap

import (
	"context"

	"github.com/Sol1tud9/taskflow/internal/domain"
	"github.com/Sol1tud9/taskflow/internal/task/publisher"
	"github.com/Sol1tud9/taskflow/internal/task/storage/postgres"
	"github.com/Sol1tud9/taskflow/internal/task/usecase"
	"github.com/Sol1tud9/taskflow/pkg/config"
	"github.com/Sol1tud9/taskflow/pkg/logger"
	"go.uber.org/zap"
)

type App struct {
	Config    *config.TaskServiceConfig
	Storage   *postgres.Storage
	Publisher *publisher.Publisher
	TaskUC    *usecase.TaskUseCase
}

func NewApp(cfg *config.TaskServiceConfig) (*App, error) {
	storage, err := postgres.NewStorage(cfg.Database)
	if err != nil {
		logger.Error("failed to init storage", zap.Error(err))
		return nil, err
	}

	pub := publisher.NewPublisher(cfg.Kafka.Brokers, cfg.Kafka.Topics)

	historyRepoAdapter := &historyRepoAdapter{storage: storage}
	taskUC := usecase.NewTaskUseCase(storage, historyRepoAdapter, pub)

	return &App{
		Config:    cfg,
		Storage:   storage,
		Publisher: pub,
		TaskUC:    taskUC,
	}, nil
}

func (a *App) Close() {
	a.Storage.Close()
	_ = a.Publisher.Close()
}

type historyRepoAdapter struct {
	storage *postgres.Storage
}

func (a *historyRepoAdapter) Create(ctx context.Context, history *domain.TaskHistory) error {
	return a.storage.CreateHistory(ctx, history)
}

func (a *historyRepoAdapter) GetByTaskID(ctx context.Context, taskID string) ([]*domain.TaskHistory, error) {
	return a.storage.GetHistoryByTaskID(ctx, taskID)
}

