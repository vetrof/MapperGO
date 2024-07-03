package handler

import (
	"fmt"
	"framrless/internal/timeutil"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func GreetHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	fmt.Fprintf(w, "Greetings, ID: %s", id)
}

func TimeHandler(w http.ResponseWriter, r *http.Request) {
	currentTime := timeutil.CurrentTime()
	fmt.Fprintf(w, "Server time: %s", currentTime)
}

func CoordHandler(w http.ResponseWriter, r *http.Request) {
	coordinates := chi.URLParam(r, "coordinates")
	fmt.Fprintf(w, "coordinates: %s", coordinates)
}
