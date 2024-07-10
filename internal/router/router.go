package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gomap/internal/handler"
	"net/http"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()

	// Включение middleware логгера
	r.Use(middleware.Logger)

	// Определение маршрутов
	r.Get("/", handler.RootHandler)
	r.Post("/set_gps", handler.SetGpsHandler)
	r.Post("/create_place", handler.CreatePlaceGpsHandler)
	r.Get("/current_map", handler.CurrentMapHandler)
	r.Get("/list_my_gps", handler.ListMyGpsMapHandler)
	r.Get("/near_place", handler.NearPlaceHandler)
	r.Get("/place/{id}", handler.PlaceDetailHandler)
	r.Get("/s_time", handler.ServerTimelHandler)
	r.Get("/login", handler.LoginDetailHandler)
	return r
}
