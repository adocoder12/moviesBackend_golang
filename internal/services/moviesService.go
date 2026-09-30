package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/adocoder12/moviesBackend_golang/internal/model"
)

// MovieStore is what the service needs from storage.
type MovieStore interface {
	GetMovies(ctx context.Context) ([]model.Movie, error)
}

type MoviesService struct {
	repo   MovieStore
	logger *slog.Logger
}

func NewMoviesService(repo MovieStore, logger *slog.Logger) *MoviesService {
	return &MoviesService{repo: repo, logger: logger}
}

func (s *MoviesService) GetMovies(ctx context.Context) ([]model.Movie, error) {
	movies, err := s.repo.GetMovies(ctx)
	if err != nil {
		return nil, fmt.Errorf("movies service: get movies: %w", err)
	}
	return movies, nil
}
