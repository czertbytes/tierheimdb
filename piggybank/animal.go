package piggybank

import (
	"fmt"
	"time"
)

func PutAnimals(animals []*Animal) (Ids, error) {
	ids := Ids{}
	for _, a := range uniqueAnimals(animals) {
		if err := PutAnimal(a); err != nil {
			return nil, err
		}

		ids = append(ids, a.Id)
	}

	return ids, nil
}

func PutAnimal(a *Animal) error {
	a.Created = time.Now().Format(time.RFC3339)

	return RedisPersistAnimal(fmt.Sprintf(REDIS_ANIMAL, a.ShelterId, a.UpdateId, a.Id), a)
}

func SearchAnimals(latLon, animalType string, pagination Pagination) (Animals, error) {
	shelters, err := GetShelters(latLon, animalType, maxPagination)
	if err != nil {
		return nil, err
	}

	animals := Animals{}
	for _, s := range shelters {
		u, err := GetLastUpdate(s.Id)
		if err != nil {
			return nil, err
		}

		as, err := GetAnimals(s.Id, u.Id, animalType, maxPagination)
		if err != nil {
			return nil, err
		}

		animals = append(animals, as...)
	}

	return animals.Paginate(pagination), nil
}

func GetAnimals(shelterId, updateId, animalType string, pagination Pagination) (Animals, error) {
	keys, err := RedisGetIndexKeys(fmt.Sprintf(animalsRedisKey(animalType), shelterId, updateId))
	if err != nil {
		return nil, err
	}

	animals, err := RedisGetAnimals(keys)
	if err != nil {
		return nil, err
	}

	return animals.Paginate(pagination), nil
}

func GetAllAnimals(shelterId, updateId string, pagination Pagination) (Animals, error) {
	return GetAnimals(shelterId, updateId, "", pagination)
}

func GetCats(shelterId, updateId string, pagination Pagination) (Animals, error) {
	return GetAnimals(shelterId, updateId, "cat", pagination)
}

func GetDogs(shelterId, updateId string, pagination Pagination) (Animals, error) {
	return GetAnimals(shelterId, updateId, "dog", pagination)
}

func animalsRedisKey(animalType string) string {
	var template string
	switch animalType {
	case "cat":
		template = REDIS_ANIMALS_CATS
	case "dog":
		template = REDIS_ANIMALS_DOGS
	default:
		template = REDIS_ANIMALS
	}

	return template
}

func GetAnimal(shelterId, updateId, id string) (Animal, error) {
	if len(shelterId) == 0 || len(updateId) == 0 || len(id) == 0 {
		return Animal{}, fmt.Errorf("Getting Animal failed! ShelterId, UpdateId, or AnimalId is not set!")
	}

	k := fmt.Sprintf(REDIS_ANIMAL, shelterId, updateId, id)
	animals, err := RedisGetAnimals(Keys{k})
	if err != nil {
		return Animal{}, err
	}

	if len(animals) == 0 {
		return Animal{}, fmt.Errorf("Getting Animal failed! Animal with key '%s' not found!", k)
	}

	return animals[0], nil
}

func DeleteAnimals(shelterId, updateId, animalType string, pagination Pagination) error {
	animals, err := GetAnimals(shelterId, updateId, animalType, pagination)
	if err != nil {
		return err
	}

	for _, a := range animals {
		if err := DeleteAnimal(shelterId, updateId, a.Id); err != nil {
			return err
		}
	}

	return nil
}

func DeleteAnimal(shelterId, updateId, animalId string) error {
	return RedisDeleteAnimal(shelterId, updateId, animalId)
}
