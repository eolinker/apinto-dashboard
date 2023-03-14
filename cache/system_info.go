package cache

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/entry/system-entry"
	"github.com/go-redis/redis/v8"
)

type ISystemInfoCache interface {
	IRedisCache[system_entry.SystemInfo]
	Key(key string) string
}

type systemInfoCache struct {
	*redisCache[system_entry.SystemInfo]
}

func (systemInfoCache) Key(key string) string {
	return fmt.Sprintf("system_info:%s", key)
}

func newSystemInfoCache(client *redis.ClusterClient) ISystemInfoCache {
	return &systemInfoCache{
		redisCache: createRedisCache[system_entry.SystemInfo](client),
	}
}
