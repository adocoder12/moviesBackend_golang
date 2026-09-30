package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/adocoder12/moviesBackend_golang/internal/dto"
	"github.com/adocoder12/moviesBackend_golang/internal/model"
)

// MovieStore is what the service needs from storage.
type MovieStore interface {
	GetMovies(ctx context.Context) ([]model.Movie, error)
	GetMovieById(ctx context.Context, id int) (*model.Movie, error)
	CreateMovie(ctx context.Context, movie *model.Movie) (*model.Movie, error)
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
	if err != nil {
		return nil, fmt.Errorf("movies service: get movies: %w", err)
	}

	s.logger.Info("get all movies service", "count", len(movies))

	response := make([]dto.ResponseMovie, 0, len(movies)) // Pre-allocate capacity for better performance
	for _, movie := range movies {
		movieDto := dto.FromModel(&movie)
		response = append(response, movieDto)
	}

	return response, nil
}

var ErrMovieNotFound = errors.New("movie not found")

func (s *MoviesService) GetMovieById(ctx context.Context, id int) (*dto.ResponseMovie, error) {
	movie, err := s.repo.GetMovieById(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrMovieNotFound
		}
		return nil, fmt.Errorf("movies service: get movie by id: %w", err)
	}
	response := dto.FromModel(movie)
	return &response, nil
}

func (s *MoviesService) CreateMovie(ctx context.Context, req dto.RequestMovie) (*dto.ResponseMovie, error) {
	movie := req.ToModel()
	newMovie, err := s.repo.CreateMovie(ctx, movie)
	if err != nil {
		return nil, fmt.Errorf("movies service: create movie: %w", err)
	}
	response := dto.FromModel(newMovie)
	return &response, nil
}
