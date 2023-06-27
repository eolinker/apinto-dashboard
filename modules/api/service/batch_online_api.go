package api_service

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/cache"
	apimodel "github.com/eolinker/apinto-dashboard/modules/api/model"
	"time"
)

type IBatchOnlineApiTaskCache interface {
	cache.IRedisCache[apimodel.BatchOnlineCheckTask, string]
}

type batchOnlineApiTaskCache struct {
	cache.IRedisCache[apimodel.BatchOnlineCheckTask, string]
}

func formatKey(token string) string {
	return fmt.Sprintf("batch_online_api_token:%s", token)
}

func newBatchOnlineTaskCache() IBatchOnlineApiTaskCache {

	return &batchOnlineApiTaskCache{
		IRedisCache: cache.CreateRedisCache[apimodel.BatchOnlineCheckTask](time.Hour*8, formatKey),
	}

}
