package cache

import (
	"fmt"
	apimodel "github.com/eolinker/apinto-dashboard/modules/api/model"
	"github.com/go-redis/redis/v8"
)

type IImportApiCache interface {
	IRedisCache[apimodel.ImportAPIRedisData]
	Key(token string) string
}

type importApiCache struct {
	*redisCache[apimodel.ImportAPIRedisData]
}

func (i *importApiCache) Key(token string) string {
	return fmt.Sprintf("import_api_token:%s", token)
}

func newImportCache(client *redis.ClusterClient) IImportApiCache {
	cache := &importApiCache{
		redisCache: createRedisCache[apimodel.ImportAPIRedisData](client),
	}
	return cache

}
