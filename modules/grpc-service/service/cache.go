package service

import (
	"github.com/eolinker/apinto-dashboard/cache"
	grpc_service "github.com/eolinker/apinto-dashboard/grpc-service"
	"time"
)

type INavigationModulesCache interface {
	cache.IRedisCache[grpc_service.NavigationModulesResp, string]
}

func navigationModulesCacheKey(string) string {
	return "navigation_modules"
}

func newNavigationModulesCache() INavigationModulesCache {
	return cache.CreateRedisCache[grpc_service.NavigationModulesResp](time.Hour, navigationModulesCacheKey)
}
