package gps_utils

import (
	"encoding/json"
	"net/http"
)

type GpsCoordinates struct {
	N string `json:"n"`
	E string `json:"e"`
}

// ParseRequestBody
func ParseRequestBody(r *http.Request, v interface{}) error {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	return decoder.Decode(v)
}

type CoordinateRequest struct {
	Name string `json:"name"`
	Lat  string `json:"lat"`
	Lng  string `json:"lng"`
}
