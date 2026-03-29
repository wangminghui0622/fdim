package cache

import (
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/stores/cache"
)

func NewRedisClient(conf cache.CacheConf) *redis.Client {
	addr := "127.0.0.1:6379"
	password := "test625"

	if conf != nil && len(conf) > 0 {
		addr = conf[0].Host
		password = conf[0].Pass
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})
	return rdb
}
