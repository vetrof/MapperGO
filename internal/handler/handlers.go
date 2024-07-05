package handler

import (
	"fmt"
	"github.com/google/uuid"
	"gomap/internal/db"
	"gomap/internal/gps_utils"
	"gomap/internal/timeutil"
	"net/http"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "RootHandler")
}

func SetGpsHandler(w http.ResponseWriter, r *http.Request) {
	var coords gps_utils.GpsCoordinates
	if err := gps_utils.ParseRequestBody(r, &coords); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Создание новой записи для сохранения координат
	userPlace := db.UserPlace{
		//ID:   generateID(), // Функция для генерации уникального ID
		N:    coords.N,
		E:    coords.E,
		Info: "some info", // Замените это необходимой информацией
	}

	if _, err := db.CreateUserPlace(&userPlace); err != nil {
		http.Error(w, "Failed to save coordinates", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "SetGpsHandler N: %s, E: %s", coords.N, coords.E)
}

func CurrentMapHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func ListMyGpsMapHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "CurrentMapHandler")
}

func NearPlaceHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "NearPlaceHandler")
}

func PlaceDetailHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "PlaceDetailHandler")
}

func ServerTimelHandler(w http.ResponseWriter, r *http.Request) {
	serverTime := timeutil.CurrentTime()
	fmt.Fprintf(w, "ServerTimelHandler %s", serverTime)
}

func LoginDetailHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func generateID() string {
	return uuid.New().String()
}
