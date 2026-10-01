package handler

import "net/http"

func (app *Application) SetupRoutes() http.Handler {
	mux := http.NewServeMux()
	//public routes
	mux.HandleFunc("GET /health", app.HealthHandler)
	mux.HandleFunc("GET /api/v1/movies", app.GetAllMoviesHandler)
	mux.HandleFunc("GET /api/v1/movies/{id}", app.GetMovieByIdHandler)
	mux.HandleFunc("POST /api/v1/movies", app.CreateMovieHandler)
	mux.HandleFunc("PUT /api/v1/movies/{id}", app.UpdateMovieHandler)
	mux.HandleFunc("DELETE /api/v1/movies/{id}", app.DeleteMovieHandler)

	return mux
}
