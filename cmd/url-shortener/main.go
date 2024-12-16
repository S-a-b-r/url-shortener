package main

import (
	"log/slog"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/S-a-b-r/url-shortener/internal/config"
	mwLogger "github.com/S-a-b-r/url-shortener/internal/http-server/middleware"
	"github.com/S-a-b-r/url-shortener/internal/lib/logger/handlers/slogpretty"
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

	db, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("error opening db", sl.Err(err))
		os.Exit(1)
	}

	res, err := db.SaveURL("https://testUrl3rwer3243333", "https://te3st3")
	if err != nil {
		log.Error("error saving url", sl.Err(err))
		os.Exit(1)
	}
	log.Debug("res", slog.Int64("resSaveUrl", res))

	router := chi.NewRouter()

	// middleware

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	// TODO: run server
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envDev:
		log = setupPrettySlog()
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
