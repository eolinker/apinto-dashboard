package service

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/cache"
	"github.com/eolinker/apinto-dashboard/modules/module-plugin/model"
	"github.com/go-redis/redis/v8"
)

type IInstalledCache interface {
	cache.IRedisCache[model.PluginInstalledStatus]
	Key(pluginID string) string
}

type installedCache struct {
	cache.IRedisCache[model.PluginInstalledStatus]
}

func (i *installedCache) Key(pluginID string) string {
	return fmt.Sprintf("plugin_installed_%s", pluginID)
}

func newIInstalledCache(client *redis.ClusterClient) IInstalledCache {
	return &installedCache{
		IRedisCache: cache.CreateRedisCache[model.PluginInstalledStatus](client),
	}

}
