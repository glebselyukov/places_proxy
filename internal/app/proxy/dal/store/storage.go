package store

import (
	"context"

	"github.com/prospik/places_proxy/internal/app/proxy/dal/dao"
)

// Storage data access layer
type Storage interface {
	PlacesDAL
	Connect() error
	Disconnect() error
	Ping() error
}

// PlacesDAL layer for places caching
type PlacesDAL interface {
	CheckKey(ctx context.Context, key string) bool
	GetPlaces(ctx context.Context, key string) (*dao.ResultData, error)
	SavePlaces(ctx context.Context, key string, checksum uint64, places []dao.Places)
}
