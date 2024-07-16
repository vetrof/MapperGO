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
	"strings"
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

func NewPlaceHandler(response http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var newPlaceReq gps_utils.NewPlace
	err := decoder.Decode(&newPlaceReq)
	if err != nil {
		http.Error(response, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Split point into lat and lng
	coords := strings.Split(newPlaceReq.Point, ",")
	if len(coords) != 2 {
		http.Error(response, "Invalid coordinates format", http.StatusBadRequest)
		return
	}
	lat := strings.TrimSpace(coords[0])
	lng := strings.TrimSpace(coords[1])

	fmt.Println("-->>", lat)
	fmt.Println("-->>", lng)

	// Create a new Place object
	place := &db.Place{
		Name: newPlaceReq.Name,
		Geom: "POINT(" + lng + " " + lat + ")",
		Desc: newPlaceReq.Desc,
	}

	// Call function to create the place in the database
	_, err = db.CreatePlace(place)
	if err != nil {
		log.Println("Error creating place:", err)
		http.Error(response, "Failed to save place", http.StatusInternalServerError)
		return
	}

	response.Write([]byte("Place saved successfully"))
}

func CurrentMapHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func ListMyGpsMapHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "CurrentMapHandler")
}

func NearPlaceHandler(response http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var coordinate gps_utils.GpsCoordinates
	err := decoder.Decode(&coordinate)
	if err != nil {
		http.Error(response, "Invalid request body", http.StatusBadRequest)
		return
	}

	places, err := db.GetNearPlaces(coordinate)
	if err != nil {
		http.Error(response, "Error getting near places", http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(places)
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
