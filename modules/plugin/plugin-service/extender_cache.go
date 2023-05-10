package plugin_service

import (
	"github.com/eolinker/apinto-dashboard/cache"
	"github.com/eolinker/apinto-dashboard/modules/plugin/plugin-model"
	"github.com/go-redis/redis/v8"
	"time"
)

type IExtenderCache interface {
	cache.IRedisCacheNoKey[plugin_model.ExtenderInfo]
}

func newIExtenderCache(client *redis.ClusterClient) IExtenderCache {
	return cache.CreateRedisCacheNoKey[plugin_model.ExtenderInfo](client, time.Minute*5, "extender")

}
