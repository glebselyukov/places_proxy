package dal

import (
	"net/url"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"

	redisdb "github.com/prospik/places_proxy/internal/app/proxy/dal/redis"
	"github.com/prospik/places_proxy/internal/app/proxy/dal/store"
)

// ParseDBScheme represents a parsed database type
func ParseDBScheme(uri string) (scheme string, err error) {
	u, err := url.Parse(uri)
	if u != nil {
		scheme = u.Scheme
		return
	}
	err = errors.Wrap(err, "can't parse database url")
	return
}

// New fabric method for storage
func New(uri string) (db store.Storage, err error) {
	scheme, err := ParseDBScheme(uri)
	if err != nil {
		return
	}

	switch scheme {
	case "redis", "rediss":
		var optRedis *redis.Options
		optRedis, err = redis.ParseURL(uri)
		if err != nil {
			err = errors.Wrap(err, "fail parse redis options")
			return
		}
		cli := redis.NewClient(optRedis)
		db = redisdb.NewRedisDB(cli)
	default:
		err = errors.Errorf("not support database with scheme '%s'", scheme)
	}
	return
}
