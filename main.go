package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/S-a-b-r/url-shortener/internal/config"
	"github.com/S-a-b-r/url-shortener/internal/http-server/handlers/redirect"
	"github.com/S-a-b-r/url-shortener/internal/http-server/handlers/save"
	mwLogger "github.com/S-a-b-r/url-shortener/internal/http-server/middleware"
	"github.com/S-a-b-r/url-shortener/internal/lib/logger/handlers/slogpretty"
	"github.com/S-a-b-r/url-shortener/internal/lib/logger/sl"
	"github.com/S-a-b-r/url-shortener/internal/storage/sqlite"
)

const (
	envDev  = "dev"
	envProd = "prod"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@test.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @servers.url https://petstore.example.com/api/v3

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)

	log.Info("starting url shortener", slog.String("env", cfg.Env))

	db, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("error opening db", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()

	// middleware

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	// @tag.name Work with url

	// @success      200              {string}  string    "ok"
	// @router /url [post]
	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			cfg.HTTPServer.Username: cfg.HTTPServer.Password,
		}))

		r.Post("/", save.New(log, db))
	})

	router.Get("{alias}", redirect.New(log, db))

	log.Info("starting server", slog.String("address", cfg.HTTPServer.Address))

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err = srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("oops, something wrong")
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
