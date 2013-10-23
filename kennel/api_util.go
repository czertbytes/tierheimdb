package main

import (
	"fmt"
	"net/http"
	"strconv"

	pb "github.com/czertbytes/tierheimdb/piggybank"
)

func parseAnimalsHandlerQuery(r *http.Request) (string, string, int, int) {
	latlon := validateStringQueryParam("latlon", r)
	animalType := validateStringQueryParam("type", r)
	limit := validateIntQueryParam("limit", r)
	offset := validateIntQueryParam("offset", r)

	return latlon, animalType, limit, offset
}

func validateStringQueryParam(name string, r *http.Request) string {
	param := r.URL.Query().Get(name)
	if len(param) == 0 {
		param = ""
	}

	return param
}

func validateIntQueryParam(name string, r *http.Request) int {
	param := r.URL.Query().Get(name)
	val, err := strconv.Atoi(param)
	if err != nil {
		val = 0
	}

	return val
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
