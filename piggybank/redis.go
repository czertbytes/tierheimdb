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
	REDIS_RATE_LIMIT_QUOTA = "tdb:ratelimit:%02d"
	REDIS_RATE_LIMIT_RESET = 60
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

func RedisSetField(key, field, value string) error {
	redisConn := RedisPool.Get()
	defer redisConn.Close()

	if _, err := redisConn.Do("HSET", key, field, value); err != nil {
		log.Fatalf("Setting field value on key %s, field %s failed!\n", key, field)
		return err
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

func RedisDeleteKey(c redis.Conn, k string) error {
	if _, err := c.Do("DEL", k); err != nil {
		return fmt.Errorf("Deleting key %s failed!", k)
	}

	return nil
}

func RedisAddIndexKey(c redis.Conn, indexKey, key string) error {
	if _, err := c.Do("SADD", indexKey, key); err != nil {
		return fmt.Errorf("Adding key %s to index :%s failed!\n", key, indexKey)
	}

	return nil
}

func RedisGetIndexKeys(k string) (Keys, error) {
	c := RedisPool.Get()
	defer c.Close()

	keysReply, err := redis.Values(c.Do("SMEMBERS", k))
	if err != nil {
		return nil, fmt.Errorf("Getting keys from indexKey %s failed!\n", k)
	}

	keys := Keys{}
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

func RedisUpdateQuota(minute int, ipAddress string) int {
	redisConn := RedisPool.Get()
	defer redisConn.Close()

	key := fmt.Sprintf(REDIS_RATE_LIMIT_QUOTA, minute)
	currentQuota, err := redis.Int(redisConn.Do("HINCRBY", key, ipAddress, "1"))
	if err != nil {
		return -1
	}

	if _, err := redisConn.Do("EXPIRE", key, REDIS_RATE_LIMIT_RESET); err != nil {
		return -1
	}

	return currentQuota
}
