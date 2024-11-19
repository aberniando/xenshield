package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"github.com/aberniando/xenshield/config"
	"github.com/aberniando/xenshield/pkg/httpserver"
	loggerPkg "github.com/aberniando/xenshield/pkg/logger"
	"github.com/aberniando/xenshield/pkg/postgres"
)

func Run() {
	logger := loggerPkg.GetLogger()

	cfg, err := config.GetConfig()
	if err != nil {
		logger.Fatal("Config error: %s", err)
	}

	// Repository
	db, err := postgres.New(cfg.PG)
	if err != nil {
		logger.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer db.Close()

	repositories := InitRepositories(db)
	services := InitServices(repositories)
	handlers := InitHandlers(services, logger)

	// HTTP Server
	handler := gin.New()
	InitRouter(handler, handlers)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logger.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		logger.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		logger.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
