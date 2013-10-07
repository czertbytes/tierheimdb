package main

import (
	"net/http"
	"strconv"
)

func parseAnimalsHandlerQuery(r *http.Request) (string, string, int, int) {
	latlon := r.URL.Query().Get("latlon")
	if len(latlon) == 0 {
		latlon = ""
	}

	animalType := r.URL.Query().Get("type")
	if len(animalType) == 0 {
		animalType = ""
	}

	l := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(l)
	if err != nil {
		limit = 100
	}

	o := r.URL.Query().Get("offset")
	offset, err := strconv.Atoi(o)
	if err != nil {
		offset = 0
	}

	return latlon, animalType, limit, offset
}
