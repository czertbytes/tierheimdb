package piggybank

import (
	"fmt"
	"time"
)

type Shelter struct {
	Id             string   `json:"id" redis:"id"`
	Created        string   `json:"created" redis:"created"`
	Enabled        bool     `json:"enabled" redis:"enabled"`
	Name           string   `json:"name" redis:"name"`
	FullName       string   `json:"fullName" redis:"fullName"`
	URL            string   `json:"url" redis:"url"`
	LogoURL        string   `json:"logoUrl" redis:"logoUrl"`
	Phone          string   `json:"phone" redis:"phone"`
	Email          string   `json:"email" redis:"email"`
	ShortDesc      string   `json:"shortDesc" redis:"shortDesc"`
	LongDesc       string   `json:"longDesc" redis:"longDesc"`
	Street         string   `json:"street" redis:"street"`
	StreetNumber   string   `json:"streetNumber" redis:"streetNumber"`
	PostalCode     string   `json:"postalCode" redis:"postalCode"`
	City           string   `json:"city" redis:"city"`
	LatLon         string   `json:"latLon" redis:"latLon"`
	Note           string   `json:"note" redis:"note"`
	AnimalTypes    []string `json:"animalTypes" redis:"-"`
	IntAnimalTypes int      `json:"-" redis:"animalTypes"`
}

func (s *Shelter) SetIntAnimalTypes() {
	s.IntAnimalTypes = intAnimalTypes(s.AnimalTypes)
}

func (s *Shelter) SetAnimalTypes() {
	s.AnimalTypes = animalTypes(s.IntAnimalTypes)
}

func (s *Shelter) HasAnimalType(animalType string) bool {
	for _, aType := range s.AnimalTypes {
		if aType == animalType {
			return true
		}
	}

	return false
}

type Shelters []Shelter

func (ss Shelters) Len() int {
	return len(ss)
}

func (ss Shelters) Swap(i, j int) {
	ss[i], ss[j] = ss[j], ss[i]
}

type ByName struct {
	Shelters
}

func (s ByName) Less(i, j int) bool {
	return s.Shelters[i].Name < s.Shelters[j].Name
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

func GetSheltersNear(latLon string, maxDistance float64) (Shelters, error) {
	lat, lon, err := parseLatLon(latLon)
	if err != nil {
		return nil, err
	}

	ss, err := GetEnabledShelters()
	if err != nil {
		return nil, err
	}

	shelters := Shelters{}
	for _, s := range ss {
		sLat, sLon, err := parseLatLon(s.LatLon)
		if err != nil {
			return nil, err
		}

		distance := haversineFormula(lat, lon, sLat, sLon)
		if distance < maxDistance {
			shelters = append(shelters, s)
		}
	}

	return shelters, nil
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
