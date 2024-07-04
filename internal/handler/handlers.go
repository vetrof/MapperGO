package handler

import (
	"fmt"
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
