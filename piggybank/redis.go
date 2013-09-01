package piggybank

import (
	"fmt"
	"log"
	"time"

	"github.com/garyburd/redigo/redis"
)

const (
	REDIS_SHELTER          = "tdb:s:%s"
	REDIS_SHELTERS         = "tdb:shelters"
	REDIS_SHELTERS_ENABLED = "tdb:shelters:enabled"
	REDIS_UPDATE           = "tdb:s:%s:u:%s"
	REDIS_UPDATES          = "tdb:s:%s:updates"
	REDIS_LAST_UPDATE      = "tdb:s:%s:last-update"
	REDIS_ANIMAL           = "tdb:s:%s:u:%s:a:%s"
	REDIS_ANIMALS          = "tdb:s:%s:u:%s:animals"
	REDIS_ANIMALS_CATS     = "tdb:s:%s:u:%s:animals:cats"
	REDIS_ANIMALS_DOGS     = "tdb:s:%s:u:%s:animals:dogs"
	REDIS_IMAGE            = "tdb:s:%s:u:%s:a:%s:i:%d"
	REDIS_IMAGES           = "tdb:s:%s:u:%s:a:%s:images"
)

var (
	RedisPool *redis.Pool
)

func RedisInit() error {
	RedisPool = &redis.Pool{
		MaxIdle:     5,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				return nil, err
			}

			return c, err
		},
	}

	return nil
}

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

func RedisGetShelters(keys []string) ([]Shelter, error) {
	c := RedisPool.Get()
	defer c.Close()

	shelters := []Shelter{}
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

	return shelters, nil
}

func RedisUpdateShelter() {

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

	return c.Send("EXEC")
}

func RedisPersistUpdate(k string, u *Update) error {
	c := RedisPool.Get()
	defer c.Close()

	kExists, err := RedisKeyExists(c, k)
	if err != nil {
		return err
	}
	if kExists != false {
		return fmt.Errorf("Creating Update with key: %s failed! Update with that key already exists!", k)
	}

	c.Send("MULTI")
	if err := RedisPersistStruct(c, k, u); err != nil {
		return err
	}

	if err := RedisAddIndexKey(c, fmt.Sprintf(REDIS_UPDATES, u.ShelterId), k); err != nil {
		return err
	}

	if err := RedisPersistValue(c, fmt.Sprintf(REDIS_LAST_UPDATE, u.ShelterId), k); err != nil {
		return err
	}

	return c.Send("EXEC")
}

func RedisGetUpdates(keys []string) ([]Update, error) {
	c := RedisPool.Get()
	defer c.Close()

	updates := []Update{}
	for _, k := range keys {
		kExists, err := RedisKeyExists(c, k)
		if err != nil {
			return nil, err
		}

		if kExists != true {
			return nil, fmt.Errorf("Update with key '%s' not found!", k)
		}

		v, err := redis.Values(c.Do("HGETALL", k))
		if err != nil {
			return nil, err
		}

		var u Update
		if err := redis.ScanStruct(v, &u); err != nil {
			return nil, err
		}

		updates = append(updates, u)
	}

	return updates, nil
}

