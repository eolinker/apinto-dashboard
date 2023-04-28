package plugin_service

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/cache"
	"github.com/eolinker/apinto-dashboard/modules/plugin/plugin-model"
	"github.com/go-redis/redis/v8"
)

type IExtenderCache interface {
	cache.IRedisCacheNoKey[plugin_model.ExtenderInfo]
}

func extenderCacheKey() string {
	return fmt.Sprintf("extender")
}

func newIExtenderCache(client *redis.ClusterClient) IExtenderCache {
	return cache.CreateRedisCacheNoKey[plugin_model.ExtenderInfo](client, extenderCacheKey)

}
