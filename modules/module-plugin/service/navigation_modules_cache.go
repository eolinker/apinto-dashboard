package service

import (
	"github.com/eolinker/apinto-dashboard/cache"
	module_plugin "github.com/eolinker/apinto-dashboard/modules/module-plugin"
	"github.com/eolinker/apinto-dashboard/modules/module-plugin/entry"
	"github.com/go-redis/redis/v8"
	"time"
)

func newNavigationModulesCache(client *redis.ClusterClient) module_plugin.INavigationModulesCache {
	return cache.CreateRedisCacheNoKey[entry.EnabledModule](client, 10*time.Minute, "navigationModulesCacheKey")

}
