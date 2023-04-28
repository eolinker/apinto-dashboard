package service

import (
	"github.com/eolinker/apinto-dashboard/cache"
	grpc_service "github.com/eolinker/apinto-dashboard/grpc-service"
	"github.com/go-redis/redis/v8"
)

type INavigationModulesCache interface {
	cache.IRedisCache[grpc_service.NavigationModulesResp, string]
}

type navigationModulesCache struct {
	cache.IRedisCache[grpc_service.NavigationModulesResp, string]
}

func navigationModulesCacheKey(string) string {
	return "navigation_modules"
}

func newNavigationModulesCache(client *redis.ClusterClient) INavigationModulesCache {
	return &navigationModulesCache{
		IRedisCache: cache.CreateRedisCache[grpc_service.NavigationModulesResp](client, navigationModulesCacheKey),
	}
}
