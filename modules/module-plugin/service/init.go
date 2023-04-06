package service

import (
	"github.com/eolinker/apinto-dashboard/cache"
	"github.com/eolinker/eosc/common/bean"
	"github.com/go-redis/redis/v8"
)

func init() {
	modulePluginServiceInfo := newModulePluginService()
	bean.Injection(&modulePluginServiceInfo)

	cache.RegisterCacheInitHandler(func(client *redis.ClusterClient) {
		iExtenderCache := newIExtenderCache(client)
		bean.Injection(&iExtenderCache)
	})
}
