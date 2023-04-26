package service

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/cache"
	user_model "github.com/eolinker/apinto-dashboard/modules/user/user-model"
	"github.com/go-redis/redis/v8"
)

type IUserInfoCache interface {
	cache.IRedisCache[user_model.UserInfo]
	Key(token int) string
}

type userCache struct {
	cache.IRedisCache[user_model.UserInfo]
}

func (userCache) Key(userId int) string {
	return fmt.Sprintf("user_info:%d", userId)
}

func newUserInfoCache(client *redis.ClusterClient) IUserInfoCache {
	return &userCache{
		IRedisCache: cache.CreateRedisCache[user_model.UserInfo](client, "apinto", "user-center"),
	}
}
