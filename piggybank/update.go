package piggybank

import (
	"fmt"
	"time"

	"github.com/nu7hatch/gouuid"
)

func NewUpdate(shelterId string) *Update {
	return &Update{
		Id:        makeUpdateId(),
		ShelterId: shelterId,
	}
}

func makeUpdateId() string {
	u4, err := uuid.NewV4()
	if err != nil {
		return fmt.Sprintf("update-%d", time.Now().UnixNano())
	}
	return u4.String()
}

func PutUpdate(u *Update) error {
	u.Created = time.Now().Format(time.RFC3339)

	return RedisPersistUpdate(fmt.Sprintf(REDIS_UPDATE, u.ShelterId, u.Id), u)
}

func GetUpdates(shelterId string, pagination Pagination) (Updates, error) {
	keys, err := RedisGetIndexKeys(fmt.Sprintf(REDIS_UPDATES, shelterId))
	if err != nil {
		return nil, err
	}

	updates, err := RedisGetUpdates(keys)
	if err != nil {
		return nil, err
	}

	return updates.Paginate(pagination), nil
}

func GetUpdate(shelterId, id string) (Update, error) {
	if len(shelterId) == 0 || len(id) == 0 {
		return Update{}, fmt.Errorf("Getting Update failed! UpdateId '%s' is not valid!", id)
	}

	return getUpdate(fmt.Sprintf(REDIS_UPDATE, shelterId, id))
}

func GetLastUpdate(shelterId string) (Update, error) {
	if len(shelterId) == 0 {
		return Update{}, fmt.Errorf("Getting LastUpdate failed! ShelterId '%s' is not valid!", shelterId)
	}

	k, err := RedisGetValue(fmt.Sprintf(REDIS_LAST_UPDATE, shelterId))
	if err != nil {
		return Update{}, err
	}

	return getUpdate(k)
}

func getUpdate(k string) (Update, error) {
	updates, err := RedisGetUpdates(Keys{k})
	if err != nil {
		return Update{}, err
	}

	if len(updates) == 0 {
		return Update{}, fmt.Errorf("Getting Update failed! UpdateId '%s' not found!", k)
	}

	return updates[0], nil
}

func DeleteUpdates(shelterId string, pagination Pagination) error {
	updates, err := GetUpdates(shelterId, pagination)
	if err != nil {
		return err
	}

	for _, u := range updates {
		if err := DeleteUpdate(shelterId, u.Id); err != nil {
			return err
		}
	}

	return nil
}

func DeleteUpdate(shelterId, updateId string) error {
	return RedisDeleteUpdate(shelterId, updateId)
}
