package redis

import (
	"github.com/go-redis/redis"
	"github.com/pkg/errors"

	"github.com/prospik/places_proxy/internal/app/proxy/dal/store"
)

type redisDatabase struct {
	client *redis.Client
}

// NewRedisDB constructor for new redis connection
func NewRedisDB(client *redis.Client) store.Storage {
	return &redisDatabase{client: client}
}

var (
	_ store.Storage = (*redisDatabase)(nil)
)

func (db *redisDatabase) Connect() error {
	if db.client == nil {
		return errors.New("can't connect to redis")
	}
	err := db.Ping()
	if err != nil {
		return err
	}
	return nil
}

func (db *redisDatabase) Ping() error {
	if db.client == nil {
		return errors.New("client is not defined")
	}
	if _, err := db.client.Ping().Result(); err != nil {
		return errors.Wrap(err, "client ping result")
	}
	return nil
}

func (db *redisDatabase) Disconnect() error {
	return db.client.Close()
}

func (db *redisDatabase) Some() {}
