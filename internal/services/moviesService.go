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
	GetMovieById(ctx context.Context, id int) (*model.Movie, error)
	CreateMovie(ctx context.Context, movie *model.Movie) (*model.Movie, error)
	UpdateMovie(ctx context.Context, movie *model.Movie) (*model.Movie, error)
	DeleteMovie(ctx context.Context, id int) error
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

	response := make([]dto.ResponseMovie, 0, len(movies))
	for i := range movies {
		response = append(response, dto.FromModel(&movies[i]))
	}
	return response, nil
}

func (s *MoviesService) GetMovieById(ctx context.Context, id int) (*dto.ResponseMovie, error) {
	movie, err := s.repo.GetMovieById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("movies service: get movie by id: %w", err)
	}

	response := dto.FromModel(movie)
	return &response, nil
}

func (s *MoviesService) CreateMovie(ctx context.Context, req dto.RequestMovie) (*dto.ResponseMovie, error) {
	newMovie, err := s.repo.CreateMovie(ctx, req.ToModel())
	if err != nil {
		return nil, fmt.Errorf("movies service: create movie: %w", err)
	}

	response := dto.FromModel(newMovie)
	return &response, nil
}

func (s *MoviesService) UpdateMovie(ctx context.Context, id int, req dto.UpdateMovieRequest) (*dto.ResponseMovie, error) {
	movie, err := s.repo.GetMovieById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("movies service: update movie: %w", err)
	}

	if req.Title != nil {
		movie.Title = *req.Title
	}
	if req.Year != nil {
		movie.Year = *req.Year
	}
	if req.Director != nil {
		movie.Director = *req.Director
	}

	updatedMovie, err := s.repo.UpdateMovie(ctx, movie)
	if err != nil {
		return nil, fmt.Errorf("movies service: update movie: %w", err)
	}

	response := dto.FromModel(updatedMovie)
	return &response, nil
}

func (s *MoviesService) DeleteMovie(ctx context.Context, id int) error {
	if err := s.repo.DeleteMovie(ctx, id); err != nil {
		return fmt.Errorf("movies service: delete movie: %w", err)
	}
	return nil
}
