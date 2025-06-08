package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"log"
	"micro-service/Internal/api"
	customLogger "micro-service/Internal/logger"
	"micro-service/Internal/repo"
	"micro-service/Internal/service"
	"os"
	"os/signal"
	"syscall"

	"micro-service/Internal/config"
)

func main() {
	if err := godotenv.Load(config.EnvPath); err != nil {
		log.Fatal("Failed to load file", err)
	}
	var cfg config.AppConfig
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(errors.Wrap(err, "Failed to load config"))
	}

	logger, err := customLogger.NewLogger(cfg.LogLevel)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Failed to init logger"))
	}

	repository, err := repo.NewRepo(context.Background(), cfg.PostgreSQL)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Failed to init repository"))
	}

	serviceInstance := service.NewService(repository, logger)

	app := api.NewRouters(&api.Routers{Service: serviceInstance}, cfg.Rest.Token)

	go func() {
		logger.Infof("Starting server on port %s", cfg.Rest.ListenAddress)
		if err := app.Listen(cfg.Rest.ListenAddress); err != nil {
			logger.Fatal(errors.Wrap(err, "Failed to start server"))
		}
	}()
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan

	logger.Info("Shutting down gracefully...")
}
