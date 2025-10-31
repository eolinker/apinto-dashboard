package config

import (
	"context"
	"strings"
	"time"

	"github.com/eolinker/apinto-dashboard/cache"
	"github.com/redis/go-redis/v9"
)

func InitRedis() {
	var client redis.UniversalClient
	switch strings.ToLower(systemConfig.RedisConfig.Cluster) {
	case "no", "false":
		client = redis.NewClient(&redis.Options{
			Addr:     systemConfig.RedisConfig.Addr[0],
			Username: systemConfig.RedisConfig.UserName,
			Password: systemConfig.RedisConfig.Password,
		})
	default:
		client = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    systemConfig.RedisConfig.Addr,
			Username: systemConfig.RedisConfig.UserName,
			Password: systemConfig.RedisConfig.Password,
		})
	}

	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second)
	defer cancelFunc()
	if err := client.Ping(timeout).Err(); err != nil {
		_ = client.Close()
		panic(err)
	}

	cache.InitCache(client, systemConfig.RedisConfig.Prefix)
}
