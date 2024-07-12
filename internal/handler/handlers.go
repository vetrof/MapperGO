package handler

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"gomap/internal/db"
	"gomap/internal/gps_utils"
	"gomap/internal/timeutil"
	"log"
	"net/http"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "RootHandler")
}

func SetGpsHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var coordinate gps_utils.GpsCoordinates
	err := decoder.Decode(&coordinate)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Println(coordinate)

	// Создаем объект места для сохранения в базу данных
	place := &db.UserPlace{
		Info: "name",
		Geom: "POINT(" + coordinate.Lng + " " + coordinate.Lat + ")",
	}

	fmt.Println(place)

	// Вызываем функцию создания места в базе данных
	_, err = db.CreateUserPlace(place)
	if err != nil {
		log.Println("Error creating place:", err)
		http.Error(w, "Failed to save place", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Place saved successfully"))
}

func CreatePlaceGpsHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var coordinate gps_utils.CoordinateRequest
	err := decoder.Decode(&coordinate)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Println(coordinate)

	// Создаем объект места для сохранения в базу данных
	place := &db.Place{
		Name: coordinate.Name,
		Geom: "POINT(" + coordinate.Lng + " " + coordinate.Lat + ")",
	}

	fmt.Println(place)

	// Вызываем функцию создания места в базе данных
	_, err = db.CreatePlace(place)
	if err != nil {
		log.Println("Error creating place:", err)
		http.Error(w, "Failed to save place", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Place saved successfully"))
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
