package namespace_service

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/cache"
	"github.com/eolinker/apinto-dashboard/modules/namespace"
	namespace_model "github.com/eolinker/apinto-dashboard/modules/namespace/namespace-model"
	"github.com/eolinker/eosc/common/bean"
	"time"
)

func init() {
	v := newNamespaceService()
	var i namespace.INamespaceService = v

	v.namespaceCacheByName = cache.CreateRedisCache[namespace_model.Namespace, string](time.Minute*30, func(name string) string {
		return fmt.Sprintf("namespace:name:%s", name)
	})
	v.namespaceCacheById = cache.CreateRedisCache[namespace_model.Namespace, int](time.Minute*30, func(id int) string {
		return fmt.Sprintf("namespace:id:%d", id)
	})
	v.namespaceCacheAll = cache.CreateRedisCacheNoKey[*namespace_model.Namespace](time.Minute*30, "namespace:all")
	bean.Injection(&i)
}
