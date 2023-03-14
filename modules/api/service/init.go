package api_service

import (
	"github.com/eolinker/apinto-dashboard/cache"
	api_store "github.com/eolinker/apinto-dashboard/modules/api/store"
	"github.com/eolinker/apinto-dashboard/store"
	"github.com/eolinker/eosc/common/bean"
	"github.com/go-redis/redis/v8"
)

func init() {

	api := NewAPIService()
	store.RegisterStore(api_store.InitStoreDB)

	cache.RegisterCacheInitHandler(func(client *redis.ClusterClient) {
		importCache := newImportCache(client)
		batchOnlineTaskCache := newBatchOnlineTaskCache(client)
		bean.Injection(&importCache)
		bean.Injection(batchOnlineTaskCache)
	})
	bean.Injection(&api)
}
