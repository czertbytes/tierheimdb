package piggybank

import (
	"fmt"
)

type City struct {
	Name     string   `json:name`
	Shelters Shelters `json:shelters`
}

type Cities []City

func GetCities() (Cities, error) {
	shelters, err := GetEnabledShelters()
	if err != nil {
		return nil, err
	}

	return getCities(shelters), nil
}

func GetCity(city string) (City, error) {
	cities, err := GetCities()
	if err != nil {
		return City{}, err
	}

	for _, c := range cities {
		if c.Name == city {
			return c, nil
		}
	}

	return City{}, fmt.Errorf("City not found!")
}

func getCities(shelters Shelters) Cities {
	cityMap := make(map[string]Shelters)
	for _, s := range shelters {
		city := s.City
		if _, found := cityMap[city]; found != true {
			cityMap[city] = Shelters{}
		}

		cityMap[city] = append(cityMap[city], s)
	}

	cities := Cities{}
	for city, shelters := range cityMap {
		cities = append(cities, City{
			Name:     city,
			Shelters: shelters,
		})
	}

	return cities
}
