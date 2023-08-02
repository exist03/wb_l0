package app

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"log"
	"wb_l0/config"
	"wb_l0/internal/handlers"
	"wb_l0/internal/repository"
	"wb_l0/internal/service"
)

type App struct {
	handlers   *handlers.Handlers
	service    *service.Service
	repository *repository.Repository
	fiber      *fiber.App
}

func New(ctx context.Context, cfg *config.Config) *App {
	a := &App{}
	a.repository = repository.New(ctx, cfg.PsqlStorage)
	a.service = service.New(a.repository)
	a.handlers = handlers.New(a.service)
	err := a.repository.CacheRecovery()
	if err != nil {
		return nil
	}
	a.fiber = fiber.New()
	a.fiber.Get("/service/get/:id", a.handlers.Get)
	return a
}

func (a *App) Run(cfg *config.Config) {
	err := a.fiber.Listen(cfg.Listen.Port)
	if err != nil {
		log.Fatalln(err)
	}
}
