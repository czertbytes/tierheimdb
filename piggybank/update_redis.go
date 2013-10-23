package piggybank

import (
	"fmt"
	"sort"

	"github.com/garyburd/redigo/redis"
)

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

func RedisGetUpdates(keys Keys) (Updates, error) {
	c := RedisPool.Get()
	defer c.Close()

	updates := Updates{}
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
	sort.Sort(ByDate{updates})

	return updates, nil
}

func RedisDeleteUpdate(shelterId, updateId string) error {
	c := RedisPool.Get()
	defer c.Close()

	k := fmt.Sprintf(REDIS_UPDATE, shelterId, updateId)
	kExists, err := RedisKeyExists(c, k)
	if err != nil {
		return err
	}
	if kExists != true {
		return fmt.Errorf("Deleting Update with key: %s failed! Update with that key does not exists!", k)
	}

	c.Send("MULTI")
	u, err := GetUpdate(shelterId, updateId)
	if err != nil {
		return err
	}

	if err := RedisDeleteIndexKey(c, fmt.Sprintf(REDIS_UPDATES, u.ShelterId), k); err != nil {
		return err
	}

	//  be careful! last update is not set! make an update as soon as possible!
	if err := RedisDeleteKey(c, fmt.Sprintf(REDIS_LAST_UPDATE, u.ShelterId)); err != nil {
		return err
	}

	if err := RedisDeleteKey(c, k); err != nil {
		return err
	}

	return c.Send("EXEC")
}

func RedisExistsUpdate(shelterId, updateId string) error {
	c := RedisPool.Get()
	defer c.Close()

	k := fmt.Sprintf(REDIS_UPDATE, shelterId, updateId)
	kExists, err := RedisKeyExists(c, k)
	if err != nil {
		return err
	}
	if kExists != true {
		return fmt.Errorf("Deleting Update with key: %s failed! Update with that key does not exists!", k)
	}

	return nil
}
