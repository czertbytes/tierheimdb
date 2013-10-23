package main

import (
	"fmt"
	"net/http"
	"strconv"

	pb "github.com/czertbytes/tierheimdb/piggybank"
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

func validateShelterId(shelterId string) error {
	shelters, err := pb.GetEnabledShelters()
	if err != nil {
		return err
	}

	for _, s := range shelters {
		if s.Id == shelterId {
			return nil
		}
	}

	return fmt.Errorf("ShelterId '%s' is not valid!", shelterId)
}

func validateUpdateId(shelterId, updateId string) error {
	if err := pb.RedisExistsUpdate(shelterId, updateId); err != nil {
		return fmt.Errorf("UpdateId '%s' is not valid!", updateId)
	}

	return nil
}

func validateAnimalId(shelterId, updateId, animalId string) error {
	if err := pb.RedisExistsAnimal(shelterId, updateId, animalId); err != nil {
		return fmt.Errorf("AnimalId '%s' is not valid!", animalId)
	}

	return nil
}
