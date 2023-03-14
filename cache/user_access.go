package cache

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/model/access-model"
	"github.com/go-redis/redis/v8"
)

type IRoleAccessCache interface {
	IRedisCache[access_model.RoleAccess]
	Key(uuid string) string
}

type roleAccessCache struct {
	*redisCache[access_model.RoleAccess]
}

func (i *roleAccessCache) Key(uuid string) string {
	return fmt.Sprintf("role_access:%s", uuid)
}

func newUserAccessCache(client *redis.ClusterClient) IRoleAccessCache {
	cache := &roleAccessCache{
		redisCache: createRedisCache[access_model.RoleAccess](client),
	}
	return cache

}
