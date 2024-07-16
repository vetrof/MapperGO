package gps_utils

import (
	"encoding/json"
	"net/http"
)

type GpsCoordinates struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}

// ParseRequestBody
func ParseRequestBody(r *http.Request, v interface{}) error {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	return decoder.Decode(v)
}

type NewPlace struct {
	Point string `json:"point"`
	Name  string `json:"name"`
	Desc  string `json:"desc"`
}
