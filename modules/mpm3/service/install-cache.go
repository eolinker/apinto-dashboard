package service

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/cache"
	"github.com/eolinker/apinto-dashboard/modules/mpm3/model"
	"time"
)

type IInstalledCache cache.IRedisCache[model.PluginInstalledStatus, string]

func installedCacheKey(pluginID string) string {
	return fmt.Sprintf("plugin_installed_%s", pluginID)
}

func newIInstalledCache() IInstalledCache {
	return cache.CreateRedisCache[model.PluginInstalledStatus, string](time.Hour, installedCacheKey)
}
