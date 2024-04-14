package main

import (
	"banner-service/internal/config"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"banner-service/internal/htttp-server/handlers/banner/create"
	"banner-service/internal/htttp-server/handlers/banner/delete"
	bGet "banner-service/internal/htttp-server/handlers/banner/get"
	cGet "banner-service/internal/htttp-server/handlers/content/get"
	aGet "banner-service/internal/htttp-server/handlers/user/get"
	mwLogger "banner-service/internal/htttp-server/middleware/logger"

	"banner-service/internal/htttp-server/middleware/jwt"
	"banner-service/internal/storage/post"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	fmt.Println(cfg)

	log := setupLogger(cfg.Env)

	log.Info(
		"starting banner-service",
		slog.String("env", cfg.Env),
		slog.String("version", "123"),
	)

	log.Debug("debug messages are enabled")

	context := context.Background()
	storage, _ := post.NewPG(context, cfg.StoragePath)
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	//router.Get("/auth", aGet.New(context, log, storage, cfg.SigningKey))

	router.Group(func(r chi.Router) {
		r.Get("/auth", aGet.New(context, log, storage, cfg.SigningKey))
	  })
	router.Group(func(r chi.Router) {
		r.Use(jwt.New(cfg.SigningKey))	

		r.Post("/banner", create.New(context, log, storage))
		r.Get("/user_banner", cGet.New(context, log, storage))
		r.Get("/banner", bGet.New(context, log, storage))
		r.Delete("/banner/{id}", delete.New(context, log, storage))
	  })

	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("panic")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default: // If env config is invalid, set prod settings by default due to security
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

//connString := "postgres://postgres:postgres@localhost:5432/banner-service"
