package service

import (
	"fmt"
	"time"

	"github.com/eolinker/apinto-dashboard/cache"
	user_model "github.com/eolinker/apinto-dashboard/modules/user/user-model"
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
	return fmt.Sprintf("user_info:name:%s", name)
}
func newUserInfoIdCache() IUserInfoCacheId {
	return cache.CreateRedisCache[user_model.UserInfo, int](time.Hour, userCacheKey)

}
func newUserInfoNameCache() IUserInfoCacheName {
	return cache.CreateRedisCache[user_model.UserInfo, string](time.Hour, userCacheName)

}
