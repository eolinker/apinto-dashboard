package service

import (
	"github.com/eolinker/apinto-dashboard/cache"
	user_model "github.com/eolinker/apinto-dashboard/modules/user/user-model"
	"github.com/go-redis/redis/v8"
)

type IUserInfoCache interface {
	cache.IRedisCache[user_model.UserInfo]
}

type batchOnlineApiTaskCache struct {
	cache.IRedisCache[user_model.UserInfo]
}

func newUserInfoCache(client *redis.ClusterClient) IUserInfoCache {
	return &batchOnlineApiTaskCache{
		IRedisCache: cache.CreateRedisCache[user_model.UserInfo](client, "apinto", "user-center"),
	}
}
