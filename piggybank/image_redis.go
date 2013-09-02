package piggybank

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

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

func RedisDeleteImages(c redis.Conn, a *Animal) error {
	imagesKey := fmt.Sprintf(REDIS_IMAGES, a.ShelterId, a.UpdateId, a.Id)
	imageKeys, err := RedisGetIndexKeys(imagesKey)
	if err != nil {
		return err
	}

	for _, k := range imageKeys {
		if err := RedisDeleteKey(c, k); err != nil {
			return err
		}
	}

	return RedisDeleteKey(c, imagesKey)
}