func RedisPersistAnimal(k string, a *Animal) error {
	c := RedisPool.Get()
	defer c.Close()

	kExists, err := RedisKeyExists(c, k)
	if err != nil {
		return err
	}
	if kExists != false {
		return fmt.Errorf("Creating Update with key: %s failed! Update with that key already exists!", k)
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

func RedisPersistImages(c redis.Conn, a *Animal) error {
	if len(a.Images) > 0 {
		for _, i := range a.Images {
			k := fmt.Sprintf(REDIS_IMAGE, a.ShelterId, a.UpdateId, a.Id, time.Now().UnixNano())
			if err := RedisPersistStruct(c, k, i); err != nil {
				return err
			}

			if err := RedisAddIndexKey(c, fmt.Sprintf(REDIS_IMAGES, a.ShelterId, a.UpdateId, a.Id), k); err != nil {
				return err
			}
		}
	}

	return nil
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

func RedisGetImages(c redis.Conn, a Animal) ([]Image, error) {
	imagesKey := fmt.Sprintf(REDIS_IMAGES, a.ShelterId, a.UpdateId, a.Id)
	kExists, err := RedisKeyExists(c, imagesKey)
	if err != nil {
		return nil, err
	}

	images := []Image{}
	if kExists {
		imageKeys, err := RedisGetIndexKeys(imagesKey)
		if err != nil {
			return nil, err
		}

		for _, k := range imageKeys {
			kExists, err := RedisKeyExists(c, k)
			if err != nil {
				return nil, err
			}

			if kExists {
				v, err := redis.Values(c.Do("HGETALL", k))
				if err != nil {
					return nil, err
				}

				var i Image
				if err := redis.ScanStruct(v, &i); err != nil {
					return nil, err
				}

				images = append(images, i)
			}
		}
	}

	return images, nil
}

func RedisKeyExists(c redis.Conn, key string) (bool, error) {
	keyExistsReply, err := redis.Int(c.Do("EXISTS", key))
	if err != nil {
		log.Fatalf("Checking existence of key: %s failed!\n", key)
		return false, err
	}

	return keyExistsReply != 0, nil
}

func RedisPersistStruct(c redis.Conn, k string, i interface{}) error {
	if _, err := c.Do("HMSET", redis.Args{}.Add(k).AddFlat(i)...); err != nil {
		return fmt.Errorf("Persisting Struct with key: %s failed!", k)
	}

	return nil
}

func RedisPersistValue(c redis.Conn, k string, v string) error {
	if _, err := c.Do("SET", k, v); err != nil {
		return fmt.Errorf("Persisting Value with key: %s failed!", k)
	}

	return nil
}

func RedisGetValue(k string) (string, error) {
	c := RedisPool.Get()
	defer c.Close()

	return redis.String(c.Do("GET", k))
}

func RedisAddIndexKey(c redis.Conn, indexKey, key string) error {
	if _, err := c.Do("SADD", indexKey, key); err != nil {
		return fmt.Errorf("Adding key %s to index :%s failed!\n", key, indexKey)
	}

	return nil
}

func RedisGetIndexKeys(k string) ([]string, error) {
	c := RedisPool.Get()
	defer c.Close()

	keysReply, err := redis.Values(c.Do("SMEMBERS", k))
	if err != nil {
		return nil, fmt.Errorf("Getting keys from indexKey %s failed!\n", k)
	}

	keys := []string{}
	for len(keysReply) > 0 {
		var key string
		keysReply, err = redis.Scan(keysReply, &key)
		if err != nil {
			return nil, fmt.Errorf("Getting key from all shelters keys request failed!")
		}

		keys = append(keys, key)
	}

	return keys, nil
}

func RedisDeleteIndexKey(c redis.Conn, k, indexKey string) error {
	if _, err := c.Do("SREM", indexKey, k); err != nil {
		return fmt.Errorf("Removing key %s from index: %s failed!\n", k, indexKey)
	}

	return nil
}

/*
func RedisDeleteKey(key string) error {
	redisConn := RedisPool.Get()
	defer redisConn.Close()

	if _, err := redisConn.Do("DEL", key); err != nil {
		log.Fatalf("Deleting key %s failed!\n", key)
		return err
	}

	return nil
}

func RedisSyncLock(key string) error {
	redisConn := RedisPool.Get()
	defer redisConn.Close()

	if _, err := redisConn.Do("SET", key, "1"); err != nil {
		log.Fatalf("Setting sync lock on key %s failed!\n", key)
		return err
	}

	return nil
}

func RedisSetField(key, field, value string) error {
	redisConn := RedisPool.Get()
	defer redisConn.Close()

	if _, err := redisConn.Do("HSET", key, field, value); err != nil {
		log.Fatalf("Setting field value on key %s, field %s failed!\n", key, field)
		return err
	}

	return nil
}

func RedisUpdateQuota(hour int, ipAddress string) int {
	redisConn := RedisPool.Get()
	defer redisConn.Close()

	key := fmt.Sprintf(RATE_LIMIT_QUOTA_KEY, hour)
	currentQuota, err := redis.Int(redisConn.Do("HINCRBY", key, ipAddress, "1"))
	if err != nil {
		return -1
	}

	if _, err := redisConn.Do("EXPIRE", key, 3600); err != nil {
		return -1
	}

	return currentQuota
}
*/
