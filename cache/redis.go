package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type RedisCache struct {
	redis *redis.Client
}

var ctx = context.Background()

func NewRedisCacheRepository() *RedisCache {
	return &RedisCache{
		redis: redis.NewClient(&redis.Options{
			Addr:     "redis:6379",
			Password: "",
			DB:       0,
		}),
	}
}

func (r *RedisCache) Get(key string) (string, error) {
	val, err := r.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", redis.Nil
	} else if err != nil {
		return "", err
	}
	return val, nil
}

func (r *RedisCache) Put(key string, val string) error {
	logrus.Infof("Saving %v into key %v", val, key)
	return r.redis.Set(ctx, key, val, 0).Err()
}

func (r *RedisCache) Delete(key string) error {
	return r.redis.Del(ctx, key).Err()
}
