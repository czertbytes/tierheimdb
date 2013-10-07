package piggybank

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func intAnimalTypes(animalTypes []string) int {
	val := 0
	for _, animalType := range animalTypes {
		switch animalType {
		case "cat":
			val |= (1 << 0)
		case "dog":
			val |= (1 << 1)
		}
	}

	return val
}

func animalTypes(intAnimalTypes int) []string {
	animalTypes := []string{}
	if intAnimalTypes&(1<<0) != 0 {
		animalTypes = append(animalTypes, "cat")
	}

	if intAnimalTypes&(1<<1) != 0 {
		animalTypes = append(animalTypes, "dog")
	}

	return animalTypes
}

func parseLatLon(location string) (float64, float64, error) {
	splitted := strings.Split(location, ",")
	if len(splitted) != 2 {
		return 0, 0, fmt.Errorf("Parsing location '%s' failed!", location)
	}

	lat, err := strconv.ParseFloat(splitted[0], 10)
	if err != nil {
		return 0, 0, fmt.Errorf("Parsing location lat '%s' failed!", splitted[0])
	}

	lon, err := strconv.ParseFloat(splitted[1], 10)
	if err != nil {
		return 0, 0, fmt.Errorf("Parsing location lon '%s' failed!", splitted[1])
	}

	return lat, lon, nil
}

func haversineFormula(lat1, lon1, lat2, lon2 float64) float64 {
	R := 6371.0
	dLat := (lat2 - lat1) * (math.Pi / 180.0)
	dLon := (lon2 - lon1) * (math.Pi / 180.0)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180.0)*math.Cos(lat2*math.Pi/180.0)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Sin(math.Sqrt(a))

	return R * c
}
