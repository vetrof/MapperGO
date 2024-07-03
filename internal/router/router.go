package router

import (
	"framrless/internal/handler"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()

	// Включение middleware логгера
	r.Use(middleware.Logger)

	// Определение маршрутов
	r.Get("/", handler.RootHandler)
	r.Get("/greet/{id}", handler.GreetHandler)
	r.Get("/time", handler.TimeHandler)

	return r
}
