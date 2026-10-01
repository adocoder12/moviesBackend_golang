package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/adocoder12/moviesBackend_golang/internal/model"
	"github.com/mattn/go-sqlite3"
)

const movieColumns = "id, title, year, director, created_at, updated_at"

type MoviesRepository struct {
	db *sql.DB
}

func NewMoviesRepository(db *sql.DB) *MoviesRepository {
	return &MoviesRepository{db: db}
}

// scanner lets scanMovie work with both *sql.Row and *sql.Rows.
type scanner interface {
	Scan(dest ...any) error
}

func scanMovie(row scanner) (model.Movie, error) {
	var movie model.Movie
	err := row.Scan(
		&movie.ID,
		&movie.Title,
		&movie.Year,
		&movie.Director,
		&movie.CreatedAt,
		&movie.UpdatedAt,
	)
	return movie, err
}

// isUniqueViolation reports whether err is a SQLite UNIQUE constraint failure.
func isUniqueViolation(err error) bool {
	var sqliteErr sqlite3.Error
	return errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique
}

func (r *MoviesRepository) GetMovies(ctx context.Context) ([]model.Movie, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT "+movieColumns+" FROM movies")
	if err != nil {
		return nil, fmt.Errorf("movies repo: get movies: %w", err)
	}
	defer func() { _ = rows.Close() }()

	movies := []model.Movie{}
	for rows.Next() {
		movie, err := scanMovie(rows)
		if err != nil {
			return nil, fmt.Errorf("movies repo: scan movie: %w", err)
		}
		movies = append(movies, movie)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("movies repo: iterate movies: %w", err)
	}

	return movies, nil
}

func (r *MoviesRepository) GetMovieById(ctx context.Context, id int) (*model.Movie, error) {
	row := r.db.QueryRowContext(ctx, "SELECT "+movieColumns+" FROM movies WHERE id = ?", id)

	movie, err := scanMovie(row)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("movies repo: get movie by id: %w", err)
	}

	return &movie, nil
}

func (r *MoviesRepository) CreateMovie(ctx context.Context, movie *model.Movie) (*model.Movie, error) {
	query := `
		INSERT INTO movies (title, year, director, created_at, updated_at)
		VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query, movie.Title, movie.Year, movie.Director).
		Scan(&movie.ID, &movie.CreatedAt, &movie.UpdatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, model.ErrDuplicate
		}
		return nil, fmt.Errorf("movies repo: create movie: %w", err)
	}

	return movie, nil
}

func (r *MoviesRepository) UpdateMovie(ctx context.Context, movie *model.Movie) (*model.Movie, error) {
	query := `
		UPDATE movies
		SET title = ?, year = ?, director = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
		RETURNING updated_at`

	err := r.db.QueryRowContext(ctx, query, movie.Title, movie.Year, movie.Director, movie.ID).
		Scan(&movie.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, model.ErrNotFound
	}
	if err != nil {
		if isUniqueViolation(err) {
			return nil, model.ErrDuplicate
		}
		return nil, fmt.Errorf("movies repo: update movie: %w", err)
	}

	return movie, nil
}

func (r *MoviesRepository) DeleteMovie(ctx context.Context, id int) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM movies WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("movies repo: delete movie: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("movies repo: delete movie rows affected: %w", err)
	}
	if affected == 0 {
		return model.ErrNotFound
	}

	return nil
}
