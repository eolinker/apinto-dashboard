package cache

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/model"
	"github.com/go-redis/redis/v8"
)

type IRoleAccessCache interface {
	IRedisCache[model.RoleAccess]
	Key(uuid string) string
}

type roleAccessCache struct {
	*redisCache[model.RoleAccess]
}

func (i *roleAccessCache) Key(uuid string) string {
	return fmt.Sprintf("role_access:%s", uuid)
}

func newUserAccessCache(client *redis.ClusterClient) IRoleAccessCache {
	cache := &roleAccessCache{
		redisCache: createRedisCache[model.RoleAccess](client),
	}
	return cache

}
