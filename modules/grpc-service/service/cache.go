package service

import (
	"github.com/eolinker/apinto-dashboard/cache"
	grpc_service "github.com/eolinker/apinto-dashboard/grpc-service"
	"github.com/go-redis/redis/v8"
	"time"
)

type INavigationModulesCache interface {
	cache.IRedisCache[grpc_service.NavigationModulesResp, string]
}

func navigationModulesCacheKey(string) string {
	return "navigation_modules"
}

func newNavigationModulesCache(client *redis.ClusterClient) INavigationModulesCache {
	return cache.CreateRedisCache[grpc_service.NavigationModulesResp](client, time.Hour, navigationModulesCacheKey)
}
