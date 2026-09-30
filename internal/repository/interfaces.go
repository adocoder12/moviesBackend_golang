package repository

import (
	"context"

	"github.com/adocoder12/moviesBackend_golang/internal/model"
)

type MoviesRepositoryInterfaces interface {
	GetMovies(ctx context.Context) ([]model.Movie, error)
}
