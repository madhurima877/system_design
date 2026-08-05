package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	cache *redis.Client
}

func NewRedisCache() *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:16380",
		Password: "",
		DB:       0,
	})

	return &RedisCache{
		cache: client,
	}
}
func (r *RedisCache) Ping() error {
	return r.cache.Ping(context.Background()).Err()
}
func (r *RedisCache) Get(key string) (string, error) {
	val, err := r.cache.Get(context.Background(), key).Result()
	if err != nil {
		return "", err
	}
	return val, nil

}
func (r *RedisCache) Set(key, value string) error {
	return r.cache.Set(context.Background(), key, value, 0).Err()
}
