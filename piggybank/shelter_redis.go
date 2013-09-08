package piggybank

import (
	"fmt"
	"sort"

	"github.com/garyburd/redigo/redis"
)

func RedisPersistShelter(k string, s *Shelter) error {
	c := RedisPool.Get()
	defer c.Close()

	kExists, err := RedisKeyExists(c, k)
	if err != nil {
		return err
	}
	if kExists != false {
		return fmt.Errorf("Creating Shelter with sKey: %s failed! Shelter with that key already exists!", k)
	}

	c.Send("MULTI")
	if err := RedisPersistStruct(c, k, s); err != nil {
		return err
	}

	if err := RedisAddIndexKey(c, REDIS_SHELTERS, k); err != nil {
		return err
	}
	if s.Enabled {
		if err := RedisAddIndexKey(c, REDIS_SHELTERS_ENABLED, k); err != nil {
			return err
		}
	}

	return c.Send("EXEC")
}

func RedisGetShelters(keys []string) (Shelters, error) {
	c := RedisPool.Get()
	defer c.Close()

	shelters := Shelters{}
	for _, k := range keys {
		kExists, err := RedisKeyExists(c, k)
		if err != nil {
			return nil, err
		}

		if kExists != true {
			return nil, fmt.Errorf("Shelter with key '%s' not found!", k)
		}

		v, err := redis.Values(c.Do("HGETALL", k))
		if err != nil {
			return nil, err
		}

		var s Shelter
		if err := redis.ScanStruct(v, &s); err != nil {
			return nil, err
		}

		shelters = append(shelters, s)
	}
	sort.Sort(ByName{shelters})

	return shelters, nil
}

func RedisUpdateShelter(shelterId string, s *Shelter) error {
	return nil
}

func RedisDeleteShelter(k, id string) error {
	c := RedisPool.Get()
	defer c.Close()

	kExists, err := RedisKeyExists(c, k)
	if err != nil {
		return err
	}
	if kExists != true {
		return fmt.Errorf("Deleting Shelter with key: %s failed! Shelter with that key does not exists!", k)
	}

	c.Send("MULTI")
	s, err := GetShelter(id)
	if err != nil {
		return err
	}

	if err := RedisDeleteIndexKey(c, REDIS_SHELTERS, k); err != nil {
		return err
	}

	if s.Enabled {
		if err := RedisDeleteIndexKey(c, REDIS_SHELTERS_ENABLED, k); err != nil {
			return err
		}
	}

	if err := RedisDeleteKey(c, k); err != nil {
		return err
	}

	return c.Send("EXEC")
}
