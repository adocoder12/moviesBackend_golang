package handler

import "net/http"

func SetupRoutes() http.Handler {
	mux := http.NewServeMux()
	//public routes
	mux.HandleFunc("GET /{$}", HomeHandler)
	mux.HandleFunc("GET /health", HealthHandler)
	return mux
}
