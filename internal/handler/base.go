package handler

import (
	"fmt"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hello from home")
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "ok")
}
