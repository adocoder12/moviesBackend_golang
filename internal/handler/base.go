package handler

import (
	"fmt"
	"net/http"
)

func (app *Application) HomeHandler(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Info("home handler")

	fmt.Fprint(w, "hello from home, good ado!!\n")
}

func (app *Application) HealthHandler(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Info("health handler")
	fmt.Fprint(w, "ok\n")
}
