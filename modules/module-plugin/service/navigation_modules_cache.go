package service

import (
	"github.com/eolinker/apinto-dashboard/cache"
	module_plugin "github.com/eolinker/apinto-dashboard/modules/module-plugin"
	"github.com/eolinker/apinto-dashboard/modules/module-plugin/entry"
	"time"
)

func newNavigationModulesCache() module_plugin.INavigationModulesCache {
	return cache.CreateRedisCacheNoKey[entry.EnabledModule](10*time.Minute, "navigationModulesCacheKey")

}
