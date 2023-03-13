package cache

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/model"
	"github.com/go-redis/redis/v8"
)

type IBatchOnlineApiTaskCache interface {
	IRedisCache[model.BatchOnlineCheckTask]
	Key(token string) string
}

type batchOnlineApiTaskCache struct {
	*redisCache[model.BatchOnlineCheckTask]
}

func (i *batchOnlineApiTaskCache) Key(token string) string {
	return fmt.Sprintf("batch_online_api_token:%s", token)
}

func newBatchOnlineTaskCache(client *redis.ClusterClient) IBatchOnlineApiTaskCache {
	cache := &batchOnlineApiTaskCache{
		redisCache: createRedisCache[model.BatchOnlineCheckTask](client),
	}
	return cache

}
