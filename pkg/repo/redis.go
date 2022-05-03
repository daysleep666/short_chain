package repo

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	REDIS_INSTANCE *redis.Client
)

func InitRedis(addr, pwd string) (err error) {
	REDIS_INSTANCE, err = initRedis(addr, pwd)
	return
}

func initRedis(addr, pwd string) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pwd,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return rdb, nil
}
