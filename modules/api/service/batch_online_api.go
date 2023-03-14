package api_service

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/cache"
	apimodel "github.com/eolinker/apinto-dashboard/modules/api/model"
	"github.com/go-redis/redis/v8"
)

type IBatchOnlineApiTaskCache interface {
	cache.IRedisCache[apimodel.BatchOnlineCheckTask]
	Key(token string) string
}

type batchOnlineApiTaskCache struct {
	cache.IRedisCache[apimodel.BatchOnlineCheckTask]
}

func (i *batchOnlineApiTaskCache) Key(token string) string {
	return fmt.Sprintf("batch_online_api_token:%s", token)
}

func newBatchOnlineTaskCache(client *redis.ClusterClient) IBatchOnlineApiTaskCache {

	return &batchOnlineApiTaskCache{
		IRedisCache: cache.CreateRedisCache[apimodel.BatchOnlineCheckTask](client),
	}

}
