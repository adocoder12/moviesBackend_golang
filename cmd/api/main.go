// wiring only: load env, open DB, start server
package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/adocoder12/moviesBackend_golang/internal/db"
	"github.com/adocoder12/moviesBackend_golang/internal/handler"
	"github.com/adocoder12/moviesBackend_golang/internal/repository"
	"github.com/adocoder12/moviesBackend_golang/internal/services"
	"github.com/joho/godotenv"
)

func main() {
	// logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	// load env vars
	if err := godotenv.Load(); err != nil {
		logger.Info("no .env file found, using system environment")
	}

	// development mode
	if os.Getenv("APP_ENV") == "development" {
		logger.Info("development mode enabled")
	}

	// port

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "movies.db"
	}

	conn, err := db.Connect(dbPath, logger)
	if err != nil {
		logger.Error("failed to connect to database", "err", err)
		os.Exit(1)
	}
	defer conn.Close()

	repo := repository.NewMoviesRepository(conn)
	moviesService := services.NewMoviesService(repo, logger)
	app := handler.NewApplication(logger, moviesService)
	routes := app.SetupRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      routes,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Info("server running", "port", port)

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("server failed", "err", err)
		os.Exit(1)
	}
}
