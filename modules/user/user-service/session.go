package service

import (
	"fmt"

	"github.com/eolinker/apinto-dashboard/modules/user"
	user_model "github.com/eolinker/apinto-dashboard/modules/user/user-model"

	"github.com/eolinker/apinto-dashboard/cache"

	"github.com/go-redis/redis/v8"
)

type sessionCache struct {
	cache.IRedisCache[user_model.Session, string]
}

func sessionCacheKey(session string) string {
	return fmt.Sprintf("session:%s", session)
}

func newSessionCache(client *redis.ClusterClient) user.ISessionCache {
	return sessionCache{
		IRedisCache: cache.CreateRedisCache[user_model.Session](client, sessionCacheKey, "apinto", "session"),
	}
}
