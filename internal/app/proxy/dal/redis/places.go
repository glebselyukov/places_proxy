package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis"

	"github.com/prospik/places_proxy/internal/app/proxy/dal/dao"
	"github.com/prospik/places_proxy/internal/app/proxy/dal/errors"
)

func (db *redisDatabase) CheckKey(ctx context.Context, key string) bool {
	cli := db.client.WithContext(ctx)

	value, err := cli.Keys(key).Result()
	if err != nil {
		return true
	}

	return len(value) > 0 && value[0] == key
}

func (db *redisDatabase) GetPlaces(ctx context.Context, key string) (*dao.ResultData, error) {
	cli := db.client.WithContext(ctx)

	result := &dao.ResultData{}

	places := &dao.CachedData{}
	if err := cli.HGet(key, placesField).Scan(places); err != nil {
		return result, err
	}

	var checksum uint64
	if err := cli.HGet(key, checksumField).Scan(&checksum); err != nil {
		return result, err
	}

	times, err := cli.HGet(key, timeField).Result()
	if err != nil {
		return result, err
	}

	bytes, err := json.Marshal(places.Data)
	if err != nil {
		return result, err
	}

	result.Data = bytes
	result.Checksum = checksum
	result.Time = times

	return result, nil
}

func (db *redisDatabase) SavePlaces(ctx context.Context, key string, checksum uint64, places []dao.Places) {
	cli := db.client.WithContext(ctx)
	if key == "" {
		return
	}

	_ = cli.Watch(func(tx *redis.Tx) (err error) {
		err = tx.HMSet(key, cachedData(checksum, places)).Err()
		if err != nil {
			return errors.ErrPlacesInternal
		}
		return
	}, key)
}

func cachedData(checksum uint64, places []dao.Places) map[string]interface{} {

	rfc := time.Now().In(time.UTC).Format(time.RFC3339Nano)

	data := &dao.CachedData{
		Data: places,
	}
	return map[string]interface{}{
		placesField:   data,
		checksumField: checksum,
		timeField:     rfc,
	}
}
