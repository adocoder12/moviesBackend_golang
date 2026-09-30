package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/adocoder12/moviesBackend_golang/internal/dto"
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

func (s *MoviesService) GetMovies(ctx context.Context) ([]dto.ResponseMovie, error) {
	movies, err := s.repo.GetMovies(ctx)
	response := make([]dto.ResponseMovie, 0)
	s.logger.Info("get all movies service", "count", len(movies))

	for _, movie := range movies {
		movieDto := dto.FromModel(&movie)
		response = append(response, movieDto)
	}
	if err != nil {
		return nil, fmt.Errorf("movies service: get movies: %w", err)
	}
	return response, nil
}
