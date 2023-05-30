package api_service

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/cache"
	apimodel "github.com/eolinker/apinto-dashboard/modules/api/model"
	"github.com/go-redis/redis/v8"
	"time"
)

type IImportApiCache interface {
	cache.IRedisCache[apimodel.ImportAPIRedisData, string]
}

func importKey(token string) string {
	return fmt.Sprintf("import_api_token:%s", token)
}

func newImportCache(client *redis.ClusterClient) IImportApiCache {
	cacheInfo := cache.CreateRedisCache[apimodel.ImportAPIRedisData, string](client, time.Hour*8, importKey)

	return cacheInfo

}
