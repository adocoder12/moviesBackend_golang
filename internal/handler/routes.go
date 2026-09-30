package handler

import "net/http"

func (app *Application) SetupRoutes() http.Handler {
	mux := http.NewServeMux()
	//public routes
	mux.HandleFunc("GET /{$}", app.HomeHandler)
	mux.HandleFunc("GET /health", app.HealthHandler)
	mux.HandleFunc("GET /movies", app.GetAllMoviesHandler)
	return mux
}
