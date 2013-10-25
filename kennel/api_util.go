package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	pb "github.com/czertbytes/tierheimdb/piggybank"
)

func parseTypedParams(r *http.Request) string {
	animalType := validateStringQueryParam("type", r)

	for _, t := range []string{"cat", "dog"} {
		if animalType == t {
			return animalType
		}
	}

	return ""
}

func parseLatLonedParams(r *http.Request) string {
	return validateStringQueryParam("latlon", r)
}

func parsePaginationParams(r *http.Request) pb.Pagination {
	offset := validateIntQueryParam("offset", r)
	limit := validateIntQueryParam("limit", r)

	if offset < 0 {
		offset = 0
	}

	if limit == 0 {
		limit = 100
	}
	if limit > 300 {
		limit = 300
	}

	return pb.Pagination{
		offset,
		limit,
	}
}

func validateStringQueryParam(name string, r *http.Request) string {
	param := r.URL.Query().Get(name)
	if len(param) == 0 {
		param = ""
	}

	return strings.ToLower(param)
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
