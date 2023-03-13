package cache

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/model"
	"github.com/go-redis/redis/v8"
)

type ISessionCache interface {
	IRedisCache[model.Session]
	Key(string) string
}

type sessionCache struct {
	*redisCache[model.Session]
}

func (sessionCache) Key(session string) string {
	return fmt.Sprintf("session:%s", session)
}

func newSessionCache(client *redis.ClusterClient) ISessionCache {
	return sessionCache{
		redisCache: createRedisCache[model.Session](client),
	}

}
