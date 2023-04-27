package service

import (
	"github.com/eolinker/apinto-dashboard/cache"
	grpc_service "github.com/eolinker/apinto-dashboard/grpc-service"
	"github.com/go-redis/redis/v8"
)

type INavigationModulesCache interface {
	cache.IRedisCache[grpc_service.NavigationModulesResp]
	Key() string
}

type navigationModulesCache struct {
	cache.IRedisCache[grpc_service.NavigationModulesResp]
}

func (i *navigationModulesCache) Key() string {
	return "navigation_modules"
}

func newNavigationModulesCache(client *redis.ClusterClient) INavigationModulesCache {
	return &navigationModulesCache{
		IRedisCache: cache.CreateRedisCache[grpc_service.NavigationModulesResp](client),
	}
}
