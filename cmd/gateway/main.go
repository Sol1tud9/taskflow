package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Sol1tud9/taskflow/internal/gateway/bootstrap"
	"github.com/Sol1tud9/taskflow/pkg/config"
	"github.com/Sol1tud9/taskflow/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "configs/gateway.yaml"
	}

	cfg, err := config.Load[config.GatewayConfig](configPath)
	if err != nil {
		panic(err)
	}

	if err := logger.Init(cfg.App.LogLevel); err != nil {
		panic(err)
	}
	defer logger.Sync()

	logger.Info("starting api-gateway", zap.String("name", cfg.App.Name))

	app, err := bootstrap.NewApp(cfg)
	if err != nil {
		logger.Fatal("failed to init app", zap.Error(err))
	}
	defer app.Close()

	addr := fmt.Sprintf(":%d", cfg.Server.HTTPPort)
	logger.Info("api-gateway started", zap.String("addr", addr))

	if err := http.ListenAndServe(addr, app.Handler.Router()); err != nil {
		logger.Fatal("failed to start server", zap.Error(err))
	}
}

