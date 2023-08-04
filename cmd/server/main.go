package main

import (
	"context"
	"wb_l0/config"
	"wb_l0/internal/app"
	"wb_l0/pkg/logger"
)

func main() {
	ctx := context.Background()
	cfg := config.GetConfigYml()
	logger := logger.GetLogger()
	a, err := app.New(ctx, cfg)
	if err != nil {
		logger.Err(err)
		logger.Fatal()
	}
	a.Run(cfg)
}
