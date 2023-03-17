package plugin_service

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/cache"
	"github.com/eolinker/apinto-dashboard/modules/plugin/plugin-model"
	"github.com/go-redis/redis/v8"
)

type IExtenderCache interface {
	cache.IRedisCache[plugin_model.ExtenderInfo]
	Key() string
}

type extenderCache struct {
	cache.IRedisCache[plugin_model.ExtenderInfo]
}

func (i *extenderCache) Key() string {
	return fmt.Sprintf("extender")
}

func newIExtenderCache(client *redis.ClusterClient) IExtenderCache {
	return &extenderCache{
		IRedisCache: cache.CreateRedisCache[plugin_model.ExtenderInfo](client),
	}

}
