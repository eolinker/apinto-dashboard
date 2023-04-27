package service

import (
	"github.com/eolinker/apinto-dashboard/cache"
	module_plugin "github.com/eolinker/apinto-dashboard/modules/module-plugin"
	"github.com/eolinker/apinto-dashboard/modules/module-plugin/entry"
	"github.com/go-redis/redis/v8"
)

type navigationModulesCache struct {
	cache.IRedisCache[entry.EnabledModule]
}

func (i *navigationModulesCache) Key() string {
	return "navigation_modules"
}

func newNavigationModulesCache(client *redis.ClusterClient) module_plugin.INavigationModulesCache {
	return &navigationModulesCache{
		IRedisCache: cache.CreateRedisCache[entry.EnabledModule](client),
	}
}
