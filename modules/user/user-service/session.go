package service

import (
	"fmt"
	"time"

	"github.com/eolinker/apinto-dashboard/modules/user"
	user_model "github.com/eolinker/apinto-dashboard/modules/user/user-model"

	"github.com/eolinker/apinto-dashboard/cache"

	"github.com/go-redis/redis/v8"
)

func sessionCacheKey(session string) string {
	return fmt.Sprintf("session:%s", session)
}

func newSessionCache(client *redis.ClusterClient) user.ISessionCache {
	return cache.CreateRedisCache[user_model.Session](client, time.Hour*24*7, sessionCacheKey, "apinto", "session")

}
