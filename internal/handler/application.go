package handler

import (
	"log/slog"

	"github.com/adocoder12/moviesBackend_golang/internal/services"
)

// Application holds the shared dependencies of the handlers.
type Application struct {
	logger *slog.Logger
	movies *services.MoviesService
}

func NewApplication(logger *slog.Logger, movies *services.MoviesService) *Application {
	return &Application{logger: logger, movies: movies}
}
