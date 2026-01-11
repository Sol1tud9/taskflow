package main

import (
	"context"
	"os"

	"github.com/Sol1tud9/taskflow/internal/activity/bootstrap"
	"github.com/Sol1tud9/taskflow/pkg/config"
	"github.com/Sol1tud9/taskflow/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "configs/activity.yaml"
	}

	cfg, err := config.Load[config.ActivityServiceConfig](configPath)
	if err != nil {
		panic(err)
	}

	if err := logger.Init(cfg.App.LogLevel); err != nil {
		panic(err)
	}
	defer logger.Sync()

	logger.Info("starting activity-service", zap.String("name", cfg.App.Name))

	app, err := bootstrap.NewApp(cfg)
	if err != nil {
		logger.Fatal("failed to init app", zap.Error(err))
	}
	defer app.Close()

	ctx := context.Background()
	app.Consumer.Start(ctx)

	logger.Info("activity-service started successfully")

	select {}
}

