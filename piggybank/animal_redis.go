package piggybank

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
)

func RedisPersistAnimal(k string, a *Animal) error {
	c := RedisPool.Get()
	defer c.Close()

	kExists, err := RedisKeyExists(c, k)
	if err != nil {
		return err
	}
	if kExists != false {
		return fmt.Errorf("Creating Animal with key: %s failed! Animal with that key already exists!", k)
	}

	c.Send("MULTI")
	if err := RedisPersistStruct(c, k, a); err != nil {
		return err
	}

	if err := RedisPersistImages(c, a); err != nil {
		return err
	}

	if err := RedisAddIndexKey(c, fmt.Sprintf(REDIS_ANIMALS, a.ShelterId, a.UpdateId), k); err != nil {
		return err
	}

	switch a.Type {
	case "cat":
		if err := RedisAddIndexKey(c, fmt.Sprintf(REDIS_ANIMALS_CATS, a.ShelterId, a.UpdateId), k); err != nil {
			return err
		}
	case "dog":
		if err := RedisAddIndexKey(c, fmt.Sprintf(REDIS_ANIMALS_DOGS, a.ShelterId, a.UpdateId), k); err != nil {
			return err
		}
	}

	return c.Send("EXEC")
}

func RedisGetAnimals(keys []string) ([]Animal, error) {
	c := RedisPool.Get()
	defer c.Close()

	animals := []Animal{}
	for _, k := range keys {
		kExists, err := RedisKeyExists(c, k)
		if err != nil {
			return nil, err
		}

		if kExists != true {
			return nil, fmt.Errorf("Animal with key '%s' not found!", k)
		}

		v, err := redis.Values(c.Do("HGETALL", k))
		if err != nil {
			return nil, err
		}

		var a Animal
		if err := redis.ScanStruct(v, &a); err != nil {
			return nil, err
		}

		images, err := RedisGetImages(c, a)
		if err != nil {
			return nil, err
		}
		a.Images = images

		animals = append(animals, a)
	}

	return animals, nil
}

func RedisUpdateAnimal(shelterId, updateId, animalId string, a *Animal) error {
	return nil
}

func RedisDeleteAnimal(shelterId, updateId, animalId string) error {
	c := RedisPool.Get()
	defer c.Close()

	k := fmt.Sprintf(REDIS_ANIMAL, shelterId, updateId, animalId)
	kExists, err := RedisKeyExists(c, k)
	if err != nil {
		return err
	}
	if kExists != true {
		return fmt.Errorf("Deleting Shelter with key: %s failed! Shelter with that key does not exists!", k)
	}

	c.Send("MULTI")
	a, err := GetAnimal(shelterId, updateId, animalId)
	if err != nil {
		return err
	}

	if err := RedisDeleteImages(c, &a); err != nil {
		return err
	}

	if err := RedisDeleteIndexKey(c, REDIS_ANIMALS, k); err != nil {
		return err
	}

	switch a.Type {
	case "cat":
		if err := RedisDeleteIndexKey(c, fmt.Sprintf(REDIS_ANIMALS_CATS, a.ShelterId, a.UpdateId), k); err != nil {
			return err
		}
	case "dog":
		if err := RedisDeleteIndexKey(c, fmt.Sprintf(REDIS_ANIMALS_DOGS, a.ShelterId, a.UpdateId), k); err != nil {
			return err
		}
	}

	if err := RedisDeleteKey(c, k); err != nil {
		return err
	}

	return c.Send("EXEC")
}
