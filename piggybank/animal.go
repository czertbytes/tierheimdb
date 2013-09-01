package piggybank

import (
	"fmt"
	"time"
)

type Animal struct {
	Id        string  `json:"id" redis:"id"`
	Created   string  `json:"created" redis:"created"`
	Name      string  `json:"name" redis:"name"`
	URL       string  `json:"url" redis:"url"`
	Priority  int     `json:"priority" redis:"priority"`
	Type      string  `json:"type" redis:"type"`
	Breed     string  `json:"breed" redis:"breed"`
	Sex       string  `json:"sex" redis:"sex"`
	ShortDesc string  `json:"shortDesc" redis:"shortDesc"`
	LongDesc  string  `json:"longDesc" redis:"longDesc"`
	Images    []Image `json:"images" redis:"-"`
	ShelterId string  `json:"shelterId" redis:"shelterId"`
	UpdateId  string  `json:"updateId" redis:"updateId"`
}

type Image struct {
	Width   int    `json:"width" redis:"width"`
	Height  int    `json:"height" redis:"height"`
	URL     string `json:"url" redis:"url"`
	Comment string `json:"comment" redis:"comment"`
}

func PutAnimals(animals []*Animal) ([]string, error) {
	ids := []string{}

	names := make(map[string]bool)
	for _, a := range animals {
		if _, found := names[a.Id]; found == true {
			a.Id = fmt.Sprintf("%s-%d", a.Id, time.Now().UnixNano())
		}

		if err := PutAnimal(a); err != nil {
			return nil, err
		}
		names[a.Id] = true

		ids = append(ids, a.Id)
	}

	return ids, nil
}

func PutAnimal(a *Animal) error {
	a.Created = time.Now().Format(time.RFC3339)

	return RedisPersistAnimal(fmt.Sprintf(REDIS_ANIMAL, a.ShelterId, a.UpdateId, a.Id), a)
}

func GetAnimals(shelterId, updateId string) ([]Animal, error) {
	keys, err := RedisGetIndexKeys(fmt.Sprintf(REDIS_ANIMALS, shelterId, updateId))
	if err != nil {
		return nil, err
	}

	return RedisGetAnimals(keys)
}

func GetAnimal(shelterId, updateId, id string) (Animal, error) {
	if len(shelterId) == 0 || len(updateId) == 0 || len(id) == 0 {
		return Animal{}, fmt.Errorf("Getting Animal failed! ShelterId, UpdateId, or AnimalId is not set!")
	}

	k := fmt.Sprintf(REDIS_ANIMAL, shelterId, updateId, id)
	animals, err := RedisGetAnimals([]string{k})
	if err != nil {
		return Animal{}, err
	}

	if len(animals) == 0 {
		return Animal{}, fmt.Errorf("Getting Animal failed! Animal with key '%s' not found!", k)
	}

	return animals[0], nil
}

func UpdateAnimal() {

}

func DeleteAnimal() {

}
