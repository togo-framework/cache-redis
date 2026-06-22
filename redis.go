// Package rediscache is a Redis driver for togo cache (CACHE_DRIVER=redis).
// Install: `togo install togo-framework/cache-redis`.
package rediscache

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/togo-framework/cache"
	"github.com/togo-framework/togo"
)

func init() {
	cache.RegisterDriver("redis", func(k *togo.Kernel) (togo.Cache, error) {
		url := os.Getenv("REDIS_URL")
		if url == "" {
			url = "redis://localhost:6379/0"
		}
		opt, err := redis.ParseURL(url)
		if err != nil {
			return nil, err
		}
		return &store{rdb: redis.NewClient(opt)}, nil
	})
}

type store struct{ rdb *redis.Client }

func (s *store) Get(key string) (any, bool) {
	v, err := s.rdb.Get(context.Background(), key).Result()
	if err != nil {
		return nil, false
	}
	var out any
	if json.Unmarshal([]byte(v), &out) != nil {
		return v, true
	}
	return out, true
}

func (s *store) Set(key string, value any, ttl time.Duration) {
	b, err := json.Marshal(value)
	if err != nil {
		return
	}
	_ = s.rdb.Set(context.Background(), key, b, ttl).Err()
}

func (s *store) Delete(key string) {
	_ = s.rdb.Del(context.Background(), key).Err()
}
