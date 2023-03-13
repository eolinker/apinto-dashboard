package cache

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/entry"
	"github.com/go-redis/redis/v8"
)

type IUserInfoCache interface {
	IRedisCache[entry.UserInfo]
	Key(userId int) string
}

type userInfoCache struct {
	*redisCache[entry.UserInfo]
}

func (userInfoCache) Key(userId int) string {
	return fmt.Sprintf("user_info:%d", userId)
}

func newUserInfoCache(client *redis.ClusterClient) IUserInfoCache {
	return &userInfoCache{
		redisCache: createRedisCache[entry.UserInfo](client),
	}
}
