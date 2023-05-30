package cache

import (
	"sync"

	"github.com/eolinker/eosc/common/bean"
	"github.com/go-redis/redis/v8"
)

type handlerFunc func(client *redis.ClusterClient)

var (
	client *redis.ClusterClient
	lock   sync.Mutex

	handlers []handlerFunc
)

func RegisterCacheInitHandler(h handlerFunc) {
	lock.Lock()
	defer lock.Unlock()
	if client != nil {
		h(client)
	} else {
		handlers = append(handlers, h)
	}

}
func InitCache(c *redis.ClusterClient) {
	//iUserInfoCache := newUserInfoCache(client)
	//iUserAccessCache := newUserAccessCache(client)
	//iSessionCache := newSessionCache(client)
	lock.Lock()
	defer lock.Unlock()
	client = c
	iCommonCache := newCommonCache(client)

	bean.Injection(&iCommonCache)

	for _, h := range handlers {
		h(client)
	}
	handlers = handlers[:0]
}
