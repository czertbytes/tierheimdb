package piggybank

import (
	"fmt"
	"time"
)

func PutShelters(shelters []*Shelter) (Ids, error) {
	ids := Ids{}
	for _, s := range shelters {
		if err := PutShelter(s); err != nil {
			return nil, err
		}

		ids = append(ids, s.Id)
	}

	return ids, nil
}

func PutShelter(s *Shelter) error {
	s.Created = time.Now().Format(time.RFC3339)

	return RedisPersistShelter(fmt.Sprintf(REDIS_SHELTER, s.Id), s)
}

func GetAllShelters() (Shelters, error) {
	keys, err := RedisGetIndexKeys(REDIS_SHELTERS)
	if err != nil {
		return nil, err
	}

	return RedisGetShelters(keys)
}

func GetEnabledShelters() (Shelters, error) {
	keys, err := RedisGetIndexKeys(fmt.Sprintf(REDIS_SHELTERS_ENABLED))
	if err != nil {
		return nil, err
	}

	return RedisGetShelters(keys)
}

func sheltersWithAnimalType(shelters Shelters, animalType string) Shelters {
	if len(animalType) > 0 {
		sheltersWithType := Shelters{}
		for _, s := range shelters {
			if s.HasAnimalType(animalType) {
				sheltersWithType = append(sheltersWithType, s)
			}
		}

		shelters = sheltersWithType
	}

	return shelters
}

func sheltersNear(shelters Shelters, latLon string) (Shelters, error) {
	if len(latLon) > 0 {
		lat, lon, err := parseLatLon(latLon)
		if err != nil {
			return nil, err
		}

		sheltersNear := Shelters{}
		for _, s := range shelters {
			sLat, sLon, err := parseLatLon(s.LatLon)
			if err != nil {
				return nil, err
			}

			distance := haversineFormula(lat, lon, sLat, sLon)
			if distance < 50.0 {
				sheltersNear = append(sheltersNear, s)
			}
		}

		shelters = sheltersNear
	}

	return shelters, nil
}

func GetShelters(latLon, animalType string, pagination Pagination) (Shelters, error) {
	shelters, err := GetEnabledShelters()
	if err != nil {
		return nil, err
	}

	shelters = sheltersWithAnimalType(shelters, animalType)
	shelters, err = sheltersNear(shelters, latLon)
	if err != nil {
		return nil, err
	}

	return shelters.Paginate(pagination), nil
}

func GetShelter(id string) (Shelter, error) {
	if len(id) == 0 {
		return Shelter{}, fmt.Errorf("Getting Shelter failed! ShelterId not set!")
	}

	k := fmt.Sprintf(REDIS_SHELTER, id)
	shelters, err := RedisGetShelters(Keys{k})
	if err != nil {
		return Shelter{}, err
	}

	if len(shelters) == 0 {
		return Shelter{}, fmt.Errorf("Getting Shelter failed! ShelterId '%s' not found!", k)
	}

	return shelters[0], nil
}

func DeleteEnabledShelters(latLon, animalType string, pagination Pagination) error {
	shelters, err := GetShelters(latLon, animalType, pagination)
	if err != nil {
		return err
	}

	for _, s := range shelters {
		if err := DeleteShelter(s.Id); err != nil {
			return err
		}
	}

	return nil
}

func DeleteShelter(id string) error {
	return RedisDeleteShelter(fmt.Sprintf(REDIS_SHELTER, id), id)
}
