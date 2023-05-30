package service

import (
	"github.com/eolinker/apinto-dashboard/cache"
	"github.com/eolinker/eosc/common/bean"
	"github.com/go-redis/redis/v8"
)

func init() {
	iModulePluginService := newModulePluginService()
	iModulePlugin := newModulePlugin()
	bean.Injection(&iModulePluginService)
	bean.Injection(&iModulePlugin)

	cache.RegisterCacheInitHandler(func(client *redis.ClusterClient) {
		iInstalledCache := newIInstalledCache(client)
		iNavigationModulesCache := newNavigationModulesCache(client)
		bean.Injection(&iInstalledCache)
		bean.Injection(&iNavigationModulesCache)
	})
}
