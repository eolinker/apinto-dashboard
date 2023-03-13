package cache

import (
	"fmt"
	apimodel "github.com/eolinker/apinto-dashboard/modules/api/model"
	"github.com/go-redis/redis/v8"
)

type IBatchOnlineApiTaskCache interface {
	IRedisCache[apimodel.BatchOnlineCheckTask]
	Key(token string) string
}

type batchOnlineApiTaskCache struct {
	*redisCache[apimodel.BatchOnlineCheckTask]
}

func (i *batchOnlineApiTaskCache) Key(token string) string {
	return fmt.Sprintf("batch_online_api_token:%s", token)
}

func newBatchOnlineTaskCache(client *redis.ClusterClient) IBatchOnlineApiTaskCache {
	cache := &batchOnlineApiTaskCache{
		redisCache: createRedisCache[apimodel.BatchOnlineCheckTask](client),
	}
	return cache

}
