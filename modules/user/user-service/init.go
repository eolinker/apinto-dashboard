package service

import (
	"github.com/eolinker/apinto-dashboard/cache"
	"github.com/eolinker/eosc/common/bean"
	"github.com/go-redis/redis/v8"
)

func init() {
	infoService := newUserInfoService()
	bean.Injection(&infoService)

	cache.RegisterCacheInitHandler(func(client *redis.ClusterClient) {
		userInfo := newUserInfoCache(client)
		session := newSessionCache(client)
		bean.Injection(&userInfo)
		bean.Injection(&session)
	})
}
