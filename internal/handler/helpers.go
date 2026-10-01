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

// notFound sends a 404.
func (app *Application) notFound(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

// badRequest logs the problem and sends a 400 with the error message.
// Always pass a non-nil error.
func (app *Application) badRequest(w http.ResponseWriter, r *http.Request, msg string, err error) {
	app.logger.Warn("bad request",
		"msg", msg,
		"err", err,
		"method", r.Method,
		"uri", r.URL.RequestURI(),
	)
	http.Error(w, msg, http.StatusBadRequest)
}
