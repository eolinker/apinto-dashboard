package service

import (
	"github.com/eolinker/apinto-dashboard/cache"
	"github.com/eolinker/eosc/common/bean"
	"github.com/go-redis/redis/v8"
)

func init() {
	cache.RegisterCacheInitHandler(func(client *redis.ClusterClient) {
		iInstalledCache := newNavigationModulesCache(client)
		bean.Injection(&iInstalledCache)
	})
}
