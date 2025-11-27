package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/DenHax/subscription-manager/docs"
	"github.com/DenHax/subscription-manager/internal/http/handler"
	"github.com/DenHax/subscription-manager/internal/http/server"
	"github.com/DenHax/subscription-manager/internal/logger/slogger"
	"github.com/DenHax/subscription-manager/internal/repo"
	"github.com/DenHax/subscription-manager/internal/service"
	storage "github.com/DenHax/subscription-manager/internal/storage/postgres"
)

// @title API for Subscription aggregation
// @version 1.0
// @description API for Subscription aggregation for Effective Mobile
// @host localhost:8080
// @BasePath /
func main() {
	slogger.InitLogging()

	storageConfig, err := storage.SetupConfig()
	if err != nil {
		slog.Error("fail to get connection url", slog.String("error", err.Error()))
		os.Exit(1)
	}
	storage, err := storage.New(*storageConfig)
	if err != nil {
		slog.Error("failed to init storage", slog.String("error", err.Error()))
		os.Exit(2)
	}
	slog.Debug("db connection", slog.String("connection url", storageConfig.URL))

	repos := repo.NewRepository(storage)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	serverConfig, err := server.SetupConfig()
	if err != nil {
		slog.Error("Failed to setup config", slog.String("error", err.Error()))
		os.Exit(3)
	}
	slog.Info("starting server", slog.String("address", serverConfig.Address))
	srv := server.New(*serverConfig, handlers.Init())

	go func() {
		if err := srv.Run(); err != nil {
			slog.Error("failed to stop server", slog.String("error", err.Error()))
		}
	}()

	slog.Info("server started")

	<-done
	slog.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("failed to stop server", slog.String("error", err.Error()))
		return
	}
}
