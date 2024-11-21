package main

import (
	"log/slog"
	"os"

	"github.com/S-a-b-r/url-shortener/internal/config"
	"github.com/S-a-b-r/url-shortener/internal/lib/logger/sl"
	"github.com/S-a-b-r/url-shortener/internal/storage/sqlite"
)

const (
	envDev  = "dev"
	envProd = "prod"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)

	log.Info("starting url shortener", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	db, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("error opening db", sl.Err(err))
		os.Exit(1)
	}
	_ = db

	// TODO: init config: cleanenv
	// TODO: init logger: slog
	// TODO: init storage: sqlight
	// TODO: init router : chi
	// TODO: run server
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envDev:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
