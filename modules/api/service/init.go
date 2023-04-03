package api_service

import (
	"github.com/eolinker/apinto-dashboard/cache"
	"github.com/eolinker/apinto-dashboard/modules/api/driver"
	"github.com/eolinker/eosc/common/bean"

	api_store "github.com/eolinker/apinto-dashboard/modules/api/store"

	"github.com/eolinker/apinto-dashboard/store"
	"github.com/go-redis/redis/v8"
)

func init() {

	api := NewAPIService()
	store.RegisterStore(api_store.InitStoreDB)
	apiDriverManager := newAPIDriverManager()
	apiHttp := driver.CreateAPIHttp("http")
	apiDriverManager.RegisterDriver(DriverApiHTTP, apiHttp)
	bean.Injection(&apiDriverManager)
	cache.RegisterCacheInitHandler(func(client *redis.ClusterClient) {
		importCache := newImportCache(client)
		batchOnlineTaskCache := newBatchOnlineTaskCache(client)
		bean.Injection(&importCache)
		bean.Injection(&batchOnlineTaskCache)
	})
	bean.Injection(&api)
}
