package namespace_service

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/cache"
	"github.com/eolinker/apinto-dashboard/modules/namespace"
	namespace_model "github.com/eolinker/apinto-dashboard/modules/namespace/namespace-model"
	"github.com/eolinker/eosc/common/bean"
	"github.com/go-redis/redis/v8"
	"time"
)

func init() {
	v := newNamespaceService()
	var i namespace.INamespaceService = v
	bean.Injection(&i)
	cache.RegisterCacheInitHandler(func(client *redis.ClusterClient) {
		v.namespaceCacheByName = cache.CreateRedisCache[namespace_model.Namespace, string](client, time.Minute*30, func(name string) string {
			return fmt.Sprintf("namespace:name:%s", name)
		})
		v.namespaceCacheById = cache.CreateRedisCache[namespace_model.Namespace, int](client, time.Minute*30, func(id int) string {
			return fmt.Sprintf("namespace:id:%d", id)
		})
		v.namespaceCacheAll = cache.CreateRedisCacheNoKey[namespace_model.Namespace](client, time.Minute*30, "namespace:all")
	})
}
