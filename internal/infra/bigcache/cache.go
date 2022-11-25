package bigcache

import (
	"context"
	"time"

	"github.com/allegro/bigcache/v3"
)

func New() *bigcache.BigCache {
	cache, err := bigcache.New(context.Background(), bigcache.DefaultConfig(1*time.Hour))
	if err != nil {
		panic(err)
	}
	return cache
}
