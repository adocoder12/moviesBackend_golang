package handler

import "net/http"

// serverError logs the real error and sends a generic 500 to the client.
func (app *Application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error("server error",
		"err", err,
		"method", r.Method,
		"uri", r.URL.RequestURI(),
	)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// clientError sends the given status (400, 404, ...) with its standard text.
func (app *Application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// notFound is a shortcut for a 404.
func (app *Application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}
