package main

import (
	"context"
	"wb_l0/config"
	"wb_l0/internal/app"
)

func main() {
	ctx := context.Background()
	cfg := config.GetConfigYml()

	a := app.New(ctx, cfg)
	a.Run(cfg)
}
