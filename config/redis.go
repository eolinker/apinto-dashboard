package config

import (
	"context"
	"time"

	"github.com/eolinker/apinto-dashboard/cache"
	"github.com/redis/go-redis/v9"
)

func InitRedis() {

	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    getRedisAddr(),
		Username: getRedisUserName(),
		Password: getRedisPwd(),
	})

	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second)
	defer cancelFunc()
	if err := client.Ping(timeout).Err(); err != nil {
		_ = client.Close()
		panic(err)
	}

	cache.InitCache(client)
}
