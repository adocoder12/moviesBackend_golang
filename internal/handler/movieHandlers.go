package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/adocoder12/moviesBackend_golang/internal/dto"
	"github.com/adocoder12/moviesBackend_golang/internal/services" // Import services package for ErrMovieNotFound
)

func (app *Application) GetAllMoviesHandler(w http.ResponseWriter, r *http.Request) {
	movies, err := app.movies.GetMovies(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(movies); err != nil {
		app.logger.Error("encode response", "err", err)
	}
	app.logger.Info("get all movies handler", "count", len(movies))
}

func (app *Application) GetMovieByIdHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		app.badRequest(w, r, err) // Use 400 Bad Request for malformed IDs
		return
	}

	movie, err := app.movies.GetMovieById(r.Context(), id)
	if err != nil {
		// Check for service-level domain error instead of sql.ErrNoRows directly
		if errors.Is(err, services.ErrMovieNotFound) {
			app.notFound(w) // Use 404 Not Found when the movie doesn't exist
			return
		}
		app.serverError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(movie); err != nil {
		app.logger.Error("encode response", "err", err)
		return
	}

	app.logger.Info("get movie by id handler", "id", id)
}

func (app *Application) CreateMovieHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.RequestMovie
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		app.badRequest(w, r, err) // Fixed: use badRequest (400) for malformed client JSON
		return
	}

	createdMovie, err := app.movies.CreateMovie(r.Context(), req)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Set headers FIRST, then write status code and body
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(createdMovie); err != nil {
		app.logger.Error("encode response", "err", err)
		return
	}

	app.logger.Info("create movie handler", "id", createdMovie.ID)
}
