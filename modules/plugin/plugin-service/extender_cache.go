package plugin_service

import (
	"github.com/eolinker/apinto-dashboard/cache"
	"github.com/eolinker/apinto-dashboard/modules/plugin/plugin-model"
	"time"
)

type IExtenderCache interface {
	cache.IRedisCacheNoKey[*plugin_model.ExtenderInfo]
}

func newIExtenderCache() IExtenderCache {
	return cache.CreateRedisCacheNoKey[*plugin_model.ExtenderInfo](time.Minute*5, "extender")

}
