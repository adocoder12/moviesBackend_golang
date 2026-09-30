package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/adocoder12/moviesBackend_golang/internal/model"
)

type MoviesRepository struct {
	db *sql.DB
}

func NewMoviesRepository(db *sql.DB) *MoviesRepository {
	return &MoviesRepository{db: db}
}

func (r *MoviesRepository) GetMovies(ctx context.Context) ([]model.Movie, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, title, year, director, created_at, updated_at FROM movies")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	movies := []model.Movie{}
	for rows.Next() {
		var movie model.Movie
		if err := rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.Year,
			&movie.Director, // Added missing director column
			&movie.CreatedAt,
			&movie.UpdatedAt,
		); err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}

func (r *MoviesRepository) GetMovieById(ctx context.Context, id int) (*model.Movie, error) {
	query := `
		SELECT id, title, year, director, created_at, updated_at
		FROM movies
		WHERE id = ?`

	row := r.db.QueryRowContext(ctx, query, id)

	var movie model.Movie
	if err := row.Scan(
		&movie.ID,
		&movie.Title,
		&movie.Year,
		&movie.Director,
		&movie.CreatedAt,
		&movie.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("movies repo: get movie by id: %w", err)
	}

	return &movie, nil
}

func (r *MoviesRepository) CreateMovie(ctx context.Context, movie *model.Movie) (*model.Movie, error) {
	query := `
		INSERT INTO movies (title, year, director, created_at, updated_at)
		VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(
		ctx,
		query,
		movie.Title, movie.Year, movie.Director,
	).Scan(&movie.ID, &movie.CreatedAt, &movie.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return movie, nil
}
