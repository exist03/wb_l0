package app

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"wb_l0/config"
	"wb_l0/internal/handlers"
	"wb_l0/internal/repository"
	"wb_l0/internal/service"
	"wb_l0/pkg/logger"
)

type App struct {
	handlers   *handlers.Handlers
	service    *service.Service
	repository *repository.Repository
	router     *fiber.App
}

func New(ctx context.Context, cfg *config.Config) (*App, error) {
	logger := logger.GetLogger()
	a := &App{}
	a.repository = repository.New(ctx, cfg.PsqlStorage)
	a.service = service.New(a.repository)
	a.handlers = handlers.New(a.service)
	err := a.repository.CacheRecovery()
	if err != nil {
		logger.Err(err).Msg("Problems with recovery")
		return nil, err
	}
	a.router = fiber.New()
	a.router.Get("/service/get/:id", a.handlers.Get)
	return a, nil
}

func (a *App) Run(cfg *config.Config) {
	logger := logger.GetLogger()
	go func() {
		err := a.service.ConsumeMessages()
		if err != nil {
			logger.Err(err)
		}
	}()

	//Graceful	Shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig
		log.Println("Gracefully shutdown")
		service.ShutdownStream()
		if err := a.router.ShutdownWithTimeout(30 * time.Second); err != nil {
			log.Fatalln("server shutdown error: ", err)
		}
	}()

	err := a.router.Listen(cfg.Listen.Port)
	if err != nil {
		log.Fatalln(err)
	}
}
