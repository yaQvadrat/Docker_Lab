package app

import (
	"app/config"
	httpapi "app/internal/controller/http/v1"
	"app/internal/repo"
	"app/internal/service"
	"app/pkg/httpserver"
	"app/pkg/postgres"
	"app/pkg/validator"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func Run() {
	// Config
	cfg, err := config.New()
	if err != nil {
		log.Fatal(fmt.Errorf("app - config.New: %w", err))
	}

	// Logger
	setLogrus(cfg.Log.Level)
	log.Info("Config read successfully...")

	// Postgres
	log.Info("Initializing postgres...")
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.MaxPoolSize))
	if err != nil {
		log.Fatal(fmt.Errorf("app - postgres.New: %w", err))
	}
	defer pg.Close()

	// Services and repos
	log.Info("Initializing services and repos...")
	services := service.NewServices(service.ServicesDependencies{
		Repos: repo.NewPostgresRepo(pg),
	})

	// Echo handler
	log.Info("Initializing handlers and routes...")
	handler := echo.New()
	httpapi.ConfigureRouter(handler, services)
	handler.Validator = validator.NewCustomValidator()

	// HttpServer
	log.Info("Starting HTTP server...")
	log.Debugf("Server address: %s", cfg.HTTP.Address)
	httpServer := httpserver.New(handler, httpserver.Address(cfg.HTTP.Address))

	// Finish waiting
	log.Info("Configuring graceful shutdown...")
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case s := <-signalChan:
		log.Infof("app - Run - signal: %s", s)
	case err := <-httpServer.Notify():
		log.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Graceful shutdown
	log.Info("Graceful shutdown...")
	if err := httpServer.Shutdown(); err != nil {
		log.Error(fmt.Errorf("app - Run - httpSever.Shutdown: %w", err))
	}
}
