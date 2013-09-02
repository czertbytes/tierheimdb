package piggybank

import (
	"fmt"
	"time"
)

type Shelter struct {
	Id           string `json:"id" redis:"id"`
	Created      string `json:"created" redis:"created"`
	Enabled      bool   `json:"enabled" redis:"enabled"`
	Name         string `json:"name" redis:"name"`
	URL          string `json:"url" redis:"url"`
	Phone        string `json:"phone" redis:"phone"`
	Email        string `json:"email" redis:"email"`
	ShortDesc    string `json:"shortDesc" redis:"shortDesc"`
	LongDesc     string `json:"longDesc" redis:"longDesc"`
	Street       string `json:"street" redis:"street"`
	StreetNumber string `json:"streetNumber" redis:"streetNumber"`
	PostalCode   string `json:"postalCode" redis:"postalCode"`
	City         string `json:"city" redis:"city"`
	LatLon       string `json:"latLon" redis:"latLon"`
	Note         string `json:"note" redis:"note"`
}

func PutShelters(shelters []*Shelter) ([]string, error) {
	ids := []string{}
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

func GetAllShelters() ([]Shelter, error) {
	keys, err := RedisGetIndexKeys(REDIS_SHELTERS)
	if err != nil {
		return nil, err
	}

	return RedisGetShelters(keys)
}

func GetEnabledShelters() ([]Shelter, error) {
	keys, err := RedisGetIndexKeys(fmt.Sprintf(REDIS_SHELTERS_ENABLED))
	if err != nil {
		return nil, err
	}

	return RedisGetShelters(keys)
}

func GetShelter(id string) (Shelter, error) {
	if len(id) == 0 {
		return Shelter{}, fmt.Errorf("Getting Shelter failed! ShelterId not set!")
	}

	k := fmt.Sprintf(REDIS_SHELTER, id)
	shelters, err := RedisGetShelters([]string{k})
	if err != nil {
		return Shelter{}, err
	}

	if len(shelters) == 0 {
		return Shelter{}, fmt.Errorf("Getting Shelter failed! ShelterId '%s' not found!", k)
	}

	return shelters[0], nil
}

func DeleteEnabledShelters() error {
	shelters, err := GetEnabledShelters()
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
