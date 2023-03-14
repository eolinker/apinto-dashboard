package api_service

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/cache"
	apimodel "github.com/eolinker/apinto-dashboard/modules/api/model"
	"github.com/go-redis/redis/v8"
)

type IImportApiCache interface {
	cache.IRedisCache[apimodel.ImportAPIRedisData]
	Key(token string) string
}

type importApiCache struct {
	cache.IRedisCache[apimodel.ImportAPIRedisData]
}

func (i *importApiCache) Key(token string) string {
	return fmt.Sprintf("import_api_token:%s", token)
}

func newImportCache(client *redis.ClusterClient) IImportApiCache {
	cache := &importApiCache{
		IRedisCache: cache.CreateRedisCache[apimodel.ImportAPIRedisData](client),
	}
	return cache

}
