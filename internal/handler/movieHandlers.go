package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/adocoder12/moviesBackend_golang/internal/dto"
	"github.com/adocoder12/moviesBackend_golang/internal/model"
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
		return
	}
	app.logger.Info("get all movies handler", "count", len(movies))
}

func (app *Application) GetMovieByIdHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		app.badRequest(w, r, "invalid movie id", err)
		return
	}

	movie, err := app.movies.GetMovieById(r.Context(), id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			app.notFound(w)
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
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MB

	var req dto.RequestMovie
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		app.badRequest(w, r, "malformed client JSON", err)
		return
	}
	if strings.TrimSpace(req.Title) == "" || req.Year < 1888 {
		app.badRequest(w, r, "title is required and year must be 1888 or later", nil)
		return
	}

	createdMovie, err := app.movies.CreateMovie(r.Context(), req)
	if err != nil {
		if errors.Is(err, model.ErrDuplicate) {
			http.Error(w, "movie already exists", http.StatusConflict)
			return
		}
		app.serverError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createdMovie); err != nil {
		app.logger.Error("encode response", "err", err)
		return
	}
	app.logger.Info("create movie handler", "id", createdMovie.ID)
}

func (app *Application) UpdateMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		app.badRequest(w, r, "invalid movie id", err)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MB

	var req dto.UpdateMovieRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		app.badRequest(w, r, "malformed client JSON", err)
		return
	}

	updatedMovie, err := app.movies.UpdateMovie(r.Context(), id, req)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrNotFound):
			app.notFound(w)
		case errors.Is(err, model.ErrDuplicate):
			http.Error(w, "movie already exists", http.StatusConflict)
		default:
			app.serverError(w, r, err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedMovie); err != nil {
		app.logger.Error("encode response", "err", err)
		return
	}
	app.logger.Info("update movie handler", "id", id)
}

func (app *Application) DeleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		app.badRequest(w, r, "invalid movie id", err)
		return
	}

	if err := app.movies.DeleteMovie(r.Context(), id); err != nil {
		if errors.Is(err, model.ErrNotFound) {
			app.notFound(w)
			return
		}
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	app.logger.Info("delete movie handler", "id", id)
}
