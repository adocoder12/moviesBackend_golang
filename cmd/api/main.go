// wiring only: load env, open DB, start server
package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/adocoder12/moviesBackend_golang/internal/handler"
	"github.com/joho/godotenv"
)

func main() {
	// logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	// load env vars
	if err := godotenv.Load(); err != nil {
		slog.Info("no .env file found, using system environment")
	}

	// development mode
	if os.Getenv("APP_ENV") == "development" {
		slog.Info("development mode enabled")
	}

	// port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handler.SetupRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	slog.Info("server starting", "port", port)

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("server failed", "err", err)
		os.Exit(1)
	}
}
