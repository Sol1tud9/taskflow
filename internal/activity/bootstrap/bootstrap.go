package bootstrap

import (
	"github.com/Sol1tud9/taskflow/internal/activity/consumer"
	"github.com/Sol1tud9/taskflow/internal/activity/storage/sharded"
	"github.com/Sol1tud9/taskflow/internal/activity/usecase"
	"github.com/Sol1tud9/taskflow/pkg/config"
	"github.com/Sol1tud9/taskflow/pkg/logger"
	"go.uber.org/zap"
)

type App struct {
	Config     *config.ActivityServiceConfig
	Storage    *sharded.ShardedStorage
	Consumer   *consumer.EventConsumer
	ActivityUC *usecase.ActivityUseCase
}

func NewApp(cfg *config.ActivityServiceConfig) (*App, error) {
	storage, err := sharded.NewShardedStorage(cfg.Sharding)
	if err != nil {
		logger.Error("failed to init sharded storage", zap.Error(err))
		return nil, err
	}

	activityUC := usecase.NewActivityUseCase(storage)

	groupID := cfg.Kafka.ConsumerGroups["activity_consumer"]
	eventConsumer := consumer.NewEventConsumer(cfg.Kafka.Brokers, cfg.Kafka.Topics, groupID, activityUC)

	return &App{
		Config:     cfg,
		Storage:    storage,
		Consumer:   eventConsumer,
		ActivityUC: activityUC,
	}, nil
}

func (a *App) Close() {
	a.Storage.Close()
	_ = a.Consumer.Close()
}

