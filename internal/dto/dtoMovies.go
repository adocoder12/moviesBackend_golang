package dto

import (
	"time"

	"github.com/adocoder12/moviesBackend_golang/internal/model"
)

type ResponseMovie struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Year      int       `json:"year"`
	Director  string    `json:"director"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RequestMovie struct {
	Title    string `json:"title"`
	Year     int    `json:"year"`
	Director string `json:"director"`
}

type UpdateMovieRequest struct {
	Title    *string `json:"title"`
	Year     *int    `json:"year"`
	Director *string `json:"director"`
}

func (r *RequestMovie) ToModel() *model.Movie {
	return &model.Movie{
		Title:    r.Title,
		Year:     r.Year,
		Director: r.Director,
	}
}

func FromModel(m *model.Movie) ResponseMovie {
	response := ResponseMovie{
		ID:        m.ID,
		Title:     m.Title,
		Year:      m.Year,
		Director:  m.Director,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}

	return response
}
