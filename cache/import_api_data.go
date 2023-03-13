package cache

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/model"
	"github.com/go-redis/redis/v8"
)

type IImportApiCache interface {
	IRedisCache[model.ImportAPIRedisData]
	Key(token string) string
}

type importApiCache struct {
	*redisCache[model.ImportAPIRedisData]
}

func (i *importApiCache) Key(token string) string {
	return fmt.Sprintf("import_api_token:%s", token)
}

func newImportCache(client *redis.ClusterClient) IImportApiCache {
	cache := &importApiCache{
		redisCache: createRedisCache[model.ImportAPIRedisData](client),
	}
	return cache

}
