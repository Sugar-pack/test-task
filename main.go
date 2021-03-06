package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Sugar-pack/test-task/internal/sender"

	"github.com/Sugar-pack/test-task/internal/api"
	"github.com/Sugar-pack/test-task/internal/config"
	"github.com/Sugar-pack/test-task/internal/logging"
	"github.com/Sugar-pack/test-task/internal/migration"
	"github.com/Sugar-pack/test-task/internal/repository"
)

func main() {
	ctx := context.Background()
	appConfig, err := config.GetAppConfig()
	if err != nil {
		log.Fatal(err)

		return
	}

	logger := logging.GetLogger()
	ctx = logging.WithContext(ctx, logger)
	err = migration.Apply(ctx, appConfig.Db)
	if err != nil {
		log.Fatal(err)

		return
	}

	dbConn, err := migration.Connect(ctx, appConfig.Db)
	if err != nil {
		log.Fatal(err)

		return
	}

	repo := repository.NewPsqlRepository(dbConn)

	rabbitMQ := sender.NewRMQProducer(appConfig.RabbitMQ.Queue, appConfig.RabbitMQ.URL)
	companyHandler := api.NewCompanyHandler(repo, rabbitMQ)
	qualifier := api.NewIPAPI(appConfig.API.Countries)
	router := api.SetupRouter(logger, companyHandler, qualifier)
	server := http.Server{
		Addr:    appConfig.API.Bind,
		Handler: router,
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(shutdown)

	go func() {
		logger.Info("Server is listening on ", appConfig.API.Bind)
		errLaS := server.ListenAndServe()
		if errLaS != nil && errors.Is(errLaS, http.ErrServerClosed) {
			logger.Fatal(errLaS)
		}
	}()

	<-shutdown

	logger.Info("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), appConfig.API.ShutdownTimeout)
	defer func() {
		cancel()
	}()

	if errShutdown := server.Shutdown(ctx); errShutdown != nil {
		logger.WithError(errShutdown).Fatal("Server shutdown error")
	}

	logger.Info("Server stopped gracefully")
}
