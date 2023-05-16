package dao

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisCache struct {
	Rdb *redis.Client
}

var Rc *RedisCache

//func init() {
//	rdb := redis.NewClient(config.C.ReadRedisConfig())
//	Rc = &RedisCache{
//		rdb: rdb,
//	}
//}

func (rc *RedisCache) Put(ctx context.Context, key, value string, expire time.Duration) error {
	err := rc.Rdb.Set(ctx, key, value, expire).Err()
	return err
}

func (rc *RedisCache) Get(ctx context.Context, key string) (string, error) {
	result, err := rc.Rdb.Get(ctx, key).Result()
	return result, err
}
