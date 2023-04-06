package service

import (
	"github.com/eolinker/apinto-dashboard/cache"
	"github.com/eolinker/eosc/common/bean"
	"github.com/go-redis/redis/v8"
)

func init() {
	pluginServiceInfo := newPluginService()
	clusterPluginServiceInfo := newClusterPluginService()
	bean.Injection(&pluginServiceInfo)
	bean.Injection(&clusterPluginServiceInfo)

	cache.RegisterCacheInitHandler(func(client *redis.ClusterClient) {
		iExtenderCache := newIExtenderCache(client)
		bean.Injection(&iExtenderCache)
	})
}
