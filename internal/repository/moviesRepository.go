package repository

import (
	"context"
	"database/sql"

	"github.com/adocoder12/moviesBackend_golang/internal/model"
)

type MoviesRepository struct {
	db *sql.DB
}

func NewMoviesRepository(db *sql.DB) *MoviesRepository {
	return &MoviesRepository{db: db}
}

func (r *MoviesRepository) GetMovies(ctx context.Context) ([]model.Movie, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, title, year FROM movies")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	movies := []model.Movie{}
	for rows.Next() {
		var movie model.Movie
		if err := rows.Scan(&movie.ID, &movie.Title, &movie.Year); err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}
