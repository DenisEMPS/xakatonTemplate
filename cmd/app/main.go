package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"xaxaton/config"
	"xaxaton/internal/handler"
	"xaxaton/internal/repository"
	"xaxaton/internal/repository/postgres"
	"xaxaton/internal/server"
	"xaxaton/internal/usecase"
	"xaxaton/pkg/logger"
)

func main() {
	cfg := config.MustLoad()
	log := logger.MustLoad(os.Stdout, cfg.Logger.Level)

	log.Info("starting application")

	pgPool, err := postgres.NewPool(cfg.Postgres)
	if err != nil {
		log.Error("failed to init pool connection", slog.Any("error", err))
		os.Exit(1)
	}
	defer pgPool.Close()

	repo := repository.New(pgPool)
	usecase := usecase.New(repo, log)
	handler := handler.New(usecase)

	srv := server.NewServer()
	srv.InitRoutes(handler)

	go func() {
		if err := srv.Run(cfg.Server.Addr); err != nil {
			log.Error("failed to start server", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	log.Info("application successfuly started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	if err := srv.Shutdown(); err != nil {
		log.Error("failed to stop server", slog.Any("error", err))
	}
}
