package cluster_service

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/cache"
	"github.com/eolinker/apinto-dashboard/modules/cluster/cluster-model"
	"time"
)

type INodeCache interface {
	cache.IRedisCache[[]*cluster_model.Node, int]
}

func installedCacheKey(clusterId int) string {
	return fmt.Sprintf("cluster_%d_nodes", clusterId)
}

func newINodeCache() INodeCache {
	return cache.CreateRedisCache[[]*cluster_model.Node, int](5*time.Minute, installedCacheKey)
}
