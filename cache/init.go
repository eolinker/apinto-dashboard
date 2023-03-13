package cache

import (
	"github.com/eolinker/eosc/common/bean"
	"github.com/go-redis/redis/v8"
)

func InitCache(client *redis.ClusterClient) {
	iUserInfoCache := newUserInfoCache(client)
	iUserAccessCache := newUserAccessCache(client)
	iSessionCache := newSessionCache(client)
	iSystemInfoCache := newSystemInfoCache(client)
	iImportApiCache := newImportCache(client)
	iBatchApiCache := newBatchOnlineTaskCache(client)
	iCommonCache := newCommonCache(client)

	bean.Injection(&iUserInfoCache)
	bean.Injection(&iUserAccessCache)
	bean.Injection(&iSessionCache)
	bean.Injection(&iSystemInfoCache)
	bean.Injection(&iImportApiCache)
	bean.Injection(&iBatchApiCache)
	bean.Injection(&iCommonCache)
}
