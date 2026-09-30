package handler

import (
	"encoding/json"
	"net/http"
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
