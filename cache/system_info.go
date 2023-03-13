package cache

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/entry"
	"github.com/go-redis/redis/v8"
)

type ISystemInfoCache interface {
	IRedisCache[entry.SystemInfo]
	Key(key string) string
}

type systemInfoCache struct {
	*redisCache[entry.SystemInfo]
}

func (systemInfoCache) Key(key string) string {
	return fmt.Sprintf("system_info:%s", key)
}

func newSystemInfoCache(client *redis.ClusterClient) ISystemInfoCache {
	return &systemInfoCache{
		redisCache: createRedisCache[entry.SystemInfo](client),
	}
}
