package service

import (
	"fmt"

	"github.com/eolinker/apinto-dashboard/cache"
	user_model "github.com/eolinker/apinto-dashboard/modules/user/user-model"
	"github.com/go-redis/redis/v8"
)

type IUserInfoCacheId interface {
	cache.IRedisCache[user_model.UserInfo, int]
}
type IUserInfoCacheName interface {
	cache.IRedisCache[user_model.UserInfo, string]
}

func userCacheKey(userId int) string {
	return fmt.Sprintf("user_info:id:%d", userId)
}
func userCacheName(name string) string {
	return fmt.Sprintf("user_info:name:%d", name)
}
func newUserInfoIdCache(client *redis.ClusterClient) IUserInfoCacheId {
	return cache.CreateRedisCache[user_model.UserInfo, int](client, userCacheKey, "apinto", "user-center")

}
func newUserInfoNameCache(client *redis.ClusterClient) IUserInfoCacheName {
	return cache.CreateRedisCache[user_model.UserInfo, string](client, userCacheName, "apinto", "user-center")

}
